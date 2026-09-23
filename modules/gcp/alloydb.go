package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/alloydb/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetAlloyDBClusterAttrs returns the settings Google Cloud holds for the given AlloyDB cluster, so
// a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAlloyDBClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, clusterID string) *alloydb.Cluster {
	cluster, err := GetAlloyDBClusterAttrsE(t, ctx, projectID, region, clusterID)
	require.NoError(t, err)

	return cluster
}

// GetAlloyDBClusterAttrsE returns the settings Google Cloud holds for the given AlloyDB cluster.
// The ctx parameter supports cancellation and timeouts.
func GetAlloyDBClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, clusterID string) (*alloydb.Cluster, error) {
	logger.Default.Logf(t, "Getting settings for AlloyDB cluster %s in region %s in project %s", clusterID, region, projectID)

	service, err := NewAlloyDBServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAlloyDBClusterAttrsWithClient(ctx, service, projectID, region, clusterID)
}

// GetAlloyDBClusterAttrsWithClient returns the settings Google Cloud holds for the given AlloyDB
// cluster using the supplied *alloydb.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see alloydb_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAlloyDBClusterAttrsWithClient(ctx context.Context, service *alloydb.Service, projectID string, region string, clusterID string) (*alloydb.Cluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/clusters/%s", projectID, region, clusterID)

	cluster, err := service.Projects.Locations.Clusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the AlloyDB cluster %s does not exist in region %s in project %s", clusterID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for AlloyDB cluster %s in region %s in project %s: %w", clusterID, region, projectID, err)
	}

	return cluster, nil
}

// NewAlloyDBServiceE creates a AlloyDB service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewAlloyDBServiceE(t testing.TestingT, ctx context.Context) (*alloydb.Service, error) {
	return alloydb.NewService(ctx, append(withOptions(), option.WithScopes(alloydb.CloudPlatformScope))...)
}
