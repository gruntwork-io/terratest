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

// GetCosmosDBAccountClient is a helper function that will setup a CosmosDB account client. This function would fail the test if there is an error.
func GetCosmosDBAccountClient(t testing.TestingT, subscriptionID string) *documentdb.DatabaseAccountsClient {
	cosmosDBAccount, err := GetCosmosDBAccountClientE(subscriptionID)
	require.NoError(t, err)

	return cosmosDBAccount
}

// GetCosmosDBAccount is a helper function that gets the database account. This function would fail the test if there is an error.
func GetCosmosDBAccount(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string) *documentdb.DatabaseAccountGetResults {
	cosmosDBAccount, err := GetCosmosDBAccountE(t, subscriptionID, resourceGroupName, accountName)
	require.NoError(t, err)

	return cosmosDBAccount
}

// GetCosmosDBAccountE is a helper function that gets the database account.
func GetCosmosDBAccountE(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string) (*documentdb.DatabaseAccountGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBAccountClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding database account
	cosmosDBAccount, err := cosmosClient.Get(context.Background(), resourceGroupName, accountName)
	if err != nil {
		return nil, err
	}

	//Return DB
	return &cosmosDBAccount, nil
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

// GetCosmosDBSQLClient is a helper function that will setup a CosmosDB SQL client. This function would fail the test if there is an error.
func GetCosmosDBSQLClient(t testing.TestingT, subscriptionID string) *documentdb.SQLResourcesClient {
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
	require.NoError(t, err)

	return cosmosClient
}

// GetCosmosDBSQLDatabase is a helper function that gets a SQL database. This function would fail the test if there is an error.
func GetCosmosDBSQLDatabase(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string) *documentdb.SQLDatabaseGetResults {
	cosmosSQLDB, err := GetCosmosDBSQLDatabaseE(t, subscriptionID, resourceGroupName, accountName, databaseName)
	require.NoError(t, err)

	return cosmosSQLDB
}

// GetCosmosDBSQLDatabaseE is a helper function that gets a SQL database.
func GetCosmosDBSQLDatabaseE(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string) (*documentdb.SQLDatabaseGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding database
	cosmosSQLDB, err := cosmosClient.GetSQLDatabase(context.Background(), resourceGroupName, accountName, databaseName)
	if err != nil {
		return nil, err
	}

	//Return DB
	return &cosmosSQLDB, nil
}

// GetCosmosDBSQLContainer is a helper function that gets a SQL container. This function would fail the test if there is an error.
func GetCosmosDBSQLContainer(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) *documentdb.SQLContainerGetResults {
	cosmosSQLContainer, err := GetCosmosDBSQLContainerE(t, subscriptionID, resourceGroupName, accountName, databaseName, containerName)
	require.NoError(t, err)

	return cosmosSQLContainer
}

// GetCosmosDBSQLContainerE is a helper function that gets a SQL container.
func GetCosmosDBSQLContainerE(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) (*documentdb.SQLContainerGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding SQL container
	cosmosSQLContainer, err := cosmosClient.GetSQLContainer(context.Background(), resourceGroupName, accountName, databaseName, containerName)
	if err != nil {
		return nil, err
	}

	//Return container
	return &cosmosSQLContainer, nil
}

// GetCosmosDBSQLDatabaseThroughput is a helper function that gets a SQL database throughput configuration. This function would fail the test if there is an error.
func GetCosmosDBSQLDatabaseThroughput(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string) *documentdb.ThroughputSettingsGetResults {
	cosmosSQLDBThroughput, err := GetCosmosDBSQLDatabaseThroughputE(t, subscriptionID, resourceGroupName, accountName, databaseName)
	require.NoError(t, err)

	return cosmosSQLDBThroughput
}

// GetCosmosDBSQLDatabaseThroughputE is a helper function that gets a SQL database throughput configuration.
func GetCosmosDBSQLDatabaseThroughputE(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string) (*documentdb.ThroughputSettingsGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding database throughput config
	cosmosSQLDBThroughput, err := cosmosClient.GetSQLDatabaseThroughput(context.Background(), resourceGroupName, accountName, databaseName)
	if err != nil {
		return nil, err
	}

	//Return throughput config
	return &cosmosSQLDBThroughput, nil
}

// GetCosmosDBSQLContainerThroughput is a helper function that gets a SQL container throughput configuration. This function would fail the test if there is an error.
func GetCosmosDBSQLContainerThroughput(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) *documentdb.ThroughputSettingsGetResults {
	cosmosSQLCtrThroughput, err := GetCosmosDBSQLContainerThroughputE(t, subscriptionID, resourceGroupName, accountName, databaseName, containerName)
	require.NoError(t, err)

	return cosmosSQLCtrThroughput
}

// GetCosmosDBSQLContainerThroughputE is a helper function that gets a SQL container throughput configuration.
func GetCosmosDBSQLContainerThroughputE(t testing.TestingT, subscriptionID string, resourceGroupName string, accountName string, databaseName string, containerName string) (*documentdb.ThroughputSettingsGetResults, error) {
	// Create a CosmosDB client
	cosmosClient, err := GetCosmosDBSQLClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding container throughput config
	cosmosSQLCtrThroughput, err := cosmosClient.GetSQLContainerThroughput(context.Background(), resourceGroupName, accountName, databaseName, containerName)
	if err != nil {
		return nil, err
	}

	//Return throughput config
	return &cosmosSQLCtrThroughput, nil
}
