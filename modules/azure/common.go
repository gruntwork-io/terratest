package azure

import (
	"os"
)

const (
	// AzureSubscriptionID is an optional env variable supported by the `azurerm` Terraform provider to
	// designate a target Azure subscription ID
	AzureSubscriptionID = "ARM_SUBSCRIPTION_ID"

	// AzureResGroupName is an optional env variable custom to Terratest to designate a target Azure resource group
	AzureResGroupName = "AZURE_RES_GROUP_NAME"
)

func getTargetAzureSubscription(subscriptionID string) (string, error) {
	if subscriptionID == "" {
		if id, exists := os.LookupEnv(AzureSubscriptionID); exists {
			return id, nil
		}

		return "", SubscriptionIDNotFound{}
	}

	return subscriptionID, nil
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

// safePtrToString converts a string pointer to a non-pointer string value, or to "" if the pointer is nil.
func safePtrToString(raw *string) string {
	if raw == nil {
		return ""
	}

	return *raw
}

// safePtrToInt32 converts a int32 pointer to a non-pointer int32 value, or to 0 if the pointer is nil.
func safePtrToInt32(raw *int32) int32 {
	if raw == nil {
		return 0
	}

	return *raw
}

// safePtrToList converts a []*string slice to a []string slice, dereferencing each element.
// Returns an empty slice if the input is nil.
func safePtrToList(raw []*string) []string {
	if raw == nil {
		return []string{}
	}

	result := make([]string, len(raw))
	for i, s := range raw {
		if s != nil {
			result[i] = *s
		}
	}

	return result
}
