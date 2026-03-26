package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/cosmos-db/mgmt/documentdb"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetCosmosDBAccountClientE is a helper function that will setup a CosmosDB account client.
func GetCosmosDBAccountClientE(subscriptionID string) (*documentdb.DatabaseAccountsClient, error) {
	// Create a CosmosDB client
	cosmosClient, err := CreateCosmosDBAccountClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	cosmosClient.Authorizer = *authorizer

	return cosmosClient, nil
}

// GetCosmosDBAccountClient is a helper function that will setup a CosmosDB account client.
// This function would fail the test if there is an error.
func GetCosmosDBAccountClient(t testing.TestingT, subscriptionID string) *documentdb.DatabaseAccountsClient {
	cosmosDBAccount, err := GetCosmosDBAccountClientE(subscriptionID)
	require.NoError(t, err)

	return cosmosDBAccount
}

// GetCosmosDBAccountContext is a helper function that gets the database account.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBAccountContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroupName string, accountName string) *documentdb.DatabaseAccountGetResults {
	cosmosDBAccount, err := GetCosmosDBAccountContextE(ctx, subscriptionID, resourceGroupName, accountName)
	require.NoError(t, err)

	return cosmosDBAccount
}

// GetCosmosDBAccountContextE is a helper function that gets the database account.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBAccountContextE(ctx context.Context, subscriptionID string, resourceGroupName string, accountName string) (*documentdb.DatabaseAccountGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBAccountClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding database account
	cosmosDBAccount, err := cosmosClient.Get(ctx, resourceGroupName, accountName)
	if err != nil {
		return nil, err
	}

	// Return DB
	return &cosmosDBAccount, nil
}

// GetCosmosDBAccount is a helper function that gets the database account.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetCosmosDBAccountContext] instead.
func GetCosmosDBAccount(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string) *documentdb.DatabaseAccountGetResults {
	return GetCosmosDBAccountContext(t, context.Background(), subscriptionID, resourceGroupName, accountName) //nolint:staticcheck
}

// GetCosmosDBAccountE is a helper function that gets the database account.
//
// Deprecated: Use [GetCosmosDBAccountContextE] instead.
func GetCosmosDBAccountE(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string) (*documentdb.DatabaseAccountGetResults, error) {
	return GetCosmosDBAccountContextE(context.Background(), subscriptionID, resourceGroupName, accountName)
}

// GetCosmosDBSQLClientE is a helper function that will setup a CosmosDB SQL client.
func GetCosmosDBSQLClientE(subscriptionID string) (*documentdb.SQLResourcesClient, error) {
	// Create a CosmosDB client
	cosmosClient, err := CreateCosmosDBSQLClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	cosmosClient.Authorizer = *authorizer

	return cosmosClient, nil
}

// GetCosmosDBSQLClient is a helper function that will setup a CosmosDB SQL client.
// This function would fail the test if there is an error.
func GetCosmosDBSQLClient(t testing.TestingT, subscriptionID string) *documentdb.SQLResourcesClient {
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
	require.NoError(t, err)

	return cosmosClient
}

// GetCosmosDBSQLDatabaseContext is a helper function that gets a SQL database.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLDatabaseContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string) *documentdb.SQLDatabaseGetResults {
	cosmosSQLDB, err := GetCosmosDBSQLDatabaseContextE(ctx, subscriptionID, resourceGroupName, accountName, databaseName)
	require.NoError(t, err)

	return cosmosSQLDB
}

// GetCosmosDBSQLDatabaseContextE is a helper function that gets a SQL database.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLDatabaseContextE(ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string) (*documentdb.SQLDatabaseGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding database
	cosmosSQLDB, err := cosmosClient.GetSQLDatabase(ctx, resourceGroupName, accountName, databaseName)
	if err != nil {
		return nil, err
	}

	// Return DB
	return &cosmosSQLDB, nil
}

// GetCosmosDBSQLDatabase is a helper function that gets a SQL database.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetCosmosDBSQLDatabaseContext] instead.
func GetCosmosDBSQLDatabase(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string) *documentdb.SQLDatabaseGetResults {
	return GetCosmosDBSQLDatabaseContext(t, context.Background(), subscriptionID, resourceGroupName, accountName, databaseName) //nolint:staticcheck
}

// GetCosmosDBSQLDatabaseE is a helper function that gets a SQL database.
//
// Deprecated: Use [GetCosmosDBSQLDatabaseContextE] instead.
func GetCosmosDBSQLDatabaseE(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string) (*documentdb.SQLDatabaseGetResults, error) {
	return GetCosmosDBSQLDatabaseContextE(context.Background(), subscriptionID, resourceGroupName, accountName, databaseName)
}

// GetCosmosDBSQLContainerContext is a helper function that gets a SQL container.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLContainerContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) *documentdb.SQLContainerGetResults {
	cosmosSQLContainer, err := GetCosmosDBSQLContainerContextE(ctx, subscriptionID, resourceGroupName, accountName, databaseName, containerName)
	require.NoError(t, err)

	return cosmosSQLContainer
}

// GetCosmosDBSQLContainerContextE is a helper function that gets a SQL container.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLContainerContextE(ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) (*documentdb.SQLContainerGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding SQL container
	cosmosSQLContainer, err := cosmosClient.GetSQLContainer(ctx, resourceGroupName, accountName, databaseName, containerName)
	if err != nil {
		return nil, err
	}

	// Return container
	return &cosmosSQLContainer, nil
}

// GetCosmosDBSQLContainer is a helper function that gets a SQL container.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetCosmosDBSQLContainerContext] instead.
func GetCosmosDBSQLContainer(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) *documentdb.SQLContainerGetResults {
	return GetCosmosDBSQLContainerContext(t, context.Background(), subscriptionID, resourceGroupName, accountName, databaseName, containerName) //nolint:staticcheck
}

// GetCosmosDBSQLContainerE is a helper function that gets a SQL container.
//
// Deprecated: Use [GetCosmosDBSQLContainerContextE] instead.
func GetCosmosDBSQLContainerE(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) (*documentdb.SQLContainerGetResults, error) {
	return GetCosmosDBSQLContainerContextE(context.Background(), subscriptionID, resourceGroupName, accountName, databaseName, containerName)
}

// GetCosmosDBSQLDatabaseThroughputContext is a helper function that gets a SQL database throughput configuration.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLDatabaseThroughputContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string) *documentdb.ThroughputSettingsGetResults {
	cosmosSQLDBThroughput, err := GetCosmosDBSQLDatabaseThroughputContextE(ctx, subscriptionID, resourceGroupName, accountName, databaseName)
	require.NoError(t, err)

	return cosmosSQLDBThroughput
}

// GetCosmosDBSQLDatabaseThroughputContextE is a helper function that gets a SQL database throughput configuration.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLDatabaseThroughputContextE(ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string) (*documentdb.ThroughputSettingsGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding database throughput config
	cosmosSQLDBThroughput, err := cosmosClient.GetSQLDatabaseThroughput(ctx, resourceGroupName, accountName, databaseName)
	if err != nil {
		return nil, err
	}

	// Return throughput config
	return &cosmosSQLDBThroughput, nil
}

// GetCosmosDBSQLDatabaseThroughput is a helper function that gets a SQL database throughput configuration.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetCosmosDBSQLDatabaseThroughputContext] instead.
func GetCosmosDBSQLDatabaseThroughput(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string) *documentdb.ThroughputSettingsGetResults {
	return GetCosmosDBSQLDatabaseThroughputContext(t, context.Background(), subscriptionID, resourceGroupName, accountName, databaseName) //nolint:staticcheck
}

// GetCosmosDBSQLDatabaseThroughputE is a helper function that gets a SQL database throughput configuration.
//
// Deprecated: Use [GetCosmosDBSQLDatabaseThroughputContextE] instead.
func GetCosmosDBSQLDatabaseThroughputE(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string) (*documentdb.ThroughputSettingsGetResults, error) {
	return GetCosmosDBSQLDatabaseThroughputContextE(context.Background(), subscriptionID, resourceGroupName, accountName, databaseName)
}

// GetCosmosDBSQLContainerThroughputContext is a helper function that gets a SQL container throughput configuration.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLContainerThroughputContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) *documentdb.ThroughputSettingsGetResults {
	cosmosSQLCtrThroughput, err := GetCosmosDBSQLContainerThroughputContextE(ctx, subscriptionID, resourceGroupName, accountName, databaseName, containerName)
	require.NoError(t, err)

	return cosmosSQLCtrThroughput
}

// GetCosmosDBSQLContainerThroughputContextE is a helper function that gets a SQL container throughput configuration.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLContainerThroughputContextE(ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) (*documentdb.ThroughputSettingsGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding container throughput config
	cosmosSQLCtrThroughput, err := cosmosClient.GetSQLContainerThroughput(ctx, resourceGroupName, accountName, databaseName, containerName)
	if err != nil {
		return nil, err
	}

	// Return throughput config
	return &cosmosSQLCtrThroughput, nil
}

// GetCosmosDBSQLContainerThroughput is a helper function that gets a SQL container throughput configuration.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetCosmosDBSQLContainerThroughputContext] instead.
func GetCosmosDBSQLContainerThroughput(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) *documentdb.ThroughputSettingsGetResults {
	return GetCosmosDBSQLContainerThroughputContext(t, context.Background(), subscriptionID, resourceGroupName, accountName, databaseName, containerName) //nolint:staticcheck
}

// GetCosmosDBSQLContainerThroughputE is a helper function that gets a SQL container throughput configuration.
//
// Deprecated: Use [GetCosmosDBSQLContainerThroughputContextE] instead.
func GetCosmosDBSQLContainerThroughputE(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) (*documentdb.ThroughputSettingsGetResults, error) {
	return GetCosmosDBSQLContainerThroughputContextE(context.Background(), subscriptionID, resourceGroupName, accountName, databaseName, containerName)
}
