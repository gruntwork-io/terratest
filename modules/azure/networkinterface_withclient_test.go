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

func newFakeInterfacesClient(t *testing.T, srv *networkfake.InterfacesServer) *armnetwork.InterfacesClient {
	t.Helper()

	transport := networkfake.NewInterfacesServerTransport(srv)
	client, err := armnetwork.NewInterfacesClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func TestGetNetworkInterfaceWithClient(t *testing.T) {
	t.Parallel()

	srv := &networkfake.InterfacesServer{
		Get: func(_ context.Context, _ string, nicName string, _ *armnetwork.InterfacesClientGetOptions) (resp azfake.Responder[armnetwork.InterfacesClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.InterfacesClientGetResponse{
				Interface: armnetwork.Interface{
					Name: to.Ptr(nicName),
					Properties: &armnetwork.InterfacePropertiesFormat{
						IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
							{
								Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
									PrivateIPAddress: to.Ptr("10.0.0.4"),
								},
							},
						},
					},
				},
			}, nil)

			return
		},
	}
	client := newFakeInterfacesClient(t, srv)

	nic, err := azure.GetNetworkInterfaceWithClient(t.Context(), client, "rg", "my-nic")
	require.NoError(t, err)
	assert.Equal(t, "my-nic", *nic.Name)
}

func TestExtractNetworkInterfacePrivateIPs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		configs  []*armnetwork.InterfaceIPConfiguration
		expected []string
	}{
		{
			name: "single IP",
			configs: []*armnetwork.InterfaceIPConfiguration{
				{Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{PrivateIPAddress: to.Ptr("10.0.0.4")}},
			},
			expected: []string{"10.0.0.4"},
		},
		{
			name: "multiple IPs",
			configs: []*armnetwork.InterfaceIPConfiguration{
				{Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{PrivateIPAddress: to.Ptr("10.0.0.4")}},
				{Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{PrivateIPAddress: to.Ptr("10.0.0.5")}},
				{Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{PrivateIPAddress: to.Ptr("10.0.0.6")}},
			},
			expected: []string{"10.0.0.4", "10.0.0.5", "10.0.0.6"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			nic := &armnetwork.Interface{
				Properties: &armnetwork.InterfacePropertiesFormat{
					IPConfigurations: tc.configs,
				},
			}

			ips := azure.ExtractNetworkInterfacePrivateIPs(nic)
			assert.Equal(t, tc.expected, ips)
		})
	}
}
