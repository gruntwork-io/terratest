package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerservice/armcontainerservice/v6"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetManagedClusterContext returns a ManagedCluster for the specified cluster in the given resource group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetManagedClusterContext(t testing.TestingT, ctx context.Context, resourceGroupName, clusterName, subscriptionID string) *armcontainerservice.ManagedCluster {
	t.Helper()

	cluster, err := GetManagedClusterContextE(t, ctx, resourceGroupName, clusterName, subscriptionID)
	require.NoError(t, err)

	return cluster
}

// GetManagedClusterE returns a ManagedCluster for the specified cluster in the given resource group.
//
// Deprecated: Use [GetManagedClusterContextE] instead.
func GetManagedClusterE(t testing.TestingT, resourceGroupName, clusterName, subscriptionID string) (*armcontainerservice.ManagedCluster, error) {
	return GetManagedClusterContextE(t, context.Background(), resourceGroupName, clusterName, subscriptionID)
}

// GetManagedClusterContextE returns a ManagedCluster for the specified cluster in the given resource group.
// The ctx parameter supports cancellation and timeouts.
func GetManagedClusterContextE(t testing.TestingT, ctx context.Context, resourceGroupName, clusterName, subscriptionID string) (*armcontainerservice.ManagedCluster, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	client, err := CreateManagedClustersClientContextE(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	return GetManagedClusterWithClient(ctx, client, resourceGroupName, clusterName)
}

// GetManagedClusterWithClient returns a ManagedCluster using the provided ManagedClustersClient.
func GetManagedClusterWithClient(ctx context.Context, client *armcontainerservice.ManagedClustersClient, resourceGroupName string, clusterName string) (*armcontainerservice.ManagedCluster, error) {
	resp, err := client.Get(ctx, resourceGroupName, clusterName, nil)
	if err != nil {
		return nil, err
	}

	return &resp.ManagedCluster, nil
}
