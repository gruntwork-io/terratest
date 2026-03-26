package azure

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/stretchr/testify/require"
)

// ResourceGroupExistsV2Context indicates whether a resource group exists within a subscription; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ResourceGroupExistsV2Context(t *testing.T, ctx context.Context, resourceGroupName string, subscriptionID string) bool {
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

// ResourceGroupExistsV2 indicates whether a resource group exists within a subscription; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ResourceGroupExistsV2Context] instead.
func ResourceGroupExistsV2(t *testing.T, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return ResourceGroupExistsV2Context(t, context.Background(), resourceGroupName, subscriptionID)
}

// ResourceGroupExistsV2E indicates whether a resource group exists within a subscription.
//
// Deprecated: Use [ResourceGroupExistsV2ContextE] instead.
func ResourceGroupExistsV2E(resourceGroupName, subscriptionID string) (bool, error) {
	return ResourceGroupExistsV2ContextE(context.Background(), resourceGroupName, subscriptionID)
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

// GetResourceGroupV2E checks whether a resource group name matches the one retrieved from the subscription.
//
// Deprecated: Use [GetResourceGroupV2ContextE] instead.
func GetResourceGroupV2E(resourceGroupName, subscriptionID string) (bool, error) {
	return GetResourceGroupV2ContextE(context.Background(), resourceGroupName, subscriptionID)
}

// GetAResourceGroupV2Context returns a resource group within a subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAResourceGroupV2Context(t *testing.T, ctx context.Context, resourceGroupName string, subscriptionID string) *armresources.ResourceGroup {
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

// GetAResourceGroupV2 returns a resource group within a subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetAResourceGroupV2Context] instead.
func GetAResourceGroupV2(t *testing.T, resourceGroupName string, subscriptionID string) *armresources.ResourceGroup {
	t.Helper()

	return GetAResourceGroupV2Context(t, context.Background(), resourceGroupName, subscriptionID)
}

// GetAResourceGroupV2E gets a resource group within a subscription.
//
// Deprecated: Use [GetAResourceGroupV2ContextE] instead.
func GetAResourceGroupV2E(resourceGroupName, subscriptionID string) (*armresources.ResourceGroup, error) {
	return GetAResourceGroupV2ContextE(context.Background(), resourceGroupName, subscriptionID)
}

// ListResourceGroupsByTagV2Context returns a resource group list within a subscription based on a tag key.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListResourceGroupsByTagV2Context(t *testing.T, ctx context.Context, tag, subscriptionID string) []*armresources.ResourceGroup {
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

// ListResourceGroupsByTagV2 returns a resource group list within a subscription based on a tag key.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListResourceGroupsByTagV2Context] instead.
func ListResourceGroupsByTagV2(t *testing.T, tag, subscriptionID string) []*armresources.ResourceGroup {
	t.Helper()

	return ListResourceGroupsByTagV2Context(t, context.Background(), tag, subscriptionID)
}

// ListResourceGroupsByTagV2E returns a resource group list within a subscription based on a tag key.
//
// Deprecated: Use [ListResourceGroupsByTagV2ContextE] instead.
func ListResourceGroupsByTagV2E(tag string, subscriptionID string) (rg []*armresources.ResourceGroup, err error) {
	return ListResourceGroupsByTagV2ContextE(context.Background(), tag, subscriptionID)
}
