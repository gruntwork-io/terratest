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

func newFakeLoadBalancersClient(t *testing.T, srv *networkfake.LoadBalancersServer) *armnetwork.LoadBalancersClient {
	t.Helper()

	transport := networkfake.NewLoadBalancersServerTransport(srv)
	client, err := armnetwork.NewLoadBalancersClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func newFakeLBFrontendIPConfigClient(t *testing.T, srv *networkfake.LoadBalancerFrontendIPConfigurationsServer) *armnetwork.LoadBalancerFrontendIPConfigurationsClient {
	t.Helper()

	transport := networkfake.NewLoadBalancerFrontendIPConfigurationsServerTransport(srv)
	client, err := armnetwork.NewLoadBalancerFrontendIPConfigurationsClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func newFakePublicIPAddressesClient(t *testing.T, srv *networkfake.PublicIPAddressesServer) *armnetwork.PublicIPAddressesClient {
	t.Helper()

	transport := networkfake.NewPublicIPAddressesServerTransport(srv)
	client, err := armnetwork.NewPublicIPAddressesClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func TestGetLoadBalancerWithClient(t *testing.T) {
	t.Parallel()

	srv := &networkfake.LoadBalancersServer{
		Get: func(_ context.Context, _ string, lbName string, _ *armnetwork.LoadBalancersClientGetOptions) (resp azfake.Responder[armnetwork.LoadBalancersClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.LoadBalancersClientGetResponse{
				LoadBalancer: armnetwork.LoadBalancer{
					Name:       to.Ptr(lbName),
					Properties: &armnetwork.LoadBalancerPropertiesFormat{},
				},
			}, nil)

			return
		},
	}
	client := newFakeLoadBalancersClient(t, srv)

	lb, err := azure.GetLoadBalancerWithClient(t.Context(), client, "rg", "my-lb")
	require.NoError(t, err)
	assert.Equal(t, "my-lb", *lb.Name)
}

func TestGetLoadBalancerFrontendIPConfigWithClient(t *testing.T) {
	t.Parallel()

	srv := &networkfake.LoadBalancerFrontendIPConfigurationsServer{
		Get: func(_ context.Context, _ string, _ string, feConfigName string, _ *armnetwork.LoadBalancerFrontendIPConfigurationsClientGetOptions) (resp azfake.Responder[armnetwork.LoadBalancerFrontendIPConfigurationsClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.LoadBalancerFrontendIPConfigurationsClientGetResponse{
				FrontendIPConfiguration: armnetwork.FrontendIPConfiguration{
					Name: to.Ptr(feConfigName),
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						PrivateIPAddress: to.Ptr("10.0.0.5"),
					},
				},
			}, nil)

			return
		},
	}
	client := newFakeLBFrontendIPConfigClient(t, srv)

	feConfig, err := azure.GetLoadBalancerFrontendIPConfigWithClient(t.Context(), client, "rg", "lb", "fe-config")
	require.NoError(t, err)
	assert.Equal(t, "fe-config", *feConfig.Name)
}

func TestExtractLoadBalancerFrontendIPConfigNames(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		configs  []*armnetwork.FrontendIPConfiguration
		expected []string
	}{
		{
			name:     "no configs",
			configs:  nil,
			expected: nil,
		},
		{
			name: "multiple configs",
			configs: []*armnetwork.FrontendIPConfiguration{
				{Name: to.Ptr("fe-1")},
				{Name: to.Ptr("fe-2")},
			},
			expected: []string{"fe-1", "fe-2"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			lb := &armnetwork.LoadBalancer{
				Properties: &armnetwork.LoadBalancerPropertiesFormat{
					FrontendIPConfigurations: tc.configs,
				},
			}

			names := azure.ExtractLoadBalancerFrontendIPConfigNames(lb)
			assert.Equal(t, tc.expected, names)
		})
	}
}

func TestGetIPOfLoadBalancerFrontendIPConfigWithClient(t *testing.T) {
	t.Parallel()

	happyPIPServer := &networkfake.PublicIPAddressesServer{
		Get: func(_ context.Context, _ string, _ string, _ *armnetwork.PublicIPAddressesClientGetOptions) (resp azfake.Responder[armnetwork.PublicIPAddressesClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.PublicIPAddressesClientGetResponse{
				PublicIPAddress: armnetwork.PublicIPAddress{
					Name: to.Ptr("my-pip"),
					Properties: &armnetwork.PublicIPAddressPropertiesFormat{
						IPAddress: to.Ptr("40.50.60.70"),
					},
				},
			}, nil)

			return
		},
	}

	notFoundPIPServer := &networkfake.PublicIPAddressesServer{
		Get: func(_ context.Context, _ string, _ string, _ *armnetwork.PublicIPAddressesClientGetOptions) (resp azfake.Responder[armnetwork.PublicIPAddressesClientGetResponse], errResp azfake.ErrorResponder) {
			errResp.SetResponseError(http.StatusNotFound, "PublicIPAddressNotFound")

			return
		},
	}

	unassignedPIPServer := &networkfake.PublicIPAddressesServer{
		Get: func(_ context.Context, _ string, _ string, _ *armnetwork.PublicIPAddressesClientGetOptions) (resp azfake.Responder[armnetwork.PublicIPAddressesClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.PublicIPAddressesClientGetResponse{
				PublicIPAddress: armnetwork.PublicIPAddress{
					Name:       to.Ptr("my-pip"),
					Properties: &armnetwork.PublicIPAddressPropertiesFormat{},
				},
			}, nil)

			return
		},
	}

	tests := []struct {
		name      string
		feConfig  *armnetwork.FrontendIPConfiguration
		pipServer *networkfake.PublicIPAddressesServer
		wantIP    string
		wantType  azure.LoadBalancerIPType
		wantErr   bool
	}{
		{
			name: "private IP",
			feConfig: &armnetwork.FrontendIPConfiguration{
				Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
					PrivateIPAddress: to.Ptr("10.0.0.5"),
				},
			},
			pipServer: happyPIPServer,
			wantIP:    "10.0.0.5",
			wantType:  azure.PrivateIP,
		},
		{
			name: "public IP",
			feConfig: &armnetwork.FrontendIPConfiguration{
				Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
					PublicIPAddress: &armnetwork.PublicIPAddress{
						ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Network/publicIPAddresses/my-pip"),
					},
				},
			},
			pipServer: happyPIPServer,
			wantIP:    "40.50.60.70",
			wantType:  azure.PublicIP,
		},
		{
			name: "no private or public IP",
			feConfig: &armnetwork.FrontendIPConfiguration{
				Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{},
			},
			pipServer: happyPIPServer,
			wantErr:   true,
		},
		{
			name: "public IP lookup fails",
			feConfig: &armnetwork.FrontendIPConfiguration{
				Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
					PublicIPAddress: &armnetwork.PublicIPAddress{
						ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Network/publicIPAddresses/gone-pip"),
					},
				},
			},
			pipServer: notFoundPIPServer,
			wantErr:   true,
		},
		{
			name: "public IP exists but unassigned",
			feConfig: &armnetwork.FrontendIPConfiguration{
				Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
					PublicIPAddress: &armnetwork.PublicIPAddress{
						ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Network/publicIPAddresses/my-pip"),
					},
				},
			},
			pipServer: unassignedPIPServer,
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pipClient := newFakePublicIPAddressesClient(t, tc.pipServer)

			ip, ipType, err := azure.GetIPOfLoadBalancerFrontendIPConfigWithClient(t.Context(), tc.feConfig, pipClient, "rg")

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantIP, ip)
				assert.Equal(t, tc.wantType, ipType)
			}
		})
	}
}
