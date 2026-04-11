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

func newFakeDatabaseAccountsClient(t *testing.T, srv *cosmosfake.DatabaseAccountsServer) *armcosmos.DatabaseAccountsClient {
	t.Helper()

	client, err := armcosmos.NewDatabaseAccountsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: cosmosfake.NewDatabaseAccountsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeSQLResourcesClient(t *testing.T, srv *cosmosfake.SQLResourcesServer) *armcosmos.SQLResourcesClient {
	t.Helper()

	client, err := armcosmos.NewSQLResourcesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: cosmosfake.NewSQLResourcesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetCosmosDBAccount tests (success and not-found)
// ---------------------------------------------------------------------------

func TestGetCosmosDBAccount(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    cosmosfake.DatabaseAccountsServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: cosmosfake.DatabaseAccountsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcosmos.DatabaseAccountsClientGetOptions) (resp azfake.Responder[armcosmos.DatabaseAccountsClientGetResponse], errResp azfake.ErrorResponder) {
					resp.SetResponse(http.StatusOK, armcosmos.DatabaseAccountsClientGetResponse{
						DatabaseAccountGetResults: armcosmos.DatabaseAccountGetResults{
							Name: to.Ptr("test-cosmos-account"),
							Properties: &armcosmos.DatabaseAccountGetProperties{
								DocumentEndpoint: to.Ptr("https://test-cosmos-account.documents.azure.com:443/"),
							},
						},
					}, nil)

					return
				},
			},
			wantName: "test-cosmos-account",
		},
		{
			name: "NotFound",
			server: cosmosfake.DatabaseAccountsServer{
				Get: func(_ context.Context, _ string, _ string, _ *armcosmos.DatabaseAccountsClientGetOptions) (resp azfake.Responder[armcosmos.DatabaseAccountsClientGetResponse], errResp azfake.ErrorResponder) {
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

			client := newFakeDatabaseAccountsClient(t, &tc.server)

			resp, err := client.Get(context.Background(), "rg", "test-cosmos-account", nil)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errSubstr)
				assert.True(t, azure.ResourceNotFoundErrorExists(err))

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *resp.Name)
			assert.NotNil(t, resp.Properties.DocumentEndpoint)
		})
	}
}

// ---------------------------------------------------------------------------
// GetCosmosDBSQLDatabase tests (success)
// ---------------------------------------------------------------------------

func TestGetCosmosDBSQLDatabase(t *testing.T) {
	t.Parallel()

	srv := &cosmosfake.SQLResourcesServer{
		GetSQLDatabase: func(_ context.Context, _ string, _ string, _ string, _ *armcosmos.SQLResourcesClientGetSQLDatabaseOptions) (resp azfake.Responder[armcosmos.SQLResourcesClientGetSQLDatabaseResponse], errResp azfake.ErrorResponder) {
			resp.SetResponse(http.StatusOK, armcosmos.SQLResourcesClientGetSQLDatabaseResponse{
				SQLDatabaseGetResults: armcosmos.SQLDatabaseGetResults{
					Name: to.Ptr("test-db"),
					Properties: &armcosmos.SQLDatabaseGetProperties{
						Resource: &armcosmos.SQLDatabaseGetPropertiesResource{
							ID: to.Ptr("test-db"),
						},
					},
				},
			}, nil)

			return
		},
	}

	client := newFakeSQLResourcesClient(t, srv)

	resp, err := client.GetSQLDatabase(context.Background(), "rg", "test-cosmos-account", "test-db", nil)
	require.NoError(t, err)
	assert.Equal(t, "test-db", *resp.Name)
	assert.Equal(t, "test-db", *resp.Properties.Resource.ID)
}

// ---------------------------------------------------------------------------
// GetCosmosDBSQLDatabase NotFound test
// ---------------------------------------------------------------------------

func TestGetCosmosDBSQLDatabase_NotFound(t *testing.T) {
	t.Parallel()

	srv := &cosmosfake.SQLResourcesServer{
		GetSQLDatabase: func(_ context.Context, _ string, _ string, _ string, _ *armcosmos.SQLResourcesClientGetSQLDatabaseOptions) (resp azfake.Responder[armcosmos.SQLResourcesClientGetSQLDatabaseResponse], errResp azfake.ErrorResponder) {
			errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

			return
		},
	}

	client := newFakeSQLResourcesClient(t, srv)

	_, err := client.GetSQLDatabase(context.Background(), "rg", "test-cosmos-account", "missing-db", nil)
	require.Error(t, err)
	assert.True(t, azure.ResourceNotFoundErrorExists(err))
}
