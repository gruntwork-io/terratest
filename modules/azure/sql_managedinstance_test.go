//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.
package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure SQL DB, these tests can be extended
*/

func TestSQLManagedInstanceExists(t *testing.T) {
	t.Parallel()

	managedInstanceName := ""
	resourceGroupName := ""
	subscriptionID := ""

	exists, err := SQLManagedInstanceExistsE(managedInstanceName, resourceGroupName, subscriptionID)

	require.False(t, exists)
	require.Error(t, err)
}

func TestGetManagedInstanceE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	managedInstanceName := ""
	subscriptionID := ""

	_, err := GetManagedInstanceE(subscriptionID, resGroupName, managedInstanceName)
	require.Error(t, err)
}

func TestGetManagedInstanceDatabasesE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	managedInstanceName := ""
	databaseName := ""
	subscriptionID := ""

	_, err := GetManagedInstanceDatabaseE(t, subscriptionID, resGroupName, managedInstanceName, databaseName)
	require.Error(t, err)
}
