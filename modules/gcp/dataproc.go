package gcp

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDataprocClusterAttrs returns the settings Google Cloud holds for the given Dataproc cluster, so
// a test can assert on what was actually created rather than only that it exists. A cluster is read
// from the region it was created in, which is part of its address rather than a filter.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, clusterID string) *dataproc.Cluster {
	cluster, err := GetDataprocClusterAttrsE(t, ctx, projectID, region, clusterID)
	require.NoError(t, err)

	return cluster
}

// GetDataprocClusterAttrsE returns the settings Google Cloud holds for the given Dataproc cluster.
// The ctx parameter supports cancellation and timeouts.
func GetDataprocClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, clusterID string) (*dataproc.Cluster, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc cluster %s in region %s in project %s", clusterID, region, projectID)

	service, err := NewDataprocServiceE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	return GetDataprocClusterAttrsWithClient(ctx, service, projectID, region, clusterID)
}

// GetDataprocClusterAttrsWithClient returns the settings Google Cloud holds for the given Dataproc
// cluster using the supplied *dataproc.Service, which has to point at that region's host. Prefer
// this variant in unit tests where the service is backed by an httptest fake server (see
// dataproc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDataprocClusterAttrsWithClient(ctx context.Context, service *dataproc.Service, projectID string, region string, clusterID string) (*dataproc.Cluster, error) {
	cluster, err := service.Projects.Regions.Clusters.Get(projectID, region, clusterID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc cluster %s does not exist in region %s in project %s", clusterID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc cluster %s in region %s in project %s: %w", clusterID, region, projectID, err)
	}

	return cluster, nil
}

// dataprocRegionPattern is what a Google Cloud region identifier may contain. It is checked before
// a region reaches an endpoint, because a value carrying a slash, a colon or an at sign would build
// a URL pointing at a host of the caller's choosing rather than at Google.
var dataprocRegionPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// NewDataprocServiceE creates a Dataproc service authenticated the same way every other client in
// this module is. Dataproc answers on a per-region host for every region but `global`, so the
// region the caller is asking about decides which one this talks to, and a region that is not a
// plain identifier is refused.
// The ctx parameter supports cancellation and timeouts.
func NewDataprocServiceE(t testing.TestingT, ctx context.Context, region string) (*dataproc.Service, error) {
	if !dataprocRegionPattern.MatchString(region) {
		return nil, fmt.Errorf("%q is not a valid region: a region may hold only lowercase letters, digits and hyphens", region)
	}

	opts := append(withOptions(), option.WithScopes(dataproc.CloudPlatformScope))
	if region != "global" {
		opts = append(opts, option.WithEndpoint(fmt.Sprintf("https://%s-dataproc.googleapis.com/", region)))
	}

	return dataproc.NewService(ctx, opts...)
}
