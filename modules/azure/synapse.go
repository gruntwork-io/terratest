package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse"
	"github.com/gruntwork-io/terratest/modules/testing/v2"
	"github.com/stretchr/testify/require"
)

// GetSynapseWorkspace is a helper function that gets the synapse workspace.
// This function would fail the test if there is an error.
func GetSynapseWorkspace(t testing.TestingT, resGroupName string, workspaceName string, subscriptionID string) *armsynapse.Workspace {
	Workspace, err := GetSynapseWorkspaceE(t, subscriptionID, resGroupName, workspaceName)
	require.NoError(t, err)

	return Workspace
}

// GetSynapseSqlPool is a helper function that gets the synapse sql pool.
// This function would fail the test if there is an error.
func GetSynapseSqlPool(t testing.TestingT, resGroupName string, workspaceName string, sqlPoolName string, subscriptionID string) *armsynapse.SQLPool {
	SQLPool, err := GetSynapseSqlPoolE(t, subscriptionID, resGroupName, workspaceName, sqlPoolName)
	require.NoError(t, err)

	return SQLPool
}

// GetSynapseWorkspaceE is a helper function that gets the workspace.
func GetSynapseWorkspaceE(t testing.TestingT, subscriptionID string, resGroupName string, workspaceName string) (*armsynapse.Workspace, error) {
	// Create a synapse client
	synapseClient, err := CreateSynapseWorkspaceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding synapse workspace
	resp, err := synapseClient.Get(context.Background(), resGroupName, workspaceName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Workspace, nil
}

// GetSynapseSqlPoolE is a helper function that gets the synapse sql pool.
func GetSynapseSqlPoolE(t testing.TestingT, subscriptionID string, resGroupName string, workspaceName string, sqlPoolName string) (*armsynapse.SQLPool, error) {
	// Create a synapse sql pool client
	synapseSqlPoolClient, err := CreateSynapseSqlPoolClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding synapse sql pool
	resp, err := synapseSqlPoolClient.Get(context.Background(), resGroupName, workspaceName, sqlPoolName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.SQLPool, nil
}
