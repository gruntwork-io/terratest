package azure

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/require"
)

// ResourceGroupExists indicates whether a resource group exists within a subscription; otherwise false
// This function would fail the test if there is an error.
func ResourceGroupExistsV2(t *testing.T, resourceGroupName string, subscriptionID string) bool {
	result, err := ResourceGroupExistsV2E(resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// ResourceGroupExistsE indicates whether a resource group exists within a subscription
func ResourceGroupExistsV2E(resourceGroupName, subscriptionID string) (bool, error) {
	exists, err := GetResourceGroupV2E(resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return exists, nil

}

// GetResourceGroupE gets a resource group within a subscription
func GetResourceGroupV2E(resourceGroupName, subscriptionID string) (bool, error) {
	rg, err := GetAResourceGroupV2E(resourceGroupName, subscriptionID)
	if err != nil {
		return false, err
	}
	return (resourceGroupName == *rg.Name), nil
}

// GetAResourceGroup returns a resource group within a subscription
// This function would fail the test if there is an error.
func GetAResourceGroupV2(t *testing.T, resourceGroupName string, subscriptionID string) *armresources.ResourceGroup {
	rg, err := GetAResourceGroupV2E(resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return rg
}

// GetAResourceGroupE gets a resource group within a subscription
func GetAResourceGroupV2E(resourceGroupName, subscriptionID string) (*armresources.ResourceGroup, error) {
	client, err := CreateResourceGroupClientV2E(subscriptionID)
	if err != nil {
		return nil, err
	}

	rg, err := client.Get(context.Background(), resourceGroupName, &armresources.ResourceGroupsClientGetOptions{})
	if err != nil {
		return nil, err
	}
	return &rg.ResourceGroup, nil
}

// ListResourceGroupsByTag returns a resource group list within a subscription based on a tag key
// This function would fail the test if there is an error.
func ListResourceGroupsByTagV2(t *testing.T, tag, subscriptionID string) []*armresources.ResourceGroup {
	rg, err := ListResourceGroupsByTagV2E(tag, subscriptionID)
	require.NoError(t, err)
	return rg
}

// ListResourceGroupsByTagE returns a resource group list within a subscription based on a tag key
func ListResourceGroupsByTagV2E(tag string, subscriptionID string) (rg []*armresources.ResourceGroup, err error) {
	client, err := CreateResourceGroupClientV2E(subscriptionID)
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf("tagName eq '%s'", tag)
	pager := client.NewListPager(&armresources.ResourceGroupsClientListOptions{
		Filter: &filter,
	})
	ctx := context.Background()
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		rg = append(rg, page.ResourceGroupListResult.Value...)
	}

	return
}
