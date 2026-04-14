package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerregistry/armcontainerregistry"
	"github.com/gruntwork-io/terratest/modules/testing"

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

// ContainerRegistryExists indicates whether the specified container registry exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ContainerRegistryExistsContext] instead.
func ContainerRegistryExists(t testing.TestingT, registryName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return ContainerRegistryExistsContext(t, context.Background(), registryName, resourceGroupName, subscriptionID)
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

// ContainerRegistryExistsE indicates whether the specified container registry exists.
//
// Deprecated: Use [ContainerRegistryExistsContextE] instead.
func ContainerRegistryExistsE(registryName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return ContainerRegistryExistsContextE(context.Background(), registryName, resourceGroupName, subscriptionID)
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

// GetContainerRegistry gets the container registry object.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetContainerRegistryContext] instead.
func GetContainerRegistry(t testing.TestingT, registryName string, resGroupName string, subscriptionID string) *armcontainerregistry.Registry {
	t.Helper()

	return GetContainerRegistryContext(t, context.Background(), registryName, resGroupName, subscriptionID)
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

// GetContainerRegistryE gets the container registry object.
//
// Deprecated: Use [GetContainerRegistryContextE] instead.
func GetContainerRegistryE(registryName string, resGroupName string, subscriptionID string) (*armcontainerregistry.Registry, error) {
	return GetContainerRegistryContextE(context.Background(), registryName, resGroupName, subscriptionID)
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

// ContainerInstanceExists indicates whether the specified container instance exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ContainerInstanceExistsContext] instead.
func ContainerInstanceExists(t testing.TestingT, instanceName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return ContainerInstanceExistsContext(t, context.Background(), instanceName, resourceGroupName, subscriptionID)
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

// ContainerInstanceExistsE indicates whether the specified container instance exists.
//
// Deprecated: Use [ContainerInstanceExistsContextE] instead.
func ContainerInstanceExistsE(instanceName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return ContainerInstanceExistsContextE(context.Background(), instanceName, resourceGroupName, subscriptionID)
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

// GetContainerInstance gets the container instance object.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetContainerInstanceContext] instead.
func GetContainerInstance(t testing.TestingT, instanceName string, resGroupName string, subscriptionID string) *armcontainerinstance.ContainerGroup {
	t.Helper()

	return GetContainerInstanceContext(t, context.Background(), instanceName, resGroupName, subscriptionID)
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

// GetContainerInstanceE gets the container instance object.
//
// Deprecated: Use [GetContainerInstanceContextE] instead.
func GetContainerInstanceE(instanceName string, resGroupName string, subscriptionID string) (*armcontainerinstance.ContainerGroup, error) {
	return GetContainerInstanceContextE(context.Background(), instanceName, resGroupName, subscriptionID)
}

// GetContainerInstanceWithClient gets a container instance using the provided ContainerGroupsClient.
func GetContainerInstanceWithClient(ctx context.Context, client *armcontainerinstance.ContainerGroupsClient, resGroupName string, instanceName string) (*armcontainerinstance.ContainerGroup, error) {
	resp, err := client.Get(ctx, resGroupName, instanceName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ContainerGroup, nil
}
