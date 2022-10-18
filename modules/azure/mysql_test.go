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
If/when CRUD methods are introduced for Azure MySQL server and database, these tests can be extended
*/

func TestGetMYSQLServerE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""

	_, err := GetMYSQLServerE(t, subscriptionID, resGroupName, serverName)
	require.Error(t, err)
}

func TestGetMYSQLDBE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""
	dbName := ""

	_, err := GetMYSQLDBE(t, subscriptionID, resGroupName, serverName, dbName)
	require.Error(t, err)
}
