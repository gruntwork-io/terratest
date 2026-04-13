package azure

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/operationalinsights/mgmt/2020-03-01-preview/operationalinsights"
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

// LogAnalyticsWorkspaceExists indicates whether the operational insights workspace exists.
// This function would fail the test if there is an error.
//
// Deprecated: Use [LogAnalyticsWorkspaceExistsContext] instead.
func LogAnalyticsWorkspaceExists(t testing.TestingT, workspaceName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return LogAnalyticsWorkspaceExistsContext(t, context.Background(), workspaceName, resourceGroupName, subscriptionID) //nolint:staticcheck
}

// GetLogAnalyticsWorkspaceContext gets an operational insights workspace if it exists in a subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetLogAnalyticsWorkspaceContext(t testing.TestingT, ctx context.Context, workspaceName string, resourceGroupName string, subscriptionID string) *operationalinsights.Workspace {
	t.Helper()

	ws, err := GetLogAnalyticsWorkspaceContextE(ctx, workspaceName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return ws
}

// GetLogAnalyticsWorkspace gets an operational insights workspace if it exists in a subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetLogAnalyticsWorkspaceContext] instead.
func GetLogAnalyticsWorkspace(t testing.TestingT, workspaceName string, resourceGroupName string, subscriptionID string) *operationalinsights.Workspace {
	t.Helper()

	return GetLogAnalyticsWorkspaceContext(t, context.Background(), workspaceName, resourceGroupName, subscriptionID) //nolint:staticcheck
}

// GetLogAnalyticsWorkspaceContextE gets an operational insights workspace if it exists in a subscription.
// The ctx parameter supports cancellation and timeouts.
func GetLogAnalyticsWorkspaceContextE(ctx context.Context, workspaceName, resoureGroupName, subscriptionID string) (*operationalinsights.Workspace, error) {
	client, err := GetLogAnalyticsWorkspacesClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	ws, err := client.Get(ctx, resoureGroupName, workspaceName)
	if err != nil {
		return nil, err
	}

	return &ws, nil
}

// GetLogAnalyticsWorkspaceE gets an operational insights workspace if it exists in a subscription.
//
// Deprecated: Use [GetLogAnalyticsWorkspaceContextE] instead.
func GetLogAnalyticsWorkspaceE(workspaceName, resoureGroupName, subscriptionID string) (*operationalinsights.Workspace, error) {
	return GetLogAnalyticsWorkspaceContextE(context.Background(), workspaceName, resoureGroupName, subscriptionID)
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

// LogAnalyticsWorkspaceExistsE indicates whether the operational insights workspace exists and may return an error.
//
// Deprecated: Use [LogAnalyticsWorkspaceExistsContextE] instead.
func LogAnalyticsWorkspaceExistsE(workspaceName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return LogAnalyticsWorkspaceExistsContextE(context.Background(), workspaceName, resourceGroupName, subscriptionID)
}

// GetLogAnalyticsWorkspacesClientContextE returns a workspaces client; otherwise error.
// The ctx parameter supports cancellation and timeouts.
func GetLogAnalyticsWorkspacesClientContextE(_ context.Context, subscriptionID string) (*operationalinsights.WorkspacesClient, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		fmt.Println("Workspace client error getting subscription")

		return nil, err
	}

	client := operationalinsights.NewWorkspacesClient(subscriptionID)

	authorizer, err := NewAuthorizer()
	if err != nil {
		fmt.Println("authorizer error")

		return nil, err
	}

	client.Authorizer = *authorizer

	return &client, nil
}

// GetLogAnalyticsWorkspacesClientE returns a workspaces client; otherwise error.
//
// Deprecated: Use [GetLogAnalyticsWorkspacesClientContextE] instead.
func GetLogAnalyticsWorkspacesClientE(subscriptionID string) (*operationalinsights.WorkspacesClient, error) {
	return GetLogAnalyticsWorkspacesClientContextE(context.Background(), subscriptionID)
}
