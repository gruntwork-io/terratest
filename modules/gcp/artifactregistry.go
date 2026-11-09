package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/artifactregistry/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetArtifactRegistryRepositoryAttrs returns the settings Google Cloud holds for the given Artifact
// Registry repository, so a test can assert on what was actually created rather than only that it
// exists. A repository is regional, so the location it was created in has to be given.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetArtifactRegistryRepositoryAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, repositoryID string) *artifactregistry.Repository {
	repository, err := GetArtifactRegistryRepositoryAttrsE(t, ctx, projectID, location, repositoryID)
	require.NoError(t, err)

	return repository
}

// GetArtifactRegistryRepositoryAttrsE returns the settings Google Cloud holds for the given
// Artifact Registry repository.
// The ctx parameter supports cancellation and timeouts.
func GetArtifactRegistryRepositoryAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, repositoryID string) (*artifactregistry.Repository, error) {
	logger.Default.Logf(t, "Getting settings for Artifact Registry repository %s in project %s", repositoryID, projectID)

	service, err := NewArtifactRegistryServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetArtifactRegistryRepositoryAttrsWithClient(ctx, service, projectID, location, repositoryID)
}

// GetArtifactRegistryRepositoryAttrsWithClient returns the settings Google Cloud holds for the
// given Artifact Registry repository using the supplied *artifactregistry.Service. Prefer this
// variant in unit tests where the service is backed by an httptest fake server (see
// artifactregistry_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetArtifactRegistryRepositoryAttrsWithClient(ctx context.Context, service *artifactregistry.Service, projectID string, location string, repositoryID string) (*artifactregistry.Repository, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/repositories/%s", projectID, location, repositoryID)

	repository, err := service.Projects.Locations.Repositories.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("Artifact Registry repository %s does not exist in project %s", repositoryID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Artifact Registry repository %s in project %s: %w", repositoryID, projectID, err)
	}

	return repository, nil
}

// GetArtifactRegistryProjectConfigAttrs returns the settings Google Cloud holds for Artifact
// Registry in the given project and location, so a test can assert on what was actually created
// rather than only that it exists.
// The config is one object per project and location rather than a resource with a name, and what it
// holds is whether Artifact Registry writes platform logs there and at which severity.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetArtifactRegistryProjectConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string) *artifactregistry.ProjectConfig {
	result, err := GetArtifactRegistryProjectConfigAttrsE(t, ctx, projectID, location)
	require.NoError(t, err)

	return result
}

// GetArtifactRegistryProjectConfigAttrsE returns the settings Google Cloud holds for Artifact
// Registry in the given project and location.
// The ctx parameter supports cancellation and timeouts.
func GetArtifactRegistryProjectConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string) (*artifactregistry.ProjectConfig, error) {
	logger.Default.Logf(t, "Getting the Artifact Registry config for %s of project %s", location, projectID)

	service, err := NewArtifactRegistryServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetArtifactRegistryProjectConfigAttrsWithClient(ctx, service, projectID, location)
}

// GetArtifactRegistryProjectConfigAttrsWithClient returns the settings Google Cloud holds for
// Artifact Registry in the given project and location using the supplied *artifactregistry.Service.
// Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see artifactregistry_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetArtifactRegistryProjectConfigAttrsWithClient(ctx context.Context, service *artifactregistry.Service, projectID string, location string) (*artifactregistry.ProjectConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/projectConfig", projectID, location)

	result, err := service.Projects.Locations.GetProjectConfig(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the project %s has no Artifact Registry config in %s", projectID, location)
		}

		return nil, fmt.Errorf("failed to get the Artifact Registry config for %s of project %s: %w", location, projectID, err)
	}

	return result, nil
}

// NewArtifactRegistryServiceE creates an Artifact Registry service authenticated the same way every
// other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewArtifactRegistryServiceE(t testing.TestingT, ctx context.Context) (*artifactregistry.Service, error) {
	return artifactregistry.NewService(ctx, append(withOptions(), option.WithScopes(artifactregistry.CloudPlatformScope))...)
}
