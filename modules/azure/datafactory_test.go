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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v9"
	datafactoryfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v9/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakeFactoriesClient(t *testing.T, srv *datafactoryfake.FactoriesServer) *armdatafactory.FactoriesClient {
	t.Helper()

	client, err := armdatafactory.NewFactoriesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: datafactoryfake.NewFactoriesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func TestGetDataFactoryWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    datafactoryfake.FactoriesServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: datafactoryfake.FactoriesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armdatafactory.FactoriesClientGetOptions) (resp azfake.Responder[armdatafactory.FactoriesClientGetResponse], errResp azfake.ErrorResponder) {
					result := armdatafactory.FactoriesClientGetResponse{
						Factory: armdatafactory.Factory{
							Name: to.Ptr("test-factory"),
						},
					}
					resp.SetResponse(http.StatusOK, result, nil)

					return
				},
			},
			wantName: "test-factory",
		},
		{
			name: "NotFound",
			server: datafactoryfake.FactoriesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armdatafactory.FactoriesClientGetOptions) (resp azfake.Responder[armdatafactory.FactoriesClientGetResponse], errResp azfake.ErrorResponder) {
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

			client := newFakeFactoriesClient(t, &tc.server)

			factory, err := azure.GetDataFactoryWithClient(context.Background(), client, "rg", "test-factory")
			if tc.wantErr {
				require.Error(t, err)

				var respErr *azcore.ResponseError
				require.ErrorAs(t, err, &respErr)
				assert.Equal(t, tc.errSubstr, respErr.ErrorCode)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *factory.Name)
		})
	}
}
