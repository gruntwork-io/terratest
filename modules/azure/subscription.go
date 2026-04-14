package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions"
)

// CreateSubscriptionsClientContextE returns a subscriptions client.
// The ctx parameter supports cancellation and timeouts.
func CreateSubscriptionsClientContextE(_ context.Context) (*armsubscriptions.Client, error) {
	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	return armsubscriptions.NewClient(cred, opts)
}

// GetSubscriptionClientContextE is a helper function that will setup an Azure Subscription client on your behalf.
// The ctx parameter supports cancellation and timeouts.
//
// Deprecated: Use [CreateSubscriptionsClientContextE] instead.
func GetSubscriptionClientContextE(ctx context.Context) (*armsubscriptions.Client, error) {
	return CreateSubscriptionsClientContextE(ctx)
}

// GetSubscriptionClientE is a helper function that will setup an Azure Subscription client on your behalf.
//
// Deprecated: Use [CreateSubscriptionsClientContextE] instead.
func GetSubscriptionClientE() (*armsubscriptions.Client, error) {
	return CreateSubscriptionsClientContextE(context.Background())
}
