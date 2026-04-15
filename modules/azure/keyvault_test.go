package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault"
	keyvaultfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakeKeyVaultVaultsClient(t *testing.T, srv *keyvaultfake.VaultsServer) *armkeyvault.VaultsClient {
	t.Helper()

	client, err := armkeyvault.NewVaultsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: keyvaultfake.NewVaultsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func TestGetKeyVaultWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    keyvaultfake.VaultsServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: keyvaultfake.VaultsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armkeyvault.VaultsClientGetOptions) (resp azfake.Responder[armkeyvault.VaultsClientGetResponse], errResp azfake.ErrorResponder) {
					result := armkeyvault.VaultsClientGetResponse{
						Vault: armkeyvault.Vault{
							Name: to.Ptr("test-vault"),
						},
					}
					resp.SetResponse(http.StatusOK, result, nil)

					return
				},
			},
			wantName: "test-vault",
		},
		{
			name: "NotFound",
			server: keyvaultfake.VaultsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armkeyvault.VaultsClientGetOptions) (resp azfake.Responder[armkeyvault.VaultsClientGetResponse], errResp azfake.ErrorResponder) {
					errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

					return
				},
			},
			wantErr:   true,
			errSubstr: "ResourceNotFound",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client := newFakeKeyVaultVaultsClient(t, &tc.server)

			vault, err := azure.GetKeyVaultWithClient(context.Background(), client, "rg", "test-vault")
			if tc.wantErr {
				require.Error(t, err)

				var respErr *azcore.ResponseError
				require.ErrorAs(t, err, &respErr)
				assert.Equal(t, tc.errSubstr, respErr.ErrorCode)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *vault.Name)
		})
	}
}
