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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/frontdoor/armfrontdoor"
	frontdoorfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/frontdoor/armfrontdoor/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Fake client helpers
// ---------------------------------------------------------------------------

func newFakeFrontDoorsClient(t *testing.T, srv *frontdoorfake.FrontDoorsServer) *armfrontdoor.FrontDoorsClient {
	t.Helper()

	client, err := armfrontdoor.NewFrontDoorsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: frontdoorfake.NewFrontDoorsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeFrontendEndpointsClient(t *testing.T, srv *frontdoorfake.FrontendEndpointsServer) *armfrontdoor.FrontendEndpointsClient {
	t.Helper()

	client, err := armfrontdoor.NewFrontendEndpointsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: frontdoorfake.NewFrontendEndpointsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetFrontDoorWithClient tests
// ---------------------------------------------------------------------------

func TestGetFrontDoorWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    frontdoorfake.FrontDoorsServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: frontdoorfake.FrontDoorsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armfrontdoor.FrontDoorsClientGetOptions) (resp azfake.Responder[armfrontdoor.FrontDoorsClientGetResponse], errResp azfake.ErrorResponder) {
					result := armfrontdoor.FrontDoorsClientGetResponse{
						FrontDoor: armfrontdoor.FrontDoor{
							Name: to.Ptr("test-frontdoor"),
						},
					}
					resp.SetResponse(http.StatusOK, result, nil)

					return
				},
			},
			wantName: "test-frontdoor",
		},
		{
			name: "NotFound",
			server: frontdoorfake.FrontDoorsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armfrontdoor.FrontDoorsClientGetOptions) (resp azfake.Responder[armfrontdoor.FrontDoorsClientGetResponse], errResp azfake.ErrorResponder) {
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

			client := newFakeFrontDoorsClient(t, &tc.server)

			fd, err := azure.GetFrontDoorWithClient(context.Background(), client, "rg", "fd")
			if tc.wantErr {
				require.Error(t, err)

				var respErr *azcore.ResponseError
				require.ErrorAs(t, err, &respErr)
				assert.Equal(t, tc.errSubstr, respErr.ErrorCode)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *fd.Name)
		})
	}
}

// ---------------------------------------------------------------------------
// GetFrontDoorFrontendEndpointWithClient tests
// ---------------------------------------------------------------------------

func TestGetFrontDoorFrontendEndpointWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    frontdoorfake.FrontendEndpointsServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: frontdoorfake.FrontendEndpointsServer{
				Get: func(_ context.Context, _ string, _ string, _ string, _ *armfrontdoor.FrontendEndpointsClientGetOptions) (resp azfake.Responder[armfrontdoor.FrontendEndpointsClientGetResponse], errResp azfake.ErrorResponder) {
					result := armfrontdoor.FrontendEndpointsClientGetResponse{
						FrontendEndpoint: armfrontdoor.FrontendEndpoint{
							Name: to.Ptr("test-endpoint"),
						},
					}
					resp.SetResponse(http.StatusOK, result, nil)

					return
				},
			},
			wantName: "test-endpoint",
		},
		{
			name: "NotFound",
			server: frontdoorfake.FrontendEndpointsServer{
				Get: func(_ context.Context, _ string, _ string, _ string, _ *armfrontdoor.FrontendEndpointsClientGetOptions) (resp azfake.Responder[armfrontdoor.FrontendEndpointsClientGetResponse], errResp azfake.ErrorResponder) {
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

			client := newFakeFrontendEndpointsClient(t, &tc.server)

			ep, err := azure.GetFrontDoorFrontendEndpointWithClient(context.Background(), client, "rg", "fd", "ep")
			if tc.wantErr {
				require.Error(t, err)

				var respErr *azcore.ResponseError
				require.ErrorAs(t, err, &respErr)
				assert.Equal(t, tc.errSubstr, respErr.ErrorCode)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *ep.Name)
		})
	}
}
