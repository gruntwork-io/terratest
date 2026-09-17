package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
)

// GetOAuthClientAttrs returns the settings Google Cloud holds for the given IAM OAuth client, so a
// test can assert on what was actually created rather than only that it exists. A client lives in
// a location, which has to be given, and the id is the one the caller chose when creating it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetOAuthClientAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, clientID string) *iam.OauthClient {
	client, err := GetOAuthClientAttrsE(t, ctx, projectID, location, clientID)
	require.NoError(t, err)

	return client
}

// GetOAuthClientAttrsE returns the settings Google Cloud holds for the given IAM OAuth client.
// The ctx parameter supports cancellation and timeouts.
func GetOAuthClientAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, clientID string) (*iam.OauthClient, error) {
	logger.Default.Logf(t, "Getting settings for OAuth client %s in project %s", clientID, projectID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetOAuthClientAttrsWithClient(ctx, service, projectID, location, clientID)
}

// GetOAuthClientAttrsWithClient returns the settings Google Cloud holds for the given IAM OAuth
// client using the supplied *iam.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see oauthclient_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetOAuthClientAttrsWithClient(ctx context.Context, service *iam.Service, projectID string, location string, clientID string) (*iam.OauthClient, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/oauthClients/%s", projectID, location, clientID)

	client, err := service.Projects.Locations.OauthClients.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("OAuth client %s does not exist in project %s", clientID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for OAuth client %s in project %s: %w", clientID, projectID, err)
	}

	return client, nil
}
