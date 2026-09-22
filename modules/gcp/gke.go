package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetGKEClusterAttrs returns the settings Google Cloud holds for the given GKE cluster, so a test
// can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) *container.Cluster {
	cluster, err := GetGKEClusterAttrsE(t, ctx, projectID, location, clusterID)
	require.NoError(t, err)

	return cluster
}

// GetGKEClusterAttrsE returns the settings Google Cloud holds for the given GKE cluster.
// The ctx parameter supports cancellation and timeouts.
func GetGKEClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) (*container.Cluster, error) {
	logger.Default.Logf(t, "Getting settings for GKE cluster %s in location %s in project %s", clusterID, location, projectID)

	service, err := NewGKEServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEClusterAttrsWithClient(ctx, service, projectID, location, clusterID)
}

// GetGKEClusterAttrsWithClient returns the settings Google Cloud holds for the given GKE cluster
// using the supplied *container.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see gke_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEClusterAttrsWithClient(ctx context.Context, service *container.Service, projectID string, location string, clusterID string) (*container.Cluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, location, clusterID)

	cluster, err := service.Projects.Locations.Clusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the GKE cluster %s does not exist in location %s in project %s", clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for GKE cluster %s in location %s in project %s: %w", clusterID, location, projectID, err)
	}

	return cluster, nil
}

// GetGKENodePoolAttrs returns the settings Google Cloud holds for the given GKE node pool, so a
// test can assert on what was actually created rather than only that it exists. A pool is named by
// its cluster as well as its own id.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKENodePoolAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, nodePoolID string) *container.NodePool {
	nodePool, err := GetGKENodePoolAttrsE(t, ctx, projectID, location, clusterID, nodePoolID)
	require.NoError(t, err)

	return nodePool
}

// GetGKENodePoolAttrsE returns the settings Google Cloud holds for the given GKE node pool.
// The ctx parameter supports cancellation and timeouts.
func GetGKENodePoolAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, nodePoolID string) (*container.NodePool, error) {
	logger.Default.Logf(t, "Getting settings for GKE node pool %s in cluster %s in location %s in project %s", nodePoolID, clusterID, location, projectID)

	service, err := NewGKEServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKENodePoolAttrsWithClient(ctx, service, projectID, location, clusterID, nodePoolID)
}

// GetGKENodePoolAttrsWithClient returns the settings Google Cloud holds for the given GKE node pool
// using the supplied *container.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see gke_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKENodePoolAttrsWithClient(ctx context.Context, service *container.Service, projectID string, location string, clusterID string, nodePoolID string) (*container.NodePool, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s/nodePools/%s", projectID, location, clusterID, nodePoolID)

	nodePool, err := service.Projects.Locations.Clusters.NodePools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the GKE node pool %s does not exist in cluster %s in location %s in project %s", nodePoolID, clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for GKE node pool %s in cluster %s in location %s in project %s: %w", nodePoolID, clusterID, location, projectID, err)
	}

	return nodePool, nil
}

// NewGKEServiceE creates a GKE service authenticated the same way every other client in this module
// is.
// The ctx parameter supports cancellation and timeouts.
func NewGKEServiceE(t testing.TestingT, ctx context.Context) (*container.Service, error) {
	return container.NewService(ctx, append(withOptions(), option.WithScopes(container.CloudPlatformScope))...)
}
