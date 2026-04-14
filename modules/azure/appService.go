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

// AppExists indicates whether the specified application exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [AppExistsContext] instead.
func AppExists(t testing.TestingT, appName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return AppExistsContext(t, context.Background(), appName, resourceGroupName, subscriptionID)
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

// AppExistsE indicates whether the specified application exists.
//
// Deprecated: Use [AppExistsContextE] instead.
func AppExistsE(appName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return AppExistsContextE(context.Background(), appName, resourceGroupName, subscriptionID)
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

// GetAppService gets the App service object for the specified application.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetAppServiceContext] instead.
func GetAppService(t testing.TestingT, appName string, resGroupName string, subscriptionID string) *armappservice.Site {
	t.Helper()

	return GetAppServiceContext(t, context.Background(), appName, resGroupName, subscriptionID)
}

// GetAppServiceContextE gets the App service object for the specified application.
// The ctx parameter supports cancellation and timeouts.
func GetAppServiceContextE(ctx context.Context, appName string, resGroupName string, subscriptionID string) (*armappservice.Site, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetAppServiceClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetAppServiceWithClient(ctx, client, rgName, appName)
}

// GetAppServiceWithClient gets the App service object using the provided client.
// This variant is useful for testing with fake clients.
func GetAppServiceWithClient(ctx context.Context, client *armappservice.WebAppsClient, resourceGroupName string, appName string) (*armappservice.Site, error) {
	resp, err := client.Get(ctx, resourceGroupName, appName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Site, nil
}

// GetAppServiceE gets the App service object for the specified application.
//
// Deprecated: Use [GetAppServiceContextE] instead.
func GetAppServiceE(appName string, resGroupName string, subscriptionID string) (*armappservice.Site, error) {
	return GetAppServiceContextE(context.Background(), appName, resGroupName, subscriptionID)
}

// GetAppServiceClientContextE creates and returns an App Service web apps client.
// The ctx parameter supports cancellation and timeouts.
func GetAppServiceClientContextE(_ context.Context, subscriptionID string) (*armappservice.WebAppsClient, error) {
	clientFactory, err := getArmAppServiceClientFactory(subscriptionID)
	if err != nil {
		return nil, err
	}

	return clientFactory.NewWebAppsClient(), nil
}

// GetAppServiceClientE creates and returns an App Service web apps client.
//
// Deprecated: Use [GetAppServiceClientContextE] instead.
func GetAppServiceClientE(subscriptionID string) (*armappservice.WebAppsClient, error) {
	return GetAppServiceClientContextE(context.Background(), subscriptionID)
}
