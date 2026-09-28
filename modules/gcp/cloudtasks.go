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
	queue, err := GetTasksQueueAttrsE(t, ctx, projectID, region, queueID)
	require.NoError(t, err)

	return queue
}

// GetTasksQueueAttrsE returns the settings Google Cloud holds for the given Cloud Tasks queue.
// The ctx parameter supports cancellation and timeouts.
func GetTasksQueueAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, queueID string) (*cloudtasks.Queue, error) {
	logger.Default.Logf(t, "Getting settings for Cloud Tasks queue %s in region %s in project %s", queueID, region, projectID)

	service, err := NewCloudTasksServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetTasksQueueAttrsWithClient(ctx, service, projectID, region, queueID)
}

// GetTasksQueueAttrsWithClient returns the settings Google Cloud holds for the given Cloud Tasks
// queue using the supplied *cloudtasks.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see cloudtasks_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetTasksQueueAttrsWithClient(ctx context.Context, service *cloudtasks.Service, projectID string, region string, queueID string) (*cloudtasks.Queue, error) {
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

// GetCloudTasksQueueIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Cloud
// Tasks queue, so a test can assert on who was actually granted access to it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudTasksQueueIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, queueID string) *cloudtasks.Policy {
	policy, err := GetCloudTasksQueueIamPolicyAttrsE(t, ctx, projectID, location, queueID)
	require.NoError(t, err)

	return policy
}

// GetCloudTasksQueueIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given queue.
// The ctx parameter supports cancellation and timeouts.
func GetCloudTasksQueueIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, queueID string) (*cloudtasks.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for queue %s in %s in project %s", queueID, location, projectID)

	service, err := NewCloudTasksServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudTasksQueueIamPolicyAttrsWithClient(ctx, service, projectID, location, queueID)
}

// GetCloudTasksQueueIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given
// queue using the supplied *cloudtasks.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see cloudtasks_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudTasksQueueIamPolicyAttrsWithClient(ctx context.Context, service *cloudtasks.Service, projectID string, location string, queueID string) (*cloudtasks.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, location, queueID)

	// This call takes a request body rather than a plain resource name, which is why it is written out
	// here rather than following the shape of the reads around it. A policy carrying a conditional
	// binding is only returned in full at version 3, so that is what is asked for.
	policy, err := service.Projects.Locations.Queues.GetIamPolicy(resource,
		&cloudtasks.GetIamPolicyRequest{
			Options: &cloudtasks.GetPolicyOptions{RequestedPolicyVersion: iamPolicyVersionWithConditions},
		}).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the queue %s does not exist in %s in project %s", queueID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for queue %s in %s in project %s: %w", queueID, location, projectID, err)
	}

	return policy, nil
}
