//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureRecoveryServicesExample(t *testing.T) {
	t.Parallel()

	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	subscriptionID := ""
	uniquePostfix := random.UniqueID()

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-recoveryservices-example",
		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	vaultName := terraform.OutputContext(t, t.Context(), terraformOptions, "recovery_service_vault_name")
	policyVMName := terraform.OutputContext(t, t.Context(), terraformOptions, "backup_policy_vm_name")

	// website::tag::4:: Verify the recovery services resources
	exists := azure.RecoveryServicesVaultExistsContext(t, t.Context(), vaultName, resourceGroupName, subscriptionID)
	assert.True(t, exists, "vault does not exist")

	policyList := azure.GetRecoveryServicesVaultBackupPolicyListContext(t, t.Context(), vaultName, resourceGroupName, subscriptionID)
	assert.NotNil(t, policyList, "vault backup policy list is nil")

	vmPolicyList := azure.GetRecoveryServicesVaultBackupProtectedVMListContext(t, t.Context(), policyVMName, vaultName, resourceGroupName, subscriptionID)
	assert.NotNil(t, vmPolicyList, "vault backup policy list for protected vm is nil")
}
