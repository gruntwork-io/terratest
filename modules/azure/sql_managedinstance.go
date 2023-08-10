package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/gruntwork-io/terratest/modules/testing"
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

// GetManagedInstance is a helper function that gets the sql server object.
// This function would fail the test if there is an error.
func GetManagedInstance(t testing.TestingT, resGroupName string, managedInstanceName string, subscriptionID string) *sql.ManagedInstance {
	managedInstance, err := GetManagedInstanceE(subscriptionID, resGroupName, managedInstanceName)
	require.NoError(t, err)

	return managedInstance
}

// GetManagedInstanceDatabase is a helper function that gets the sql server object.
// This function would fail the test if there is an error.
func GetManagedInstanceDatabase(t testing.TestingT, resGroupName string, managedInstanceName string, databaseName string, subscriptionID string) *sql.ManagedDatabase {
	managedDatabase, err := GetManagedInstanceDatabaseE(t, subscriptionID, resGroupName, managedInstanceName, databaseName)
	require.NoError(t, err)

	return managedDatabase
}

// GetManagedInstanceE is a helper function that gets the sql server object.
func GetManagedInstanceE(subscriptionID string, resGroupName string, managedInstanceName string) (*sql.ManagedInstance, error) {
	// Create a SQl Server client
	sqlmiClient, err := CreateSQLMangedInstanceClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	//Get the corresponding server client
	sqlmi, err := sqlmiClient.Get(context.Background(), resGroupName, managedInstanceName)
	if err != nil {
		return nil, err
	}

	//Return sql mi
	return &sqlmi, nil
}

// GetManagedInstanceDatabaseE is a helper function that gets the sql server object.
func GetManagedInstanceDatabaseE(t testing.TestingT, subscriptionID string, resGroupName string, managedInstanceName string, databaseName string) (*sql.ManagedDatabase, error) {
	// Create a SQlMI db client
	sqlmiDbClient, err := CreateSQLMangedDatabasesClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	//Get the corresponding  client
	sqlmidb, err := sqlmiDbClient.Get(context.Background(), resGroupName, managedInstanceName, databaseName)
	if err != nil {
		return nil, err
	}

	//Return sql mi db
	return &sqlmidb, nil
}
