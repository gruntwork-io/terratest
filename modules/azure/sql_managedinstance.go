package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
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
	sqlmiClient, err := CreateSQLManagedInstanceClientContext(ctx, subscriptionID)
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
	sqlmiDBClient, err := CreateSQLManagedDatabasesClientContext(ctx, subscriptionID)
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
