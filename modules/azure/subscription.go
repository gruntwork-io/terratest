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
