package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	identitytoolkit "google.golang.org/api/identitytoolkit/v2"
	"google.golang.org/api/option"
)

// Identity Platform is served by the Identity Toolkit API, which is why these reads say Identity
// Platform in their names and the client says identitytoolkit.

// GetIdentityPlatformConfigAttrs returns the settings Google Cloud holds for the given Identity Platform config, so a test can assert on what was
// actually created rather than only that it exists.
// The config is the project's own sign-in settings, so which methods are allowed and whether sign-up is open are the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformConfigAttrs(t testing.TestingT, ctx context.Context, projectID string) *identitytoolkit.GoogleCloudIdentitytoolkitAdminV2Config {
	attrs, err := GetIdentityPlatformConfigAttrsE(t, ctx, projectID)
	require.NoError(t, err)

	return attrs
}

// GetIdentityPlatformConfigAttrsE returns the settings Google Cloud holds for the given Identity Platform config.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2Config, error) {
	logger.Default.Logf(t, "Getting settings for Identity Platform config for project %s", projectID)

	service, err := NewIdentityPlatformServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIdentityPlatformConfigAttrsWithClient(ctx, service, projectID)
}

// GetIdentityPlatformConfigAttrsWithClient returns the settings Google Cloud holds for the given Identity Platform config using the supplied
// *identitytoolkit.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see identitytoolkit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformConfigAttrsWithClient(ctx context.Context, service *identitytoolkit.Service, projectID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2Config, error) {
	name := fmt.Sprintf("projects/%s/config", projectID)

	attrs, err := service.Projects.GetConfig(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Identity Platform config for project %s does not exist", projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Identity Platform config for project %s: %w", projectID, err)
	}

	return attrs, nil
}

// GetIdentityPlatformTenantAttrs returns the settings Google Cloud holds for the given Identity Platform tenant, so a test can assert on what was
// actually created rather than only that it exists.
// A tenant is a separate user pool inside one project, so whether password sign-in and email link sign-in are on decides how its users get in.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformTenantAttrs(t testing.TestingT, ctx context.Context, projectID string, tenantID string) *identitytoolkit.GoogleCloudIdentitytoolkitAdminV2Tenant {
	attrs, err := GetIdentityPlatformTenantAttrsE(t, ctx, projectID, tenantID)
	require.NoError(t, err)

	return attrs
}

// GetIdentityPlatformTenantAttrsE returns the settings Google Cloud holds for the given Identity Platform tenant.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformTenantAttrsE(t testing.TestingT, ctx context.Context, projectID string, tenantID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2Tenant, error) {
	logger.Default.Logf(t, "Getting settings for Identity Platform tenant %s in project %s", tenantID, projectID)

	service, err := NewIdentityPlatformServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIdentityPlatformTenantAttrsWithClient(ctx, service, projectID, tenantID)
}

// GetIdentityPlatformTenantAttrsWithClient returns the settings Google Cloud holds for the given Identity Platform tenant using the supplied
// *identitytoolkit.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see identitytoolkit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformTenantAttrsWithClient(ctx context.Context, service *identitytoolkit.Service, projectID string, tenantID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2Tenant, error) {
	name := fmt.Sprintf("projects/%s/tenants/%s", projectID, tenantID)

	attrs, err := service.Projects.Tenants.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Identity Platform tenant %s in project %s does not exist", tenantID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Identity Platform tenant %s in project %s: %w", tenantID, projectID, err)
	}

	return attrs, nil
}

// GetIdentityPlatformDefaultSupportedIdpConfigAttrs returns the settings Google Cloud holds for the given Identity Platform default IdP config, so a test can assert on what was
// actually created rather than only that it exists.
// This is how one of Google's built-in providers is turned on, so the client id it carries and whether it is enabled are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformDefaultSupportedIdpConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, idpID string) *identitytoolkit.GoogleCloudIdentitytoolkitAdminV2DefaultSupportedIdpConfig {
	attrs, err := GetIdentityPlatformDefaultSupportedIdpConfigAttrsE(t, ctx, projectID, idpID)
	require.NoError(t, err)

	return attrs
}

// GetIdentityPlatformDefaultSupportedIdpConfigAttrsE returns the settings Google Cloud holds for the given Identity Platform default IdP config.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformDefaultSupportedIdpConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, idpID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2DefaultSupportedIdpConfig, error) {
	logger.Default.Logf(t, "Getting settings for Identity Platform default IdP config %s in project %s", idpID, projectID)

	service, err := NewIdentityPlatformServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIdentityPlatformDefaultSupportedIdpConfigAttrsWithClient(ctx, service, projectID, idpID)
}

// GetIdentityPlatformDefaultSupportedIdpConfigAttrsWithClient returns the settings Google Cloud holds for the given Identity Platform default IdP config using the supplied
// *identitytoolkit.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see identitytoolkit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformDefaultSupportedIdpConfigAttrsWithClient(ctx context.Context, service *identitytoolkit.Service, projectID string, idpID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2DefaultSupportedIdpConfig, error) {
	name := fmt.Sprintf("projects/%s/defaultSupportedIdpConfigs/%s", projectID, idpID)

	attrs, err := service.Projects.DefaultSupportedIdpConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Identity Platform default IdP config %s in project %s does not exist", idpID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Identity Platform default IdP config %s in project %s: %w", idpID, projectID, err)
	}

	return attrs, nil
}

// GetIdentityPlatformInboundSamlConfigAttrs returns the settings Google Cloud holds for the given Identity Platform inbound SAML config, so a test can assert on what was
// actually created rather than only that it exists.
// A SAML config is a federation with an outside identity provider, so the issuer and the SSO URL it names decide who can sign in.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformInboundSamlConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, configID string) *identitytoolkit.GoogleCloudIdentitytoolkitAdminV2InboundSamlConfig {
	attrs, err := GetIdentityPlatformInboundSamlConfigAttrsE(t, ctx, projectID, configID)
	require.NoError(t, err)

	return attrs
}

// GetIdentityPlatformInboundSamlConfigAttrsE returns the settings Google Cloud holds for the given Identity Platform inbound SAML config.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformInboundSamlConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, configID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2InboundSamlConfig, error) {
	logger.Default.Logf(t, "Getting settings for Identity Platform inbound SAML config %s in project %s", configID, projectID)

	service, err := NewIdentityPlatformServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIdentityPlatformInboundSamlConfigAttrsWithClient(ctx, service, projectID, configID)
}

// GetIdentityPlatformInboundSamlConfigAttrsWithClient returns the settings Google Cloud holds for the given Identity Platform inbound SAML config using the supplied
// *identitytoolkit.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see identitytoolkit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformInboundSamlConfigAttrsWithClient(ctx context.Context, service *identitytoolkit.Service, projectID string, configID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2InboundSamlConfig, error) {
	name := fmt.Sprintf("projects/%s/inboundSamlConfigs/%s", projectID, configID)

	attrs, err := service.Projects.InboundSamlConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Identity Platform inbound SAML config %s in project %s does not exist", configID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Identity Platform inbound SAML config %s in project %s: %w", configID, projectID, err)
	}

	return attrs, nil
}

// GetIdentityPlatformOAuthIdpConfigAttrs returns the settings Google Cloud holds for the given Identity Platform OAuth IdP config, so a test can assert on what was
// actually created rather than only that it exists.
// An OAuth config is a federation with an outside OIDC provider, so the issuer and client id it names decide who can sign in.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformOAuthIdpConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, configID string) *identitytoolkit.GoogleCloudIdentitytoolkitAdminV2OAuthIdpConfig {
	attrs, err := GetIdentityPlatformOAuthIdpConfigAttrsE(t, ctx, projectID, configID)
	require.NoError(t, err)

	return attrs
}

// GetIdentityPlatformOAuthIdpConfigAttrsE returns the settings Google Cloud holds for the given Identity Platform OAuth IdP config.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformOAuthIdpConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, configID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2OAuthIdpConfig, error) {
	logger.Default.Logf(t, "Getting settings for Identity Platform OAuth IdP config %s in project %s", configID, projectID)

	service, err := NewIdentityPlatformServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIdentityPlatformOAuthIdpConfigAttrsWithClient(ctx, service, projectID, configID)
}

// GetIdentityPlatformOAuthIdpConfigAttrsWithClient returns the settings Google Cloud holds for the given Identity Platform OAuth IdP config using the supplied
// *identitytoolkit.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see identitytoolkit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformOAuthIdpConfigAttrsWithClient(ctx context.Context, service *identitytoolkit.Service, projectID string, configID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2OAuthIdpConfig, error) {
	name := fmt.Sprintf("projects/%s/oauthIdpConfigs/%s", projectID, configID)

	attrs, err := service.Projects.OauthIdpConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Identity Platform OAuth IdP config %s in project %s does not exist", configID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Identity Platform OAuth IdP config %s in project %s: %w", configID, projectID, err)
	}

	return attrs, nil
}

// GetIdentityPlatformTenantDefaultSupportedIdpConfigAttrs returns the settings Google Cloud holds for the given Identity Platform tenant default IdP config, so a test can assert on what was
// actually created rather than only that it exists.
// A tenant turns on built-in providers separately from the project, so this is a different resource from the project-level one.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformTenantDefaultSupportedIdpConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, tenantID string, idpID string) *identitytoolkit.GoogleCloudIdentitytoolkitAdminV2DefaultSupportedIdpConfig {
	attrs, err := GetIdentityPlatformTenantDefaultSupportedIdpConfigAttrsE(t, ctx, projectID, tenantID, idpID)
	require.NoError(t, err)

	return attrs
}

// GetIdentityPlatformTenantDefaultSupportedIdpConfigAttrsE returns the settings Google Cloud holds for the given Identity Platform tenant default IdP config.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformTenantDefaultSupportedIdpConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, tenantID string, idpID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2DefaultSupportedIdpConfig, error) {
	logger.Default.Logf(t, "Getting settings for Identity Platform tenant default IdP config %s on tenant %s in project %s", idpID, tenantID, projectID)

	service, err := NewIdentityPlatformServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIdentityPlatformTenantDefaultSupportedIdpConfigAttrsWithClient(ctx, service, projectID, tenantID, idpID)
}

// GetIdentityPlatformTenantDefaultSupportedIdpConfigAttrsWithClient returns the settings Google Cloud holds for the given Identity Platform tenant default IdP config using the supplied
// *identitytoolkit.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see identitytoolkit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformTenantDefaultSupportedIdpConfigAttrsWithClient(ctx context.Context, service *identitytoolkit.Service, projectID string, tenantID string, idpID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2DefaultSupportedIdpConfig, error) {
	name := fmt.Sprintf("projects/%s/tenants/%s/defaultSupportedIdpConfigs/%s", projectID, tenantID, idpID)

	attrs, err := service.Projects.Tenants.DefaultSupportedIdpConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Identity Platform tenant default IdP config %s on tenant %s in project %s does not exist", idpID, tenantID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Identity Platform tenant default IdP config %s on tenant %s in project %s: %w", idpID, tenantID, projectID, err)
	}

	return attrs, nil
}

// GetIdentityPlatformTenantInboundSamlConfigAttrs returns the settings Google Cloud holds for the given Identity Platform tenant inbound SAML config, so a test can assert on what was
// actually created rather than only that it exists.
// A tenant federates separately from the project, so its SAML configs are their own resources.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformTenantInboundSamlConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, tenantID string, configID string) *identitytoolkit.GoogleCloudIdentitytoolkitAdminV2InboundSamlConfig {
	attrs, err := GetIdentityPlatformTenantInboundSamlConfigAttrsE(t, ctx, projectID, tenantID, configID)
	require.NoError(t, err)

	return attrs
}

// GetIdentityPlatformTenantInboundSamlConfigAttrsE returns the settings Google Cloud holds for the given Identity Platform tenant inbound SAML config.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformTenantInboundSamlConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, tenantID string, configID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2InboundSamlConfig, error) {
	logger.Default.Logf(t, "Getting settings for Identity Platform tenant inbound SAML config %s on tenant %s in project %s", configID, tenantID, projectID)

	service, err := NewIdentityPlatformServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIdentityPlatformTenantInboundSamlConfigAttrsWithClient(ctx, service, projectID, tenantID, configID)
}

// GetIdentityPlatformTenantInboundSamlConfigAttrsWithClient returns the settings Google Cloud holds for the given Identity Platform tenant inbound SAML config using the supplied
// *identitytoolkit.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see identitytoolkit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformTenantInboundSamlConfigAttrsWithClient(ctx context.Context, service *identitytoolkit.Service, projectID string, tenantID string, configID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2InboundSamlConfig, error) {
	name := fmt.Sprintf("projects/%s/tenants/%s/inboundSamlConfigs/%s", projectID, tenantID, configID)

	attrs, err := service.Projects.Tenants.InboundSamlConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Identity Platform tenant inbound SAML config %s on tenant %s in project %s does not exist", configID, tenantID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Identity Platform tenant inbound SAML config %s on tenant %s in project %s: %w", configID, tenantID, projectID, err)
	}

	return attrs, nil
}

// GetIdentityPlatformTenantOAuthIdpConfigAttrs returns the settings Google Cloud holds for the given Identity Platform tenant OAuth IdP config, so a test can assert on what was
// actually created rather than only that it exists.
// A tenant federates separately from the project, so its OIDC configs are their own resources.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformTenantOAuthIdpConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, tenantID string, configID string) *identitytoolkit.GoogleCloudIdentitytoolkitAdminV2OAuthIdpConfig {
	attrs, err := GetIdentityPlatformTenantOAuthIdpConfigAttrsE(t, ctx, projectID, tenantID, configID)
	require.NoError(t, err)

	return attrs
}

// GetIdentityPlatformTenantOAuthIdpConfigAttrsE returns the settings Google Cloud holds for the given Identity Platform tenant OAuth IdP config.
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformTenantOAuthIdpConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, tenantID string, configID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2OAuthIdpConfig, error) {
	logger.Default.Logf(t, "Getting settings for Identity Platform tenant OAuth IdP config %s on tenant %s in project %s", configID, tenantID, projectID)

	service, err := NewIdentityPlatformServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetIdentityPlatformTenantOAuthIdpConfigAttrsWithClient(ctx, service, projectID, tenantID, configID)
}

// GetIdentityPlatformTenantOAuthIdpConfigAttrsWithClient returns the settings Google Cloud holds for the given Identity Platform tenant OAuth IdP config using the supplied
// *identitytoolkit.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see identitytoolkit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetIdentityPlatformTenantOAuthIdpConfigAttrsWithClient(ctx context.Context, service *identitytoolkit.Service, projectID string, tenantID string, configID string) (*identitytoolkit.GoogleCloudIdentitytoolkitAdminV2OAuthIdpConfig, error) {
	name := fmt.Sprintf("projects/%s/tenants/%s/oauthIdpConfigs/%s", projectID, tenantID, configID)

	attrs, err := service.Projects.Tenants.OauthIdpConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Identity Platform tenant OAuth IdP config %s on tenant %s in project %s does not exist", configID, tenantID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Identity Platform tenant OAuth IdP config %s on tenant %s in project %s: %w", configID, tenantID, projectID, err)
	}

	return attrs, nil
}

// NewIdentityPlatformServiceE creates an Identity Platform service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewIdentityPlatformServiceE(t testing.TestingT, ctx context.Context) (*identitytoolkit.Service, error) {
	return identitytoolkit.NewService(ctx, append(withOptions(), option.WithScopes(identitytoolkit.CloudPlatformScope))...)
}
