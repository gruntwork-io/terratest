package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql"
	"github.com/gruntwork-io/terratest/modules/testing/v2"
	"github.com/stretchr/testify/require"
)

// GetPostgreSQLServerClientE is a helper function that will setup a postgresql server client.
func GetPostgreSQLServerClientE(subscriptionID string) (*armpostgresql.ServersClient, error) {
	clientFactory, err := getArmPostgreSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}
	return clientFactory.NewServersClient(), nil
}

// GetPostgreSQLServer is a helper function that gets the server.
// This function would fail the test if there is an error.
func GetPostgreSQLServer(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) *armpostgresql.Server {
	postgresqlServer, err := GetPostgreSQLServerE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return postgresqlServer
}

// GetPostgreSQLServerE is a helper function that gets the server.
func GetPostgreSQLServerE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (*armpostgresql.Server, error) {
	// Create a postgresql Server client
	postgresqlClient, err := GetPostgreSQLServerClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding server
	resp, err := postgresqlClient.Get(context.Background(), resGroupName, serverName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Server, nil
}

// GetPostgreSQLDBClientE is a helper function that will setup a postgresql DB client.
func GetPostgreSQLDBClientE(subscriptionID string) (*armpostgresql.DatabasesClient, error) {
	clientFactory, err := getArmPostgreSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}
	return clientFactory.NewDatabasesClient(), nil
}

// GetPostgreSQLDB is a helper function that gets the database.
// This function would fail the test if there is an error.
func GetPostgreSQLDB(t testing.TestingT, resGroupName string, serverName string, dbName string, subscriptionID string) *armpostgresql.Database {
	database, err := GetPostgreSQLDBE(t, subscriptionID, resGroupName, serverName, dbName)
	require.NoError(t, err)

	return database
}

// GetPostgreSQLDBE is a helper function that gets the database.
func GetPostgreSQLDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string, dbName string) (*armpostgresql.Database, error) {
	// Create a postgresql db client
	postgresqldbClient, err := GetPostgreSQLDBClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding db
	resp, err := postgresqldbClient.Get(context.Background(), resGroupName, serverName, dbName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Database, nil
}

// ListPostgreSQLDB is a helper function that gets all databases per server.
func ListPostgreSQLDB(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) []*armpostgresql.Database {
	dblist, err := ListPostgreSQLDBE(t, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return dblist
}

// ListPostgreSQLDBE is a helper function that gets all databases per server.
func ListPostgreSQLDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) ([]*armpostgresql.Database, error) {
	// Create a postgresql db client
	postgresqldbClient, err := GetPostgreSQLDBClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the databases using pager
	pager := postgresqldbClient.NewListByServerPager(resGroupName, serverName, nil)
	var databases []*armpostgresql.Database
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		databases = append(databases, page.Value...)
	}

	return databases, nil
}
