package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/firebaseappcheck/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetFirebaseAppCheckAppAttestConfigAttrs returns the settings Google Cloud holds for the given Firebase App Check App Attest config, so a test can assert on what was
// actually created rather than only that it exists.
// The config decides how long an App Attest token stays valid, so that time-to-live is the whole setting.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckAppAttestConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, appID string) *firebaseappcheck.GoogleFirebaseAppcheckV1AppAttestConfig {
	attrs, err := GetFirebaseAppCheckAppAttestConfigAttrsE(t, ctx, projectID, appID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseAppCheckAppAttestConfigAttrsE returns the settings Google Cloud holds for the given Firebase App Check App Attest config.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckAppAttestConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, appID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1AppAttestConfig, error) {
	logger.Default.Logf(t, "Getting settings for Firebase App Check App Attest config for app %s in project %s", appID, projectID)

	service, err := NewFirebaseAppCheckServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseAppCheckAppAttestConfigAttrsWithClient(ctx, service, projectID, appID)
}

// GetFirebaseAppCheckAppAttestConfigAttrsWithClient returns the settings Google Cloud holds for the given Firebase App Check App Attest config using the supplied
// *firebaseappcheck.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebaseappcheck_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckAppAttestConfigAttrsWithClient(ctx context.Context, service *firebaseappcheck.Service, projectID string, appID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1AppAttestConfig, error) {
	name := fmt.Sprintf("projects/%s/apps/%s/appAttestConfig", projectID, appID)

	attrs, err := service.Projects.Apps.AppAttestConfig.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase App Check App Attest config for app %s in project %s does not exist", appID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase App Check App Attest config for app %s in project %s: %w", appID, projectID, err)
	}

	return attrs, nil
}

// GetFirebaseAppCheckDeviceCheckConfigAttrs returns the settings Google Cloud holds for the given Firebase App Check DeviceCheck config, so a test can assert on what was
// actually created rather than only that it exists.
// The config holds the Apple key App Check verifies with, so the key and team it names decide which devices are trusted.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckDeviceCheckConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, appID string) *firebaseappcheck.GoogleFirebaseAppcheckV1DeviceCheckConfig {
	attrs, err := GetFirebaseAppCheckDeviceCheckConfigAttrsE(t, ctx, projectID, appID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseAppCheckDeviceCheckConfigAttrsE returns the settings Google Cloud holds for the given Firebase App Check DeviceCheck config.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckDeviceCheckConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, appID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1DeviceCheckConfig, error) {
	logger.Default.Logf(t, "Getting settings for Firebase App Check DeviceCheck config for app %s in project %s", appID, projectID)

	service, err := NewFirebaseAppCheckServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseAppCheckDeviceCheckConfigAttrsWithClient(ctx, service, projectID, appID)
}

// GetFirebaseAppCheckDeviceCheckConfigAttrsWithClient returns the settings Google Cloud holds for the given Firebase App Check DeviceCheck config using the supplied
// *firebaseappcheck.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebaseappcheck_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckDeviceCheckConfigAttrsWithClient(ctx context.Context, service *firebaseappcheck.Service, projectID string, appID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1DeviceCheckConfig, error) {
	name := fmt.Sprintf("projects/%s/apps/%s/deviceCheckConfig", projectID, appID)

	attrs, err := service.Projects.Apps.DeviceCheckConfig.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase App Check DeviceCheck config for app %s in project %s does not exist", appID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase App Check DeviceCheck config for app %s in project %s: %w", appID, projectID, err)
	}

	return attrs, nil
}

// GetFirebaseAppCheckPlayIntegrityConfigAttrs returns the settings Google Cloud holds for the given Firebase App Check Play Integrity config, so a test can assert on what was
// actually created rather than only that it exists.
// The config decides how long a Play Integrity token stays valid, so that time-to-live is the whole setting.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckPlayIntegrityConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, appID string) *firebaseappcheck.GoogleFirebaseAppcheckV1PlayIntegrityConfig {
	attrs, err := GetFirebaseAppCheckPlayIntegrityConfigAttrsE(t, ctx, projectID, appID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseAppCheckPlayIntegrityConfigAttrsE returns the settings Google Cloud holds for the given Firebase App Check Play Integrity config.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckPlayIntegrityConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, appID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1PlayIntegrityConfig, error) {
	logger.Default.Logf(t, "Getting settings for Firebase App Check Play Integrity config for app %s in project %s", appID, projectID)

	service, err := NewFirebaseAppCheckServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseAppCheckPlayIntegrityConfigAttrsWithClient(ctx, service, projectID, appID)
}

// GetFirebaseAppCheckPlayIntegrityConfigAttrsWithClient returns the settings Google Cloud holds for the given Firebase App Check Play Integrity config using the supplied
// *firebaseappcheck.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebaseappcheck_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckPlayIntegrityConfigAttrsWithClient(ctx context.Context, service *firebaseappcheck.Service, projectID string, appID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1PlayIntegrityConfig, error) {
	name := fmt.Sprintf("projects/%s/apps/%s/playIntegrityConfig", projectID, appID)

	attrs, err := service.Projects.Apps.PlayIntegrityConfig.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase App Check Play Integrity config for app %s in project %s does not exist", appID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase App Check Play Integrity config for app %s in project %s: %w", appID, projectID, err)
	}

	return attrs, nil
}

// GetFirebaseAppCheckRecaptchaEnterpriseConfigAttrs returns the settings Google Cloud holds for the given Firebase App Check reCAPTCHA Enterprise config, so a test can assert on what was
// actually created rather than only that it exists.
// The config names the site key App Check checks against, so that key decides which web clients pass.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckRecaptchaEnterpriseConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, appID string) *firebaseappcheck.GoogleFirebaseAppcheckV1RecaptchaEnterpriseConfig {
	attrs, err := GetFirebaseAppCheckRecaptchaEnterpriseConfigAttrsE(t, ctx, projectID, appID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseAppCheckRecaptchaEnterpriseConfigAttrsE returns the settings Google Cloud holds for the given Firebase App Check reCAPTCHA Enterprise config.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckRecaptchaEnterpriseConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, appID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1RecaptchaEnterpriseConfig, error) {
	logger.Default.Logf(t, "Getting settings for Firebase App Check reCAPTCHA Enterprise config for app %s in project %s", appID, projectID)

	service, err := NewFirebaseAppCheckServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseAppCheckRecaptchaEnterpriseConfigAttrsWithClient(ctx, service, projectID, appID)
}

// GetFirebaseAppCheckRecaptchaEnterpriseConfigAttrsWithClient returns the settings Google Cloud holds for the given Firebase App Check reCAPTCHA Enterprise config using the supplied
// *firebaseappcheck.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebaseappcheck_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckRecaptchaEnterpriseConfigAttrsWithClient(ctx context.Context, service *firebaseappcheck.Service, projectID string, appID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1RecaptchaEnterpriseConfig, error) {
	name := fmt.Sprintf("projects/%s/apps/%s/recaptchaEnterpriseConfig", projectID, appID)

	attrs, err := service.Projects.Apps.RecaptchaEnterpriseConfig.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase App Check reCAPTCHA Enterprise config for app %s in project %s does not exist", appID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase App Check reCAPTCHA Enterprise config for app %s in project %s: %w", appID, projectID, err)
	}

	return attrs, nil
}

// GetFirebaseAppCheckRecaptchaV3ConfigAttrs returns the settings Google Cloud holds for the given Firebase App Check reCAPTCHA v3 config, so a test can assert on what was
// actually created rather than only that it exists.
// The config holds the v3 secret App Check verifies with, so whether a secret is set decides whether web clients can attest at all.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckRecaptchaV3ConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, appID string) *firebaseappcheck.GoogleFirebaseAppcheckV1RecaptchaV3Config {
	attrs, err := GetFirebaseAppCheckRecaptchaV3ConfigAttrsE(t, ctx, projectID, appID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseAppCheckRecaptchaV3ConfigAttrsE returns the settings Google Cloud holds for the given Firebase App Check reCAPTCHA v3 config.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckRecaptchaV3ConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, appID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1RecaptchaV3Config, error) {
	logger.Default.Logf(t, "Getting settings for Firebase App Check reCAPTCHA v3 config for app %s in project %s", appID, projectID)

	service, err := NewFirebaseAppCheckServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseAppCheckRecaptchaV3ConfigAttrsWithClient(ctx, service, projectID, appID)
}

// GetFirebaseAppCheckRecaptchaV3ConfigAttrsWithClient returns the settings Google Cloud holds for the given Firebase App Check reCAPTCHA v3 config using the supplied
// *firebaseappcheck.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebaseappcheck_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckRecaptchaV3ConfigAttrsWithClient(ctx context.Context, service *firebaseappcheck.Service, projectID string, appID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1RecaptchaV3Config, error) {
	name := fmt.Sprintf("projects/%s/apps/%s/recaptchaV3Config", projectID, appID)

	attrs, err := service.Projects.Apps.RecaptchaV3Config.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase App Check reCAPTCHA v3 config for app %s in project %s does not exist", appID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase App Check reCAPTCHA v3 config for app %s in project %s: %w", appID, projectID, err)
	}

	return attrs, nil
}

// GetFirebaseAppCheckDebugTokenAttrs returns the settings Google Cloud holds for the given Firebase App Check debug token, so a test can assert on what was
// actually created rather than only that it exists.
// A debug token lets one untrusted client through, so its display name is the only way to tell which one it is.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckDebugTokenAttrs(t testing.TestingT, ctx context.Context, projectID string, appID string) *firebaseappcheck.GoogleFirebaseAppcheckV1DebugToken {
	attrs, err := GetFirebaseAppCheckDebugTokenAttrsE(t, ctx, projectID, appID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseAppCheckDebugTokenAttrsE returns the settings Google Cloud holds for the given Firebase App Check debug token.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckDebugTokenAttrsE(t testing.TestingT, ctx context.Context, projectID string, appID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1DebugToken, error) {
	logger.Default.Logf(t, "Getting settings for Firebase App Check debug token for app %s in project %s", appID, projectID)

	service, err := NewFirebaseAppCheckServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseAppCheckDebugTokenAttrsWithClient(ctx, service, projectID, appID)
}

// GetFirebaseAppCheckDebugTokenAttrsWithClient returns the settings Google Cloud holds for the given Firebase App Check debug token using the supplied
// *firebaseappcheck.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebaseappcheck_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckDebugTokenAttrsWithClient(ctx context.Context, service *firebaseappcheck.Service, projectID string, appID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1DebugToken, error) {
	name := fmt.Sprintf("projects/%s/apps/%s/debugTokens/gw-library-test", projectID, appID)

	attrs, err := service.Projects.Apps.DebugTokens.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase App Check debug token for app %s in project %s does not exist", appID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase App Check debug token for app %s in project %s: %w", appID, projectID, err)
	}

	return attrs, nil
}

// GetFirebaseAppCheckServiceConfigAttrs returns the settings Google Cloud holds for the given Firebase App Check service config, so a test can assert on what was
// actually created rather than only that it exists.
// The config decides whether App Check is enforced for one Google service, so its enforcement mode is exactly what it does.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckServiceConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, serviceID string) *firebaseappcheck.GoogleFirebaseAppcheckV1Service {
	attrs, err := GetFirebaseAppCheckServiceConfigAttrsE(t, ctx, projectID, serviceID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseAppCheckServiceConfigAttrsE returns the settings Google Cloud holds for the given Firebase App Check service config.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckServiceConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, serviceID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1Service, error) {
	logger.Default.Logf(t, "Getting settings for Firebase App Check service config for %s in project %s", serviceID, projectID)

	service, err := NewFirebaseAppCheckServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseAppCheckServiceConfigAttrsWithClient(ctx, service, projectID, serviceID)
}

// GetFirebaseAppCheckServiceConfigAttrsWithClient returns the settings Google Cloud holds for the given Firebase App Check service config using the supplied
// *firebaseappcheck.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebaseappcheck_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckServiceConfigAttrsWithClient(ctx context.Context, service *firebaseappcheck.Service, projectID string, serviceID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1Service, error) {
	name := fmt.Sprintf("projects/%s/services/%s", projectID, serviceID)

	attrs, err := service.Projects.Services.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase App Check service config for %s in project %s does not exist", serviceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase App Check service config for %s in project %s: %w", serviceID, projectID, err)
	}

	return attrs, nil
}

// GetFirebaseAppCheckResourcePolicyAttrs returns the settings Google Cloud holds for the given Firebase App Check resource policy, so a test can assert on what was
// actually created rather than only that it exists.
// A resource policy overrides enforcement for one resource inside a service, so its mode and the resource it names are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckResourcePolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, serviceID string, policyID string) *firebaseappcheck.GoogleFirebaseAppcheckV1ResourcePolicy {
	attrs, err := GetFirebaseAppCheckResourcePolicyAttrsE(t, ctx, projectID, serviceID, policyID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseAppCheckResourcePolicyAttrsE returns the settings Google Cloud holds for the given Firebase App Check resource policy.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckResourcePolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, serviceID string, policyID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1ResourcePolicy, error) {
	logger.Default.Logf(t, "Getting settings for Firebase App Check resource policy %s for %s in project %s", policyID, serviceID, projectID)

	service, err := NewFirebaseAppCheckServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseAppCheckResourcePolicyAttrsWithClient(ctx, service, projectID, serviceID, policyID)
}

// GetFirebaseAppCheckResourcePolicyAttrsWithClient returns the settings Google Cloud holds for the given Firebase App Check resource policy using the supplied
// *firebaseappcheck.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebaseappcheck_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppCheckResourcePolicyAttrsWithClient(ctx context.Context, service *firebaseappcheck.Service, projectID string, serviceID string, policyID string) (*firebaseappcheck.GoogleFirebaseAppcheckV1ResourcePolicy, error) {
	name := fmt.Sprintf("projects/%s/services/%s/resourcePolicies/%s", projectID, serviceID, policyID)

	attrs, err := service.Projects.Services.ResourcePolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase App Check resource policy %s for %s in project %s does not exist", policyID, serviceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase App Check resource policy %s for %s in project %s: %w", policyID, serviceID, projectID, err)
	}

	return attrs, nil
}

// NewFirebaseAppCheckServiceE creates a Firebase App Check service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewFirebaseAppCheckServiceE(t testing.TestingT, ctx context.Context) (*firebaseappcheck.Service, error) {
	return firebaseappcheck.NewService(ctx, append(withOptions(), option.WithScopes(firebaseappcheck.CloudPlatformScope))...)
}
