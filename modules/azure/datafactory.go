package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v9"
	"github.com/gruntwork-io/terratest/modules/testing"
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

// DataFactoryExists indicates whether the Data Factory exists for the subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [DataFactoryExistsContext] instead.
func DataFactoryExists(t testing.TestingT, dataFactoryName string, resourceGroupName string, subscriptionID string) bool {
	t.Helper()

	return DataFactoryExistsContext(t, context.Background(), dataFactoryName, resourceGroupName, subscriptionID) //nolint:staticcheck
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

// DataFactoryExistsE indicates whether the specified Data Factory exists and may return an error.
//
// Deprecated: Use [DataFactoryExistsContextE] instead.
func DataFactoryExistsE(dataFactoryName string, resourceGroupName string, subscriptionID string) (bool, error) {
	return DataFactoryExistsContextE(context.Background(), dataFactoryName, resourceGroupName, subscriptionID)
}

// GetDataFactoryContext returns the Data Factory object.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataFactoryContext(t testing.TestingT, ctx context.Context, resGroupName string, factoryName string, subscriptionID string) *armdatafactory.Factory {
	t.Helper()

	factory, err := GetDataFactoryContextE(ctx, subscriptionID, resGroupName, factoryName)
	require.NoError(t, err)

	return factory
}

// GetDataFactory is a helper function that gets the data factory.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetDataFactoryContext] instead.
func GetDataFactory(t testing.TestingT, resGroupName string, factoryName string, subscriptionID string) *armdatafactory.Factory {
	t.Helper()

	return GetDataFactoryContext(t, context.Background(), resGroupName, factoryName, subscriptionID) //nolint:staticcheck
}

// GetDataFactoryContextE returns the Data Factory object.
// The ctx parameter supports cancellation and timeouts.
func GetDataFactoryContextE(ctx context.Context, subscriptionID string, resGroupName string, factoryName string) (*armdatafactory.Factory, error) {
	// Create a datafactory client
	datafactoryClient, err := CreateDataFactoriesClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get the corresponding data factory
	resp, err := datafactoryClient.Get(ctx, resGroupName, factoryName, nil)
	if err != nil {
		return nil, err
	}

	// Return data factory
	return &resp.Factory, nil
}

// GetDataFactoryE is a helper function that gets the data factory.
//
// Deprecated: Use [GetDataFactoryContextE] instead.
func GetDataFactoryE(subscriptionID string, resGroupName string, factoryName string) (*armdatafactory.Factory, error) {
	return GetDataFactoryContextE(context.Background(), subscriptionID, resGroupName, factoryName)
}
