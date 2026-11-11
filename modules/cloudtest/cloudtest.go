// Package cloudtest ties a cloud resource's lifetime to a test, and turns a
// partial read into a failure. It holds no cloud SDK and names no resource
// type, because a shared inventory library owns those.
package cloudtest

import (
	"context"
	"fmt"
	"time"

	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	inv "github.com/james00012/cloud-inventory"
)

// TestingTB is TestingT plus Cleanup, which it predates. Create needs Cleanup,
// so it asks for this and the compiler enforces it rather than the resource
// leaking when the method is absent. *testing.T satisfies it.
type TestingTB interface {
	testing.TestingT
	Cleanup(func())
}

// Managed is a resource type a test can both make and remove, since a test
// that creates one must be able to clean it up.
type Managed[T any, Spec any, S any] interface {
	inv.Creator[T, Spec, S]
	inv.Deleter
}

// CleanupTimeout bounds one resource's deletion. A test whose resource is slow
// to delete raises it.
var CleanupTimeout = 15 * time.Minute

// Create makes the resource, ties its lifetime to the test, and returns the
// provider's own struct. Register Terraform's teardown with Cleanup rather than
// defer, so the two unwind on one stack and this resource goes first: every
// deferred call in a test body runs before any cleanup.
func Create[T inv.Resource, Spec any, S any](t TestingTB, ctx context.Context, c Managed[T, Spec, S], scope S, spec Spec) T {
	t.Helper()

	created, err := c.Create(ctx, scope, spec)
	if err != nil {
		t.Fatalf("cloudtest: creating the resource for %s: %v", t.Name(), err)
		return created
	}

	id := created.ResourceID()
	kind := created.ResourceType()

	t.Cleanup(func() {
		// A cleanup runs after the test's own context is done, so it needs a
		// deadline of its own or every cleanup after a timeout leaks.
		cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), CleanupTimeout)
		defer cancel()

		if err := c.Delete(cleanupCtx, id); err != nil {
			t.Errorf("cloudtest: deleting %s %s: %v", kind, id, err)
		}
	})

	return created
}

// Get reads one resource and fails the test if it is not there. The id comes
// from the provider's own constructor, so no caller writes the format by hand.
func Get[T any](t testing.TestingT, ctx context.Context, g inv.Getter[T], id inv.ID) T {
	t.Helper()

	got, err := g.Get(ctx, id)
	if err != nil {
		t.Fatalf("cloudtest: reading %s: %v", id, err)
	}

	return got
}

// List reads every resource in the scope, and fails the test if any part of the
// scope could not be read. A check that reads three of five regions and finds
// nothing has proved nothing, but it looks exactly like a pass.
func List[T any, S any](t testing.TestingT, ctx context.Context, l inv.Lister[T, S], scope S) []T {
	t.Helper()

	found, incomplete, err := l.List(ctx, scope)
	if err != nil {
		t.Fatalf("cloudtest: listing: %v", err)
	}
	if incomplete.Any() {
		t.Fatalf("cloudtest: could not read %s, so this listing proves nothing: %v",
			incomplete.String(), incomplete.Cause())
	}

	return found
}

// WaitUntil reads the resource until it satisfies want, and fails the test if it
// never does. The retrying is this repository's own, so the logging and the
// error text match every other wait in a suite.
func WaitUntil[T any](t testing.TestingT, ctx context.Context, g inv.Getter[T], id inv.ID, description string, want func(T) bool, maxRetries int, sleepBetweenRetries time.Duration) T {
	t.Helper()

	var zero T

	answer, err := retry.DoWithRetryInterfaceContextE(t, ctx,
		fmt.Sprintf("waiting for %s to be %s", id, description),
		maxRetries, sleepBetweenRetries,
		func() (any, error) {
			got, err := g.Get(ctx, id)
			if err != nil {
				return nil, err
			}
			if !want(got) {
				return nil, fmt.Errorf("%s is not %s yet", id, description)
			}
			return got, nil
		})
	if err != nil {
		t.Fatalf("cloudtest: %s was not %s: %v", id, description, err)
		return zero
	}

	typed, ok := answer.(T)
	if !ok {
		t.Fatalf("cloudtest: waiting for %s answered with %T", id, answer)
		return zero
	}

	return typed
}

// AssertNone fails the test if anything in the scope matches, which is what a
// leak check asserts. A scope it could not read fails too, through List.
func AssertNone[T inv.Resource, S any](t testing.TestingT, ctx context.Context, l inv.Lister[T, S], scope S, matches func(T) bool) {
	t.Helper()

	for _, found := range List(t, ctx, l, scope) {
		if matches(found) {
			t.Errorf("cloudtest: %s %s was left behind", found.ResourceType(), found.ResourceID())
		}
	}
}
