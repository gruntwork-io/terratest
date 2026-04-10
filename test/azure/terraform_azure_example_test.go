//go:build azure || (azureslim && compute)
// +build azure azureslim,compute

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureExample(t *testing.T) {
	t.Parallel()

	uniquePostfix := random.UniqueID()

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	vmName := terraform.OutputContext(t, t.Context(), terraformOptions, "vm_name")
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")

	// website::tag::4:: Look up the size of the given Virtual Machine and ensure it matches the output.
	actualVMSize := azure.GetSizeOfVirtualMachineContext(t, t.Context(), vmName, resourceGroupName, "")
	expectedVMSize := armcompute.VirtualMachineSizeTypes("Standard_B1s")
	assert.Equal(t, expectedVMSize, actualVMSize)
}
