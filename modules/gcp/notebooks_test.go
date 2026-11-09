package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/notebooks/v1"
	"google.golang.org/api/option"
)

// newFakeNotebooksService points a real Notebooks client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeNotebooksService(t *testing.T, handler http.Handler) *notebooks.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := notebooks.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetNotebooksEnvironmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an environment the terraform-google-ml
	// notebooks environment module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/environments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/environments/gw-library-test","displayName":"terratest environment","description":"created by terratest","vmImage":{"project":"deeplearning-platform-release","imageFamily":"common-cpu"}}`))
	})

	environment, err := gcp.GetNotebooksEnvironmentAttrsWithClient(context.Background(), newFakeNotebooksService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest environment", environment.DisplayName)
	assert.Equal(t, "created by terratest", environment.Description)
	require.NotNil(t, environment.VmImage)
	assert.Equal(t, "common-cpu", environment.VmImage.ImageFamily)
}

func TestGetNotebooksEnvironmentAttrsWithClientReportsAMissingEnvironment(t *testing.T) {
	t.Parallel()

	// A caller who asks for an environment that is not there should be told that, rather than be
	// handed the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetNotebooksEnvironmentAttrsWithClient(context.Background(), newFakeNotebooksService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNotebooksRuntimeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a runtime the managed notebooks runtime
	// module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/runtimes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/runtimes/gw-library-test","state":"ACTIVE","healthState":"HEALTHY","accessConfig":{"accessType":"SINGLE_USER","runtimeOwner":"terratest@example.com","proxyUri":"gw-library-test-dot-us-central1.notebooks.googleusercontent.com"},"virtualMachine":{"virtualMachineConfig":{"machineType":"n1-standard-4"}}}`))
	})

	runtime, err := gcp.GetNotebooksRuntimeAttrsWithClient(context.Background(), newFakeNotebooksService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", runtime.State)
	require.NotNil(t, runtime.AccessConfig)
	assert.Equal(t, "SINGLE_USER", runtime.AccessConfig.AccessType)
	assert.Contains(t, runtime.AccessConfig.ProxyUri, "notebooks.googleusercontent.com")
}

func TestGetNotebooksRuntimeAttrsWithClientReportsAMissingRuntime(t *testing.T) {
	t.Parallel()

	// A caller who asks for a runtime that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetNotebooksRuntimeAttrsWithClient(context.Background(), newFakeNotebooksService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNotebooksRuntimeIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The policy is shaped like the one Google returns after the runtime IAM policy module has applied
	// a binding, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/runtimes/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		assert.Equal(t, "3", r.URL.Query().Get("options.requestedPolicyVersion"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":1,"etag":"BwYabc","bindings":[{"role":"roles/notebooks.viewer","members":["serviceAccount:terratest@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.GetNotebooksRuntimeIamPolicyAttrsWithClient(context.Background(), newFakeNotebooksService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/notebooks.viewer", policy.Bindings[0].Role)
}

func TestGetNotebooksRuntimeIamPolicyAttrsWithClientReportsAMissingRuntime(t *testing.T) {
	t.Parallel()

	// Asking for the policy of a runtime that is not there should name the runtime rather than hand
	// back the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetNotebooksRuntimeIamPolicyAttrsWithClient(context.Background(), newFakeNotebooksService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
