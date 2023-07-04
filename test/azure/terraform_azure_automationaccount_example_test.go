// +build azureautomation

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

func TestTerraformAzureAutomationAccountExample(t *testing.T) {
	t.Parallel()

	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	subscriptionID := ""
	uniquePostfix := random.UniqueId()
	expectedAutomationAccountName := "terratest-AutomationAccount"
	expectedSampleDscName := "SampleDSC"
	expectedSampleDscConfigurationName := "SampleDSC.NotWebServer"
	expectedVMNodeHostName := "dscnode"
	expectedRunAsAccountName := "terratest-AutomationRunAsConnectionName"
	expectedRunAsType := "AzureServicePrincipal"
	expectedRunAsCertificateName := "terratest-AutomationConnectionCertificateName"
	expectedRunAsCertificateThumbprint := `env:"TF_VAR_automation_run_as_certificate_thumbprint"`
	// Construct options for TF apply
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-automationaccount-example",
		Vars: map[string]interface{}{
			"postfix":                                  uniquePostfix,
			"automation_account_name":                  expectedAutomationAccountName,
			"sample_dsc_name":                          expectedSampleDscName,
			"automation_run_as_connection_name":        expectedRunAsAccountName,
			"automation_run_as_connection_type":        expectedRunAsType,
			"automation_run_as_certificate_name":       expectedRunAsCertificateName,
			"automation_run_as_certificate_thumbprint": expectedRunAsCertificateThumbprint,
			"vm_host_name":                             expectedVMNodeHostName,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	automationAccountName := terraform.Output(t, terraformOptions, "automation_account_name")

	// Check that the automation account deployed successfully
	actualAutomationAccountExists := azure.AutomationAccountExists(t, automationAccountName, resourceGroupName, subscriptionID)
	assert.True(t, actualAutomationAccountExists)
	// Check that the Run As configuraiton is valid
	runAsAccountValidates := azure.AutomationAccountRunAsConnectionExists(t, expectedRunAsAccountName+"-"+uniquePostfix, expectedRunAsType, expectedRunAsCertificateThumbprint, automationAccountName, resourceGroupName, subscriptionID)
	assert.True(t, runAsAccountValidates)
	// Check that the sample DSC was uploaded successfully into the deployed automation account
	actualDSCExists := azure.AutomationAccountDscExists(t, expectedSampleDscName, automationAccountName, resourceGroupName, subscriptionID)
	assert.True(t, actualDSCExists)
	// Check that the DSC in the automation account successfully compiled
	dscCompiled := azure.WaitUntilDscCompiled(t, expectedSampleDscName, automationAccountName, resourceGroupName, subscriptionID)
	assert.True(t, dscCompiled)
	// Check that the DSC was successfully configured on the VM node
	dscConfiguredOnNode := azure.AutomationAccountDscAppliedSuccessfullyToVM(t, expectedSampleDscConfigurationName, expectedVMNodeHostName+"-"+uniquePostfix, automationAccountName, resourceGroupName, subscriptionID)
	assert.True(t, dscConfiguredOnNode)
}
