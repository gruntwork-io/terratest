package azure

import (
	"context"

	//"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/Azure/azure-sdk-for-go/services/eventhub/mgmt/2017-04-01/eventhub"
	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetPostgreSQLServerClientE is a helper function that will setup a postgresql server client.
func GetEventHubClientE(subscriptionID string) (*eventhub.EventHubsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a postgresql server client
	eventHubClient := eventhub.NewEventHubsClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	eventHubClient.Authorizer = *authorizer

	return &eventHubClient, nil
}

// GetPostgreSQLServerE is a helper function that gets the server.
func GetEventHubNamespaceE(t testing.TestingT, subscriptionID string, resGroupName string, namespace string) (*.Server, error) {
	// Create a postgresql Server client
	eventHubClient, err := GetEventHubClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding server client
	postgresqlServer, err := eventHubClient.Get(context.Background(), resGroupName, namespace)
	if err != nil {
		return nil, err
	}

	// TODO: temp
	require.NoError(t, err)

	//Return server
	return &postgresqlServer, nil
}
