package azure

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/containerinstance/mgmt/2018-10-01/containerinstance"
	"github.com/Azure/azure-sdk-for-go/services/containerregistry/mgmt/2019-05-01/containerregistry"

	"github.com/stretchr/testify/require"
)

// ContainerRegistryExists indicates whether the specified container registry exists.
// This function would fail the test if there is an error.
func ContainerRegistryExists(t *testing.T, registryName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := ContainerRegistryExistsE(registryName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// ContainerRegistryExistsE indicates whether the specified container registry exists.
func ContainerRegistryExistsE(registryName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetContainerRegistryE(registryName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetContainerRegistry gets the container registry object
// This function would fail the test if there is an error.
func GetContainerRegistry(t *testing.T, registryName string, resGroupName string, subscriptionID string) *containerregistry.Registry {
	resource, err := GetContainerRegistryE(registryName, resGroupName, subscriptionID)

	require.NoError(t, err)

	return resource
}

// GetContainerRegistryE gets the container registry object
func GetContainerRegistryE(registryName string, resGroupName string, subscriptionID string) (*containerregistry.Registry, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetContainerRegistryClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resource, err := client.Get(context.Background(), rgName, registryName)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

// GetContainerRegistryClientE is a helper function that will setup an Azure Container Registry client on your behalf
func GetContainerRegistryClientE(subscriptionID string) (*containerregistry.RegistriesClient, error) {
	// Create an ACR client
	registryClient, err := CreateContainerRegistryClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	registryClient.Authorizer = *authorizer
	return registryClient, nil
}

// ContainerInstanceExists indicates whether the specified container instance exists.
// This function would fail the test if there is an error.
func ContainerInstanceExists(t *testing.T, instanceName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := ContainerInstanceExistsE(instanceName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// ContainerInstanceExistsE indicates whether the specified container instance exists.
func ContainerInstanceExistsE(instanceName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetContainerInstanceE(instanceName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetContainerInstance gets the container instance object
// This function would fail the test if there is an error.
func GetContainerInstance(t *testing.T, instanceName string, resGroupName string, subscriptionID string) *containerinstance.ContainerGroup {
	instance, err := GetContainerInstanceE(instanceName, resGroupName, subscriptionID)

	require.NoError(t, err)

	return instance
}

// GetContainerInstanceE gets the container instance object
func GetContainerInstanceE(instanceName string, resGroupName string, subscriptionID string) (*containerinstance.ContainerGroup, error) {
	rgName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetContainerInstanceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	instance, err := client.Get(context.Background(), rgName, instanceName)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

// GetContainerInstanceClientE is a helper function that will setup an Azure Container Instance client on your behalf
func GetContainerInstanceClientE(subscriptionID string) (*containerinstance.ContainerGroupsClient, error) {
	// Create an ACI client
	instanceClient, err := CreateContainerInstanceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	instanceClient.Authorizer = *authorizer
	return instanceClient, nil
}
