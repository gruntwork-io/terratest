//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/azure"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure Synapse, these tests can be extended
*/

func TestGetSynapseWorkspaceContextE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	subscriptionID := ""
	workspaceName := ""

	_, err := azure.GetSynapseWorkspaceContextE(t.Context(), subscriptionID, resGroupName, workspaceName)
	require.Error(t, err)
}

func TestGetSynapseSQLPoolContextE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	subscriptionID := ""
	workspaceName := ""
	sqlPoolName := ""

	_, err := azure.GetSynapseSQLPoolContextE(t.Context(), subscriptionID, resGroupName, workspaceName, sqlPoolName)
	require.Error(t, err)
}
