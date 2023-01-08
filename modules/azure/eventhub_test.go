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
If/when CRUD methods are introduced for Azure PostgreSQL server and database, these tests can be extended
*/

func TestGetEventHubE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	namespace := ""
	subscriptionID := ""

	_, err := GetEventHubE(t, subscriptionID, resGroupName, namespace)
	require.Error(t, err)
}

// func TestGetPostgreSQLDBE(t *testing.T) {
// 	t.Parallel()

// 	resGroupName := ""
// 	serverName := ""
// 	subscriptionID := ""
// 	dbName := ""

// 	_, err := GetPostgreSQLDBE(t, subscriptionID, resGroupName, serverName, dbName)
// 	require.Error(t, err)
// }
