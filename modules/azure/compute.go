package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetVirtualMachineClientContext is a helper function that will setup an Azure Virtual Machine client on your behalf.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineClientContext(t testing.TestingT, ctx context.Context, subscriptionID string) *armcompute.VirtualMachinesClient {
	t.Helper()

	vmClient, err := GetVirtualMachineClientContextE(ctx, subscriptionID)
	require.NoError(t, err)

	return vmClient
}

// GetVirtualMachineClient is a helper function that will setup an Azure Virtual Machine client on your behalf.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineClientContext] instead.
func GetVirtualMachineClient(t testing.TestingT, subscriptionID string) *armcompute.VirtualMachinesClient {
	t.Helper()

	return GetVirtualMachineClientContext(t, context.Background(), subscriptionID) //nolint:staticcheck
}

// GetVirtualMachineClientContextE is a helper function that will setup an Azure Virtual Machine client on your behalf.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineClientContextE(ctx context.Context, subscriptionID string) (*armcompute.VirtualMachinesClient, error) {
	vmClient, err := CreateVirtualMachinesClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return vmClient, nil
}

// GetVirtualMachineClientE is a helper function that will setup an Azure Virtual Machine client on your behalf.
//
// Deprecated: Use [GetVirtualMachineClientContextE] instead.
func GetVirtualMachineClientE(subscriptionID string) (*armcompute.VirtualMachinesClient, error) {
	return GetVirtualMachineClientContextE(context.Background(), subscriptionID)
}

// VirtualMachineExists indicates whether the specified Azure Virtual Machine exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [VirtualMachineExistsContext] instead.
func VirtualMachineExists(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	return VirtualMachineExistsContext(t, context.Background(), vmName, resGroupName, subscriptionID)
}

// VirtualMachineExistsE indicates whether the specified Azure Virtual Machine exists.
//
// Deprecated: Use [VirtualMachineExistsContextE] instead.
func VirtualMachineExistsE(vmName string, resGroupName string, subscriptionID string) (bool, error) {
	return VirtualMachineExistsContextE(context.Background(), vmName, resGroupName, subscriptionID)
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

// GetVirtualMachineNics gets a list of Network Interface names for a specified Azure Virtual Machine.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineNicsContext] instead.
func GetVirtualMachineNics(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	return GetVirtualMachineNicsContext(t, context.Background(), vmName, resGroupName, subscriptionID)
}

// GetVirtualMachineNicsE gets a list of Network Interface names for a specified Azure Virtual Machine.
//
// Deprecated: Use [GetVirtualMachineNicsContextE] instead.
func GetVirtualMachineNicsE(vmName string, resGroupName string, subscriptionID string) ([]string, error) {
	return GetVirtualMachineNicsContextE(context.Background(), vmName, resGroupName, subscriptionID)
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
	if vm == nil || vm.Properties == nil || vm.Properties.NetworkProfile == nil {
		return nil, nil
	}

	vmNICs := vm.Properties.NetworkProfile.NetworkInterfaces

	var nics []string

	for _, nic := range vmNICs {
		if nic == nil || nic.ID == nil {
			continue
		}

		nicName, err := GetNameFromResourceIDE(*nic.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse NIC resource ID %q: %w", *nic.ID, err)
		}

		nics = append(nics, nicName)
	}

	return nics, nil
}

// GetVirtualMachineManagedDisks gets the list of Managed Disk names of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineManagedDisksContext] instead.
func GetVirtualMachineManagedDisks(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	return GetVirtualMachineManagedDisksContext(t, context.Background(), vmName, resGroupName, subscriptionID)
}

// GetVirtualMachineManagedDisksE gets the list of Managed Disk names of the specified Azure Virtual Machine.
//
// Deprecated: Use [GetVirtualMachineManagedDisksContextE] instead.
func GetVirtualMachineManagedDisksE(vmName string, resGroupName string, subscriptionID string) ([]string, error) {
	return GetVirtualMachineManagedDisksContextE(context.Background(), vmName, resGroupName, subscriptionID)
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
	if vm == nil || vm.Properties == nil || vm.Properties.StorageProfile == nil {
		return nil
	}

	vmDisks := vm.Properties.StorageProfile.DataDisks

	diskNames := make([]string, 0, len(vmDisks))

	for _, v := range vmDisks {
		if v == nil || v.Name == nil {
			continue
		}

		diskNames = append(diskNames, *v.Name)
	}

	return diskNames
}

// GetVirtualMachineOSDiskName gets the OS Disk name of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineOSDiskNameContext] instead.
func GetVirtualMachineOSDiskName(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	t.Helper()

	return GetVirtualMachineOSDiskNameContext(t, context.Background(), vmName, resGroupName, subscriptionID)
}

// GetVirtualMachineOSDiskNameE gets the OS Disk name of the specified Azure Virtual Machine.
//
// Deprecated: Use [GetVirtualMachineOSDiskNameContextE] instead.
func GetVirtualMachineOSDiskNameE(vmName string, resGroupName string, subscriptionID string) (string, error) {
	return GetVirtualMachineOSDiskNameContextE(context.Background(), vmName, resGroupName, subscriptionID)
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
	if vm == nil || vm.Properties == nil || vm.Properties.StorageProfile == nil ||
		vm.Properties.StorageProfile.OSDisk == nil || vm.Properties.StorageProfile.OSDisk.Name == nil {
		return ""
	}

	return *vm.Properties.StorageProfile.OSDisk.Name
}

// GetVirtualMachineAvailabilitySetID gets the Availability Set ID of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineAvailabilitySetIDContext] instead.
func GetVirtualMachineAvailabilitySetID(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	t.Helper()

	return GetVirtualMachineAvailabilitySetIDContext(t, context.Background(), vmName, resGroupName, subscriptionID)
}

// GetVirtualMachineAvailabilitySetIDE gets the Availability Set ID of the specified Azure Virtual Machine.
//
// Deprecated: Use [GetVirtualMachineAvailabilitySetIDContextE] instead.
func GetVirtualMachineAvailabilitySetIDE(vmName string, resGroupName string, subscriptionID string) (string, error) {
	return GetVirtualMachineAvailabilitySetIDContextE(context.Background(), vmName, resGroupName, subscriptionID)
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
	if vm == nil || vm.Properties == nil || vm.Properties.AvailabilitySet == nil ||
		vm.Properties.AvailabilitySet.ID == nil {
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

// GetVirtualMachineImage gets the Image of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineImageContext] instead.
func GetVirtualMachineImage(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) *VMImage {
	t.Helper()

	return GetVirtualMachineImageContext(t, context.Background(), vmName, resGroupName, subscriptionID)
}

// GetVirtualMachineImageE gets the Image of the specified Azure Virtual Machine.
//
// Deprecated: Use [GetVirtualMachineImageContextE] instead.
func GetVirtualMachineImageE(vmName string, resGroupName string, subscriptionID string) (*VMImage, error) {
	return GetVirtualMachineImageContextE(context.Background(), vmName, resGroupName, subscriptionID)
}

// GetVirtualMachineImageContext gets the Image of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineImageContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) *VMImage {
	t.Helper()

	vmImage, err := GetVirtualMachineImageContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vmImage
}

// GetVirtualMachineImageContextE gets the Image of the specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineImageContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (*VMImage, error) {
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		return nil, err
	}

	return extractVMImage(vm), nil
}

// extractVMImage extracts the Image reference from a Virtual Machine object.
// For custom images where Publisher/Offer/SKU/Version may be nil, empty strings are returned.
func extractVMImage(vm *armcompute.VirtualMachine) *VMImage {
	img := &VMImage{}

	if vm == nil || vm.Properties == nil || vm.Properties.StorageProfile == nil ||
		vm.Properties.StorageProfile.ImageReference == nil {
		return img
	}

	ref := vm.Properties.StorageProfile.ImageReference

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

// GetSizeOfVirtualMachine gets the Size Type of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetSizeOfVirtualMachineContext] instead.
func GetSizeOfVirtualMachine(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) armcompute.VirtualMachineSizeTypes {
	t.Helper()

	return GetSizeOfVirtualMachineContext(t, context.Background(), vmName, resGroupName, subscriptionID)
}

// GetSizeOfVirtualMachineE gets the Size Type of the specified Azure Virtual Machine.
//
// Deprecated: Use [GetSizeOfVirtualMachineContextE] instead.
func GetSizeOfVirtualMachineE(vmName string, resGroupName string, subscriptionID string) (armcompute.VirtualMachineSizeTypes, error) {
	return GetSizeOfVirtualMachineContextE(context.Background(), vmName, resGroupName, subscriptionID)
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
	if vm == nil || vm.Properties == nil || vm.Properties.HardwareProfile == nil ||
		vm.Properties.HardwareProfile.VMSize == nil {
		return ""
	}

	return *vm.Properties.HardwareProfile.VMSize
}

// GetVirtualMachineTags gets the Tags of the specified Virtual Machine as a map.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineTagsContext] instead.
func GetVirtualMachineTags(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) map[string]string {
	t.Helper()

	return GetVirtualMachineTagsContext(t, context.Background(), vmName, resGroupName, subscriptionID)
}

// GetVirtualMachineTagsE gets the Tags of the specified Virtual Machine as a map.
//
// Deprecated: Use [GetVirtualMachineTagsContextE] instead.
func GetVirtualMachineTagsE(vmName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	return GetVirtualMachineTagsContextE(context.Background(), vmName, resGroupName, subscriptionID)
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

	if vm == nil || vm.Tags == nil {
		return tags
	}

	for k, v := range vm.Tags {
		if v == nil {
			continue
		}

		tags[k] = *v
	}

	return tags
}

// ***************************************************** //
// Get multiple Virtual Machines from a Resource Group
// ***************************************************** //

// ListVirtualMachinesForResourceGroup gets a list of all Virtual Machine names in the specified Resource Group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListVirtualMachinesForResourceGroupContext] instead.
func ListVirtualMachinesForResourceGroup(t testing.TestingT, resGroupName string, subscriptionID string) []string {
	t.Helper()

	return ListVirtualMachinesForResourceGroupContext(t, context.Background(), resGroupName, subscriptionID)
}

// ListVirtualMachinesForResourceGroupE gets a list of all Virtual Machine names in the specified Resource Group.
//
// Deprecated: Use [ListVirtualMachinesForResourceGroupContextE] instead.
func ListVirtualMachinesForResourceGroupE(resGroupName string, subscriptionID string) ([]string, error) {
	return ListVirtualMachinesForResourceGroupContextE(context.Background(), resGroupName, subscriptionID)
}

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
	vmClient, err := GetVirtualMachineClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return ListVirtualMachinesForResourceGroupWithClient(ctx, vmClient, resourceGroupName)
}

// ListVirtualMachinesForResourceGroupWithClient gets a list of all Virtual Machine names
// in the specified Resource Group using the provided client.
// This variant is useful for testing with fake clients.
func ListVirtualMachinesForResourceGroupWithClient(ctx context.Context, client *armcompute.VirtualMachinesClient, resourceGroupName string) ([]string, error) {
	return listVirtualMachineNames(ctx, client, resourceGroupName)
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
			if v == nil || v.Name == nil {
				continue
			}

			vmDetails = append(vmDetails, *v.Name)
		}
	}

	return vmDetails, nil
}

// GetVirtualMachinesForResourceGroup gets all Virtual Machine objects in the specified Resource Group. Each
// VM Object represents the entire set of VM compute properties accessible by using the VM name as the map key.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachinesForResourceGroupContext] instead.
func GetVirtualMachinesForResourceGroup(t testing.TestingT, resGroupName string, subscriptionID string) map[string]armcompute.VirtualMachineProperties {
	t.Helper()

	return GetVirtualMachinesForResourceGroupContext(t, context.Background(), resGroupName, subscriptionID)
}

// GetVirtualMachinesForResourceGroupE gets all Virtual Machine objects in the specified Resource Group. Each
// VM Object represents the entire set of VM compute properties accessible by using the VM name as the map key.
//
// Deprecated: Use [GetVirtualMachinesForResourceGroupContextE] instead.
func GetVirtualMachinesForResourceGroupE(resGroupName string, subscriptionID string) (map[string]armcompute.VirtualMachineProperties, error) {
	return GetVirtualMachinesForResourceGroupContextE(context.Background(), resGroupName, subscriptionID)
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
	vmClient, err := GetVirtualMachineClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetVirtualMachinesForResourceGroupWithClient(ctx, vmClient, resourceGroupName)
}

// GetVirtualMachinesForResourceGroupWithClient gets all Virtual Machine objects in the specified Resource Group
// using the provided client. Each VM Object represents the entire set of VM compute properties accessible by
// using the VM name as the map key.
// This variant is useful for testing with fake clients.
func GetVirtualMachinesForResourceGroupWithClient(ctx context.Context, client *armcompute.VirtualMachinesClient, resourceGroupName string) (map[string]armcompute.VirtualMachineProperties, error) {
	return listVirtualMachineProperties(ctx, client, resourceGroupName)
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
			if v == nil || v.Name == nil || v.Properties == nil {
				continue
			}

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
	if vm == nil || vm.VirtualMachine == nil || vm.Properties == nil ||
		vm.Properties.HardwareProfile == nil || vm.Properties.HardwareProfile.VMSize == nil {
		return ""
	}

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

// GetVirtualMachine gets a Virtual Machine in the specified Azure Resource Group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineContext] instead.
func GetVirtualMachine(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) *armcompute.VirtualMachine {
	t.Helper()

	return GetVirtualMachineContext(t, context.Background(), vmName, resGroupName, subscriptionID)
}

// GetVirtualMachineE gets a Virtual Machine in the specified Azure Resource Group.
//
// Deprecated: Use [GetVirtualMachineContextE] instead.
func GetVirtualMachineE(vmName string, resGroupName string, subscriptionID string) (*armcompute.VirtualMachine, error) {
	return GetVirtualMachineContextE(context.Background(), vmName, resGroupName, subscriptionID)
}

// GetVirtualMachineContextE gets a Virtual Machine in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (*armcompute.VirtualMachine, error) {
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := GetVirtualMachineClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetVirtualMachineWithClient(ctx, client, resGroupName, vmName)
}

// GetVirtualMachineWithClient gets a Virtual Machine using the provided client.
// This variant is useful for testing with fake clients.
func GetVirtualMachineWithClient(ctx context.Context, client *armcompute.VirtualMachinesClient, resGroupName string, vmName string) (*armcompute.VirtualMachine, error) {
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
