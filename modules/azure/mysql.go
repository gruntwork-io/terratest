package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysql"
	"github.com/gruntwork-io/terratest/modules/testing/v2"
	"github.com/stretchr/testify/require"
)

// GetMYSQLServerClientE is a helper function that will setup a mysql server client.
func GetMYSQLServerClientE(subscriptionID string) (*armmysql.ServersClient, error) {
	clientFactory, err := getArmMySQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}
	return clientFactory.NewServersClient(), nil
}

// GetMYSQLServer is a helper function that gets the server.
// This function would fail the test if there is an error.
func GetMYSQLServer(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) *armmysql.Server {
	mysqlServer, err := GetMYSQLServerE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return mysqlServer
}

// GetMYSQLServerE is a helper function that gets the server.
func GetMYSQLServerE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (*armmysql.Server, error) {
	// Create a MySQL Server client
	mysqlClient, err := CreateMySQLServerClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding server
	resp, err := mysqlClient.Get(context.Background(), resGroupName, serverName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Server, nil
}

// GetMYSQLDBClientE is a helper function that will setup a mysql DB client.
func GetMYSQLDBClientE(subscriptionID string) (*armmysql.DatabasesClient, error) {
	clientFactory, err := getArmMySQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}
	return clientFactory.NewDatabasesClient(), nil
}

// GetMYSQLDB is a helper function that gets the database.
// This function would fail the test if there is an error.
func GetMYSQLDB(t testing.TestingT, resGroupName string, serverName string, dbName string, subscriptionID string) *armmysql.Database {
	database, err := GetMYSQLDBE(t, subscriptionID, resGroupName, serverName, dbName)
	require.NoError(t, err)

	return database
}

// GetMYSQLDBE is a helper function that gets the database.
func GetMYSQLDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string, dbName string) (*armmysql.Database, error) {
	// Create a MySQL db client
	mysqldbClient, err := GetMYSQLDBClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding db
	resp, err := mysqldbClient.Get(context.Background(), resGroupName, serverName, dbName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Database, nil
}

// ListMySQLDB is a helper function that gets all databases per server.
func ListMySQLDB(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) []*armmysql.Database {
	dblist, err := ListMySQLDBE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return dblist
}

// ListMySQLDBE is a helper function that gets all databases per server.
func ListMySQLDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) ([]*armmysql.Database, error) {
	// Create a MySQL db client
	mysqldbClient, err := GetMYSQLDBClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the databases using pager
	pager := mysqldbClient.NewListByServerPager(resGroupName, serverName, nil)
	var databases []*armmysql.Database
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		databases = append(databases, page.Value...)
	}

	return databases, nil
}
