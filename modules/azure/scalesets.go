package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetVirtualMachineScaleSetsClient is a helper function that will setup an Azure Virtual Machine Scale Sets client on your behalf.
// This function would fail the test if there is an error.
func GetVirtualMachineScaleSetsClient(t testing.TestingT, subscriptionID string) *compute.VirtualMachineScaleSetsClient {
	vmssClient, err := GetVirtualMachineScaleSetsClientE(subscriptionID)
	require.NoError(t, err)
	return vmssClient
}

// GetVirtualMachineScaleSetsClientE is a helper function that will setup an Azure Virtual Machine Scale Sets client on your behalf.
func GetVirtualMachineScaleSetsClientE(subscriptionID string) (*compute.VirtualMachineScaleSetsClient, error) {
	// Create a VM Scale Sets client
	vmssClient, err := CreateVirtualMachineScaleSetsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	vmssClient.Authorizer = *authorizer
	return vmssClient, nil
}

// VirtualMachineScaleSetExists indicates whether the specifcied Azure Virtual Machine Scale Set exists.
// This function would fail the test if there is an error.
func VirtualMachineScaleSetExists(t testing.TestingT, vmssName string, resGroupName string, subscriptionID string) bool {
	exists, err := VirtualMachineScaleSetExistsE(vmssName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// VirtualMachineScaleSetsExistsE indicates whether the specifcied Azure Virtual Machine exists.
func VirtualMachineScaleSetExistsE(vmssName string, resGroupName string, subscriptionID string) (bool, error) {
	// Get VM Object
	_, err := GetVirtualMachineScaleSetE(vmssName, resGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetVirtualMachineScaleSetTags gets the Tags of the specified Virtual Machine Scale Set as a map.
// This function would fail the test if there is an error.
func GetVirtualMachineScaleSetTags(t testing.TestingT, vmssName string, resGroupName string, subscriptionID string) map[string]string {
	tags, err := GetVirtualMachineScaleSetTagsE(vmssName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return tags
}

// GetVirtualMachineScaleSetTagsE gets the Tags of the specified Virtual Machine Scale Set as a map.
func GetVirtualMachineScaleSetTagsE(vmssName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	// Setup a blank map to populate and return
	tags := make(map[string]string)

	// Get VM Scale Set Object
	vmss, err := GetVirtualMachineScaleSetE(vmssName, resGroupName, subscriptionID)
	if err != nil {
		return tags, err
	}

	// Range through existing tags and populate above map accordingly
	for k, v := range vmss.Tags {
		tags[k] = *v
	}

	return tags, nil
}

// GetVirtualMachineScaleSet gets a Virtual Machine Scale Set in the specified Azure Resource Group.
// This function would fail the test if there is an error.
func GetVirtualMachineScaleSet(t testing.TestingT, vmssName string, resGroupName string, subscriptionID string) *compute.VirtualMachineScaleSet {
	vmss, err := GetVirtualMachineScaleSetE(vmssName, resGroupName, subscriptionID)
	require.NoError(t, err)
	return vmss
}

// GetVirtualMachineScaleSetE gets a Virtual Machine Scale Set in the specified Azure Resource Group.
func GetVirtualMachineScaleSetE(vmssName string, resGroupName string, subscriptionID string) (*compute.VirtualMachineScaleSet, error) {
	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return nil, err
	}

	// Get the client reference
	client, err := GetVirtualMachineScaleSetsClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	vmss, err := client.Get(context.Background(), resGroupName, vmssName)
	if err != nil {
		return nil, err
	}

	return &vmss, nil
}
