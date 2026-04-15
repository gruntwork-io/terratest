package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysql"
	mysqlfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysql/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Fake client helpers
// ---------------------------------------------------------------------------

func newFakeMySQLServersClient(t *testing.T, srv *mysqlfake.ServersServer) *armmysql.ServersClient {
	t.Helper()

	client, err := armmysql.NewServersClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: mysqlfake.NewServersServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeMySQLDatabasesClient(t *testing.T, srv *mysqlfake.DatabasesServer) *armmysql.DatabasesClient {
	t.Helper()

	client, err := armmysql.NewDatabasesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: mysqlfake.NewDatabasesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetMYSQLServerWithClient tests
// ---------------------------------------------------------------------------

func TestGetMYSQLServerWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  mysqlfake.ServersServer
		wantErr bool
	}{
		{
			name: "Success",
			server: mysqlfake.ServersServer{
				Get: func(_ context.Context, _ string, serverName string, _ *armmysql.ServersClientGetOptions) (resp azfake.Responder[armmysql.ServersClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armmysql.ServersClientGetResponse{
						Server: armmysql.Server{
							Name: to.Ptr(serverName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: mysqlfake.ServersServer{
				Get: func(_ context.Context, _ string, _ string, _ *armmysql.ServersClientGetOptions) (resp azfake.Responder[armmysql.ServersClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeMySQLServersClient(t, &srv)

			server, err := azure.GetMYSQLServerWithClient(context.Background(), client, "rg", "my-server")
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
// GetMYSQLDBWithClient tests
// ---------------------------------------------------------------------------

func TestGetMYSQLDBWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  mysqlfake.DatabasesServer
		wantErr bool
	}{
		{
			name: "Success",
			server: mysqlfake.DatabasesServer{
				Get: func(_ context.Context, _ string, _ string, dbName string, _ *armmysql.DatabasesClientGetOptions) (resp azfake.Responder[armmysql.DatabasesClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armmysql.DatabasesClientGetResponse{
						Database: armmysql.Database{
							Name: to.Ptr(dbName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: mysqlfake.DatabasesServer{
				Get: func(_ context.Context, _ string, _ string, _ string, _ *armmysql.DatabasesClientGetOptions) (resp azfake.Responder[armmysql.DatabasesClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeMySQLDatabasesClient(t, &srv)

			db, err := azure.GetMYSQLDBWithClient(context.Background(), client, "rg", "my-server", "my-db")
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
// ListMySQLDBWithClient tests
// ---------------------------------------------------------------------------

func TestListMySQLDBWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  mysqlfake.DatabasesServer
		want    int
		wantErr bool
	}{
		{
			name: "TwoDatabases",
			server: mysqlfake.DatabasesServer{
				NewListByServerPager: func(_ string, _ string, _ *armmysql.DatabasesClientListByServerOptions) (resp azfake.PagerResponder[armmysql.DatabasesClientListByServerResponse]) {
					resp.AddPage(http.StatusOK, armmysql.DatabasesClientListByServerResponse{
						DatabaseListResult: armmysql.DatabaseListResult{
							Value: []*armmysql.Database{
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
			server: mysqlfake.DatabasesServer{
				NewListByServerPager: func(_ string, _ string, _ *armmysql.DatabasesClientListByServerOptions) (resp azfake.PagerResponder[armmysql.DatabasesClientListByServerResponse]) {
					resp.AddPage(http.StatusOK, armmysql.DatabasesClientListByServerResponse{
						DatabaseListResult: armmysql.DatabaseListResult{
							Value: []*armmysql.Database{},
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
			client := newFakeMySQLDatabasesClient(t, &srv)

			dbs, err := azure.ListMySQLDBWithClient(context.Background(), client, "rg", "my-server")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, dbs, tc.want)
		})
	}
}
