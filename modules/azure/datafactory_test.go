package azure_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	azure "github.com/gruntwork-io/terratest/modules/azure"
)

/*
The below tests are currently stubbed out, with the expectation that they will throw errors.
If/when CRUD methods are introduced for Azure Synapse, these tests can be extended
*/
func TestDataFactoryExists(t *testing.T) {
	t.Parallel()

	dataFactoryName := ""
	resourceGroupName := ""
	subscriptionID := ""

	exists, err := azure.DataFactoryExistsContextE(t.Context(), dataFactoryName, resourceGroupName, subscriptionID)

	require.False(t, exists)
	require.Error(t, err)
}

func TestGetDataFactoryContextE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	subscriptionID := ""
	dataFactoryName := ""

	_, err := azure.GetDataFactoryContextE(t.Context(), subscriptionID, resGroupName, dataFactoryName)
	require.Error(t, err)
}
