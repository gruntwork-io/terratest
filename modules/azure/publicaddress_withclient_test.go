package azure_test

import (
	"context"
	"net/http"
	"testing"

	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	networkfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPublicIPAddressWithClient(t *testing.T) {
	t.Parallel()

	srv := &networkfake.PublicIPAddressesServer{
		Get: func(_ context.Context, _ string, pipName string, _ *armnetwork.PublicIPAddressesClientGetOptions) (resp azfake.Responder[armnetwork.PublicIPAddressesClientGetResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armnetwork.PublicIPAddressesClientGetResponse{
				PublicIPAddress: armnetwork.PublicIPAddress{
					Name: to.Ptr(pipName),
					Properties: &armnetwork.PublicIPAddressPropertiesFormat{
						IPAddress: to.Ptr("52.168.1.1"),
					},
				},
			}, nil)

			return
		},
	}
	client := newFakePublicIPAddressesClient(t, srv)

	pip, err := azure.GetPublicIPAddressWithClient(t.Context(), client, "rg", "my-pip")
	require.NoError(t, err)
	assert.Equal(t, "my-pip", *pip.Name)
	assert.Equal(t, "52.168.1.1", *pip.Properties.IPAddress)
}

func TestExtractIPOfPublicIPAddress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		pip     *armnetwork.PublicIPAddress
		wantIP  string
		wantErr bool
	}{
		{
			name: "has IP",
			pip: &armnetwork.PublicIPAddress{
				Name: to.Ptr("pip-1"),
				Properties: &armnetwork.PublicIPAddressPropertiesFormat{
					IPAddress: to.Ptr("52.168.1.1"),
				},
			},
			wantIP: "52.168.1.1",
		},
		{
			name: "nil properties",
			pip: &armnetwork.PublicIPAddress{
				Name: to.Ptr("pip-nil"),
			},
			wantErr: true,
		},
		{
			name: "nil IP address",
			pip: &armnetwork.PublicIPAddress{
				Name:       to.Ptr("pip-no-ip"),
				Properties: &armnetwork.PublicIPAddressPropertiesFormat{},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ip, err := azure.ExtractIPOfPublicIPAddress(tc.pip)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.wantIP, ip)
			}
		})
	}
}
