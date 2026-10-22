package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	notebooksv2 "google.golang.org/api/notebooks/v2"
	"google.golang.org/api/option"
)

// GetWorkbenchInstanceAttrs returns the settings Google Cloud holds for the given Workbench instance, so a test can assert on
// what was actually created rather than only that it exists.
// Workbench instances are served by the second version of the Notebooks API, while the managed
// runtimes in notebooks.go come from the first.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkbenchInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) *notebooksv2.Instance {
	result, err := GetWorkbenchInstanceAttrsE(t, ctx, projectID, location, instanceID)
	require.NoError(t, err)

	return result
}

// GetWorkbenchInstanceAttrsE returns the settings Google Cloud holds for the given Workbench instance.
// The ctx parameter supports cancellation and timeouts.
func GetWorkbenchInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) (*notebooksv2.Instance, error) {
	logger.Default.Logf(t, "Getting settings for Workbench instance %s in %s of project %s", instanceID, location, projectID)

	service, err := NewWorkbenchServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkbenchInstanceAttrsWithClient(ctx, service, projectID, location, instanceID)
}

// GetWorkbenchInstanceAttrsWithClient returns the settings Google Cloud holds for the given Workbench instance using the supplied
// *notebooksv2.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see notebooksv2_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkbenchInstanceAttrsWithClient(ctx context.Context, service *notebooksv2.Service, projectID string, location string, instanceID string) (*notebooksv2.Instance, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/instances/%s", projectID, location, instanceID)

	result, err := service.Projects.Locations.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Workbench instance %s does not exist in %s of project %s", instanceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Workbench instance %s in %s of project %s: %w", instanceID, location, projectID, err)
	}

	return result, nil
}

// NewWorkbenchServiceE creates a Workbench service authenticated the same way every other client in
// this module is.
// Workbench instances live behind the second version of the Notebooks API, which is why this client
// sits beside the one notebooks.go builds.
// The ctx parameter supports cancellation and timeouts.
func NewWorkbenchServiceE(t testing.TestingT, ctx context.Context) (*notebooksv2.Service, error) {
	return notebooksv2.NewService(ctx, append(withOptions(), option.WithScopes(notebooksv2.CloudPlatformScope))...)
}
