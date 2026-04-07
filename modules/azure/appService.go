package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appservice/armappservice/v2"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// AppExistsContext indicates whether the specified application exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func AppExistsContext(t testing.TestingT, ctx context.Context, appName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := AppExistsContextE(ctx, appName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// AppExistsContextE indicates whether the specified application exists.
// The ctx parameter supports cancellation and timeouts.
func AppExistsContextE(ctx context.Context, appName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetAppServiceContextE(ctx, appName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetAppServiceContext gets the App service object for the specified application.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAppServiceContext(t testing.TestingT, ctx context.Context, appName string, resGroupName string, subscriptionID string) *armappservice.Site {
	t.Helper()

	site, err := GetAppServiceContextE(ctx, appName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return site
}

// GetAppServiceContextE gets the App service object for the specified application.
// The ctx parameter supports cancellation and timeouts.
func GetAppServiceContextE(ctx context.Context, appName string, resGroupName string, subscriptionID string) (*armappservice.Site, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetAppServiceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, rgName, appName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Site, nil
}

// GetAppServiceClientE creates and returns an App Service web apps client.
func GetAppServiceClientE(subscriptionID string) (*armappservice.WebAppsClient, error) {
	clientFactory, err := getArmAppServiceClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewWebAppsClient(), nil
}
