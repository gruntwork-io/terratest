package azure

import (
	"context"

	machinelearning "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning/v3"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// MachinelearningWorkspaceExists indicates whether a machine learning workspace exists within a subscription; otherwise false
// This function would fail the test if there is an error.
func MachinelearningWorkspaceExists(t testing.TestingT, resourceGroupName string, workspaceName string, subscriptionID string) bool {
	result, err := MachinelearningWorkspaceExistsE(t, resourceGroupName, workspaceName, subscriptionID)
	require.NoError(t, err)
	return result
}

// MachinelearningWorkspaceExistsE indicates whether a machine learning workspace exists within a subscription
func MachinelearningWorkspaceExistsE(t testing.TestingT, resourceGroupName, workspaceName, subscriptionID string) (bool, error) {
	workspace, err := GetMachinelearningWorkspaceE(t, resourceGroupName, workspaceName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return (workspaceName == *workspace.Name), nil

}

// GetMachinelearningWorkspace is a helper function that gets the machinelearning workspace.
// This function would fail the test if there is an error.
func GetMachinelearningWorkspace(t testing.TestingT, resGroupName string, workspaceName string, subscriptionID string) *machinelearning.Workspace {
	workspace, err := GetMachinelearningWorkspaceE(t, resGroupName, workspaceName, subscriptionID)
	require.NoError(t, err)

	return workspace
}

// GetMachinelearningWorkspaceE is a helper function that gets the machinelearning workspace.
func GetMachinelearningWorkspaceE(t testing.TestingT, resGroupName string, workspaceName string, subscriptionID string) (*machinelearning.Workspace, error) {
	// Create a ml workspace client
	workspaceClient, err := GetMachinelearningWorkspaceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding workspace client
	workspace, err := workspaceClient.Get(context.Background(), resGroupName, workspaceName, nil)
	if err != nil {
		return nil, err
	}

	return &workspace.Workspace, nil
}

// GetMachinelearningWorkspaceClientE is a helper function that will setup a machine learning workspace client.
func GetMachinelearningWorkspaceClientE(subscriptionID string) (*machinelearning.WorkspacesClient, error) {
	// Create a new Subnet client from client factory
	client, err := CreateMachinelearningWorkspaceClientE(subscriptionID, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// MachinelearningComputeExists indicates whether a machine learning compute clusterexists within a machine learning workspace; otherwise false
// This function would fail the test if there is an error.
func MachinelearningComputeExists(t testing.TestingT, resourceGroupName string, workspaceName string, computeName string, subscriptionID string) bool {
	result, err := MachinelearningComputeExistsE(t, resourceGroupName, workspaceName, computeName, subscriptionID)
	require.NoError(t, err)
	return result
}

// MachinelearningComputeExistsE indicates whether a machine learning compute clusterexists exists within a machine learning workspace
func MachinelearningComputeExistsE(t testing.TestingT, resourceGroupName, workspaceName, computeName, subscriptionID string) (bool, error) {
	compute, err := GetMachinelearningComputeE(t, resourceGroupName, workspaceName, computeName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return (computeName == *compute.Name), nil

}

// GetMachinelearningCompute is a helper function that gets the machinelearning compute.
// This function would fail the test if there is an error.
func GetMachinelearningCompute(t testing.TestingT, resGroupName string, workspaceName string, computeName string, subscriptionID string) *machinelearning.ComputeResource {
	compute, err := GetMachinelearningComputeE(t, resGroupName, workspaceName, computeName, subscriptionID)
	require.NoError(t, err)

	return compute
}

// GetMachinelearningComputeE is a helper function that gets the machinelearning workspace compute.
func GetMachinelearningComputeE(t testing.TestingT, resGroupName string, workspaceName string, computeName string, subscriptionID string) (*machinelearning.ComputeResource, error) {
	// Create a ml workspace client
	computeClient, err := GetMachinelearningComputeClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding workspace client
	compute, err := computeClient.Get(context.Background(), resGroupName, workspaceName, computeName, nil)
	if err != nil {
		return nil, err
	}

	return &compute.ComputeResource, nil
}

// GetMachinelearningComputeClientE is a helper function that will setup a machine learning compute client.
func GetMachinelearningComputeClientE(subscriptionID string) (*machinelearning.ComputeClient, error) {
	// Create a new Subnet client from client factory
	client, err := CreateMachinelearningComputeClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	return client, nil
}
