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
	"google.golang.org/api/networkservices/v1"
	"google.golang.org/api/option"
)

// newFakeNetworkServicesService points a real Network Services client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeNetworkServicesService(t *testing.T, handler http.Handler) *networkservices.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := networkservices.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetNetworkServicesMeshAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a mesh the terraform-google-networking mesh module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/meshes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/meshes/gw-library-test","description":"created by terratest","interceptionPort":15001,"labels":{"purpose":"terratest"}}`))
	})

	mesh, err := gcp.GetNetworkServicesMeshAttrsWithClient(context.Background(), newFakeNetworkServicesService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(15001), mesh.InterceptionPort)
	assert.Equal(t, "created by terratest", mesh.Description)
}

func TestGetNetworkServicesGRPCRouteAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a route the terraform-google-networking gRPC route module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/grpcRoutes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/grpcRoutes/gw-library-test","hostnames":["terratest.example.com"],"rules":[{"matches":[{"method":{"grpcService":"terratest.Service","grpcMethod":"Call"}}],"action":{"retryPolicy":{"numRetries":3,"retryConditions":["unavailable"]}}}]}`))
	})

	route, err := gcp.GetNetworkServicesGRPCRouteAttrsWithClient(context.Background(), newFakeNetworkServicesService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, []string{"terratest.example.com"}, route.Hostnames)
	require.Len(t, route.Rules, 1)
	require.NotNil(t, route.Rules[0].Action)
	require.NotNil(t, route.Rules[0].Action.RetryPolicy)
	assert.Equal(t, int64(3), route.Rules[0].Action.RetryPolicy.NumRetries)
}

func TestGetNetworkServicesHTTPRouteAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a route the terraform-google-networking HTTP route module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/httpRoutes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/httpRoutes/gw-library-test","hostnames":["terratest.example.com"],"rules":[{"matches":[{"prefixMatch":"/terratest"}],"action":{"redirect":{"httpsRedirect":true,"responseCode":"MOVED_PERMANENTLY_DEFAULT"}}}]}`))
	})

	route, err := gcp.GetNetworkServicesHTTPRouteAttrsWithClient(context.Background(), newFakeNetworkServicesService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, route.Rules, 1)
	require.Len(t, route.Rules[0].Matches, 1)
	require.NotNil(t, route.Rules[0].Action)
	require.NotNil(t, route.Rules[0].Action.Redirect)
	assert.Equal(t, "/terratest", route.Rules[0].Matches[0].PrefixMatch)
	assert.True(t, route.Rules[0].Action.Redirect.HttpsRedirect)
}

func TestGetNetworkServicesTCPRouteAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a route the terraform-google-networking TCP route module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/tcpRoutes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/tcpRoutes/gw-library-test","rules":[{"matches":[{"address":"0.0.0.0/0","port":"8080"}],"action":{"originalDestination":true}}]}`))
	})

	route, err := gcp.GetNetworkServicesTCPRouteAttrsWithClient(context.Background(), newFakeNetworkServicesService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, route.Rules, 1)
	require.Len(t, route.Rules[0].Matches, 1)
	require.NotNil(t, route.Rules[0].Action)
	assert.Equal(t, "8080", route.Rules[0].Matches[0].Port)
	assert.True(t, route.Rules[0].Action.OriginalDestination)
}

func TestGetNetworkServicesTLSRouteAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a route the terraform-google-networking TLS route module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/tlsRoutes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/tlsRoutes/gw-library-test","rules":[{"matches":[{"sniHost":["terratest.example.com"],"alpn":["h2"]}],"action":{"destinations":[{"serviceName":"projects/gw-library-test-project/global/backendServices/gw-library-test","weight":1}]}}]}`))
	})

	route, err := gcp.GetNetworkServicesTLSRouteAttrsWithClient(context.Background(), newFakeNetworkServicesService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	require.Len(t, route.Rules, 1)
	require.Len(t, route.Rules[0].Matches, 1)
	require.NotNil(t, route.Rules[0].Action)
	require.Len(t, route.Rules[0].Action.Destinations, 1)
	assert.Equal(t, []string{"terratest.example.com"}, route.Rules[0].Matches[0].SniHost)
	assert.Equal(t, int64(1), route.Rules[0].Action.Destinations[0].Weight)
}

func TestGetNetworkServicesEndpointPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a policy the terraform-google-networking endpoint policy module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/endpointPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/endpointPolicies/gw-library-test","type":"SIDECAR_PROXY","endpointMatcher":{"metadataLabelMatcher":{"metadataLabelMatchCriteria":"MATCH_ANY","metadataLabels":[{"labelName":"purpose","labelValue":"terratest"}]}}}`))
	})

	policy, err := gcp.GetNetworkServicesEndpointPolicyAttrsWithClient(context.Background(), newFakeNetworkServicesService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "SIDECAR_PROXY", policy.Type)
	require.NotNil(t, policy.EndpointMatcher)
	require.NotNil(t, policy.EndpointMatcher.MetadataLabelMatcher)
	assert.Equal(t, "MATCH_ANY", policy.EndpointMatcher.MetadataLabelMatcher.MetadataLabelMatchCriteria)
}

func TestGetNetworkServicesMeshAttrsWithClientReportsAMissingMesh(t *testing.T) {
	t.Parallel()

	// A caller who asks for a mesh that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetNetworkServicesMeshAttrsWithClient(context.Background(), newFakeNetworkServicesService(t, handler), "gw-library-test-project", "global", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
