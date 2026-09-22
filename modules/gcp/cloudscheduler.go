package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/cloudscheduler/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetSchedulerJobAttrs returns the settings Google Cloud holds for the given Cloud Scheduler job,
// so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSchedulerJobAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, jobID string) *cloudscheduler.Job {
	job, err := GetSchedulerJobAttrsE(t, ctx, projectID, region, jobID)
	require.NoError(t, err)

	return job
}

// GetSchedulerJobAttrsE returns the settings Google Cloud holds for the given Cloud Scheduler job.
// The ctx parameter supports cancellation and timeouts.
func GetSchedulerJobAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, jobID string) (*cloudscheduler.Job, error) {
	logger.Default.Logf(t, "Getting settings for Cloud Scheduler job %s in region %s in project %s", jobID, region, projectID)

	service, err := NewCloudSchedulerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSchedulerJobAttrsWithClient(ctx, service, projectID, region, jobID)
}

// GetSchedulerJobAttrsWithClient returns the settings Google Cloud holds for the given Cloud
// Scheduler job using the supplied *cloudscheduler.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see cloudscheduler_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSchedulerJobAttrsWithClient(ctx context.Context, service *cloudscheduler.Service, projectID string, region string, jobID string) (*cloudscheduler.Job, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectID, region, jobID)

	job, err := service.Projects.Locations.Jobs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Scheduler job %s does not exist in region %s in project %s", jobID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud Scheduler job %s in region %s in project %s: %w", jobID, region, projectID, err)
	}

	return job, nil
}

// NewCloudSchedulerServiceE creates a Cloud Scheduler service authenticated the same way every
// other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCloudSchedulerServiceE(t testing.TestingT, ctx context.Context) (*cloudscheduler.Service, error) {
	return cloudscheduler.NewService(ctx, append(withOptions(), option.WithScopes(cloudscheduler.CloudPlatformScope))...)
}
