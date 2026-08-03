package k8s_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/k8s/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/rest"
)

// fatalRecorder captures a Fatalf instead of failing the enclosing test. FailNow semantics require that the call
// does not return, so Fatalf ends the goroutine the way testing.T would.
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

// runAndRecover runs fn on its own goroutine so a FailNow inside it does not end the calling test.
func runAndRecover(fn func()) {
	done := make(chan struct{})

	go func() {
		defer close(done)
		fn()
	}()

	<-done
}

// TestSaveKubectlOptionsRejectsRestConfig pins that saving options built with a RestConfig fails loudly rather
// than writing a file whose reload would silently target the ambient kubeconfig cluster.
func TestSaveKubectlOptionsRejectsRestConfig(t *testing.T) {
	t.Parallel()

	folder := t.TempDir()
	options := k8s.NewKubectlOptionsWithRestConfig(&rest.Config{Host: "https://example.com"}, "default")

	recorder := &fatalRecorder{}
	runAndRecover(func() { k8s.SaveKubectlOptions(recorder, folder, options) })

	assert.True(t, recorder.failed, "saving options with a RestConfig must fail the test")
	assert.Contains(t, recorder.msg, "RestConfig", "the message must name the offending field")
	assert.NoFileExists(t, filepath.Join(folder, ".test-data", "KubectlOptions.json"),
		"no file may be written when the save is rejected")
}

// TestSaveKubectlOptionsAcceptsKubeconfigOptions is the companion: the ordinary path is untouched.
func TestSaveKubectlOptionsAcceptsKubeconfigOptions(t *testing.T) {
	t.Parallel()

	folder := t.TempDir()
	options := k8s.NewKubectlOptions("terratest-context", "~/.kube/config", "default")

	k8s.SaveKubectlOptions(t, folder, options)
	require.FileExists(t, filepath.Join(folder, ".test-data", "KubectlOptions.json"))

	loaded := k8s.LoadKubectlOptions(t, folder)
	assert.Equal(t, "terratest-context", loaded.ContextName)
	assert.Equal(t, "default", loaded.Namespace)
	assert.Nil(t, loaded.RestConfig)
}

// TestSaveKubectlOptionsMessageIsActionable keeps the failure message useful rather than just a marshal error.
func TestSaveKubectlOptionsMessageIsActionable(t *testing.T) {
	t.Parallel()

	recorder := &fatalRecorder{}
	options := k8s.NewKubectlOptionsWithRestConfig(&rest.Config{Host: "https://example.com"}, "default")
	runAndRecover(func() { k8s.SaveKubectlOptions(recorder, t.TempDir(), options) })

	assert.NotContains(t, strings.ToLower(recorder.msg), "unsupported type",
		"the caller should not be shown a raw encoding/json error")
	assert.Contains(t, recorder.msg, "kubeconfig", "the message should point at the supported alternative")
}
