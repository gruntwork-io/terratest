package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3"
	"github.com/gruntwork-io/terratest/modules/testing"
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

// ManagedEnvironmentExists indicates whether the specified Managed Environment exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ManagedEnvironmentExistsContext] instead.
func ManagedEnvironmentExists(t testing.TestingT, environmentName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return ManagedEnvironmentExistsContext(t, context.Background(), environmentName, resourceGroupName, subscriptionID) //nolint:staticcheck
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

// ManagedEnvironmentExistsE indicates whether the specified Managed Environment exists.
//
// Deprecated: Use [ManagedEnvironmentExistsContextE] instead.
func ManagedEnvironmentExistsE(environmentName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return ManagedEnvironmentExistsContextE(context.Background(), environmentName, resourceGroupName, subscriptionID)
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

// GetManagedEnvironment gets the Managed Environment object.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetManagedEnvironmentContext] instead.
func GetManagedEnvironment(t testing.TestingT, environmentName string, resourceGroupName string, subscriptionID string) *armappcontainers.ManagedEnvironment {
	t.Helper()

	return GetManagedEnvironmentContext(t, context.Background(), environmentName, resourceGroupName, subscriptionID) //nolint:staticcheck
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

// GetManagedEnvironmentE gets the Managed Environment object.
//
// Deprecated: Use [GetManagedEnvironmentContextE] instead.
func GetManagedEnvironmentE(environmentName string, resourceGroupName string, subscriptionID string) (*armappcontainers.ManagedEnvironment, error) {
	return GetManagedEnvironmentContextE(context.Background(), environmentName, resourceGroupName, subscriptionID)
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

// ContainerAppExists indicates whether the Container App exists for the subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ContainerAppExistsContext] instead.
func ContainerAppExists(t testing.TestingT, containerAppName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return ContainerAppExistsContext(t, context.Background(), containerAppName, resourceGroupName, subscriptionID) //nolint:staticcheck
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

// ContainerAppExistsE indicates whether the Container App exists for the subscription.
//
// Deprecated: Use [ContainerAppExistsContextE] instead.
func ContainerAppExistsE(containerAppName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return ContainerAppExistsContextE(context.Background(), containerAppName, resourceGroupName, subscriptionID)
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

// GetContainerApp gets the Container App object.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetContainerAppContext] instead.
func GetContainerApp(t testing.TestingT, containerAppName string, resourceGroupName string, subscriptionID string) *armappcontainers.ContainerApp {
	t.Helper()

	return GetContainerAppContext(t, context.Background(), containerAppName, resourceGroupName, subscriptionID) //nolint:staticcheck
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

// GetContainerAppE gets the Container App object.
//
// Deprecated: Use [GetContainerAppContextE] instead.
func GetContainerAppE(containerAppName string, resourceGroupName string, subscriptionID string) (*armappcontainers.ContainerApp, error) {
	return GetContainerAppContextE(context.Background(), containerAppName, resourceGroupName, subscriptionID)
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

// ContainerAppJobExists indicates whether the Container App Job exists for the subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ContainerAppJobExistsContext] instead.
func ContainerAppJobExists(t testing.TestingT, containerAppName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return ContainerAppJobExistsContext(t, context.Background(), containerAppName, resourceGroupName, subscriptionID) //nolint:staticcheck
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

// ContainerAppJobExistsE indicates whether the Container App Job exists for the subscription.
//
// Deprecated: Use [ContainerAppJobExistsContextE] instead.
func ContainerAppJobExistsE(containerAppName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return ContainerAppJobExistsContextE(context.Background(), containerAppName, resourceGroupName, subscriptionID)
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

// GetContainerAppJob gets the Container App Job object.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetContainerAppJobContext] instead.
func GetContainerAppJob(t testing.TestingT, containerAppName string, resourceGroupName string, subscriptionID string) *armappcontainers.Job {
	t.Helper()

	return GetContainerAppJobContext(t, context.Background(), containerAppName, resourceGroupName, subscriptionID) //nolint:staticcheck
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

// GetContainerAppJobE gets the Container App Job object.
//
// Deprecated: Use [GetContainerAppJobContextE] instead.
func GetContainerAppJobE(containerAppName string, resourceGroupName string, subscriptionID string) (*armappcontainers.Job, error) {
	return GetContainerAppJobContextE(context.Background(), containerAppName, resourceGroupName, subscriptionID)
}
