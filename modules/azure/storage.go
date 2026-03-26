package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-06-01/storage"
	"github.com/Azure/go-autorest/autorest/azure"
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

// StorageAccountExists indicates whether the storage account name exactly matches; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [StorageAccountExistsContext] instead.
func StorageAccountExists(t testing.TestingT, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return StorageAccountExistsContext(t, context.Background(), storageAccountName, resourceGroupName, subscriptionID)
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

// StorageBlobContainerExists returns true if the container name exactly matches; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [StorageBlobContainerExistsContext] instead.
func StorageBlobContainerExists(t testing.TestingT, containerName string, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return StorageBlobContainerExistsContext(t, context.Background(), containerName, storageAccountName, resourceGroupName, subscriptionID)
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

// StorageFileShareExists returns true if the file share name exactly matches; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [StorageFileShareExistsContext] instead.
func StorageFileShareExists(t testing.TestingT, fileShareName string, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return StorageFileShareExistsContext(t, context.Background(), fileShareName, storageAccountName, resourceGroupName, subscriptionID)
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

// StorageFileShareExistsE returns true if the file share name exactly matches; otherwise false.
//
// Deprecated: Use [StorageFileShareExistsContextE] instead.
func StorageFileShareExistsE(fileShareName string, storageAccountName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return StorageFileShareExistsContextE(context.Background(), fileShareName, storageAccountName, resourceGroupName, subscriptionID)
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

// GetStorageBlobContainerPublicAccess indicates whether a storage container has public access; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetStorageBlobContainerPublicAccessContext] instead.
func GetStorageBlobContainerPublicAccess(t testing.TestingT, containerName string, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return GetStorageBlobContainerPublicAccessContext(t, context.Background(), containerName, storageAccountName, resourceGroupName, subscriptionID)
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

// GetStorageAccountKind returns one of Storage, StorageV2, BlobStorage, FileStorage, or BlockBlobStorage.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetStorageAccountKindContext] instead.
func GetStorageAccountKind(t testing.TestingT, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	t.Helper()

	return GetStorageAccountKindContext(t, context.Background(), storageAccountName, resourceGroupName, subscriptionID)
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

// GetStorageAccountSkuTier returns the storage account sku tier as Standard or Premium.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetStorageAccountSkuTierContext] instead.
func GetStorageAccountSkuTier(t testing.TestingT, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	t.Helper()

	return GetStorageAccountSkuTierContext(t, context.Background(), storageAccountName, resourceGroupName, subscriptionID)
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

// GetStorageDNSString builds and returns the storage account dns string if the storage account exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetStorageDNSStringContext] instead.
func GetStorageDNSString(t testing.TestingT, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	t.Helper()

	return GetStorageDNSStringContext(t, context.Background(), storageAccountName, resourceGroupName, subscriptionID)
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

// StorageAccountExistsE indicates whether the storage account name exists; otherwise false.
//
// Deprecated: Use [StorageAccountExistsContextE] instead.
func StorageAccountExistsE(storageAccountName, resourceGroupName, subscriptionID string) (bool, error) {
	return StorageAccountExistsContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageAccountContextE gets a storage account; otherwise error.
// See https://docs.microsoft.com/rest/api/storagerp/storageaccounts/getproperties for more information.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (*storage.Account, error) {
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

// GetStorageAccountE gets a storage account; otherwise error.
// See https://docs.microsoft.com/rest/api/storagerp/storageaccounts/getproperties for more information.
//
// Deprecated: Use [GetStorageAccountContextE] instead.
func GetStorageAccountE(storageAccountName, resourceGroupName, subscriptionID string) (*storage.Account, error) {
	return GetStorageAccountContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
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

// StorageBlobContainerExistsE returns true if the container name exists; otherwise false.
//
// Deprecated: Use [StorageBlobContainerExistsContextE] instead.
func StorageBlobContainerExistsE(containerName, storageAccountName, resourceGroupName, subscriptionID string) (bool, error) {
	return StorageBlobContainerExistsContextE(context.Background(), containerName, storageAccountName, resourceGroupName, subscriptionID)
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

	return (string(container.PublicAccess) != "None"), nil
}

// GetStorageBlobContainerPublicAccessE indicates whether a storage container has public access; otherwise false.
//
// Deprecated: Use [GetStorageBlobContainerPublicAccessContextE] instead.
func GetStorageBlobContainerPublicAccessE(containerName, storageAccountName, resourceGroupName, subscriptionID string) (bool, error) {
	return GetStorageBlobContainerPublicAccessContextE(context.Background(), containerName, storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageAccountKindContextE returns one of Storage, StorageV2, BlobStorage, FileStorage, or BlockBlobStorage.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountKindContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return string(storageAccount.Kind), nil
}

// GetStorageAccountKindE returns one of Storage, StorageV2, BlobStorage, FileStorage, or BlockBlobStorage.
//
// Deprecated: Use [GetStorageAccountKindContextE] instead.
func GetStorageAccountKindE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	return GetStorageAccountKindContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageAccountSkuTierContextE returns the storage account sku tier as Standard or Premium.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountSkuTierContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return string(storageAccount.Sku.Tier), nil
}

// GetStorageAccountSkuTierE returns the storage account sku tier as Standard or Premium.
//
// Deprecated: Use [GetStorageAccountSkuTierContextE] instead.
func GetStorageAccountSkuTierE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	return GetStorageAccountSkuTierContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageBlobContainerContextE returns the Blob container client.
// The ctx parameter supports cancellation and timeouts.
func GetStorageBlobContainerContextE(ctx context.Context, containerName, storageAccountName, resourceGroupName, subscriptionID string) (*storage.BlobContainer, error) {
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

	container, err := client.Get(ctx, resourceGroupName, storageAccountName, containerName)
	if err != nil {
		return nil, err
	}

	return &container, nil
}

// GetStorageBlobContainerE returns the Blob container client.
//
// Deprecated: Use [GetStorageBlobContainerContextE] instead.
func GetStorageBlobContainerE(containerName, storageAccountName, resourceGroupName, subscriptionID string) (*storage.BlobContainer, error) {
	return GetStorageBlobContainerContextE(context.Background(), containerName, storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageAccountPropertyContextE returns StorageAccount properties.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountPropertyContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (*storage.Account, error) {
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

	account, err := client.GetProperties(ctx, resourceGroupName, storageAccountName, "")
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// GetStorageAccountPropertyE returns StorageAccount properties.
//
// Deprecated: Use [GetStorageAccountPropertyContextE] instead.
func GetStorageAccountPropertyE(storageAccountName, resourceGroupName, subscriptionID string) (*storage.Account, error) {
	return GetStorageAccountPropertyContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageFileShareContext returns the specified file share.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageFileShareContext(t testing.TestingT, ctx context.Context, fileShareName, storageAccountName, resourceGroupName, subscriptionID string) *storage.FileShare {
	t.Helper()

	fileShare, err := GetStorageFileShareContextE(ctx, fileShareName, storageAccountName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return fileShare
}

// GetStorageFileShare returns the specified file share.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetStorageFileShareContext] instead.
func GetStorageFileShare(t testing.TestingT, fileShareName, storageAccountName, resourceGroupName, subscriptionID string) *storage.FileShare {
	t.Helper()

	return GetStorageFileShareContext(t, context.Background(), fileShareName, storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageFileShareContextE returns the specified file share.
// The ctx parameter supports cancellation and timeouts.
func GetStorageFileShareContextE(ctx context.Context, fileShareName, storageAccountName, resourceGroupName, subscriptionID string) (*storage.FileShare, error) {
	resourceGroupName, err2 := getTargetAzureResourceGroupName(resourceGroupName)
	if err2 != nil {
		return nil, err2
	}

	client, err := CreateStorageFileSharesClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	fileShare, err := client.Get(ctx, resourceGroupName, storageAccountName, fileShareName, "stats")
	if err != nil {
		return nil, err
	}

	return &fileShare, nil
}

// GetStorageFileShareE returns the specified file share.
//
// Deprecated: Use [GetStorageFileShareContextE] instead.
func GetStorageFileShareE(fileShareName, storageAccountName, resourceGroupName, subscriptionID string) (*storage.FileShare, error) {
	return GetStorageFileShareContextE(context.Background(), fileShareName, storageAccountName, resourceGroupName, subscriptionID)
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

// GetStorageAccountPrimaryBlobEndpointContextE gets the storage account blob endpoint as URI string.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountPrimaryBlobEndpointContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return *storageAccount.AccountProperties.PrimaryEndpoints.Blob, nil
}

// GetStorageAccountPrimaryBlobEndpointE gets the storage account blob endpoint as URI string.
//
// Deprecated: Use [GetStorageAccountPrimaryBlobEndpointContextE] instead.
func GetStorageAccountPrimaryBlobEndpointE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	return GetStorageAccountPrimaryBlobEndpointContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
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

// GetStorageDNSStringE builds and returns the storage account dns string if the storage account exists.
//
// Deprecated: Use [GetStorageDNSStringContextE] instead.
func GetStorageDNSStringE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	return GetStorageDNSStringContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
}
