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
If/when CRUD methods are introduced for Azure Virtual Machines, these tests can be extended.
*/

func TestManagedEnvironmentExists(t *testing.T) {
	t.Parallel()

	environmentName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := ManagedEnvironmentExistsE(environmentName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetManagedEnvironmentE(t *testing.T) {
	t.Parallel()

	environmentName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := GetManagedEnvironmentE(environmentName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}

func TestContainerAppExists(t *testing.T) {
	t.Parallel()

	environmentName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := ContainerAppExistsE(environmentName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetContainerAppE(t *testing.T) {
	t.Parallel()

	environmentName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := GetContainerAppE(environmentName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}

func TestContainerAppJobExists(t *testing.T) {
	t.Parallel()

	environmentName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := ContainerAppJobExistsE(environmentName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetContainerJobAppE(t *testing.T) {
	t.Parallel()

	environmentName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := GetContainerAppJobE(environmentName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}
