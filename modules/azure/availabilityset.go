package azure

import (
	"context"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// AvailabilitySetExists indicates whether the specified Azure Availability Set exists.
// This function would fail the test if there is an error.
func AvailabilitySetExists(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) bool {
	exists, err := AvailabilitySetExistsE(t, avsName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// AvailabilitySetExistsE indicates whether the specified Azure Availability Set exists
func AvailabilitySetExistsE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (bool, error) {
	_, err := GetAvailabilitySetE(t, avsName, resGroupName, subscriptionID)
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
func CheckAvailabilitySetContainsVM(t testing.TestingT, vmName string, avsName string, resGroupName string, subscriptionID string) bool {
	success, err := CheckAvailabilitySetContainsVME(t, vmName, avsName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return success
}

// CheckAvailabilitySetContainsVME checks if the Virtual Machine is contained in the Availability Set VMs
func CheckAvailabilitySetContainsVME(t testing.TestingT, vmName string, avsName string, resGroupName string, subscriptionID string) (bool, error) {
	client, err := CreateAvailabilitySetClientE(subscriptionID)
	if err != nil {
		return false, err
	}

	// Get the Availability Set
	avs, err := client.Get(context.Background(), resGroupName, avsName)
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

// GetAvailabilitySetVMNamesInCaps gets a list of VM names in the specified Azure Availability Set.
// This function would fail the test if there is an error.
func GetAvailabilitySetVMNamesInCaps(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) []string {
	vms, err := GetAvailabilitySetVMNamesInCapsE(t, avsName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return vms
}

// GetAvailabilitySetVMNamesInCapsE gets a list of VM names in the specified Azure Availability Set
func GetAvailabilitySetVMNamesInCapsE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) ([]string, error) {
	client, err := CreateAvailabilitySetClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	avs, err := client.Get(context.Background(), resGroupName, avsName)
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

// GetAvailabilitySetFaultDomainCount gets the Fault Domain Count for the specified Azure Availability Set.
// This function would fail the test if there is an error.
func GetAvailabilitySetFaultDomainCount(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) int32 {
	avsFaultDomainCount, err := GetAvailabilitySetFaultDomainCountE(t, avsName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return avsFaultDomainCount
}

// GetAvailabilitySetFaultDomainCountE gets the Fault Domain Count for the specified Azure Availability Set
func GetAvailabilitySetFaultDomainCountE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (int32, error) {
	avs, err := GetAvailabilitySetE(t, avsName, resGroupName, subscriptionID)
	if err != nil {
		return -1, err
	}
	return *avs.PlatformFaultDomainCount, nil
}

// GetAvailabilitySetE gets an Availability Set in the specified Azure Resource Group
func GetAvailabilitySetE(t testing.TestingT, avsName string, resGroupName string, subscriptionID string) (*compute.AvailabilitySet, error) {
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
	avs, err := client.Get(context.Background(), resGroupName, avsName)
	if err != nil {
		return nil, err
	}

	return &avs, nil
}

// GetAvailabilitySetClientE gets a new Availability Set client in the specified Azure Subscription
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
