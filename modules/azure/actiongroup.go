package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/preview/preview/monitor/mgmt/insights"
	"github.com/stretchr/testify/require"
)

// GetActionGroupResource gets the ActionGroupResource.
// ruleName - required to find the ActionGroupResource.
// resGroupName - use an empty string if you have the AZURE_RES_GROUP_NAME environment variable set
// subscriptionId - use an empty string if you have the ARM_SUBSCRIPTION_ID environment variable set
func GetActionGroupResource(t *testing.T, ruleName string, resGroupName string, subscriptionID string) *insights.ActionGroupResource {
	actionGroupResource, err := GetActionGroupResourceE(ruleName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return actionGroupResource
}

// GetActionGroupResourceE gets the ActionGroupResource with Error details on error.
// ruleName - required to find the ActionGroupResource.
// resGroupName - use an empty string if you have the AZURE_RES_GROUP_NAME environment variable set
// subscriptionId - use an empty string if you have the ARM_SUBSCRIPTION_ID environment variable set
func GetActionGroupResourceE(ruleName string, resGroupName string, subscriptionID string) (*insights.ActionGroupResource, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := CreateActionGroupClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	actionGroup, err := client.Get(context.Background(), rgName, ruleName)
	if err != nil {
		return nil, err
	}

	return &actionGroup, nil
}

// TODO: remove in next version
func getActionGroupClient(subscriptionID string) (*insights.ActionGroupsClient, error) {
	subID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	metricAlertsClient := insights.NewActionGroupsClient(subID)

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	metricAlertsClient.Authorizer = *authorizer

	return &metricAlertsClient, nil
}
