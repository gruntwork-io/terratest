package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when methods to create and delete log analytics resources are added, these tests can be extended.
*/

func TestLogAnalyticsWorkspace(t *testing.T) {
	t.Parallel()

	_, err := LogAnalyticsWorkspaceExistsE("fake", "", "")
	assert.Error(t, err, "Workspace")
}

func TestGetLogAnalyticsWorkspaceE(t *testing.T) {
	t.Parallel()
	workspaceName := ""
	resourceGroupName := ""
	subscriptionID := ""

	_, err := GetLogAnalyticsWorkspaceE(workspaceName, resourceGroupName, subscriptionID)
	require.Error(t, err)
}
