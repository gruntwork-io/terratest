package azure

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2016-06-01/recoveryservices"
	"github.com/Azure/azure-sdk-for-go/services/recoveryservices/mgmt/2020-02-02/backup"
	"github.com/stretchr/testify/require"
)

// RecoveryServicesVaultExists indicates whether a recovery services vault exists; otherwise false.
// This function would fail the test if there is an error.
func RecoveryServicesVaultExists(t *testing.T, vaultName, resourceGroupName, subscriptionID string) bool {
	exists, err := RecoveryServicesVaultExistsE(vaultName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// GetRecoveryServicesVaultBackupPolicyList returns a list of backup policies for the given vault.
// This function would fail the test if there is an error.
func GetRecoveryServicesVaultBackupPolicyList(t *testing.T, vaultName, resourceGroupName, subscriptionID string) map[string]backup.ProtectionPolicyResource {
	list, err := GetRecoveryServicesVaultBackupPolicyListE(vaultName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return list
}

// GetRecoveryServicesVaultBackupProtectedVMList returns a list of protected VM's on the given vault/policy.
// This function would fail the test if there is an error.
func GetRecoveryServicesVaultBackupProtectedVMList(t *testing.T, policyName, vaultName, resourceGroupName, subscriptionID string) map[string]backup.AzureIaaSComputeVMProtectedItem {
	list, err := GetRecoveryServicesVaultBackupProtectedVMListE(policyName, vaultName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return list
}

// RecoveryServicesVaultExistsE indicates whether a recovery services vault exists; otherwise false or error.
func RecoveryServicesVaultExistsE(vaultName, resourceGroupName, subscriptionID string) (bool, error) {
	_, err := GetRecoveryServicesVaultE(vaultName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetRecoveryServicesVaultE returns a vault instance.
func GetRecoveryServicesVaultE(vaultName, resourceGroupName, subscriptionID string) (*recoveryservices.Vault, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	resourceGroupName, err2 := getTargetAzureResourceGroupName((resourceGroupName))
	if err2 != nil {
		return nil, err2
	}

	client := recoveryservices.NewVaultsClient(subscriptionID)
	// setup auth and create request params
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer
	vault, err := client.Get(context.Background(), resourceGroupName, vaultName)
	if err != nil {
		return nil, err
	}
	return &vault, nil
}

// GetRecoveryServicesVaultBackupPolicyListE returns a list of backup policies for the given vault.
func GetRecoveryServicesVaultBackupPolicyListE(vaultName, resourceGroupName, subscriptionID string) (map[string]backup.ProtectionPolicyResource, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	resourceGroupName, err2 := getTargetAzureResourceGroupName(resourceGroupName)
	if err2 != nil {
		return nil, err2
	}

	client := backup.NewPoliciesClient(subscriptionID)
	// setup authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer
	listIter, err := client.ListComplete(context.Background(), vaultName, resourceGroupName, "")
	if err != nil {
		return nil, err
	}

	policyMap := make(map[string]backup.ProtectionPolicyResource)
	for listIter.NotDone() {
		v := listIter.Value()
		policyMap[*v.Name] = v
		err := listIter.NextWithContext(context.Background())
		if err != nil {
			return nil, err
		}

	}
	return policyMap, nil
}

// GetRecoveryServicesVaultBackupProtectedVMListE returns a list of protected VM's on the given vault/policy.
func GetRecoveryServicesVaultBackupProtectedVMListE(policyName, vaultName, resourceGroupName, subscriptionID string) (map[string]backup.AzureIaaSComputeVMProtectedItem, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	resourceGroupName, err = getTargetAzureResourceGroupName(resourceGroupName)
	if err != nil {
		return nil, err
	}

	client := backup.NewProtectedItemsGroupClient(subscriptionID)
	// setup authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}
	client.Authorizer = *authorizer
	// Build a filter string to narrow down results to just VM's
	filter := fmt.Sprintf("backupManagementType eq 'AzureIaasVM' and itemType eq 'VM' and policyName eq '%s'", policyName)
	listIter, err := client.ListComplete(context.Background(), vaultName, resourceGroupName, filter, "")
	if err != nil {
		return nil, err
	}
	// Prep the return container
	vmList := make(map[string]backup.AzureIaaSComputeVMProtectedItem)
	// First iterator check
	for listIter.NotDone() {
		currentVM, _ := listIter.Value().Properties.AsAzureIaaSComputeVMProtectedItem()
		vmList[*currentVM.FriendlyName] = *currentVM
		err := listIter.NextWithContext(context.Background())
		if err != nil {
			return nil, err
		}
	}
	return vmList, nil
}
