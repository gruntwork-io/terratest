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
	"google.golang.org/api/workflows/v1"
)

// GetWorkflowAttrs returns the settings Google Cloud holds for the given workflow, so a test can
// assert on what was actually created rather than only that it exists. A workflow is regional, so
// the region it was created in has to be given.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkflowAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, workflowName string) *workflows.Workflow {
	workflow, err := GetWorkflowAttrsE(t, ctx, projectID, region, workflowName)
	require.NoError(t, err)

	return workflow
}

// GetWorkflowAttrsE returns the settings Google Cloud holds for the given workflow.
// The ctx parameter supports cancellation and timeouts.
func GetWorkflowAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, workflowName string) (*workflows.Workflow, error) {
	logger.Default.Logf(t, "Getting settings for workflow %s in project %s", workflowName, projectID)

	service, err := NewWorkflowsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkflowAttrsWithClient(ctx, service, projectID, region, workflowName)
}

// GetWorkflowAttrsWithClient returns the settings Google Cloud holds for the given workflow using
// the supplied *workflows.Service. Prefer this variant in unit tests where the service is backed by
// an httptest fake server (see workflows_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkflowAttrsWithClient(ctx context.Context, service *workflows.Service, projectID string, region string, workflowName string) (*workflows.Workflow, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/workflows/%s", projectID, region, workflowName)

	workflow, err := service.Projects.Locations.Workflows.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("workflow %s does not exist in project %s", workflowName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for workflow %s in project %s: %w", workflowName, projectID, err)
	}

	return workflow, nil
}

// NewWorkflowsServiceE creates a Workflows service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewWorkflowsServiceE(t testing.TestingT, ctx context.Context) (*workflows.Service, error) {
	return workflows.NewService(ctx, append(withOptions(), option.WithScopes(workflows.CloudPlatformScope))...)
}
