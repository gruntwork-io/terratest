package azure_test

import (
	"context"
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

func newFakeVirtualNetworksClient(t *testing.T, srv *networkfake.VirtualNetworksServer) *armnetwork.VirtualNetworksClient {
	t.Helper()

	transport := networkfake.NewVirtualNetworksServerTransport(srv)
	client, err := armnetwork.NewVirtualNetworksClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func newFakeSubnetsClient(t *testing.T, srv *networkfake.SubnetsServer) *armnetwork.SubnetsClient {
	t.Helper()

	transport := networkfake.NewSubnetsServerTransport(srv)
	client, err := armnetwork.NewSubnetsClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func TestGetVirtualNetworkWithClient(t *testing.T) {
	t.Parallel()

	srv := &networkfake.VirtualNetworksServer{
		Get: func(_ context.Context, _ string, vnetName string, _ *armnetwork.VirtualNetworksClientGetOptions) (resp azfake.Responder[armnetwork.VirtualNetworksClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.VirtualNetworksClientGetResponse{
				VirtualNetwork: armnetwork.VirtualNetwork{
					Name:       to.Ptr(vnetName),
					Properties: &armnetwork.VirtualNetworkPropertiesFormat{},
				},
			}, nil)

			return
		},
	}
	client := newFakeVirtualNetworksClient(t, srv)

	vnet, err := azure.GetVirtualNetworkWithClient(t.Context(), client, "rg", "my-vnet")
	require.NoError(t, err)
	assert.Equal(t, "my-vnet", *vnet.Name)
}

func TestGetSubnetWithClient(t *testing.T) {
	t.Parallel()

	srv := &networkfake.SubnetsServer{
		Get: func(_ context.Context, _ string, _ string, subnetName string, _ *armnetwork.SubnetsClientGetOptions) (resp azfake.Responder[armnetwork.SubnetsClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.SubnetsClientGetResponse{
				Subnet: armnetwork.Subnet{
					Name: to.Ptr(subnetName),
					Properties: &armnetwork.SubnetPropertiesFormat{
						AddressPrefix: to.Ptr("10.0.1.0/24"),
					},
				},
			}, nil)

			return
		},
	}
	client := newFakeSubnetsClient(t, srv)

	subnet, err := azure.GetSubnetWithClient(t.Context(), client, "rg", "my-vnet", "my-subnet")
	require.NoError(t, err)
	assert.Equal(t, "my-subnet", *subnet.Name)
	assert.Equal(t, "10.0.1.0/24", *subnet.Properties.AddressPrefix)
}

func TestGetVirtualNetworkSubnetsWithClient(t *testing.T) {
	t.Parallel()

	srv := &networkfake.SubnetsServer{
		NewListPager: func(_ string, _ string, _ *armnetwork.SubnetsClientListOptions) (resp azfake.PagerResponder[armnetwork.SubnetsClientListResponse]) {
			resp.AddPage(http.StatusOK, armnetwork.SubnetsClientListResponse{
				SubnetListResult: armnetwork.SubnetListResult{
					Value: []*armnetwork.Subnet{
						{
							Name:       to.Ptr("subnet-a"),
							Properties: &armnetwork.SubnetPropertiesFormat{AddressPrefix: to.Ptr("10.0.1.0/24")},
						},
					},
				},
			}, nil)
			resp.AddPage(http.StatusOK, armnetwork.SubnetsClientListResponse{
				SubnetListResult: armnetwork.SubnetListResult{
					Value: []*armnetwork.Subnet{
						{
							Name:       to.Ptr("subnet-b"),
							Properties: &armnetwork.SubnetPropertiesFormat{AddressPrefix: to.Ptr("10.0.2.0/24")},
						},
					},
				},
			}, nil)

			return
		},
	}
	client := newFakeSubnetsClient(t, srv)

	subnets, err := azure.GetVirtualNetworkSubnetsWithClient(t.Context(), client, "rg", "my-vnet")
	require.NoError(t, err)
	assert.Equal(t, map[string]string{
		"subnet-a": "10.0.1.0/24",
		"subnet-b": "10.0.2.0/24",
	}, subnets)
}

func TestCheckSubnetContainsIPWithClient(t *testing.T) {
	t.Parallel()

	subnetServer := &networkfake.SubnetsServer{
		Get: func(_ context.Context, _ string, _ string, _ string, _ *armnetwork.SubnetsClientGetOptions) (resp azfake.Responder[armnetwork.SubnetsClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.SubnetsClientGetResponse{
				Subnet: armnetwork.Subnet{
					Name: to.Ptr("my-subnet"),
					Properties: &armnetwork.SubnetPropertiesFormat{
						AddressPrefix: to.Ptr("10.0.1.0/24"),
					},
				},
			}, nil)

			return
		},
	}

	tests := []struct {
		name    string
		ip      string
		want    bool
		wantErr bool
	}{
		{
			name: "IP in range",
			ip:   "10.0.1.5",
			want: true,
		},
		{
			name: "IP out of range",
			ip:   "10.0.2.5",
			want: false,
		},
		{
			name:    "invalid IP",
			ip:      "not-an-ip",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newFakeSubnetsClient(t, subnetServer)

			got, err := azure.CheckSubnetContainsIPWithClient(t.Context(), client, tc.ip, "my-subnet", "my-vnet", "rg")
			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, got)
			}
		})
	}
}

func TestExtractVirtualNetworkDNSServerIPs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		vnet     *armnetwork.VirtualNetwork
		expected []string
	}{
		{
			name: "has DNS servers",
			vnet: &armnetwork.VirtualNetwork{
				Properties: &armnetwork.VirtualNetworkPropertiesFormat{
					DhcpOptions: &armnetwork.DhcpOptions{
						DNSServers: []*string{to.Ptr("8.8.8.8"), to.Ptr("8.8.4.4")},
					},
				},
			},
			expected: []string{"8.8.8.8", "8.8.4.4"},
		},
		{
			name: "nil DhcpOptions",
			vnet: &armnetwork.VirtualNetwork{
				Properties: &armnetwork.VirtualNetworkPropertiesFormat{},
			},
			expected: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := azure.ExtractVirtualNetworkDNSServerIPs(tc.vnet)
			assert.Equal(t, tc.expected, result)
		})
	}
}
