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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	appservicefake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Fake client helper
// ---------------------------------------------------------------------------

func newFakeWebAppsClient(t *testing.T, srv *appservicefake.WebAppsServer) *armappservice.WebAppsClient {
	t.Helper()

	client, err := armappservice.NewWebAppsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: appservicefake.NewWebAppsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetAppServiceWithClient tests
// ---------------------------------------------------------------------------

func TestGetAppServiceWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    appservicefake.WebAppsServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: appservicefake.WebAppsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armappservice.WebAppsClientGetOptions) (resp azfake.Responder[armappservice.WebAppsClientGetResponse], errResp azfake.ErrorResponder) {
					result := armappservice.WebAppsClientGetResponse{
						Site: armappservice.Site{
							Name: to.Ptr("test-app"),
						},
					}
					resp.SetResponse(http.StatusOK, result, nil)

					return
				},
			},
			wantName: "test-app",
		},
		{
			name: "NotFound",
			server: appservicefake.WebAppsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armappservice.WebAppsClientGetOptions) (resp azfake.Responder[armappservice.WebAppsClientGetResponse], errResp azfake.ErrorResponder) {
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

			client := newFakeWebAppsClient(t, &tc.server)

			site, err := azure.GetAppServiceWithClient(context.Background(), client, "rg", "app")
			if tc.wantErr {
				require.Error(t, err)

				var respErr *azcore.ResponseError
				require.ErrorAs(t, err, &respErr)
				assert.Equal(t, tc.errSubstr, respErr.ErrorCode)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *site.Name)
		})
	}
}
