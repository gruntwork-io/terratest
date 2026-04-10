//go:build kubeall || kubernetes
// +build kubeall kubernetes

// NOTE: we have build tags to differentiate kubernetes tests from non-kubernetes tests. This is done because minikube
// is heavy and can interfere with docker related tests in terratest. Specifically, many of the tests start to fail with
// `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes tests and helm
// tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.  We
// recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package k8s_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/k8s"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test that RunKubectlAndGetOutputE will run kubectl and return the output by running a can-i command call.
func TestRunKubectlAndGetOutputReturnsOutput(t *testing.T) {
	t.Parallel()

	namespaceName := "kubectl-test-" + strings.ToLower(random.UniqueID())
	options := k8s.NewKubectlOptions("", "", namespaceName)
	output, err := k8s.RunKubectlAndGetOutputE(t, options, "auth", "can-i", "get", "pods")
	require.NoError(t, err)
	require.Equal(t, "yes", output)
}

//nolint:paralleltest,tparallel // subtests share mutable parsedTimeout via http server
func TestKubectlRequestTimeout(t *testing.T) {
	t.Parallel()

	var parsedTimeout time.Duration

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parsedTimeout, _ = time.ParseDuration(r.URL.Query().Get("timeout"))
		select {
		case <-time.After(3 * time.Second):
		case <-r.Context().Done():
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("dummy-error"))
	}))

	config := fmt.Sprintf(`
apiVersion: v1
kind: Config
clusters:
- name: dummy-cluster
  cluster:
    server: %s
users:
- name: dummy-user
  user:
    token: dummy-token
contexts:
- name: dummy-context
  context:
    cluster: dummy-cluster
    user: dummy-user
current-context: dummy-context
`, server.URL)

	t.Run("WithoutTimeout", func(t *testing.T) {
		options := &k8s.KubectlOptions{
			ContextName: "dummy-context",
			ConfigPath:  k8s.StoreConfigToTempFile(t, config),
		}
		_, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "pods")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "dummy-error")
		assert.NotContains(t, err.Error(), "Client.Timeout exceeded while awaiting headers")
	})

	t.Run("WithTimeout", func(t *testing.T) {
		options := &k8s.KubectlOptions{
			ContextName:    "dummy-context",
			ConfigPath:     k8s.StoreConfigToTempFile(t, config),
			RequestTimeout: time.Second,
		}
		_, err := k8s.RunKubectlAndGetOutputE(t, options, "get", "pods")
		require.Error(t, err)
		assert.Equal(t, options.RequestTimeout, parsedTimeout)
		assert.NotContains(t, err.Error(), "dummy-error")
		assert.Contains(t, err.Error(), "Client.Timeout exceeded while awaiting headers")
	})
}
