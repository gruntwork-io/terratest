package azure

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/require"
)

// StorageAccountExists indicates whether the storage account name exactly matches; otherwise false.
// This function would fail the test if there is an error.
func StorageAccountExists(t *testing.T, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	result, err := StorageAccountExistsE(storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// StorageBlobContainerExists returns true if the container name exactly matches; otherwise false
// This function would fail the test if there is an error.
func StorageBlobContainerExists(t *testing.T, containerName string, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	result, err := StorageBlobContainerExistsE(containerName, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// GetStorageBlobContainerPublicAccess indicates whether a storage container has public access; otherwise false.
// This function would fail the test if there is an error.
func GetStorageBlobContainerPublicAccess(t *testing.T, containerName string, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	result, err := GetStorageBlobContainerPublicAccessE(containerName, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// GetStorageAccountKind returns one of Storage, StorageV2, BlobStorage, FileStorage, or BlockBlobStorage.
// This function would fail the test if there is an error.
func GetStorageAccountKind(t *testing.T, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	result, err := GetStorageAccountKindE(storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// GetStorageAccountSkuTier returns the storage account sku tier as Standard or Premium.
// This function would fail the test if there is an error.
func GetStorageAccountSkuTier(t *testing.T, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	result, err := GetStorageAccountSkuTierE(storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// GetStorageDNSString builds and returns the storage account dns string if the storage account exists.
// This function would fail the test if there is an error.
func GetStorageDNSString(t *testing.T, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	result, err := GetStorageDNSStringE(storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return result
}

// StorageAccountExistsE indicates whether the storage account name exists; otherwise false.
func StorageAccountExistsE(storageAccountName, resourceGroupName, subscriptionID string) (bool, error) {
	_, err := GetStorageAccountE(storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetStorageAccountE gets a storage account; otherwise error.  See https://docs.microsoft.com/rest/api/storagerp/storageaccounts/getproperties for more information.
func GetStorageAccountE(storageAccountName, resourceGroupName, subscriptionID string) (*storage.Account, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}
	resourceGroupName, err2 := getTargetAzureResourceGroupName((resourceGroupName))
	if err2 != nil {
		return nil, err2
	}
	storageAccount, err3 := GetStorageAccountPropertyE(storageAccountName, resourceGroupName, subscriptionID)
	if err3 != nil {
		return nil, err3
	}
	return storageAccount, nil
}

// StorageBlobContainerExistsE returns true if the container name exists; otherwise false.
func StorageBlobContainerExistsE(containerName, storageAccountName, resourceGroupName, subscriptionID string) (bool, error) {
	_, err := GetStorageBlobContainerE(containerName, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetStorageBlobContainerPublicAccessE indicates whether a storage container has public access; otherwise false.
func GetStorageBlobContainerPublicAccessE(containerName, storageAccountName, resourceGroupName, subscriptionID string) (bool, error) {
	container, err := GetStorageBlobContainerE(containerName, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}

	return (string(container.PublicAccess) != "None"), nil
}

// GetStorageAccountKindE returns one of Storage, StorageV2, BlobStorage, FileStorage, or BlockBlobStorage.
func GetStorageAccountKindE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {

	storageAccount, err := GetStorageAccountPropertyE(storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}
	return string(storageAccount.Kind), nil
}

// GetStorageAccountSkuTierE returns the storage account sku tier as Standard or Premium.
func GetStorageAccountSkuTierE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyE(storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}
	return string(storageAccount.Sku.Tier), nil
}

// GetStorageBlobContainerE returns Blob container client.
func GetStorageBlobContainerE(containerName, storageAccountName, resourceGroupName, subscriptionID string) (*storage.BlobContainer, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}
	resourceGroupName, err2 := getTargetAzureResourceGroupName((resourceGroupName))
	if err2 != nil {
		return nil, err2
	}
	client, err := CreateStorageBlobContainerClientE(subscriptionID)
	if err != nil {
		return nil, err
	}
	container, err := client.Get(context.Background(), resourceGroupName, storageAccountName, containerName)
	if err != nil {
		return nil, err
	}
	return &container, nil
}

// GetStorageAccountPropertyE returns StorageAccount properties.
func GetStorageAccountPropertyE(storageAccountName, resourceGroupName, subscriptionID string) (*storage.Account, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}
	resourceGroupName, err2 := getTargetAzureResourceGroupName((resourceGroupName))
	if err2 != nil {
		return nil, err2
	}
	client, err := CreateStorageAccountClientE(subscriptionID)
	if err != nil {
		return nil, err
	}
	account, err := client.GetProperties(context.Background(), resourceGroupName, storageAccountName, "")
	if err != nil {
		return nil, err
	}
	return &account, nil
}

// GetStorageAccountClientE creates a storage account client.
// TODO: remove in next version
func GetStorageAccountClientE(subscriptionID string) (*storage.AccountsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	storageAccountClient := storage.NewAccountsClient(subscriptionID)
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	storageAccountClient.Authorizer = *authorizer
	return &storageAccountClient, nil
}

// GetStorageBlobContainerClientE creates a storage container client.
// TODO: remove in next version
func GetStorageBlobContainerClientE(subscriptionID string) (*storage.BlobContainersClient, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	blobContainerClient := storage.NewBlobContainersClient(subscriptionID)
	authorizer, err := NewAuthorizer()

	if err != nil {
		return nil, err
	}
	blobContainerClient.Authorizer = *authorizer
	return &blobContainerClient, nil
}

// GetStorageURISuffixE returns the proper storage URI suffix for the configured Azure environment.
func GetStorageURISuffixE() (string, error) {
	envName := "AzurePublicCloud"
	env, err := azure.EnvironmentFromName(envName)
	if err != nil {
		return "", err
	}
	return env.StorageEndpointSuffix, nil
}

// GetStorageAccountPrimaryBlobEndpointE gets the storage account blob endpoint as URI string.
func GetStorageAccountPrimaryBlobEndpointE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyE(storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return *storageAccount.AccountProperties.PrimaryEndpoints.Blob, nil
}

// GetStorageDNSStringE builds and returns the storage account dns string if the storage account exists.
func GetStorageDNSStringE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	retval, err := StorageAccountExistsE(storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}
	if retval {
		storageSuffix, err2 := GetStorageURISuffixE()
		if err2 != nil {
			return "", err2
		}
		return fmt.Sprintf("https://%s.blob.%s/", storageAccountName, storageSuffix), nil
	}

	return "", NewNotFoundError("storage account", storageAccountName, "")
}
