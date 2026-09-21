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

func TestFetchSubnetworkWithClient(t *testing.T) {
	t.Parallel()

	// Each response below is shaped like the one Google returns for a resource the
	// terraform-google-networking modules created, not a copy of any one fixture's values.
	handler := respond(t, http.MethodGet, "/projects/gw-library-test-project/regions/us-central1/subnetworks/gw-library-test", http.StatusOK, `{
		"kind":"compute#subnetwork",
		"name":"gw-library-test",
		"ipCidrRange":"10.10.0.0/24",
		"privateIpGoogleAccess":true,
		"secondaryIpRanges":[{"rangeName":"pods","ipCidrRange":"10.20.0.0/20"}],
		"logConfig":{"enable":true,"flowSampling":0.5}
	}`)

	subnetwork, err := gcp.FetchSubnetworkWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "10.10.0.0/24", subnetwork.IpCidrRange)
	assert.True(t, subnetwork.PrivateIpGoogleAccess)
	require.Len(t, subnetwork.SecondaryIpRanges, 1)
	assert.Equal(t, "pods", subnetwork.SecondaryIpRanges[0].RangeName)
	require.NotNil(t, subnetwork.LogConfig)
	assert.InDelta(t, 0.5, subnetwork.LogConfig.FlowSampling, 1e-9)
}

func TestFetchSubnetworkWithClientMissingSubnetwork(t *testing.T) {
	t.Parallel()

	// The error names the call, the project, the region and the subnet, in the shape the other
	// compute reads here use.
	_, err := gcp.FetchSubnetworkWithClient(context.Background(), newFakeComputeService(t, respond(t, "", "", http.StatusNotFound, "")), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "Subnetworks.Get(gw-library-test-project, us-central1, gone)")
}

func TestFetchRouterWithClient(t *testing.T) {
	t.Parallel()

	// A NAT gateway comes back inside its router, so the router read is also the NAT read.
	handler := respond(t, http.MethodGet, "/projects/gw-library-test-project/regions/us-central1/routers/gw-library-test", http.StatusOK, `{
		"kind":"compute#router",
		"name":"gw-library-test",
		"bgp":{"asn":64514},
		"nats":[{
			"name":"gw-library-test-nat",
			"natIpAllocateOption":"AUTO_ONLY",
			"sourceSubnetworkIpRangesToNat":"ALL_SUBNETWORKS_ALL_IP_RANGES",
			"logConfig":{"enable":true,"filter":"ERRORS_ONLY"}
		}]
	}`)

	router, err := gcp.FetchRouterWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.NotNil(t, router.Bgp)
	assert.Equal(t, int64(64514), router.Bgp.Asn)
	require.Len(t, router.Nats, 1)
	assert.Equal(t, "gw-library-test-nat", router.Nats[0].Name)
	assert.Equal(t, "AUTO_ONLY", router.Nats[0].NatIpAllocateOption)
	assert.Equal(t, "ALL_SUBNETWORKS_ALL_IP_RANGES", router.Nats[0].SourceSubnetworkIpRangesToNat)
	require.NotNil(t, router.Nats[0].LogConfig)
	assert.Equal(t, "ERRORS_ONLY", router.Nats[0].LogConfig.Filter)
}

func TestFetchRouterWithClientMissingRouter(t *testing.T) {
	t.Parallel()

	_, err := gcp.FetchRouterWithClient(context.Background(), newFakeComputeService(t, respond(t, "", "", http.StatusNotFound, "")), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "Routers.Get(gw-library-test-project, us-central1, gone)")
}

func TestFetchAddressWithClient(t *testing.T) {
	t.Parallel()

	handler := respond(t, http.MethodGet, "/projects/gw-library-test-project/regions/us-central1/addresses/gw-library-test", http.StatusOK, `{
		"kind":"compute#address",
		"name":"gw-library-test",
		"description":"created by terratest",
		"address":"10.10.0.5",
		"addressType":"INTERNAL",
		"purpose":"GCE_ENDPOINT",
		"status":"RESERVED"
	}`)

	address, err := gcp.FetchAddressWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", address.Description)
	assert.Equal(t, "10.10.0.5", address.Address)
	assert.Equal(t, "INTERNAL", address.AddressType)
	assert.Equal(t, "GCE_ENDPOINT", address.Purpose)
	assert.Equal(t, "RESERVED", address.Status)
}

func TestFetchAddressWithClientMissingAddress(t *testing.T) {
	t.Parallel()

	_, err := gcp.FetchAddressWithClient(context.Background(), newFakeComputeService(t, respond(t, "", "", http.StatusNotFound, "")), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "Addresses.Get(gw-library-test-project, us-central1, gone)")
}

func TestFetchGlobalAddressWithClient(t *testing.T) {
	t.Parallel()

	// A global address is read from its own endpoint, with no region in the path.
	handler := respond(t, http.MethodGet, "/projects/gw-library-test-project/global/addresses/gw-library-test", http.StatusOK, `{
		"kind":"compute#address",
		"name":"gw-library-test",
		"addressType":"INTERNAL",
		"purpose":"VPC_PEERING",
		"prefixLength":16,
		"network":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/networks/gw-library-test"
	}`)

	address, err := gcp.FetchGlobalAddressWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "INTERNAL", address.AddressType)
	assert.Equal(t, "VPC_PEERING", address.Purpose)
	assert.Equal(t, int64(16), address.PrefixLength)
	assert.True(t, strings.HasSuffix(address.Network, "/global/networks/gw-library-test"))
}

func TestFetchGlobalAddressWithClientMissingAddress(t *testing.T) {
	t.Parallel()

	_, err := gcp.FetchGlobalAddressWithClient(context.Background(), newFakeComputeService(t, respond(t, "", "", http.StatusNotFound, "")), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "GlobalAddresses.Get(gw-library-test-project, gone)")
}
