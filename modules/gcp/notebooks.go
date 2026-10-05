package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/notebooks/v1"
	"google.golang.org/api/option"
)

// GetNotebooksEnvironmentAttrs returns the settings Google Cloud holds for the given Notebooks
// environment, so a test can assert on the image an instance would boot from rather than only that
// the environment exists. An environment describes a machine; it does not run one.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNotebooksEnvironmentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, environmentID string) *notebooks.Environment {
	environment, err := GetNotebooksEnvironmentAttrsE(t, ctx, projectID, location, environmentID)
	require.NoError(t, err)

	return environment
}

// GetNotebooksEnvironmentAttrsE returns the settings Google Cloud holds for the given Notebooks
// environment.
// The ctx parameter supports cancellation and timeouts.
func GetNotebooksEnvironmentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, environmentID string) (*notebooks.Environment, error) {
	logger.Default.Logf(t, "Getting settings for Notebooks environment %s in %s in project %s", environmentID, location, projectID)

	service, err := NewNotebooksServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNotebooksEnvironmentAttrsWithClient(ctx, service, projectID, location, environmentID)
}

// GetNotebooksEnvironmentAttrsWithClient returns the settings Google Cloud holds for the given
// Notebooks environment using the supplied *notebooks.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see notebooks_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNotebooksEnvironmentAttrsWithClient(ctx context.Context, service *notebooks.Service, projectID string, location string, environmentID string) (*notebooks.Environment, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/environments/%s", projectID, location, environmentID)

	environment, err := service.Projects.Locations.Environments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Notebooks environment %s in %s in project %s does not exist", environmentID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Notebooks environment %s in %s in project %s: %w", environmentID, location, projectID, err)
	}

	return environment, nil
}

// NewNotebooksServiceE creates a Notebooks service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewNotebooksServiceE(t testing.TestingT, ctx context.Context) (*notebooks.Service, error) {
	return notebooks.NewService(ctx, append(withOptions(), option.WithScopes(notebooks.CloudPlatformScope))...)
}
