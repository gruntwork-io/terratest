package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// GetSynapseWorkspaceContext retrieves the synapse workspace for the given subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSynapseWorkspaceContext(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, workspaceName string) *armsynapse.Workspace {
	t.Helper()

	workspace, err := GetSynapseWorkspaceContextE(ctx, subscriptionID, resGroupName, workspaceName)
	require.NoError(t, err)

	return workspace
}

// GetSynapseWorkspaceContextE retrieves the synapse workspace for the given subscription.
// The ctx parameter supports cancellation and timeouts.
func GetSynapseWorkspaceContextE(ctx context.Context, subscriptionID string, resGroupName string, workspaceName string) (*armsynapse.Workspace, error) {
	synapseClient, err := CreateSynapseWorkspaceClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetSynapseWorkspaceWithClient(ctx, synapseClient, resGroupName, workspaceName)
}

// GetSynapseWorkspaceWithClient retrieves the synapse workspace using the provided WorkspacesClient.
func GetSynapseWorkspaceWithClient(ctx context.Context, client *armsynapse.WorkspacesClient, resGroupName string, workspaceName string) (*armsynapse.Workspace, error) {
	resp, err := client.Get(ctx, resGroupName, workspaceName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Workspace, nil
}

// GetSynapseSQLPoolContext retrieves the synapse SQL pool for the given subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSynapseSQLPoolContext(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, workspaceName string, sqlPoolName string) *armsynapse.SQLPool {
	t.Helper()

	sqlPool, err := GetSynapseSQLPoolContextE(ctx, subscriptionID, resGroupName, workspaceName, sqlPoolName)
	require.NoError(t, err)

	return sqlPool
}

// GetSynapseSQLPoolContextE retrieves the synapse SQL pool for the given subscription.
// The ctx parameter supports cancellation and timeouts.
func GetSynapseSQLPoolContextE(ctx context.Context, subscriptionID string, resGroupName string, workspaceName string, sqlPoolName string) (*armsynapse.SQLPool, error) {
	synapseSQLPoolClient, err := CreateSynapseSQLPoolClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetSynapseSQLPoolWithClient(ctx, synapseSQLPoolClient, resGroupName, workspaceName, sqlPoolName)
}

// GetSynapseSQLPoolWithClient retrieves the synapse SQL pool using the provided SQLPoolsClient.
func GetSynapseSQLPoolWithClient(ctx context.Context, client *armsynapse.SQLPoolsClient, resGroupName string, workspaceName string, sqlPoolName string) (*armsynapse.SQLPool, error) {
	resp, err := client.Get(ctx, resGroupName, workspaceName, sqlPoolName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.SQLPool, nil
}
