package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// ResourceGroupExistsV2Context indicates whether a resource group exists within a subscription; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ResourceGroupExistsV2Context(t testing.TestingT, ctx context.Context, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	result, err := ResourceGroupExistsV2ContextE(ctx, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return result
}

// ResourceGroupExistsV2ContextE indicates whether a resource group exists within a subscription.
// The ctx parameter supports cancellation and timeouts.
func ResourceGroupExistsV2ContextE(ctx context.Context, resourceGroupName, subscriptionID string) (bool, error) {
	exists, err := GetResourceGroupV2ContextE(ctx, resourceGroupName, subscriptionID)
	if err != nil {
		if resourceGroupNotFoundError(err) {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}

// GetResourceGroupV2ContextE checks whether a resource group name matches the one retrieved from the subscription.
// The ctx parameter supports cancellation and timeouts.
func GetResourceGroupV2ContextE(ctx context.Context, resourceGroupName, subscriptionID string) (bool, error) {
	rg, err := GetAResourceGroupV2ContextE(ctx, resourceGroupName, subscriptionID)
	if err != nil {
		return false, err
	}

	return (resourceGroupName == *rg.Name), nil
}

// GetAResourceGroupV2Context returns a resource group within a subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAResourceGroupV2Context(t testing.TestingT, ctx context.Context, resourceGroupName string, subscriptionID string) *armresources.ResourceGroup {
	t.Helper()

	rg, err := GetAResourceGroupV2ContextE(ctx, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return rg
}

// GetAResourceGroupV2ContextE gets a resource group within a subscription.
// The ctx parameter supports cancellation and timeouts.
func GetAResourceGroupV2ContextE(ctx context.Context, resourceGroupName, subscriptionID string) (*armresources.ResourceGroup, error) {
	client, err := CreateResourceGroupClientV2E(subscriptionID)
	if err != nil {
		return nil, err
	}

	rg, err := client.Get(ctx, resourceGroupName, &armresources.ResourceGroupsClientGetOptions{})
	if err != nil {
		return nil, err
	}

	return &rg.ResourceGroup, nil
}

// ListResourceGroupsByTagV2Context returns a resource group list within a subscription based on a tag key.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListResourceGroupsByTagV2Context(t testing.TestingT, ctx context.Context, tag, subscriptionID string) []*armresources.ResourceGroup {
	t.Helper()

	rg, err := ListResourceGroupsByTagV2ContextE(ctx, tag, subscriptionID)
	require.NoError(t, err)

	return rg
}

// ListResourceGroupsByTagV2ContextE returns a resource group list within a subscription based on a tag key.
// The ctx parameter supports cancellation and timeouts.
func ListResourceGroupsByTagV2ContextE(ctx context.Context, tag string, subscriptionID string) (rg []*armresources.ResourceGroup, err error) {
	client, err := CreateResourceGroupClientV2E(subscriptionID)
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf("tagName eq '%s'", tag)
	pager := client.NewListPager(&armresources.ResourceGroupsClientListOptions{
		Filter: &filter,
	})

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		rg = append(rg, page.ResourceGroupListResult.Value...)
	}

	return
}
