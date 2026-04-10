package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v2"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// LogAnalyticsWorkspaceExistsContext indicates whether the operational insights workspace exists.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func LogAnalyticsWorkspaceExistsContext(t testing.TestingT, ctx context.Context, workspaceName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := LogAnalyticsWorkspaceExistsContextE(ctx, workspaceName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// GetLogAnalyticsWorkspaceContext gets an operational insights workspace if it exists in a subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogAnalyticsWorkspaceContext(t testing.TestingT, ctx context.Context, workspaceName string, resourceGroupName string, subscriptionID string) *armoperationalinsights.Workspace {
	t.Helper()

	ws, err := GetLogAnalyticsWorkspaceContextE(ctx, workspaceName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return ws
}

// GetLogAnalyticsWorkspaceContextE gets an operational insights workspace if it exists in a subscription.
// The ctx parameter supports cancellation and timeouts.
func GetLogAnalyticsWorkspaceContextE(ctx context.Context, workspaceName, resoureGroupName, subscriptionID string) (*armoperationalinsights.Workspace, error) {
	client, err := CreateLogAnalyticsWorkspacesClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Get(ctx, resoureGroupName, workspaceName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Workspace, nil
}

// LogAnalyticsWorkspaceExistsContextE indicates whether the operational insights workspace exists and may return an error.
// The ctx parameter supports cancellation and timeouts.
func LogAnalyticsWorkspaceExistsContextE(ctx context.Context, workspaceName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetLogAnalyticsWorkspaceContextE(ctx, workspaceName, resourceGroupName, subscriptionID)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetLogAnalyticsWorkspacesClientE returns a workspaces client; otherwise error.
func GetLogAnalyticsWorkspacesClientE(subscriptionID string) (*armoperationalinsights.WorkspacesClient, error) {
	return CreateLogAnalyticsWorkspacesClientE(subscriptionID)
}
