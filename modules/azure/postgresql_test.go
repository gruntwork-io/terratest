package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql"
	postgresqlfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Fake client helpers
// ---------------------------------------------------------------------------

func newFakePostgreSQLServersClient(t *testing.T, srv *postgresqlfake.ServersServer) *armpostgresql.ServersClient {
	t.Helper()

	client, err := armpostgresql.NewServersClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: postgresqlfake.NewServersServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakePostgreSQLDatabasesClient(t *testing.T, srv *postgresqlfake.DatabasesServer) *armpostgresql.DatabasesClient {
	t.Helper()

	client, err := armpostgresql.NewDatabasesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: postgresqlfake.NewDatabasesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetPostgreSQLServerWithClient tests
// ---------------------------------------------------------------------------

func TestGetPostgreSQLServerWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  postgresqlfake.ServersServer
		wantErr bool
	}{
		{
			name: "Success",
			server: postgresqlfake.ServersServer{
				Get: func(_ context.Context, _ string, serverName string, _ *armpostgresql.ServersClientGetOptions) (resp azfake.Responder[armpostgresql.ServersClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armpostgresql.ServersClientGetResponse{
						Server: armpostgresql.Server{
							Name: to.Ptr(serverName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: postgresqlfake.ServersServer{
				Get: func(_ context.Context, _ string, _ string, _ *armpostgresql.ServersClientGetOptions) (resp azfake.Responder[armpostgresql.ServersClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakePostgreSQLServersClient(t, &srv)

			server, err := azure.GetPostgreSQLServerWithClient(context.Background(), client, "rg", "my-server")
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
// GetPostgreSQLDBWithClient tests
// ---------------------------------------------------------------------------

func TestGetPostgreSQLDBWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  postgresqlfake.DatabasesServer
		wantErr bool
	}{
		{
			name: "Success",
			server: postgresqlfake.DatabasesServer{
				Get: func(_ context.Context, _ string, _ string, dbName string, _ *armpostgresql.DatabasesClientGetOptions) (resp azfake.Responder[armpostgresql.DatabasesClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armpostgresql.DatabasesClientGetResponse{
						Database: armpostgresql.Database{
							Name: to.Ptr(dbName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: postgresqlfake.DatabasesServer{
				Get: func(_ context.Context, _ string, _ string, _ string, _ *armpostgresql.DatabasesClientGetOptions) (resp azfake.Responder[armpostgresql.DatabasesClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakePostgreSQLDatabasesClient(t, &srv)

			db, err := azure.GetPostgreSQLDBWithClient(context.Background(), client, "rg", "my-server", "my-db")
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
// ListPostgreSQLDBWithClient tests
// ---------------------------------------------------------------------------

func TestListPostgreSQLDBWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name    string
		server  postgresqlfake.DatabasesServer
		want    int
		wantErr bool
	}{
		{
			name: "TwoDatabases",
			server: postgresqlfake.DatabasesServer{
				NewListByServerPager: func(_ string, _ string, _ *armpostgresql.DatabasesClientListByServerOptions) (resp azfake.PagerResponder[armpostgresql.DatabasesClientListByServerResponse]) {
					resp.AddPage(http.StatusOK, armpostgresql.DatabasesClientListByServerResponse{
						DatabaseListResult: armpostgresql.DatabaseListResult{
							Value: []*armpostgresql.Database{
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
			server: postgresqlfake.DatabasesServer{
				NewListByServerPager: func(_ string, _ string, _ *armpostgresql.DatabasesClientListByServerOptions) (resp azfake.PagerResponder[armpostgresql.DatabasesClientListByServerResponse]) {
					resp.AddPage(http.StatusOK, armpostgresql.DatabasesClientListByServerResponse{
						DatabaseListResult: armpostgresql.DatabaseListResult{
							Value: []*armpostgresql.Database{},
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
			client := newFakePostgreSQLDatabasesClient(t, &srv)

			dbs, err := azure.ListPostgreSQLDBWithClient(context.Background(), client, "rg", "my-server")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, dbs, tc.want)
		})
	}
}
