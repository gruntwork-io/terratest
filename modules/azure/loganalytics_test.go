package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v2"
	oifake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v2/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakeWorkspacesClient(t *testing.T, srv *oifake.WorkspacesServer) *armoperationalinsights.WorkspacesClient {
	t.Helper()

	transport := oifake.NewWorkspacesServerTransport(srv)
	client, err := armoperationalinsights.NewWorkspacesClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func TestGetLogAnalyticsWorkspaceWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		server  oifake.WorkspacesServer
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			server: oifake.WorkspacesServer{
				Get: func(_ context.Context, _ string, workspaceName string, _ *armoperationalinsights.WorkspacesClientGetOptions) (resp azfake.Responder[armoperationalinsights.WorkspacesClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armoperationalinsights.WorkspacesClientGetResponse{
						Workspace: armoperationalinsights.Workspace{
							Name: to.Ptr(workspaceName),
							Properties: &armoperationalinsights.WorkspaceProperties{
								CustomerID: to.Ptr("workspace-id-123"),
							},
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: oifake.WorkspacesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armoperationalinsights.WorkspacesClientGetOptions) (resp azfake.Responder[armoperationalinsights.WorkspacesClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeWorkspacesClient(t, &srv)

			ws, err := azure.GetLogAnalyticsWorkspaceWithClient(t.Context(), client, "rg", "my-workspace")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-workspace", *ws.Name)
			assert.Equal(t, "workspace-id-123", *ws.Properties.CustomerID)
		})
	}
}
