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
	"google.golang.org/api/gkehub/v1"
	"google.golang.org/api/option"
)

// newFakeGKEHubService points a real GKE Hub client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeGKEHubService(t *testing.T, handler http.Handler) *gkehub.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := gkehub.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetGKEHubScopeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// A fleet spans regions, so its scopes live in the global collection rather than in one region.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/scopes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/scopes/gw-library-test","labels":{"purpose":"terratest"},"namespaceLabels":{"team":"terratest"},"state":{"code":"READY"}}`))
	})

	scope, err := gcp.GetGKEHubScopeAttrsWithClient(context.Background(), newFakeGKEHubService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, map[string]string{"purpose": "terratest"}, scope.Labels)
	assert.Equal(t, map[string]string{"team": "terratest"}, scope.NamespaceLabels)
}

func TestGetGKEHubNamespaceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// A namespace belongs to a scope, so the path names both.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/scopes/gw-library-scope/namespaces/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/scopes/gw-library-scope/namespaces/gw-library-test","scope":"projects/gw-library-test-project/locations/global/scopes/gw-library-scope","namespaceLabels":{"team":"terratest"}}`))
	})

	namespace, err := gcp.GetGKEHubNamespaceAttrsWithClient(context.Background(), newFakeGKEHubService(t, handler), "gw-library-test-project", "gw-library-scope", "gw-library-test")
	require.NoError(t, err)

	assert.Contains(t, namespace.Scope, "scopes/gw-library-scope")
	assert.Equal(t, map[string]string{"team": "terratest"}, namespace.NamespaceLabels)
}

func TestGetGKEHubScopeAttrsWithClientReportsAMissingScope(t *testing.T) {
	t.Parallel()

	// A caller who asks for a scope that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetGKEHubScopeAttrsWithClient(context.Background(), newFakeGKEHubService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetGKEHubMembershipAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a fleet membership a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/memberships/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"description":"terratest membership","externalId":"8f0bcb5a-6d54-4b5f-9a1e-0d0d2a2b3c4d","endpoint":{"gkeCluster":{"resourceLink":"//container.googleapis.com/projects/gw-library-test-project/locations/us-central1/clusters/gw-library-test"}},"state":{"code":"READY"}}`))
	})

	result, err := gcp.GetGKEHubMembershipAttrsWithClient(context.Background(), newFakeGKEHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest membership", result.Description)
	require.NotNil(t, result.Endpoint)
	require.NotNil(t, result.Endpoint.GkeCluster)
	assert.Contains(t, result.Endpoint.GkeCluster.ResourceLink, "/clusters/gw-library-test")
	require.NotNil(t, result.State)
	assert.Equal(t, "READY", result.State.Code)
}

func TestGetGKEHubMembershipBindingAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a fleet membership binding a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/memberships/gw-library-test/bindings/gw-library-binding"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"scope":"projects/gw-library-test-project/locations/global/scopes/gw-library-test","state":{"code":"READY"},"labels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.GetGKEHubMembershipBindingAttrsWithClient(context.Background(), newFakeGKEHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-binding")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/locations/global/scopes/gw-library-test", result.Scope)
	require.NotNil(t, result.State)
	assert.Equal(t, "READY", result.State.Code)
}

func TestGetGKEHubMembershipAttrsWithClientReportsAMissingFleetMembership(t *testing.T) {
	t.Parallel()

	// A caller who names a fleet membership that is not there should read a sentence about that fleet membership, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKEHubMembershipAttrsWithClient(context.Background(), newFakeGKEHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
