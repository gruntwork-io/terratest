package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/datafactory/mgmt/2018-06-01/datafactory"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DataFactoryExists indicates whether the Data Factory exists for the subscription.
// This function would fail the test if there is an error.
func DataFactoryExists(t testing.TestingT, dataFactoryName string, resourceGroupName string, subscriptionID string) bool {
	exists, err := DataFactoryExistsE(dataFactoryName, resourceGroupName, subscriptionID)
	require.NoError(t, err)
	return exists
}

// DataFactoryExistsE indicates whether the specified Data Factory exists and may return an error.
func DataFactoryExistsE(dataFactoryName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetDataFactoryE(subscriptionID, resourceGroupName, dataFactoryName)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetDataFactory is a helper function that gets the synapse workspace.
// This function would fail the test if there is an error.
func GetDataFactory(t testing.TestingT, resGroupName string, factoryName string, subscriptionID string) *datafactory.Factory {
	Workspace, err := GetDataFactoryE(subscriptionID, resGroupName, factoryName)
	require.NoError(t, err)

	return Workspace
}

// GetDataFactoryE is a helper function that gets the workspace.
func GetDataFactoryE(subscriptionID string, resGroupName string, factoryName string) (*datafactory.Factory, error) {
	// Create a datafactory client
	datafactoryClient, err := CreateDataFactoriesClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding synapse workspace
	dataFactory, err := datafactoryClient.Get(context.Background(), resGroupName, factoryName, "")
	if err != nil {
		return nil, err
	}

	//Return synapse workspace
	return &dataFactory, nil
}
