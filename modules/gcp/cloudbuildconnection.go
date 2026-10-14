package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	cloudbuildv2 "google.golang.org/api/cloudbuild/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetCloudBuildConnectionIamPolicyAttrs returns the IAM policy Google Cloud holds for the Cloud Build connection, so a test can assert on what was
// actually created rather than only that it exists.
// The connection resources live on the v2 surface, which is a separate client from the v1 one the
// build trigger and worker pool reads use.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudBuildConnectionIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) *cloudbuildv2.Policy {
	result, err := GetCloudBuildConnectionIamPolicyAttrsE(t, ctx, projectID, location, connectionID)
	require.NoError(t, err)

	return result
}

// GetCloudBuildConnectionIamPolicyAttrsE returns the IAM policy Google Cloud holds for the Cloud Build connection.
// The ctx parameter supports cancellation and timeouts.
func GetCloudBuildConnectionIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, connectionID string) (*cloudbuildv2.Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for Cloud Build connection %s in %s in project %s", connectionID, location, projectID)

	service, err := NewCloudBuildV2ServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudBuildConnectionIamPolicyAttrsWithClient(ctx, service, projectID, location, connectionID)
}

// GetCloudBuildConnectionIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the Cloud Build connection using the supplied
// *cloudbuildv2.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see cloudbuildconnection_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudBuildConnectionIamPolicyAttrsWithClient(ctx context.Context, service *cloudbuildv2.Service, projectID string, location string, connectionID string) (*cloudbuildv2.Policy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/connections/%s", projectID, location, connectionID)

	result, err := service.Projects.Locations.Connections.GetIamPolicy(name).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Build connection %s in %s in project %s does not exist", connectionID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for Cloud Build connection %s in %s in project %s: %w", connectionID, location, projectID, err)
	}

	return result, nil
}

// NewCloudBuildV2ServiceE creates a Cloud Build v2 service authenticated the same way every other
// client in this module is. Connections and repositories are v2 resources, and the v1 client the
// build trigger and worker pool reads use does not carry them.
// The ctx parameter supports cancellation and timeouts.
func NewCloudBuildV2ServiceE(t testing.TestingT, ctx context.Context) (*cloudbuildv2.Service, error) {
	return cloudbuildv2.NewService(ctx, append(withOptions(), option.WithScopes(cloudbuildv2.CloudPlatformScope))...)
}
