package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	sqlfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Fake client helpers
// ---------------------------------------------------------------------------

func newFakeManagedInstancesClient(t *testing.T, srv *sqlfake.ManagedInstancesServer) *armsql.ManagedInstancesClient {
	t.Helper()

	client, err := armsql.NewManagedInstancesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: sqlfake.NewManagedInstancesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeManagedDatabasesClient(t *testing.T, srv *sqlfake.ManagedDatabasesServer) *armsql.ManagedDatabasesClient {
	t.Helper()

	client, err := armsql.NewManagedDatabasesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: sqlfake.NewManagedDatabasesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetManagedInstanceWithClient tests
// ---------------------------------------------------------------------------

func TestGetManagedInstanceWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  sqlfake.ManagedInstancesServer
		wantErr bool
	}{
		{
			name: "Success",
			server: sqlfake.ManagedInstancesServer{
				Get: func(_ context.Context, _ string, instanceName string, _ *armsql.ManagedInstancesClientGetOptions) (resp azfake.Responder[armsql.ManagedInstancesClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armsql.ManagedInstancesClientGetResponse{
						ManagedInstance: armsql.ManagedInstance{
							Name: to.Ptr(instanceName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: sqlfake.ManagedInstancesServer{
				Get: func(_ context.Context, _ string, _ string, _ *armsql.ManagedInstancesClientGetOptions) (resp azfake.Responder[armsql.ManagedInstancesClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeManagedInstancesClient(t, &srv)

			instance, err := azure.GetManagedInstanceWithClient(context.Background(), client, "rg", "my-instance")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-instance", *instance.Name)
		})
	}
}

// ---------------------------------------------------------------------------
// GetManagedInstanceDatabaseWithClient tests
// ---------------------------------------------------------------------------

func TestGetManagedInstanceDatabaseWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  sqlfake.ManagedDatabasesServer
		wantErr bool
	}{
		{
			name: "Success",
			server: sqlfake.ManagedDatabasesServer{
				Get: func(_ context.Context, _ string, _ string, dbName string, _ *armsql.ManagedDatabasesClientGetOptions) (resp azfake.Responder[armsql.ManagedDatabasesClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armsql.ManagedDatabasesClientGetResponse{
						ManagedDatabase: armsql.ManagedDatabase{
							Name: to.Ptr(dbName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: sqlfake.ManagedDatabasesServer{
				Get: func(_ context.Context, _ string, _ string, _ string, _ *armsql.ManagedDatabasesClientGetOptions) (resp azfake.Responder[armsql.ManagedDatabasesClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeManagedDatabasesClient(t, &srv)

			db, err := azure.GetManagedInstanceDatabaseWithClient(context.Background(), client, "rg", "my-instance", "my-db")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-db", *db.Name)
		})
	}
}
