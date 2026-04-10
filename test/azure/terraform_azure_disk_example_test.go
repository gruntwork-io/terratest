//go:build azure
// +build azure

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

func TestTerraformAzureDiskExample(t *testing.T) {
	t.Parallel()

	// Subscription ID, leave blank if available as an Environment Var
	subID := ""
	uniquePostfix := random.UniqueID()

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-disk-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	expectedDiskName := terraform.OutputContext(t, t.Context(), terraformOptions, "disk_name")
	expectedDiskType := terraform.OutputContext(t, t.Context(), terraformOptions, "disk_type")

	// Check the Disk Type
	actualDisk := azure.GetDiskContext(t, t.Context(), expectedDiskName, resourceGroupName, subID)
	assert.Equal(t, armcompute.DiskStorageAccountTypes(expectedDiskType), *actualDisk.SKU.Name)
}
