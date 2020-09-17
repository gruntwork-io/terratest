package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetVirtualMachineClient is a helper function that will setup an Azure Virtual Machine client on your behalf
func GetVirtualMachineClient(subscriptionID string) (*compute.VirtualMachinesClient, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Create a VM client
	vmClient := compute.NewVirtualMachinesClient(subscriptionID)

	// Create an authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	// Attach authorizer to the client
	vmClient.Authorizer = *authorizer

	return &vmClient, nil
}

// GetSizeOfVirtualMachine gets the size type of the given Azure Virtual Machine
func GetSizeOfVirtualMachine(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) compute.VirtualMachineSizeTypes {
	size, err := GetSizeOfVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return size
}

// GetSizeOfVirtualMachineE gets the size type of the given Azure Virtual Machine
func GetSizeOfVirtualMachineE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (compute.VirtualMachineSizeTypes, error) {
	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return "", err
	}

	// Create a VM client
	vmClient, err := GetVirtualMachineClient(subscriptionID)
	if err != nil {
		return "", err
	}

	// Get the details of the target virtual machine
	vm, err := vmClient.Get(context.Background(), resGroupName, vmName, compute.InstanceView)
	if err != nil {
		return "", err
	}

	return vm.VirtualMachineProperties.HardwareProfile.VMSize, nil
}

// GetTagsForVirtualMachine gets the tags of the given Virtual Machine as a map
func GetTagsForVirtualMachine(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) map[string]string {
	tags, err := GetTagsForVirtualMachineE(t, vmName, resGroupName, subscriptionID)
	require.NoError(t, err)

	return tags
}

// GetTagsForVirtualMachineE gets the tags of the given Virtual Machine as a map
func GetTagsForVirtualMachineE(t testing.TestingT, vmName string, resGroupName string, subscriptionID string) (map[string]string, error) {
	// Setup a blank map to populate and return
	tags := make(map[string]string)

	// Validate resource group name and subscription ID
	resGroupName, err := getTargetAzureResourceGroupName(resGroupName)
	if err != nil {
		return tags, err
	}

	// Create a VM client
	vmClient, err := GetVirtualMachineClient(subscriptionID)
	if err != nil {
		return tags, err
	}

	// Get the details of the target virtual machine
	vm, err := vmClient.Get(context.Background(), resGroupName, vmName, compute.InstanceView)
	if err != nil {
		return tags, err
	}

	// Range through existing tags and populate above map accordingly
	for k, v := range vm.Tags {
		tags[k] = *v
	}

	return tags, nil
}
