//nolint:dupl // structural duplication with different Azure SDK types
package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysql"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetMYSQLServerClientContextE is a helper function that will setup a mysql server client.
// The ctx parameter supports cancellation and timeouts.
func GetMYSQLServerClientContextE(_ context.Context, subscriptionID string) (*armmysql.ServersClient, error) {
	clientFactory, err := getArmMySQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewServersClient(), nil
}

// GetMYSQLServerClientE is a helper function that will setup a mysql server client.
//
// Deprecated: Use [GetMYSQLServerClientContextE] instead.
func GetMYSQLServerClientE(subscriptionID string) (*armmysql.ServersClient, error) {
	return GetMYSQLServerClientContextE(context.Background(), subscriptionID)
}

// GetMYSQLServerContext is a helper function that gets the server.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMYSQLServerContext(t testing.TestingT, ctx context.Context, resGroupName string, serverName string, subscriptionID string) *armmysql.Server {
	t.Helper()

	mysqlServer, err := GetMYSQLServerContextE(t, ctx, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return mysqlServer
}

// GetMYSQLServer is a helper function that gets the server.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetMYSQLServerContext] instead.
func GetMYSQLServer(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) *armmysql.Server {
	t.Helper()

	return GetMYSQLServerContext(t, context.Background(), resGroupName, serverName, subscriptionID) //nolint:staticcheck
}

// GetMYSQLServerContextE is a helper function that gets the server.
// The ctx parameter supports cancellation and timeouts.
func GetMYSQLServerContextE(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string) (*armmysql.Server, error) {
	// Create a MySQL Server client
	mysqlClient, err := CreateMySQLServerClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding server
	resp, err := mysqlClient.Get(ctx, resGroupName, serverName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Server, nil
}

// GetMYSQLServerE is a helper function that gets the server.
//
// Deprecated: Use [GetMYSQLServerContextE] instead.
func GetMYSQLServerE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (*armmysql.Server, error) {
	return GetMYSQLServerContextE(t, context.Background(), subscriptionID, resGroupName, serverName)
}

// GetMYSQLDBClientContextE is a helper function that will setup a mysql DB client.
// The ctx parameter supports cancellation and timeouts.
func GetMYSQLDBClientContextE(_ context.Context, subscriptionID string) (*armmysql.DatabasesClient, error) {
	clientFactory, err := getArmMySQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDatabasesClient(), nil
}

// GetMYSQLDBClientE is a helper function that will setup a mysql DB client.
//
// Deprecated: Use [GetMYSQLDBClientContextE] instead.
func GetMYSQLDBClientE(subscriptionID string) (*armmysql.DatabasesClient, error) {
	return GetMYSQLDBClientContextE(context.Background(), subscriptionID)
}

// GetMYSQLDBContext is a helper function that gets the database.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMYSQLDBContext(t testing.TestingT, ctx context.Context, resGroupName string, serverName string, dbName string, subscriptionID string) *armmysql.Database {
	t.Helper()

	database, err := GetMYSQLDBContextE(t, ctx, subscriptionID, resGroupName, serverName, dbName)
	require.NoError(t, err)

	return database
}

// GetMYSQLDB is a helper function that gets the database.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetMYSQLDBContext] instead.
func GetMYSQLDB(t testing.TestingT, resGroupName string, serverName string, dbName string, subscriptionID string) *armmysql.Database {
	t.Helper()

	return GetMYSQLDBContext(t, context.Background(), resGroupName, serverName, dbName, subscriptionID) //nolint:staticcheck
}

// GetMYSQLDBContextE is a helper function that gets the database.
// The ctx parameter supports cancellation and timeouts.
func GetMYSQLDBContextE(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string, dbName string) (*armmysql.Database, error) {
	// Create a MySQL db client
	mysqldbClient, err := GetMYSQLDBClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding db
	resp, err := mysqldbClient.Get(ctx, resGroupName, serverName, dbName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Database, nil
}

// GetMYSQLDBE is a helper function that gets the database.
//
// Deprecated: Use [GetMYSQLDBContextE] instead.
func GetMYSQLDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string, dbName string) (*armmysql.Database, error) {
	return GetMYSQLDBContextE(t, context.Background(), subscriptionID, resGroupName, serverName, dbName)
}

// ListMySQLDBContext is a helper function that gets all databases per server.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListMySQLDBContext(t testing.TestingT, ctx context.Context, resGroupName string, serverName string, subscriptionID string) []*armmysql.Database {
	t.Helper()

	dblist, err := ListMySQLDBContextE(t, ctx, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return dblist
}

// ListMySQLDB is a helper function that gets all databases per server.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListMySQLDBContext] instead.
func ListMySQLDB(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) []*armmysql.Database {
	t.Helper()

	return ListMySQLDBContext(t, context.Background(), resGroupName, serverName, subscriptionID) //nolint:staticcheck
}

// ListMySQLDBContextE is a helper function that gets all databases per server.
// The ctx parameter supports cancellation and timeouts.
func ListMySQLDBContextE(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string) ([]*armmysql.Database, error) {
	// Create a MySQL db client
	mysqldbClient, err := GetMYSQLDBClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the databases using pager
	pager := mysqldbClient.NewListByServerPager(resGroupName, serverName, nil)

	var databases []*armmysql.Database

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		databases = append(databases, page.Value...)
	}

	return databases, nil
}

// ListMySQLDBE is a helper function that gets all databases per server.
//
// Deprecated: Use [ListMySQLDBContextE] instead.
func ListMySQLDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) ([]*armmysql.Database, error) {
	return ListMySQLDBContextE(t, context.Background(), subscriptionID, resGroupName, serverName)
}
