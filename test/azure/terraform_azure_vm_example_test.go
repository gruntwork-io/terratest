//go:build azure || (azureslim && compute)
// +build azure azureslim,compute

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test_test

import (
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureVmExample(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	uniquePostfix := random.UniqueID()

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
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// Run tests for the Virtual Machine.
	testStrategiesForVMs(t, terraformOptions, subscriptionID)
	testMultipleVMs(t, terraformOptions, subscriptionID)
	testInformationOfVM(t, terraformOptions, subscriptionID)
	testDisksOfVM(t, terraformOptions, subscriptionID)
	testNetworkOfVM(t, terraformOptions, subscriptionID)
}

// These 3 tests check for the same property but illustrate different testing strategies for
// retrieving the data. The first strategy is used in the other tests of this module while
// the other two can be extended by the user as needed.
func testStrategiesForVMs(t *testing.T, terraformOptions *terraform.Options, subscriptionID string) {
	t.Helper()

	// Run `terraform output` to get the values of output variables.
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	virtualMachineName := terraform.OutputContext(t, t.Context(), terraformOptions, "vm_name")
	expectedVMSize := armcompute.VirtualMachineSizeTypes(terraform.OutputContext(t, t.Context(), terraformOptions, "vm_size"))

	// 1. Check the VM Size directly. This strategy gets one specific property of the VM per method.
	actualVMSize := azure.GetSizeOfVirtualMachineContext(t, t.Context(), virtualMachineName, resourceGroupName, subscriptionID)
	assert.Equal(t, expectedVMSize, actualVMSize)

	// 2. Check the VM size by reference. This strategy is beneficial when checking multiple properties
	// by using one VM reference. Optional parameters have to be checked first to avoid nil panics.
	vmByRef := azure.GetVirtualMachineContext(t, t.Context(), virtualMachineName, resourceGroupName, subscriptionID)
	actualVMSize = *vmByRef.Properties.HardwareProfile.VMSize
	assert.Equal(t, expectedVMSize, actualVMSize)

	// 3. Check the VM size by instance. This strategy is beneficial when checking multiple properties
	// by using one VM instance and making calls against it with the added benefit of property check abstraction.
	vmInstance := azure.Instance{VirtualMachine: vmByRef}
	actualVMSize = vmInstance.GetVirtualMachineInstanceSize()
	assert.Equal(t, expectedVMSize, actualVMSize)
}

// These tests check for the multiple Virtual Machines in a Resource Group.
func testMultipleVMs(t *testing.T, terraformOptions *terraform.Options, subscriptionID string) {
	t.Helper()

	// Run `terraform output` to get the values of output variables.
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	expectedVMName := terraform.OutputContext(t, t.Context(), terraformOptions, "vm_name")
	expectedVMSize := armcompute.VirtualMachineSizeTypes(terraform.OutputContext(t, t.Context(), terraformOptions, "vm_size"))
	expectedAvsName := terraform.OutputContext(t, t.Context(), terraformOptions, "availability_set_name")

	// Check against all VM names in a Resource Group.
	vmList := azure.ListVirtualMachinesForResourceGroupContext(t, t.Context(), resourceGroupName, subscriptionID)
	expectedVMCount := 1
	assert.Len(t, vmList, expectedVMCount)
	assert.Contains(t, vmList, expectedVMName)

	// Check Availability Set for multiple VMs.
	actualVMsInAvs := azure.GetAvailabilitySetVMNamesInCapsContext(t, t.Context(), expectedAvsName, resourceGroupName, subscriptionID)
	assert.Contains(t, actualVMsInAvs, strings.ToUpper(expectedVMName))

	// Get all VMs in a Resource Group, including their properties, therefore avoiding
	// multiple SDK calls. The penalty for this approach is introducing direct references
	// which need to be checked for nil for optional configurations.
	vmsByRef := azure.GetVirtualMachinesForResourceGroupContext(t, t.Context(), resourceGroupName, subscriptionID)
	thisVM := vmsByRef[expectedVMName]
	assert.Equal(t, expectedVMSize, *thisVM.HardwareProfile.VMSize)

	// Check for the VM negative test.
	fakeVM := "vm-" + random.UniqueID()
	assert.Nil(t, vmsByRef[fakeVM].VMID)
}

// These tests check information directly related to the specified Azure Virtual Machine.
func testInformationOfVM(t *testing.T, terraformOptions *terraform.Options, subscriptionID string) {
	t.Helper()

	// Run `terraform output` to get the values of output variables.
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	virtualMachineName := terraform.OutputContext(t, t.Context(), terraformOptions, "vm_name")
	expectedVMAdminUser := terraform.OutputListContext(t, t.Context(), terraformOptions, "vm_admin_username")
	expectedImageSKU := terraform.OutputListContext(t, t.Context(), terraformOptions, "vm_image_sku")
	expectedImageVersion := terraform.OutputListContext(t, t.Context(), terraformOptions, "vm_image_version")
	expectedAvsName := terraform.OutputContext(t, t.Context(), terraformOptions, "availability_set_name")
	expectedVMTags := terraform.OutputMapContext(t, t.Context(), terraformOptions, "vm_tags")

	// Check if the Virtual Machine exists.
	assert.True(t, azure.VirtualMachineExistsContext(t, t.Context(), virtualMachineName, resourceGroupName, subscriptionID))

	// Check the Admin User of the VM.
	actualVM := azure.GetVirtualMachineContext(t, t.Context(), virtualMachineName, resourceGroupName, subscriptionID)
	actualVMAdminUser := *actualVM.Properties.OSProfile.AdminUsername
	assert.Equal(t, expectedVMAdminUser[0], actualVMAdminUser)

	// Check the Storage Image properties of the VM.
	actualImage := azure.GetVirtualMachineImageContext(t, t.Context(), virtualMachineName, resourceGroupName, subscriptionID)
	assert.Contains(t, expectedImageSKU[0], actualImage.SKU)
	assert.Contains(t, expectedImageVersion[0], actualImage.Version)

	// Check the Availability Set of the VM.
	// The AVS ID returned from the VM is always CAPS so ignoring case in the assertion.
	actualexpectedAvsName := azure.GetVirtualMachineAvailabilitySetIDContext(t, t.Context(), virtualMachineName, resourceGroupName, subscriptionID)
	assert.True(t, strings.EqualFold(expectedAvsName, actualexpectedAvsName))

	// Check the assigned Tags of the VM, assert empty if no tags.
	actualVMTags := azure.GetVirtualMachineTagsContext(t, t.Context(), virtualMachineName, resourceGroupName, "")
	assert.Equal(t, expectedVMTags, actualVMTags)
}

// These tests check the OS Disk and Attached Managed Disks for the Azure Virtual Machine.
// The following Terratest Azure module is utilized in addition to the compute module:
// - disk
// See the terraform_azure_disk_example_test.go for other related tests.
func testDisksOfVM(t *testing.T, terraformOptions *terraform.Options, subscriptionID string) {
	t.Helper()

	// Run `terraform output` to get the values of output variables.
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	virtualMachineName := terraform.OutputContext(t, t.Context(), terraformOptions, "vm_name")
	expectedOSDiskName := terraform.OutputContext(t, t.Context(), terraformOptions, "os_disk_name")
	expectedDiskName := terraform.OutputContext(t, t.Context(), terraformOptions, "managed_disk_name")
	expectedDiskType := terraform.OutputContext(t, t.Context(), terraformOptions, "managed_disk_type")

	// Check the OS Disk name of the VM.
	actualOSDiskName := azure.GetVirtualMachineOSDiskNameContext(t, t.Context(), virtualMachineName, resourceGroupName, subscriptionID)
	assert.Equal(t, expectedOSDiskName, actualOSDiskName)

	// Check the VM Managed Disk exists in the list of all VM Managed Disks.
	actualManagedDiskNames := azure.GetVirtualMachineManagedDisksContext(t, t.Context(), virtualMachineName, resourceGroupName, subscriptionID)
	assert.Contains(t, actualManagedDiskNames, expectedDiskName)

	// Check the Managed Disk count of the VM.
	expectedManagedDiskCount := 1
	assert.Len(t, actualManagedDiskNames, expectedManagedDiskCount)

	// Check the Disk Type of the Managed Disk of the VM.
	// This does not apply to VHD disks saved under a storage account.
	actualDisk := azure.GetDiskContext(t, t.Context(), expectedDiskName, resourceGroupName, subscriptionID)
	actualDiskType := actualDisk.SKU.Name
	assert.Equal(t, armcompute.DiskStorageAccountTypes(expectedDiskType), *actualDiskType)
}

// These tests check the underlying Virtual Network, Network Interface and associated Public IP Address.
// The following Terratest Azure modules are utilized in addition to the compute module:
// - networkinterface
// - publicaddress
// - virtualnetwork
// See the terraform_azure_network_example_test.go for other related tests.
func testNetworkOfVM(t *testing.T, terraformOptions *terraform.Options, subscriptionID string) {
	t.Helper()

	// Run `terraform output` to get the values of output variables.
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	virtualMachineName := terraform.OutputContext(t, t.Context(), terraformOptions, "vm_name")
	expectedVNetName := terraform.OutputContext(t, t.Context(), terraformOptions, "virtual_network_name")
	expectedSubnetName := terraform.OutputContext(t, t.Context(), terraformOptions, "subnet_name")
	expectedPublicAddressName := terraform.OutputContext(t, t.Context(), terraformOptions, "public_ip_name")
	expectedNicName := terraform.OutputContext(t, t.Context(), terraformOptions, "network_interface_name")
	expectedPrivateIPAddress := terraform.OutputContext(t, t.Context(), terraformOptions, "private_ip")

	// VirtualNetwork and Subnet tests
	// Check the Subnet exists in the Virtual Network.
	actualVnetSubnets := azure.GetVirtualNetworkSubnetsContext(t, t.Context(), expectedVNetName, resourceGroupName, subscriptionID)
	assert.NotNil(t, actualVnetSubnets[expectedVNetName])

	// Check the Private IP is in the Subnet Range.
	actualVMNicIPInSubnet := azure.CheckSubnetContainsIPContext(t, t.Context(), expectedPrivateIPAddress, expectedSubnetName, expectedVNetName, resourceGroupName, subscriptionID)
	assert.True(t, actualVMNicIPInSubnet)

	// Network Interface Card tests
	// Check the VM Network Interface exists in the list of all VM Network Interfaces.
	actualNics := azure.GetVirtualMachineNicsContext(t, t.Context(), virtualMachineName, resourceGroupName, subscriptionID)
	assert.Contains(t, actualNics, expectedNicName)

	// Check the Network Interface count of the VM.
	expectedNICCount := 1
	assert.Len(t, actualNics, expectedNICCount)

	// Check for the Private IP in the NICs IP list.
	actualPrivateIPAddress := azure.GetNetworkInterfacePrivateIPsContext(t, t.Context(), expectedNicName, resourceGroupName, subscriptionID)
	assert.Contains(t, actualPrivateIPAddress, expectedPrivateIPAddress)

	// Public IP Address test
	// Check for the Public IP for the NIC. No expected value since it is assigned runtime.
	actualPublicIP := azure.GetIPOfPublicIPAddressByNameContext(t, t.Context(), expectedPublicAddressName, resourceGroupName, subscriptionID)
	assert.NotNil(t, actualPublicIP)
}
