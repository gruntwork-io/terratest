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
	"google.golang.org/api/storagebatchoperations/v1"
)

// GetStorageBatchOperationsJobAttrs returns the settings Google Cloud holds for the given Storage batch operations job, so a test can assert on what was
// actually created rather than only that it exists.
// A job rewrites or relabels many objects at once, so what it was told to do and how it ended are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetStorageBatchOperationsJobAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, jobID string) *storagebatchoperations.Job {
	attrs, err := GetStorageBatchOperationsJobAttrsE(t, ctx, projectID, location, jobID)
	require.NoError(t, err)

	return attrs
}

// GetStorageBatchOperationsJobAttrsE returns the settings Google Cloud holds for the given Storage batch operations job.
// The ctx parameter supports cancellation and timeouts.
func GetStorageBatchOperationsJobAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, jobID string) (*storagebatchoperations.Job, error) {
	logger.Default.Logf(t, "Getting settings for Storage batch operations job %s in %s in project %s", jobID, location, projectID)

	service, err := NewStorageBatchOperationsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetStorageBatchOperationsJobAttrsWithClient(ctx, service, projectID, location, jobID)
}

// GetStorageBatchOperationsJobAttrsWithClient returns the settings Google Cloud holds for the given Storage batch operations job using the supplied
// *storagebatchoperations.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see storagebatchoperations_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetStorageBatchOperationsJobAttrsWithClient(ctx context.Context, service *storagebatchoperations.Service, projectID string, location string, jobID string) (*storagebatchoperations.Job, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/jobs/%s", projectID, location, jobID)

	attrs, err := service.Projects.Locations.Jobs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Storage batch operations job %s in %s in project %s does not exist", jobID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Storage batch operations job %s in %s in project %s: %w", jobID, location, projectID, err)
	}

	return attrs, nil
}

// NewStorageBatchOperationsServiceE creates a Storage Batch Operations service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewStorageBatchOperationsServiceE(t testing.TestingT, ctx context.Context) (*storagebatchoperations.Service, error) {
	return storagebatchoperations.NewService(ctx, append(withOptions(), option.WithScopes(storagebatchoperations.CloudPlatformScope))...)
}
