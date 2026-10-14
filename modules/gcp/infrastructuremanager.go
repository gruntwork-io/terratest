package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/config/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetInfraManagerDeploymentAttrs returns the settings Google Cloud holds for the Infrastructure Manager deployment, so a test can assert on what was
// actually created rather than only that it exists.
// The state and the latest revision are what say whether Infrastructure Manager actually applied the
// configuration the deployment points at.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetInfraManagerDeploymentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, deploymentID string) *config.Deployment {
	result, err := GetInfraManagerDeploymentAttrsE(t, ctx, projectID, location, deploymentID)
	require.NoError(t, err)

	return result
}

// GetInfraManagerDeploymentAttrsE returns the settings Google Cloud holds for the Infrastructure Manager deployment.
// The ctx parameter supports cancellation and timeouts.
func GetInfraManagerDeploymentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, deploymentID string) (*config.Deployment, error) {
	logger.Default.Logf(t, "Getting settings for Infrastructure Manager deployment %s in %s in project %s", deploymentID, location, projectID)

	service, err := NewInfraManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetInfraManagerDeploymentAttrsWithClient(ctx, service, projectID, location, deploymentID)
}

// GetInfraManagerDeploymentAttrsWithClient returns the settings Google Cloud holds for the Infrastructure Manager deployment using the supplied
// *config.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see infrastructuremanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetInfraManagerDeploymentAttrsWithClient(ctx context.Context, service *config.Service, projectID string, location string, deploymentID string) (*config.Deployment, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/deployments/%s", projectID, location, deploymentID)

	result, err := service.Projects.Locations.Deployments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Infrastructure Manager deployment %s in %s in project %s does not exist", deploymentID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Infrastructure Manager deployment %s in %s in project %s: %w", deploymentID, location, projectID, err)
	}

	return result, nil
}

// NewInfraManagerServiceE creates a Infrastructure Manager service authenticated the same way every other client in this
// module is.
// Google serves Infrastructure Manager from the config API, which is why the import reads config
// rather than anything with infra in the name.
// The ctx parameter supports cancellation and timeouts.
func NewInfraManagerServiceE(t testing.TestingT, ctx context.Context) (*config.Service, error) {
	return config.NewService(ctx, append(withOptions(), option.WithScopes(config.CloudPlatformScope))...)
}
