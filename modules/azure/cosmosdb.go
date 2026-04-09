package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetCosmosDBAccountClientE is a helper function that will setup a CosmosDB account client.
func GetCosmosDBAccountClientE(subscriptionID string) (*armcosmos.DatabaseAccountsClient, error) {
	return CreateCosmosDBAccountClientE(subscriptionID)
}

// GetCosmosDBAccountClient is a helper function that will setup a CosmosDB account client.
// This function would fail the test if there is an error.
func GetCosmosDBAccountClient(t testing.TestingT, subscriptionID string) *armcosmos.DatabaseAccountsClient {
	cosmosDBAccount, err := GetCosmosDBAccountClientE(subscriptionID)
	require.NoError(t, err)

	return cosmosDBAccount
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
	cosmosClient, err := GetCosmosDBAccountClientE(subscriptionID)
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

// GetCosmosDBSQLClientE is a helper function that will setup a CosmosDB SQL client.
func GetCosmosDBSQLClientE(subscriptionID string) (*armcosmos.SQLResourcesClient, error) {
	return CreateCosmosDBSQLClientE(subscriptionID)
}

// GetCosmosDBSQLClient is a helper function that will setup a CosmosDB SQL client.
// This function would fail the test if there is an error.
func GetCosmosDBSQLClient(t testing.TestingT, subscriptionID string) *armcosmos.SQLResourcesClient {
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
	require.NoError(t, err)

	return cosmosClient
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
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
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
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
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
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
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
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
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
