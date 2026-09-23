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

// GetBigtableTableAttrs returns the settings Google Cloud holds for the given Bigtable table, so a test can assert
// on what was actually created rather than only that it exists. The full view is asked for, because the default
// answer carries no column families and a caller asserting on one would read an empty map.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigtableTableAttrs(t testing.TestingT, ctx context.Context, projectID string, instanceID string, tableID string) *bigtableadmin.Table {
	table, err := GetBigtableTableAttrsE(t, ctx, projectID, instanceID, tableID)
	require.NoError(t, err)

	return table
}

// GetBigtableTableAttrsE returns the settings Google Cloud holds for the given Bigtable table.
// The ctx parameter supports cancellation and timeouts.
func GetBigtableTableAttrsE(t testing.TestingT, ctx context.Context, projectID string, instanceID string, tableID string) (*bigtableadmin.Table, error) {
	logger.Default.Logf(t, "Getting settings for Bigtable table %s in instance %s in project %s", tableID, instanceID, projectID)

	service, err := NewBigtableAdminServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigtableTableAttrsWithClient(ctx, service, projectID, instanceID, tableID)
}

// GetBigtableTableAttrsWithClient returns the settings Google Cloud holds for the given Bigtable table using the
// supplied *bigtableadmin.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see bigtable_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigtableTableAttrsWithClient(ctx context.Context, service *bigtableadmin.Service, projectID string, instanceID string, tableID string) (*bigtableadmin.Table, error) {
	name := fmt.Sprintf("projects/%s/instances/%s/tables/%s", projectID, instanceID, tableID)

	table, err := service.Projects.Instances.Tables.Get(name).View("FULL").Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Bigtable table %s in instance %s in project %s does not exist", tableID, instanceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Bigtable table %s in instance %s in project %s: %w", tableID, instanceID, projectID, err)
	}

	return table, nil
}

// GetBigtableAppProfileAttrs returns the settings Google Cloud holds for the given Bigtable app profile, so a test can assert
// on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigtableAppProfileAttrs(t testing.TestingT, ctx context.Context, projectID string, instanceID string, profileID string) *bigtableadmin.AppProfile {
	profile, err := GetBigtableAppProfileAttrsE(t, ctx, projectID, instanceID, profileID)
	require.NoError(t, err)

	return profile
}

// GetBigtableAppProfileAttrsE returns the settings Google Cloud holds for the given Bigtable app profile.
// The ctx parameter supports cancellation and timeouts.
func GetBigtableAppProfileAttrsE(t testing.TestingT, ctx context.Context, projectID string, instanceID string, profileID string) (*bigtableadmin.AppProfile, error) {
	logger.Default.Logf(t, "Getting settings for Bigtable app profile %s in instance %s in project %s", profileID, instanceID, projectID)

	service, err := NewBigtableAdminServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigtableAppProfileAttrsWithClient(ctx, service, projectID, instanceID, profileID)
}

// GetBigtableAppProfileAttrsWithClient returns the settings Google Cloud holds for the given Bigtable app profile using the
// supplied *bigtableadmin.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see bigtable_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigtableAppProfileAttrsWithClient(ctx context.Context, service *bigtableadmin.Service, projectID string, instanceID string, profileID string) (*bigtableadmin.AppProfile, error) {
	name := fmt.Sprintf("projects/%s/instances/%s/appProfiles/%s", projectID, instanceID, profileID)

	profile, err := service.Projects.Instances.AppProfiles.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Bigtable app profile %s in instance %s in project %s does not exist", profileID, instanceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Bigtable app profile %s in instance %s in project %s: %w", profileID, instanceID, projectID, err)
	}

	return profile, nil
}

// GetBigtableAuthorizedViewAttrs returns the settings Google Cloud holds for the given Bigtable authorized view, so a test can assert
// on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigtableAuthorizedViewAttrs(t testing.TestingT, ctx context.Context, projectID string, instanceID string, tableID string, viewID string) *bigtableadmin.AuthorizedView {
	view, err := GetBigtableAuthorizedViewAttrsE(t, ctx, projectID, instanceID, tableID, viewID)
	require.NoError(t, err)

	return view
}

// GetBigtableAuthorizedViewAttrsE returns the settings Google Cloud holds for the given Bigtable authorized view.
// The ctx parameter supports cancellation and timeouts.
func GetBigtableAuthorizedViewAttrsE(t testing.TestingT, ctx context.Context, projectID string, instanceID string, tableID string, viewID string) (*bigtableadmin.AuthorizedView, error) {
	logger.Default.Logf(t, "Getting settings for Bigtable authorized view %s on table %s in instance %s in project %s", viewID, tableID, instanceID, projectID)

	service, err := NewBigtableAdminServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigtableAuthorizedViewAttrsWithClient(ctx, service, projectID, instanceID, tableID, viewID)
}

// GetBigtableAuthorizedViewAttrsWithClient returns the settings Google Cloud holds for the given Bigtable authorized view using the
// supplied *bigtableadmin.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see bigtable_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigtableAuthorizedViewAttrsWithClient(ctx context.Context, service *bigtableadmin.Service, projectID string, instanceID string, tableID string, viewID string) (*bigtableadmin.AuthorizedView, error) {
	name := fmt.Sprintf("projects/%s/instances/%s/tables/%s/authorizedViews/%s", projectID, instanceID, tableID, viewID)

	view, err := service.Projects.Instances.Tables.AuthorizedViews.Get(name).View("FULL").Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Bigtable authorized view %s on table %s in instance %s in project %s does not exist", viewID, tableID, instanceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Bigtable authorized view %s on table %s in instance %s in project %s: %w", viewID, tableID, instanceID, projectID, err)
	}

	return view, nil
}

// GetBigtableTableIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Bigtable
// table, so a test can assert on who may act on it. That is a different question from what the
// table holds.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigtableTableIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, instanceID string, tableID string) *bigtableadmin.Policy {
	policy, err := GetBigtableTableIamPolicyAttrsE(t, ctx, projectID, instanceID, tableID)
	require.NoError(t, err)

	return policy
}

// GetBigtableTableIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given Bigtable
// table.
// The ctx parameter supports cancellation and timeouts.
func GetBigtableTableIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, instanceID string, tableID string) (*bigtableadmin.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy of Bigtable table %s in instance %s in project %s", tableID, instanceID, projectID)

	service, err := NewBigtableAdminServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigtableTableIamPolicyAttrsWithClient(ctx, service, projectID, instanceID, tableID)
}

// GetBigtableTableIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given
// Bigtable table using the supplied *bigtableadmin.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see bigtable_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigtableTableIamPolicyAttrsWithClient(ctx context.Context, service *bigtableadmin.Service, projectID string, instanceID string, tableID string) (*bigtableadmin.Policy, error) {
	name := fmt.Sprintf("projects/%s/instances/%s/tables/%s", projectID, instanceID, tableID)

	// This call takes a request body even when nothing is being asked for beyond the policy.
	policy, err := service.Projects.Instances.Tables.GetIamPolicy(name, &bigtableadmin.GetIamPolicyRequest{}).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Bigtable table %s does not exist in instance %s in project %s", tableID, instanceID, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy of Bigtable table %s in instance %s in project %s: %w", tableID, instanceID, projectID, err)
	}

	return policy, nil
}
