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
If/when CRUD methods are introduced for Azure PostgreSQL server and database, these tests can be extended
*/

func TestGetPostgreSQLServerContextE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""

	_, err := azure.GetPostgreSQLServerContextE(t, t.Context(), subscriptionID, resGroupName, serverName)
	require.Error(t, err)
}

func TestGetPostgreSQLDBContextE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	serverName := ""
	subscriptionID := ""
	dbName := ""

	_, err := azure.GetPostgreSQLDBContextE(t, t.Context(), subscriptionID, resGroupName, serverName, dbName)
	require.Error(t, err)
}
