package azure_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	azure "github.com/gruntwork-io/terratest/modules/azure"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods to create and delete recovery services resources are added, these tests can be extended.
*/

func TestRecoveryServicesVaultName(t *testing.T) {
	t.Parallel()

	_, err := azure.GetRecoveryServicesVaultE("", "", "")
	require.Error(t, err, "vault")
}

func TestRecoveryServicesVaultExists(t *testing.T) {
	t.Parallel()

	_, err := azure.RecoveryServicesVaultExistsE("", "", "")
	require.Error(t, err, "vault exists")
}

func TestRecoveryServicesVaultBackupPolicyList(t *testing.T) {
	t.Parallel()

	_, err := azure.GetRecoveryServicesVaultBackupPolicyListE("", "", "")
	require.Error(t, err, "Backup policy list not faulted")
}

func TestRecoveryServicesVaultBackupProtectedVMList(t *testing.T) {
	t.Parallel()

	_, err := azure.GetRecoveryServicesVaultBackupProtectedVMListE("", "", "", "")
	require.Error(t, err, "Backup policy protected vm list not faulted")
}
