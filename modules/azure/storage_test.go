package azure_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	azure "github.com/gruntwork-io/terratest/modules/azure"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods to create and delete storage accounts are added, these tests can be extended.
*/

func TestStorageAccountExists(t *testing.T) {
	t.Parallel()

	_, err := azure.StorageAccountExistsContextE(context.Background(), "", "", "")
	require.Error(t, err)
}

func TestStorageBlobContainerExists(t *testing.T) {
	t.Parallel()

	_, err := azure.StorageBlobContainerExistsContextE(context.Background(), "", "", "", "")
	require.Error(t, err)
}

func TestStorageBlobContainerPublicAccess(t *testing.T) {
	t.Parallel()

	_, err := azure.GetStorageBlobContainerPublicAccessContextE(context.Background(), "", "", "", "")
	require.Error(t, err)
}

func TestGetStorageAccountKind(t *testing.T) {
	t.Parallel()

	_, err := azure.GetStorageAccountKindContextE(context.Background(), "", "", "")
	require.Error(t, err)
}

func TestGetStorageAccountSkuTier(t *testing.T) {
	t.Parallel()

	_, err := azure.GetStorageAccountSkuTierContextE(context.Background(), "", "", "")
	require.Error(t, err)
}

func TestGetStorageDNSString(t *testing.T) {
	t.Parallel()

	_, err := azure.GetStorageDNSStringContextE(context.Background(), "", "", "")
	require.Error(t, err)
}
