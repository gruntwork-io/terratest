package azure

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	storagefake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage/fake"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Fake client helpers
// ---------------------------------------------------------------------------

func newFakeStorageAccountsClient(t *testing.T, srv storagefake.AccountsServer) *armstorage.AccountsClient {
	t.Helper()
	client, err := armstorage.NewAccountsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: storagefake.NewAccountsServerTransport(&srv),
		}})
	require.NoError(t, err)
	return client
}

func newFakeBlobContainersClient(t *testing.T, srv storagefake.BlobContainersServer) *armstorage.BlobContainersClient {
	t.Helper()
	client, err := armstorage.NewBlobContainersClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: storagefake.NewBlobContainersServerTransport(&srv),
		}})
	require.NoError(t, err)
	return client
}

func newFakeFileSharesClient(t *testing.T, srv storagefake.FileSharesServer) *armstorage.FileSharesClient {
	t.Helper()
	client, err := armstorage.NewFileSharesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: storagefake.NewFileSharesServerTransport(&srv),
		}})
	require.NoError(t, err)
	return client
}

// ---------------------------------------------------------------------------
// fetchStorageAccountProperties tests
// ---------------------------------------------------------------------------

func TestFetchStorageAccountProperties_Success(t *testing.T) {
	t.Parallel()

	srv := storagefake.AccountsServer{
		GetProperties: func(ctx context.Context, resourceGroupName, accountName string, options *armstorage.AccountsClientGetPropertiesOptions) (resp azfake.Responder[armstorage.AccountsClientGetPropertiesResponse], errResp azfake.ErrorResponder) {
			result := armstorage.AccountsClientGetPropertiesResponse{
				Account: armstorage.Account{
					Name: to.Ptr("teststorage"),
					Kind: to.Ptr(armstorage.KindStorageV2),
					SKU:  &armstorage.SKU{Tier: to.Ptr(armstorage.SKUTierStandard)},
				},
			}
			resp.SetResponse(http.StatusOK, result, nil)
			return
		},
	}

	client := newFakeStorageAccountsClient(t, srv)
	account, err := fetchStorageAccountProperties(context.Background(), client, "rg", "teststorage")
	require.NoError(t, err)
	assert.Equal(t, "teststorage", *account.Name)
}

func TestFetchStorageAccountProperties_NotFound(t *testing.T) {
	t.Parallel()

	srv := storagefake.AccountsServer{
		GetProperties: func(ctx context.Context, resourceGroupName, accountName string, options *armstorage.AccountsClientGetPropertiesOptions) (resp azfake.Responder[armstorage.AccountsClientGetPropertiesResponse], errResp azfake.ErrorResponder) {
			errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")
			return
		},
	}

	client := newFakeStorageAccountsClient(t, srv)
	_, err := fetchStorageAccountProperties(context.Background(), client, "rg", "missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// ---------------------------------------------------------------------------
// fetchBlobContainer tests
// ---------------------------------------------------------------------------

func TestFetchBlobContainer_Success(t *testing.T) {
	t.Parallel()

	srv := storagefake.BlobContainersServer{
		Get: func(ctx context.Context, resourceGroupName, accountName, containerName string, options *armstorage.BlobContainersClientGetOptions) (resp azfake.Responder[armstorage.BlobContainersClientGetResponse], errResp azfake.ErrorResponder) {
			result := armstorage.BlobContainersClientGetResponse{
				BlobContainer: armstorage.BlobContainer{
					Name: to.Ptr("testcontainer"),
				},
			}
			resp.SetResponse(http.StatusOK, result, nil)
			return
		},
	}

	client := newFakeBlobContainersClient(t, srv)
	container, err := fetchBlobContainer(context.Background(), client, "rg", "teststorage", "testcontainer")
	require.NoError(t, err)
	assert.Equal(t, "testcontainer", *container.Name)
}

// ---------------------------------------------------------------------------
// fetchFileShare tests
// ---------------------------------------------------------------------------

func TestFetchFileShare_Success(t *testing.T) {
	t.Parallel()

	srv := storagefake.FileSharesServer{
		Get: func(ctx context.Context, resourceGroupName, accountName, shareName string, options *armstorage.FileSharesClientGetOptions) (resp azfake.Responder[armstorage.FileSharesClientGetResponse], errResp azfake.ErrorResponder) {
			result := armstorage.FileSharesClientGetResponse{
				FileShare: armstorage.FileShare{
					Name: to.Ptr("testshare"),
				},
			}
			resp.SetResponse(http.StatusOK, result, nil)
			return
		},
	}

	client := newFakeFileSharesClient(t, srv)
	share, err := fetchFileShare(context.Background(), client, "rg", "teststorage", "testshare")
	require.NoError(t, err)
	assert.Equal(t, "testshare", *share.Name)
}

// ---------------------------------------------------------------------------
// extractBlobContainerPublicAccess tests
// ---------------------------------------------------------------------------

func TestExtractBlobContainerPublicAccess_None(t *testing.T) {
	t.Parallel()

	container := &armstorage.BlobContainer{
		ContainerProperties: &armstorage.ContainerProperties{
			PublicAccess: to.Ptr(armstorage.PublicAccessNone),
		},
	}
	assert.False(t, extractBlobContainerPublicAccess(container))
}

func TestExtractBlobContainerPublicAccess_Blob(t *testing.T) {
	t.Parallel()

	container := &armstorage.BlobContainer{
		ContainerProperties: &armstorage.ContainerProperties{
			PublicAccess: to.Ptr(armstorage.PublicAccessBlob),
		},
	}
	assert.True(t, extractBlobContainerPublicAccess(container))
}

// ---------------------------------------------------------------------------
// extractStorageAccountKind tests
// ---------------------------------------------------------------------------

func TestExtractStorageAccountKind(t *testing.T) {
	t.Parallel()

	account := &armstorage.Account{
		Kind: to.Ptr(armstorage.KindStorageV2),
	}
	assert.Equal(t, "StorageV2", extractStorageAccountKind(account))
}

// ---------------------------------------------------------------------------
// extractStorageAccountSkuTier tests
// ---------------------------------------------------------------------------

func TestExtractStorageAccountSkuTier(t *testing.T) {
	t.Parallel()

	account := &armstorage.Account{
		SKU: &armstorage.SKU{
			Tier: to.Ptr(armstorage.SKUTierStandard),
		},
	}
	assert.Equal(t, "Standard", extractStorageAccountSkuTier(account))
}
