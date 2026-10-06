package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/gkebackup/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetGKEBackupPlanAttrs returns the settings Google Cloud holds for the given backup plan, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEBackupPlanAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) *gkebackup.BackupPlan {
	plan, err := GetGKEBackupPlanAttrsE(t, ctx, projectID, location, planID)
	require.NoError(t, err)

	return plan
}

// GetGKEBackupPlanAttrsE returns the settings Google Cloud holds for the given backup plan.
// The ctx parameter supports cancellation and timeouts.
func GetGKEBackupPlanAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) (*gkebackup.BackupPlan, error) {
	logger.Default.Logf(t, "Getting settings for backup plan %s in %s in project %s", planID, location, projectID)

	service, err := NewGKEBackupServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEBackupPlanAttrsWithClient(ctx, service, projectID, location, planID)
}

// GetGKEBackupPlanAttrsWithClient returns the settings Google Cloud holds for the given backup plan using the supplied
// *gkebackup.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see gkebackup_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEBackupPlanAttrsWithClient(ctx context.Context, service *gkebackup.Service, projectID string, location string, planID string) (*gkebackup.BackupPlan, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backupPlans/%s", projectID, location, planID)

	plan, err := service.Projects.Locations.BackupPlans.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the backup plan %s does not exist in %s in project %s", planID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for backup plan %s in %s in project %s: %w", planID, location, projectID, err)
	}

	return plan, nil
}

// GetGKERestorePlanAttrs returns the settings Google Cloud holds for the given restore plan, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKERestorePlanAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) *gkebackup.RestorePlan {
	plan, err := GetGKERestorePlanAttrsE(t, ctx, projectID, location, planID)
	require.NoError(t, err)

	return plan
}

// GetGKERestorePlanAttrsE returns the settings Google Cloud holds for the given restore plan.
// The ctx parameter supports cancellation and timeouts.
func GetGKERestorePlanAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) (*gkebackup.RestorePlan, error) {
	logger.Default.Logf(t, "Getting settings for restore plan %s in %s in project %s", planID, location, projectID)

	service, err := NewGKEBackupServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKERestorePlanAttrsWithClient(ctx, service, projectID, location, planID)
}

// GetGKERestorePlanAttrsWithClient returns the settings Google Cloud holds for the given restore plan using the supplied
// *gkebackup.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see gkebackup_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKERestorePlanAttrsWithClient(ctx context.Context, service *gkebackup.Service, projectID string, location string, planID string) (*gkebackup.RestorePlan, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/restorePlans/%s", projectID, location, planID)

	plan, err := service.Projects.Locations.RestorePlans.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the restore plan %s does not exist in %s in project %s", planID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for restore plan %s in %s in project %s: %w", planID, location, projectID, err)
	}

	return plan, nil
}

// GetGKEBackupChannelAttrs returns the settings Google Cloud holds for the given GKE backup channel, so a test can assert on what was
// actually created rather than only that it exists.
// A channel is what lets a backup plan write into another project, so the project it names is the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEBackupChannelAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *gkebackup.BackupChannel {
	attrs, err := GetGKEBackupChannelAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetGKEBackupChannelAttrsE returns the settings Google Cloud holds for the given GKE backup channel.
// The ctx parameter supports cancellation and timeouts.
func GetGKEBackupChannelAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*gkebackup.BackupChannel, error) {
	logger.Default.Logf(t, "Getting settings for GKE backup channel %s in %s in project %s", id, location, projectID)

	service, err := NewGKEBackupServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEBackupChannelAttrsWithClient(ctx, service, projectID, location, id)
}

// GetGKEBackupChannelAttrsWithClient returns the settings Google Cloud holds for the given GKE backup channel using the supplied
// *gkebackup.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see gkebackup_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEBackupChannelAttrsWithClient(ctx context.Context, service *gkebackup.Service, projectID string, location string, id string) (*gkebackup.BackupChannel, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backupChannels/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.BackupChannels.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the GKE backup channel %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for GKE backup channel %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetGKERestoreChannelAttrs returns the settings Google Cloud holds for the given GKE restore channel, so a test can assert on what was
// actually created rather than only that it exists.
// A channel is what lets a restore plan read from another project, so the project it names is the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKERestoreChannelAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *gkebackup.RestoreChannel {
	attrs, err := GetGKERestoreChannelAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetGKERestoreChannelAttrsE returns the settings Google Cloud holds for the given GKE restore channel.
// The ctx parameter supports cancellation and timeouts.
func GetGKERestoreChannelAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*gkebackup.RestoreChannel, error) {
	logger.Default.Logf(t, "Getting settings for GKE restore channel %s in %s in project %s", id, location, projectID)

	service, err := NewGKEBackupServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKERestoreChannelAttrsWithClient(ctx, service, projectID, location, id)
}

// GetGKERestoreChannelAttrsWithClient returns the settings Google Cloud holds for the given GKE restore channel using the supplied
// *gkebackup.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see gkebackup_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKERestoreChannelAttrsWithClient(ctx context.Context, service *gkebackup.Service, projectID string, location string, id string) (*gkebackup.RestoreChannel, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/restoreChannels/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.RestoreChannels.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the GKE restore channel %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for GKE restore channel %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetGKEBackupPlanIamPolicyAttrs returns the IAM policy Google Cloud holds for the given GKE backup plan, so a test can assert on what was
// actually created rather than only that it exists.
// Who may run or read a backup is what the policy decides, and a plan Google created on its own carries no binding at all.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEBackupPlanIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) *gkebackup.Policy {
	policy, err := GetGKEBackupPlanIamPolicyAttrsE(t, ctx, projectID, location, planID)
	require.NoError(t, err)

	return policy
}

// GetGKEBackupPlanIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given GKE backup plan.
// The ctx parameter supports cancellation and timeouts.
func GetGKEBackupPlanIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) (*gkebackup.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for GKE backup plan %s in %s in project %s", planID, location, projectID)

	service, err := NewGKEBackupServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEBackupPlanIamPolicyAttrsWithClient(ctx, service, projectID, location, planID)
}

// GetGKEBackupPlanIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given GKE backup plan using the supplied
// *gkebackup.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see gkebackup_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEBackupPlanIamPolicyAttrsWithClient(ctx context.Context, service *gkebackup.Service, projectID string, location string, planID string) (*gkebackup.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/backupPlans/%s", projectID, location, planID)

	policy, err := service.Projects.Locations.BackupPlans.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the GKE backup plan %s in %s in project %s does not exist", planID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for GKE backup plan %s in %s in project %s: %w", planID, location, projectID, err)
	}

	return policy, nil
}

// GetGKERestorePlanIamPolicyAttrs returns the IAM policy Google Cloud holds for the given GKE restore plan, so a test can assert on what was
// actually created rather than only that it exists.
// Who may run a restore is what the policy decides, and a restore is the one operation that overwrites a running cluster.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKERestorePlanIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) *gkebackup.Policy {
	policy, err := GetGKERestorePlanIamPolicyAttrsE(t, ctx, projectID, location, planID)
	require.NoError(t, err)

	return policy
}

// GetGKERestorePlanIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given GKE restore plan.
// The ctx parameter supports cancellation and timeouts.
func GetGKERestorePlanIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, planID string) (*gkebackup.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for GKE restore plan %s in %s in project %s", planID, location, projectID)

	service, err := NewGKEBackupServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKERestorePlanIamPolicyAttrsWithClient(ctx, service, projectID, location, planID)
}

// GetGKERestorePlanIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given GKE restore plan using the supplied
// *gkebackup.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see gkebackup_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKERestorePlanIamPolicyAttrsWithClient(ctx context.Context, service *gkebackup.Service, projectID string, location string, planID string) (*gkebackup.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/restorePlans/%s", projectID, location, planID)

	policy, err := service.Projects.Locations.RestorePlans.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the GKE restore plan %s in %s in project %s does not exist", planID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for GKE restore plan %s in %s in project %s: %w", planID, location, projectID, err)
	}

	return policy, nil
}

// NewGKEBackupServiceE creates a Backup for GKE service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewGKEBackupServiceE(t testing.TestingT, ctx context.Context) (*gkebackup.Service, error) {
	return gkebackup.NewService(ctx, append(withOptions(), option.WithScopes(gkebackup.CloudPlatformScope))...)
}
