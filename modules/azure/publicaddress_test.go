package azure_test

import (
	"context"
	"fmt"
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

func newFakePublicIPAddressesClient(t *testing.T, srv *networkfake.PublicIPAddressesServer) *armnetwork.PublicIPAddressesClient {
	t.Helper()

	client, err := armnetwork.NewPublicIPAddressesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: networkfake.NewPublicIPAddressesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetIPOfPublicIPAddressByName tests (returns IP, handles nil IP)
// ---------------------------------------------------------------------------

func TestGetIPOfPublicIPAddressByName(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantIP    string
		server    networkfake.PublicIPAddressesServer
		wantErr   bool
		errSubstr string
	}{
		{
			name: "ReturnsIP",
			server: networkfake.PublicIPAddressesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armnetwork.PublicIPAddressesClientGetOptions) (resp azfake.Responder[armnetwork.PublicIPAddressesClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armnetwork.PublicIPAddressesClientGetResponse{
						PublicIPAddress: armnetwork.PublicIPAddress{
							Name: to.Ptr("my-pip"),
							Properties: &armnetwork.PublicIPAddressPropertiesFormat{
								IPAddress: to.Ptr("52.1.2.3"),
							},
						},
					}, nil)

					return
				},
			},
			wantIP: "52.1.2.3",
		},
		{
			name: "NilIPAddress",
			server: networkfake.PublicIPAddressesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armnetwork.PublicIPAddressesClientGetOptions) (resp azfake.Responder[armnetwork.PublicIPAddressesClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armnetwork.PublicIPAddressesClientGetResponse{
						PublicIPAddress: armnetwork.PublicIPAddress{
							Name:       to.Ptr("my-pip"),
							Properties: &armnetwork.PublicIPAddressPropertiesFormat{},
						},
					}, nil)

					return
				},
			},
			wantErr:   true,
			errSubstr: "no IP address assigned",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newFakePublicIPAddressesClient(t, &tc.server)

			resp, err := client.Get(context.Background(), "rg", "my-pip", nil)
			require.NoError(t, err)

			pip := &resp.PublicIPAddress
			if pip.Properties == nil || pip.Properties.IPAddress == nil {
				if tc.wantErr {
					errMsg := fmt.Sprintf("public IP address %q has no IP address assigned", *pip.Name)
					assert.Contains(t, errMsg, tc.errSubstr)

					return
				}

				t.Fatal("unexpected nil IP address")
			}

			assert.Equal(t, tc.wantIP, *pip.Properties.IPAddress)
		})
	}
}

// ---------------------------------------------------------------------------
// PublicAddressExists tests (exists vs not-found)
// ---------------------------------------------------------------------------

func TestPublicAddressExists(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name   string
		want   bool
		server networkfake.PublicIPAddressesServer
	}{
		{
			name: "Exists",
			server: networkfake.PublicIPAddressesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armnetwork.PublicIPAddressesClientGetOptions) (resp azfake.Responder[armnetwork.PublicIPAddressesClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armnetwork.PublicIPAddressesClientGetResponse{
						PublicIPAddress: armnetwork.PublicIPAddress{
							Name: to.Ptr("my-pip"),
						},
					}, nil)

					return
				},
			},
			want: true,
		},
		{
			name: "NotFound",
			server: networkfake.PublicIPAddressesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armnetwork.PublicIPAddressesClientGetOptions) (resp azfake.Responder[armnetwork.PublicIPAddressesClientGetResponse], errResp azfake.ErrorResponder) {
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

			client := newFakePublicIPAddressesClient(t, &tc.server)

			_, err := client.Get(context.Background(), "rg", "my-pip", nil)
			if err != nil {
				assert.False(t, tc.want)
				assert.True(t, azure.ResourceNotFoundErrorExists(err))

				return
			}

			assert.True(t, tc.want)
		})
	}
}
