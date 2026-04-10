package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// AvailabilitySetExistsContext indicates whether the specified Azure Availability Set exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func AvailabilitySetExistsContext(t testing.TestingT, ctx context.Context, avsName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := AvailabilitySetExistsContextE(t, ctx, avsName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// AvailabilitySetExistsContextE indicates whether the specified Azure Availability Set exists.
// The ctx parameter supports cancellation and timeouts.
func AvailabilitySetExistsContextE(t testing.TestingT, ctx context.Context, avsName string, resGroupName string, subscriptionID string) (bool, error) {
	_, err := GetAvailabilitySetContextE(t, ctx, avsName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// CheckAvailabilitySetContainsVMContext checks if the Virtual Machine is contained in the Availability Set VMs.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CheckAvailabilitySetContainsVMContext(t testing.TestingT, ctx context.Context, vmName string, avsName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	success, err := CheckAvailabilitySetContainsVMContextE(t, ctx, vmName, avsName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return success
}

// CheckAvailabilitySetContainsVMContextE checks if the Virtual Machine is contained in the Availability Set VMs.
// The ctx parameter supports cancellation and timeouts.
func CheckAvailabilitySetContainsVMContextE(t testing.TestingT, ctx context.Context, vmName string, avsName string, resGroupName string, subscriptionID string) (bool, error) {
	client, err := CreateAvailabilitySetClientE(subscriptionID)
	if err != nil {
		return false, err
	}

	// Get the Availability Set
	resp, err := client.Get(ctx, resGroupName, avsName, nil)
	if err != nil {
		return false, err
	}

	return AvsContainsVM(vmName, avsName, resp.Properties.VirtualMachines)
}

// GetAvailabilitySetVMNamesInCapsContext gets a list of VM names in the specified Azure Availability Set.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAvailabilitySetVMNamesInCapsContext(t testing.TestingT, ctx context.Context, avsName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	vms, err := GetAvailabilitySetVMNamesInCapsContextE(t, ctx, avsName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return vms
}

// GetAvailabilitySetVMNamesInCapsContextE gets a list of VM names in the specified Azure Availability Set.
// The ctx parameter supports cancellation and timeouts.
func GetAvailabilitySetVMNamesInCapsContextE(t testing.TestingT, ctx context.Context, avsName string, resGroupName string, subscriptionID string) ([]string, error) {
	client, err := CreateAvailabilitySetClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, resGroupName, avsName, nil)
	if err != nil {
		return nil, err
	}

	return ExtractAvsVMNames(resp.Properties.VirtualMachines), nil
}

// GetAvailabilitySetFaultDomainCountContext gets the Fault Domain Count for the specified Azure Availability Set.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAvailabilitySetFaultDomainCountContext(t testing.TestingT, ctx context.Context, avsName string, resGroupName string, subscriptionID string) int32 {
	t.Helper()

	avsFaultDomainCount, err := GetAvailabilitySetFaultDomainCountContextE(t, ctx, avsName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return avsFaultDomainCount
}

// GetAvailabilitySetFaultDomainCountContextE gets the Fault Domain Count for the specified Azure Availability Set.
// The ctx parameter supports cancellation and timeouts.
func GetAvailabilitySetFaultDomainCountContextE(t testing.TestingT, ctx context.Context, avsName string, resGroupName string, subscriptionID string) (int32, error) {
	avs, err := GetAvailabilitySetContextE(t, ctx, avsName, resGroupName, subscriptionID)
	if err != nil {
		return -1, err
	}

	return *avs.Properties.PlatformFaultDomainCount, nil
}

// GetAvailabilitySetContextE gets an Availability Set in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetAvailabilitySetContextE(t testing.TestingT, ctx context.Context, avsName string, resGroupName string, subscriptionID string) (*armcompute.AvailabilitySet, error) {
	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := CreateAvailabilitySetClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Availability Set
	resp, err := client.Get(ctx, resGroupName, avsName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.AvailabilitySet, nil
}

// AvsContainsVM checks if the named VM is in the Availability Set's VM list using case-insensitive matching.
func AvsContainsVM(vmName string, avsName string, vms []*armcompute.SubResource) (bool, error) {
	for _, vm := range vms {
		// VM IDs are always ALL CAPS in this property so ignoring case
		if strings.EqualFold(vmName, GetNameFromResourceID(*vm.ID)) {
			return true, nil
		}
	}

	return false, NewNotFoundError("Virtual Machine", vmName, avsName)
}

// ExtractAvsVMNames extracts the VM names from an Availability Set's VM list.
func ExtractAvsVMNames(vms []*armcompute.SubResource) []string {
	var names []string

	for _, vm := range vms {
		// IDs are returned in ALL CAPS for this property
		if vmName := GetNameFromResourceID(*vm.ID); len(vmName) > 0 {
			names = append(names, vmName)
		}
	}

	return names
}
