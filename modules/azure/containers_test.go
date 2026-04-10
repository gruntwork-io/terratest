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

func TestContainerRegistryExistsContextE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	registryName := ""
	subscriptionID := ""

	_, err := azure.ContainerRegistryExistsContextE(t.Context(), registryName, resGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetContainerRegistryContextE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	registryName := ""
	subscriptionID := ""

	_, err := azure.GetContainerRegistryContextE(t.Context(), registryName, resGroupName, subscriptionID)
	require.Error(t, err)
}

func TestContainerInstanceExistsContextE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	instanceName := ""
	subscriptionID := ""

	_, err := azure.ContainerInstanceExistsContextE(t.Context(), instanceName, resGroupName, subscriptionID)
	require.Error(t, err)
}

func TestGetContainerInstanceContextE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	instanceName := ""
	subscriptionID := ""

	_, err := azure.GetContainerInstanceContextE(t.Context(), instanceName, resGroupName, subscriptionID)
	require.Error(t, err)
}
