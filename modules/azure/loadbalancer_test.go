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

	client, err := armnetwork.NewLoadBalancersClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: networkfake.NewLoadBalancersServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeLBFrontendIPConfigClient(t *testing.T, srv *networkfake.LoadBalancerFrontendIPConfigurationsServer) *armnetwork.LoadBalancerFrontendIPConfigurationsClient {
	t.Helper()

	client, err := armnetwork.NewLoadBalancerFrontendIPConfigurationsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: networkfake.NewLoadBalancerFrontendIPConfigurationsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetLoadBalancerFrontendIPConfigNames tests
// ---------------------------------------------------------------------------

func TestGetLoadBalancerFrontendIPConfigNames(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name   string
		server networkfake.LoadBalancersServer
		want   []string
	}{
		{
			name: "TwoConfigs",
			server: networkfake.LoadBalancersServer{
				Get: func(_ context.Context, _ string, _ string, _ *armnetwork.LoadBalancersClientGetOptions) (resp azfake.Responder[armnetwork.LoadBalancersClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armnetwork.LoadBalancersClientGetResponse{
						LoadBalancer: armnetwork.LoadBalancer{
							Name: to.Ptr("test-lb"),
							Properties: &armnetwork.LoadBalancerPropertiesFormat{
								FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{
									{Name: to.Ptr("config-public")},
									{Name: to.Ptr("config-private")},
								},
							},
						},
					}, nil)

					return
				},
			},
			want: []string{"config-public", "config-private"},
		},
		{
			name: "EmptyConfigs",
			server: networkfake.LoadBalancersServer{
				Get: func(_ context.Context, _ string, _ string, _ *armnetwork.LoadBalancersClientGetOptions) (resp azfake.Responder[armnetwork.LoadBalancersClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armnetwork.LoadBalancersClientGetResponse{
						LoadBalancer: armnetwork.LoadBalancer{
							Name: to.Ptr("test-lb"),
							Properties: &armnetwork.LoadBalancerPropertiesFormat{
								FrontendIPConfigurations: []*armnetwork.FrontendIPConfiguration{},
							},
						},
					}, nil)

					return
				},
			},
			want: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newFakeLoadBalancersClient(t, &tc.server)

			resp, err := client.Get(context.Background(), "rg", "test-lb", nil)
			require.NoError(t, err)

			feConfigs := resp.Properties.FrontendIPConfigurations
			if len(feConfigs) == 0 {
				assert.Nil(t, tc.want)

				return
			}

			configNames := make([]string, len(feConfigs))
			for i, config := range feConfigs {
				configNames[i] = *config.Name
			}

			assert.Equal(t, tc.want, configNames)
		})
	}
}

// ---------------------------------------------------------------------------
// GetIPOfLoadBalancerFrontendIPConfig tests (public vs private detection)
// ---------------------------------------------------------------------------

func TestGetIPOfLoadBalancerFrontendIPConfig_PublicIP(t *testing.T) {
	t.Parallel()

	srv := &networkfake.LoadBalancerFrontendIPConfigurationsServer{
		Get: func(_ context.Context, _ string, _ string, _ string, _ *armnetwork.LoadBalancerFrontendIPConfigurationsClientGetOptions) (resp azfake.Responder[armnetwork.LoadBalancerFrontendIPConfigurationsClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.LoadBalancerFrontendIPConfigurationsClientGetResponse{
				FrontendIPConfiguration: armnetwork.FrontendIPConfiguration{
					Name: to.Ptr("public-config"),
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						PublicIPAddress: &armnetwork.PublicIPAddress{
							ID: to.Ptr("/subscriptions/sub/resourceGroups/rg/providers/Microsoft.Network/publicIPAddresses/my-pip"),
						},
					},
				},
			}, nil)

			return
		},
	}

	client := newFakeLBFrontendIPConfigClient(t, srv)

	resp, err := client.Get(context.Background(), "rg", "test-lb", "public-config", nil)
	require.NoError(t, err)

	// Verify the response has a public IP reference (indicating PublicIP type)
	assert.NotNil(t, resp.Properties.PublicIPAddress)
	assert.Equal(t, "my-pip", azure.GetNameFromResourceID(*resp.Properties.PublicIPAddress.ID))
}

func TestGetIPOfLoadBalancerFrontendIPConfig_PrivateIP(t *testing.T) {
	t.Parallel()

	srv := &networkfake.LoadBalancerFrontendIPConfigurationsServer{
		Get: func(_ context.Context, _ string, _ string, _ string, _ *armnetwork.LoadBalancerFrontendIPConfigurationsClientGetOptions) (resp azfake.Responder[armnetwork.LoadBalancerFrontendIPConfigurationsClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.LoadBalancerFrontendIPConfigurationsClientGetResponse{
				FrontendIPConfiguration: armnetwork.FrontendIPConfiguration{
					Name: to.Ptr("private-config"),
					Properties: &armnetwork.FrontendIPConfigurationPropertiesFormat{
						PrivateIPAddress: to.Ptr("10.0.1.5"),
					},
				},
			}, nil)

			return
		},
	}

	client := newFakeLBFrontendIPConfigClient(t, srv)

	resp, err := client.Get(context.Background(), "rg", "test-lb", "private-config", nil)
	require.NoError(t, err)

	// No public IP means this is a private frontend config
	assert.Nil(t, resp.Properties.PublicIPAddress)
	assert.Equal(t, "10.0.1.5", *resp.Properties.PrivateIPAddress)
}

// ---------------------------------------------------------------------------
// LoadBalancer NotFound error handling
// ---------------------------------------------------------------------------

func TestGetLoadBalancer_NotFound(t *testing.T) {
	t.Parallel()

	srv := &networkfake.LoadBalancersServer{
		Get: func(_ context.Context, _ string, _ string, _ *armnetwork.LoadBalancersClientGetOptions) (resp azfake.Responder[armnetwork.LoadBalancersClientGetResponse], errResp azfake.ErrorResponder) {
			errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

			return
		},
	}

	client := newFakeLoadBalancersClient(t, srv)

	_, err := client.Get(context.Background(), "rg", "missing-lb", nil)
	require.Error(t, err)
	assert.True(t, azure.ResourceNotFoundErrorExists(err))
}
