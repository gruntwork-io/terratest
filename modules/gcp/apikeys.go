package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/apikeys/v2"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetAPIKeyAttrs returns the settings Google Cloud holds for the given API key, so a test can
// assert on what was actually created rather than only that it exists. The id is the one Google
// assigns, and this reads the key's settings and never its secret string, which no test needs and
// which would land in a log.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAPIKeyAttrs(t testing.TestingT, ctx context.Context, projectID string, keyID string) *apikeys.V2Key {
	key, err := GetAPIKeyAttrsE(t, ctx, projectID, keyID)
	require.NoError(t, err)

	return key
}

// GetAPIKeyAttrsE returns the settings Google Cloud holds for the given API key.
// The ctx parameter supports cancellation and timeouts.
func GetAPIKeyAttrsE(t testing.TestingT, ctx context.Context, projectID string, keyID string) (*apikeys.V2Key, error) {
	logger.Default.Logf(t, "Getting settings for API key %s in project %s", keyID, projectID)

	service, err := NewAPIKeysServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAPIKeyAttrsWithClient(ctx, service, projectID, keyID)
}

// GetAPIKeyAttrsWithClient returns the settings Google Cloud holds for the given API key using the
// supplied *apikeys.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see apikeys_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAPIKeyAttrsWithClient(ctx context.Context, service *apikeys.Service, projectID string, keyID string) (*apikeys.V2Key, error) {
	name := fmt.Sprintf("projects/%s/locations/global/keys/%s", projectID, keyID)

	key, err := service.Projects.Locations.Keys.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("API key %s does not exist in project %s", keyID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for API key %s in project %s: %w", keyID, projectID, err)
	}

	return key, nil
}

// NewAPIKeysServiceE creates an API Keys service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewAPIKeysServiceE(t testing.TestingT, ctx context.Context) (*apikeys.Service, error) {
	return apikeys.NewService(ctx, append(withOptions(), option.WithScopes(apikeys.CloudPlatformScope))...)
}
