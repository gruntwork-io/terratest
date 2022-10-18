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

func TestGetSQLServerE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""

	_, err := GetSQLServerE(t, resGroupName, serverName, subscriptionID)
	require.Error(t, err)
}

func TestGetSQLDatabaseE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	dbName := ""
	subscriptionID := ""

	_, err := GetSQLDatabaseE(t, resGroupName, serverName, dbName, subscriptionID)
	require.Error(t, err)
}
