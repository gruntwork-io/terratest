package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/bigtableadmin/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetBigtableInstanceAttrs returns the settings Google Cloud holds for the given Bigtable instance,
// so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigtableInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, instanceID string) *bigtableadmin.Instance {
	instance, err := GetBigtableInstanceAttrsE(t, ctx, projectID, instanceID)
	require.NoError(t, err)

	return instance
}

// GetBigtableInstanceAttrsE returns the settings Google Cloud holds for the given Bigtable
// instance.
// The ctx parameter supports cancellation and timeouts.
func GetBigtableInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, instanceID string) (*bigtableadmin.Instance, error) {
	logger.Default.Logf(t, "Getting settings for Bigtable instance %s in project %s", instanceID, projectID)

	service, err := NewBigtableAdminServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigtableInstanceAttrsWithClient(ctx, service, projectID, instanceID)
}

// GetBigtableInstanceAttrsWithClient returns the settings Google Cloud holds for the given Bigtable
// instance using the supplied *bigtableadmin.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see bigtable_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigtableInstanceAttrsWithClient(ctx context.Context, service *bigtableadmin.Service, projectID string, instanceID string) (*bigtableadmin.Instance, error) {
	name := fmt.Sprintf("projects/%s/instances/%s", projectID, instanceID)

	instance, err := service.Projects.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Bigtable instance %s does not exist in project %s", instanceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Bigtable instance %s in project %s: %w", instanceID, projectID, err)
	}

	return instance, nil
}

// GetBigtableClusterAttrs returns the settings Google Cloud holds for the given Bigtable cluster,
// so a test can assert on what was actually created rather than only that it exists. An instance
// carries no clusters in the response Google returns for it, so a cluster is read on its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigtableClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, instanceID string, clusterID string) *bigtableadmin.Cluster {
	cluster, err := GetBigtableClusterAttrsE(t, ctx, projectID, instanceID, clusterID)
	require.NoError(t, err)

	return cluster
}

// GetBigtableClusterAttrsE returns the settings Google Cloud holds for the given Bigtable cluster.
// The ctx parameter supports cancellation and timeouts.
func GetBigtableClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, instanceID string, clusterID string) (*bigtableadmin.Cluster, error) {
	logger.Default.Logf(t, "Getting settings for Bigtable cluster %s in instance %s in project %s", clusterID, instanceID, projectID)

	service, err := NewBigtableAdminServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigtableClusterAttrsWithClient(ctx, service, projectID, instanceID, clusterID)
}

// GetBigtableClusterAttrsWithClient returns the settings Google Cloud holds for the given Bigtable
// cluster using the supplied *bigtableadmin.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see bigtable_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigtableClusterAttrsWithClient(ctx context.Context, service *bigtableadmin.Service, projectID string, instanceID string, clusterID string) (*bigtableadmin.Cluster, error) {
	name := fmt.Sprintf("projects/%s/instances/%s/clusters/%s", projectID, instanceID, clusterID)

	cluster, err := service.Projects.Instances.Clusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Bigtable cluster %s does not exist in instance %s in project %s", clusterID, instanceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Bigtable cluster %s in instance %s in project %s: %w", clusterID, instanceID, projectID, err)
	}

	return cluster, nil
}

// NewBigtableAdminServiceE creates a Bigtable service authenticated the same way every other client
// in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewBigtableAdminServiceE(t testing.TestingT, ctx context.Context) (*bigtableadmin.Service, error) {
	return bigtableadmin.NewService(ctx, append(withOptions(), option.WithScopes(bigtableadmin.CloudPlatformScope))...)
}
