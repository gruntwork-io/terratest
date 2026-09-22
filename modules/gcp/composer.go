package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/composer/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetComposerEnvironmentAttrs returns the settings Google Cloud holds for the given Cloud Composer
// environment, so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetComposerEnvironmentAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, environmentID string) *composer.Environment {
	environment, err := GetComposerEnvironmentAttrsE(t, ctx, projectID, region, environmentID)
	require.NoError(t, err)

	return environment
}

// GetComposerEnvironmentAttrsE returns the settings Google Cloud holds for the given Cloud Composer
// environment.
// The ctx parameter supports cancellation and timeouts.
func GetComposerEnvironmentAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, environmentID string) (*composer.Environment, error) {
	logger.Default.Logf(t, "Getting settings for Cloud Composer environment %s in region %s in project %s", environmentID, region, projectID)

	service, err := NewComposerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetComposerEnvironmentAttrsWithClient(ctx, service, projectID, region, environmentID)
}

// GetComposerEnvironmentAttrsWithClient returns the settings Google Cloud holds for the given Cloud
// Composer environment using the supplied *composer.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see composer_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetComposerEnvironmentAttrsWithClient(ctx context.Context, service *composer.Service, projectID string, region string, environmentID string) (*composer.Environment, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/environments/%s", projectID, region, environmentID)

	environment, err := service.Projects.Locations.Environments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Composer environment %s does not exist in region %s in project %s", environmentID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud Composer environment %s in region %s in project %s: %w", environmentID, region, projectID, err)
	}

	return environment, nil
}

// NewComposerServiceE creates a Cloud Composer service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewComposerServiceE(t testing.TestingT, ctx context.Context) (*composer.Service, error) {
	return composer.NewService(ctx, append(withOptions(), option.WithScopes(composer.CloudPlatformScope))...)
}
