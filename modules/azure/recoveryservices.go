package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservices"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservicesbackup/v4"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// RecoveryServicesVaultExistsContext indicates whether a recovery services vault exists; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func RecoveryServicesVaultExistsContext(t testing.TestingT, ctx context.Context, vaultName, resourceGroupName, subscriptionID string) bool {
	t.Helper()

	exists, err := RecoveryServicesVaultExistsContextE(ctx, vaultName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// RecoveryServicesVaultExists indicates whether a recovery services vault exists; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [RecoveryServicesVaultExistsContext] instead.
func RecoveryServicesVaultExists(t testing.TestingT, vaultName, resourceGroupName, subscriptionID string) bool {
	t.Helper()

	return RecoveryServicesVaultExistsContext(t, context.Background(), vaultName, resourceGroupName, subscriptionID)
}

// GetRecoveryServicesVaultBackupPolicyListContext returns a list of backup policies for the given vault.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRecoveryServicesVaultBackupPolicyListContext(t testing.TestingT, ctx context.Context, vaultName, resourceGroupName, subscriptionID string) map[string]armrecoveryservicesbackup.ProtectionPolicyResource {
	t.Helper()

	list, err := GetRecoveryServicesVaultBackupPolicyListContextE(ctx, vaultName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return list
}

// GetRecoveryServicesVaultBackupPolicyList returns a list of backup policies for the given vault.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetRecoveryServicesVaultBackupPolicyListContext] instead.
func GetRecoveryServicesVaultBackupPolicyList(t testing.TestingT, vaultName, resourceGroupName, subscriptionID string) map[string]armrecoveryservicesbackup.ProtectionPolicyResource {
	t.Helper()

	return GetRecoveryServicesVaultBackupPolicyListContext(t, context.Background(), vaultName, resourceGroupName, subscriptionID)
}

// GetRecoveryServicesVaultBackupProtectedVMListContext returns a list of protected VMs on the given vault and policy.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRecoveryServicesVaultBackupProtectedVMListContext(t testing.TestingT, ctx context.Context, policyName, vaultName, resourceGroupName, subscriptionID string) map[string]armrecoveryservicesbackup.AzureIaaSComputeVMProtectedItem {
	t.Helper()

	list, err := GetRecoveryServicesVaultBackupProtectedVMListContextE(ctx, policyName, vaultName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return list
}

// GetRecoveryServicesVaultBackupProtectedVMList returns a list of protected VM's on the given vault/policy.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetRecoveryServicesVaultBackupProtectedVMListContext] instead.
func GetRecoveryServicesVaultBackupProtectedVMList(t testing.TestingT, policyName, vaultName, resourceGroupName, subscriptionID string) map[string]armrecoveryservicesbackup.AzureIaaSComputeVMProtectedItem {
	t.Helper()

	return GetRecoveryServicesVaultBackupProtectedVMListContext(t, context.Background(), policyName, vaultName, resourceGroupName, subscriptionID)
}

// RecoveryServicesVaultExistsContextE indicates whether a recovery services vault exists; otherwise false or error.
// The ctx parameter supports cancellation and timeouts.
func RecoveryServicesVaultExistsContextE(ctx context.Context, vaultName, resourceGroupName, subscriptionID string) (bool, error) {
	_, err := GetRecoveryServicesVaultContextE(ctx, vaultName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// RecoveryServicesVaultExistsE indicates whether a recovery services vault exists; otherwise false or error.
//
// Deprecated: Use [RecoveryServicesVaultExistsContextE] instead.
func RecoveryServicesVaultExistsE(vaultName, resourceGroupName, subscriptionID string) (bool, error) {
	return RecoveryServicesVaultExistsContextE(context.Background(), vaultName, resourceGroupName, subscriptionID)
}

// GetRecoveryServicesVaultContextE returns a recovery services vault instance.
// The ctx parameter supports cancellation and timeouts.
func GetRecoveryServicesVaultContextE(ctx context.Context, vaultName, resourceGroupName, subscriptionID string) (*armrecoveryservices.Vault, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	resourceGroupName, err = getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	client, err := armrecoveryservices.NewVaultsClient(subscriptionID, cred, opts)
	if err != nil {
		return nil, err
	}

	return GetRecoveryServicesVaultWithClient(ctx, client, resourceGroupName, vaultName)
}

// GetRecoveryServicesVaultE returns a vault instance.
//
// Deprecated: Use [GetRecoveryServicesVaultContextE] instead.
func GetRecoveryServicesVaultE(vaultName, resourceGroupName, subscriptionID string) (*armrecoveryservices.Vault, error) {
	return GetRecoveryServicesVaultContextE(context.Background(), vaultName, resourceGroupName, subscriptionID)
}

// GetRecoveryServicesVaultBackupPolicyListContextE returns a list of backup policies for the given vault.
// The ctx parameter supports cancellation and timeouts.
func GetRecoveryServicesVaultBackupPolicyListContextE(ctx context.Context, vaultName, resourceGroupName, subscriptionID string) (map[string]armrecoveryservicesbackup.ProtectionPolicyResource, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	resourceGroupName, err = getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	client, err := armrecoveryservicesbackup.NewBackupPoliciesClient(subscriptionID, cred, opts)
	if err != nil {
		return nil, err
	}

	return GetBackupPolicyListWithClient(ctx, client, vaultName, resourceGroupName)
}

// GetRecoveryServicesVaultBackupPolicyListE returns a list of backup policies for the given vault.
//
// Deprecated: Use [GetRecoveryServicesVaultBackupPolicyListContextE] instead.
func GetRecoveryServicesVaultBackupPolicyListE(vaultName, resourceGroupName, subscriptionID string) (map[string]armrecoveryservicesbackup.ProtectionPolicyResource, error) {
	return GetRecoveryServicesVaultBackupPolicyListContextE(context.Background(), vaultName, resourceGroupName, subscriptionID)
}

// GetRecoveryServicesVaultBackupProtectedVMListContextE returns a list of protected VMs on the given vault and policy.
// The ctx parameter supports cancellation and timeouts.
func GetRecoveryServicesVaultBackupProtectedVMListContextE(ctx context.Context, policyName, vaultName, resourceGroupName, subscriptionID string) (map[string]armrecoveryservicesbackup.AzureIaaSComputeVMProtectedItem, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	resourceGroupName, err = getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	client, err := armrecoveryservicesbackup.NewBackupProtectedItemsClient(subscriptionID, cred, opts)
	if err != nil {
		return nil, err
	}

	return GetBackupProtectedVMListWithClient(ctx, client, vaultName, resourceGroupName, policyName)
}

// GetRecoveryServicesVaultBackupProtectedVMListE returns a list of protected VM's on the given vault/policy.
//
// Deprecated: Use [GetRecoveryServicesVaultBackupProtectedVMListContextE] instead.
func GetRecoveryServicesVaultBackupProtectedVMListE(policyName, vaultName, resourceGroupName, subscriptionID string) (map[string]armrecoveryservicesbackup.AzureIaaSComputeVMProtectedItem, error) {
	return GetRecoveryServicesVaultBackupProtectedVMListContextE(context.Background(), policyName, vaultName, resourceGroupName, subscriptionID)
}

// GetRecoveryServicesVaultWithClient retrieves a recovery services vault using the provided client.
func GetRecoveryServicesVaultWithClient(ctx context.Context, client *armrecoveryservices.VaultsClient, resourceGroupName, vaultName string) (*armrecoveryservices.Vault, error) {
	resp, err := client.Get(ctx, resourceGroupName, vaultName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Vault, nil
}

// GetBackupPolicyListWithClient retrieves all backup policies for a vault using the provided client.
func GetBackupPolicyListWithClient(ctx context.Context, client *armrecoveryservicesbackup.BackupPoliciesClient, vaultName, resourceGroupName string) (map[string]armrecoveryservicesbackup.ProtectionPolicyResource, error) {
	pager := client.NewListPager(vaultName, resourceGroupName, nil)
	policyMap := make(map[string]armrecoveryservicesbackup.ProtectionPolicyResource)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Value {
			if v == nil || v.Name == nil {
				continue
			}

			policyMap[*v.Name] = *v
		}
	}

	return policyMap, nil
}

// GetBackupProtectedVMListWithClient retrieves all protected VMs matching the given policy using the provided client.
func GetBackupProtectedVMListWithClient(ctx context.Context, client *armrecoveryservicesbackup.BackupProtectedItemsClient, vaultName, resourceGroupName, policyName string) (map[string]armrecoveryservicesbackup.AzureIaaSComputeVMProtectedItem, error) {
	filter := fmt.Sprintf("backupManagementType eq 'AzureIaasVM' and itemType eq 'VM' and policyName eq '%s'", policyName)

	pager := client.NewListPager(vaultName, resourceGroupName, &armrecoveryservicesbackup.BackupProtectedItemsClientListOptions{
		Filter: &filter,
	})

	vmList := make(map[string]armrecoveryservicesbackup.AzureIaaSComputeVMProtectedItem)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, item := range page.Value {
			if item == nil || item.Properties == nil {
				continue
			}

			vmItem, ok := item.Properties.(*armrecoveryservicesbackup.AzureIaaSComputeVMProtectedItem)
			if !ok || vmItem == nil || vmItem.FriendlyName == nil {
				continue
			}

			vmList[*vmItem.FriendlyName] = *vmItem
		}
	}

	return vmList, nil
}
