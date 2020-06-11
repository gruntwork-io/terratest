package azure

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetSkuOfAppServicePlan returns the SkuDescription from the App Service Plan
func GetSkuOfAppServicePlan(t testing.TestingT, planName string, resGroupName string, subscriptionID string) *web.SkuDescription {
	plan := GetAppServicePlan(t, planName, resGroupName, subscriptionID)

	return plan.Sku
}

// GetTagsOfAppServicePlan gets the tags of the given App Service Plan as a map
func GetTagsOfAppServicePlan(t testing.TestingT, planName string, resGroupName string, subscriptionID string) map[string]string {
	plan := GetAppServicePlan(t, planName, resGroupName, subscriptionID)

	tags := make(map[string]string)
	for k, v := range plan.Tags {
		tags[k] = *v
	}

	return tags
}

// GetAppServicePlan gets the AppServicePlan after checking for a
// valid SubscriptionID, Resource Group Name, and authenticating the App Service Plan Client
func GetAppServicePlan(t testing.TestingT, planName string, resGroupName string, subscriptionID string) *web.AppServicePlan {
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