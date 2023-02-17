package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2019-08-01/web"
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
func GetAppService(t *testing.T, appName string, resGroupName string, subscriptionID string) *web.Site {
	site, err := GetAppServiceE(appName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return site
}

// GetAppServiceE gets the App service object
func GetAppServiceE(appName string, resGroupName string, subscriptionID string) (*web.Site, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetAppServiceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resource, err := client.Get(context.Background(), rgName, appName)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

func GetAppServiceClientE(subscriptionID string) (*web.AppsClient, error) {
	// Create an Apps client
	appsClient, err := CreateAppServiceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	appsClient.Authorizer = *authorizer
	return appsClient, nil
}
