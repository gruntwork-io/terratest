package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetPostgreSQLServerClientE is a helper function that will setup a postgresql server client.
func GetPostgreSQLServerClientE(subscriptionID string) (*postgresql.ServersClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a postgresql server client
	postgresqlClient := postgresql.NewServersClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	postgresqlClient.Authorizer = *authorizer

	return &postgresqlClient, nil
}

// GetPostgreSQLServer is a helper function that gets the server.
// This function would fail the test if there is an error.
func GetPostgreSQLServer(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) *postgresql.Server {
	postgresqlServer, err := GetPostgreSQLServerE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return postgresqlServer
}

// GetPostgreSQLServerE is a helper function that gets the server.
func GetPostgreSQLServerE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (*postgresql.Server, error) {
	// Create a postgresql Server client
	postgresqlClient, err := GetPostgreSQLServerClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding server client
	postgresqlServer, err := postgresqlClient.Get(context.Background(), resGroupName, serverName)
	if err != nil {
		return nil, err
	}

	//Return server
	return &postgresqlServer, nil
}

// GetPostgreSQLDBClientE is a helper function that will setup a postgresql DB client.
func GetPostgreSQLDBClientE(subscriptionID string) (*postgresql.DatabasesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a postgresql db client
	postgresqlDBClient := postgresql.NewDatabasesClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	postgresqlDBClient.Authorizer = *authorizer

	return &postgresqlDBClient, nil
}

//GetPostgreSQLDB is a helper function that gets the database.
// This function would fail the test if there is an error.
func GetPostgreSQLDB(t testing.TestingT, resGroupName string, serverName string, dbName string, subscriptionID string) *postgresql.Database {
	database, err := GetPostgreSQLDBE(t, subscriptionID, resGroupName, serverName, dbName)
	require.NoError(t, err)

	return database
}

//GetPostgreSQLDBE is a helper function that gets the database.
func GetPostgreSQLDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string, dbName string) (*postgresql.Database, error) {
	// Create a postgresql db client
	postgresqldbClient, err := GetPostgreSQLDBClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding db client
	postgresqlDb, err := postgresqldbClient.Get(context.Background(), resGroupName, serverName, dbName)
	if err != nil {
		return nil, err
	}

	//Return DB
	return &postgresqlDb, nil
}

//ListPostgreSQLDB is a helper function that gets all databases per server.
func ListPostgreSQLDB(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) []postgresql.Database {
	dblist, err := ListPostgreSQLDBE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return dblist
}

//ListPostgreSQLDBE is a helper function that gets all databases per server.
func ListPostgreSQLDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) ([]postgresql.Database, error) {
	// Create a postgresql db client
	postgresqldbClient, err := GetPostgreSQLDBClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding db client
	postgresqlDbs, err := postgresqldbClient.ListByServer(context.Background(), resGroupName, serverName)
	if err != nil {
		return nil, err
	}

	//Return DB lists
	return *postgresqlDbs.Value, nil
}
