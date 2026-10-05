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
	"google.golang.org/api/networkconnectivity/v1"
	"google.golang.org/api/option"
)

// newFakeNetworkConnectivityService points a real *networkconnectivity.Service at a local test
// server, so the Google transport is exercised rather than a hand-written stand-in for a type we
// do not own.
func newFakeNetworkConnectivityService(t *testing.T, handler http.Handler) *networkconnectivity.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := networkconnectivity.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetNetworkConnectivityHubAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking hub module sets, because the point of
	// reading settings back is asserting a module configured the hub it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/hubs/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/global/hubs/gw-library-test",
			"description":"created by terratest",
			"state":"ACTIVE",
			"policyMode":"PRESET",
			"presetTopology":"MESH",
			"exportPsc":true,
			"labels":{"managed-by":"terratest"}
		}`))
	})

	hub, err := gcp.GetNetworkConnectivityHubAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", hub.Description)
	assert.Equal(t, "ACTIVE", hub.State)
	assert.Equal(t, "MESH", hub.PresetTopology)
	assert.True(t, hub.ExportPsc)
	assert.Equal(t, "terratest", hub.Labels["managed-by"])
}

func TestGetNetworkConnectivityHubAttrsWithClientMissingHub(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the hub and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetNetworkConnectivityHubAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetNetworkConnectivitySpokeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking spoke module sets, because the point
	// of reading settings back is asserting a module configured the spoke it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/spokes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/spokes/gw-library-test","description":"created by terratest","hub":"projects/gw-library-test-project/locations/global/hubs/gw-library-test","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	spoke, err := gcp.GetNetworkConnectivitySpokeAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", spoke.Description)
	assert.True(t, strings.HasSuffix(spoke.Hub, "/hubs/gw-library-test"))
	assert.Equal(t, "ACTIVE", spoke.State)
	assert.Equal(t, "terratest", spoke.Labels["purpose"])
}

func TestGetPolicyBasedRouteAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking policy based route module sets.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/policyBasedRoutes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/policyBasedRoutes/gw-library-test","description":"created by terratest","priority":900,"network":"projects/gw-library-test-project/global/networks/gw-library-test","filter":{"protocolVersion":"IPV4","destRange":"192.0.2.0/24","ipProtocol":"TCP"},"nextHopOtherRoutes":"DEFAULT_ROUTING"}`))
	})

	route, err := gcp.GetPolicyBasedRouteAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", route.Description)
	assert.Equal(t, int64(900), route.Priority)
	require.NotNil(t, route.Filter)
	assert.Equal(t, "192.0.2.0/24", route.Filter.DestRange)
	assert.Equal(t, "TCP", route.Filter.IpProtocol)
}

func TestGetInternalRangeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a range the terraform-google-networking internal range module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/internalRanges/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/internalRanges/gw-library-test","description":"created by terratest","ipCidrRange":"10.20.0.0/24","usage":"FOR_VPC","peering":"FOR_SELF","labels":{"purpose":"terratest"}}`))
	})

	internalRange, err := gcp.GetInternalRangeAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "10.20.0.0/24", internalRange.IpCidrRange)
	assert.Equal(t, "FOR_VPC", internalRange.Usage)
	assert.Equal(t, "FOR_SELF", internalRange.Peering)
}

func TestGetNetworkConnectivityMulticloudDataTransferConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a multicloud data transfer config the terraform-google-networking multicloud data transfer config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/multicloudDataTransferConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/multicloudDataTransferConfigs/gw-library-test","description":"created by terratest","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkConnectivityMulticloudDataTransferConfigAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetNetworkConnectivityMulticloudDataTransferConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a multicloud data transfer config that is not there should read a sentence about that multicloud data transfer config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkConnectivityMulticloudDataTransferConfigAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkConnectivityRegionalEndpointAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a regional endpoint the terraform-google-networking regional endpoint module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/regionalEndpoints/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/regionalEndpoints/gw-library-test","description":"created by terratest","accessType":"GLOBAL","targetGoogleApi":"storage.googleapis.com","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkConnectivityRegionalEndpointAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "GLOBAL", attrs.AccessType)
	assert.Equal(t, "storage.googleapis.com", attrs.TargetGoogleApi)
}

func TestGetNetworkConnectivityRegionalEndpointAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a regional endpoint that is not there should read a sentence about that regional endpoint, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkConnectivityRegionalEndpointAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkConnectivityServiceConnectionPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a service connection policy the terraform-google-networking service connection policy module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/serviceConnectionPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/serviceConnectionPolicies/gw-library-test","description":"created by terratest","network":"projects/gw-library-test-project/global/networks/gw-library-test","serviceClass":"gcp-memorystore","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkConnectivityServiceConnectionPolicyAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "gcp-memorystore", attrs.ServiceClass)
}

func TestGetNetworkConnectivityServiceConnectionPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a service connection policy that is not there should read a sentence about that service connection policy, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkConnectivityServiceConnectionPolicyAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkConnectivityTransportAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a transport the terraform-google-networking transport module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/transports/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/transports/gw-library-test","description":"created by terratest","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkConnectivityTransportAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetNetworkConnectivityTransportAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a transport that is not there should read a sentence about that transport, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkConnectivityTransportAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkConnectivityHubGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a hub group the terraform-google-networking hub group module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/hubs/gw-library-parent/groups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/hubs/gw-library-parent/groups/gw-library-test","description":"created by terratest","state":"ACTIVE","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkConnectivityHubGroupAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetNetworkConnectivityHubGroupAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a hub group that is not there should read a sentence about that hub group, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkConnectivityHubGroupAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetNetworkConnectivityTransferDestinationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a transfer destination the terraform-google-networking transfer destination module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/multicloudDataTransferConfigs/gw-library-parent/destinations/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/multicloudDataTransferConfigs/gw-library-parent/destinations/gw-library-test","description":"created by terratest","ipPrefix":"198.51.100.0/24","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetNetworkConnectivityTransferDestinationAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "198.51.100.0/24", attrs.IpPrefix)
}

func TestGetNetworkConnectivityTransferDestinationAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a transfer destination that is not there should read a sentence about that transfer destination, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetNetworkConnectivityTransferDestinationAttrsWithClient(context.Background(), newFakeNetworkConnectivityService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
