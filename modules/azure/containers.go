package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance"
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

	return ContainerRegistryExistsContext(t, context.Background(), registryName, resourceGroupName, subscriptionID) //nolint:staticcheck
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

	return GetContainerRegistryContext(t, context.Background(), registryName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetContainerRegistryContextE gets the container registry object.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:dupl
func GetContainerRegistryContextE(ctx context.Context, registryName string, resGroupName string, subscriptionID string) (*armcontainerregistry.Registry, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	subID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	client, err := armcontainerregistry.NewRegistriesClient(subID, cred, opts)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, rgName, registryName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Registry, nil
}

// GetContainerRegistryE gets the container registry object.
//
// Deprecated: Use [GetContainerRegistryContextE] instead.
func GetContainerRegistryE(registryName string, resGroupName string, subscriptionID string) (*armcontainerregistry.Registry, error) {
	return GetContainerRegistryContextE(context.Background(), registryName, resGroupName, subscriptionID)
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

	return ContainerInstanceExistsContext(t, context.Background(), instanceName, resourceGroupName, subscriptionID) //nolint:staticcheck
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

	return GetContainerInstanceContext(t, context.Background(), instanceName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetContainerInstanceContextE gets the container instance object.
// The ctx parameter supports cancellation and timeouts.
//
//nolint:dupl
func GetContainerInstanceContextE(ctx context.Context, instanceName string, resGroupName string, subscriptionID string) (*armcontainerinstance.ContainerGroup, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	subID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	client, err := armcontainerinstance.NewContainerGroupsClient(subID, cred, opts)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, rgName, instanceName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ContainerGroup, nil
}

// GetContainerInstanceE gets the container instance object.
//
// Deprecated: Use [GetContainerInstanceContextE] instead.
func GetContainerInstanceE(instanceName string, resGroupName string, subscriptionID string) (*armcontainerinstance.ContainerGroup, error) {
	return GetContainerInstanceContextE(context.Background(), instanceName, resGroupName, subscriptionID)
}
