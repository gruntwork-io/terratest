package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/cloudtasks/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetCloudTasksQueueAttrs returns the settings Google Cloud holds for the given Cloud Tasks queue, so a
// test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudTasksQueueAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, queueID string) *cloudtasks.Queue {
	queue, err := GetCloudTasksQueueAttrsE(t, ctx, projectID, region, queueID)
	require.NoError(t, err)

	return queue
}

// GetCloudTasksQueueAttrsE returns the settings Google Cloud holds for the given Cloud Tasks queue.
// The ctx parameter supports cancellation and timeouts.
func GetCloudTasksQueueAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, queueID string) (*cloudtasks.Queue, error) {
	logger.Default.Logf(t, "Getting settings for Cloud Tasks queue %s in region %s in project %s", queueID, region, projectID)

	service, err := NewCloudTasksServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudTasksQueueAttrsWithClient(ctx, service, projectID, region, queueID)
}

// GetCloudTasksQueueAttrsWithClient returns the settings Google Cloud holds for the given Cloud Tasks
// queue using the supplied *cloudtasks.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see cloudtasks_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudTasksQueueAttrsWithClient(ctx context.Context, service *cloudtasks.Service, projectID string, region string, queueID string) (*cloudtasks.Queue, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, region, queueID)

	queue, err := service.Projects.Locations.Queues.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Tasks queue %s does not exist in region %s in project %s", queueID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud Tasks queue %s in region %s in project %s: %w", queueID, region, projectID, err)
	}

	return queue, nil
}

// NewCloudTasksServiceE creates a Cloud Tasks service authenticated the same way every other client
// in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewCloudTasksServiceE(t testing.TestingT, ctx context.Context) (*cloudtasks.Service, error) {
	return cloudtasks.NewService(ctx, append(withOptions(), option.WithScopes(cloudtasks.CloudPlatformScope))...)
}
