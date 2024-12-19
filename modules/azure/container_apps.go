package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/appcontainers/armappcontainers/v3"
	"github.com/stretchr/testify/require"
)

// ManagedEnvironmentExists indicates whether the specified Managed Environment exists.
// This function would fail the test if there is an error.
func ManagedEnvironmentExists(t *testing.T, environmentName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := ManagedEnvironmentExistsE(environmentName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// ManagedEnvironmentExistsE indicates whether the specified Managed Environment exists.
func ManagedEnvironmentExistsE(environmentName string, resourceGroupName string, subscriptionID string) (bool, error) {
	client, err := CreateManagedEnvironmentsClientE(subscriptionID)
	if err != nil {
		return false, err
	}
	_, err = client.Get(context.Background(), resourceGroupName, environmentName, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetManagedEnvironment gets the Managed Environment object
// This function would fail the test if there is an error.
func GetManagedEnvironment(t *testing.T, environmentName string, resourceGroupName string, subscriptionID string) *armappcontainers.ManagedEnvironment {
	env, err := GetManagedEnvironmentE(environmentName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return env
}

// GetManagedEnvironmentE gets the Managed Environment object
func GetManagedEnvironmentE(environmentName string, resourceGroupName string, subscriptionID string) (*armappcontainers.ManagedEnvironment, error) {
	client, err := CreateManagedEnvironmentsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}
	env, err := client.Get(context.Background(), resourceGroupName, environmentName, nil)
	if err != nil {
		return nil, err
	}
	return &env.ManagedEnvironment, nil
}

// ContainerAppExists indicates whether the Container App exists for the subscription.
// This function would fail the test if there is an error.
func ContainerAppExists(t *testing.T, containerAppName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := ContainerAppExistsE(containerAppName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// ContainerAppExistsE indicates whether the Container App exists for the subscription.
func ContainerAppExistsE(containerAppName string, resourceGroupName string, subscriptionID string) (bool, error) {
	client, err := CreateContainerAppsClientE(subscriptionID)
	if err != nil {
		return false, err
	}
	_, err = client.Get(context.Background(), resourceGroupName, containerAppName, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetContainerApp gets the Container App object
// This function would fail the test if there is an error.
func GetContainerApp(t *testing.T, containerAppName string, resourceGroupName string, subscriptionID string) *armappcontainers.ContainerApp {
	app, err := GetContainerAppE(containerAppName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return app
}

// GetContainerAppE gets the Container App object
func GetContainerAppE(environmentName string, resourceGroupName string, subscriptionID string) (*armappcontainers.ContainerApp, error) {
	client, err := CreateContainerAppsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}
	app, err := client.Get(context.Background(), resourceGroupName, environmentName, nil)
	if err != nil {
		return nil, err
	}
	return &app.ContainerApp, nil
}

// ContainerAppJobExists indicates whether the Container App Job exists for the subscription.
// This function would fail the test if there is an error.
func ContainerAppJobExists(t *testing.T, containerAppName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := ContainerAppJobExistsE(containerAppName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// ContainerAppJobExistsE indicates whether the Container App Job exists for the subscription.
func ContainerAppJobExistsE(containerAppName string, resourceGroupName string, subscriptionID string) (bool, error) {
	client, err := CreateContainerAppJobsClientE(subscriptionID)
	if err != nil {
		return false, err
	}
	_, err = client.Get(context.Background(), resourceGroupName, containerAppName, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetContainerAppJob gets the Container App Job object
// This function would fail the test if there is an error.
func GetContainerAppJob(t *testing.T, containerAppName string, resourceGroupName string, subscriptionID string) *armappcontainers.Job {
	app, err := GetContainerAppJobE(containerAppName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return app
}

// GetContainerAppJobE gets the Container App Job object
func GetContainerAppJobE(environmentName string, resourceGroupName string, subscriptionID string) (*armappcontainers.Job, error) {
	client, err := CreateContainerAppJobsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}
	app, err := client.Get(context.Background(), resourceGroupName, environmentName, nil)
	if err != nil {
		return nil, err
	}
	return &app.Job, nil
}
