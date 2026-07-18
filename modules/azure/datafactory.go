package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v9"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// DataFactoryExistsContext indicates whether the Data Factory exists for the subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DataFactoryExistsContext(t testing.TestingT, ctx context.Context, dataFactoryName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	exists, err := DataFactoryExistsContextE(ctx, dataFactoryName, resourceGroupName, subscriptionID)
	require.NoError(t, err)

	return exists
}

// DataFactoryExistsContextE indicates whether the specified Data Factory exists and may return an error.
// The ctx parameter supports cancellation and timeouts.
func DataFactoryExistsContextE(ctx context.Context, dataFactoryName string, resourceGroupName string, subscriptionID string) (bool, error) {
	_, err := GetDataFactoryContextE(ctx, subscriptionID, resourceGroupName, dataFactoryName)
	if err != nil {
		if ResourceNotFoundErrorExists(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// GetDataFactoryContext returns the Data Factory object.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataFactoryContext(t testing.TestingT, ctx context.Context, subscriptionID string, resGroupName string, factoryName string) *armdatafactory.Factory {
	t.Helper()

	factory, err := GetDataFactoryContextE(ctx, subscriptionID, resGroupName, factoryName)
	require.NoError(t, err)

	return factory
}

// GetDataFactoryContextE returns the Data Factory object.
// The ctx parameter supports cancellation and timeouts.
func GetDataFactoryContextE(ctx context.Context, subscriptionID string, resGroupName string, factoryName string) (*armdatafactory.Factory, error) {

	datafactoryClient, err := CreateDataFactoriesClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetDataFactoryWithClient(ctx, datafactoryClient, resGroupName, factoryName)
}

// GetDataFactoryWithClient returns a Data Factory using the provided FactoriesClient.
// This variant is useful for testing with fake clients.
func GetDataFactoryWithClient(ctx context.Context, client *armdatafactory.FactoriesClient, resGroupName string, factoryName string) (*armdatafactory.Factory, error) {
	resp, err := client.Get(ctx, resGroupName, factoryName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.Factory, nil
}
