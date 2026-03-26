package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/synapse/armsynapse"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetSynapseWorkspaceContext retrieves the synapse workspace for the given subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSynapseWorkspaceContext(t testing.TestingT, ctx context.Context, resGroupName string, workspaceName string, subscriptionID string) *armsynapse.Workspace {
	t.Helper()

	workspace, err := GetSynapseWorkspaceContextE(ctx, subscriptionID, resGroupName, workspaceName)
	require.NoError(t, err)

	return workspace
}

// GetSynapseWorkspaceContextE retrieves the synapse workspace for the given subscription.
// The ctx parameter supports cancellation and timeouts.
func GetSynapseWorkspaceContextE(ctx context.Context, subscriptionID string, resGroupName string, workspaceName string) (*armsynapse.Workspace, error) {
	synapseClient, err := CreateSynapseWorkspaceClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := synapseClient.Get(ctx, resGroupName, workspaceName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Workspace, nil
}

// GetSynapseWorkspace retrieves the synapse workspace for the given subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetSynapseWorkspaceContext] instead.
func GetSynapseWorkspace(t testing.TestingT, resGroupName string, workspaceName string, subscriptionID string) *armsynapse.Workspace {
	t.Helper()

	return GetSynapseWorkspaceContext(t, context.Background(), resGroupName, workspaceName, subscriptionID)
}

// GetSynapseWorkspaceE retrieves the synapse workspace for the given subscription.
//
// Deprecated: Use [GetSynapseWorkspaceContextE] instead.
func GetSynapseWorkspaceE(t testing.TestingT, subscriptionID string, resGroupName string, workspaceName string) (*armsynapse.Workspace, error) { //nolint:unparam // t kept for API compatibility
	return GetSynapseWorkspaceContextE(context.Background(), subscriptionID, resGroupName, workspaceName)
}

// GetSynapseSQLPoolContext retrieves the synapse SQL pool for the given subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSynapseSQLPoolContext(t testing.TestingT, ctx context.Context, resGroupName string, workspaceName string, sqlPoolName string, subscriptionID string) *armsynapse.SQLPool {
	t.Helper()

	sqlPool, err := GetSynapseSQLPoolContextE(ctx, subscriptionID, resGroupName, workspaceName, sqlPoolName)
	require.NoError(t, err)

	return sqlPool
}

// GetSynapseSQLPoolContextE retrieves the synapse SQL pool for the given subscription.
// The ctx parameter supports cancellation and timeouts.
func GetSynapseSQLPoolContextE(ctx context.Context, subscriptionID string, resGroupName string, workspaceName string, sqlPoolName string) (*armsynapse.SQLPool, error) {
	synapseSQLPoolClient, err := CreateSynapseSqlPoolClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := synapseSQLPoolClient.Get(ctx, resGroupName, workspaceName, sqlPoolName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.SQLPool, nil
}

// GetSynapseSQLPool retrieves the synapse SQL pool for the given subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetSynapseSQLPoolContext] instead.
func GetSynapseSQLPool(t testing.TestingT, resGroupName string, workspaceName string, sqlPoolName string, subscriptionID string) *armsynapse.SQLPool {
	t.Helper()

	return GetSynapseSQLPoolContext(t, context.Background(), resGroupName, workspaceName, sqlPoolName, subscriptionID)
}

// GetSynapseSQLPoolE retrieves the synapse SQL pool for the given subscription.
//
// Deprecated: Use [GetSynapseSQLPoolContextE] instead.
func GetSynapseSQLPoolE(subscriptionID string, resGroupName string, workspaceName string, sqlPoolName string) (*armsynapse.SQLPool, error) {
	return GetSynapseSQLPoolContextE(context.Background(), subscriptionID, resGroupName, workspaceName, sqlPoolName)
}

// GetSynapseSqlPool retrieves the synapse SQL pool for the given subscription.
// This function would fail the test if there is an error.
//
//nolint:staticcheck,revive // Deprecated: Use [GetSynapseSQLPoolContext] instead.
func GetSynapseSqlPool(t testing.TestingT, resGroupName string, workspaceName string, sqlPoolName string, subscriptionID string) *armsynapse.SQLPool {
	t.Helper()

	return GetSynapseSQLPoolContext(t, context.Background(), resGroupName, workspaceName, sqlPoolName, subscriptionID)
}

// GetSynapseSqlPoolE retrieves the synapse SQL pool for the given subscription.
//
//nolint:staticcheck,revive // Deprecated: Use [GetSynapseSQLPoolContextE] instead.
func GetSynapseSqlPoolE(t testing.TestingT, subscriptionID string, resGroupName string, workspaceName string, sqlPoolName string) (*armsynapse.SQLPool, error) { //nolint:unparam // t kept for API compatibility
	return GetSynapseSQLPoolContextE(context.Background(), subscriptionID, resGroupName, workspaceName, sqlPoolName)
}
