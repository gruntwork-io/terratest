package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/recoveryservices/armrecoveryservices/v2"
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

// GetRecoveryServicesVaultBackupPolicyListContext returns a list of backup policies for the given vault.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRecoveryServicesVaultBackupPolicyListContext(t testing.TestingT, ctx context.Context, vaultName, resourceGroupName, subscriptionID string) map[string]armrecoveryservicesbackup.ProtectionPolicyResource {
	t.Helper()

	list, err := GetRecoveryServicesVaultBackupPolicyListContextE(ctx, vaultName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return list
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

	client, err := CreateRecoveryServicesVaultsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, resourceGroupName, vaultName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Vault, nil
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

	client, err := CreateBackupPoliciesClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	pager := client.NewListPager(vaultName, resourceGroupName, nil)
	policyMap := make(map[string]armrecoveryservicesbackup.ProtectionPolicyResource)

	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Value {
			policyMap[*v.Name] = *v
		}
	}

	return policyMap, nil
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

	client, err := CreateBackupProtectedItemsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

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
			if vmItem, ok := item.Properties.(*armrecoveryservicesbackup.AzureIaaSComputeVMProtectedItem); ok {
				vmList[*vmItem.FriendlyName] = *vmItem
			}
		}
	}

	return vmList, nil
}
