package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure Synapse, these tests can be extended
*/

func TestGetSynapseWorkspaceE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	subscriptionID := ""
	workspaceName := ""

	_, err := GetSynapseWorkspaceE(t, subscriptionID, resGroupName, workspaceName)
	require.Error(t, err)
}

func TestGetSynapseSqlPoolE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	subscriptionID := ""
	workspaceName := ""
	sqlPoolName := ""

	_, err := GetSynapseSqlPoolE(t, subscriptionID, resGroupName, workspaceName, sqlPoolName)
	require.Error(t, err)
}
