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

// GetPostgresqlServer is a helper function that gets the server.
// This function would fail the test if there is an error.
func GetPostgresqlServer(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) *postgresql.Server {
	postgresqlServer, err := GetPostgresqlServerE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return postgresqlServer
}

// GetPostgresqlServerE is a helper function that gets the server.
func GetPostgresqlServerE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (*postgresql.Server, error) {
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

// GetPostgresqlDBClientE is a helper function that will setup a postgresql DB client.
func GetPostgresqlDBClientE(subscriptionID string) (*postgresql.DatabasesClient, error) {
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

//GetPostgresqlDB is a helper function that gets the database.
// This function would fail the test if there is an error.
func GetPostgresqlDB(t testing.TestingT, resGroupName string, serverName string, dbName string, subscriptionID string) *postgresql.Database {
	database, err := GetPostgresqlDBE(t, subscriptionID, resGroupName, serverName, dbName)
	require.NoError(t, err)

	return database
}

//GetPostgresqlDBE is a helper function that gets the database.
func GetPostgresqlDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string, dbName string) (*postgresql.Database, error) {
	// Create a postgresql db client
	postgresqldbClient, err := GetPostgresqlDBClientE(subscriptionID)
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

//ListPostgresqlDB is a helper function that gets all databases per server.
func ListPostgresqlDB(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) []postgresql.Database {
	dblist, err := ListPostgresqlDBE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return dblist
}

//ListPostgresqlDBE is a helper function that gets all databases per server.
func ListPostgresqlDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) ([]postgresql.Database, error) {
	// Create a postgresql db client
	postgresqldbClient, err := GetPostgresqlDBClientE(subscriptionID)
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
