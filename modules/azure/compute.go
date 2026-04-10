package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetVirtualMachineClient is a helper function that will setup an Azure Virtual Machine client on your behalf.
// This function would fail the test if there is an error.
func GetVirtualMachineClient(t testing.TestingT, subscriptionID string) *armcompute.VirtualMachinesClient {
	t.Helper()

	vmClient, err := GetVirtualMachineClientE(subscriptionID)
	require.NoError(t, err)

	return vmClient
}

// GetVirtualMachineClientE is a helper function that will setup an Azure Virtual Machine client on your behalf.
func GetVirtualMachineClientE(subscriptionID string) (*armcompute.VirtualMachinesClient, error) {
	// snippet-tag-start::client_factory_example.helper
	vmClient, err := CreateVirtualMachinesClientE(subscriptionID)
	if err != nil {
		return nil, err
	}
	// snippet-tag-end::client_factory_example.helper

	return vmClient, nil
}

// VirtualMachineExistsContext indicates whether the specified Azure Virtual Machine exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func VirtualMachineExistsContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := VirtualMachineExistsContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// VirtualMachineExistsContextE indicates whether the specified Azure Virtual Machine exists.
// The ctx parameter supports cancellation and timeouts.
func VirtualMachineExistsContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (bool, error) {
	_, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetVirtualMachineNicsContext gets a list of Network Interface names for a specified Azure Virtual Machine.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineNicsContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	nicList, err := GetVirtualMachineNicsContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return nicList
}

// GetVirtualMachineNicsContextE gets a list of Network Interface names for a specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineNicsContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) ([]string, error) {
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	return extractVMNics(vm)
}

// extractVMNics extracts the Network Interface names from a Virtual Machine object.
func extractVMNics(vm *armcompute.VirtualMachine) ([]string, error) {
	if vm.Properties == nil || vm.Properties.NetworkProfile == nil {
		return nil, nil
	}

	vmNICs := vm.Properties.NetworkProfile.NetworkInterfaces

	var nics []string

	for _, nic := range vmNICs {
		nicName, err := GetNameFromResourceIDE(*nic.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse NIC resource ID %q: %w", *nic.ID, err)
		}

		nics = append(nics, nicName)
	}

	return nics, nil
}

// GetVirtualMachineManagedDisksContext gets the list of Managed Disk names of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineManagedDisksContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	diskNames, err := GetVirtualMachineManagedDisksContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return diskNames
}

// GetVirtualMachineManagedDisksContextE gets the list of Managed Disk names of the specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineManagedDisksContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) ([]string, error) {
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	return extractVMManagedDisks(vm), nil
}

// extractVMManagedDisks extracts the Managed Disk names from a Virtual Machine object.
func extractVMManagedDisks(vm *armcompute.VirtualMachine) []string {
	vmDisks := vm.Properties.StorageProfile.DataDisks

	diskNames := make([]string, len(vmDisks))
	for i, v := range vmDisks {
		diskNames[i] = *v.Name
	}

	return diskNames
}

// GetVirtualMachineOSDiskNameContext gets the OS Disk name of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineOSDiskNameContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) string {
	t.Helper()

	osDiskName, err := GetVirtualMachineOSDiskNameContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return osDiskName
}

// GetVirtualMachineOSDiskNameContextE gets the OS Disk name of the specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineOSDiskNameContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (string, error) {
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return extractVMOSDiskName(vm), nil
}

// extractVMOSDiskName extracts the OS Disk name from a Virtual Machine object.
func extractVMOSDiskName(vm *armcompute.VirtualMachine) string {
	return *vm.Properties.StorageProfile.OSDisk.Name
}

// GetVirtualMachineAvailabilitySetIDContext gets the Availability Set ID of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineAvailabilitySetIDContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) string {
	t.Helper()

	avsID, err := GetVirtualMachineAvailabilitySetIDContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return avsID
}

// GetVirtualMachineAvailabilitySetIDContextE gets the Availability Set ID of the specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineAvailabilitySetIDContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (string, error) {
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return extractVMAvailabilitySetID(vm)
}

// extractVMAvailabilitySetID extracts the Availability Set ID from a Virtual Machine object.
func extractVMAvailabilitySetID(vm *armcompute.VirtualMachine) (string, error) {
	if vm.Properties.AvailabilitySet == nil {
		return "", nil
	}

	avs, err := GetNameFromResourceIDE(*vm.Properties.AvailabilitySet.ID)
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

// GetVirtualMachineImageContext gets the Image of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineImageContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) VMImage {
	t.Helper()

	vmImage, err := GetVirtualMachineImageContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vmImage
}

// GetVirtualMachineImageContextE gets the Image of the specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineImageContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (VMImage, error) {
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		return VMImage{}, err
	}

	return extractVMImage(vm), nil
}

// extractVMImage extracts the Image reference from a Virtual Machine object.
// For custom images where Publisher/Offer/SKU/Version may be nil, empty strings are returned.
func extractVMImage(vm *armcompute.VirtualMachine) VMImage {
	ref := vm.Properties.StorageProfile.ImageReference

	var img VMImage

	if ref.Publisher != nil {
		img.Publisher = *ref.Publisher
	}

	if ref.Offer != nil {
		img.Offer = *ref.Offer
	}

	if ref.SKU != nil {
		img.SKU = *ref.SKU
	}

	if ref.Version != nil {
		img.Version = *ref.Version
	}

	return img
}

// GetSizeOfVirtualMachineContext gets the Size Type of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSizeOfVirtualMachineContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) armcompute.VirtualMachineSizeTypes {
	t.Helper()

	size, err := GetSizeOfVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return size
}

// GetSizeOfVirtualMachineContextE gets the Size Type of the specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetSizeOfVirtualMachineContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (armcompute.VirtualMachineSizeTypes, error) {
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return extractVMSize(vm), nil
}

// extractVMSize extracts the VM size from a Virtual Machine object.
func extractVMSize(vm *armcompute.VirtualMachine) armcompute.VirtualMachineSizeTypes {
	return *vm.Properties.HardwareProfile.VMSize
}

// GetVirtualMachineTagsContext gets the Tags of the specified Virtual Machine as a map.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineTagsContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) map[string]string {
	t.Helper()

	tags, err := GetVirtualMachineTagsContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return tags
}

// GetVirtualMachineTagsContextE gets the Tags of the specified Virtual Machine as a map.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineTagsContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		return make(map[string]string), err
	}

	return extractVMTags(vm), nil
}

// extractVMTags extracts the Tags from a Virtual Machine object as a map of string to string.
func extractVMTags(vm *armcompute.VirtualMachine) map[string]string {
	tags := make(map[string]string)

	if vm.Tags == nil {
		return tags
	}

	for k, v := range vm.Tags {
		tags[k] = *v
	}

	return tags
}

// ***************************************************** //
// Get multiple Virtual Machines from a Resource Group
// ***************************************************** //

// ListVirtualMachinesForResourceGroupContext gets a list of all Virtual Machine names in the specified Resource Group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ListVirtualMachinesForResourceGroupContext(t testing.TestingT, ctx context.Context, resGroupName string, subscriptionID string) []string {
	t.Helper()

	vms, err := ListVirtualMachinesForResourceGroupContextE(ctx, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vms
}

// ListVirtualMachinesForResourceGroupContextE gets a list of all Virtual Machine names in the specified Resource Group.
// The ctx parameter supports cancellation and timeouts.
func ListVirtualMachinesForResourceGroupContextE(ctx context.Context, resourceGroupName string, subscriptionID string) ([]string, error) {
	vmClient, err := GetVirtualMachineClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	return listVirtualMachineNames(ctx, vmClient, resourceGroupName)
}

// listVirtualMachineNames pages through VMs in a resource group and returns their names.
func listVirtualMachineNames(ctx context.Context, client *armcompute.VirtualMachinesClient, resourceGroupName string) ([]string, error) {
	var vmDetails []string

	pager := client.NewListPager(resourceGroupName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Value {
			vmDetails = append(vmDetails, *v.Name)
		}
	}

	return vmDetails, nil
}

// GetVirtualMachinesForResourceGroupContext gets all Virtual Machine objects in the specified Resource Group. Each
// VM Object represents the entire set of VM compute properties accessible by using the VM name as the map key.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachinesForResourceGroupContext(t testing.TestingT, ctx context.Context, resGroupName string, subscriptionID string) map[string]armcompute.VirtualMachineProperties {
	t.Helper()

	vms, err := GetVirtualMachinesForResourceGroupContextE(ctx, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vms
}

// GetVirtualMachinesForResourceGroupContextE gets all Virtual Machine objects in the specified Resource Group. Each
// VM Object represents the entire set of VM compute properties accessible by using the VM name as the map key.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachinesForResourceGroupContextE(ctx context.Context, resourceGroupName string, subscriptionID string) (map[string]armcompute.VirtualMachineProperties, error) {
	vmClient, err := GetVirtualMachineClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	return listVirtualMachineProperties(ctx, vmClient, resourceGroupName)
}

// listVirtualMachineProperties pages through VMs in a resource group and returns a map of VM name to properties.
func listVirtualMachineProperties(ctx context.Context, client *armcompute.VirtualMachinesClient, resourceGroupName string) (map[string]armcompute.VirtualMachineProperties, error) {
	vmDetails := make(map[string]armcompute.VirtualMachineProperties)

	pager := client.NewListPager(resourceGroupName, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Value {
			vmDetails[*v.Name] = *v.Properties
		}
	}

	return vmDetails, nil
}

// ******************************************************************** //
// Get VM using Instance and Instance property get, reducing SKD calls
// ******************************************************************** //

// Instance of the VM
type Instance struct {
	*armcompute.VirtualMachine
}

// GetVirtualMachineInstanceSize gets the size of the Virtual Machine.
func (vm *Instance) GetVirtualMachineInstanceSize() armcompute.VirtualMachineSizeTypes {
	return *vm.Properties.HardwareProfile.VMSize
}

// *********************** //
// Get the base VM Object
// *********************** //

// GetVirtualMachineContext gets a Virtual Machine in the specified Azure Resource Group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) *armcompute.VirtualMachine {
	t.Helper()

	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vm
}

// GetVirtualMachineContextE gets a Virtual Machine in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (*armcompute.VirtualMachine, error) {
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetVirtualMachineClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	return fetchVirtualMachine(ctx, client, resGroupName, vmName)
}

// fetchVirtualMachine retrieves a single Virtual Machine from Azure using the provided client.
func fetchVirtualMachine(ctx context.Context, client *armcompute.VirtualMachinesClient, resGroupName, vmName string) (*armcompute.VirtualMachine, error) {
	resp, err := client.Get(ctx, resGroupName, vmName, &armcompute.VirtualMachinesClientGetOptions{
		Expand: to.Ptr(armcompute.InstanceViewTypesInstanceView),
	})
	if err != nil {
		return nil, err
	}

	return &resp.VirtualMachine, nil
}
