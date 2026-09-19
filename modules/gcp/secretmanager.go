package gcp

import (
	"context"
	"errors"
	"fmt"
	"regexp"

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

// GetRegionalSecretAttrs returns the settings Google Cloud holds for the given regional Secret
// Manager secret, so a test can assert on what was actually created rather than only that it
// exists. A regional secret lives in a location, which has to be given: it answers on a
// per-location host, unlike a global secret.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRegionalSecretAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, secretID string) *secretmanager.Secret {
	secret, err := GetRegionalSecretAttrsE(t, ctx, projectID, location, secretID)
	require.NoError(t, err)

	return secret
}

// GetRegionalSecretAttrsE returns the settings Google Cloud holds for the given regional Secret
// Manager secret.
// The ctx parameter supports cancellation and timeouts.
func GetRegionalSecretAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, secretID string) (*secretmanager.Secret, error) {
	logger.Default.Logf(t, "Getting settings for regional secret %s in %s in project %s", secretID, location, projectID)

	service, err := NewRegionalSecretManagerServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetRegionalSecretAttrsWithClient(ctx, service, projectID, location, secretID)
}

// GetRegionalSecretAttrsWithClient returns the settings Google Cloud holds for the given regional
// Secret Manager secret using the supplied *secretmanager.Service, which has to point at that
// location's host. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see secretmanager_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetRegionalSecretAttrsWithClient(ctx context.Context, service *secretmanager.Service, projectID string, location string, secretID string) (*secretmanager.Secret, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/secrets/%s", projectID, location, secretID)

	secret, err := service.Projects.Locations.Secrets.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("regional secret %s does not exist in %s in project %s", secretID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for regional secret %s in %s in project %s: %w", secretID, location, projectID, err)
	}

	return secret, nil
}

// secretLocationPattern is what a Google Cloud location identifier may contain. It is checked
// before a location reaches an endpoint, because a value carrying a slash, a colon or an at sign
// would build a URL pointing at a host of the caller's choosing rather than at Google.
var secretLocationPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// NewRegionalSecretManagerServiceE creates a Secret Manager service pointed at the given location's
// own host, authenticated the same way every other client in this module is. Regional secrets are
// served only from that host, so a service built against the global one returns 404 for a regional
// secret that exists, and a location that is not a plain identifier is refused.
// The ctx parameter supports cancellation and timeouts.
func NewRegionalSecretManagerServiceE(t testing.TestingT, ctx context.Context, location string) (*secretmanager.Service, error) {
	if !secretLocationPattern.MatchString(location) {
		return nil, fmt.Errorf("%q is not a valid location: a location may hold only lowercase letters, digits and hyphens", location)
	}

	opts := append(withOptions(), option.WithScopes(secretmanager.CloudPlatformScope),
		option.WithEndpoint(fmt.Sprintf("https://secretmanager.%s.rep.googleapis.com/", location)))

	return secretmanager.NewService(ctx, opts...)
}

// NewSecretManagerServiceE creates a Secret Manager service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewSecretManagerServiceE(t testing.TestingT, ctx context.Context) (*secretmanager.Service, error) {
	return secretmanager.NewService(ctx, append(withOptions(), option.WithScopes(secretmanager.CloudPlatformScope))...)
}
