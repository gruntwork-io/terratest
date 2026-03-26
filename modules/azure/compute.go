package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetVirtualMachineClient is a helper function that will setup an Azure Virtual Machine client on your behalf.
// This function would fail the test if there is an error.
func GetVirtualMachineClient(t testing.TestingT, subscriptionID string) *compute.VirtualMachinesClient {
	t.Helper()

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

// VirtualMachineExistsContext indicates whether the specified Azure Virtual Machine exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func VirtualMachineExistsContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := VirtualMachineExistsContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// VirtualMachineExists indicates whether the specified Azure Virtual Machine exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [VirtualMachineExistsContext] instead.
func VirtualMachineExists(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	return VirtualMachineExistsContext(t, context.Background(), vmName, resGroupName, subscriptionID) //nolint:staticcheck
}

// VirtualMachineExistsContextE indicates whether the specified Azure Virtual Machine exists.
// The ctx parameter supports cancellation and timeouts.
func VirtualMachineExistsContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get VM Object
	_, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// VirtualMachineExistsE indicates whether the specified Azure Virtual Machine exists.
//
// Deprecated: Use [VirtualMachineExistsContextE] instead.
func VirtualMachineExistsE(vmName string, resGroupName string, subscriptionID string) (bool, error) {
	return VirtualMachineExistsContextE(context.Background(), vmName, resGroupName, subscriptionID)
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

// GetVirtualMachineNics gets a list of Network Interface names for a specified Azure Virtual Machine.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineNicsContext] instead.
func GetVirtualMachineNics(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	return GetVirtualMachineNicsContext(t, context.Background(), vmName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetVirtualMachineNicsContextE gets a list of Network Interface names for a specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineNicsContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) ([]string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
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

// GetVirtualMachineNicsE gets a list of Network Interface names for a specified Azure Virtual Machine.
//
// Deprecated: Use [GetVirtualMachineNicsContextE] instead.
func GetVirtualMachineNicsE(vmName string, resGroupName string, subscriptionID string) ([]string, error) {
	return GetVirtualMachineNicsContextE(context.Background(), vmName, resGroupName, subscriptionID)
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

// GetVirtualMachineManagedDisks gets the list of Managed Disk names of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineManagedDisksContext] instead.
func GetVirtualMachineManagedDisks(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	return GetVirtualMachineManagedDisksContext(t, context.Background(), vmName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetVirtualMachineManagedDisksContextE gets the list of Managed Disk names of the specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineManagedDisksContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) ([]string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
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

// GetVirtualMachineManagedDisksE gets the list of Managed Disk names of the specified Azure Virtual Machine.
//
// Deprecated: Use [GetVirtualMachineManagedDisksContextE] instead.
func GetVirtualMachineManagedDisksE(vmName string, resGroupName string, subscriptionID string) ([]string, error) {
	return GetVirtualMachineManagedDisksContextE(context.Background(), vmName, resGroupName, subscriptionID)
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

// GetVirtualMachineOSDiskName gets the OS Disk name of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineOSDiskNameContext] instead.
func GetVirtualMachineOSDiskName(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	t.Helper()

	return GetVirtualMachineOSDiskNameContext(t, context.Background(), vmName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetVirtualMachineOSDiskNameContextE gets the OS Disk name of the specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineOSDiskNameContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return *vm.StorageProfile.OsDisk.Name, nil
}

// GetVirtualMachineOSDiskNameE gets the OS Disk name of the specified Azure Virtual Machine.
//
// Deprecated: Use [GetVirtualMachineOSDiskNameContextE] instead.
func GetVirtualMachineOSDiskNameE(vmName string, resGroupName string, subscriptionID string) (string, error) {
	return GetVirtualMachineOSDiskNameContextE(context.Background(), vmName, resGroupName, subscriptionID)
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

// GetVirtualMachineAvailabilitySetID gets the Availability Set ID of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineAvailabilitySetIDContext] instead.
func GetVirtualMachineAvailabilitySetID(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) string {
	t.Helper()

	return GetVirtualMachineAvailabilitySetIDContext(t, context.Background(), vmName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetVirtualMachineAvailabilitySetIDContextE gets the Availability Set ID of the specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineAvailabilitySetIDContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (string, error) {
	// Get VM Object
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
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

// GetVirtualMachineAvailabilitySetIDE gets the Availability Set ID of the specified Azure Virtual Machine.
//
// Deprecated: Use [GetVirtualMachineAvailabilitySetIDContextE] instead.
func GetVirtualMachineAvailabilitySetIDE(vmName string, resGroupName string, subscriptionID string) (string, error) {
	return GetVirtualMachineAvailabilitySetIDContextE(context.Background(), vmName, resGroupName, subscriptionID)
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

// GetVirtualMachineImage gets the Image of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineImageContext] instead.
func GetVirtualMachineImage(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) VMImage {
	t.Helper()

	return GetVirtualMachineImageContext(t, context.Background(), vmName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetVirtualMachineImageContextE gets the Image of the specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineImageContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (VMImage, error) {
	var vmImage VMImage

	// Get VM Object
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
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

// GetVirtualMachineImageE gets the Image of the specified Azure Virtual Machine.
//
// Deprecated: Use [GetVirtualMachineImageContextE] instead.
func GetVirtualMachineImageE(vmName string, resGroupName string, subscriptionID string) (VMImage, error) {
	return GetVirtualMachineImageContextE(context.Background(), vmName, resGroupName, subscriptionID)
}

// GetSizeOfVirtualMachineContext gets the Size Type of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSizeOfVirtualMachineContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) compute.VirtualMachineSizeTypes {
	t.Helper()

	size, err := GetSizeOfVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return size
}

// GetSizeOfVirtualMachine gets the Size Type of the specified Azure Virtual Machine.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetSizeOfVirtualMachineContext] instead.
func GetSizeOfVirtualMachine(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) compute.VirtualMachineSizeTypes {
	t.Helper()

	return GetSizeOfVirtualMachineContext(t, context.Background(), vmName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetSizeOfVirtualMachineContextE gets the Size Type of the specified Azure Virtual Machine.
// The ctx parameter supports cancellation and timeouts.
func GetSizeOfVirtualMachineContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (compute.VirtualMachineSizeTypes, error) {
	// Get VM Object
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		return "", err
	}

	return vm.VirtualMachineProperties.HardwareProfile.VMSize, nil
}

// GetSizeOfVirtualMachineE gets the Size Type of the specified Azure Virtual Machine.
//
// Deprecated: Use [GetSizeOfVirtualMachineContextE] instead.
func GetSizeOfVirtualMachineE(vmName string, resGroupName string, subscriptionID string) (compute.VirtualMachineSizeTypes, error) {
	return GetSizeOfVirtualMachineContextE(context.Background(), vmName, resGroupName, subscriptionID)
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

// GetVirtualMachineTags gets the Tags of the specified Virtual Machine as a map.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineTagsContext] instead.
func GetVirtualMachineTags(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) map[string]string {
	t.Helper()

	return GetVirtualMachineTagsContext(t, context.Background(), vmName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetVirtualMachineTagsContextE gets the Tags of the specified Virtual Machine as a map.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineTagsContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	// Setup a blank map to populate and return
	tags := make(map[string]string)

	// Get VM Object
	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	if err != nil {
		return tags, err
	}

	// Range through existing tags and populate above map accordingly
	for k, v := range vm.Tags {
		tags[k] = *v
	}

	return tags, nil
}

// GetVirtualMachineTagsE gets the Tags of the specified Virtual Machine as a map.
//
// Deprecated: Use [GetVirtualMachineTagsContextE] instead.
func GetVirtualMachineTagsE(vmName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	return GetVirtualMachineTagsContextE(context.Background(), vmName, resGroupName, subscriptionID)
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

// ListVirtualMachinesForResourceGroup gets a list of all Virtual Machine names in the specified Resource Group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [ListVirtualMachinesForResourceGroupContext] instead.
func ListVirtualMachinesForResourceGroup(t testing.TestingT, resGroupName string, subscriptionID string) []string {
	t.Helper()

	return ListVirtualMachinesForResourceGroupContext(t, context.Background(), resGroupName, subscriptionID) //nolint:staticcheck
}

// ListVirtualMachinesForResourceGroupContextE gets a list of all Virtual Machine names in the specified Resource Group.
// The ctx parameter supports cancellation and timeouts.
func ListVirtualMachinesForResourceGroupContextE(ctx context.Context, resourceGroupName string, subscriptionID string) ([]string, error) {
	var vmDetails []string

	vmClient, err := GetVirtualMachineClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	vms, err := vmClient.List(ctx, resourceGroupName)
	if err != nil {
		return nil, err
	}

	for _, v := range vms.Values() {
		vmDetails = append(vmDetails, *v.Name)
	}

	return vmDetails, nil
}

// ListVirtualMachinesForResourceGroupE gets a list of all Virtual Machine names in the specified Resource Group.
//
// Deprecated: Use [ListVirtualMachinesForResourceGroupContextE] instead.
func ListVirtualMachinesForResourceGroupE(resourceGroupName string, subscriptionID string) ([]string, error) {
	return ListVirtualMachinesForResourceGroupContextE(context.Background(), resourceGroupName, subscriptionID)
}

// GetVirtualMachinesForResourceGroupContext gets all Virtual Machine objects in the specified Resource Group. Each
// VM Object represents the entire set of VM compute properties accessible by using the VM name as the map key.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachinesForResourceGroupContext(t testing.TestingT, ctx context.Context, resGroupName string, subscriptionID string) map[string]compute.VirtualMachineProperties {
	t.Helper()

	vms, err := GetVirtualMachinesForResourceGroupContextE(ctx, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vms
}

// GetVirtualMachinesForResourceGroup gets all Virtual Machine objects in the specified Resource Group. Each
// VM Object represents the entire set of VM compute properties accessible by using the VM name as the map key.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachinesForResourceGroupContext] instead.
func GetVirtualMachinesForResourceGroup(t testing.TestingT, resGroupName string, subscriptionID string) map[string]compute.VirtualMachineProperties {
	t.Helper()

	return GetVirtualMachinesForResourceGroupContext(t, context.Background(), resGroupName, subscriptionID) //nolint:staticcheck
}

// GetVirtualMachinesForResourceGroupContextE gets all Virtual Machine objects in the specified Resource Group. Each
// VM Object represents the entire set of VM compute properties accessible by using the VM name as the map key.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachinesForResourceGroupContextE(ctx context.Context, resourceGroupName string, subscriptionID string) (map[string]compute.VirtualMachineProperties, error) {
	// Create VM Client
	vmClient, err := GetVirtualMachineClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the list of VMs in the Resource Group
	vms, err := vmClient.List(ctx, resourceGroupName)
	if err != nil {
		return nil, err
	}

	// Get the VMs in the Resource Group.
	vmDetails := make(map[string]compute.VirtualMachineProperties, len(vms.Values()))

	for _, v := range vms.Values() {
		// VM name and machine properties are required for each VM, no nil check required.
		vmDetails[*v.Name] = *v.VirtualMachineProperties
	}

	return vmDetails, nil
}

// GetVirtualMachinesForResourceGroupE gets all Virtual Machine objects in the specified Resource Group. Each
// VM Object represents the entire set of VM compute properties accessible by using the VM name as the map key.
//
// Deprecated: Use [GetVirtualMachinesForResourceGroupContextE] instead.
func GetVirtualMachinesForResourceGroupE(resourceGroupName string, subscriptionID string) (map[string]compute.VirtualMachineProperties, error) {
	return GetVirtualMachinesForResourceGroupContextE(context.Background(), resourceGroupName, subscriptionID)
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

// GetVirtualMachineContext gets a Virtual Machine in the specified Azure Resource Group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineContext(t testing.TestingT, ctx context.Context, vmName string, resGroupName string, subscriptionID string) *compute.VirtualMachine {
	t.Helper()

	vm, err := GetVirtualMachineContextE(ctx, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vm
}

// GetVirtualMachine gets a Virtual Machine in the specified Azure Resource Group.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetVirtualMachineContext] instead.
func GetVirtualMachine(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) *compute.VirtualMachine {
	t.Helper()

	return GetVirtualMachineContext(t, context.Background(), vmName, resGroupName, subscriptionID) //nolint:staticcheck
}

// GetVirtualMachineContextE gets a Virtual Machine in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetVirtualMachineContextE(ctx context.Context, vmName string, resGroupName string, subscriptionID string) (*compute.VirtualMachine, error) {
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

	vm, err := client.Get(ctx, resGroupName, vmName, compute.InstanceView)
	if err != nil {
		return nil, err
	}

	return &vm, nil
}

// GetVirtualMachineE gets a Virtual Machine in the specified Azure Resource Group.
//
// Deprecated: Use [GetVirtualMachineContextE] instead.
func GetVirtualMachineE(vmName string, resGroupName string, subscriptionID string) (*compute.VirtualMachine, error) {
	return GetVirtualMachineContextE(context.Background(), vmName, resGroupName, subscriptionID)
}
