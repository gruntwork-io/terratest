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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservices"
	recoveryfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservices/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup/v4"
	backupfake "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup/v4/fake"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Fake client helpers
// ---------------------------------------------------------------------------

func newFakeVaultsClient(t *testing.T, srv *recoveryfake.VaultsServer) *armrecoveryservices.VaultsClient {
	t.Helper()

	client, err := armrecoveryservices.NewVaultsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: recoveryfake.NewVaultsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeBackupPoliciesClient(t *testing.T, srv *backupfake.BackupPoliciesServer) *armrecoveryservicesbackup.BackupPoliciesClient {
	t.Helper()

	client, err := armrecoveryservicesbackup.NewBackupPoliciesClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: backupfake.NewBackupPoliciesServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

func newFakeBackupProtectedItemsClient(t *testing.T, srv *backupfake.BackupProtectedItemsServer) *armrecoveryservicesbackup.BackupProtectedItemsClient {
	t.Helper()

	client, err := armrecoveryservicesbackup.NewBackupProtectedItemsClient("fake-sub", &azfake.TokenCredential{},
		&arm.ClientOptions{ClientOptions: policy.ClientOptions{
			Transport: backupfake.NewBackupProtectedItemsServerTransport(srv),
		}})
	require.NoError(t, err)

	return client
}

// ---------------------------------------------------------------------------
// GetRecoveryServicesVaultWithClient tests
// ---------------------------------------------------------------------------

func TestGetRecoveryServicesVaultWithClient_Success(t *testing.T) {
	t.Parallel()

	srv := &recoveryfake.VaultsServer{
		Get: func(_ context.Context, _, _ string, _ *armrecoveryservices.VaultsClientGetOptions) (resp azfake.Responder[armrecoveryservices.VaultsClientGetResponse], errResp azfake.ErrorResponder) {
			result := armrecoveryservices.VaultsClientGetResponse{
				Vault: armrecoveryservices.Vault{
					Name: to.Ptr("testvault"),
				},
			}
			resp.SetResponse(http.StatusOK, result, nil)

			return
		},
	}

	client := newFakeVaultsClient(t, srv)
	vault, err := azure.GetRecoveryServicesVaultWithClient(t.Context(), client, "rg", "testvault")

	require.NoError(t, err)
	assert.Equal(t, "testvault", *vault.Name)
}

func TestGetRecoveryServicesVaultWithClient_NotFound(t *testing.T) {
	t.Parallel()

	srv := &recoveryfake.VaultsServer{
		Get: func(_ context.Context, _, _ string, _ *armrecoveryservices.VaultsClientGetOptions) (resp azfake.Responder[armrecoveryservices.VaultsClientGetResponse], errResp azfake.ErrorResponder) {
			errResp.SetResponseError(http.StatusNotFound, "ResourceNotFound")

			return
		},
	}

	client := newFakeVaultsClient(t, srv)
	_, err := azure.GetRecoveryServicesVaultWithClient(t.Context(), client, "rg", "missing")

	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	assert.Equal(t, "ResourceNotFound", respErr.ErrorCode)
}

// ---------------------------------------------------------------------------
// GetBackupPolicyListWithClient tests
// ---------------------------------------------------------------------------

func TestGetBackupPolicyListWithClient_OnePage(t *testing.T) {
	t.Parallel()

	srv := &backupfake.BackupPoliciesServer{
		NewListPager: func(_, _ string, _ *armrecoveryservicesbackup.BackupPoliciesClientListOptions) (resp azfake.PagerResponder[armrecoveryservicesbackup.BackupPoliciesClientListResponse]) {
			resp.AddPage(http.StatusOK, armrecoveryservicesbackup.BackupPoliciesClientListResponse{
				ProtectionPolicyResourceList: armrecoveryservicesbackup.ProtectionPolicyResourceList{
					Value: []*armrecoveryservicesbackup.ProtectionPolicyResource{
						{Name: to.Ptr("policy1")},
						{Name: to.Ptr("policy2")},
					},
				},
			}, nil)

			return
		},
	}

	client := newFakeBackupPoliciesClient(t, srv)
	policies, err := azure.GetBackupPolicyListWithClient(t.Context(), client, "vault", "rg")

	require.NoError(t, err)
	assert.Len(t, policies, 2)
	assert.Contains(t, policies, "policy1")
	assert.Contains(t, policies, "policy2")
}

func TestGetBackupPolicyListWithClient_MultiplePages(t *testing.T) {
	t.Parallel()

	srv := &backupfake.BackupPoliciesServer{
		NewListPager: func(_, _ string, _ *armrecoveryservicesbackup.BackupPoliciesClientListOptions) (resp azfake.PagerResponder[armrecoveryservicesbackup.BackupPoliciesClientListResponse]) {
			resp.AddPage(http.StatusOK, armrecoveryservicesbackup.BackupPoliciesClientListResponse{
				ProtectionPolicyResourceList: armrecoveryservicesbackup.ProtectionPolicyResourceList{
					Value: []*armrecoveryservicesbackup.ProtectionPolicyResource{
						{Name: to.Ptr("policy1")},
					},
				},
			}, nil)
			resp.AddPage(http.StatusOK, armrecoveryservicesbackup.BackupPoliciesClientListResponse{
				ProtectionPolicyResourceList: armrecoveryservicesbackup.ProtectionPolicyResourceList{
					Value: []*armrecoveryservicesbackup.ProtectionPolicyResource{
						{Name: to.Ptr("policy2")},
					},
				},
			}, nil)

			return
		},
	}

	client := newFakeBackupPoliciesClient(t, srv)
	policies, err := azure.GetBackupPolicyListWithClient(t.Context(), client, "vault", "rg")

	require.NoError(t, err)
	assert.Len(t, policies, 2)
	assert.Contains(t, policies, "policy1")
	assert.Contains(t, policies, "policy2")
}

func TestGetBackupPolicyListWithClient_Empty(t *testing.T) {
	t.Parallel()

	srv := &backupfake.BackupPoliciesServer{
		NewListPager: func(_, _ string, _ *armrecoveryservicesbackup.BackupPoliciesClientListOptions) (resp azfake.PagerResponder[armrecoveryservicesbackup.BackupPoliciesClientListResponse]) {
			resp.AddPage(http.StatusOK, armrecoveryservicesbackup.BackupPoliciesClientListResponse{
				ProtectionPolicyResourceList: armrecoveryservicesbackup.ProtectionPolicyResourceList{
					Value: []*armrecoveryservicesbackup.ProtectionPolicyResource{},
				},
			}, nil)

			return
		},
	}

	client := newFakeBackupPoliciesClient(t, srv)
	policies, err := azure.GetBackupPolicyListWithClient(t.Context(), client, "vault", "rg")

	require.NoError(t, err)
	assert.Empty(t, policies)
}

// ---------------------------------------------------------------------------
// GetBackupProtectedVMListWithClient tests
// ---------------------------------------------------------------------------

func TestGetBackupProtectedVMListWithClient_WithVMItems(t *testing.T) {
	t.Parallel()

	srv := &backupfake.BackupProtectedItemsServer{
		NewListPager: func(_, _ string, _ *armrecoveryservicesbackup.BackupProtectedItemsClientListOptions) (resp azfake.PagerResponder[armrecoveryservicesbackup.BackupProtectedItemsClientListResponse]) {
			resp.AddPage(http.StatusOK, armrecoveryservicesbackup.BackupProtectedItemsClientListResponse{
				ProtectedItemResourceList: armrecoveryservicesbackup.ProtectedItemResourceList{
					Value: []*armrecoveryservicesbackup.ProtectedItemResource{
						{
							Properties: &armrecoveryservicesbackup.AzureIaaSComputeVMProtectedItem{
								FriendlyName: to.Ptr("myVM"),
							},
						},
						{
							Properties: &armrecoveryservicesbackup.AzureIaaSComputeVMProtectedItem{
								FriendlyName: to.Ptr("myVM2"),
							},
						},
					},
				},
			}, nil)

			return
		},
	}

	client := newFakeBackupProtectedItemsClient(t, srv)
	vms, err := azure.GetBackupProtectedVMListWithClient(t.Context(), client, "vault", "rg", "policy1")

	require.NoError(t, err)
	assert.Len(t, vms, 2)
	assert.Contains(t, vms, "myVM")
	assert.Contains(t, vms, "myVM2")
}

func TestGetBackupProtectedVMListWithClient_NonVMItemsSkipped(t *testing.T) {
	t.Parallel()

	srv := &backupfake.BackupProtectedItemsServer{
		NewListPager: func(_, _ string, _ *armrecoveryservicesbackup.BackupProtectedItemsClientListOptions) (resp azfake.PagerResponder[armrecoveryservicesbackup.BackupProtectedItemsClientListResponse]) {
			resp.AddPage(http.StatusOK, armrecoveryservicesbackup.BackupProtectedItemsClientListResponse{
				ProtectedItemResourceList: armrecoveryservicesbackup.ProtectedItemResourceList{
					Value: []*armrecoveryservicesbackup.ProtectedItemResource{
						{
							Properties: &armrecoveryservicesbackup.AzureFileshareProtectedItem{
								FriendlyName: to.Ptr("myShare"),
							},
						},
					},
				},
			}, nil)

			return
		},
	}

	client := newFakeBackupProtectedItemsClient(t, srv)
	vms, err := azure.GetBackupProtectedVMListWithClient(t.Context(), client, "vault", "rg", "policy1")

	require.NoError(t, err)
	assert.Empty(t, vms)
}

func TestGetBackupProtectedVMListWithClient_Empty(t *testing.T) {
	t.Parallel()

	srv := &backupfake.BackupProtectedItemsServer{
		NewListPager: func(_, _ string, _ *armrecoveryservicesbackup.BackupProtectedItemsClientListOptions) (resp azfake.PagerResponder[armrecoveryservicesbackup.BackupProtectedItemsClientListResponse]) {
			resp.AddPage(http.StatusOK, armrecoveryservicesbackup.BackupProtectedItemsClientListResponse{
				ProtectedItemResourceList: armrecoveryservicesbackup.ProtectedItemResourceList{
					Value: []*armrecoveryservicesbackup.ProtectedItemResource{},
				},
			}, nil)

			return
		},
	}

	client := newFakeBackupProtectedItemsClient(t, srv)
	vms, err := azure.GetBackupProtectedVMListWithClient(t.Context(), client, "vault", "rg", "policy1")

	require.NoError(t, err)
	assert.Empty(t, vms)
}
