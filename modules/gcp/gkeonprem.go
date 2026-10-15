package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/gkeonprem/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetGKEOnPremBareMetalAdminClusterAttrs returns the settings Google Cloud holds for the bare metal admin cluster, so a test can assert on
// what was actually created rather than only that it exists.
// An admin cluster is registered rather than created by Google, so its state and the version it
// reports are the cluster's own answer and not an echo of the request.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremBareMetalAdminClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) *gkeonprem.BareMetalAdminCluster {
	result, err := GetGKEOnPremBareMetalAdminClusterAttrsE(t, ctx, projectID, location, clusterID)
	require.NoError(t, err)

	return result
}

// GetGKEOnPremBareMetalAdminClusterAttrsE returns the settings Google Cloud holds for the bare metal admin cluster.
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremBareMetalAdminClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) (*gkeonprem.BareMetalAdminCluster, error) {
	logger.Default.Logf(t, "Getting settings for bare metal admin cluster %s in %s in project %s", clusterID, location, projectID)

	service, err := NewGKEOnPremServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEOnPremBareMetalAdminClusterAttrsWithClient(ctx, service, projectID, location, clusterID)
}

// GetGKEOnPremBareMetalAdminClusterAttrsWithClient returns the settings Google Cloud holds for the bare metal admin cluster using the supplied
// *gkeonprem.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see gkeonprem_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremBareMetalAdminClusterAttrsWithClient(ctx context.Context, service *gkeonprem.Service, projectID string, location string, clusterID string) (*gkeonprem.BareMetalAdminCluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/bareMetalAdminClusters/%s", projectID, location, clusterID)

	result, err := service.Projects.Locations.BareMetalAdminClusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the bare metal admin cluster %s does not exist in %s in project %s", clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for bare metal admin cluster %s in %s in project %s: %w", clusterID, location, projectID, err)
	}

	return result, nil
}

// GetGKEOnPremBareMetalClusterAttrs returns the settings Google Cloud holds for the bare metal user cluster, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremBareMetalClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) *gkeonprem.BareMetalCluster {
	result, err := GetGKEOnPremBareMetalClusterAttrsE(t, ctx, projectID, location, clusterID)
	require.NoError(t, err)

	return result
}

// GetGKEOnPremBareMetalClusterAttrsE returns the settings Google Cloud holds for the bare metal user cluster.
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremBareMetalClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) (*gkeonprem.BareMetalCluster, error) {
	logger.Default.Logf(t, "Getting settings for bare metal user cluster %s in %s in project %s", clusterID, location, projectID)

	service, err := NewGKEOnPremServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEOnPremBareMetalClusterAttrsWithClient(ctx, service, projectID, location, clusterID)
}

// GetGKEOnPremBareMetalClusterAttrsWithClient returns the settings Google Cloud holds for the bare metal user cluster using the supplied
// *gkeonprem.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see gkeonprem_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremBareMetalClusterAttrsWithClient(ctx context.Context, service *gkeonprem.Service, projectID string, location string, clusterID string) (*gkeonprem.BareMetalCluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/bareMetalClusters/%s", projectID, location, clusterID)

	result, err := service.Projects.Locations.BareMetalClusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the bare metal user cluster %s does not exist in %s in project %s", clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for bare metal user cluster %s in %s in project %s: %w", clusterID, location, projectID, err)
	}

	return result, nil
}

// GetGKEOnPremBareMetalNodePoolAttrs returns the settings Google Cloud holds for the bare metal node pool, so a test can assert on
// what was actually created rather than only that it exists.
// A node pool hangs off a user cluster rather than off the project, so that cluster is named as well.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremBareMetalNodePoolAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, nodePoolID string) *gkeonprem.BareMetalNodePool {
	result, err := GetGKEOnPremBareMetalNodePoolAttrsE(t, ctx, projectID, location, clusterID, nodePoolID)
	require.NoError(t, err)

	return result
}

// GetGKEOnPremBareMetalNodePoolAttrsE returns the settings Google Cloud holds for the bare metal node pool.
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremBareMetalNodePoolAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, nodePoolID string) (*gkeonprem.BareMetalNodePool, error) {
	logger.Default.Logf(t, "Getting settings for bare metal node pool %s of cluster %s in %s in project %s", nodePoolID, clusterID, location, projectID)

	service, err := NewGKEOnPremServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEOnPremBareMetalNodePoolAttrsWithClient(ctx, service, projectID, location, clusterID, nodePoolID)
}

// GetGKEOnPremBareMetalNodePoolAttrsWithClient returns the settings Google Cloud holds for the bare metal node pool using the supplied
// *gkeonprem.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see gkeonprem_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremBareMetalNodePoolAttrsWithClient(ctx context.Context, service *gkeonprem.Service, projectID string, location string, clusterID string, nodePoolID string) (*gkeonprem.BareMetalNodePool, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/bareMetalClusters/%s/bareMetalNodePools/%s", projectID, location, clusterID, nodePoolID)

	result, err := service.Projects.Locations.BareMetalClusters.BareMetalNodePools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the bare metal node pool %s does not exist on cluster %s in %s in project %s", nodePoolID, clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for bare metal node pool %s of cluster %s in %s in project %s: %w", nodePoolID, clusterID, location, projectID, err)
	}

	return result, nil
}

// GetGKEOnPremVMwareAdminClusterAttrs returns the settings Google Cloud holds for the VMware admin cluster, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremVMwareAdminClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) *gkeonprem.VmwareAdminCluster {
	result, err := GetGKEOnPremVMwareAdminClusterAttrsE(t, ctx, projectID, location, clusterID)
	require.NoError(t, err)

	return result
}

// GetGKEOnPremVMwareAdminClusterAttrsE returns the settings Google Cloud holds for the VMware admin cluster.
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremVMwareAdminClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) (*gkeonprem.VmwareAdminCluster, error) {
	logger.Default.Logf(t, "Getting settings for VMware admin cluster %s in %s in project %s", clusterID, location, projectID)

	service, err := NewGKEOnPremServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEOnPremVMwareAdminClusterAttrsWithClient(ctx, service, projectID, location, clusterID)
}

// GetGKEOnPremVMwareAdminClusterAttrsWithClient returns the settings Google Cloud holds for the VMware admin cluster using the supplied
// *gkeonprem.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see gkeonprem_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremVMwareAdminClusterAttrsWithClient(ctx context.Context, service *gkeonprem.Service, projectID string, location string, clusterID string) (*gkeonprem.VmwareAdminCluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/vmwareAdminClusters/%s", projectID, location, clusterID)

	result, err := service.Projects.Locations.VmwareAdminClusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the VMware admin cluster %s does not exist in %s in project %s", clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for VMware admin cluster %s in %s in project %s: %w", clusterID, location, projectID, err)
	}

	return result, nil
}

// GetGKEOnPremVMwareClusterAttrs returns the settings Google Cloud holds for the VMware user cluster, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremVMwareClusterAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) *gkeonprem.VmwareCluster {
	result, err := GetGKEOnPremVMwareClusterAttrsE(t, ctx, projectID, location, clusterID)
	require.NoError(t, err)

	return result
}

// GetGKEOnPremVMwareClusterAttrsE returns the settings Google Cloud holds for the VMware user cluster.
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremVMwareClusterAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string) (*gkeonprem.VmwareCluster, error) {
	logger.Default.Logf(t, "Getting settings for VMware user cluster %s in %s in project %s", clusterID, location, projectID)

	service, err := NewGKEOnPremServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEOnPremVMwareClusterAttrsWithClient(ctx, service, projectID, location, clusterID)
}

// GetGKEOnPremVMwareClusterAttrsWithClient returns the settings Google Cloud holds for the VMware user cluster using the supplied
// *gkeonprem.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see gkeonprem_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremVMwareClusterAttrsWithClient(ctx context.Context, service *gkeonprem.Service, projectID string, location string, clusterID string) (*gkeonprem.VmwareCluster, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/vmwareClusters/%s", projectID, location, clusterID)

	result, err := service.Projects.Locations.VmwareClusters.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the VMware user cluster %s does not exist in %s in project %s", clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for VMware user cluster %s in %s in project %s: %w", clusterID, location, projectID, err)
	}

	return result, nil
}

// GetGKEOnPremVMwareNodePoolAttrs returns the settings Google Cloud holds for the VMware node pool, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremVMwareNodePoolAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, nodePoolID string) *gkeonprem.VmwareNodePool {
	result, err := GetGKEOnPremVMwareNodePoolAttrsE(t, ctx, projectID, location, clusterID, nodePoolID)
	require.NoError(t, err)

	return result
}

// GetGKEOnPremVMwareNodePoolAttrsE returns the settings Google Cloud holds for the VMware node pool.
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremVMwareNodePoolAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clusterID string, nodePoolID string) (*gkeonprem.VmwareNodePool, error) {
	logger.Default.Logf(t, "Getting settings for VMware node pool %s of cluster %s in %s in project %s", nodePoolID, clusterID, location, projectID)

	service, err := NewGKEOnPremServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEOnPremVMwareNodePoolAttrsWithClient(ctx, service, projectID, location, clusterID, nodePoolID)
}

// GetGKEOnPremVMwareNodePoolAttrsWithClient returns the settings Google Cloud holds for the VMware node pool using the supplied
// *gkeonprem.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see gkeonprem_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEOnPremVMwareNodePoolAttrsWithClient(ctx context.Context, service *gkeonprem.Service, projectID string, location string, clusterID string, nodePoolID string) (*gkeonprem.VmwareNodePool, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/vmwareClusters/%s/vmwareNodePools/%s", projectID, location, clusterID, nodePoolID)

	result, err := service.Projects.Locations.VmwareClusters.VmwareNodePools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the VMware node pool %s does not exist on cluster %s in %s in project %s", nodePoolID, clusterID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for VMware node pool %s of cluster %s in %s in project %s: %w", nodePoolID, clusterID, location, projectID, err)
	}

	return result, nil
}

// NewGKEOnPremServiceE creates an Anthos on premises service authenticated the same way every other
// client in this module is.
// The API is named gkeonprem rather than anthos, and it serves both the bare metal and the VMware
// clusters a project has registered.
// The ctx parameter supports cancellation and timeouts.
func NewGKEOnPremServiceE(t testing.TestingT, ctx context.Context) (*gkeonprem.Service, error) {
	return gkeonprem.NewService(ctx, append(withOptions(), option.WithScopes(gkeonprem.CloudPlatformScope))...)
}
