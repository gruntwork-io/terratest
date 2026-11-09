package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/deploymentmanager/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetDeploymentManagerDeploymentAttrs returns the settings Google Cloud holds for the given Deployment Manager deployment, so a test can assert on
// what was actually created rather than only that it exists.
// The manifest Google built from the submitted config, and the fingerprint an update has to carry, are
// both reported on the deployment rather than set by the caller.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDeploymentManagerDeploymentAttrs(t testing.TestingT, ctx context.Context, projectID string, deploymentID string) *deploymentmanager.Deployment {
	result, err := GetDeploymentManagerDeploymentAttrsE(t, ctx, projectID, deploymentID)
	require.NoError(t, err)

	return result
}

// GetDeploymentManagerDeploymentAttrsE returns the settings Google Cloud holds for the given Deployment Manager deployment.
// The ctx parameter supports cancellation and timeouts.
func GetDeploymentManagerDeploymentAttrsE(t testing.TestingT, ctx context.Context, projectID string, deploymentID string) (*deploymentmanager.Deployment, error) {
	logger.Default.Logf(t, "Getting settings for Deployment Manager deployment %s in project %s", deploymentID, projectID)

	service, err := NewDeploymentManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDeploymentManagerDeploymentAttrsWithClient(ctx, service, projectID, deploymentID)
}

// GetDeploymentManagerDeploymentAttrsWithClient returns the settings Google Cloud holds for the given Deployment Manager deployment using the supplied
// *deploymentmanager.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see deploymentmanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDeploymentManagerDeploymentAttrsWithClient(ctx context.Context, service *deploymentmanager.Service, projectID string, deploymentID string) (*deploymentmanager.Deployment, error) {
	result, err := service.Deployments.Get(projectID, deploymentID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Deployment Manager deployment %s does not exist in project %s", deploymentID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Deployment Manager deployment %s in project %s: %w", deploymentID, projectID, err)
	}

	return result, nil
}

// NewDeploymentManagerServiceE creates a Deployment Manager service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewDeploymentManagerServiceE(t testing.TestingT, ctx context.Context) (*deploymentmanager.Service, error) {
	return deploymentmanager.NewService(ctx, append(withOptions(), option.WithScopes(deploymentmanager.CloudPlatformScope))...)
}
