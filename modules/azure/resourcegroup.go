package azure

import (
	"context"
	"errors"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"github.com/Azure/go-autorest/autorest"
	autorestAzure "github.com/Azure/go-autorest/autorest/azure"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// ResourceGroupExistsContext indicates whether a resource group exists within a subscription; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ResourceGroupExistsContext(t testing.TestingT, ctx context.Context, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	result, err := ResourceGroupExistsContextE(ctx, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return result
}

// ResourceGroupExistsContextE indicates whether a resource group exists within a subscription.
// The ctx parameter supports cancellation and timeouts.
func ResourceGroupExistsContextE(ctx context.Context, resourceGroupName, subscriptionID string) (bool, error) {
	exists, err := GetResourceGroupContextE(ctx, resourceGroupName, subscriptionID)
	if err != nil {
		if resourceGroupNotFoundError(err) {
			return false, nil
		}

		return false, err
	}

	return exists, nil
}

// ResourceGroupExists indicates whether a resource group exists within a subscription; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ResourceGroupExistsContext] instead.
func ResourceGroupExists(t testing.TestingT, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return ResourceGroupExistsContext(t, context.Background(), resourceGroupName, subscriptionID)
}

// ResourceGroupExistsE indicates whether a resource group exists within a subscription.
//
// Deprecated: Use [ResourceGroupExistsContextE] instead.
func ResourceGroupExistsE(resourceGroupName, subscriptionID string) (bool, error) {
	return ResourceGroupExistsContextE(context.Background(), resourceGroupName, subscriptionID)
}

// GetResourceGroupContextE checks whether a resource group name matches the one retrieved from the subscription.
// The ctx parameter supports cancellation and timeouts.
func GetResourceGroupContextE(ctx context.Context, resourceGroupName, subscriptionID string) (bool, error) {
	rg, err := GetAResourceGroupContextE(ctx, resourceGroupName, subscriptionID)
	if err != nil {
		return false, err
	}

	return (resourceGroupName == *rg.Name), nil
}

// GetResourceGroupE checks whether a resource group name matches the one retrieved from the subscription.
//
// Deprecated: Use [GetResourceGroupContextE] instead.
func GetResourceGroupE(resourceGroupName, subscriptionID string) (bool, error) {
	return GetResourceGroupContextE(context.Background(), resourceGroupName, subscriptionID)
}

// GetAResourceGroupContext returns a resource group within a subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAResourceGroupContext(t testing.TestingT, ctx context.Context, resourceGroupName string, subscriptionID string) *armresources.ResourceGroup {
	t.Helper()

	rg, err := GetAResourceGroupContextE(ctx, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return rg
}

// GetAResourceGroupContextE gets a resource group within a subscription.
// The ctx parameter supports cancellation and timeouts.
func GetAResourceGroupContextE(ctx context.Context, resourceGroupName, subscriptionID string) (*armresources.ResourceGroup, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	client, err := armresources.NewResourceGroupsClient(subscriptionID, cred, opts)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, resourceGroupName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ResourceGroup, nil
}

// GetAResourceGroup returns a resource group within a subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetAResourceGroupContext] instead.
func GetAResourceGroup(t testing.TestingT, resourceGroupName string, subscriptionID string) *armresources.ResourceGroup {
	t.Helper()

	return GetAResourceGroupContext(t, context.Background(), resourceGroupName, subscriptionID)
}

// GetAResourceGroupE gets a resource group within a subscription.
//
// Deprecated: Use [GetAResourceGroupContextE] instead.
func GetAResourceGroupE(resourceGroupName, subscriptionID string) (*armresources.ResourceGroup, error) {
	return GetAResourceGroupContextE(context.Background(), resourceGroupName, subscriptionID)
}

// ListResourceGroupsByTagContext returns a resource group list within a subscription based on a tag key.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListResourceGroupsByTagContext(t testing.TestingT, ctx context.Context, tag, subscriptionID string) []armresources.ResourceGroup {
	t.Helper()

	rg, err := ListResourceGroupsByTagContextE(ctx, tag, subscriptionID)
	require.NoError(t, err)

	return rg
}

// ListResourceGroupsByTagContextE returns a resource group list within a subscription based on a tag key.
// The ctx parameter supports cancellation and timeouts.
func ListResourceGroupsByTagContextE(ctx context.Context, tag string, subscriptionID string) ([]armresources.ResourceGroup, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	client, err := armresources.NewResourceGroupsClient(subscriptionID, cred, opts)
	if err != nil {
		return nil, err
	}

	filter := fmt.Sprintf("tagName eq '%s'", tag)
	pager := client.NewListPager(&armresources.ResourceGroupsClientListOptions{
		Filter: &filter,
	})

	var groups []armresources.ResourceGroup

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, rg := range page.Value {
			groups = append(groups, *rg)
		}
	}

	return groups, nil
}

// ListResourceGroupsByTag returns a resource group list within a subscription based on a tag key.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListResourceGroupsByTagContext] instead.
func ListResourceGroupsByTag(t testing.TestingT, tag, subscriptionID string) []armresources.ResourceGroup {
	t.Helper()

	return ListResourceGroupsByTagContext(t, context.Background(), tag, subscriptionID)
}

// ListResourceGroupsByTagE returns a resource group list within a subscription based on a tag key.
//
// Deprecated: Use [ListResourceGroupsByTagContextE] instead.
func ListResourceGroupsByTagE(tag string, subscriptionID string) ([]armresources.ResourceGroup, error) {
	return ListResourceGroupsByTagContextE(context.Background(), tag, subscriptionID)
}

func resourceGroupNotFoundError(err error) bool {
	if err != nil {
		var azcoreErr *azcore.ResponseError
		if errors.As(err, &azcoreErr) {
			return azcoreErr.ErrorCode == "ResourceGroupNotFound"
		}

		var autorestError autorest.DetailedError
		if errors.As(err, &autorestError) {
			var requestError *autorestAzure.RequestError
			if errors.As(autorestError.Original, &requestError) {
				return (requestError.ServiceError.Code == "ResourceGroupNotFound")
			}
		}
	}

	return false
}
