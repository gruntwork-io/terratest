//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureAvailabilitySetExample(t *testing.T) {
	t.Parallel()

	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	subscriptionID := ""
	uniquePostfix := random.UniqueId()
	expectedAvsName := fmt.Sprintf("avs-%s", uniquePostfix)
	expectedVMName := fmt.Sprintf("vm-%s", uniquePostfix)
	var expectedAvsFaultDomainCount int32 = 3

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// Relative path to the Terraform dir
		TerraformDir: "../../examples/azure/terraform-azure-availabilityset-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"postfix":                uniquePostfix,
			"avs_fault_domain_count": expectedAvsFaultDomainCount,
			// "location": "East US",
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")

	// Check the Availability Set Exists
	actualAvsExists := azure.AvailabilitySetExists(t, expectedAvsName, resourceGroupName, subscriptionID)
	assert.True(t, actualAvsExists)

	// Check the Availability Set Fault Domain Count
	actualAvsFaultDomainCount := azure.GetAvailabilitySetFaultDomainCount(t, expectedAvsName, resourceGroupName, subscriptionID)
	assert.Equal(t, expectedAvsFaultDomainCount, actualAvsFaultDomainCount)

	// Check the Availability Set for a VM
	actualVMPresent := azure.CheckAvailabilitySetContainsVM(t, expectedVMName, expectedAvsName, resourceGroupName, subscriptionID)
	assert.True(t, actualVMPresent)
}
