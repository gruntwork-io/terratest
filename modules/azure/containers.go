package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"

	"github.com/stretchr/testify/require"
)

// ContainerRegistryExistsContext indicates whether the specified container registry exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ContainerRegistryExistsContext(t testing.TestingT, ctx context.Context, registryName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := ContainerRegistryExistsContextE(ctx, registryName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// ContainerRegistryExistsContextE indicates whether the specified container registry exists.
// The ctx parameter supports cancellation and timeouts.
func ContainerRegistryExistsContextE(ctx context.Context, registryName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetContainerRegistryContextE(ctx, registryName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetContainerRegistryContext gets the container registry object.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetContainerRegistryContext(t testing.TestingT, ctx context.Context, registryName string, resGroupName string, subscriptionID string) *armcontainerregistry.Registry {
	t.Helper()

	resource, err := GetContainerRegistryContextE(ctx, registryName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return resource
}

// GetContainerRegistryContextE gets the container registry object.
// The ctx parameter supports cancellation and timeouts.
func GetContainerRegistryContextE(ctx context.Context, registryName string, resGroupName string, subscriptionID string) (*armcontainerregistry.Registry, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := CreateContainerRegistryClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetContainerRegistryWithClient(ctx, client, rgName, registryName)
}

// GetContainerRegistryWithClient gets a container registry using the provided RegistriesClient.
func GetContainerRegistryWithClient(ctx context.Context, client *armcontainerregistry.RegistriesClient, resGroupName string, registryName string) (*armcontainerregistry.Registry, error) {
	resp, err := client.Get(ctx, resGroupName, registryName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Registry, nil
}

// ContainerInstanceExistsContext indicates whether the specified container instance exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ContainerInstanceExistsContext(t testing.TestingT, ctx context.Context, instanceName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := ContainerInstanceExistsContextE(ctx, instanceName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// ContainerInstanceExistsContextE indicates whether the specified container instance exists.
// The ctx parameter supports cancellation and timeouts.
func ContainerInstanceExistsContextE(ctx context.Context, instanceName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetContainerInstanceContextE(ctx, instanceName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetContainerInstanceContext gets the container instance object.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetContainerInstanceContext(t testing.TestingT, ctx context.Context, instanceName string, resGroupName string, subscriptionID string) *armcontainerinstance.ContainerGroup {
	t.Helper()

	instance, err := GetContainerInstanceContextE(ctx, instanceName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return instance
}

// GetContainerInstanceContextE gets the container instance object.
// The ctx parameter supports cancellation and timeouts.
func GetContainerInstanceContextE(ctx context.Context, instanceName string, resGroupName string, subscriptionID string) (*armcontainerinstance.ContainerGroup, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := CreateContainerInstanceClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetContainerInstanceWithClient(ctx, client, rgName, instanceName)
}

// GetContainerInstanceWithClient gets a container instance using the provided ContainerGroupsClient.
func GetContainerInstanceWithClient(ctx context.Context, client *armcontainerinstance.ContainerGroupsClient, resGroupName string, instanceName string) (*armcontainerinstance.ContainerGroup, error) {
	resp, err := client.Get(ctx, resGroupName, instanceName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ContainerGroup, nil
}
