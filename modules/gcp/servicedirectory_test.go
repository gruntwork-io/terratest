package gcp_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	"google.golang.org/api/servicedirectory/v1"
)

// newFakeServiceDirectoryService points a real Service Directory client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeServiceDirectoryService(t *testing.T, handler http.Handler) *servicedirectory.APIService {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := servicedirectory.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetServiceDirectoryNamespaceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a namespace the terraform-google-networking namespace module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/namespaces/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/namespaces/gw-library-test","labels":{"purpose":"terratest"}}`))
	})

	namespace, err := gcp.GetServiceDirectoryNamespaceAttrsWithClient(context.Background(), newFakeServiceDirectoryService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest", namespace.Labels["purpose"])
}

func TestGetServiceDirectoryServiceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a service the terraform-google-networking service module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/namespaces/gw-library-test/services/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/namespaces/gw-library-test/services/gw-library-test","annotations":{"purpose":"terratest"}}`))
	})

	registered, err := gcp.GetServiceDirectoryServiceAttrsWithClient(context.Background(), newFakeServiceDirectoryService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest", registered.Annotations["purpose"])
}

func TestGetServiceDirectoryEndpointAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an endpoint the terraform-google-networking endpoint module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/namespaces/gw-library-test/services/gw-library-test/endpoints/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/namespaces/gw-library-test/services/gw-library-test/endpoints/gw-library-test","address":"10.0.0.7","port":8080,"annotations":{"purpose":"terratest"}}`))
	})

	endpoint, err := gcp.GetServiceDirectoryEndpointAttrsWithClient(context.Background(), newFakeServiceDirectoryService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "10.0.0.7", endpoint.Address)
	assert.Equal(t, int64(8080), endpoint.Port)
}

func TestGetServiceDirectoryNamespaceIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a namespace the terraform-google-networking namespace IAM policy module granted access
	// on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/namespaces/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)

		// A conditional binding only comes back at version 3, so the read has to ask for it, which for
		// this call means in the request body rather than in the query.
		var request struct {
			Options struct {
				RequestedPolicyVersion int64 `json:"requestedPolicyVersion"`
			} `json:"options"`
		}

		// assert rather than require: a failed require inside a handler stops the wrong goroutine.
		if assert.NoError(t, json.NewDecoder(r.Body).Decode(&request)) {
			assert.Equal(t, int64(3), request.Options.RequestedPolicyVersion)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/servicedirectory.viewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetServiceDirectoryNamespaceIamPolicyAttrsWithClient(context.Background(), newFakeServiceDirectoryService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/servicedirectory.viewer", policy.Bindings[0].Role)
	require.NotNil(t, policy.Bindings[0].Condition, "a conditional binding should keep its condition")
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, policy.Bindings[0].Condition.Expression)
}

func TestGetServiceDirectoryServiceIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a service the terraform-google-networking service IAM policy module granted access
	// on, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/namespaces/gw-library-test/services/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)

		// A conditional binding only comes back at version 3, so the read has to ask for it, which for
		// this call means in the request body rather than in the query.
		var request struct {
			Options struct {
				RequestedPolicyVersion int64 `json:"requestedPolicyVersion"`
			} `json:"options"`
		}

		// assert rather than require: a failed require inside a handler stops the wrong goroutine.
		if assert.NoError(t, json.NewDecoder(r.Body).Decode(&request)) {
			assert.Equal(t, int64(3), request.Options.RequestedPolicyVersion)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/servicedirectory.editor","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetServiceDirectoryServiceIamPolicyAttrsWithClient(context.Background(), newFakeServiceDirectoryService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/servicedirectory.editor", policy.Bindings[0].Role)
	require.NotNil(t, policy.Bindings[0].Condition, "a conditional binding should keep its condition")
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, policy.Bindings[0].Condition.Expression)
}
