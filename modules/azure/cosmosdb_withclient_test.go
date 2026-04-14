package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3"
	cosmosfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newFakeCosmosDBAccountClient(t *testing.T, srv *cosmosfake.DatabaseAccountsServer) *armcosmos.DatabaseAccountsClient {
	t.Helper()

	transport := cosmosfake.NewDatabaseAccountsServerTransport(srv)
	client, err := armcosmos.NewDatabaseAccountsClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func newFakeCosmosDBSQLClient(t *testing.T, srv *cosmosfake.SQLResourcesServer) *armcosmos.SQLResourcesClient {
	t.Helper()

	transport := cosmosfake.NewSQLResourcesServerTransport(srv)
	client, err := armcosmos.NewSQLResourcesClient("fake-sub", &azfake.TokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	require.NoError(t, err)

	return client
}

func TestGetCosmosDBAccountWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		server  cosmosfake.DatabaseAccountsServer
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			server: cosmosfake.DatabaseAccountsServer{
				Get: func(_ context.Context, _ string, accountName string, _ *armcosmos.DatabaseAccountsClientGetOptions) (resp azfake.Responder[armcosmos.DatabaseAccountsClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armcosmos.DatabaseAccountsClientGetResponse{
						DatabaseAccountGetResults: armcosmos.DatabaseAccountGetResults{
							Name: to.Ptr(accountName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: cosmosfake.DatabaseAccountsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcosmos.DatabaseAccountsClientGetOptions) (resp azfake.Responder[armcosmos.DatabaseAccountsClientGetResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeCosmosDBAccountClient(t, &srv)

			account, err := azure.GetCosmosDBAccountWithClient(t.Context(), client, "rg", "my-cosmos")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-cosmos", *account.Name)
		})
	}
}

func TestGetCosmosDBSQLDatabaseWithClient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		server  cosmosfake.SQLResourcesServer
		name    string
		wantErr bool
	}{
		{
			name: "Success",
			server: cosmosfake.SQLResourcesServer{
				GetSQLDatabase: func(_ context.Context, _ string, _ string, dbName string, _ *armcosmos.SQLResourcesClientGetSQLDatabaseOptions) (resp azfake.Responder[armcosmos.SQLResourcesClientGetSQLDatabaseResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armcosmos.SQLResourcesClientGetSQLDatabaseResponse{
						SQLDatabaseGetResults: armcosmos.SQLDatabaseGetResults{
							Name: to.Ptr(dbName),
						},
					}, nil)

					return
				},
			},
		},
		{
			name: "NotFound",
			server: cosmosfake.SQLResourcesServer{
				GetSQLDatabase: func(_ context.Context, _ string, _ string, _ string, _ *armcosmos.SQLResourcesClientGetSQLDatabaseOptions) (resp azfake.Responder[armcosmos.SQLResourcesClientGetSQLDatabaseResponse], errResp azfake.ErrorResponder) {
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
			client := newFakeCosmosDBSQLClient(t, &srv)

			db, err := azure.GetCosmosDBSQLDatabaseWithClient(t.Context(), client, "rg", "acct", "my-db")
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, "my-db", *db.Name)
		})
	}
}

func TestGetCosmosDBSQLContainerWithClient(t *testing.T) {
	t.Parallel()

	srv := &cosmosfake.SQLResourcesServer{
		GetSQLContainer: func(_ context.Context, _ string, _ string, _ string, containerName string, _ *armcosmos.SQLResourcesClientGetSQLContainerOptions) (resp azfake.Responder[armcosmos.SQLResourcesClientGetSQLContainerResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armcosmos.SQLResourcesClientGetSQLContainerResponse{
				SQLContainerGetResults: armcosmos.SQLContainerGetResults{
					Name: to.Ptr(containerName),
				},
			}, nil)

			return
		},
	}
	client := newFakeCosmosDBSQLClient(t, srv)

	container, err := azure.GetCosmosDBSQLContainerWithClient(t.Context(), client, "rg", "acct", "db", "my-container")
	require.NoError(t, err)
	assert.Equal(t, "my-container", *container.Name)
}

func TestGetCosmosDBSQLDatabaseThroughputWithClient(t *testing.T) {
	t.Parallel()

	srv := &cosmosfake.SQLResourcesServer{
		GetSQLDatabaseThroughput: func(_ context.Context, _ string, _ string, _ string, _ *armcosmos.SQLResourcesClientGetSQLDatabaseThroughputOptions) (resp azfake.Responder[armcosmos.SQLResourcesClientGetSQLDatabaseThroughputResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armcosmos.SQLResourcesClientGetSQLDatabaseThroughputResponse{
				ThroughputSettingsGetResults: armcosmos.ThroughputSettingsGetResults{
					Properties: &armcosmos.ThroughputSettingsGetProperties{
						Resource: &armcosmos.ThroughputSettingsGetPropertiesResource{
							Throughput: to.Ptr[int32](400),
						},
					},
				},
			}, nil)

			return
		},
	}
	client := newFakeCosmosDBSQLClient(t, srv)

	tp, err := azure.GetCosmosDBSQLDatabaseThroughputWithClient(t.Context(), client, "rg", "acct", "db")
	require.NoError(t, err)
	assert.Equal(t, int32(400), *tp.Properties.Resource.Throughput)
}

func TestGetCosmosDBSQLContainerThroughputWithClient(t *testing.T) {
	t.Parallel()

	srv := &cosmosfake.SQLResourcesServer{
		GetSQLContainerThroughput: func(_ context.Context, _ string, _ string, _ string, _ string, _ *armcosmos.SQLResourcesClientGetSQLContainerThroughputOptions) (resp azfake.Responder[armcosmos.SQLResourcesClientGetSQLContainerThroughputResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armcosmos.SQLResourcesClientGetSQLContainerThroughputResponse{
				ThroughputSettingsGetResults: armcosmos.ThroughputSettingsGetResults{
					Properties: &armcosmos.ThroughputSettingsGetProperties{
						Resource: &armcosmos.ThroughputSettingsGetPropertiesResource{
							Throughput: to.Ptr[int32](800),
						},
					},
				},
			}, nil)

			return
		},
	}
	client := newFakeCosmosDBSQLClient(t, srv)

	tp, err := azure.GetCosmosDBSQLContainerThroughputWithClient(t.Context(), client, "rg", "acct", "db", "ctr")
	require.NoError(t, err)
	assert.Equal(t, int32(800), *tp.Properties.Resource.Throughput)
}
