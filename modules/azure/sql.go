package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/gruntwork-io/terratest/modules/testing/v2"
	"github.com/stretchr/testify/require"
)

// GetSQLServerClient is a helper function that will setup a sql server client
func GetSQLServerClient(subscriptionID string) (*armsql.ServersClient, error) {
	return CreateSQLServerClient(subscriptionID)
}

// GetSQLServer is a helper function that gets the sql server object.
// This function would fail the test if there is an error.
func GetSQLServer(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) *armsql.Server {
	sqlServer, err := GetSQLServerE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return sqlServer
}

// GetSQLServerE is a helper function that gets the sql server object.
func GetSQLServerE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (*armsql.Server, error) {
	// Create a SQL Server client
	sqlClient, err := CreateSQLServerClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding server
	resp, err := sqlClient.Get(context.Background(), resGroupName, serverName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Server, nil
}

// GetDatabaseClient is a helper function that will setup a sql DB client
func GetDatabaseClient(subscriptionID string) (*armsql.DatabasesClient, error) {
	return CreateDatabaseClient(subscriptionID)
}

// ListSQLServerDatabases is a helper function that gets a list of databases on a sql server
func ListSQLServerDatabases(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) []*armsql.Database {
	dbList, err := ListSQLServerDatabasesE(t, resGroupName, serverName, subscriptionID)
	require.NoError(t, err)

	return dbList
}

// ListSQLServerDatabasesE is a helper function that gets a list of databases on a sql server
func ListSQLServerDatabasesE(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) ([]*armsql.Database, error) {
	// Create a SQL db client
	sqlClient, err := CreateDatabaseClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the databases using pager
	pager := sqlClient.NewListByServerPager(resGroupName, serverName, nil)
	var databases []*armsql.Database
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		databases = append(databases, page.Value...)
	}

	return databases, nil
}

// GetSQLDatabase is a helper function that gets the sql db.
// This function would fail the test if there is an error.
func GetSQLDatabase(t testing.TestingT, resGroupName string, serverName string, dbName string, subscriptionID string) *armsql.Database {
	database, err := GetSQLDatabaseE(t, subscriptionID, resGroupName, serverName, dbName)
	require.NoError(t, err)

	return database
}

// GetSQLDatabaseE is a helper function that gets the sql db.
func GetSQLDatabaseE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string, dbName string) (*armsql.Database, error) {
	// Create a SQL db client
	sqlClient, err := CreateDatabaseClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding DB
	resp, err := sqlClient.Get(context.Background(), resGroupName, serverName, dbName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Database, nil
}
