package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v2"
)

// GetCloudRunServiceAttrs returns the settings Google Cloud holds for the given Cloud Run service,
// so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunServiceAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, serviceID string) *run.GoogleCloudRunV2Service {
	runService, err := GetCloudRunServiceAttrsE(t, ctx, projectID, region, serviceID)
	require.NoError(t, err)

	return runService
}

// GetCloudRunServiceAttrsE returns the settings Google Cloud holds for the given Cloud Run service.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunServiceAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, serviceID string) (*run.GoogleCloudRunV2Service, error) {
	logger.Default.Logf(t, "Getting settings for Cloud Run service %s in region %s in project %s", serviceID, region, projectID)

	service, err := NewCloudRunServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudRunServiceAttrsWithClient(ctx, service, projectID, region, serviceID)
}

// GetCloudRunServiceAttrsWithClient returns the settings Google Cloud holds for the given Cloud Run
// service using the supplied *run.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see cloudrun_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunServiceAttrsWithClient(ctx context.Context, service *run.Service, projectID string, region string, serviceID string) (*run.GoogleCloudRunV2Service, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/services/%s", projectID, region, serviceID)

	runService, err := service.Projects.Locations.Services.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Run service %s does not exist in region %s in project %s", serviceID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud Run service %s in region %s in project %s: %w", serviceID, region, projectID, err)
	}

	return runService, nil
}

// GetCloudRunJobAttrs returns the settings Google Cloud holds for the given Cloud Run job, so a
// test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunJobAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, jobID string) *run.GoogleCloudRunV2Job {
	job, err := GetCloudRunJobAttrsE(t, ctx, projectID, region, jobID)
	require.NoError(t, err)

	return job
}

// GetCloudRunJobAttrsE returns the settings Google Cloud holds for the given Cloud Run job.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunJobAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, jobID string) (*run.GoogleCloudRunV2Job, error) {
	logger.Default.Logf(t, "Getting settings for Cloud Run job %s in region %s in project %s", jobID, region, projectID)

	service, err := NewCloudRunServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudRunJobAttrsWithClient(ctx, service, projectID, region, jobID)
}

// GetCloudRunJobAttrsWithClient returns the settings Google Cloud holds for the given Cloud Run job
// using the supplied *run.Service. Prefer this variant in unit tests where the service is backed by
// an httptest fake server (see cloudrun_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunJobAttrsWithClient(ctx context.Context, service *run.Service, projectID string, region string, jobID string) (*run.GoogleCloudRunV2Job, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectID, region, jobID)

	job, err := service.Projects.Locations.Jobs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Run job %s does not exist in region %s in project %s", jobID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud Run job %s in region %s in project %s: %w", jobID, region, projectID, err)
	}

	return job, nil
}

// NewCloudRunServiceE creates a Cloud Run service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCloudRunServiceE(t testing.TestingT, ctx context.Context) (*run.Service, error) {
	return run.NewService(ctx, append(withOptions(), option.WithScopes(run.CloudPlatformScope))...)
}
