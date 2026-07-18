package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/operationalinsights/armoperationalinsights/v2"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
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
func GetLogAnalyticsWorkspaceContextE(ctx context.Context, workspaceName, resourceGroupName, subscriptionID string) (*armoperationalinsights.Workspace, error) {
	client, err := CreateLogAnalyticsWorkspacesClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetLogAnalyticsWorkspaceWithClient(ctx, client, resourceGroupName, workspaceName)
}

// GetLogAnalyticsWorkspaceWithClient gets an operational insights workspace using the provided WorkspacesClient.
func GetLogAnalyticsWorkspaceWithClient(ctx context.Context, client *armoperationalinsights.WorkspacesClient, resourceGroupName string, workspaceName string) (*armoperationalinsights.Workspace, error) {
	resp, err := client.Get(ctx, resourceGroupName, workspaceName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Workspace, nil
}

// CreateLogAnalyticsWorkspacesClientContextE returns a workspaces client.
// The ctx parameter supports cancellation and timeouts.
func CreateLogAnalyticsWorkspacesClientContextE(_ context.Context, subscriptionID string) (*armoperationalinsights.WorkspacesClient, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	cred, err := newArmCredential()
	if err != nil {
		return nil, err
	}

	opts, err := newArmClientOptions()
	if err != nil {
		return nil, err
	}

	return armoperationalinsights.NewWorkspacesClient(subscriptionID, cred, opts)
}
