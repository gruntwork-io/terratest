package azure

import (
	"context"
	"errors"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// AvailabilitySetExists indicates whether the specified Azure Availability Set exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [AvailabilitySetExistsContext] instead.
func AvailabilitySetExists(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	return AvailabilitySetExistsContext(t, context.Background(), avsName, resGroupName, subscriptionID)
}

// AvailabilitySetExistsE indicates whether the specified Azure Availability Set exists.
//
// Deprecated: Use [AvailabilitySetExistsContextE] instead.
func AvailabilitySetExistsE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (bool, error) {
	t.Helper()

	return AvailabilitySetExistsContextE(t, context.Background(), avsName, resGroupName, subscriptionID)
}

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

// CheckAvailabilitySetContainsVM checks if the Virtual Machine is contained in the Availability Set VMs.
// This function would fail the test if there is an error.
//
// Deprecated: Use [CheckAvailabilitySetContainsVMContext] instead.
func CheckAvailabilitySetContainsVM(t testing.TestingT, vmName string, avsName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	return CheckAvailabilitySetContainsVMContext(t, context.Background(), vmName, avsName, resGroupName, subscriptionID)
}

// CheckAvailabilitySetContainsVME checks if the Virtual Machine is contained in the Availability Set VMs.
//
// Deprecated: Use [CheckAvailabilitySetContainsVMContextE] instead.
func CheckAvailabilitySetContainsVME(t testing.TestingT, vmName string, avsName string, resGroupName string, subscriptionID string) (bool, error) {
	t.Helper()

	return CheckAvailabilitySetContainsVMContextE(t, context.Background(), vmName, avsName, resGroupName, subscriptionID)
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
	client, err := CreateAvailabilitySetClientContextE(ctx, subscriptionID)
	if err != nil {
		return false, err
	}

	return CheckAvailabilitySetContainsVMWithClient(ctx, client, resGroupName, avsName, vmName)
}

// CheckAvailabilitySetContainsVMWithClient checks if the Virtual Machine is contained in the Availability Set VMs
// using the provided AvailabilitySetsClient.
func CheckAvailabilitySetContainsVMWithClient(ctx context.Context, client *armcompute.AvailabilitySetsClient, resGroupName string, avsName string, vmName string) (bool, error) {
	resp, err := client.Get(ctx, resGroupName, avsName, nil)
	if err != nil {
		return false, err
	}

	for _, vm := range resp.Properties.VirtualMachines {
		if vm.ID == nil {
			continue
		}
		// VM IDs are always ALL CAPS in this property so ignoring case
		if strings.EqualFold(vmName, GetNameFromResourceID(*vm.ID)) {
			return true, nil
		}
	}

	return false, NewNotFoundError("Virtual Machine", vmName, avsName)
}

// GetAvailabilitySetVMNamesInCaps gets a list of VM names in the specified Azure Availability Set.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetAvailabilitySetVMNamesInCapsContext] instead.
func GetAvailabilitySetVMNamesInCaps(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	return GetAvailabilitySetVMNamesInCapsContext(t, context.Background(), avsName, resGroupName, subscriptionID)
}

// GetAvailabilitySetVMNamesInCapsE gets a list of VM names in the specified Azure Availability Set.
//
// Deprecated: Use [GetAvailabilitySetVMNamesInCapsContextE] instead.
func GetAvailabilitySetVMNamesInCapsE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) ([]string, error) {
	t.Helper()

	return GetAvailabilitySetVMNamesInCapsContextE(t, context.Background(), avsName, resGroupName, subscriptionID)
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
	client, err := CreateAvailabilitySetClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetAvailabilitySetVMNamesInCapsWithClient(ctx, client, resGroupName, avsName)
}

// GetAvailabilitySetVMNamesInCapsWithClient gets a list of VM names in the specified Azure Availability Set
// using the provided AvailabilitySetsClient.
func GetAvailabilitySetVMNamesInCapsWithClient(ctx context.Context, client *armcompute.AvailabilitySetsClient, resGroupName string, avsName string) ([]string, error) {
	resp, err := client.Get(ctx, resGroupName, avsName, nil)
	if err != nil {
		return nil, err
	}

	vms := []string{}

	for _, vm := range resp.Properties.VirtualMachines {
		if vm.ID == nil {
			continue
		}
		// IDs are returned in ALL CAPS for this property
		if vmName := GetNameFromResourceID(*vm.ID); len(vmName) > 0 {
			vms = append(vms, vmName)
		}
	}

	return vms, nil
}

// GetAvailabilitySetFaultDomainCount gets the Fault Domain Count for the specified Azure Availability Set.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetAvailabilitySetFaultDomainCountContext] instead.
func GetAvailabilitySetFaultDomainCount(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) int32 {
	t.Helper()

	return GetAvailabilitySetFaultDomainCountContext(t, context.Background(), avsName, resGroupName, subscriptionID)
}

// GetAvailabilitySetFaultDomainCountE gets the Fault Domain Count for the specified Azure Availability Set.
//
// Deprecated: Use [GetAvailabilitySetFaultDomainCountContextE] instead.
func GetAvailabilitySetFaultDomainCountE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (int32, error) {
	t.Helper()

	return GetAvailabilitySetFaultDomainCountContextE(t, context.Background(), avsName, resGroupName, subscriptionID)
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

	return ExtractAvailabilitySetFaultDomainCount(avs)
}

// ExtractAvailabilitySetFaultDomainCount gets the Fault Domain Count from the provided AvailabilitySet.
func ExtractAvailabilitySetFaultDomainCount(avs *armcompute.AvailabilitySet) (int32, error) {
	if avs == nil || avs.Properties == nil || avs.Properties.PlatformFaultDomainCount == nil {
		return -1, errors.New("availability set has no fault domain count")
	}

	return *avs.Properties.PlatformFaultDomainCount, nil
}

// GetAvailabilitySetE gets an Availability Set in the specified Azure Resource Group.
//
// Deprecated: Use [GetAvailabilitySetContextE] instead.
func GetAvailabilitySetE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (*armcompute.AvailabilitySet, error) {
	t.Helper()

	return GetAvailabilitySetContextE(t, context.Background(), avsName, resGroupName, subscriptionID)
}

// GetAvailabilitySetContextE gets an Availability Set in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetAvailabilitySetContextE(t testing.TestingT, ctx context.Context, avsName string, resGroupName string, subscriptionID string) (*armcompute.AvailabilitySet, error) {
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	client, err := CreateAvailabilitySetClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetAvailabilitySetWithClient(ctx, client, resGroupName, avsName)
}

// GetAvailabilitySetWithClient gets an Availability Set using the provided AvailabilitySetsClient.
func GetAvailabilitySetWithClient(ctx context.Context, client *armcompute.AvailabilitySetsClient, resGroupName string, avsName string) (*armcompute.AvailabilitySet, error) {
	resp, err := client.Get(ctx, resGroupName, avsName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.AvailabilitySet, nil
}
