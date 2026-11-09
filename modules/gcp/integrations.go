package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/integrations/v1"
	"google.golang.org/api/option"
)

// GetIntegrationsAuthConfigAttrs returns the settings Google Cloud holds for the given Application Integration auth config, so a test can assert on
// what was actually created rather than only that it exists.
// The credentials an auth config holds are write only, so a read answers with the state, the
// expiry and the certificate info Google derived from them.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIntegrationsAuthConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) *integrations.GoogleCloudIntegrationsV1alphaAuthConfig {
	result, err := GetIntegrationsAuthConfigAttrsE(t, ctx, projectID, location, configID)
	require.NoError(t, err)

	return result
}

// GetIntegrationsAuthConfigAttrsE returns the settings Google Cloud holds for the given Application Integration auth config.
// The ctx parameter supports cancellation and timeouts.
func GetIntegrationsAuthConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, configID string) (*integrations.GoogleCloudIntegrationsV1alphaAuthConfig, error) {
	logger.Default.Logf(t, "Getting settings for Application Integration auth config %s in %s of project %s", configID, location, projectID)

	service, err := NewIntegrationsServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIntegrationsAuthConfigAttrsWithClient(ctx, service, projectID, location, configID)
}

// GetIntegrationsAuthConfigAttrsWithClient returns the settings Google Cloud holds for the given Application Integration auth config using the supplied
// *integrations.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see integrations_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIntegrationsAuthConfigAttrsWithClient(ctx context.Context, service *integrations.Service, projectID string, location string, configID string) (*integrations.GoogleCloudIntegrationsV1alphaAuthConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/authConfigs/%s", projectID, location, configID)

	result, err := service.Projects.Locations.AuthConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Application Integration auth config %s does not exist in %s of project %s", configID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Application Integration auth config %s in %s of project %s: %w", configID, location, projectID, err)
	}

	return result, nil
}

// NewIntegrationsServiceE creates a Application Integration service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewIntegrationsServiceE(t testing.TestingT, ctx context.Context) (*integrations.Service, error) {
	return integrations.NewService(ctx, append(withOptions(), option.WithScopes(integrations.CloudPlatformScope))...)
}
