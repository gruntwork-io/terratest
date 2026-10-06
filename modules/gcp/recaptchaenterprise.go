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
	"google.golang.org/api/recaptchaenterprise/v1"
)

// GetRecaptchaKeyAttrs returns the settings Google Cloud holds for the reCAPTCHA key, so a test can assert
// on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRecaptchaKeyAttrs(t testing.TestingT, ctx context.Context, projectID string, keyID string) *recaptchaenterprise.GoogleCloudRecaptchaenterpriseV1Key {
	result, err := GetRecaptchaKeyAttrsE(t, ctx, projectID, keyID)
	require.NoError(t, err)

	return result
}

// GetRecaptchaKeyAttrsE returns the settings Google Cloud holds for the reCAPTCHA key.
// The ctx parameter supports cancellation and timeouts.
func GetRecaptchaKeyAttrsE(t testing.TestingT, ctx context.Context, projectID string, keyID string) (*recaptchaenterprise.GoogleCloudRecaptchaenterpriseV1Key, error) {
	logger.Default.Logf(t, "Getting settings for reCAPTCHA key %s in project %s", keyID, projectID)

	service, err := NewRecaptchaEnterpriseServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetRecaptchaKeyAttrsWithClient(ctx, service, projectID, keyID)
}

// GetRecaptchaKeyAttrsWithClient returns the settings Google Cloud holds for the reCAPTCHA key using the
// supplied *recaptchaenterprise.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see recaptchaenterprise_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetRecaptchaKeyAttrsWithClient(ctx context.Context, service *recaptchaenterprise.Service, projectID string, keyID string) (*recaptchaenterprise.GoogleCloudRecaptchaenterpriseV1Key, error) {
	name := fmt.Sprintf("projects/%s/keys/%s", projectID, keyID)

	result, err := service.Projects.Keys.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the reCAPTCHA key %s in project %s does not exist", keyID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for reCAPTCHA key %s in project %s: %w", keyID, projectID, err)
	}

	return result, nil
}

// NewRecaptchaEnterpriseServiceE creates a reCAPTCHA Enterprise service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewRecaptchaEnterpriseServiceE(t testing.TestingT, ctx context.Context) (*recaptchaenterprise.Service, error) {
	return recaptchaenterprise.NewService(ctx, append(withOptions(), option.WithScopes(recaptchaenterprise.CloudPlatformScope))...)
}
