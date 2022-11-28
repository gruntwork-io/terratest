// +build azure azureslim,compute

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureVmScaleSetsExample(t *testing.T) {
	t.Parallel()

	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	subscriptionID := ""
	uniquePostfix := random.UniqueId()

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located.
		TerraformDir: "../../examples/azure/terraform-azure-vm-scalesets-example",

		// Variables to pass to our Terraform code using -var options.
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	scaleSetName := terraform.Output(t, terraformOptions, "scaleset_name")
	expectedVMSize := compute.VirtualMachineSizeTypes(terraform.Output(t, terraformOptions, "scaleset_vm_size"))
	expectedTags := terraform.OutputMap(t, terraformOptions, "scaleset_tags")
	expectedCapacity := terraform.Output(t, terraformOptions, "scaleset_capacity")
	expectedVmNamePrefix := terraform.Output(t, terraformOptions, "scaleset_vm_name_prefix")
	expectedPublicAddressName := terraform.Output(t, terraformOptions, "public_ip_name")
	expectedVNetName := terraform.Output(t, terraformOptions, "virtual_network_name")
	expectedSubnetName := terraform.Output(t, terraformOptions, "subnet_name")
	expectedLBPublicName := terraform.Output(t, terraformOptions, "lb_public_name")

	// Verify VM Scale Sets properties and ensure it matches the output
	t.Run("ScaleSet", func(t *testing.T) {
		vmss := azure.GetVirtualMachineScaleSet(t, scaleSetName, resourceGroupName, subscriptionID)

		// Size of VM instance in a VM Scale Set
		actualVMSize := *vmss.Sku.Name
		assert.Equal(t, expectedVMSize, actualVMSize)

		// VM Scale Set's VM instance name prefix
		assert.Equal(t, expectedVmNamePrefix, *vmss.VirtualMachineProfile.OsProfile.ComputerNamePrefix)

		// VM Scale Set tags
		actualTags := azure.GetVirtualMachineScaleSetTags(t, scaleSetName, resourceGroupName, subscriptionID)
		assert.Equal(t, expectedTags, actualTags)

		// VM Scale Set capacity
		actualCapacity := strconv.FormatInt(*vmss.Sku.Capacity, 10)
		assert.Equal(t, expectedCapacity, actualCapacity)
	})

	// Verify Network properties
	t.Run("Network", func(t *testing.T) {
		actualPublicIP := azure.GetIPOfPublicIPAddressByName(t, expectedPublicAddressName, resourceGroupName, subscriptionID)
		assert.NotNil(t, actualPublicIP)

		actualVnetSubnets := azure.GetVirtualNetworkSubnets(t, expectedVNetName, resourceGroupName, subscriptionID)
		assert.NotNil(t, actualVnetSubnets[expectedVNetName])
	})

	// Verify resource existence
	t.Run("Exists", func(t *testing.T) {
		// Check the Virtual Network exists
		assert.True(t, azure.VirtualNetworkExists(t, expectedVNetName, resourceGroupName, subscriptionID))

		// Check the Subnet exists
		assert.True(t, azure.SubnetExists(t, expectedSubnetName, expectedVNetName, resourceGroupName, subscriptionID))

		// Check the VM Scale Set exists
		assert.True(t, azure.VirtualMachineScaleSetExists(t, scaleSetName, resourceGroupName, subscriptionID))

		// Check the Load Balancer exists
		assert.True(t, azure.LoadBalancerExists(t, expectedLBPublicName, resourceGroupName, subscriptionID))
	})
}
