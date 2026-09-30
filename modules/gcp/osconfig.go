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
	"google.golang.org/api/osconfig/v1"
)

// GetPatchDeploymentAttrs returns the settings Google Cloud holds for the given OS Config patch
// deployment, so a test can assert on which instances it patches and when rather than only that it
// exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetPatchDeploymentAttrs(t testing.TestingT, ctx context.Context, projectID string, deploymentID string) *osconfig.PatchDeployment {
	deployment, err := GetPatchDeploymentAttrsE(t, ctx, projectID, deploymentID)
	require.NoError(t, err)

	return deployment
}

// GetPatchDeploymentAttrsE returns the settings Google Cloud holds for the given OS Config patch
// deployment.
// The ctx parameter supports cancellation and timeouts.
func GetPatchDeploymentAttrsE(t testing.TestingT, ctx context.Context, projectID string, deploymentID string) (*osconfig.PatchDeployment, error) {
	logger.Default.Logf(t, "Getting settings for patch deployment %s in project %s", deploymentID, projectID)

	service, err := NewOSConfigServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetPatchDeploymentAttrsWithClient(ctx, service, projectID, deploymentID)
}

// GetPatchDeploymentAttrsWithClient returns the settings Google Cloud holds for the given OS Config
// patch deployment using the supplied *osconfig.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see osconfig_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetPatchDeploymentAttrsWithClient(ctx context.Context, service *osconfig.Service, projectID string, deploymentID string) (*osconfig.PatchDeployment, error) {
	name := fmt.Sprintf("projects/%s/patchDeployments/%s", projectID, deploymentID)

	deployment, err := service.Projects.PatchDeployments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the patch deployment %s in project %s does not exist", deploymentID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for patch deployment %s in project %s: %w", deploymentID, projectID, err)
	}

	return deployment, nil
}

// GetOSPolicyAssignmentAttrs returns the settings Google Cloud holds for the given OS policy
// assignment, so a test can assert on the policies it applies and the instances it applies them to.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetOSPolicyAssignmentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, assignmentID string) *osconfig.OSPolicyAssignment {
	assignment, err := GetOSPolicyAssignmentAttrsE(t, ctx, projectID, location, assignmentID)
	require.NoError(t, err)

	return assignment
}

// GetOSPolicyAssignmentAttrsE returns the settings Google Cloud holds for the given OS policy
// assignment.
// The ctx parameter supports cancellation and timeouts.
func GetOSPolicyAssignmentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, assignmentID string) (*osconfig.OSPolicyAssignment, error) {
	logger.Default.Logf(t, "Getting settings for OS policy assignment %s in %s in project %s", assignmentID, location, projectID)

	service, err := NewOSConfigServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetOSPolicyAssignmentAttrsWithClient(ctx, service, projectID, location, assignmentID)
}

// GetOSPolicyAssignmentAttrsWithClient returns the settings Google Cloud holds for the given OS
// policy assignment using the supplied *osconfig.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see osconfig_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetOSPolicyAssignmentAttrsWithClient(ctx context.Context, service *osconfig.Service, projectID string, location string, assignmentID string) (*osconfig.OSPolicyAssignment, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/osPolicyAssignments/%s", projectID, location, assignmentID)

	assignment, err := service.Projects.Locations.OsPolicyAssignments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the OS policy assignment %s in %s in project %s does not exist", assignmentID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for OS policy assignment %s in %s in project %s: %w", assignmentID, location, projectID, err)
	}

	return assignment, nil
}

// NewOSConfigServiceE creates an OS Config service authenticated the same way every other client in
// this module is. It serves both the patch deployment calls and the OS policy assignment calls.
// The ctx parameter supports cancellation and timeouts.
func NewOSConfigServiceE(t testing.TestingT, ctx context.Context) (*osconfig.Service, error) {
	return osconfig.NewService(ctx, append(withOptions(), option.WithScopes(osconfig.CloudPlatformScope))...)
}
