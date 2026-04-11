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

	client, err := armnetwork.NewInterfacesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: networkfake.NewInterfacesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetNetworkInterfacePrivateIPs tests (extracts private IPs from configs)
// ---------------------------------------------------------------------------

func TestGetNetworkInterfacePrivateIPs(t *testing.T) {
	t.Parallel()

	srv := &networkfake.InterfacesServer{
		Get: func(_ context.Context, _ string, _ string, _ *armnetwork.InterfacesClientGetOptions) (resp azfake.Responder[armnetwork.InterfacesClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.InterfacesClientGetResponse{
				Interface: armnetwork.Interface{
					Name: to.Ptr("test-nic"),
					Properties: &armnetwork.InterfacePropertiesFormat{
						IPConfigurations: []*armnetwork.InterfaceIPConfiguration{
							{
								Name: to.Ptr("ipconfig1"),
								Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
									PrivateIPAddress: to.Ptr("10.0.0.4"),
								},
							},
							{
								Name: to.Ptr("ipconfig2"),
								Properties: &armnetwork.InterfaceIPConfigurationPropertiesFormat{
									PrivateIPAddress: to.Ptr("10.0.0.5"),
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

	resp, err := client.Get(context.Background(), "rg", "test-nic", nil)
	require.NoError(t, err)

	privateIPs := make([]string, 0, len(resp.Properties.IPConfigurations))

	for _, ipConfig := range resp.Properties.IPConfigurations {
		privateIPs = append(privateIPs, *ipConfig.Properties.PrivateIPAddress)
	}

	assert.Equal(t, []string{"10.0.0.4", "10.0.0.5"}, privateIPs)
}

// ---------------------------------------------------------------------------
// NetworkInterfaceExists tests (exists vs not-found)
// ---------------------------------------------------------------------------

func TestNetworkInterfaceExists(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name   string
		want   bool
		server networkfake.InterfacesServer
	}{
		{
			name: "Exists",
			server: networkfake.InterfacesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armnetwork.InterfacesClientGetOptions) (resp azfake.Responder[armnetwork.InterfacesClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armnetwork.InterfacesClientGetResponse{
						Interface: armnetwork.Interface{
							Name: to.Ptr("test-nic"),
						},
					}, nil)

					return
				},
			},
			want: true,
		},
		{
			name: "NotFound",
			server: networkfake.InterfacesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armnetwork.InterfacesClientGetOptions) (resp azfake.Responder[armnetwork.InterfacesClientGetResponse], errResp azfake.ErrorResponder) {
					errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

					return
				},
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newFakeInterfacesClient(t, &tc.server)

			_, err := client.Get(context.Background(), "rg", "test-nic", nil)
			if err != nil {
				assert.False(t, tc.want)
				assert.True(t, azure.ResourceNotFoundErrorExists(err))

				return
			}

			assert.True(t, tc.want)
		})
	}
}

// ---------------------------------------------------------------------------
// GetNetworkInterface success test
// ---------------------------------------------------------------------------

func TestGetNetworkInterface_Success(t *testing.T) {
	t.Parallel()

	srv := &networkfake.InterfacesServer{
		Get: func(_ context.Context, _ string, _ string, _ *armnetwork.InterfacesClientGetOptions) (resp azfake.Responder[armnetwork.InterfacesClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.InterfacesClientGetResponse{
				Interface: armnetwork.Interface{
					Name: to.Ptr("my-nic"),
					Properties: &armnetwork.InterfacePropertiesFormat{
						MacAddress: to.Ptr("00-0D-3A-12-34-56"),
					},
				},
			}, nil)

			return
		},
	}

	client := newFakeInterfacesClient(t, srv)

	resp, err := client.Get(context.Background(), "rg", "my-nic", nil)
	require.NoError(t, err)
	assert.Equal(t, "my-nic", *resp.Interface.Name)
	assert.Equal(t, "00-0D-3A-12-34-56", *resp.Properties.MacAddress)
}
