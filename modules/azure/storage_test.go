package azure_test

import (
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

	_, err := azure.StorageAccountExistsE("", "", "")
	require.Error(t, err)
}

func TestStorageBlobContainerExists(t *testing.T) {
	t.Parallel()

	_, err := azure.StorageBlobContainerExistsE("", "", "", "")
	require.Error(t, err)
}

func TestStorageBlobContainerPublicAccess(t *testing.T) {
	t.Parallel()

	_, err := azure.GetStorageBlobContainerPublicAccessE("", "", "", "")
	require.Error(t, err)
}

func TestGetStorageAccountKind(t *testing.T) {
	t.Parallel()

	_, err := azure.GetStorageAccountKindE("", "", "")
	require.Error(t, err)
}

func TestGetStorageAccountSkuTier(t *testing.T) {
	t.Parallel()

	_, err := azure.GetStorageAccountSkuTierE("", "", "")
	require.Error(t, err)
}

func TestGetStorageDNSString(t *testing.T) {
	t.Parallel()

	_, err := azure.GetStorageDNSStringE("", "", "")
	require.Error(t, err)
}
