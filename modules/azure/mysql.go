package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/mysql/mgmt/mysql"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetMYSQLServerClientE is a helper function that will setup a mysql server client.
// TODO: remove in next version
func GetMYSQLServerClientE(subscriptionID string) (*mysql.ServersClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a mysql server client
	mysqlClient := mysql.NewServersClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	mysqlClient.Authorizer = *authorizer

	return &mysqlClient, nil
}

// GetMYSQLServer is a helper function that gets the server.
// This function would fail the test if there is an error.
func GetMYSQLServer(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) *mysql.Server {
	mysqlServer, err := GetMYSQLServerE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return mysqlServer
}

// GetMYSQLServerE is a helper function that gets the server.
func GetMYSQLServerE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (*mysql.Server, error) {
	// Create a mySQl Server client
	mysqlClient, err := CreateMySQLServerClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding server client
	mysqlServer, err := mysqlClient.Get(context.Background(), resGroupName, serverName)
	if err != nil {
		return nil, err
	}

	//Return server
	return &mysqlServer, nil
}

// GetMYSQLDBClientE is a helper function that will setup a mysql DB client.
func GetMYSQLDBClientE(subscriptionID string) (*mysql.DatabasesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a mysql db client
	mysqlDBClient := mysql.NewDatabasesClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	mysqlDBClient.Authorizer = *authorizer

	return &mysqlDBClient, nil
}

// GetMYSQLDB is a helper function that gets the database.
// This function would fail the test if there is an error.
func GetMYSQLDB(t testing.TestingT, resGroupName string, serverName string, dbName string, subscriptionID string) *mysql.Database {
	database, err := GetMYSQLDBE(t, subscriptionID, resGroupName, serverName, dbName)
	require.NoError(t, err)

	return database
}

// GetMYSQLDBE is a helper function that gets the database.
func GetMYSQLDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string, dbName string) (*mysql.Database, error) {
	// Create a mySQl db client
	mysqldbClient, err := GetMYSQLDBClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding db client
	mysqlDb, err := mysqldbClient.Get(context.Background(), resGroupName, serverName, dbName)
	if err != nil {
		return nil, err
	}

	//Return DB
	return &mysqlDb, nil
}

// ListMySQLDB is a helper function that gets all databases per server.
func ListMySQLDB(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) []mysql.Database {
	dblist, err := ListMySQLDBE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return dblist
}

// ListMySQLDBE is a helper function that gets all databases per server.
func ListMySQLDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) ([]mysql.Database, error) {
	// Create a mySQl db client
	mysqldbClient, err := GetMYSQLDBClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding db client
	mysqlDbs, err := mysqldbClient.ListByServer(context.Background(), resGroupName, serverName)
	if err != nil {
		return nil, err
	}

	//Return DB lists
	return *mysqlDbs.Value, nil
}
