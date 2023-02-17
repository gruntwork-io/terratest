package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetVirtualMachineClient is a helper function that will setup an Azure Virtual Machine client on your behalf.
func GetVirtualMachineClient(t testing.TestingT, subscriptionID string) *compute.VirtualMachinesClient {
	vmClient, err := GetVirtualMachineClientE(subscriptionID)
	require.NoError(t, err)
	return vmClient
}

// GetVirtualMachineClientE is a helper function that will setup an Azure Virtual Machine client on your behalf.
func GetVirtualMachineClientE(subscriptionID string) (*compute.VirtualMachinesClient, error) {

	// snippet-tag-start::client_factory_example.helper
	// Create a VM client
	vmClient, err := CreateVirtualMachinesClientE(subscriptionID)
	if err != nil {
		return nil, err
	}
	// snippet-tag-end::client_factory_example.helper

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	vmClient.Authorizer = *authorizer
	return vmClient, nil
}

// VirtualMachineExists indicates whether the specifcied Azure Virtual Machine exists.
// This function would fail the test if there is an error.
func VirtualMachineExists(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) bool {
	exists, err := VirtualMachineExistsE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// VirtualMachineExistsE indicates whether the specifcied Azure Virtual Machine exists.
func VirtualMachineExistsE(vmName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get VM Object
	_, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetVirtualMachineNics gets a list of Network Interface names for a specifcied Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineNics(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) []string {
	nicList, err := GetVirtualMachineNicsE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return nicList
}

// GetVirtualMachineNicsE gets a list of Network Interface names for a specified Azure Virtual Machine.
func GetVirtualMachineNicsE(vmName string, resGroupName string, subscriptionID string) ([]string, error) {

	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get VM NIC(s); value always present, no nil checks needed.
	vmNICs := *vm.NetworkProfile.NetworkInterfaces

	nics := make([]string, len(vmNICs))
	for i, nic := range vmNICs {
		// Get ID from resource string.
		nicName, err := GetNameFromResourceIDE(*nic.ID)
		if err == nil {
			nics[i] = nicName
		}
	}
	return nics, nil
}

// GetVirtualMachineManagedDisks gets the list of Managed Disk names of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineManagedDisks(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) []string {
	diskNames, err := GetVirtualMachineManagedDisksE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return diskNames
}

// GetVirtualMachineManagedDisksE gets the list of Managed Disk names of the specified Azure Virtual Machine.
func GetVirtualMachineManagedDisksE(vmName string, resGroupName string, subscriptionID string) ([]string, error) {

	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get VM attached Disks; value always present even if no disks attached, no nil check needed.
	vmDisks := *vm.StorageProfile.DataDisks

	// Get the Names of the attached Managed Disks
	diskNames := make([]string, len(vmDisks))
	for i, v := range vmDisks {
		// Disk names are required, no nil check needed.
		diskNames[i] = *v.Name
	}

	return diskNames, nil
}

// GetVirtualMachineOSDiskName gets the OS Disk name of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineOSDiskName(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	osDiskName, err := GetVirtualMachineOSDiskNameE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return osDiskName
}

// GetVirtualMachineOSDiskNameE gets the OS Disk name of the specified Azure Virtual Machine.
func GetVirtualMachineOSDiskNameE(vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return *vm.StorageProfile.OsDisk.Name, nil
}

// GetVirtualMachineAvailabilitySetID gets the Availability Set ID of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineAvailabilitySetID(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	avsID, err := GetVirtualMachineAvailabilitySetIDE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return avsID
}

// GetVirtualMachineAvailabilitySetIDE gets the Availability Set ID of the specified Azure Virtual Machine.
func GetVirtualMachineAvailabilitySetIDE(vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	// Virtual Machine has no associated Availability Set
	if vm.AvailabilitySet == nil {
		return "", nil
	}

	// Get ID from resource string
	avs, err := GetNameFromResourceIDE(*vm.AvailabilitySet.ID)
	if err != nil {
		return "", err
	}

	return avs, nil
}

// VMImage represents the storage image for the specified Azure Virtual Machine.
type VMImage struct {
	Publisher string
	Offer     string
	SKU       string
	Version   string
}

// GetVirtualMachineImage gets the Image of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetVirtualMachineImage(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) VMImage {
	vmImage, err := GetVirtualMachineImageE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vmImage
}

// GetVirtualMachineImageE gets the Image  of the specified Azure Virtual Machine.
func GetVirtualMachineImageE(vmName string, resGroupName string, subscriptionID string) (VMImage, error) {
	var vmImage VMImage

	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return vmImage, err
	}

	// Populate VM Image; values always present, no nil checks needed
	vmImage.Publisher = *vm.StorageProfile.ImageReference.Publisher
	vmImage.Offer = *vm.StorageProfile.ImageReference.Offer
	vmImage.SKU = *vm.StorageProfile.ImageReference.Sku
	vmImage.Version = *vm.StorageProfile.ImageReference.Version

	return vmImage, nil
}

// GetSizeOfVirtualMachine gets the Size Type of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
func GetSizeOfVirtualMachine(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) compute.VirtualMachineSizeTypes {
	size, err := GetSizeOfVirtualMachineE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return size
}

// GetSizeOfVirtualMachineE gets the Size Type of the specified Azure Virtual Machine.
func GetSizeOfVirtualMachineE(vmName string, resGroupName string, subscriptionID string) (compute.VirtualMachineSizeTypes, error) {
	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return vm.VirtualMachineProperties.HardwareProfile.VMSize, nil
}

// GetVirtualMachineTags gets the Tags of the specified Virtual Machine as a map.
// This function would fail the test if there is an error.
func GetVirtualMachineTags(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) map[string]string {
	tags, err := GetVirtualMachineTagsE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return tags
}

// GetVirtualMachineTagsE gets the Tags of the specified Virtual Machine as a map.
func GetVirtualMachineTagsE(vmName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	// Setup a blank map to populate and return
	tags := make(map[string]string)

	// Get VM Object
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	if err != nil {
		return tags, err
	}

	// Range through existing tags and populate above map accordingly
	for k, v := range vm.Tags {
		tags[k] = *v
	}

	return tags, nil
}

// ***************************************************** //
// Get multiple Virtual Machines from a Resource Group
// ***************************************************** //

// ListVirtualMachinesForResourceGroup gets a list of all Virtual Machine names in the specified Resource Group.
// This function would fail the test if there is an error.
func ListVirtualMachinesForResourceGroup(t testing.TestingT, resGroupName string, subscriptionID string) []string {
	vms, err := ListVirtualMachinesForResourceGroupE(resGroupName, subscriptionID)
	require.NoError(t, err)
	return vms
}

// ListVirtualMachinesForResourceGroupE gets a list of all Virtual Machine names in the specified Resource Group.
func ListVirtualMachinesForResourceGroupE(resourceGroupName string, subscriptionID string) ([]string, error) {
	var vmDetails []string

	vmClient, err := GetVirtualMachineClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	vms, err := vmClient.List(context.Background(), resourceGroupName)
	if err != nil {
		return nil, err
	}

	for _, v := range vms.Values() {
		vmDetails = append(vmDetails, *v.Name)
	}
	return vmDetails, nil
}

// GetVirtualMachinesForResourceGroup gets all Virtual Machine objects in the specified Resource Group. Each
// VM Object represents the entire set of VM compute properties accessible by using the VM name as the map key.
// This function would fail the test if there is an error.
func GetVirtualMachinesForResourceGroup(t testing.TestingT, resGroupName string, subscriptionID string) map[string]compute.VirtualMachineProperties {
	vms, err := GetVirtualMachinesForResourceGroupE(resGroupName, subscriptionID)
	require.NoError(t, err)
	return vms
}

// GetVirtualMachinesForResourceGroupE gets all Virtual Machine objects in the specified Resource Group. Each
// VM Object represents the entire set of VM compute properties accessible by using the VM name as the map key.
func GetVirtualMachinesForResourceGroupE(resourceGroupName string, subscriptionID string) (map[string]compute.VirtualMachineProperties, error) {
	// Create VM Client
	vmClient, err := GetVirtualMachineClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the list of VMs in the Resource Group
	vms, err := vmClient.List(context.Background(), resourceGroupName)
	if err != nil {
		return nil, err
	}

	// Get the VMs in the Resource Group.
	vmDetails := make(map[string]compute.VirtualMachineProperties, len(vms.Values()))
	for _, v := range vms.Values() {
		// VM name and machine properties are required for each VM, no nill check required.
		vmDetails[*v.Name] = *v.VirtualMachineProperties
	}
	return vmDetails, nil
}

// ******************************************************************** //
// Get VM using Instance and Instance property get, reducing SKD calls
// ******************************************************************** //

// Instance of the VM
type Instance struct {
	*compute.VirtualMachine
}

// GetVirtualMachineInstanceSize gets the size of the Virtual Machine.
func (vm *Instance) GetVirtualMachineInstanceSize() compute.VirtualMachineSizeTypes {
	return vm.VirtualMachineProperties.HardwareProfile.VMSize
}

// *********************** //
// Get the base VM Object
// *********************** //

// GetVirtualMachine gets a Virtual Machine in the specified Azure Resource Group.
// This function would fail the test if there is an error.
func GetVirtualMachine(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) *compute.VirtualMachine {
	vm, err := GetVirtualMachineE(vmName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return vm
}

// GetVirtualMachineE gets a Virtual Machine in the specified Azure Resource Group.
func GetVirtualMachineE(vmName string, resGroupName string, subscriptionID string) (*compute.VirtualMachine, error) {
	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetVirtualMachineClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	vm, err := client.Get(context.Background(), resGroupName, vmName, compute.InstanceView)
	if err != nil {
		return nil, err
	}

	return &vm, nil
}
