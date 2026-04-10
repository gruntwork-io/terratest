package azure

import (
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

// GetSubscriptionClientE is a helper function that will setup an Azure Subscription client on your behalf.
func GetSubscriptionClientE() (*armsubscriptions.Client, error) {
	return CreateSubscriptionsClientE()
}
