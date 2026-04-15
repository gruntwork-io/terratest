package azure_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage"
	storagefake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// FetchStorageAccountProperties tests
// ---------------------------------------------------------------------------

func TestFetchStorageAccountProperties_Success(t *testing.T) {
	t.Parallel()

	srv := &storagefake.AccountsServer{
		GetProperties: func(_ context.Context, _, _ string, _ *armstorage.AccountsClientGetPropertiesOptions) (resp azfake.Responder[armstorage.AccountsClientGetPropertiesResponse], errResp azfake.ErrorResponder) {
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
	account, err := azure.FetchStorageAccountProperties(t.Context(), client, "rg", "teststorage")

	require.NoError(t, err)
	assert.Equal(t, "teststorage", *account.Name)
}

func TestFetchStorageAccountProperties_NotFound(t *testing.T) {
	t.Parallel()

	srv := &storagefake.AccountsServer{
		GetProperties: func(_ context.Context, _, _ string, _ *armstorage.AccountsClientGetPropertiesOptions) (resp azfake.Responder[armstorage.AccountsClientGetPropertiesResponse], errResp azfake.ErrorResponder) {
			errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

			return
		},
	}

	client := newFakeStorageAccountsClient(t, srv)
	_, err := azure.FetchStorageAccountProperties(t.Context(), client, "rg", "missing")

	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	assert.Equal(t, "ResourceNotFound", respErr.ErrorCode)
}

// ---------------------------------------------------------------------------
// FetchBlobContainer tests
// ---------------------------------------------------------------------------

func TestFetchBlobContainer_Success(t *testing.T) {
	t.Parallel()

	srv := &storagefake.BlobContainersServer{
		Get: func(_ context.Context, _, _, _ string, _ *armstorage.BlobContainersClientGetOptions) (resp azfake.Responder[armstorage.BlobContainersClientGetResponse], errResp azfake.ErrorResponder) {
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
	container, err := azure.FetchBlobContainer(t.Context(), client, "rg", "teststorage", "testcontainer")

	require.NoError(t, err)
	assert.Equal(t, "testcontainer", *container.Name)
}

// ---------------------------------------------------------------------------
// FetchFileShare tests
// ---------------------------------------------------------------------------

func TestFetchFileShare_Success(t *testing.T) {
	t.Parallel()

	srv := &storagefake.FileSharesServer{
		Get: func(_ context.Context, _, _, _ string, _ *armstorage.FileSharesClientGetOptions) (resp azfake.Responder[armstorage.FileSharesClientGetResponse], errResp azfake.ErrorResponder) {
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
	share, err := azure.FetchFileShare(t.Context(), client, "rg", "teststorage", "testshare")

	require.NoError(t, err)
	assert.Equal(t, "testshare", *share.Name)
}

// ---------------------------------------------------------------------------
// ExtractBlobContainerPublicAccess tests
// ---------------------------------------------------------------------------

func TestExtractBlobContainerPublicAccess_None(t *testing.T) {
	t.Parallel()

	container := &armstorage.BlobContainer{
		ContainerProperties: &armstorage.ContainerProperties{
			PublicAccess: to.Ptr(armstorage.PublicAccessNone),
		},
	}
	assert.False(t, azure.ExtractBlobContainerPublicAccess(container))
}

func TestExtractBlobContainerPublicAccess_Blob(t *testing.T) {
	t.Parallel()

	container := &armstorage.BlobContainer{
		ContainerProperties: &armstorage.ContainerProperties{
			PublicAccess: to.Ptr(armstorage.PublicAccessBlob),
		},
	}
	assert.True(t, azure.ExtractBlobContainerPublicAccess(container))
}

// ---------------------------------------------------------------------------
// ExtractStorageAccountKind tests
// ---------------------------------------------------------------------------

func TestExtractStorageAccountKind(t *testing.T) {
	t.Parallel()

	account := &armstorage.Account{
		Kind: to.Ptr(armstorage.KindStorageV2),
	}
	assert.Equal(t, "StorageV2", azure.ExtractStorageAccountKind(account))
}

// ---------------------------------------------------------------------------
// ExtractStorageAccountSkuTier tests
// ---------------------------------------------------------------------------

func TestExtractStorageAccountSkuTier(t *testing.T) {
	t.Parallel()

	account := &armstorage.Account{
		SKU: &armstorage.SKU{
			Tier: to.Ptr(armstorage.SKUTierStandard),
		},
	}
	assert.Equal(t, "Standard", azure.ExtractStorageAccountSkuTier(account))
}

// ---------------------------------------------------------------------------
// FetchStorageAccountProperties tests
// ---------------------------------------------------------------------------

func TestFetchStorageAccountProperties(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    storagefake.AccountsServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: storagefake.AccountsServer{
				GetProperties: func(_ context.Context, _, _ string, _ *armstorage.AccountsClientGetPropertiesOptions) (resp azfake.Responder[armstorage.AccountsClientGetPropertiesResponse], errResp azfake.ErrorResponder) {
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
			},
			wantName: "teststorage",
		},
		{
			name: "NotFound",
			server: storagefake.AccountsServer{
				GetProperties: func(_ context.Context, _, _ string, _ *armstorage.AccountsClientGetPropertiesOptions) (resp azfake.Responder[armstorage.AccountsClientGetPropertiesResponse], errResp azfake.ErrorResponder) {
					errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

					return
				},
			},
			wantErr:   true,
			errSubstr: "ResourceNotFound",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := tc.server
			client := newFakeStorageAccountsClient(t, &srv)

			account, err := azure.FetchStorageAccountProperties(t.Context(), client, "rg", "teststorage")
			if tc.wantErr {
				require.Error(t, err)

				var respErr *azcore.ResponseError
				require.ErrorAs(t, err, &respErr)
				assert.Equal(t, tc.errSubstr, respErr.ErrorCode)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *account.Name)
		})
	}
}

// ---------------------------------------------------------------------------
// FetchBlobContainer tests
// ---------------------------------------------------------------------------

func TestFetchBlobContainer(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    storagefake.BlobContainersServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: storagefake.BlobContainersServer{
				Get: func(_ context.Context, _, _, _ string, _ *armstorage.BlobContainersClientGetOptions) (resp azfake.Responder[armstorage.BlobContainersClientGetResponse], errResp azfake.ErrorResponder) {
					result := armstorage.BlobContainersClientGetResponse{
						BlobContainer: armstorage.BlobContainer{
							Name: to.Ptr("testcontainer"),
						},
					}
					resp.SetResponse(http.StatusOK, result, nil)

					return
				},
			},
			wantName: "testcontainer",
		},
		{
			name: "NotFound",
			server: storagefake.BlobContainersServer{
				Get: func(_ context.Context, _, _, _ string, _ *armstorage.BlobContainersClientGetOptions) (resp azfake.Responder[armstorage.BlobContainersClientGetResponse], errResp azfake.ErrorResponder) {
					errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

					return
				},
			},
			wantErr:   true,
			errSubstr: "ResourceNotFound",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := tc.server
			client := newFakeBlobContainersClient(t, &srv)

			container, err := azure.FetchBlobContainer(t.Context(), client, "rg", "teststorage", "testcontainer")
			if tc.wantErr {
				require.Error(t, err)

				var respErr *azcore.ResponseError
				require.ErrorAs(t, err, &respErr)
				assert.Equal(t, tc.errSubstr, respErr.ErrorCode)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *container.Name)
		})
	}
}

// ---------------------------------------------------------------------------
// FetchFileShare tests
// ---------------------------------------------------------------------------

func TestFetchFileShare(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name      string
		wantName  string
		errSubstr string
		server    storagefake.FileSharesServer
		wantErr   bool
	}{
		{
			name: "Success",
			server: storagefake.FileSharesServer{
				Get: func(_ context.Context, _, _, _ string, _ *armstorage.FileSharesClientGetOptions) (resp azfake.Responder[armstorage.FileSharesClientGetResponse], errResp azfake.ErrorResponder) {
					result := armstorage.FileSharesClientGetResponse{
						FileShare: armstorage.FileShare{
							Name: to.Ptr("testshare"),
						},
					}
					resp.SetResponse(http.StatusOK, result, nil)

					return
				},
			},
			wantName: "testshare",
		},
		{
			name: "NotFound",
			server: storagefake.FileSharesServer{
				Get: func(_ context.Context, _, _, _ string, _ *armstorage.FileSharesClientGetOptions) (resp azfake.Responder[armstorage.FileSharesClientGetResponse], errResp azfake.ErrorResponder) {
					errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

					return
				},
			},
			wantErr:   true,
			errSubstr: "ResourceNotFound",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srv := tc.server
			client := newFakeFileSharesClient(t, &srv)

			share, err := azure.FetchFileShare(t.Context(), client, "rg", "teststorage", "testshare")
			if tc.wantErr {
				require.Error(t, err)

				var respErr *azcore.ResponseError
				require.ErrorAs(t, err, &respErr)
				assert.Equal(t, tc.errSubstr, respErr.ErrorCode)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, *share.Name)
		})
	}
}

// ---------------------------------------------------------------------------
// Fake client helpers
// ---------------------------------------------------------------------------

func newFakeStorageAccountsClient(t *testing.T, srv *storagefake.AccountsServer) *armstorage.AccountsClient {
	t.Helper()

	client, err := armstorage.NewAccountsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: storagefake.NewAccountsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeBlobContainersClient(t *testing.T, srv *storagefake.BlobContainersServer) *armstorage.BlobContainersClient {
	t.Helper()

	client, err := armstorage.NewBlobContainersClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: storagefake.NewBlobContainersServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeFileSharesClient(t *testing.T, srv *storagefake.FileSharesServer) *armstorage.FileSharesClient {
	t.Helper()

	client, err := armstorage.NewFileSharesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: storagefake.NewFileSharesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}
