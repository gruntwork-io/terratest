package azure

import (
	"testing"

	"github.com/stretchr/testify/require"
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

	exists, err := DataFactoryExistsE(dataFactoryName, resourceGroupName, subscriptionID)

	require.False(t, exists)
	require.Error(t, err)
}

func TestGetDataFactoryE(t *testing.T) {
	t.Parallel()

	resGroupName := ""
	subscriptionID := ""
	dataFactoryName := ""

	_, err := GetDataFactoryE(subscriptionID, resGroupName, dataFactoryName)
	require.Error(t, err)
}
