//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

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
	uniquePostfix := random.UniqueId()

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
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	vaultName := terraform.Output(t, terraformOptions, "recovery_service_vault_name")
	policyVmName := terraform.Output(t, terraformOptions, "backup_policy_vm_name")

	// website::tag::4:: Verify the recovery services resources
	exists := azure.RecoveryServicesVaultExists(t, vaultName, resourceGroupName, subscriptionID)
	assert.True(t, exists, "vault does not exist")

	policyList := azure.GetRecoveryServicesVaultBackupPolicyList(t, vaultName, resourceGroupName, subscriptionID)
	assert.NotNil(t, policyList, "vault backup policy list is nil")

	vmPolicyList := azure.GetRecoveryServicesVaultBackupProtectedVMList(t, policyVmName, vaultName, resourceGroupName, subscriptionID)
	assert.NotNil(t, vmPolicyList, "vault backup policy list for protected vm is nil")
}
