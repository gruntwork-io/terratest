package azure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	azure "github.com/gruntwork-io/terratest/modules/azure"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods to create and delete log analytics resources are added, these tests can be extended.
*/

func TestLogAnalyticsWorkspace(t *testing.T) {
	t.Parallel()

	_, err := azure.LogAnalyticsWorkspaceExistsE("fake", "", "")
	assert.Error(t, err, "Workspace")
}

func TestGetLogAnalyticsWorkspaceE(t *testing.T) {
	t.Parallel()

	workspaceName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := azure.GetLogAnalyticsWorkspaceE(workspaceName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}
