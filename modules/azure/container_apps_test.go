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
If/when CRUD methods are introduced for Azure Virtual Machines, these tests can be extended.
*/

func TestManagedEnvironmentExistsContextE(t *testing.T) {
	t.Parallel()

	environmentName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := azure.ManagedEnvironmentExistsContextE(t.Context(), environmentName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetManagedEnvironmentContextE(t *testing.T) {
	t.Parallel()

	environmentName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := azure.GetManagedEnvironmentContextE(t.Context(), environmentName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}

func TestContainerAppExistsContextE(t *testing.T) {
	t.Parallel()

	environmentName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := azure.ContainerAppExistsContextE(t.Context(), environmentName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetContainerAppContextE(t *testing.T) {
	t.Parallel()

	environmentName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := azure.GetContainerAppContextE(t.Context(), environmentName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}

func TestContainerAppJobExistsContextE(t *testing.T) {
	t.Parallel()

	environmentName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := azure.ContainerAppJobExistsContextE(t.Context(), environmentName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetContainerAppJobContextE(t *testing.T) {
	t.Parallel()

	environmentName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := azure.GetContainerAppJobContextE(t.Context(), environmentName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}
