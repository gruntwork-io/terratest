package k8s_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s/v2"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
)

func TestSaveAndLoadKubectlOptions(t *testing.T) {
	t.Parallel()

	tmpFolder := t.TempDir()

	expectedData := &k8s.KubectlOptions{
		ContextName: "terratest-context",
		ConfigPath:  "~/.kube/config",
		Namespace:   "default",
		Env: map[string]string{
			"TERRATEST_ENV_VAR": "terratest",
		},
	}
	k8s.SaveKubectlOptions(t, tmpFolder, expectedData)

	actualData := k8s.LoadKubectlOptions(t, tmpFolder)
	assert.Equal(t, expectedData, actualData)
}

// fatalRecorder captures a Fatalf instead of failing the enclosing test.
//
// FailNow must not return. The failure originates inside teststate.Save, which calls t.Fatalf on a marshal error
// and then carries on to os.WriteFile; a recorder that returned would let it write an empty file. Goexit gives the
// same semantics as testing.T, so the recorder has to run on its own goroutine.
type fatalRecorder struct {
	failed bool
	msg    string
}

func (r *fatalRecorder) Fail()                 { r.failed = true }
func (r *fatalRecorder) FailNow()              { r.failed = true; runtime.Goexit() }
func (r *fatalRecorder) Error(args ...any)     { r.failed = true }
func (r *fatalRecorder) Errorf(string, ...any) { r.failed = true }
func (r *fatalRecorder) Fatal(args ...any)     { r.msg = fmt.Sprint(args...); r.FailNow() }
func (r *fatalRecorder) Fatalf(format string, args ...any) {
	r.msg = fmt.Sprintf(format, args...)
	r.FailNow()
}
func (r *fatalRecorder) Name() string { return "fatalRecorder" }
func (r *fatalRecorder) Helper()      {}

// runUntilExit runs fn on its own goroutine and waits for it, so a FailNow inside fn does not end the caller.
func runUntilExit(fn func()) {
	done := make(chan struct{})

	go func() { defer close(done); fn() }()

	<-done
}

// TestSaveKubectlOptionsRejectsRestConfig pins that saving options carrying a RestConfig fails and writes nothing,
// rather than dropping the config and producing a file whose reload targets the ambient kubeconfig cluster.
func TestSaveKubectlOptionsRejectsRestConfig(t *testing.T) {
	t.Parallel()

	folder := t.TempDir()
	options := k8s.NewKubectlOptionsWithRestConfig(&rest.Config{Host: "https://example.com"}, "default")

	recorder := &fatalRecorder{}
	runUntilExit(func() { k8s.SaveKubectlOptions(recorder, folder, options) })

	assert.True(t, recorder.failed, "saving options with a RestConfig must fail the test")
	assert.Contains(t, recorder.msg, "RestConfig", "the message must name the offending field")
	assert.NoFileExists(t, filepath.Join(folder, ".test-data", "KubectlOptions.json"),
		"no file may be written when the save is rejected")
}
