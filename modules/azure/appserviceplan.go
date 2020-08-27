package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/web/mgmt/web"
	"github.com/stretchr/testify/require"
)

// GetAppServicePlan gets the AppServicePlan.
// planName - required to find the AppServicePlan.
// resGroupName - use an empty string if you have the AZURE_RES_GROUP_NAME environment variable set
// subscriptionId - use an empty string if you have the ARM_SUBSCRIPTION_ID environment variable set
func GetAppServicePlan(t *testing.T, planName string, resGroupName string, subscriptionID string) *web.AppServicePlan {
	plan, err := getAppServicePlanE(planName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return plan
}

func getAppServicePlanE(planName string, resGroupName string, subscriptionID string) (*web.AppServicePlan, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := getAppServicePlanClient(subscriptionID)
	if err != nil {
		return nil, err
	}

	plan, err := client.Get(context.Background(), rgName, planName)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func getAppServicePlanClient(subscriptionID string) (*web.AppServicePlansClient, error) {
	subID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an AppServicePlanClient
	planClient := web.NewAppServicePlansClient(subID)

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	planClient.Authorizer = *authorizer

	return &planClient, nil
}
