//go:build azure || (azureslim && compute)
// +build azure azureslim,compute

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureVmExample(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	uniquePostfix := random.UniqueId()

	// Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located.
		TerraformDir: "../../examples/azure/terraform-azure-vm-example",

		// Variables to pass to our Terraform code using -var options.
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created.
	defer terraform.Destroy(t, terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run tests for the Virtual Machine.
	testStrategiesForVMs(t, terraformOptions, subscriptionID)
	testMultipleVMs(t, terraformOptions, subscriptionID)
	testInformationOfVM(t, terraformOptions, subscriptionID)
	testDisksOfVM(t, terraformOptions, subscriptionID)
	testNetworkOfVM(t, terraformOptions, subscriptionID)
}

// These 3 tests check for the same property but illustrate different testing strategies for
// retriving the data. The first strategy is used in the other tests of this module while
// the other two can be extended by the user as needed.
func testStrategiesForVMs(t *testing.T, terraformOptions *terraform.Options, subscriptionID string) {
	// Run `terraform output` to get the values of output variables.
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	virtualMachineName := terraform.Output(t, terraformOptions, "vm_name")
	expectedVMSize := compute.VirtualMachineSizeTypes(terraform.Output(t, terraformOptions, "vm_size"))

	// 1. Check the VM Size directly. This strategy gets one specific property of the VM per method.
	actualVMSize := azure.GetSizeOfVirtualMachine(t, virtualMachineName, resourceGroupName, subscriptionID)
	assert.Equal(t, expectedVMSize, actualVMSize)

	// 2. Check the VM size by reference. This strategy is beneficial when checking multiple properties
	// by using one VM reference. Optional parameters have to be checked first to avoid nil panics.
	vmByRef := azure.GetVirtualMachine(t, virtualMachineName, resourceGroupName, subscriptionID)
	actualVMSize = vmByRef.HardwareProfile.VMSize
	assert.Equal(t, expectedVMSize, actualVMSize)

	// 3. Check the VM size by instance. This strategy is beneficial when checking multiple properties
	// by using one VM instance and making calls against it with the added benefit of property check abstraction.
	vmInstance := azure.Instance{VirtualMachine: vmByRef}
	actualVMSize = vmInstance.GetVirtualMachineInstanceSize()
	assert.Equal(t, expectedVMSize, actualVMSize)
}

// These tests check for the multiple Virtual Machines in a Resource Group.
func testMultipleVMs(t *testing.T, terraformOptions *terraform.Options, subscriptionID string) {
	// Run `terraform output` to get the values of output variables.
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	expectedVMName := terraform.Output(t, terraformOptions, "vm_name")
	expectedVMSize := compute.VirtualMachineSizeTypes(terraform.Output(t, terraformOptions, "vm_size"))
	expectedAvsName := terraform.Output(t, terraformOptions, "availability_set_name")

	// Check against all VM names in a Resource Group.
	vmList := azure.ListVirtualMachinesForResourceGroup(t, resourceGroupName, subscriptionID)
	expectedVMCount := 1
	assert.Equal(t, expectedVMCount, len(vmList))
	assert.Contains(t, vmList, expectedVMName)

	// Check Availability Set for multiple VMs.
	actualVMsInAvs := azure.GetAvailabilitySetVMNamesInCaps(t, expectedAvsName, resourceGroupName, subscriptionID)
	assert.Contains(t, actualVMsInAvs, strings.ToUpper(expectedVMName))

	// Get all VMs in a Resource Group, including their properties, therefore avoiding
	// multiple SDK calls. The penalty for this approach is introducing direct references
	// which need to be checked for nil for optional configurations.
	vmsByRef := azure.GetVirtualMachinesForResourceGroup(t, resourceGroupName, subscriptionID)
	thisVM := vmsByRef[expectedVMName]
	assert.Equal(t, expectedVMSize, thisVM.HardwareProfile.VMSize)

	// Check for the VM negative test.
	fakeVM := fmt.Sprintf("vm-%s", random.UniqueId())
	assert.Nil(t, vmsByRef[fakeVM].VMID)
}

// These tests check information directly related to the specified Azure Virtual Machine.
func testInformationOfVM(t *testing.T, terraformOptions *terraform.Options, subscriptionID string) {
	// Run `terraform output` to get the values of output variables.
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	virtualMachineName := terraform.Output(t, terraformOptions, "vm_name")
	expectedVmAdminUser := terraform.OutputList(t, terraformOptions, "vm_admin_username")
	expectedImageSKU := terraform.OutputList(t, terraformOptions, "vm_image_sku")
	expectedImageVersion := terraform.OutputList(t, terraformOptions, "vm_image_version")
	expectedAvsName := terraform.Output(t, terraformOptions, "availability_set_name")
	expectedVMTags := terraform.OutputMap(t, terraformOptions, "vm_tags")

	// Check if the Virtual Machine exists.
	assert.True(t, azure.VirtualMachineExists(t, virtualMachineName, resourceGroupName, subscriptionID))

	// Check the Admin User of the VM.
	actualVM := azure.GetVirtualMachine(t, virtualMachineName, resourceGroupName, subscriptionID)
	actualVmAdminUser := *actualVM.OsProfile.AdminUsername
	assert.Equal(t, expectedVmAdminUser[0], actualVmAdminUser)

	// Check the Storage Image properties of the VM.
	actualImage := azure.GetVirtualMachineImage(t, virtualMachineName, resourceGroupName, subscriptionID)
	assert.Contains(t, expectedImageSKU[0], actualImage.SKU)
	assert.Contains(t, expectedImageVersion[0], actualImage.Version)

	// Check the Availability Set of the VM.
	// The AVS ID returned from the VM is always CAPS so ignoring case in the assertion.
	actualexpectedAvsName := azure.GetVirtualMachineAvailabilitySetID(t, virtualMachineName, resourceGroupName, subscriptionID)
	assert.True(t, strings.EqualFold(expectedAvsName, actualexpectedAvsName))

	// Check the assigned Tags of the VM, assert empty if no tags.
	actualVMTags := azure.GetVirtualMachineTags(t, virtualMachineName, resourceGroupName, "")
	assert.Equal(t, expectedVMTags, actualVMTags)
}

// These tests check the OS Disk and Attached Managed Disks for the Azure Virtual Machine.
// The following Terratest Azure module is utilized in addition to the compute module:
// - disk
// See the terraform_azure_disk_example_test.go for other related tests.
func testDisksOfVM(t *testing.T, terraformOptions *terraform.Options, subscriptionID string) {
	// Run `terraform output` to get the values of output variables.
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	virtualMachineName := terraform.Output(t, terraformOptions, "vm_name")
	expectedOSDiskName := terraform.Output(t, terraformOptions, "os_disk_name")
	expectedDiskName := terraform.Output(t, terraformOptions, "managed_disk_name")
	expectedDiskType := terraform.Output(t, terraformOptions, "managed_disk_type")

	// Check the OS Disk name of the VM.
	actualOSDiskName := azure.GetVirtualMachineOSDiskName(t, virtualMachineName, resourceGroupName, subscriptionID)
	assert.Equal(t, expectedOSDiskName, actualOSDiskName)

	// Check the VM Managed Disk exists in the list of all VM Managed Disks.
	actualManagedDiskNames := azure.GetVirtualMachineManagedDisks(t, virtualMachineName, resourceGroupName, subscriptionID)
	assert.Contains(t, actualManagedDiskNames, expectedDiskName)

	// Check the Managed Disk count of the VM.
	expectedManagedDiskCount := 1
	assert.Equal(t, expectedManagedDiskCount, len(actualManagedDiskNames))

	// Check the Disk Type of the Managed Disk of the VM.
	// This does not apply to VHD disks saved under a storage account.
	actualDisk := azure.GetDisk(t, expectedDiskName, resourceGroupName, subscriptionID)
	actualDiskType := actualDisk.Sku.Name
	assert.Equal(t, compute.DiskStorageAccountTypes(expectedDiskType), actualDiskType)
}

// These tests check the underlying Virtual Network, Network Interface and associated Public IP Address.
// The following Terratest Azure modules are utilized in addition to the compute module:
// - networkinterface
// - publicaddress
// - virtualnetwork
// See the terraform_azure_network_example_test.go for other related tests.
func testNetworkOfVM(t *testing.T, terraformOptions *terraform.Options, subscriptionID string) {
	// Run `terraform output` to get the values of output variables.
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	virtualMachineName := terraform.Output(t, terraformOptions, "vm_name")
	expectedVNetName := terraform.Output(t, terraformOptions, "virtual_network_name")
	expectedSubnetName := terraform.Output(t, terraformOptions, "subnet_name")
	expectedPublicAddressName := terraform.Output(t, terraformOptions, "public_ip_name")
	expectedNicName := terraform.Output(t, terraformOptions, "network_interface_name")
	expectedPrivateIPAddress := terraform.Output(t, terraformOptions, "private_ip")

	// VirtualNetwork and Subnet tests
	// Check the Subnet exists in the Virtual Network.
	actualVnetSubnets := azure.GetVirtualNetworkSubnets(t, expectedVNetName, resourceGroupName, subscriptionID)
	assert.NotNil(t, actualVnetSubnets[expectedVNetName])

	// Check the Private IP is in the Subnet Range.
	actualVMNicIPInSubnet := azure.CheckSubnetContainsIP(t, expectedPrivateIPAddress, expectedSubnetName, expectedVNetName, resourceGroupName, subscriptionID)
	assert.True(t, actualVMNicIPInSubnet)

	// Network Interface Card tests
	// Check the VM Network Interface exists in the list of all VM Network Interfaces.
	actualNics := azure.GetVirtualMachineNics(t, virtualMachineName, resourceGroupName, subscriptionID)
	assert.Contains(t, actualNics, expectedNicName)

	// Check the Network Interface count of the VM.
	expectedNICCount := 1
	assert.Equal(t, expectedNICCount, len(actualNics))

	// Check for the Private IP in the NICs IP list.
	actualPrivateIPAddress := azure.GetNetworkInterfacePrivateIPs(t, expectedNicName, resourceGroupName, subscriptionID)
	assert.Contains(t, actualPrivateIPAddress, expectedPrivateIPAddress)

	// Public IP Address test
	// Check for the Public IP for the NIC. No expected value since it is assigned runtime.
	actualPublicIP := azure.GetIPOfPublicIPAddressByName(t, expectedPublicAddressName, resourceGroupName, subscriptionID)
	assert.NotNil(t, actualPublicIP)

}
