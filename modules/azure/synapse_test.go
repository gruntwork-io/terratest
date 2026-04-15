package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse"
	synapsefake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Fake client helpers
// ---------------------------------------------------------------------------

func newFakeSynapseWorkspacesClient(t *testing.T, srv *synapsefake.WorkspacesServer) *armsynapse.WorkspacesClient {
	t.Helper()

	client, err := armsynapse.NewWorkspacesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: synapsefake.NewWorkspacesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeSQLPoolsClient(t *testing.T, srv *synapsefake.SQLPoolsServer) *armsynapse.SQLPoolsClient {
	t.Helper()

	client, err := armsynapse.NewSQLPoolsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: synapsefake.NewSQLPoolsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetSynapseWorkspaceWithClient tests
// ---------------------------------------------------------------------------

func TestGetSynapseWorkspaceWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  synapsefake.WorkspacesServer
		wantErr bool
	}{
		{
			name: "Success",
			server: synapsefake.WorkspacesServer{
				Get: func(_ context.Context, _ string, workspaceName string, _ *armsynapse.WorkspacesClientGetOptions) (resp azfake.Responder[armsynapse.WorkspacesClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armsynapse.WorkspacesClientGetResponse{
						Workspace: armsynapse.Workspace{
							Name: to.Ptr(workspaceName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: synapsefake.WorkspacesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armsynapse.WorkspacesClientGetOptions) (resp azfake.Responder[armsynapse.WorkspacesClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeSynapseWorkspacesClient(t, &srv)

			ws, err := azure.GetSynapseWorkspaceWithClient(context.Background(), client, "rg", "my-workspace")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-workspace", *ws.Name)
		})
	}
}

// ---------------------------------------------------------------------------
// GetSynapseSQLPoolWithClient tests
// ---------------------------------------------------------------------------

func TestGetSynapseSQLPoolWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  synapsefake.SQLPoolsServer
		wantErr bool
	}{
		{
			name: "Success",
			server: synapsefake.SQLPoolsServer{
				Get: func(_ context.Context, _ string, _ string, poolName string, _ *armsynapse.SQLPoolsClientGetOptions) (resp azfake.Responder[armsynapse.SQLPoolsClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armsynapse.SQLPoolsClientGetResponse{
						SQLPool: armsynapse.SQLPool{
							Name: to.Ptr(poolName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: synapsefake.SQLPoolsServer{
				Get: func(_ context.Context, _ string, _ string, _ string, _ *armsynapse.SQLPoolsClientGetOptions) (resp azfake.Responder[armsynapse.SQLPoolsClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeSQLPoolsClient(t, &srv)

			pool, err := azure.GetSynapseSQLPoolWithClient(context.Background(), client, "rg", "my-workspace", "my-pool")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-pool", *pool.Name)
		})
	}
}
