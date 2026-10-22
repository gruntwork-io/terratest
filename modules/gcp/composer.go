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

// GetComposerUserWorkloadsConfigMapAttrs returns the settings Google Cloud holds for the given Cloud Composer user workloads config map, so a test can assert on
// what was actually created rather than only that it exists.
// The map's data comes back as the keys and values Airflow will read, which is the whole of what the
// module sets.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetComposerUserWorkloadsConfigMapAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, environmentID string, configMapID string) *composer.UserWorkloadsConfigMap {
	result, err := GetComposerUserWorkloadsConfigMapAttrsE(t, ctx, projectID, location, environmentID, configMapID)
	require.NoError(t, err)

	return result
}

// GetComposerUserWorkloadsConfigMapAttrsE returns the settings Google Cloud holds for the given Cloud Composer user workloads config map.
// The ctx parameter supports cancellation and timeouts.
func GetComposerUserWorkloadsConfigMapAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, environmentID string, configMapID string) (*composer.UserWorkloadsConfigMap, error) {
	logger.Default.Logf(t, "Getting settings for user workloads config map %s on Composer environment %s in %s of project %s", configMapID, environmentID, location, projectID)

	service, err := NewComposerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetComposerUserWorkloadsConfigMapAttrsWithClient(ctx, service, projectID, location, environmentID, configMapID)
}

// GetComposerUserWorkloadsConfigMapAttrsWithClient returns the settings Google Cloud holds for the given Cloud Composer user workloads config map using the supplied
// *composer.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see composer_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetComposerUserWorkloadsConfigMapAttrsWithClient(ctx context.Context, service *composer.Service, projectID string, location string, environmentID string, configMapID string) (*composer.UserWorkloadsConfigMap, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/environments/%s/userWorkloadsConfigMaps/%s", projectID, location, environmentID, configMapID)

	result, err := service.Projects.Locations.Environments.UserWorkloadsConfigMaps.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the user workloads config map %s does not exist on Composer environment %s in %s of project %s", configMapID, environmentID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for user workloads config map %s on Composer environment %s in %s of project %s: %w", configMapID, environmentID, location, projectID, err)
	}

	return result, nil
}

// NewComposerServiceE creates a Cloud Composer service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewComposerServiceE(t testing.TestingT, ctx context.Context) (*composer.Service, error) {
	return composer.NewService(ctx, append(withOptions(), option.WithScopes(composer.CloudPlatformScope))...)
}
