package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/stretchr/testify/require"
)

// AppExists indicates whether the specified application exists.
// This function would fail the test if there is an error.
func AppExists(t *testing.T, appName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := AppExistsE(appName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// AppExistsE indicates whether the specified application exists.
func AppExistsE(appName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetAppServiceE(appName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetAppService gets the App service object
// This function would fail the test if there is an error.
func GetAppService(t *testing.T, appName string, resGroupName string, subscriptionID string) *armappservice.Site {
	site, err := GetAppServiceE(appName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return site
}

// GetAppServiceE gets the App service object
func GetAppServiceE(appName string, resGroupName string, subscriptionID string) (*armappservice.Site, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetAppServiceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(context.Background(), rgName, appName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Site, nil
}

// GetAppServiceClientE creates and returns an App Service web apps client
func GetAppServiceClientE(subscriptionID string) (*armappservice.WebAppsClient, error) {
	clientFactory, err := getArmAppServiceClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}
	return clientFactory.NewWebAppsClient(), nil
}
