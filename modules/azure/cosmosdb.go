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

// GetCosmosDBAccount is a helper function that gets the database account.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetCosmosDBAccountContext] instead.
func GetCosmosDBAccount(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string) *armcosmos.DatabaseAccountGetResults {
	return GetCosmosDBAccountContext(t, context.Background(), subscriptionID, resourceGroupName, accountName)
}

// GetCosmosDBAccountE is a helper function that gets the database account.
//
// Deprecated: Use [GetCosmosDBAccountContextE] instead.
func GetCosmosDBAccountE(subscriptionID string, resourceGroupName string, accountName string) (*armcosmos.DatabaseAccountGetResults, error) {
	return GetCosmosDBAccountContextE(context.Background(), subscriptionID, resourceGroupName, accountName)
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

	return GetCosmosDBAccountWithClient(ctx, cosmosClient, resourceGroupName, accountName)
}

// GetCosmosDBAccountWithClient gets a database account using the provided DatabaseAccountsClient.
func GetCosmosDBAccountWithClient(ctx context.Context, client *armcosmos.DatabaseAccountsClient, resourceGroupName string, accountName string) (*armcosmos.DatabaseAccountGetResults, error) {
	resp, err := client.Get(ctx, resourceGroupName, accountName, nil)
	if err != nil {
		return nil, err
	}

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

// GetCosmosDBSQLDatabase is a helper function that gets a SQL database.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetCosmosDBSQLDatabaseContext] instead.
func GetCosmosDBSQLDatabase(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string) *armcosmos.SQLDatabaseGetResults {
	return GetCosmosDBSQLDatabaseContext(t, context.Background(), subscriptionID, resourceGroupName, accountName, databaseName)
}

// GetCosmosDBSQLDatabaseE is a helper function that gets a SQL database.
//
// Deprecated: Use [GetCosmosDBSQLDatabaseContextE] instead.
func GetCosmosDBSQLDatabaseE(subscriptionID string, resourceGroupName string, accountName string, databaseName string) (*armcosmos.SQLDatabaseGetResults, error) {
	return GetCosmosDBSQLDatabaseContextE(context.Background(), subscriptionID, resourceGroupName, accountName, databaseName)
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

	return GetCosmosDBSQLDatabaseWithClient(ctx, cosmosClient, resourceGroupName, accountName, databaseName)
}

// GetCosmosDBSQLDatabaseWithClient gets a SQL database using the provided SQLResourcesClient.
func GetCosmosDBSQLDatabaseWithClient(ctx context.Context, client *armcosmos.SQLResourcesClient, resourceGroupName string, accountName string, databaseName string) (*armcosmos.SQLDatabaseGetResults, error) {
	resp, err := client.GetSQLDatabase(ctx, resourceGroupName, accountName, databaseName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.SQLDatabaseGetResults, nil
}

// GetCosmosDBSQLContainer is a helper function that gets a SQL container.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetCosmosDBSQLContainerContext] instead.
func GetCosmosDBSQLContainer(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) *armcosmos.SQLContainerGetResults {
	return GetCosmosDBSQLContainerContext(t, context.Background(), subscriptionID, resourceGroupName, accountName, databaseName, containerName)
}

// GetCosmosDBSQLContainerE is a helper function that gets a SQL container.
//
// Deprecated: Use [GetCosmosDBSQLContainerContextE] instead.
func GetCosmosDBSQLContainerE(subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) (*armcosmos.SQLContainerGetResults, error) {
	return GetCosmosDBSQLContainerContextE(context.Background(), subscriptionID, resourceGroupName, accountName, databaseName, containerName)
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

	return GetCosmosDBSQLContainerWithClient(ctx, cosmosClient, resourceGroupName, accountName, databaseName, containerName)
}

// GetCosmosDBSQLContainerWithClient gets a SQL container using the provided SQLResourcesClient.
func GetCosmosDBSQLContainerWithClient(ctx context.Context, client *armcosmos.SQLResourcesClient, resourceGroupName string, accountName string, databaseName string, containerName string) (*armcosmos.SQLContainerGetResults, error) {
	resp, err := client.GetSQLContainer(ctx, resourceGroupName, accountName, databaseName, containerName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.SQLContainerGetResults, nil
}

// GetCosmosDBSQLDatabaseThroughput is a helper function that gets a SQL database throughput configuration.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetCosmosDBSQLDatabaseThroughputContext] instead.
func GetCosmosDBSQLDatabaseThroughput(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string) *armcosmos.ThroughputSettingsGetResults {
	return GetCosmosDBSQLDatabaseThroughputContext(t, context.Background(), subscriptionID, resourceGroupName, accountName, databaseName)
}

// GetCosmosDBSQLDatabaseThroughputE is a helper function that gets a SQL database throughput configuration.
//
// Deprecated: Use [GetCosmosDBSQLDatabaseThroughputContextE] instead.
func GetCosmosDBSQLDatabaseThroughputE(subscriptionID string, resourceGroupName string, accountName string, databaseName string) (*armcosmos.ThroughputSettingsGetResults, error) {
	return GetCosmosDBSQLDatabaseThroughputContextE(context.Background(), subscriptionID, resourceGroupName, accountName, databaseName)
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

	return GetCosmosDBSQLDatabaseThroughputWithClient(ctx, cosmosClient, resourceGroupName, accountName, databaseName)
}

// GetCosmosDBSQLDatabaseThroughputWithClient gets a SQL database throughput configuration
// using the provided SQLResourcesClient.
func GetCosmosDBSQLDatabaseThroughputWithClient(ctx context.Context, client *armcosmos.SQLResourcesClient, resourceGroupName string, accountName string, databaseName string) (*armcosmos.ThroughputSettingsGetResults, error) {
	resp, err := client.GetSQLDatabaseThroughput(ctx, resourceGroupName, accountName, databaseName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ThroughputSettingsGetResults, nil
}

// GetCosmosDBSQLContainerThroughput is a helper function that gets a SQL container throughput configuration.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetCosmosDBSQLContainerThroughputContext] instead.
func GetCosmosDBSQLContainerThroughput(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) *armcosmos.ThroughputSettingsGetResults {
	return GetCosmosDBSQLContainerThroughputContext(t, context.Background(), subscriptionID, resourceGroupName, accountName, databaseName, containerName)
}

// GetCosmosDBSQLContainerThroughputE is a helper function that gets a SQL container throughput configuration.
//
// Deprecated: Use [GetCosmosDBSQLContainerThroughputContextE] instead.
func GetCosmosDBSQLContainerThroughputE(subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) (*armcosmos.ThroughputSettingsGetResults, error) {
	return GetCosmosDBSQLContainerThroughputContextE(context.Background(), subscriptionID, resourceGroupName, accountName, databaseName, containerName)
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

	return GetCosmosDBSQLContainerThroughputWithClient(ctx, cosmosClient, resourceGroupName, accountName, databaseName, containerName)
}

// GetCosmosDBSQLContainerThroughputWithClient gets a SQL container throughput configuration
// using the provided SQLResourcesClient.
func GetCosmosDBSQLContainerThroughputWithClient(ctx context.Context, client *armcosmos.SQLResourcesClient, resourceGroupName string, accountName string, databaseName string, containerName string) (*armcosmos.ThroughputSettingsGetResults, error) {
	resp, err := client.GetSQLContainerThroughput(ctx, resourceGroupName, accountName, databaseName, containerName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ThroughputSettingsGetResults, nil
}
