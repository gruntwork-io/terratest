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
	"google.golang.org/api/secretmanager/v1"
)

// GetSecretAttrs returns the settings Google Cloud holds for the given Secret Manager secret, so a
// test can assert on what was actually created rather than only that it exists. This reads the
// secret's settings and never its value, which no test needs and which would land in a log.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSecretAttrs(t testing.TestingT, ctx context.Context, projectID string, secretID string) *secretmanager.Secret {
	secret, err := GetSecretAttrsE(t, ctx, projectID, secretID)
	require.NoError(t, err)

	return secret
}

// GetSecretAttrsE returns the settings Google Cloud holds for the given Secret Manager secret.
// The ctx parameter supports cancellation and timeouts.
func GetSecretAttrsE(t testing.TestingT, ctx context.Context, projectID string, secretID string) (*secretmanager.Secret, error) {
	logger.Default.Logf(t, "Getting settings for secret %s in project %s", secretID, projectID)

	service, err := NewSecretManagerServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSecretAttrsWithClient(ctx, service, projectID, secretID)
}

// GetSecretAttrsWithClient returns the settings Google Cloud holds for the given Secret Manager
// secret using the supplied *secretmanager.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see secretmanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSecretAttrsWithClient(ctx context.Context, service *secretmanager.Service, projectID string, secretID string) (*secretmanager.Secret, error) {
	name := fmt.Sprintf("projects/%s/secrets/%s", projectID, secretID)

	secret, err := service.Projects.Secrets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("secret %s does not exist in project %s", secretID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for secret %s in project %s: %w", secretID, projectID, err)
	}

	return secret, nil
}

// NewSecretManagerServiceE creates a Secret Manager service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewSecretManagerServiceE(t testing.TestingT, ctx context.Context) (*secretmanager.Service, error) {
	return secretmanager.NewService(ctx, append(withOptions(), option.WithScopes(secretmanager.CloudPlatformScope))...)
}
