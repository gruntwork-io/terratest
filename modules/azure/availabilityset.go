package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
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

// AvailabilitySetExists indicates whether the specified Azure Availability Set exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [AvailabilitySetExistsContext] instead.
func AvailabilitySetExists(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	return AvailabilitySetExistsContext(t, context.Background(), avsName, resGroupName, subscriptionID)
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

// AvailabilitySetExistsE indicates whether the specified Azure Availability Set exists.
//
// Deprecated: Use [AvailabilitySetExistsContextE] instead.
func AvailabilitySetExistsE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (bool, error) {
	return AvailabilitySetExistsContextE(t, context.Background(), avsName, resGroupName, subscriptionID)
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

// CheckAvailabilitySetContainsVM checks if the Virtual Machine is contained in the Availability Set VMs.
// This function would fail the test if there is an error.
//
// Deprecated: Use [CheckAvailabilitySetContainsVMContext] instead.
func CheckAvailabilitySetContainsVM(t testing.TestingT, vmName string, avsName string, resGroupName string, subscriptionID string) bool {
	t.Helper()

	return CheckAvailabilitySetContainsVMContext(t, context.Background(), vmName, avsName, resGroupName, subscriptionID)
}

// CheckAvailabilitySetContainsVMContextE checks if the Virtual Machine is contained in the Availability Set VMs.
// The ctx parameter supports cancellation and timeouts.
func CheckAvailabilitySetContainsVMContextE(t testing.TestingT, ctx context.Context, vmName string, avsName string, resGroupName string, subscriptionID string) (bool, error) {
	client, err := CreateAvailabilitySetClientE(subscriptionID)
	if err != nil {
		return false, err
	}

	// Get the Availability Set
	avs, err := client.Get(ctx, resGroupName, avsName)
	if err != nil {
		return false, err
	}

	// Check if the VM is found in the AVS VM collection and return true
	for _, vm := range *avs.VirtualMachines {
		// VM IDs are always ALL CAPS in this property so ignoring case
		if strings.EqualFold(vmName, GetNameFromResourceID(*vm.ID)) {
			return true, nil
		}
	}

	return false, NewNotFoundError("Virtual Machine", vmName, avsName)
}

// CheckAvailabilitySetContainsVME checks if the Virtual Machine is contained in the Availability Set VMs.
//
// Deprecated: Use [CheckAvailabilitySetContainsVMContextE] instead.
func CheckAvailabilitySetContainsVME(t testing.TestingT, vmName string, avsName string, resGroupName string, subscriptionID string) (bool, error) {
	return CheckAvailabilitySetContainsVMContextE(t, context.Background(), vmName, avsName, resGroupName, subscriptionID)
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

// GetAvailabilitySetVMNamesInCaps gets a list of VM names in the specified Azure Availability Set.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetAvailabilitySetVMNamesInCapsContext] instead.
func GetAvailabilitySetVMNamesInCaps(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) []string {
	t.Helper()

	return GetAvailabilitySetVMNamesInCapsContext(t, context.Background(), avsName, resGroupName, subscriptionID)
}

// GetAvailabilitySetVMNamesInCapsContextE gets a list of VM names in the specified Azure Availability Set.
// The ctx parameter supports cancellation and timeouts.
func GetAvailabilitySetVMNamesInCapsContextE(t testing.TestingT, ctx context.Context, avsName string, resGroupName string, subscriptionID string) ([]string, error) {
	client, err := CreateAvailabilitySetClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	avs, err := client.Get(ctx, resGroupName, avsName)
	if err != nil {
		return nil, err
	}

	vms := []string{}

	// Get the names for all VMs in the Availability Set
	for _, vm := range *avs.VirtualMachines {
		// IDs are returned in ALL CAPS for this property
		if vmName := GetNameFromResourceID(*vm.ID); len(vmName) > 0 {
			vms = append(vms, vmName)
		}
	}

	return vms, nil
}

// GetAvailabilitySetVMNamesInCapsE gets a list of VM names in the specified Azure Availability Set.
//
// Deprecated: Use [GetAvailabilitySetVMNamesInCapsContextE] instead.
func GetAvailabilitySetVMNamesInCapsE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) ([]string, error) {
	return GetAvailabilitySetVMNamesInCapsContextE(t, context.Background(), avsName, resGroupName, subscriptionID)
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

// GetAvailabilitySetFaultDomainCount gets the Fault Domain Count for the specified Azure Availability Set.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetAvailabilitySetFaultDomainCountContext] instead.
func GetAvailabilitySetFaultDomainCount(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) int32 {
	t.Helper()

	return GetAvailabilitySetFaultDomainCountContext(t, context.Background(), avsName, resGroupName, subscriptionID)
}

// GetAvailabilitySetFaultDomainCountContextE gets the Fault Domain Count for the specified Azure Availability Set.
// The ctx parameter supports cancellation and timeouts.
func GetAvailabilitySetFaultDomainCountContextE(t testing.TestingT, ctx context.Context, avsName string, resGroupName string, subscriptionID string) (int32, error) {
	avs, err := GetAvailabilitySetContextE(t, ctx, avsName, resGroupName, subscriptionID)
	if err != nil {
		return -1, err
	}

	return *avs.PlatformFaultDomainCount, nil
}

// GetAvailabilitySetFaultDomainCountE gets the Fault Domain Count for the specified Azure Availability Set.
//
// Deprecated: Use [GetAvailabilitySetFaultDomainCountContextE] instead.
func GetAvailabilitySetFaultDomainCountE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (int32, error) {
	return GetAvailabilitySetFaultDomainCountContextE(t, context.Background(), avsName, resGroupName, subscriptionID)
}

// GetAvailabilitySetContextE gets an Availability Set in the specified Azure Resource Group.
// The ctx parameter supports cancellation and timeouts.
func GetAvailabilitySetContextE(t testing.TestingT, ctx context.Context, avsName string, resGroupName string, subscriptionID string) (*compute.AvailabilitySet, error) {
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
	avs, err := client.Get(ctx, resGroupName, avsName)
	if err != nil {
		return nil, err
	}

	return &avs, nil
}

// GetAvailabilitySetE gets an Availability Set in the specified Azure Resource Group.
//
// Deprecated: Use [GetAvailabilitySetContextE] instead.
func GetAvailabilitySetE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (*compute.AvailabilitySet, error) {
	return GetAvailabilitySetContextE(t, context.Background(), avsName, resGroupName, subscriptionID)
}

// GetAvailabilitySetClientE gets a new Availability Set client in the specified Azure Subscription.
// TODO: remove in next version
func GetAvailabilitySetClientE(subscriptionID string) (*compute.AvailabilitySetsClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the Availability Set client
	client := compute.NewAvailabilitySetsClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	return &client, nil
}
