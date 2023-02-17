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

func TestAppExistsE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	appName := ""
	subscriptionID := ""

	_, err := AppExistsE(appName, resGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetAppServiceE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	appName := ""
	subscriptionID := ""

	_, err := GetAppServiceE(appName, resGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetAppServiceClientE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""

	_, err := GetAppServiceClientE(subscriptionID)
	require.NoError(t, err)
}
