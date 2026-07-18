package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// ManagedEnvironmentExistsContext indicates whether the specified Managed Environment exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ManagedEnvironmentExistsContext(t testing.TestingT, ctx context.Context, environmentName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := ManagedEnvironmentExistsContextE(ctx, environmentName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// ManagedEnvironmentExistsContextE indicates whether the specified Managed Environment exists.
// The ctx parameter supports cancellation and timeouts.
func ManagedEnvironmentExistsContextE(ctx context.Context, environmentName string, resourceGroupName string, subscriptionID string) (bool, error) {
	client, err := CreateManagedEnvironmentsClientContextE(ctx, subscriptionID)
	if err != nil {
		return false, err
	}

	_, err = client.Get(ctx, resourceGroupName, environmentName, nil)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetManagedEnvironmentContext returns the Managed Environment object.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetManagedEnvironmentContext(t testing.TestingT, ctx context.Context, environmentName string, resourceGroupName string, subscriptionID string) *armappcontainers.ManagedEnvironment {
	t.Helper()

	env, err := GetManagedEnvironmentContextE(ctx, environmentName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return env
}

// GetManagedEnvironmentContextE returns the Managed Environment object.
// The ctx parameter supports cancellation and timeouts.
func GetManagedEnvironmentContextE(ctx context.Context, environmentName string, resourceGroupName string, subscriptionID string) (*armappcontainers.ManagedEnvironment, error) {
	client, err := CreateManagedEnvironmentsClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetManagedEnvironmentWithClient(ctx, client, resourceGroupName, environmentName)
}

// GetManagedEnvironmentWithClient returns a Managed Environment using the provided ManagedEnvironmentsClient.
// This variant is useful for testing with fake clients.
func GetManagedEnvironmentWithClient(ctx context.Context, client *armappcontainers.ManagedEnvironmentsClient, resourceGroupName string, environmentName string) (*armappcontainers.ManagedEnvironment, error) {
	resp, err := client.Get(ctx, resourceGroupName, environmentName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ManagedEnvironment, nil
}

// ContainerAppExistsContext indicates whether the Container App exists for the subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ContainerAppExistsContext(t testing.TestingT, ctx context.Context, containerAppName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := ContainerAppExistsContextE(ctx, containerAppName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// ContainerAppExistsContextE indicates whether the Container App exists for the subscription.
// The ctx parameter supports cancellation and timeouts.
func ContainerAppExistsContextE(ctx context.Context, containerAppName string, resourceGroupName string, subscriptionID string) (bool, error) {
	client, err := CreateContainerAppsClientContextE(ctx, subscriptionID)
	if err != nil {
		return false, err
	}

	_, err = client.Get(ctx, resourceGroupName, containerAppName, nil)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetContainerAppContext returns the Container App object.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetContainerAppContext(t testing.TestingT, ctx context.Context, containerAppName string, resourceGroupName string, subscriptionID string) *armappcontainers.ContainerApp {
	t.Helper()

	app, err := GetContainerAppContextE(ctx, containerAppName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return app
}

// GetContainerAppContextE returns the Container App object.
// The ctx parameter supports cancellation and timeouts.
func GetContainerAppContextE(ctx context.Context, containerAppName string, resourceGroupName string, subscriptionID string) (*armappcontainers.ContainerApp, error) {
	client, err := CreateContainerAppsClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetContainerAppWithClient(ctx, client, resourceGroupName, containerAppName)
}

// GetContainerAppWithClient returns a Container App using the provided ContainerAppsClient.
// This variant is useful for testing with fake clients.
func GetContainerAppWithClient(ctx context.Context, client *armappcontainers.ContainerAppsClient, resourceGroupName string, containerAppName string) (*armappcontainers.ContainerApp, error) {
	resp, err := client.Get(ctx, resourceGroupName, containerAppName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ContainerApp, nil
}

// ContainerAppJobExistsContext indicates whether the Container App Job exists for the subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ContainerAppJobExistsContext(t testing.TestingT, ctx context.Context, containerAppName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := ContainerAppJobExistsContextE(ctx, containerAppName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// ContainerAppJobExistsContextE indicates whether the Container App Job exists for the subscription.
// The ctx parameter supports cancellation and timeouts.
func ContainerAppJobExistsContextE(ctx context.Context, containerAppName string, resourceGroupName string, subscriptionID string) (bool, error) {
	client, err := CreateContainerAppJobsClientContextE(ctx, subscriptionID)
	if err != nil {
		return false, err
	}

	_, err = client.Get(ctx, resourceGroupName, containerAppName, nil)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetContainerAppJobContext returns the Container App Job object.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetContainerAppJobContext(t testing.TestingT, ctx context.Context, containerAppName string, resourceGroupName string, subscriptionID string) *armappcontainers.Job {
	t.Helper()

	app, err := GetContainerAppJobContextE(ctx, containerAppName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return app
}

// GetContainerAppJobContextE returns the Container App Job object.
// The ctx parameter supports cancellation and timeouts.
func GetContainerAppJobContextE(ctx context.Context, containerAppName string, resourceGroupName string, subscriptionID string) (*armappcontainers.Job, error) {
	client, err := CreateContainerAppJobsClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetContainerAppJobWithClient(ctx, client, resourceGroupName, containerAppName)
}

// GetContainerAppJobWithClient returns a Container App Job using the provided JobsClient.
// This variant is useful for testing with fake clients.
func GetContainerAppJobWithClient(ctx context.Context, client *armappcontainers.JobsClient, resourceGroupName string, containerAppName string) (*armappcontainers.Job, error) {
	resp, err := client.Get(ctx, resourceGroupName, containerAppName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Job, nil
}
