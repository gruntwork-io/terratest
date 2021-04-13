package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/sql/mgmt/sql"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetSQLServerClient is a helper function that will setup a sql server client
// TODO: remove in next version
func GetSQLServerClient(subscriptionID string) (*sql.ServersClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a sql server client
	sqlClient := sql.NewServersClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	sqlClient.Authorizer = *authorizer

	return &sqlClient, nil
}

// GetSQLServer is a helper function that gets the sql server object.
// This function would fail the test if there is an error.
func GetSQLServer(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) *sql.Server {
	sqlServer, err := GetSQLServerE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return sqlServer
}

// GetSQLServerE is a helper function that gets the sql server object.
func GetSQLServerE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (*sql.Server, error) {
	// Create a SQl Server client
	sqlClient, err := CreateSQLServerClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	//Get the corresponding server client
	sqlServer, err := sqlClient.Get(context.Background(), resGroupName, serverName)
	if err != nil {
		return nil, err
	}

	//Return sql server
	return &sqlServer, nil
}

// GetDatabaseClient  is a helper function that will setup a sql DB client
// TODO: remove in next version
func GetDatabaseClient(subscriptionID string) (*sql.DatabasesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a sql DB client
	sqlDBClient := sql.NewDatabasesClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	sqlDBClient.Authorizer = *authorizer

	return &sqlDBClient, nil
}

//ListSQLServerDatabases is a helper function that gets a list of databases on a sql server
func ListSQLServerDatabases(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) *[]sql.Database {
	dbList, err := ListSQLServerDatabasesE(t, resGroupName, serverName, subscriptionID)
	require.NoError(t, err)

	return dbList
}

//ListSQLServerDatabasesE is a helper function that gets a list of databases on a sql server
func ListSQLServerDatabasesE(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) (*[]sql.Database, error) {
	// Create a SQl db client
	sqlClient, err := CreateDatabaseClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	//Get the corresponding DB client
	sqlDbs, err := sqlClient.ListByServer(context.Background(), resGroupName, serverName, "", "")
	if err != nil {
		return nil, err
	}

	// Return DB ID
	return sqlDbs.Value, nil
}

// GetSQLDatabase is a helper function that gets the sql db.
// This function would fail the test if there is an error.
func GetSQLDatabase(t testing.TestingT, resGroupName string, serverName string, dbName string, subscriptionID string) *sql.Database {
	database, err := GetSQLDatabaseE(t, subscriptionID, resGroupName, serverName, dbName)
	require.NoError(t, err)

	return database
}

// GetSQLDatabaseE is a helper function that gets the sql db.
func GetSQLDatabaseE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string, dbName string) (*sql.Database, error) {
	// Create a SQl db client
	sqlClient, err := CreateDatabaseClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	//Get the corresponding DB client
	sqlDb, err := sqlClient.Get(context.Background(), resGroupName, serverName, dbName, "")
	if err != nil {
		return nil, err
	}

	// Return DB
	return &sqlDb, nil
}
