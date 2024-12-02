package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/synapse/mgmt/2020-12-01/synapse"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetSynapseWorkspace is a helper function that gets the synapse workspace.
// This function would fail the test if there is an error.
func GetSynapseWorkspace(t testing.TestingT, resGroupName string, workspaceName string, subscriptionID string) *synapse.Workspace {
	Workspace, err := GetSynapseWorkspaceE(t, subscriptionID, resGroupName, workspaceName)
	require.NoError(t, err)

	return Workspace
}

// GetSynapseSqlPool is a helper function that gets the synapse workspace.
// This function would fail the test if there is an error.
func GetSynapseSqlPool(t testing.TestingT, resGroupName string, workspaceName string, sqlPoolName string, subscriptionID string) *synapse.SQLPool {
	SQLPool, err := GetSynapseSqlPoolE(t, subscriptionID, resGroupName, workspaceName, sqlPoolName)
	require.NoError(t, err)

	return SQLPool
}

// GetSynapseWorkspaceE is a helper function that gets the workspace.
func GetSynapseWorkspaceE(t testing.TestingT, subscriptionID string, resGroupName string, workspaceName string) (*synapse.Workspace, error) {
	// Create a synapse client
	synapseClient, err := CreateSynapseWorkspaceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding synapse workspace
	synapseWorkspace, err := synapseClient.Get(context.Background(), resGroupName, workspaceName)
	if err != nil {
		return nil, err
	}

	//Return synapse workspace
	return &synapseWorkspace, nil
}

// GetSynapseSqlPoolE is a helper function that gets the synapse sql pool.
func GetSynapseSqlPoolE(t testing.TestingT, subscriptionID string, resGroupName string, workspaceName string, sqlPoolName string) (*synapse.SQLPool, error) {
	// Create a synapse client
	synapseSqlPoolClient, err := CreateSynapseSqlPoolClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding synapse workspace
	synapseSqlPool, err := synapseSqlPoolClient.Get(context.Background(), resGroupName, workspaceName, sqlPoolName)
	if err != nil {
		return nil, err
	}

	//Return synapse workspace
	return &synapseSqlPool, nil
}
