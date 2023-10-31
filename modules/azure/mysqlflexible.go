package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/mysql/mgmt/mysqlflexibleservers"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetMySqlFlexibleServer is a helper function that gets the server.
// This function would fail the test if there is an error.
func GetMySqlFlexibleServer(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) *mysqlflexibleservers.Server {
	mysqlServer, err := GetMySqlFlexibleServerE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return mysqlServer
}

// GetMySqlFlexibleServerE is a helper function that gets the server.
func GetMySqlFlexibleServerE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (*mysqlflexibleservers.Server, error) {
	// Create a MySql Server client
	flexibleMySqlClient, err := CreateMySqlFlexibleServerClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding flexible server client
	flexibleMySqlServer, err := flexibleMySqlClient.Get(context.Background(), resGroupName, serverName)
	if err != nil {
		return nil, err
	}

	//Return flexible server
	return &flexibleMySqlServer, nil
}

// GetMySqlFlexibleDBClientE is a helper function that will setup a MySql flexible DB client.
func GetMySqlFlexibleServerDBClientE(subscriptionID string) (*mysqlflexibleservers.DatabasesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a mysql db client
	mysqlFlexibleDBClient := mysqlflexibleservers.NewDatabasesClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	mysqlFlexibleDBClient.Authorizer = *authorizer

	return &mysqlFlexibleDBClient, nil
}

// GetMySqlFlexibleServerDB is a helper function that gets the database.
// This function would fail the test if there is an error.
func GetMySqlFlexibleServerDB(t testing.TestingT, resGroupName string, serverName string, dbName string, subscriptionID string) *mysqlflexibleservers.Database {
	database, err := GetMySqlFlexibleServerDBE(t, subscriptionID, resGroupName, serverName, dbName)
	require.NoError(t, err)

	return database
}

// GetMySqlFlexibleServerDBE is a helper function that gets the database.
func GetMySqlFlexibleServerDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string, dbName string) (*mysqlflexibleservers.Database, error) {
	// Create a MySql Flexible Server DB client
	mysqlFlexibleServerDBClient, err := CreateMySqlFlexibleServerDBClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding MySql Flexible Server DB client
	mysqlFlexibleServerDB, err := mysqlFlexibleServerDBClient.Get(context.Background(), resGroupName, serverName, dbName)
	if err != nil {
		return nil, err
	}

	//Return MySql Flexible Server DB
	return &mysqlFlexibleServerDB, nil
}

// ListMySqlFlexibleServerDB is a helper function that gets all databases per server.
func ListMySqlFlexibleServerDB(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) []mysqlflexibleservers.Database {
	mysqlFlexibleServerDBList, err := ListMySqlFlexibleServerDBE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return mysqlFlexibleServerDBList
}

// ListMySqlFlexibleServerDBE is a helper function that gets all databases per server.
func ListMySqlFlexibleServerDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) ([]mysqlflexibleservers.Database, error) {
	// Create a MySql Flexible Server DB client
	mysqlFlexibleServerDBClient, err := CreateMySqlFlexibleServerDBClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding MySql Flexible Server DB client
	mysqlFlexibleServerDBs, err := mysqlFlexibleServerDBClient.ListByServer(context.Background(), resGroupName, serverName)
	if err != nil {
		return nil, err
	}

	//Return MySql Flexible Server DB list
	return mysqlFlexibleServerDBs.Values(), nil
}
