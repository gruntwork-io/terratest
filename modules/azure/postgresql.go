//nolint:dupl // structural duplication with different Azure SDK types
package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/postgresql/armpostgresql"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetPostgreSQLServerClientContextE is a helper function that will setup a postgresql server client.
// The ctx parameter supports cancellation and timeouts.
func GetPostgreSQLServerClientContextE(_ context.Context, subscriptionID string) (*armpostgresql.ServersClient, error) {
	clientFactory, err := getArmPostgreSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewServersClient(), nil
}

// GetPostgreSQLServerClientE is a helper function that will setup a postgresql server client.
//
// Deprecated: Use [GetPostgreSQLServerClientContextE] instead.
func GetPostgreSQLServerClientE(subscriptionID string) (*armpostgresql.ServersClient, error) {
	return GetPostgreSQLServerClientContextE(context.Background(), subscriptionID)
}

// GetPostgreSQLServerContext is a helper function that gets the server.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetPostgreSQLServerContext(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string) *armpostgresql.Server {
	t.Helper()

	postgresqlServer, err := GetPostgreSQLServerContextE(t, ctx, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return postgresqlServer
}

// GetPostgreSQLServer is a helper function that gets the server.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetPostgreSQLServerContext] instead.
func GetPostgreSQLServer(t testing.TestingT, resGroupName string, serverName string, subscriptionID string) *armpostgresql.Server {
	t.Helper()

	return GetPostgreSQLServerContext(t, context.Background(), subscriptionID, resGroupName, serverName) //nolint:staticcheck
}

// GetPostgreSQLServerContextE is a helper function that gets the server.
// The ctx parameter supports cancellation and timeouts.
func GetPostgreSQLServerContextE(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string) (*armpostgresql.Server, error) {
	// Create a postgresql Server client
	postgresqlClient, err := GetPostgreSQLServerClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetPostgreSQLServerWithClient(ctx, postgresqlClient, resGroupName, serverName)
}

// GetPostgreSQLServerWithClient gets the server using the provided ServersClient.
func GetPostgreSQLServerWithClient(ctx context.Context, client *armpostgresql.ServersClient, resGroupName string, serverName string) (*armpostgresql.Server, error) {
	resp, err := client.Get(ctx, resGroupName, serverName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Server, nil
}

// GetPostgreSQLServerE is a helper function that gets the server.
//
// Deprecated: Use [GetPostgreSQLServerContextE] instead.
func GetPostgreSQLServerE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) (*armpostgresql.Server, error) {
	return GetPostgreSQLServerContextE(t, context.Background(), subscriptionID, resGroupName, serverName)
}

// GetPostgreSQLDBClientContextE is a helper function that will setup a postgresql DB client.
// The ctx parameter supports cancellation and timeouts.
func GetPostgreSQLDBClientContextE(_ context.Context, subscriptionID string) (*armpostgresql.DatabasesClient, error) {
	clientFactory, err := getArmPostgreSQLClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewDatabasesClient(), nil
}

// GetPostgreSQLDBClientE is a helper function that will setup a postgresql DB client.
//
// Deprecated: Use [GetPostgreSQLDBClientContextE] instead.
func GetPostgreSQLDBClientE(subscriptionID string) (*armpostgresql.DatabasesClient, error) {
	return GetPostgreSQLDBClientContextE(context.Background(), subscriptionID)
}

// GetPostgreSQLDBContext is a helper function that gets the database.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetPostgreSQLDBContext(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string, dbName string) *armpostgresql.Database {
	t.Helper()

	database, err := GetPostgreSQLDBContextE(t, ctx, subscriptionID, resGroupName, serverName, dbName)
	require.NoError(t, err)

	return database
}

// GetPostgreSQLDB is a helper function that gets the database.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetPostgreSQLDBContext] instead.
func GetPostgreSQLDB(t testing.TestingT, resGroupName string, serverName string, dbName string, subscriptionID string) *armpostgresql.Database {
	t.Helper()

	return GetPostgreSQLDBContext(t, context.Background(), subscriptionID, resGroupName, serverName, dbName) //nolint:staticcheck
}

// GetPostgreSQLDBContextE is a helper function that gets the database.
// The ctx parameter supports cancellation and timeouts.
func GetPostgreSQLDBContextE(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string, dbName string) (*armpostgresql.Database, error) {
	// Create a postgresql db client
	postgresqldbClient, err := GetPostgreSQLDBClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetPostgreSQLDBWithClient(ctx, postgresqldbClient, resGroupName, serverName, dbName)
}

// GetPostgreSQLDBWithClient gets the database using the provided DatabasesClient.
func GetPostgreSQLDBWithClient(ctx context.Context, client *armpostgresql.DatabasesClient, resGroupName string, serverName string, dbName string) (*armpostgresql.Database, error) {
	resp, err := client.Get(ctx, resGroupName, serverName, dbName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Database, nil
}

// GetPostgreSQLDBE is a helper function that gets the database.
//
// Deprecated: Use [GetPostgreSQLDBContextE] instead.
func GetPostgreSQLDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string, dbName string) (*armpostgresql.Database, error) {
	return GetPostgreSQLDBContextE(t, context.Background(), subscriptionID, resGroupName, serverName, dbName)
}

// ListPostgreSQLDBContext is a helper function that gets all databases per server.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListPostgreSQLDBContext(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string) []*armpostgresql.Database {
	t.Helper()

	dblist, err := ListPostgreSQLDBContextE(t, ctx, subscriptionID, resGroupName, serverName)
	require.NoError(t, err)

	return dblist
}

// ListPostgreSQLDB is a helper function that gets all databases per server.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListPostgreSQLDBContext] instead.
func ListPostgreSQLDB(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) []*armpostgresql.Database {
	t.Helper()

	return ListPostgreSQLDBContext(t, context.Background(), subscriptionID, resGroupName, serverName) //nolint:staticcheck
}

// ListPostgreSQLDBContextE is a helper function that gets all databases per server.
// The ctx parameter supports cancellation and timeouts.
func ListPostgreSQLDBContextE(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, serverName string) ([]*armpostgresql.Database, error) {
	// Create a postgresql db client
	postgresqldbClient, err := GetPostgreSQLDBClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return ListPostgreSQLDBWithClient(ctx, postgresqldbClient, resGroupName, serverName)
}

// ListPostgreSQLDBWithClient lists all databases per server using the provided DatabasesClient.
func ListPostgreSQLDBWithClient(ctx context.Context, client *armpostgresql.DatabasesClient, resGroupName string, serverName string) ([]*armpostgresql.Database, error) {
	pager := client.NewListByServerPager(resGroupName, serverName, nil)

	var databases []*armpostgresql.Database

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		databases = append(databases, page.Value...)
	}

	return databases, nil
}

// ListPostgreSQLDBE is a helper function that gets all databases per server.
//
// Deprecated: Use [ListPostgreSQLDBContextE] instead.
func ListPostgreSQLDBE(t testing.TestingT, subscriptionID string, resGroupName string, serverName string) ([]*armpostgresql.Database, error) {
	return ListPostgreSQLDBContextE(t, context.Background(), subscriptionID, resGroupName, serverName)
}
