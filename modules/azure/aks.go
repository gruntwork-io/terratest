package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/services/containerservice/mgmt/2019-11-01/containerservice"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetManagedClustersClientE is a helper function that will setup an Azure ManagedClusters client on your behalf.
func GetManagedClustersClientE(subscriptionID string) (*containerservice.ManagedClustersClient, error) {
	// Create a cluster client
	client, err := CreateManagedClustersClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	// setup authorizer
	authorizer, err := NewAuthorizer()
	if err != nil {
		return nil, err
	}

	client.Authorizer = *authorizer

	return &client, nil
}

// GetManagedClusterContextE returns a ManagedCluster for the specified cluster in the given resource group.
// The ctx parameter supports cancellation and timeouts.
func GetManagedClusterContextE(t testing.TestingT, ctx context.Context, resourceGroupName, clusterName, subscriptionID string) (*containerservice.ManagedCluster, error) {
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	client, err := GetManagedClustersClientE(subscriptionID)
	if err != nil {
		return nil, err
	}

	managedCluster, err := client.Get(ctx, resourceGroupName, clusterName)
	if err != nil {
		return nil, err
	}

	return &managedCluster, nil
}

// GetManagedClusterContext returns a ManagedCluster for the specified cluster in the given resource group.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetManagedClusterContext(t testing.TestingT, ctx context.Context, resourceGroupName, clusterName, subscriptionID string) *containerservice.ManagedCluster {
	t.Helper()

	cluster, err := GetManagedClusterContextE(t, ctx, resourceGroupName, clusterName, subscriptionID)
	require.NoError(t, err)

	return cluster
}

// GetManagedClusterE returns a ManagedCluster for the specified cluster in the given resource group.
//
// Deprecated: Use [GetManagedClusterContextE] instead.
func GetManagedClusterE(t testing.TestingT, resourceGroupName, clusterName, subscriptionID string) (*containerservice.ManagedCluster, error) {
	return GetManagedClusterContextE(t, context.Background(), resourceGroupName, clusterName, subscriptionID)
}
