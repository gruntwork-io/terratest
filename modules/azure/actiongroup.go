package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetActionGroupResourceContext gets the ActionGroupResource.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
// ruleName - required to find the ActionGroupResource.
// resGroupName - use an empty string if you have the AZURE_RES_GROUP_NAME environment variable set
// subscriptionId - use an empty string if you have the ARM_SUBSCRIPTION_ID environment variable set
func GetActionGroupResourceContext(t testing.TestingT, ctx context.Context, ruleName string, resGroupName string, subscriptionID string) *insights.ActionGroupResource {
	actionGroupResource, err := GetActionGroupResourceContextE(ctx, ruleName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return actionGroupResource
}

// GetActionGroupResourceContextE gets the ActionGroupResource.
// The ctx parameter supports cancellation and timeouts.
// ruleName - required to find the ActionGroupResource.
// resGroupName - use an empty string if you have the AZURE_RES_GROUP_NAME environment variable set
// subscriptionId - use an empty string if you have the ARM_SUBSCRIPTION_ID environment variable set
func GetActionGroupResourceContextE(ctx context.Context, ruleName string, resGroupName string, subscriptionID string) (*insights.ActionGroupResource, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := CreateActionGroupClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	actionGroup, err := client.Get(ctx, rgName, ruleName)
	if err != nil {
		return nil, err
	}

	return &actionGroup, nil
}

// GetActionGroupResource gets the ActionGroupResource.
// This function would fail the test if there is an error.
// ruleName - required to find the ActionGroupResource.
// resGroupName - use an empty string if you have the AZURE_RES_GROUP_NAME environment variable set
// subscriptionId - use an empty string if you have the ARM_SUBSCRIPTION_ID environment variable set
//
// Deprecated: Use [GetActionGroupResourceContext] instead.
func GetActionGroupResource(t testing.TestingT, ruleName string, resGroupName string, subscriptionID string) *insights.ActionGroupResource {
	return GetActionGroupResourceContext(t, context.Background(), ruleName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetActionGroupResourceE gets the ActionGroupResource.
// ruleName - required to find the ActionGroupResource.
// resGroupName - use an empty string if you have the AZURE_RES_GROUP_NAME environment variable set
// subscriptionId - use an empty string if you have the ARM_SUBSCRIPTION_ID environment variable set
//
// Deprecated: Use [GetActionGroupResourceContextE] instead.
func GetActionGroupResourceE(ruleName string, resGroupName string, subscriptionID string) (*insights.ActionGroupResource, error) {
	return GetActionGroupResourceContextE(context.Background(), ruleName, resGroupName, subscriptionID)
}
