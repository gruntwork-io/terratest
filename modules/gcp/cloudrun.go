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

	service, err := NewCloudRunClientE(t, ctx)
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

	service, err := NewCloudRunClientE(t, ctx)
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

// NewCloudRunClientE creates a Cloud Run client authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCloudRunClientE(t testing.TestingT, ctx context.Context) (*run.Service, error) {
	return run.NewService(ctx, append(withOptions(), option.WithScopes(run.CloudPlatformScope))...)
}

// GetCloudRunWorkerPoolAttrs returns the settings Google Cloud holds for the given Cloud Run worker
// pool, so a test can assert on what was actually created rather than only that it exists. A worker
// pool runs containers with no request path in front of them, so it has no URL and no traffic split.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunWorkerPoolAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string) *run.GoogleCloudRunV2WorkerPool {
	pool, err := GetCloudRunWorkerPoolAttrsE(t, ctx, projectID, location, poolID)
	require.NoError(t, err)

	return pool
}

// GetCloudRunWorkerPoolAttrsE returns the settings Google Cloud holds for the given worker pool.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunWorkerPoolAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string) (*run.GoogleCloudRunV2WorkerPool, error) {
	logger.Default.Logf(t, "Getting settings for worker pool %s in %s in project %s", poolID, location, projectID)

	service, err := NewCloudRunClientE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudRunWorkerPoolAttrsWithClient(ctx, service, projectID, location, poolID)
}

// GetCloudRunWorkerPoolAttrsWithClient returns the settings Google Cloud holds for the given worker
// pool using the supplied *run.Service. Prefer this variant in unit tests where the service is backed
// by an httptest fake server (see cloudrun_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunWorkerPoolAttrsWithClient(ctx context.Context, service *run.Service, projectID string, location string, poolID string) (*run.GoogleCloudRunV2WorkerPool, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/workerPools/%s", projectID, location, poolID)

	pool, err := service.Projects.Locations.WorkerPools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the worker pool %s does not exist in %s in project %s", poolID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for worker pool %s in %s in project %s: %w", poolID, location, projectID, err)
	}

	return pool, nil
}

// GetCloudRunServiceIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Cloud Run
// service, so a test can assert on who was actually granted access to it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunServiceIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string) *run.GoogleIamV1Policy {
	policy, err := GetCloudRunServiceIamPolicyAttrsE(t, ctx, projectID, location, serviceID)
	require.NoError(t, err)

	return policy
}

// GetCloudRunServiceIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given service.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunServiceIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string) (*run.GoogleIamV1Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for service %s in %s in project %s", serviceID, location, projectID)

	service, err := NewCloudRunClientE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudRunServiceIamPolicyAttrsWithClient(ctx, service, projectID, location, serviceID)
}

// GetCloudRunServiceIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given
// service using the supplied *run.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see cloudrun_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunServiceIamPolicyAttrsWithClient(ctx context.Context, service *run.Service, projectID string, location string, serviceID string) (*run.GoogleIamV1Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/services/%s", projectID, location, serviceID)

	policy, err := service.Projects.Locations.Services.GetIamPolicy(resource).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the service %s does not exist in %s in project %s", serviceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for service %s in %s in project %s: %w", serviceID, location, projectID, err)
	}

	return policy, nil
}

// GetCloudRunJobIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Cloud Run
// job, so a test can assert on who was actually granted access to it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunJobIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, jobID string) *run.GoogleIamV1Policy {
	policy, err := GetCloudRunJobIamPolicyAttrsE(t, ctx, projectID, location, jobID)
	require.NoError(t, err)

	return policy
}

// GetCloudRunJobIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given job.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunJobIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, jobID string) (*run.GoogleIamV1Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for job %s in %s in project %s", jobID, location, projectID)

	service, err := NewCloudRunClientE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudRunJobIamPolicyAttrsWithClient(ctx, service, projectID, location, jobID)
}

// GetCloudRunJobIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given
// job using the supplied *run.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see cloudrun_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunJobIamPolicyAttrsWithClient(ctx context.Context, service *run.Service, projectID string, location string, jobID string) (*run.GoogleIamV1Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectID, location, jobID)

	policy, err := service.Projects.Locations.Jobs.GetIamPolicy(resource).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the job %s does not exist in %s in project %s", jobID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for job %s in %s in project %s: %w", jobID, location, projectID, err)
	}

	return policy, nil
}

// GetCloudRunWorkerPoolIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Cloud Run
// worker pool, so a test can assert on who was actually granted access to it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunWorkerPoolIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string) *run.GoogleIamV1Policy {
	policy, err := GetCloudRunWorkerPoolIamPolicyAttrsE(t, ctx, projectID, location, poolID)
	require.NoError(t, err)

	return policy
}

// GetCloudRunWorkerPoolIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given worker pool.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunWorkerPoolIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string) (*run.GoogleIamV1Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for worker pool %s in %s in project %s", poolID, location, projectID)

	service, err := NewCloudRunClientE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudRunWorkerPoolIamPolicyAttrsWithClient(ctx, service, projectID, location, poolID)
}

// GetCloudRunWorkerPoolIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given
// worker pool using the supplied *run.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see cloudrun_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunWorkerPoolIamPolicyAttrsWithClient(ctx context.Context, service *run.Service, projectID string, location string, poolID string) (*run.GoogleIamV1Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/workerPools/%s", projectID, location, poolID)

	policy, err := service.Projects.Locations.WorkerPools.GetIamPolicy(resource).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the worker pool %s does not exist in %s in project %s", poolID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for worker pool %s in %s in project %s: %w", poolID, location, projectID, err)
	}

	return policy, nil
}
