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

func newFakeSQLServersClient(t *testing.T, srv *sqlfake.ServersServer) *armsql.ServersClient {
	t.Helper()

	client, err := armsql.NewServersClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: sqlfake.NewServersServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeSQLDatabasesClient(t *testing.T, srv *sqlfake.DatabasesServer) *armsql.DatabasesClient {
	t.Helper()

	client, err := armsql.NewDatabasesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: sqlfake.NewDatabasesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetSQLServerWithClient tests
// ---------------------------------------------------------------------------

func TestGetSQLServerWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  sqlfake.ServersServer
		wantErr bool
	}{
		{
			name: "Success",
			server: sqlfake.ServersServer{
				Get: func(_ context.Context, _ string, serverName string, _ *armsql.ServersClientGetOptions) (resp azfake.Responder[armsql.ServersClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armsql.ServersClientGetResponse{
						Server: armsql.Server{
							Name: to.Ptr(serverName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: sqlfake.ServersServer{
				Get: func(_ context.Context, _ string, _ string, _ *armsql.ServersClientGetOptions) (resp azfake.Responder[armsql.ServersClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeSQLServersClient(t, &srv)

			server, err := azure.GetSQLServerWithClient(context.Background(), client, "rg", "my-server")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-server", *server.Name)
		})
	}
}

// ---------------------------------------------------------------------------
// GetSQLDatabaseWithClient tests
// ---------------------------------------------------------------------------

func TestGetSQLDatabaseWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  sqlfake.DatabasesServer
		wantErr bool
	}{
		{
			name: "Success",
			server: sqlfake.DatabasesServer{
				Get: func(_ context.Context, _ string, _ string, dbName string, _ *armsql.DatabasesClientGetOptions) (resp azfake.Responder[armsql.DatabasesClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armsql.DatabasesClientGetResponse{
						Database: armsql.Database{
							Name: to.Ptr(dbName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: sqlfake.DatabasesServer{
				Get: func(_ context.Context, _ string, _ string, _ string, _ *armsql.DatabasesClientGetOptions) (resp azfake.Responder[armsql.DatabasesClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeSQLDatabasesClient(t, &srv)

			db, err := azure.GetSQLDatabaseWithClient(context.Background(), client, "rg", "my-server", "my-db")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-db", *db.Name)
		})
	}
}

// ---------------------------------------------------------------------------
// ListSQLServerDatabasesWithClient tests
// ---------------------------------------------------------------------------

func TestListSQLServerDatabasesWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  sqlfake.DatabasesServer
		want    int
		wantErr bool
	}{
		{
			name: "TwoDatabases",
			server: sqlfake.DatabasesServer{
				NewListByServerPager: func(_ string, _ string, _ *armsql.DatabasesClientListByServerOptions) (resp azfake.PagerResponder[armsql.DatabasesClientListByServerResponse]) {
					resp.AddPage(http.StatusOK, armsql.DatabasesClientListByServerResponse{
						DatabaseListResult: armsql.DatabaseListResult{
							Value: []*armsql.Database{
								{Name: to.Ptr("db1")},
								{Name: to.Ptr("db2")},
							},
						},
					}, nil)

					return
				},
			},
			want: 2,
		},
		{
			name: "Empty",
			server: sqlfake.DatabasesServer{
				NewListByServerPager: func(_ string, _ string, _ *armsql.DatabasesClientListByServerOptions) (resp azfake.PagerResponder[armsql.DatabasesClientListByServerResponse]) {
					resp.AddPage(http.StatusOK, armsql.DatabasesClientListByServerResponse{
						DatabaseListResult: armsql.DatabaseListResult{
							Value: []*armsql.Database{},
						},
					}, nil)

					return
				},
			},
			want: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := tc.server
			client := newFakeSQLDatabasesClient(t, &srv)

			dbs, err := azure.ListSQLServerDatabasesWithClient(context.Background(), client, "rg", "my-server")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, dbs, tc.want)
		})
	}
}
