//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
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

	exists, err := azure.SQLManagedInstanceExistsContextE(t.Context(), managedInstanceName, resourceGroupName, subscriptionID)

	require.False(t, exists)
	require.Error(t, err)
}

func TestGetManagedInstanceE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	managedInstanceName := ""
	subscriptionID := ""

	_, err := azure.GetManagedInstanceContextE(t.Context(), subscriptionID, resGroupName, managedInstanceName)
	require.Error(t, err)
}

func TestGetManagedInstanceDatabasesE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	managedInstanceName := ""
	databaseName := ""
	subscriptionID := ""

	_, err := azure.GetManagedInstanceDatabaseContextE(t.Context(), subscriptionID, resGroupName, managedInstanceName, databaseName)
	require.Error(t, err)
}
