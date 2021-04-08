package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-05-01/resources"
	"github.com/stretchr/testify/require"
)

// ResourceGroupExists indicates whether a resource group exists within a subscription; otherwise false
// This function would fail the test if there is an error.
func ResourceGroupExists(t *testing.T, resourceGroupName string, subscriptionID string) bool {
	result, err := ResourceGroupExistsE(resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// ResourceGroupExistsE indicates whether a resource group exists within a subscription
func ResourceGroupExistsE(resourceGroupName, subscriptionID string) (bool, error) {
	exists, err := GetResourceGroupE(resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return exists, nil

}

// GetResourceGroupE gets a resource group within a subscription
func GetResourceGroupE(resourceGroupName, subscriptionID string) (bool, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return false, err
	}
	client, err := CreateResourceGroupClientE(subscriptionID)
	if err != nil {
		return false, err
	}
	rg, err := client.Get(context.Background(), resourceGroupName)
	if err != nil {
		return false, err
	}
	return (resourceGroupName == *rg.Name), nil
}

// GetResourceGroupClientE gets a resource group client in a subscription
// TODO: remove in next version
func GetResourceGroupClientE(subscriptionID string) (*resources.GroupsClient, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}
	resourceGroupClient := resources.NewGroupsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	resourceGroupClient.Authorizer = *authorizer
	return &resourceGroupClient, nil
}
