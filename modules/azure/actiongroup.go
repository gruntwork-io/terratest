package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetActionGroupResourceContext gets the ActionGroupResource.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
// ruleName - required to find the ActionGroupResource.
// resGroupName - use an empty string if you have the AZURE_RES_GROUP_NAME environment variable set
// subscriptionId - use an empty string if you have the ARM_SUBSCRIPTION_ID environment variable set
func GetActionGroupResourceContext(t testing.TestingT, ctx context.Context, ruleName string, resGroupName string, subscriptionID string) *armmonitor.ActionGroupResource {
	actionGroupResource, err := GetActionGroupResourceContextE(ctx, ruleName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return actionGroupResource
}

// GetActionGroupResourceContextE gets the ActionGroupResource.
// The ctx parameter supports cancellation and timeouts.
// ruleName - required to find the ActionGroupResource.
// resGroupName - use an empty string if you have the AZURE_RES_GROUP_NAME environment variable set
// subscriptionId - use an empty string if you have the ARM_SUBSCRIPTION_ID environment variable set
func GetActionGroupResourceContextE(ctx context.Context, ruleName string, resGroupName string, subscriptionID string) (*armmonitor.ActionGroupResource, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := CreateActionGroupClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, rgName, ruleName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ActionGroupResource, nil
}

