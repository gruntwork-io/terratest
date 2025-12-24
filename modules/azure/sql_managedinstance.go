package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql"
	"github.com/gruntwork-io/terratest/modules/testing/v2"
	"github.com/stretchr/testify/require"
)

// SQLManagedInstanceExists indicates whether the SQL Managed Instance exists for the subscription.
// This function would fail the test if there is an error.
func SQLManagedInstanceExists(t testing.TestingT, managedInstanceName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := SQLManagedInstanceExistsE(managedInstanceName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// SQLManagedInstanceExistsE indicates whether the specified SQL Managed Instance exists and may return an error.
func SQLManagedInstanceExistsE(managedInstanceName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetManagedInstanceE(subscriptionID, resourceGroupName, managedInstanceName)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetManagedInstance is a helper function that gets the sql managed instance object.
// This function would fail the test if there is an error.
func GetManagedInstance(t testing.TestingT, resGroupName string, managedInstanceName string, subscriptionID string) *armsql.ManagedInstance {
	managedInstance, err := GetManagedInstanceE(subscriptionID, resGroupName, managedInstanceName)
	require.NoError(t, err)

	return managedInstance
}

// GetManagedInstanceDatabase is a helper function that gets the sql managed database object.
// This function would fail the test if there is an error.
func GetManagedInstanceDatabase(t testing.TestingT, resGroupName string, managedInstanceName string, databaseName string, subscriptionID string) *armsql.ManagedDatabase {
	managedDatabase, err := GetManagedInstanceDatabaseE(t, subscriptionID, resGroupName, managedInstanceName, databaseName)
	require.NoError(t, err)

	return managedDatabase
}

// GetManagedInstanceE is a helper function that gets the sql managed instance object.
func GetManagedInstanceE(subscriptionID string, resGroupName string, managedInstanceName string) (*armsql.ManagedInstance, error) {
	// Create a SQL Managed Instance client
	sqlmiClient, err := CreateSQLMangedInstanceClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding managed instance
	resp, err := sqlmiClient.Get(context.Background(), resGroupName, managedInstanceName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ManagedInstance, nil
}

// GetManagedInstanceDatabaseE is a helper function that gets the sql managed database object.
func GetManagedInstanceDatabaseE(t testing.TestingT, subscriptionID string, resGroupName string, managedInstanceName string, databaseName string) (*armsql.ManagedDatabase, error) {
	// Create a SQL MI db client
	sqlmiDbClient, err := CreateSQLMangedDatabasesClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding database
	resp, err := sqlmiDbClient.Get(context.Background(), resGroupName, managedInstanceName, databaseName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ManagedDatabase, nil
}
