package gcp_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchHAVPNGatewayWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking HA VPN gateway module sets, because the point
	// of reading settings back is asserting a module configured the gateway it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/vpnGateways/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","stackType":"IPV4_ONLY","network":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/networks/gw-library-test","vpnInterfaces":[{"id":0},{"id":1}]}`))
	})

	gateway, err := gcp.FetchHAVPNGatewayWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", gateway.Name)
	assert.Equal(t, "IPV4_ONLY", gateway.StackType)
	require.Len(t, gateway.VpnInterfaces, 2, "an HA gateway carries two interfaces")
	assert.True(t, strings.HasSuffix(gateway.Network, "/networks/gw-library-test"))
}

func TestFetchTargetVPNGatewayWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking classic VPN gateway module sets, because the point
	// of reading settings back is asserting a module configured the gateway it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/targetVpnGateways/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","status":"READY","network":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/networks/gw-library-test"}`))
	})

	gateway, err := gcp.FetchTargetVPNGatewayWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", gateway.Name)
	assert.Equal(t, "created by terratest", gateway.Description)
	assert.Equal(t, "READY", gateway.Status)
}

func TestFetchExternalVPNGatewayWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking external VPN gateway module sets, because the point
	// of reading settings back is asserting a module configured the gateway it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/externalVpnGateways/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","redundancyType":"SINGLE_IP_INTERNALLY_REDUNDANT","interfaces":[{"id":0,"ipAddress":"192.0.2.1"}]}`))
	})

	gateway, err := gcp.FetchExternalVPNGatewayWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", gateway.Name)
	assert.Equal(t, "SINGLE_IP_INTERNALLY_REDUNDANT", gateway.RedundancyType)
	require.Len(t, gateway.Interfaces, 1)
	assert.Equal(t, "192.0.2.1", gateway.Interfaces[0].IpAddress)
}

func TestFetchVPNTunnelWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking VPN tunnel module sets, because the point
	// of reading settings back is asserting a module configured the tunnel it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/vpnTunnels/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","ikeVersion":2,"status":"FIRST_HANDSHAKE","peerExternalGatewayInterface":0,"sharedSecret":"*************","router":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/regions/us-central1/routers/gw-library-test"}`))
	})

	tunnel, err := gcp.FetchVPNTunnelWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", tunnel.Name)
	assert.Equal(t, int64(2), tunnel.IkeVersion)
	assert.True(t, strings.HasSuffix(tunnel.Router, "/routers/gw-library-test"))
	// Google answers with a mask rather than the secret it was given.
	assert.Equal(t, "*************", tunnel.SharedSecret, "a shared secret should never come back in clear")
}

func TestFetchRouterRoutePolicyWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking router route policy module sets, because the point
	// of reading settings back is asserting a module configured the policy it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/routers/gw-library-test/getRoutePolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"resource":{"name":"gw-library-test","type":"ROUTE_POLICY_TYPE_EXPORT","terms":[{"priority":1,"match":{"expression":"true"}}]}}`))
	})

	policy, err := gcp.FetchRouterRoutePolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", policy.Name)
	assert.Equal(t, "ROUTE_POLICY_TYPE_EXPORT", policy.Type)
	require.Len(t, policy.Terms, 1)
	assert.Equal(t, int64(1), policy.Terms[0].Priority)
}
