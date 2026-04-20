package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// StorageAccountExists indicates whether the storage account name exactly matches; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [StorageAccountExistsContext] instead.
func StorageAccountExists(t testing.TestingT, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return StorageAccountExistsContext(t, context.Background(), storageAccountName, resourceGroupName, subscriptionID)
}

// StorageAccountExistsE indicates whether the storage account name exactly matches; otherwise false.
//
// Deprecated: Use [StorageAccountExistsContextE] instead.
func StorageAccountExistsE(storageAccountName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return StorageAccountExistsContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
}

// StorageAccountExistsContext indicates whether the storage account name exactly matches; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func StorageAccountExistsContext(t testing.TestingT, ctx context.Context, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	result, err := StorageAccountExistsContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
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

// StorageBlobContainerExistsE returns true if the container name exactly matches; otherwise false.
//
// Deprecated: Use [StorageBlobContainerExistsContextE] instead.
func StorageBlobContainerExistsE(containerName string, storageAccountName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return StorageBlobContainerExistsContextE(context.Background(), containerName, storageAccountName, resourceGroupName, subscriptionID)
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

// StorageFileShareExists returns true if the file share name exactly matches; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [StorageFileShareExistsContext] instead.
func StorageFileShareExists(t testing.TestingT, fileShareName string, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return StorageFileShareExistsContext(t, context.Background(), fileShareName, storageAccountName, resourceGroupName, subscriptionID)
}

// StorageFileShareExistsE returns true if the file share name exactly matches; otherwise false.
//
// Deprecated: Use [StorageFileShareExistsContextE] instead.
func StorageFileShareExistsE(fileShareName string, storageAccountName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return StorageFileShareExistsContextE(context.Background(), fileShareName, storageAccountName, resourceGroupName, subscriptionID)
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

// GetStorageBlobContainerPublicAccess indicates whether a storage container has public access; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetStorageBlobContainerPublicAccessContext] instead.
func GetStorageBlobContainerPublicAccess(t testing.TestingT, containerName string, storageAccountName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return GetStorageBlobContainerPublicAccessContext(t, context.Background(), containerName, storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageBlobContainerPublicAccessE indicates whether a storage container has public access; otherwise false.
//
// Deprecated: Use [GetStorageBlobContainerPublicAccessContextE] instead.
func GetStorageBlobContainerPublicAccessE(containerName string, storageAccountName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return GetStorageBlobContainerPublicAccessContextE(context.Background(), containerName, storageAccountName, resourceGroupName, subscriptionID)
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

// GetStorageAccountKind returns one of Storage, StorageV2, BlobStorage, FileStorage, or BlockBlobStorage.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetStorageAccountKindContext] instead.
func GetStorageAccountKind(t testing.TestingT, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	t.Helper()

	return GetStorageAccountKindContext(t, context.Background(), storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageAccountKindE returns one of Storage, StorageV2, BlobStorage, FileStorage, or BlockBlobStorage.
//
// Deprecated: Use [GetStorageAccountKindContextE] instead.
func GetStorageAccountKindE(storageAccountName string, resourceGroupName string, subscriptionID string) (string, error) {
	return GetStorageAccountKindContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
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

// GetStorageAccountSkuTier returns the storage account sku tier as Standard or Premium.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetStorageAccountSkuTierContext] instead.
func GetStorageAccountSkuTier(t testing.TestingT, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	t.Helper()

	return GetStorageAccountSkuTierContext(t, context.Background(), storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageAccountSkuTierE returns the storage account sku tier as Standard or Premium.
//
// Deprecated: Use [GetStorageAccountSkuTierContextE] instead.
func GetStorageAccountSkuTierE(storageAccountName string, resourceGroupName string, subscriptionID string) (string, error) {
	return GetStorageAccountSkuTierContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
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

// GetStorageDNSString builds and returns the storage account dns string if the storage account exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetStorageDNSStringContext] instead.
func GetStorageDNSString(t testing.TestingT, storageAccountName string, resourceGroupName string, subscriptionID string) string {
	t.Helper()

	return GetStorageDNSStringContext(t, context.Background(), storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageDNSStringE builds and returns the storage account dns string if the storage account exists.
//
// Deprecated: Use [GetStorageDNSStringContextE] instead.
func GetStorageDNSStringE(storageAccountName string, resourceGroupName string, subscriptionID string) (string, error) {
	return GetStorageDNSStringContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
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

// GetStorageAccountE gets a storage account; otherwise error.
//
// Deprecated: Use [GetStorageAccountContextE] instead.
func GetStorageAccountE(storageAccountName, resourceGroupName, subscriptionID string) (*armstorage.Account, error) {
	return GetStorageAccountContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
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

	return ExtractBlobContainerPublicAccess(container), nil
}

// GetStorageAccountKindContextE returns one of Storage, StorageV2, BlobStorage, FileStorage, or BlockBlobStorage.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountKindContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return ExtractStorageAccountKind(storageAccount), nil
}

// GetStorageAccountSkuTierContextE returns the storage account sku tier as Standard or Premium.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountSkuTierContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return ExtractStorageAccountSkuTier(storageAccount), nil
}

// GetStorageBlobContainerE returns the Blob container client.
//
// Deprecated: Use [GetStorageBlobContainerContextE] instead.
func GetStorageBlobContainerE(containerName, storageAccountName, resourceGroupName, subscriptionID string) (*armstorage.BlobContainer, error) {
	return GetStorageBlobContainerContextE(context.Background(), containerName, storageAccountName, resourceGroupName, subscriptionID)
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

	client, err := CreateStorageBlobContainerClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return FetchBlobContainer(ctx, client, resourceGroupName, storageAccountName, containerName)
}

// GetStorageAccountPropertyE returns StorageAccount properties.
//
// Deprecated: Use [GetStorageAccountPropertyContextE] instead.
func GetStorageAccountPropertyE(storageAccountName, resourceGroupName, subscriptionID string) (*armstorage.Account, error) {
	return GetStorageAccountPropertyContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
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

	client, err := CreateStorageAccountClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return FetchStorageAccountProperties(ctx, client, resourceGroupName, storageAccountName)
}

// GetStorageFileShare returns the specified file share.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetStorageFileShareContext] instead.
func GetStorageFileShare(t testing.TestingT, fileShareName, storageAccountName, resourceGroupName, subscriptionID string) *armstorage.FileShare {
	t.Helper()

	return GetStorageFileShareContext(t, context.Background(), fileShareName, storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageFileShareE returns the specified file share.
//
// Deprecated: Use [GetStorageFileShareContextE] instead.
func GetStorageFileShareE(fileShareName, storageAccountName, resourceGroupName, subscriptionID string) (*armstorage.FileShare, error) {
	return GetStorageFileShareContextE(context.Background(), fileShareName, storageAccountName, resourceGroupName, subscriptionID)
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

	client, err := CreateStorageFileSharesClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return FetchFileShare(ctx, client, resourceGroupName, storageAccountName, fileShareName)
}

// FetchStorageAccountProperties retrieves the storage account properties using the provided client.
func FetchStorageAccountProperties(ctx context.Context, client *armstorage.AccountsClient, resourceGroupName, storageAccountName string) (*armstorage.Account, error) {
	resp, err := client.GetProperties(ctx, resourceGroupName, storageAccountName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Account, nil
}

// FetchBlobContainer retrieves a blob container using the provided client.
func FetchBlobContainer(ctx context.Context, client *armstorage.BlobContainersClient, resourceGroupName, storageAccountName, containerName string) (*armstorage.BlobContainer, error) {
	resp, err := client.Get(ctx, resourceGroupName, storageAccountName, containerName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.BlobContainer, nil
}

// FetchFileShare retrieves a file share using the provided client with stats expansion.
func FetchFileShare(ctx context.Context, client *armstorage.FileSharesClient, resourceGroupName, storageAccountName, fileShareName string) (*armstorage.FileShare, error) {
	expand := "stats"

	resp, err := client.Get(ctx, resourceGroupName, storageAccountName, fileShareName, &armstorage.FileSharesClientGetOptions{
		Expand: &expand,
	})
	if err != nil {
		return nil, err
	}

	return &resp.FileShare, nil
}

// ExtractBlobContainerPublicAccess returns true if the container has public access other than "None".
func ExtractBlobContainerPublicAccess(container *armstorage.BlobContainer) bool {
	if container == nil || container.ContainerProperties == nil || container.ContainerProperties.PublicAccess == nil {
		return false
	}

	return *container.ContainerProperties.PublicAccess != armstorage.PublicAccessNone
}

// ExtractStorageAccountKind returns the storage account kind as a string.
func ExtractStorageAccountKind(account *armstorage.Account) string {
	if account == nil || account.Kind == nil {
		return ""
	}

	return string(*account.Kind)
}

// ExtractStorageAccountSkuTier returns the storage account SKU tier as a string.
func ExtractStorageAccountSkuTier(account *armstorage.Account) string {
	if account == nil || account.SKU == nil || account.SKU.Tier == nil {
		return ""
	}

	return string(*account.SKU.Tier)
}

// GetStorageAccountPrimaryBlobEndpointE gets the storage account blob endpoint as URI string.
//
// Deprecated: Use [GetStorageAccountPrimaryBlobEndpointContextE] instead.
func GetStorageAccountPrimaryBlobEndpointE(storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	return GetStorageAccountPrimaryBlobEndpointContextE(context.Background(), storageAccountName, resourceGroupName, subscriptionID)
}

// GetStorageAccountPrimaryBlobEndpointContextE gets the storage account blob endpoint as URI string.
// The ctx parameter supports cancellation and timeouts.
func GetStorageAccountPrimaryBlobEndpointContextE(ctx context.Context, storageAccountName, resourceGroupName, subscriptionID string) (string, error) {
	storageAccount, err := GetStorageAccountPropertyContextE(ctx, storageAccountName, resourceGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	if storageAccount == nil || storageAccount.Properties == nil ||
		storageAccount.Properties.PrimaryEndpoints == nil ||
		storageAccount.Properties.PrimaryEndpoints.Blob == nil {
		return "", NewNotFoundError("primary blob endpoint", storageAccountName, "")
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
		storageSuffix, err2 := GetStorageURISuffixE() //nolint:contextcheck
		if err2 != nil {
			return "", err2
		}

		return fmt.Sprintf("https://%s.blob.%s/", storageAccountName, storageSuffix), nil
	}

	return "", NewNotFoundError("storage account", storageAccountName, "")
}
