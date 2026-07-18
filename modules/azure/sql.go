package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// GetSQLServerClientContext is a helper function that will setup a sql server client.
// The ctx parameter supports cancellation and timeouts.
func GetSQLServerClientContext(ctx context.Context, subscriptionID string) (*armsql.ServersClient, error) {
	return CreateSQLServerClientContext(ctx, subscriptionID)
}

// GetSQLServerContext is a helper function that gets the sql server object.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSQLServerContext(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string) *armsql.Server {
	t.Helper()

	sqlServer, err := GetSQLServerContextE(t, ctx, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return sqlServer
}

// GetSQLServerContextE is a helper function that gets the sql server object.
// The ctx parameter supports cancellation and timeouts.
func GetSQLServerContextE(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string) (*armsql.Server, error) {

	sqlClient, err := CreateSQLServerClientContext(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetSQLServerWithClient(ctx, sqlClient, resGroupName, serverName)
}

// GetSQLServerWithClient gets the sql server object using the provided ServersClient.
func GetSQLServerWithClient(ctx context.Context, client *armsql.ServersClient, resGroupName string, serverName string) (*armsql.Server, error) {
	resp, err := client.Get(ctx, resGroupName, serverName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Server, nil
}

// GetDatabaseClientContext is a helper function that will setup a sql DB client.
// The ctx parameter supports cancellation and timeouts.
func GetDatabaseClientContext(ctx context.Context, subscriptionID string) (*armsql.DatabasesClient, error) {
	return CreateDatabaseClientContext(ctx, subscriptionID)
}

// ListSQLServerDatabasesContext is a helper function that gets a list of databases on a sql server.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListSQLServerDatabasesContext(t testing.TestingT, ctx context.Context, resGroupName string, serverName string, subscriptionID string) []*armsql.Database {
	t.Helper()

	dbList, err := ListSQLServerDatabasesContextE(t, ctx, resGroupName, serverName, subscriptionID)
	require.NoError(t, err)

	return dbList
}

// ListSQLServerDatabasesContextE is a helper function that gets a list of databases on a sql server.
// The ctx parameter supports cancellation and timeouts.
func ListSQLServerDatabasesContextE(t testing.TestingT, ctx context.Context, resGroupName string, serverName string, subscriptionID string) ([]*armsql.Database, error) {

	sqlClient, err := CreateDatabaseClientContext(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return ListSQLServerDatabasesWithClient(ctx, sqlClient, resGroupName, serverName)
}

// ListSQLServerDatabasesWithClient lists databases on a sql server using the provided DatabasesClient.
func ListSQLServerDatabasesWithClient(ctx context.Context, client *armsql.DatabasesClient, resGroupName string, serverName string) ([]*armsql.Database, error) {
	pager := client.NewListByServerPager(resGroupName, serverName, nil)

	var databases []*armsql.Database

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		databases = append(databases, page.Value...)
	}

	return databases, nil
}

// GetSQLDatabaseContext is a helper function that gets the sql db.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSQLDatabaseContext(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string, dbName string) *armsql.Database {
	t.Helper()

	database, err := GetSQLDatabaseContextE(t, ctx, subscriptionID, resGroupName, serverName, dbName)
	require.NoError(t, err)

	return database
}

// GetSQLDatabaseContextE is a helper function that gets the sql db.
// The ctx parameter supports cancellation and timeouts.
func GetSQLDatabaseContextE(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string, dbName string) (*armsql.Database, error) {

	sqlClient, err := CreateDatabaseClientContext(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetSQLDatabaseWithClient(ctx, sqlClient, resGroupName, serverName, dbName)
}

// GetSQLDatabaseWithClient gets the sql db using the provided DatabasesClient.
func GetSQLDatabaseWithClient(ctx context.Context, client *armsql.DatabasesClient, resGroupName string, serverName string, dbName string) (*armsql.Database, error) {
	resp, err := client.Get(ctx, resGroupName, serverName, dbName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Database, nil
}
