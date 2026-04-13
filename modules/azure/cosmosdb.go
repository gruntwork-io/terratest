package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetCosmosDBAccountClientContextE is a helper function that will setup a CosmosDB account client.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBAccountClientContextE(ctx context.Context, subscriptionID string) (*armcosmos.DatabaseAccountsClient, error) {
	return CreateCosmosDBAccountClientContextE(ctx, subscriptionID)
}

// GetCosmosDBAccountClientE is a helper function that will setup a CosmosDB account client.
//
// Deprecated: Use [GetCosmosDBAccountClientContextE] instead.
func GetCosmosDBAccountClientE(subscriptionID string) (*armcosmos.DatabaseAccountsClient, error) {
	return GetCosmosDBAccountClientContextE(context.Background(), subscriptionID)
}

// GetCosmosDBAccountClientContext is a helper function that will setup a CosmosDB account client.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBAccountClientContext(t testing.TestingT, ctx context.Context, subscriptionID string) *armcosmos.DatabaseAccountsClient {
	cosmosDBAccount, err := GetCosmosDBAccountClientContextE(ctx, subscriptionID)
	require.NoError(t, err)

	return cosmosDBAccount
}

// GetCosmosDBAccountClient is a helper function that will setup a CosmosDB account client.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetCosmosDBAccountClientContext] instead.
func GetCosmosDBAccountClient(t testing.TestingT, subscriptionID string) *armcosmos.DatabaseAccountsClient {
	return GetCosmosDBAccountClientContext(t, context.Background(), subscriptionID) //nolint:staticcheck
}

// GetCosmosDBAccountContext is a helper function that gets the database account.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBAccountContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroupName string, accountName string) *armcosmos.DatabaseAccountGetResults {
	cosmosDBAccount, err := GetCosmosDBAccountContextE(ctx, subscriptionID, resourceGroupName, accountName)
	require.NoError(t, err)

	return cosmosDBAccount
}

// GetCosmosDBAccountContextE is a helper function that gets the database account.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBAccountContextE(ctx context.Context, subscriptionID string, resourceGroupName string, accountName string) (*armcosmos.DatabaseAccountGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBAccountClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding database account
	resp, err := cosmosClient.Get(ctx, resourceGroupName, accountName, nil)
	if err != nil {
		return nil, err
	}

	// Return DB
	return &resp.DatabaseAccountGetResults, nil
}

// GetCosmosDBSQLClientContextE is a helper function that will setup a CosmosDB SQL client.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLClientContextE(ctx context.Context, subscriptionID string) (*armcosmos.SQLResourcesClient, error) {
	return CreateCosmosDBSQLClientContextE(ctx, subscriptionID)
}

// GetCosmosDBSQLClientE is a helper function that will setup a CosmosDB SQL client.
//
// Deprecated: Use [GetCosmosDBSQLClientContextE] instead.
func GetCosmosDBSQLClientE(subscriptionID string) (*armcosmos.SQLResourcesClient, error) {
	return GetCosmosDBSQLClientContextE(context.Background(), subscriptionID)
}

// GetCosmosDBSQLClientContext is a helper function that will setup a CosmosDB SQL client.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLClientContext(t testing.TestingT, ctx context.Context, subscriptionID string) *armcosmos.SQLResourcesClient {
	cosmosClient, err := GetCosmosDBSQLClientContextE(ctx, subscriptionID)
	require.NoError(t, err)

	return cosmosClient
}

// GetCosmosDBSQLClient is a helper function that will setup a CosmosDB SQL client.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetCosmosDBSQLClientContext] instead.
func GetCosmosDBSQLClient(t testing.TestingT, subscriptionID string) *armcosmos.SQLResourcesClient {
	return GetCosmosDBSQLClientContext(t, context.Background(), subscriptionID) //nolint:staticcheck
}

// GetCosmosDBSQLDatabaseContext is a helper function that gets a SQL database.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLDatabaseContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string) *armcosmos.SQLDatabaseGetResults {
	cosmosSQLDB, err := GetCosmosDBSQLDatabaseContextE(ctx, subscriptionID, resourceGroupName, accountName, databaseName)
	require.NoError(t, err)

	return cosmosSQLDB
}

// GetCosmosDBSQLDatabaseContextE is a helper function that gets a SQL database.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLDatabaseContextE(ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string) (*armcosmos.SQLDatabaseGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBSQLClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding database
	resp, err := cosmosClient.GetSQLDatabase(ctx, resourceGroupName, accountName, databaseName, nil)
	if err != nil {
		return nil, err
	}

	// Return DB
	return &resp.SQLDatabaseGetResults, nil
}

// GetCosmosDBSQLContainerContext is a helper function that gets a SQL container.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLContainerContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) *armcosmos.SQLContainerGetResults {
	cosmosSQLContainer, err := GetCosmosDBSQLContainerContextE(ctx, subscriptionID, resourceGroupName, accountName, databaseName, containerName)
	require.NoError(t, err)

	return cosmosSQLContainer
}

// GetCosmosDBSQLContainerContextE is a helper function that gets a SQL container.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLContainerContextE(ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) (*armcosmos.SQLContainerGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBSQLClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding SQL container
	resp, err := cosmosClient.GetSQLContainer(ctx, resourceGroupName, accountName, databaseName, containerName, nil)
	if err != nil {
		return nil, err
	}

	// Return container
	return &resp.SQLContainerGetResults, nil
}

// GetCosmosDBSQLDatabaseThroughputContext is a helper function that gets a SQL database throughput configuration.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLDatabaseThroughputContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string) *armcosmos.ThroughputSettingsGetResults {
	cosmosSQLDBThroughput, err := GetCosmosDBSQLDatabaseThroughputContextE(ctx, subscriptionID, resourceGroupName, accountName, databaseName)
	require.NoError(t, err)

	return cosmosSQLDBThroughput
}

// GetCosmosDBSQLDatabaseThroughputContextE is a helper function that gets a SQL database throughput configuration.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLDatabaseThroughputContextE(ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string) (*armcosmos.ThroughputSettingsGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBSQLClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding database throughput config
	resp, err := cosmosClient.GetSQLDatabaseThroughput(ctx, resourceGroupName, accountName, databaseName, nil)
	if err != nil {
		return nil, err
	}

	// Return throughput config
	return &resp.ThroughputSettingsGetResults, nil
}

// GetCosmosDBSQLContainerThroughputContext is a helper function that gets a SQL container throughput configuration.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLContainerThroughputContext(t testing.TestingT, ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) *armcosmos.ThroughputSettingsGetResults {
	cosmosSQLCtrThroughput, err := GetCosmosDBSQLContainerThroughputContextE(ctx, subscriptionID, resourceGroupName, accountName, databaseName, containerName)
	require.NoError(t, err)

	return cosmosSQLCtrThroughput
}

// GetCosmosDBSQLContainerThroughputContextE is a helper function that gets a SQL container throughput configuration.
// The ctx parameter supports cancellation and timeouts.
func GetCosmosDBSQLContainerThroughputContextE(ctx context.Context, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) (*armcosmos.ThroughputSettingsGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBSQLClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding container throughput config
	resp, err := cosmosClient.GetSQLContainerThroughput(ctx, resourceGroupName, accountName, databaseName, containerName, nil)
	if err != nil {
		return nil, err
	}

	// Return throughput config
	return &resp.ThroughputSettingsGetResults, nil
}
