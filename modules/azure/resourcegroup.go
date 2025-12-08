package azure

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/resources/mgmt/resources"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
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
		if resourceGroupNotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return exists, nil

}

// GetResourceGroupE gets a resource group within a subscription
func GetResourceGroupE(resourceGroupName, subscriptionID string) (bool, error) {

	rg, err := GetAResourceGroupE(resourceGroupName, subscriptionID)
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

// GetAResourceGroup returns a resource group within a subscription
// This function would fail the test if there is an error.
func GetAResourceGroup(t *testing.T, resourceGroupName string, subscriptionID string) *resources.Group {
	rg, err := GetAResourceGroupE(resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return rg
}

// GetAResourceGroupE gets a resource group within a subscription
func GetAResourceGroupE(resourceGroupName, subscriptionID string) (*resources.Group, error) {
	client, err := CreateResourceGroupClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	rg, err := client.Get(context.Background(), resourceGroupName)
	if err != nil {
		return nil, err
	}
	return &rg, nil
}

// ListResourceGroupsByTag returns a resource group list within a subscription based on a tag key
// This function would fail the test if there is an error.
func ListResourceGroupsByTag(t *testing.T, tag, subscriptionID string) []resources.Group {
	rg, err := ListResourceGroupsByTagE(tag, subscriptionID)
	require.NoError(t, err)
	return rg
}

// ListResourceGroupsByTagE returns a resource group list within a subscription based on a tag key
func ListResourceGroupsByTagE(tag string, subscriptionID string) ([]resources.Group, error) {
	client, err := CreateResourceGroupClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	rg, err := client.List(context.Background(), fmt.Sprintf("tagName eq '%s'", tag), nil)
	if err != nil {
		return nil, err
	}
	return rg.Values(), nil
}

func resourceGroupNotFoundError(err error) bool {
	if err != nil {
		if autorestError, ok := err.(autorest.DetailedError); ok {
			if requestError, ok := autorestError.Original.(*azure.RequestError); ok {
				return (requestError.ServiceError.Code == "ResourceGroupNotFound")
			}
		}
		if azcoreErr, ok := err.(*azcore.ResponseError); ok {
			return azcoreErr.ErrorCode == "ResourceGroupNotFound"
		}
	}
	return false
}
