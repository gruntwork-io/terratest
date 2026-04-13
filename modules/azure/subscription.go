package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-06-01/subscriptions"
)

// GetSubscriptionClientContextE is a helper function that will setup an Azure Subscription client on your behalf.
// The ctx parameter supports cancellation and timeouts.
func GetSubscriptionClientContextE(ctx context.Context) (*subscriptions.Client, error) {
	// Create a Subscription client
	client, err := CreateSubscriptionsClientContextE(ctx)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	client.Authorizer = *authorizer

	return &client, nil
}

// GetSubscriptionClientE is a helper function that will setup an Azure Subscription client on your behalf.
//
// Deprecated: Use [GetSubscriptionClientContextE] instead.
func GetSubscriptionClientE() (*subscriptions.Client, error) {
	return GetSubscriptionClientContextE(context.Background())
}
