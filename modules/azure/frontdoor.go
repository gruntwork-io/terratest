package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/frontdoor/mgmt/frontdoor"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// FrontDoorExists indicates whether the Front Door exists for the subscription.
// This function would fail the test if there is an error.
func FrontDoorExists(t testing.TestingT, frontDoorName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := FrontDoorExistsE(frontDoorName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// GetFrontDoor gets a Front Door by name if it exists for the subscription.
// This function would fail the test if there is an error.
func GetFrontDoor(t testing.TestingT, frontDoorName string, resourceGroupName string, subscriptionID string) *frontdoor.FrontDoor {
	fd, err := GetFrontDoorE(frontDoorName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return fd
}

// FrontDoorFrontendEndpointExists indicates whether the frontend endpoint exists for the provided Front Door.
// This function would fail the test if there is an error.
func FrontDoorFrontendEndpointExists(t testing.TestingT, endpointName string, frontDoorName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := FrontDoorFrontendEndpointExistsE(endpointName, frontDoorName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// GetFrontDoorFrontendEndpoint gets a frontend endpoint by name for the provided Front Door if it exists for the subscription.
// This function would fail the test if there is an error.
func GetFrontDoorFrontendEndpoint(t testing.TestingT, endpointName string, frontDoorName string, resourceGroupName string, subscriptionID string) *frontdoor.FrontendEndpoint {
	ep, err := GetFrontDoorFrontendEndpointE(endpointName, frontDoorName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return ep
}

// FrontDoorExistsE indicates whether the specified Front Door exists and may return an error.
func FrontDoorExistsE(frontDoorName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetFrontDoorE(frontDoorName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// FrontDoorFrontendEndpointExistsE indicates whether the specified endpoint exists for the provided Front Door and may return an error.
func FrontDoorFrontendEndpointExistsE(endpointName string, frontDoorName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetFrontDoorFrontendEndpointE(endpointName, frontDoorName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetFrontDoorE gets the specified Front Door if it exists and may return an error.
func GetFrontDoorE(frontDoorName, resoureGroupName, subscriptionID string) (*frontdoor.FrontDoor, error) {
	client, err := GetFrontDoorClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	fd, err := client.Get(context.Background(), resoureGroupName, frontDoorName)
	if err != nil {
		return nil, err
	}

	return &fd, nil
}

// GetFrontDoorFrontendEndpointE gets the specified Frontend Endpoint for the provided Front Door if it exists and may return an error.
func GetFrontDoorFrontendEndpointE(endpointName, frontDoorName, resourceGroupName, subscriptionID string) (*frontdoor.FrontendEndpoint, error) {
	client, err := GetFrontDoorFrontendEndpointClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	endpoint, err := client.Get(context.Background(), resourceGroupName, frontDoorName, endpointName)
	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}

// GetFrontDoorClientE return a front door client; otherwise error.
func GetFrontDoorClientE(subscriptionID string) (*frontdoor.FrontDoorsClient, error) {
	client, err := CreateFrontDoorClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer
	return client, nil
}

// GetFrontDoorFrontendEndpointClientE returns a front door frontend endpoints client; otherwise error.
func GetFrontDoorFrontendEndpointClientE(subscriptionID string) (*frontdoor.FrontendEndpointsClient, error) {
	client, err := CreateFrontDoorFrontendEndpointClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer
	return client, nil
}
