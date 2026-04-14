package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	monitorfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakeActionGroupsClient(t *testing.T, srv *monitorfake.ActionGroupsServer) *armmonitor.ActionGroupsClient {
	t.Helper()

	transport := monitorfake.NewActionGroupsServerTransport(srv)
	client, err := armmonitor.NewActionGroupsClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func TestGetActionGroupResourceWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		server  monitorfake.ActionGroupsServer
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			server: monitorfake.ActionGroupsServer{
				Get: func(_ context.Context, _ string, actionGroupName string, _ *armmonitor.ActionGroupsClientGetOptions) (resp azfake.Responder[armmonitor.ActionGroupsClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armmonitor.ActionGroupsClientGetResponse{
						ActionGroupResource: armmonitor.ActionGroupResource{
							Name: to.Ptr(actionGroupName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: monitorfake.ActionGroupsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armmonitor.ActionGroupsClientGetOptions) (resp azfake.Responder[armmonitor.ActionGroupsClientGetResponse], errResp azfake.ErrorResponder) {
					errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")
					return
				},
			},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := tc.server
			client := newFakeActionGroupsClient(t, &srv)

			resource, err := azure.GetActionGroupResourceWithClient(t.Context(), client, "rg", "my-action-group")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-action-group", *resource.Name)
		})
	}
}
