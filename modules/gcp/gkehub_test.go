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

func TestGetGKEHubFeatureAttrsWithClient(t *testing.T) {
	t.Parallel()

	// One read answers for two modules, so the response carries both the feature's own spec and the per membership
	// spec that the feature membership module writes.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/features/configmanagement"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/features/configmanagement","labels":{"owner":"terratest"},"membershipSpecs":{"projects/gw-library-test-project/locations/global/memberships/gw-library-test":{"configmanagement":{"version":"1.19.0"}}},"resourceState":{"state":"ACTIVE"}}`))
	})

	feature, err := gcp.GetGKEHubFeatureAttrsWithClient(context.Background(), newFakeGKEHubService(t, handler), "gw-library-test-project", "global", "configmanagement")
	require.NoError(t, err)

	assert.Equal(t, map[string]string{"owner": "terratest"}, feature.Labels)
	require.Len(t, feature.MembershipSpecs, 1)
	require.NotNil(t, feature.ResourceState)
	assert.Equal(t, "ACTIVE", feature.ResourceState.State)
}

func TestGetGKEHubFeatureAttrsWithClientReportsAMissingFeature(t *testing.T) {
	t.Parallel()

	// A caller who names a feature that is not enabled should read a sentence about that feature, not a status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKEHubFeatureAttrsWithClient(context.Background(), newFakeGKEHubService(t, handler), "gw-library-test-project", "global", "multiclusteringress")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "multiclusteringress")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetGKEHubFleetAttrsWithClient(t *testing.T) {
	t.Parallel()

	// A project holds one fleet per location and the API names it default, so the path is fixed and the caller passes
	// no id.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/fleets/default"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/fleets/default","displayName":"terratest fleet","state":{"code":"READY"},"labels":{"owner":"terratest"}}`))
	})

	fleet, err := gcp.GetGKEHubFleetAttrsWithClient(context.Background(), newFakeGKEHubService(t, handler), "gw-library-test-project", "global")
	require.NoError(t, err)

	assert.Equal(t, "terratest fleet", fleet.DisplayName)
	require.NotNil(t, fleet.State)
	assert.Equal(t, "READY", fleet.State.Code)
}

func TestGetGKEHubFleetAttrsWithClientReportsAProjectWithNoFleet(t *testing.T) {
	t.Parallel()

	// A project that has never created a fleet answers 404, and the sentence has to say the project holds no fleet
	// rather than imply the project is gone.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetGKEHubFleetAttrsWithClient(context.Background(), newFakeGKEHubService(t, handler), "gw-library-test-project", "global")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-test-project")
	assert.Contains(t, err.Error(), "has no fleet")
}
