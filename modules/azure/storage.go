package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// StorageAccountExistsContext indicates whether the storage account name exactly matches; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func StorageAccountExistsContext(t testing.TestingT, ctx context.Context, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	result, err := StorageAccountExistsContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return result
}

// StorageBlobContainerExistsContext returns true if the container name exactly matches; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func StorageBlobContainerExistsContext(t testing.TestingT, ctx context.Context, containerName string, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	result, err := StorageBlobContainerExistsContextE(ctx, containerName, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return result
}

// StorageFileShareExistsContext returns true if the file share name exactly matches; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func StorageFileShareExistsContext(t testing.TestingT, ctx context.Context, fileShareName string, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	result, err := StorageFileShareExistsContextE(ctx, fileShareName, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return result
}

// StorageFileShareExistsContextE returns true if the file share name exactly matches; otherwise false.
// The ctx parameter supports cancellation and timeouts.
func StorageFileShareExistsContextE(ctx context.Context, fileShareName string, storageAccountName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetStorageFileShareContextE(ctx, fileShareName, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetStorageBlobContainerPublicAccessContext indicates whether a storage container has public access; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageBlobContainerPublicAccessContext(t testing.TestingT, ctx context.Context, containerName string, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	result, err := GetStorageBlobContainerPublicAccessContextE(ctx, containerName, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return result
}

// GetStorageAccountKindContext returns one of Storage, StorageV2, BlobStorage, FileStorage, or BlockBlobStorage.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountKindContext(t testing.TestingT, ctx context.Context, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	t.Helper()

	result, err := GetStorageAccountKindContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return result
}

// GetStorageAccountSkuTierContext returns the storage account sku tier as Standard or Premium.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountSkuTierContext(t testing.TestingT, ctx context.Context, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	t.Helper()

	result, err := GetStorageAccountSkuTierContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return result
}

// GetStorageDNSStringContext builds and returns the storage account dns string if the storage account exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageDNSStringContext(t testing.TestingT, ctx context.Context, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	t.Helper()

	result, err := GetStorageDNSStringContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return result
}

// StorageAccountExistsContextE indicates whether the storage account name exists; otherwise false.
// The ctx parameter supports cancellation and timeouts.
func StorageAccountExistsContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (bool, error) {
	_, err := GetStorageAccountContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetStorageAccountContextE gets a storage account; otherwise error.
// See https://docs.microsoft.com/rest/api/storagerp/storageaccounts/getproperties for more information.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (*armstorage.Account, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	resourceGroupName, err2 := getTargetAzureResourceGroupName(resourceGroupName)
	if err2 != nil {
		return nil, err2
	}

	storageAccount, err3 := GetStorageAccountPropertyContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	if err3 != nil {
		return nil, err3
	}

	return storageAccount, nil
}

// StorageBlobContainerExistsContextE returns true if the container name exists; otherwise false.
// The ctx parameter supports cancellation and timeouts.
func StorageBlobContainerExistsContextE(ctx context.Context, containerName, storageAccountName, resourceGroupName, subscriptionID string) (bool, error) {
	_, err := GetStorageBlobContainerContextE(ctx, containerName, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetStorageBlobContainerPublicAccessContextE indicates whether a storage container has public access; otherwise false.
// The ctx parameter supports cancellation and timeouts.
func GetStorageBlobContainerPublicAccessContextE(ctx context.Context, containerName, storageAccountName, resourceGroupName, subscriptionID string) (bool, error) {
	container, err := GetStorageBlobContainerContextE(ctx, containerName, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return extractBlobContainerPublicAccess(container), nil
}

// GetStorageAccountKindContextE returns one of Storage, StorageV2, BlobStorage, FileStorage, or BlockBlobStorage.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountKindContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return extractStorageAccountKind(storageAccount), nil
}

// GetStorageAccountSkuTierContextE returns the storage account sku tier as Standard or Premium.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountSkuTierContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return extractStorageAccountSkuTier(storageAccount), nil
}

// GetStorageBlobContainerContextE returns the Blob container client.
// The ctx parameter supports cancellation and timeouts.
func GetStorageBlobContainerContextE(ctx context.Context, containerName, storageAccountName, resourceGroupName, subscriptionID string) (*armstorage.BlobContainer, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	resourceGroupName, err2 := getTargetAzureResourceGroupName(resourceGroupName)
	if err2 != nil {
		return nil, err2
	}

	client, err := CreateStorageBlobContainerClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	return fetchBlobContainer(ctx, client, resourceGroupName, storageAccountName, containerName)
}

// GetStorageAccountPropertyContextE returns StorageAccount properties.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountPropertyContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (*armstorage.Account, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	resourceGroupName, err2 := getTargetAzureResourceGroupName(resourceGroupName)
	if err2 != nil {
		return nil, err2
	}

	client, err := CreateStorageAccountClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	return fetchStorageAccountProperties(ctx, client, resourceGroupName, storageAccountName)
}

// GetStorageFileShareContext returns the specified file share.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageFileShareContext(t testing.TestingT, ctx context.Context, fileShareName, storageAccountName, resourceGroupName, subscriptionID string) *armstorage.FileShare {
	t.Helper()

	fileShare, err := GetStorageFileShareContextE(ctx, fileShareName, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return fileShare
}

// GetStorageFileShareContextE returns the specified file share.
// The ctx parameter supports cancellation and timeouts.
func GetStorageFileShareContextE(ctx context.Context, fileShareName, storageAccountName, resourceGroupName, subscriptionID string) (*armstorage.FileShare, error) {
	resourceGroupName, err2 := getTargetAzureResourceGroupName(resourceGroupName)
	if err2 != nil {
		return nil, err2
	}

	client, err := CreateStorageFileSharesClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	return fetchFileShare(ctx, client, resourceGroupName, storageAccountName, fileShareName)
}

// fetchStorageAccountProperties retrieves the storage account properties using the provided client.
func fetchStorageAccountProperties(ctx context.Context, client *armstorage.AccountsClient, resourceGroupName, storageAccountName string) (*armstorage.Account, error) {
	resp, err := client.GetProperties(ctx, resourceGroupName, storageAccountName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Account, nil
}

// fetchBlobContainer retrieves a blob container using the provided client.
func fetchBlobContainer(ctx context.Context, client *armstorage.BlobContainersClient, resourceGroupName, storageAccountName, containerName string) (*armstorage.BlobContainer, error) {
	resp, err := client.Get(ctx, resourceGroupName, storageAccountName, containerName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.BlobContainer, nil
}

// fetchFileShare retrieves a file share using the provided client with stats expansion.
func fetchFileShare(ctx context.Context, client *armstorage.FileSharesClient, resourceGroupName, storageAccountName, fileShareName string) (*armstorage.FileShare, error) {
	expand := "stats"

	resp, err := client.Get(ctx, resourceGroupName, storageAccountName, fileShareName, &armstorage.FileSharesClientGetOptions{
		Expand: &expand,
	})
	if err != nil {
		return nil, err
	}

	return &resp.FileShare, nil
}

// extractBlobContainerPublicAccess returns true if the container has public access other than "None".
func extractBlobContainerPublicAccess(container *armstorage.BlobContainer) bool {
	return string(*container.ContainerProperties.PublicAccess) != "None"
}

// extractStorageAccountKind returns the storage account kind as a string.
func extractStorageAccountKind(account *armstorage.Account) string {
	return string(*account.Kind)
}

// extractStorageAccountSkuTier returns the storage account SKU tier as a string.
func extractStorageAccountSkuTier(account *armstorage.Account) string {
	return string(*account.SKU.Tier)
}

// GetStorageAccountPrimaryBlobEndpointContextE gets the storage account blob endpoint as URI string.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountPrimaryBlobEndpointContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return *storageAccount.Properties.PrimaryEndpoints.Blob, nil
}

// GetStorageDNSStringContextE builds and returns the storage account dns string if the storage account exists.
// The ctx parameter supports cancellation and timeouts.
func GetStorageDNSStringContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	retval, err := StorageAccountExistsContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
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
