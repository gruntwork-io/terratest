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
If/when CRUD methods are introduced for Azure MySQL server and database, these tests can be extended
*/

func TestContainerRegistryExistsE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	registryName := ""
	subscriptionID := ""

	_, err := azure.ContainerRegistryExistsE(registryName, resGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetContainerRegistryE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	registryName := ""
	subscriptionID := ""

	_, err := azure.GetContainerRegistryE(registryName, resGroupName, subscriptionID)
	require.Error(t, err)
}

func TestCreateContainerRegistryClientE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""

	_, err := azure.CreateContainerRegistryClientE(subscriptionID)
	require.NoError(t, err)
}

func TestContainerInstanceExistsE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	instanceName := ""
	subscriptionID := ""

	_, err := azure.ContainerInstanceExistsE(instanceName, resGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetContainerInstanceE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	instanceName := ""
	subscriptionID := ""

	_, err := azure.GetContainerInstanceE(instanceName, resGroupName, subscriptionID)
	require.Error(t, err)
}

func TestCreateContainerInstanceClientE(t *testing.T) {
	t.Parallel()

	subscriptionID := ""

	_, err := azure.CreateContainerInstanceClientE(subscriptionID)
	require.NoError(t, err)
}
