package cloudtest

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	inv "github.com/james00012/cloud-inventory"
)

// *testing.T satisfies the interface Create asks for, so the Cleanup
// requirement costs a caller nothing.
var _ TestingTB = (*testing.T)(nil)

// res is a stand-in for a provider's resource type.
type res struct {
	id   inv.ID
	name string
	up   bool
}

func (r res) ResourceID() inv.ID      { return r.id }
func (r res) ResourceType() string    { return "fake-thing" }
func (r res) DisplayName() string     { return r.name }
func (r res) CreatedAt() time.Time    { return time.Time{} }
func (r res) Tags() map[string]string { return nil }

// fake implements every call interface the wrapper uses.
type fake struct {
	createErr error
	created   []inv.ID
	deleted   []inv.ID
	getSeq    []res
	list      []res
	incomp    inv.Incomplete
	listErr   error

	// recorded at call time: inspecting the context afterwards sees the
	// cleanup's own defer cancel() instead of what Delete was handed.
	deleteCtxErr      error
	deleteHadDeadline bool
}

func (f *fake) Create(ctx context.Context, scope string, spec string) (res, error) {
	if f.createErr != nil {
		return res{}, f.createErr
	}
	id := inv.ID(spec)
	f.created = append(f.created, id)
	return res{id: id, name: spec}, nil
}

func (f *fake) Delete(ctx context.Context, id inv.ID) error {
	f.deleteCtxErr = ctx.Err()
	_, f.deleteHadDeadline = ctx.Deadline()
	f.deleted = append(f.deleted, id)
	return nil
}

func (f *fake) Get(ctx context.Context, id inv.ID) (res, error) {
	if len(f.getSeq) == 0 {
		return res{}, errors.New("not found")
	}
	next := f.getSeq[0]
	f.getSeq = f.getSeq[1:]
	return next, nil
}

func (f *fake) List(ctx context.Context, scope string) ([]res, inv.Incomplete, error) {
	return f.list, f.incomp, f.listErr
}

// recorder stands in for *testing.T so a test can assert on what the wrapper
// reported. It satisfies Terratest's TestingT, plus Cleanup.
type recorder struct {
	errorf   []string
	fatalf   []string
	cleanups []func()
	failed   bool
}

func (r *recorder) Helper()      {}
func (r *recorder) Name() string { return "TestRecorded" }
func (r *recorder) Fail()        { r.failed = true }
func (r *recorder) FailNow()     { r.failed = true }
func (r *recorder) Fatal(args ...any) {
	r.fatalf = append(r.fatalf, "fatal")
	r.failed = true
}
func (r *recorder) Fatalf(format string, args ...any) {
	r.fatalf = append(r.fatalf, format)
	r.failed = true
}
func (r *recorder) Error(args ...any) {
	r.errorf = append(r.errorf, "error")
	r.failed = true
}
func (r *recorder) Errorf(format string, args ...any) {
	r.errorf = append(r.errorf, format)
	r.failed = true
}
func (r *recorder) Cleanup(fn func()) { r.cleanups = append(r.cleanups, fn) }

// runCleanups unwinds the stack the way the testing package does.
func (r *recorder) runCleanups() {
	for i := len(r.cleanups) - 1; i >= 0; i-- {
		r.cleanups[i]()
	}
}

func TestCreateTiesTheResourceToTheTest(t *testing.T) {
	rec := &recorder{}
	f := &fake{}

	got := Create[res, string, string](rec, context.Background(), f, "scope", "thing-1")

	if got.ResourceID() != "thing-1" {
		t.Fatalf("got id %q, want thing-1", got.ResourceID())
	}
	if len(f.deleted) != 0 {
		t.Fatalf("deleted before the test finished: %v", f.deleted)
	}
	if len(rec.cleanups) != 1 {
		t.Fatalf("registered %d cleanups, want 1", len(rec.cleanups))
	}

	rec.runCleanups()

	if len(f.deleted) != 1 || f.deleted[0] != "thing-1" {
		t.Fatalf("cleanup deleted %v, want [thing-1]", f.deleted)
	}
	if rec.failed {
		t.Fatalf("reported a failure: %v %v", rec.errorf, rec.fatalf)
	}
}

func TestCleanupGetsItsOwnDeadlineAfterTheTestContextIsDone(t *testing.T) {
	rec := &recorder{}
	f := &fake{}

	ctx, cancel := context.WithCancel(context.Background())
	Create[res, string, string](rec, ctx, f, "scope", "thing-1")

	// The test's context dies, as it does on a timeout or a failure.
	cancel()
	rec.runCleanups()

	if len(f.deleted) != 1 {
		t.Fatalf("cleanup did not delete after the test context was cancelled: %v", f.deleted)
	}
	if f.deleteCtxErr != nil {
		t.Fatalf("cleanup ran with a dead context: %v", f.deleteCtxErr)
	}
	if !f.deleteHadDeadline {
		t.Fatalf("cleanup context carried no deadline of its own")
	}
}

// TestCleanupOrderAgainstTerraformTeardown pins the ordering the Create doc
// comment tells callers to use. Registering Terraform's teardown with t.Cleanup
// puts both on one stack, so the resource goes before the network it sits in.
func TestCleanupOrderAgainstTerraformTeardown(t *testing.T) {
	rec := &recorder{}
	f := &fake{}
	var order []string

	// What a test body does first: apply, and register the teardown.
	rec.Cleanup(func() { order = append(order, "terraform destroy") })

	// Then create a resource inside what Terraform made.
	Create[res, string, string](rec, context.Background(), f, "scope", "thing-1")
	rec.Cleanup(func() { order = append(order, "instance deleted") })

	rec.runCleanups()

	if order[0] != "instance deleted" || order[len(order)-1] != "terraform destroy" {
		t.Fatalf("teardown ran in the wrong order: %v", order)
	}
}

func TestListFailsTheTestWhenPartOfTheScopeCouldNotBeRead(t *testing.T) {
	rec := &recorder{}
	f := &fake{
		list: []res{{id: "found-1"}},
		incomp: inv.Incomplete{
			Asked: 3,
			Gaps:  []inv.Gap{{Scope: "eu-west-1", Fault: inv.FaultDenied, Err: errors.New("denied")}},
		},
	}

	List[res, string](rec, context.Background(), f, "scope")

	if len(rec.fatalf) != 1 || !strings.Contains(rec.fatalf[0], "proves nothing") {
		t.Fatalf("a partial read did not fail the test: %v %v", rec.fatalf, rec.errorf)
	}
}

func TestAssertNoneReportsALeak(t *testing.T) {
	rec := &recorder{}
	f := &fake{list: []res{{id: "leaked-1", name: "e2e-abc"}, {id: "other", name: "prod"}}}

	AssertNone[res, string](rec, context.Background(), f, "scope", func(r res) bool {
		return strings.HasPrefix(r.DisplayName(), "e2e-")
	})

	if len(rec.errorf) != 1 || !strings.Contains(rec.errorf[0], "left behind") {
		t.Fatalf("did not report the leak: %v", rec.errorf)
	}
}

func TestWaitUntilRetriesUntilTheResourceIsReady(t *testing.T) {
	rec := &recorder{}
	f := &fake{getSeq: []res{{id: "x", up: false}, {id: "x", up: false}, {id: "x", up: true}}}

	got := WaitUntil[res](rec, context.Background(), f, "x", "up",
		func(r res) bool { return r.up }, 5, time.Millisecond)

	if !got.up {
		t.Fatalf("returned a resource that was not up")
	}
	if rec.failed {
		t.Fatalf("failed the test while retrying: %v %v", rec.errorf, rec.fatalf)
	}
}

func TestWaitUntilFailsWhenTheResourceNeverBecomesReady(t *testing.T) {
	rec := &recorder{}
	f := &fake{getSeq: []res{{id: "x"}, {id: "x"}, {id: "x"}}}

	WaitUntil[res](rec, context.Background(), f, "x", "up",
		func(r res) bool { return r.up }, 3, time.Millisecond)

	if len(rec.fatalf) == 0 {
		t.Fatalf("did not fail after exhausting retries")
	}
}

func TestGetReadsOneResource(t *testing.T) {
	rec := &recorder{}
	f := &fake{getSeq: []res{{id: "thing-1", name: "api"}}}

	got := Get[res](rec, context.Background(), f, "thing-1")

	if got.DisplayName() != "api" {
		t.Fatalf("got %q, want api", got.DisplayName())
	}
	if rec.failed {
		t.Fatalf("failed a successful read: %v %v", rec.errorf, rec.fatalf)
	}
}

func TestGetFailsTheTestWhenTheResourceIsNotThere(t *testing.T) {
	rec := &recorder{}
	f := &fake{}

	Get[res](rec, context.Background(), f, "missing-1")

	if len(rec.fatalf) != 1 || !strings.Contains(rec.fatalf[0], "reading") {
		t.Fatalf("a missing resource did not fail the test: %v", rec.fatalf)
	}
}

func TestCreateFailsTheTestWhenTheProviderRefuses(t *testing.T) {
	rec := &recorder{}
	f := &fake{createErr: errors.New("quota exceeded")}

	Create[res, string, string](rec, context.Background(), f, "scope", "thing-1")

	if len(rec.fatalf) != 1 || !strings.Contains(rec.fatalf[0], "creating") {
		t.Fatalf("a refused create did not fail the test: %v", rec.fatalf)
	}
	if len(rec.cleanups) != 0 {
		t.Fatalf("registered cleanup for a resource that was never made")
	}
}

func TestListFailsTheTestWhenTheReadItselfFails(t *testing.T) {
	rec := &recorder{}
	f := &fake{listErr: errors.New("no credentials")}

	List[res, string](rec, context.Background(), f, "scope")

	if len(rec.fatalf) != 1 || !strings.Contains(rec.fatalf[0], "listing") {
		t.Fatalf("a failed read did not fail the test: %v", rec.fatalf)
	}
}
