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
If/when methods to create and delete front door are added, these tests can be extended.
*/

func TestFrontDoorExists(t *testing.T) {
	t.Parallel()

	frontDoorName := "TestFrontDoor"
	resourceGroupName := "TestResourceGroup"
	subscriptionID := ""

	exists, err := azure.FrontDoorExistsContextE(t.Context(), frontDoorName, resourceGroupName, subscriptionID)

	require.False(t, exists)
	require.Error(t, err)
}

func TestGetFrontDoor(t *testing.T) {
	t.Parallel()

	frontDoorName := "TestFrontDoor"
	resourceGroupName := "TestResourceGroup"
	subscriptionID := ""

	instance, err := azure.GetFrontDoorContextE(t.Context(), frontDoorName, resourceGroupName, subscriptionID)

	require.Nil(t, instance)
	require.Error(t, err)
}

func TestFrontDoorFrontendEndpointExists(t *testing.T) {
	t.Parallel()

	endpointName := "TestFrontendEndpoint"
	frontDoorName := "TestFrontDoor"
	resourceGroupName := "TestResourceGroup"
	subscriptionID := ""

	endpoint, err := azure.FrontDoorFrontendEndpointExistsContextE(t.Context(), endpointName, frontDoorName, resourceGroupName, subscriptionID)

	require.False(t, endpoint)
	require.Error(t, err)
}

func TestGetFrontDoorFrontendEndpoint(t *testing.T) {
	t.Parallel()

	endpointName := "TestFrontendEndpoint"
	frontDoorName := "TestFrontDoor"
	resourceGroupName := "TestResourceGroup"
	subscriptionID := ""

	endpoint, err := azure.GetFrontDoorFrontendEndpointContextE(t.Context(), endpointName, frontDoorName, resourceGroupName, subscriptionID)

	require.Nil(t, endpoint)
	require.Error(t, err)
}
