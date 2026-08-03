package k8s_test

import (
	"fmt"
	"path/filepath"
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

// fatalRecorder captures a Fatalf instead of failing the enclosing test. SaveKubectlOptions returns after its
// Fatalf, so FailNow does not need Goexit semantics here.
type fatalRecorder struct {
	failed bool
	msg    string
}

func (r *fatalRecorder) Fail()                 { r.failed = true }
func (r *fatalRecorder) FailNow()              { r.failed = true }
func (r *fatalRecorder) Error(args ...any)     { r.failed = true }
func (r *fatalRecorder) Errorf(string, ...any) { r.failed = true }
func (r *fatalRecorder) Fatal(args ...any)     { r.msg = fmt.Sprint(args...); r.FailNow() }
func (r *fatalRecorder) Fatalf(format string, args ...any) {
	r.msg = fmt.Sprintf(format, args...)
	r.FailNow()
}
func (r *fatalRecorder) Name() string { return "fatalRecorder" }
func (r *fatalRecorder) Helper()      {}

// TestSaveKubectlOptionsRejectsRestConfig pins that saving options built with a RestConfig fails loudly and writes
// nothing. Before the json:"-" tag this also failed, but with a raw encoding/json error; the tag would otherwise
// have turned it into a silent, lossy write whose reload targets the ambient kubeconfig cluster.
func TestSaveKubectlOptionsRejectsRestConfig(t *testing.T) {
	t.Parallel()

	folder := t.TempDir()
	options := k8s.NewKubectlOptionsWithRestConfig(&rest.Config{Host: "https://example.com"}, "default")

	recorder := &fatalRecorder{}
	k8s.SaveKubectlOptions(recorder, folder, options)

	assert.True(t, recorder.failed, "saving options with a RestConfig must fail the test")
	assert.Contains(t, recorder.msg, "RestConfig", "the message must name the offending field")
	assert.Contains(t, recorder.msg, "kubeconfig", "the message must point at the supported alternative")
	assert.NoFileExists(t, filepath.Join(folder, ".test-data", "KubectlOptions.json"),
		"no file may be written when the save is rejected")
}
