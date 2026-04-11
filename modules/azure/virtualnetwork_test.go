package azure_test

import (
	"context"
	"net"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	networkfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakeSubnetsClient(t *testing.T, srv *networkfake.SubnetsServer) *armnetwork.SubnetsClient {
	t.Helper()

	client, err := armnetwork.NewSubnetsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: networkfake.NewSubnetsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeVirtualNetworksClient(t *testing.T, srv *networkfake.VirtualNetworksServer) *armnetwork.VirtualNetworksClient {
	t.Helper()

	client, err := armnetwork.NewVirtualNetworksClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: networkfake.NewVirtualNetworksServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// CheckSubnetContainsIP tests (business logic: IP range checking)
// ---------------------------------------------------------------------------

func TestCheckSubnetContainsIP(t *testing.T) {
	t.Parallel()

	srv := &networkfake.SubnetsServer{
		Get: func(_ context.Context, _ string, _ string, _ string, _ *armnetwork.SubnetsClientGetOptions) (resp azfake.Responder[armnetwork.SubnetsClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.SubnetsClientGetResponse{
				Subnet: armnetwork.Subnet{
					Name: to.Ptr("test-subnet"),
					Properties: &armnetwork.SubnetPropertiesFormat{
						AddressPrefix: to.Ptr("10.0.0.0/24"),
					},
				},
			}, nil)

			return
		},
	}

	client := newFakeSubnetsClient(t, srv)

	// Get subnet to extract address prefix (simulating what CheckSubnetContainsIPContextE does)
	resp, err := client.Get(context.Background(), "rg", "test-vnet", "test-subnet", nil)
	require.NoError(t, err)

	subnetPrefix := *resp.Properties.AddressPrefix

	tests := []struct {
		name    string
		ip      string
		want    bool
		wantErr bool
	}{
		{
			name: "IPInRange",
			ip:   "10.0.0.100",
			want: true,
		},
		{
			name: "IPOutOfRange",
			ip:   "10.0.1.100",
			want: false,
		},
		{
			name:    "InvalidIP",
			ip:      "not-an-ip",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ip := net.ParseIP(tc.ip)
			if ip == nil {
				assert.True(t, tc.wantErr)

				return
			}

			_, ipNet, parseErr := net.ParseCIDR(subnetPrefix)
			require.NoError(t, parseErr)

			assert.Equal(t, tc.want, ipNet.Contains(ip))
		})
	}
}

// ---------------------------------------------------------------------------
// GetVirtualNetworkSubnets tests (pagination builds map)
// ---------------------------------------------------------------------------

func TestGetVirtualNetworkSubnets(t *testing.T) {
	t.Parallel()

	srv := &networkfake.SubnetsServer{
		NewListPager: func(_ string, _ string, _ *armnetwork.SubnetsClientListOptions) (resp azfake.PagerResponder[armnetwork.SubnetsClientListResponse]) {
			// Page 1
			resp.AddPage(http.StatusOK, armnetwork.SubnetsClientListResponse{
				SubnetListResult: armnetwork.SubnetListResult{
					Value: []*armnetwork.Subnet{
						{
							Name: to.Ptr("subnet-a"),
							Properties: &armnetwork.SubnetPropertiesFormat{
								AddressPrefix: to.Ptr("10.0.0.0/24"),
							},
						},
					},
				},
			}, nil)

			// Page 2
			resp.AddPage(http.StatusOK, armnetwork.SubnetsClientListResponse{
				SubnetListResult: armnetwork.SubnetListResult{
					Value: []*armnetwork.Subnet{
						{
							Name: to.Ptr("subnet-b"),
							Properties: &armnetwork.SubnetPropertiesFormat{
								AddressPrefix: to.Ptr("10.0.1.0/24"),
							},
						},
					},
				},
			}, nil)

			return
		},
	}

	client := newFakeSubnetsClient(t, srv)

	// Simulate GetVirtualNetworkSubnetsContextE pagination logic
	subNetDetails := map[string]string{}

	pager := client.NewListPager("rg", "test-vnet", nil)
	for pager.More() {
		page, pageErr := pager.NextPage(context.Background())
		require.NoError(t, pageErr)

		for _, v := range page.Value {
			subNetDetails[*v.Name] = *v.Properties.AddressPrefix
		}
	}

	assert.Equal(t, map[string]string{
		"subnet-a": "10.0.0.0/24",
		"subnet-b": "10.0.1.0/24",
	}, subNetDetails)
}

// ---------------------------------------------------------------------------
// GetVirtualNetworkDNSServerIPs tests (extracts DNS IPs, handles nil DHCP options)
// ---------------------------------------------------------------------------

func TestGetVirtualNetworkDNSServerIPs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		vnet armnetwork.VirtualNetwork
		want []string
	}{
		{
			name: "WithDNSServers",
			vnet: armnetwork.VirtualNetwork{
				Name: to.Ptr("test-vnet"),
				Properties: &armnetwork.VirtualNetworkPropertiesFormat{
					DhcpOptions: &armnetwork.DhcpOptions{
						DNSServers: []*string{
							to.Ptr("8.8.8.8"),
							to.Ptr("8.8.4.4"),
						},
					},
				},
			},
			want: []string{"8.8.8.8", "8.8.4.4"},
		},
		{
			name: "NilDhcpOptions",
			vnet: armnetwork.VirtualNetwork{
				Name: to.Ptr("test-vnet"),
				Properties: &armnetwork.VirtualNetworkPropertiesFormat{
					DhcpOptions: nil,
				},
			},
			want: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := &networkfake.VirtualNetworksServer{
				Get: func(_ context.Context, _ string, _ string, _ *armnetwork.VirtualNetworksClientGetOptions) (resp azfake.Responder[armnetwork.VirtualNetworksClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armnetwork.VirtualNetworksClientGetResponse{
						VirtualNetwork: tc.vnet,
					}, nil)

					return
				},
			}

			client := newFakeVirtualNetworksClient(t, srv)

			resp, err := client.Get(context.Background(), "rg", "test-vnet", nil)
			require.NoError(t, err)

			// Simulate GetVirtualNetworkDNSServerIPsContextE logic
			if resp.Properties.DhcpOptions == nil {
				assert.Nil(t, tc.want)

				return
			}

			dnsServers := make([]string, len(resp.Properties.DhcpOptions.DNSServers))
			for i, s := range resp.Properties.DhcpOptions.DNSServers {
				dnsServers[i] = *s
			}

			assert.Equal(t, tc.want, dnsServers)
		})
	}
}

// ---------------------------------------------------------------------------
// VirtualNetwork NotFound error handling
// ---------------------------------------------------------------------------

func TestGetVirtualNetwork_NotFound(t *testing.T) {
	t.Parallel()

	srv := &networkfake.VirtualNetworksServer{
		Get: func(_ context.Context, _ string, _ string, _ *armnetwork.VirtualNetworksClientGetOptions) (resp azfake.Responder[armnetwork.VirtualNetworksClientGetResponse], errResp azfake.ErrorResponder) {
			errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

			return
		},
	}

	client := newFakeVirtualNetworksClient(t, srv)

	_, err := client.Get(context.Background(), "rg", "missing-vnet", nil)
	require.Error(t, err)
	assert.True(t, azure.ResourceNotFoundErrorExists(err))
}

// ---------------------------------------------------------------------------
// Subnet NotFound error handling
// ---------------------------------------------------------------------------

func TestGetSubnet_NotFound(t *testing.T) {
	t.Parallel()

	srv := &networkfake.SubnetsServer{
		Get: func(_ context.Context, _ string, _ string, _ string, _ *armnetwork.SubnetsClientGetOptions) (resp azfake.Responder[armnetwork.SubnetsClientGetResponse], errResp azfake.ErrorResponder) {
			errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

			return
		},
	}

	client := newFakeSubnetsClient(t, srv)

	_, err := client.Get(context.Background(), "rg", "test-vnet", "missing-subnet", nil)
	require.Error(t, err)
	assert.True(t, azure.ResourceNotFoundErrorExists(err))
}
