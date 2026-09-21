package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetWorkerPoolAttrs returns the settings Google Cloud holds for the given Cloud Build private
// worker pool, so a test can assert on what was actually created rather than only that it exists.
// A pool lives in a location, which has to be given.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkerPoolAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, poolName string) *cloudbuild.WorkerPool {
	pool, err := GetWorkerPoolAttrsE(t, ctx, projectID, location, poolName)
	require.NoError(t, err)

	return pool
}

// GetWorkerPoolAttrsE returns the settings Google Cloud holds for the given Cloud Build private
// worker pool.
// The ctx parameter supports cancellation and timeouts.
func GetWorkerPoolAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, poolName string) (*cloudbuild.WorkerPool, error) {
	logger.Default.Logf(t, "Getting settings for worker pool %s in project %s", poolName, projectID)

	service, err := NewCloudBuildRESTServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkerPoolAttrsWithClient(ctx, service, projectID, location, poolName)
}

// GetWorkerPoolAttrsWithClient returns the settings Google Cloud holds for the given Cloud Build
// private worker pool using the supplied *cloudbuild.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see workerpool_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkerPoolAttrsWithClient(ctx context.Context, service *cloudbuild.Service, projectID string, location string, poolName string) (*cloudbuild.WorkerPool, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/workerPools/%s", projectID, location, poolName)

	pool, err := service.Projects.Locations.WorkerPools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("worker pool %s does not exist in project %s", poolName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for worker pool %s in project %s: %w", poolName, projectID, err)
	}

	return pool, nil
}

// NewCloudBuildRESTServiceE creates a Cloud Build REST service authenticated the same way every
// other client in this module is. The build functions in cloudbuild.go use the gRPC client, which
// has no worker pool or trigger read that a local test server can stand in for; this is the REST
// client for the same API, and buildtrigger.go uses it too.
// The ctx parameter supports cancellation and timeouts.
func NewCloudBuildRESTServiceE(t testing.TestingT, ctx context.Context) (*cloudbuild.Service, error) {
	return cloudbuild.NewService(ctx, append(withOptions(), option.WithScopes(cloudbuild.CloudPlatformScope))...)
}
