package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// SQLManagedInstanceExistsContext indicates whether the SQL Managed Instance exists for the subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func SQLManagedInstanceExistsContext(t testing.TestingT, ctx context.Context, managedInstanceName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := SQLManagedInstanceExistsContextE(ctx, managedInstanceName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// SQLManagedInstanceExistsContextE indicates whether the specified SQL Managed Instance exists.
// The ctx parameter supports cancellation and timeouts.
func SQLManagedInstanceExistsContextE(ctx context.Context, managedInstanceName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetManagedInstanceContextE(ctx, subscriptionID, resourceGroupName, managedInstanceName)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// SQLManagedInstanceExists indicates whether the SQL Managed Instance exists for the subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [SQLManagedInstanceExistsContext] instead.
func SQLManagedInstanceExists(t testing.TestingT, managedInstanceName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return SQLManagedInstanceExistsContext(t, context.Background(), managedInstanceName, resourceGroupName, subscriptionID)
}

// SQLManagedInstanceExistsE indicates whether the specified SQL Managed Instance exists and may return an error.
//
// Deprecated: Use [SQLManagedInstanceExistsContextE] instead.
func SQLManagedInstanceExistsE(managedInstanceName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return SQLManagedInstanceExistsContextE(context.Background(), managedInstanceName, resourceGroupName, subscriptionID)
}

// GetManagedInstanceContext retrieves the SQL managed instance object for the given subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetManagedInstanceContext(t testing.TestingT, ctx context.Context, resGroupName string, managedInstanceName string, subscriptionID string) *armsql.ManagedInstance {
	t.Helper()

	managedInstance, err := GetManagedInstanceContextE(ctx, subscriptionID, resGroupName, managedInstanceName)
	require.NoError(t, err)

	return managedInstance
}

// GetManagedInstanceContextE retrieves the SQL managed instance object for the given subscription.
// The ctx parameter supports cancellation and timeouts.
func GetManagedInstanceContextE(ctx context.Context, subscriptionID string, resGroupName string, managedInstanceName string) (*armsql.ManagedInstance, error) {
	sqlmiClient, err := CreateSQLMangedInstanceClientContext(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetManagedInstanceWithClient(ctx, sqlmiClient, resGroupName, managedInstanceName)
}

// GetManagedInstanceWithClient retrieves the SQL managed instance using the provided ManagedInstancesClient.
func GetManagedInstanceWithClient(ctx context.Context, client *armsql.ManagedInstancesClient, resGroupName string, managedInstanceName string) (*armsql.ManagedInstance, error) {
	resp, err := client.Get(ctx, resGroupName, managedInstanceName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ManagedInstance, nil
}

// GetManagedInstance retrieves the SQL managed instance object for the given subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetManagedInstanceContext] instead.
func GetManagedInstance(t testing.TestingT, resGroupName string, managedInstanceName string, subscriptionID string) *armsql.ManagedInstance {
	t.Helper()

	return GetManagedInstanceContext(t, context.Background(), resGroupName, managedInstanceName, subscriptionID)
}

// GetManagedInstanceE retrieves the SQL managed instance object for the given subscription.
//
// Deprecated: Use [GetManagedInstanceContextE] instead.
func GetManagedInstanceE(subscriptionID string, resGroupName string, managedInstanceName string) (*armsql.ManagedInstance, error) {
	return GetManagedInstanceContextE(context.Background(), subscriptionID, resGroupName, managedInstanceName)
}

// GetManagedInstanceDatabaseContext retrieves the SQL managed database object for the given subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetManagedInstanceDatabaseContext(t testing.TestingT, ctx context.Context, resGroupName string, managedInstanceName string, databaseName string, subscriptionID string) *armsql.ManagedDatabase {
	t.Helper()

	managedDatabase, err := GetManagedInstanceDatabaseContextE(ctx, subscriptionID, resGroupName, managedInstanceName, databaseName)
	require.NoError(t, err)

	return managedDatabase
}

// GetManagedInstanceDatabaseContextE retrieves the SQL managed database object for the given subscription.
// The ctx parameter supports cancellation and timeouts.
func GetManagedInstanceDatabaseContextE(ctx context.Context, subscriptionID string, resGroupName string, managedInstanceName string, databaseName string) (*armsql.ManagedDatabase, error) {
	sqlmiDBClient, err := CreateSQLMangedDatabasesClientContext(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetManagedInstanceDatabaseWithClient(ctx, sqlmiDBClient, resGroupName, managedInstanceName, databaseName)
}

// GetManagedInstanceDatabaseWithClient retrieves the SQL managed database using the provided ManagedDatabasesClient.
func GetManagedInstanceDatabaseWithClient(ctx context.Context, client *armsql.ManagedDatabasesClient, resGroupName string, managedInstanceName string, databaseName string) (*armsql.ManagedDatabase, error) {
	resp, err := client.Get(ctx, resGroupName, managedInstanceName, databaseName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ManagedDatabase, nil
}

// GetManagedInstanceDatabase retrieves the SQL managed database object for the given subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetManagedInstanceDatabaseContext] instead.
func GetManagedInstanceDatabase(t testing.TestingT, resGroupName string, managedInstanceName string, databaseName string, subscriptionID string) *armsql.ManagedDatabase {
	t.Helper()

	return GetManagedInstanceDatabaseContext(t, context.Background(), resGroupName, managedInstanceName, databaseName, subscriptionID)
}

// GetManagedInstanceDatabaseE retrieves the SQL managed database object for the given subscription.
//
// Deprecated: Use [GetManagedInstanceDatabaseContextE] instead.
func GetManagedInstanceDatabaseE(t testing.TestingT, subscriptionID string, resGroupName string, managedInstanceName string, databaseName string) (*armsql.ManagedDatabase, error) { //nolint:unparam // t kept for API compatibility
	return GetManagedInstanceDatabaseContextE(context.Background(), subscriptionID, resGroupName, managedInstanceName, databaseName)
}
