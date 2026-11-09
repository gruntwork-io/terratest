package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/cloudbuild/v1"
	"google.golang.org/api/googleapi"
)

// GetCloudBuildBitbucketServerConfigAttrs returns the settings Google Cloud holds for the given Cloud Build Bitbucket Server config, so a test can assert on
// what was actually created rather than only that it exists.
// The secrets a config is created with are write only, so a read answers with the host URI, the
// username and the webhook key rather than with what was sent.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudBuildBitbucketServerConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) *cloudbuild.BitbucketServerConfig {
	result, err := GetCloudBuildBitbucketServerConfigAttrsE(t, ctx, projectID, location, configID)
	require.NoError(t, err)

	return result
}

// GetCloudBuildBitbucketServerConfigAttrsE returns the settings Google Cloud holds for the given Cloud Build Bitbucket Server config.
// The ctx parameter supports cancellation and timeouts.
func GetCloudBuildBitbucketServerConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) (*cloudbuild.BitbucketServerConfig, error) {
	logger.Default.Logf(t, "Getting settings for Cloud Build Bitbucket Server config %s in %s of project %s", configID, location, projectID)

	service, err := NewCloudBuildRESTServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetCloudBuildBitbucketServerConfigAttrsWithClient(ctx, service, projectID, location, configID)
}

// GetCloudBuildBitbucketServerConfigAttrsWithClient returns the settings Google Cloud holds for the given Cloud Build Bitbucket Server config using the supplied
// *cloudbuild.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see cloudbuildbitbucket_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudBuildBitbucketServerConfigAttrsWithClient(ctx context.Context, service *cloudbuild.Service, projectID string, location string, configID string) (*cloudbuild.BitbucketServerConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/bitbucketServerConfigs/%s", projectID, location, configID)

	result, err := service.Projects.Locations.BitbucketServerConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Build Bitbucket Server config %s does not exist in %s of project %s", configID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud Build Bitbucket Server config %s in %s of project %s: %w", configID, location, projectID, err)
	}

	return result, nil
}
