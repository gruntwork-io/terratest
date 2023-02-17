//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods to create and delete network resources are added, these tests can be extended.
*/

func TestGetActionGroupResourceEWithMissingResourceGroupName(t *testing.T) {
	t.Parallel()

	ruleName := "Hello"
	resGroupName := ""
	subscriptionID := ""

	_, err := GetActionGroupResourceE(ruleName, resGroupName, subscriptionID)

	require.Error(t, err)
}

func TestGetActionGroupResourceEWithInvalidResourceGroupName(t *testing.T) {
	t.Parallel()

	ruleName := ""
	resGroupName := "Hello"
	subscriptionID := ""

	_, err := GetActionGroupResourceE(ruleName, resGroupName, subscriptionID)

	require.Error(t, err)
}

func TestGetActionGroupClient(t *testing.T) {
	t.Parallel()

	subscriptionID := ""

	client, err := getActionGroupClient(subscriptionID)

	require.NoError(t, err)
	assert.NotEmpty(t, *client)
}
