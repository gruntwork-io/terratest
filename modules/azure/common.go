package azure

import (
	"fmt"
	"os"
)

const (
	// AzureSubscriptionID is an optional env variable supported by the `azurerm` Terraform provider to
	// designate a target Azure subscription ID
	AzureSubscriptionID = "ARM_SUBSCRIPTION_ID"

	// AzureResGroupName is an optional env variable custom to Terratest to designate a target Azure resource group
	AzureResGroupName = "AZURE_RES_GROUP_NAME"
)

// GetTargetAzureSubscription is a helper function to find the correct target Azure Subscription ID,
// with provided arguments taking precedence over environment variables
func GetTargetAzureSubscription(subscriptionID string) (string, error) {
	return getTargetAzureSubscription(subscriptionID)
}

func getTargetAzureSubscription(subscriptionID string) (string, error) {
	fmt.Printf("Initial subscription ID is %s\n", subscriptionID)
	if subscriptionID == "" {
		if id, exists := os.LookupEnv(AzureSubscriptionID); exists {
			return id, nil
		}

		return "", SubscriptionIDNotFound{}
	}

	fmt.Printf("Final subscription ID is %s\n", subscriptionID)

	return subscriptionID, nil
}

// GetTargetAzureResourceGroupName is a helper function to find the correct target Azure Resource Group name,
// with provided arguments taking precedence over environment variables
func GetTargetAzureResourceGroupName(resourceGroupName string) (string, error) {
	return getTargetAzureResourceGroupName(resourceGroupName)
}

func getTargetAzureResourceGroupName(resourceGroupName string) (string, error) {
	if resourceGroupName == "" {
		if name, exists := os.LookupEnv(AzureResGroupName); exists {
			return name, nil
		}

		return "", ResourceGroupNameNotFound{}
	}

	return resourceGroupName, nil
}
