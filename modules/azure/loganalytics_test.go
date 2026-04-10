//go:build azure
// +build azure

package azure_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods to create and delete log analytics resources are added, these tests can be extended.
*/

func TestLogAnalyticsWorkspaceExistsContextE(t *testing.T) {
	t.Parallel()

	_, err := azure.LogAnalyticsWorkspaceExistsContextE(t.Context(), "fake", "", "")
	require.Error(t, err)
}

func TestGetLogAnalyticsWorkspaceContextE(t *testing.T) {
	t.Parallel()

	workspaceName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := azure.GetLogAnalyticsWorkspaceContextE(t.Context(), workspaceName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}
