package azure

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2020-02-02/backup"
	"github.com/stretchr/testify/require"
)

// RecoveryServicesVaultExistsContext indicates whether a recovery services vault exists; otherwise false.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func RecoveryServicesVaultExistsContext(t *testing.T, ctx context.Context, vaultName, resourceGroupName, subscriptionID string) bool {
	t.Helper()

	exists, err := RecoveryServicesVaultExistsContextE(ctx, vaultName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// RecoveryServicesVaultExists indicates whether a recovery services vault exists; otherwise false.
// This function would fail the test if there is an error.
//
// Deprecated: Use [RecoveryServicesVaultExistsContext] instead.
func RecoveryServicesVaultExists(t *testing.T, vaultName, resourceGroupName, subscriptionID string) bool {
	t.Helper()

	return RecoveryServicesVaultExistsContext(t, context.Background(), vaultName, resourceGroupName, subscriptionID)
}

// GetRecoveryServicesVaultBackupPolicyListContext returns a list of backup policies for the given vault.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRecoveryServicesVaultBackupPolicyListContext(t *testing.T, ctx context.Context, vaultName, resourceGroupName, subscriptionID string) map[string]backup.ProtectionPolicyResource {
	t.Helper()

	list, err := GetRecoveryServicesVaultBackupPolicyListContextE(ctx, vaultName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return list
}

// GetRecoveryServicesVaultBackupPolicyList returns a list of backup policies for the given vault.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetRecoveryServicesVaultBackupPolicyListContext] instead.
func GetRecoveryServicesVaultBackupPolicyList(t *testing.T, vaultName, resourceGroupName, subscriptionID string) map[string]backup.ProtectionPolicyResource {
	t.Helper()

	return GetRecoveryServicesVaultBackupPolicyListContext(t, context.Background(), vaultName, resourceGroupName, subscriptionID)
}

// GetRecoveryServicesVaultBackupProtectedVMListContext returns a list of protected VMs on the given vault and policy.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRecoveryServicesVaultBackupProtectedVMListContext(t *testing.T, ctx context.Context, policyName, vaultName, resourceGroupName, subscriptionID string) map[string]backup.AzureIaaSComputeVMProtectedItem {
	t.Helper()

	list, err := GetRecoveryServicesVaultBackupProtectedVMListContextE(ctx, policyName, vaultName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return list
}

// GetRecoveryServicesVaultBackupProtectedVMList returns a list of protected VM's on the given vault/policy.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetRecoveryServicesVaultBackupProtectedVMListContext] instead.
func GetRecoveryServicesVaultBackupProtectedVMList(t *testing.T, policyName, vaultName, resourceGroupName, subscriptionID string) map[string]backup.AzureIaaSComputeVMProtectedItem {
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
func GetRecoveryServicesVaultContextE(ctx context.Context, vaultName, resourceGroupName, subscriptionID string) (*recoveryservices.Vault, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	resourceGroupName, err = getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	client := recoveryservices.NewVaultsClient(subscriptionID)

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	vault, err := client.Get(ctx, resourceGroupName, vaultName)
	if err != nil {
		return nil, err
	}

	return &vault, nil
}

// GetRecoveryServicesVaultE returns a vault instance.
//
// Deprecated: Use [GetRecoveryServicesVaultContextE] instead.
func GetRecoveryServicesVaultE(vaultName, resourceGroupName, subscriptionID string) (*recoveryservices.Vault, error) {
	return GetRecoveryServicesVaultContextE(context.Background(), vaultName, resourceGroupName, subscriptionID)
}

// GetRecoveryServicesVaultBackupPolicyListContextE returns a list of backup policies for the given vault.
// The ctx parameter supports cancellation and timeouts.
func GetRecoveryServicesVaultBackupPolicyListContextE(ctx context.Context, vaultName, resourceGroupName, subscriptionID string) (map[string]backup.ProtectionPolicyResource, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	resourceGroupName, err = getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	client := backup.NewPoliciesClient(subscriptionID)

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	listIter, err := client.ListComplete(ctx, vaultName, resourceGroupName, "")
	if err != nil {
		return nil, err
	}

	policyMap := make(map[string]backup.ProtectionPolicyResource)

	for listIter.NotDone() {
		v := listIter.Value()
		policyMap[*v.Name] = v

		err := listIter.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
	}

	return policyMap, nil
}

// GetRecoveryServicesVaultBackupPolicyListE returns a list of backup policies for the given vault.
//
// Deprecated: Use [GetRecoveryServicesVaultBackupPolicyListContextE] instead.
func GetRecoveryServicesVaultBackupPolicyListE(vaultName, resourceGroupName, subscriptionID string) (map[string]backup.ProtectionPolicyResource, error) {
	return GetRecoveryServicesVaultBackupPolicyListContextE(context.Background(), vaultName, resourceGroupName, subscriptionID)
}

// GetRecoveryServicesVaultBackupProtectedVMListContextE returns a list of protected VMs on the given vault and policy.
// The ctx parameter supports cancellation and timeouts.
func GetRecoveryServicesVaultBackupProtectedVMListContextE(ctx context.Context, policyName, vaultName, resourceGroupName, subscriptionID string) (map[string]backup.AzureIaaSComputeVMProtectedItem, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	resourceGroupName, err = getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	client := backup.NewProtectedItemsGroupClient(subscriptionID)

	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	filter := fmt.Sprintf("backupManagementType eq 'AzureIaasVM' and itemType eq 'VM' and policyName eq '%s'", policyName)

	listIter, err := client.ListComplete(ctx, vaultName, resourceGroupName, filter, "")
	if err != nil {
		return nil, err
	}

	vmList := make(map[string]backup.AzureIaaSComputeVMProtectedItem)

	for listIter.NotDone() {
		currentVM, _ := listIter.Value().Properties.AsAzureIaaSComputeVMProtectedItem()
		vmList[*currentVM.FriendlyName] = *currentVM

		err := listIter.NextWithContext(ctx)
		if err != nil {
			return nil, err
		}
	}

	return vmList, nil
}

// GetRecoveryServicesVaultBackupProtectedVMListE returns a list of protected VM's on the given vault/policy.
//
// Deprecated: Use [GetRecoveryServicesVaultBackupProtectedVMListContextE] instead.
func GetRecoveryServicesVaultBackupProtectedVMListE(policyName, vaultName, resourceGroupName, subscriptionID string) (map[string]backup.AzureIaaSComputeVMProtectedItem, error) {
	return GetRecoveryServicesVaultBackupProtectedVMListContextE(context.Background(), policyName, vaultName, resourceGroupName, subscriptionID)
}
