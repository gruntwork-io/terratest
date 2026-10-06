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
	"google.golang.org/api/option"
)

// GetServiceAccountAttrs returns the settings Google Cloud holds for the given service account, so
// a test can assert on what was actually created rather than only that it exists. The email is the
// full one, `<account-id>@<project>.iam.gserviceaccount.com`.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetServiceAccountAttrs(t testing.TestingT, ctx context.Context, projectID string, email string) *iam.ServiceAccount {
	account, err := GetServiceAccountAttrsE(t, ctx, projectID, email)
	require.NoError(t, err)

	return account
}

// GetServiceAccountAttrsE returns the settings Google Cloud holds for the given service account.
// The ctx parameter supports cancellation and timeouts.
func GetServiceAccountAttrsE(t testing.TestingT, ctx context.Context, projectID string, email string) (*iam.ServiceAccount, error) {
	logger.Default.Logf(t, "Getting settings for service account %s in project %s", email, projectID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetServiceAccountAttrsWithClient(ctx, service, projectID, email)
}

// GetServiceAccountAttrsWithClient returns the settings Google Cloud holds for the given service
// account using the supplied *iam.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see iam_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetServiceAccountAttrsWithClient(ctx context.Context, service *iam.Service, projectID string, email string) (*iam.ServiceAccount, error) {
	name := fmt.Sprintf("projects/%s/serviceAccounts/%s", projectID, email)

	account, err := service.Projects.ServiceAccounts.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("service account %s does not exist in project %s", email, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for service account %s in project %s: %w", email, projectID, err)
	}

	return account, nil
}

// GetWorkloadIdentityPoolAttrs returns the settings Google Cloud holds for the given workload
// identity pool, so a test can assert on what was actually created rather than only that it exists.
// The pool lives in a location, which is `global` for every pool the console creates.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkloadIdentityPoolAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string) *iam.WorkloadIdentityPool {
	pool, err := GetWorkloadIdentityPoolAttrsE(t, ctx, projectID, location, poolID)
	require.NoError(t, err)

	return pool
}

// GetWorkloadIdentityPoolAttrsE returns the settings Google Cloud holds for the given workload
// identity pool.
// The ctx parameter supports cancellation and timeouts.
func GetWorkloadIdentityPoolAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string) (*iam.WorkloadIdentityPool, error) {
	logger.Default.Logf(t, "Getting settings for workload identity pool %s in project %s", poolID, projectID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkloadIdentityPoolAttrsWithClient(ctx, service, projectID, location, poolID)
}

// GetWorkloadIdentityPoolAttrsWithClient returns the settings Google Cloud holds for the given
// workload identity pool using the supplied *iam.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see iam_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkloadIdentityPoolAttrsWithClient(ctx context.Context, service *iam.Service, projectID string, location string, poolID string) (*iam.WorkloadIdentityPool, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/workloadIdentityPools/%s", projectID, location, poolID)

	pool, err := service.Projects.Locations.WorkloadIdentityPools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("workload identity pool %s does not exist in project %s", poolID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for workload identity pool %s in project %s: %w", poolID, projectID, err)
	}

	return pool, nil
}

// GetWorkloadIdentityPoolProviderAttrs returns the settings Google Cloud holds for the given
// workload identity pool provider, so a test can assert on what was actually created rather than
// only that it exists. A provider is named by its pool as well as its own id.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkloadIdentityPoolProviderAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string, providerID string) *iam.WorkloadIdentityPoolProvider {
	provider, err := GetWorkloadIdentityPoolProviderAttrsE(t, ctx, projectID, location, poolID, providerID)
	require.NoError(t, err)

	return provider
}

// GetWorkloadIdentityPoolProviderAttrsE returns the settings Google Cloud holds for the given
// workload identity pool provider.
// The ctx parameter supports cancellation and timeouts.
func GetWorkloadIdentityPoolProviderAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string, providerID string) (*iam.WorkloadIdentityPoolProvider, error) {
	logger.Default.Logf(t, "Getting settings for workload identity pool provider %s in pool %s in project %s", providerID, poolID, projectID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkloadIdentityPoolProviderAttrsWithClient(ctx, service, projectID, location, poolID, providerID)
}

// GetWorkloadIdentityPoolProviderAttrsWithClient returns the settings Google Cloud holds for the
// given workload identity pool provider using the supplied *iam.Service. Prefer this variant in unit
// tests where the service is backed by an httptest fake server (see iam_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkloadIdentityPoolProviderAttrsWithClient(ctx context.Context, service *iam.Service, projectID string, location string, poolID string, providerID string) (*iam.WorkloadIdentityPoolProvider, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/workloadIdentityPools/%s/providers/%s", projectID, location, poolID, providerID)

	provider, err := service.Projects.Locations.WorkloadIdentityPools.Providers.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("workload identity pool provider %s does not exist in pool %s in project %s", providerID, poolID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for workload identity pool provider %s in pool %s in project %s: %w", providerID, poolID, projectID, err)
	}

	return provider, nil
}

// GetWorkforcePoolAttrs returns the settings Google Cloud holds for the given workforce pool, so a test can assert on what was
// actually created rather than only that it exists.
// A pool is where an outside identity provider's users are represented, so the parent it belongs to and whether it is disabled decide who can sign in.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolAttrs(t testing.TestingT, ctx context.Context, poolID string) *iam.WorkforcePool {
	attrs, err := GetWorkforcePoolAttrsE(t, ctx, poolID)
	require.NoError(t, err)

	return attrs
}

// GetWorkforcePoolAttrsE returns the settings Google Cloud holds for the given workforce pool.
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolAttrsE(t testing.TestingT, ctx context.Context, poolID string) (*iam.WorkforcePool, error) {
	logger.Default.Logf(t, "Getting settings for workforce pool %s", poolID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkforcePoolAttrsWithClient(ctx, service, poolID)
}

// GetWorkforcePoolAttrsWithClient returns the settings Google Cloud holds for the given workforce pool using the supplied
// *iam.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see iam_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolAttrsWithClient(ctx context.Context, service *iam.Service, poolID string) (*iam.WorkforcePool, error) {
	name := "locations/global/workforcePools/" + poolID

	attrs, err := service.Locations.WorkforcePools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the workforce pool %s does not exist", poolID)
		}

		return nil, fmt.Errorf("failed to get settings for workforce pool %s: %w", poolID, err)
	}

	return attrs, nil
}

// GetWorkforcePoolProviderAttrs returns the settings Google Cloud holds for the given workforce pool provider, so a test can assert on what was
// actually created rather than only that it exists.
// The provider holds the protocol and the attribute mapping, so what it maps and whether it is disabled decide who gets in and as what.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolProviderAttrs(t testing.TestingT, ctx context.Context, poolID string, providerID string) *iam.WorkforcePoolProvider {
	attrs, err := GetWorkforcePoolProviderAttrsE(t, ctx, poolID, providerID)
	require.NoError(t, err)

	return attrs
}

// GetWorkforcePoolProviderAttrsE returns the settings Google Cloud holds for the given workforce pool provider.
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolProviderAttrsE(t testing.TestingT, ctx context.Context, poolID string, providerID string) (*iam.WorkforcePoolProvider, error) {
	logger.Default.Logf(t, "Getting settings for workforce pool provider %s in pool %s", providerID, poolID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkforcePoolProviderAttrsWithClient(ctx, service, poolID, providerID)
}

// GetWorkforcePoolProviderAttrsWithClient returns the settings Google Cloud holds for the given workforce pool provider using the supplied
// *iam.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see iam_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolProviderAttrsWithClient(ctx context.Context, service *iam.Service, poolID string, providerID string) (*iam.WorkforcePoolProvider, error) {
	name := fmt.Sprintf("locations/global/workforcePools/%s/providers/%s", poolID, providerID)

	attrs, err := service.Locations.WorkforcePools.Providers.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the workforce pool provider %s in pool %s does not exist", providerID, poolID)
		}

		return nil, fmt.Errorf("failed to get settings for workforce pool provider %s in pool %s: %w", providerID, poolID, err)
	}

	return attrs, nil
}

// GetWorkforcePoolProviderKeyAttrs returns the settings Google Cloud holds for the given workforce pool provider key, so a test can assert on what was
// actually created rather than only that it exists.
// The key is what verifies the provider's assertions, so its use and whether it has expired decide whether a sign-in is accepted.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolProviderKeyAttrs(t testing.TestingT, ctx context.Context, poolID string, providerID string, keyID string) *iam.WorkforcePoolProviderKey {
	attrs, err := GetWorkforcePoolProviderKeyAttrsE(t, ctx, poolID, providerID, keyID)
	require.NoError(t, err)

	return attrs
}

// GetWorkforcePoolProviderKeyAttrsE returns the settings Google Cloud holds for the given workforce pool provider key.
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolProviderKeyAttrsE(t testing.TestingT, ctx context.Context, poolID string, providerID string, keyID string) (*iam.WorkforcePoolProviderKey, error) {
	logger.Default.Logf(t, "Getting settings for workforce pool provider key %s on provider %s in pool %s", keyID, providerID, poolID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkforcePoolProviderKeyAttrsWithClient(ctx, service, poolID, providerID, keyID)
}

// GetWorkforcePoolProviderKeyAttrsWithClient returns the settings Google Cloud holds for the given workforce pool provider key using the supplied
// *iam.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see iam_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolProviderKeyAttrsWithClient(ctx context.Context, service *iam.Service, poolID string, providerID string, keyID string) (*iam.WorkforcePoolProviderKey, error) {
	name := fmt.Sprintf("locations/global/workforcePools/%s/providers/%s/keys/%s", poolID, providerID, keyID)

	attrs, err := service.Locations.WorkforcePools.Providers.Keys.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the workforce pool provider key %s on provider %s in pool %s does not exist", keyID, providerID, poolID)
		}

		return nil, fmt.Errorf("failed to get settings for workforce pool provider key %s on provider %s in pool %s: %w", keyID, providerID, poolID, err)
	}

	return attrs, nil
}

// GetWorkforcePoolProviderScimTenantAttrs returns the settings Google Cloud holds for the given workforce pool provider SCIM tenant, so a test can assert on what was
// actually created rather than only that it exists.
// The tenant is what an outside directory pushes users into, so its claim mapping and whether it is disabled decide which accounts appear.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolProviderScimTenantAttrs(t testing.TestingT, ctx context.Context, poolID string, providerID string, tenantID string) *iam.WorkforcePoolProviderScimTenant {
	attrs, err := GetWorkforcePoolProviderScimTenantAttrsE(t, ctx, poolID, providerID, tenantID)
	require.NoError(t, err)

	return attrs
}

// GetWorkforcePoolProviderScimTenantAttrsE returns the settings Google Cloud holds for the given workforce pool provider SCIM tenant.
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolProviderScimTenantAttrsE(t testing.TestingT, ctx context.Context, poolID string, providerID string, tenantID string) (*iam.WorkforcePoolProviderScimTenant, error) {
	logger.Default.Logf(t, "Getting settings for workforce pool provider SCIM tenant %s on provider %s in pool %s", tenantID, providerID, poolID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkforcePoolProviderScimTenantAttrsWithClient(ctx, service, poolID, providerID, tenantID)
}

// GetWorkforcePoolProviderScimTenantAttrsWithClient returns the settings Google Cloud holds for the given workforce pool provider SCIM tenant using the supplied
// *iam.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see iam_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolProviderScimTenantAttrsWithClient(ctx context.Context, service *iam.Service, poolID string, providerID string, tenantID string) (*iam.WorkforcePoolProviderScimTenant, error) {
	name := fmt.Sprintf("locations/global/workforcePools/%s/providers/%s/scimTenants/%s", poolID, providerID, tenantID)

	attrs, err := service.Locations.WorkforcePools.Providers.ScimTenants.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the workforce pool provider SCIM tenant %s on provider %s in pool %s does not exist", tenantID, providerID, poolID)
		}

		return nil, fmt.Errorf("failed to get settings for workforce pool provider SCIM tenant %s on provider %s in pool %s: %w", tenantID, providerID, poolID, err)
	}

	return attrs, nil
}

// GetWorkforcePoolProviderScimTokenAttrs returns the settings Google Cloud holds for the given workforce pool provider SCIM token, so a test can assert on what was
// actually created rather than only that it exists.
// The token is the credential the outside directory authenticates with, so its state says whether that directory can still push.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolProviderScimTokenAttrs(t testing.TestingT, ctx context.Context, poolID string, providerID string, tenantID string, tokenID string) *iam.WorkforcePoolProviderScimToken {
	attrs, err := GetWorkforcePoolProviderScimTokenAttrsE(t, ctx, poolID, providerID, tenantID, tokenID)
	require.NoError(t, err)

	return attrs
}

// GetWorkforcePoolProviderScimTokenAttrsE returns the settings Google Cloud holds for the given workforce pool provider SCIM token.
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolProviderScimTokenAttrsE(t testing.TestingT, ctx context.Context, poolID string, providerID string, tenantID string, tokenID string) (*iam.WorkforcePoolProviderScimToken, error) {
	logger.Default.Logf(t, "Getting settings for workforce pool provider SCIM token %s on tenant %s on provider %s in pool %s", tokenID, tenantID, providerID, poolID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkforcePoolProviderScimTokenAttrsWithClient(ctx, service, poolID, providerID, tenantID, tokenID)
}

// GetWorkforcePoolProviderScimTokenAttrsWithClient returns the settings Google Cloud holds for the given workforce pool provider SCIM token using the supplied
// *iam.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see iam_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkforcePoolProviderScimTokenAttrsWithClient(ctx context.Context, service *iam.Service, poolID string, providerID string, tenantID string, tokenID string) (*iam.WorkforcePoolProviderScimToken, error) {
	name := fmt.Sprintf("locations/global/workforcePools/%s/providers/%s/scimTenants/%s/tokens/%s", poolID, providerID, tenantID, tokenID)

	attrs, err := service.Locations.WorkforcePools.Providers.ScimTenants.Tokens.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the workforce pool provider SCIM token %s on tenant %s on provider %s in pool %s does not exist", tokenID, tenantID, providerID, poolID)
		}

		return nil, fmt.Errorf("failed to get settings for workforce pool provider SCIM token %s on tenant %s on provider %s in pool %s: %w", tokenID, tenantID, providerID, poolID, err)
	}

	return attrs, nil
}

// GetWorkloadIdentityPoolManagedIdentityAttrs returns the settings Google Cloud holds for the given
// managed identity, so a test can assert on what was actually created rather than only that it
// exists. A managed identity sits in a namespace inside a workload identity pool, so it is named by
// both as well as by itself.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkloadIdentityPoolManagedIdentityAttrs(t testing.TestingT, ctx context.Context, projectID string, poolID string, namespaceID string, identityID string) *iam.WorkloadIdentityPoolManagedIdentity {
	identity, err := GetWorkloadIdentityPoolManagedIdentityAttrsE(t, ctx, projectID, poolID, namespaceID, identityID)
	require.NoError(t, err)

	return identity
}

// GetWorkloadIdentityPoolManagedIdentityAttrsE returns the settings Google Cloud holds for the given
// managed identity.
// The ctx parameter supports cancellation and timeouts.
func GetWorkloadIdentityPoolManagedIdentityAttrsE(t testing.TestingT, ctx context.Context, projectID string, poolID string, namespaceID string, identityID string) (*iam.WorkloadIdentityPoolManagedIdentity, error) {
	logger.Default.Logf(t, "Getting settings for managed identity %s in namespace %s in pool %s in project %s", identityID, namespaceID, poolID, projectID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkloadIdentityPoolManagedIdentityAttrsWithClient(ctx, service, projectID, poolID, namespaceID, identityID)
}

// GetWorkloadIdentityPoolManagedIdentityAttrsWithClient returns the settings Google Cloud holds for
// the given managed identity using the supplied *iam.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see iam_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkloadIdentityPoolManagedIdentityAttrsWithClient(ctx context.Context, service *iam.Service, projectID string, poolID string, namespaceID string, identityID string) (*iam.WorkloadIdentityPoolManagedIdentity, error) {
	name := fmt.Sprintf("projects/%s/locations/global/workloadIdentityPools/%s/namespaces/%s/managedIdentities/%s", projectID, poolID, namespaceID, identityID)

	identity, err := service.Projects.Locations.WorkloadIdentityPools.Namespaces.ManagedIdentities.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the managed identity %s does not exist in namespace %s in pool %s in project %s", identityID, namespaceID, poolID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for managed identity %s in namespace %s in pool %s in project %s: %w", identityID, namespaceID, poolID, projectID, err)
	}

	return identity, nil
}

// GetOAuthClientCredentialAttrs returns the settings Google Cloud holds for the given OAuth client
// credential, so a test can assert on what was actually created rather than only that it exists. The
// secret is cleared before the credential is returned, so a test that logs what it read cannot put a
// working credential in a log.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetOAuthClientCredentialAttrs(t testing.TestingT, ctx context.Context, projectID string, clientID string, credentialID string) *iam.OauthClientCredential {
	credential, err := GetOAuthClientCredentialAttrsE(t, ctx, projectID, clientID, credentialID)
	require.NoError(t, err)

	return credential
}

// GetOAuthClientCredentialAttrsE returns the settings Google Cloud holds for the given OAuth client
// credential.
// The ctx parameter supports cancellation and timeouts.
func GetOAuthClientCredentialAttrsE(t testing.TestingT, ctx context.Context, projectID string, clientID string, credentialID string) (*iam.OauthClientCredential, error) {
	logger.Default.Logf(t, "Getting settings for credential %s on OAuth client %s in project %s", credentialID, clientID, projectID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetOAuthClientCredentialAttrsWithClient(ctx, service, projectID, clientID, credentialID)
}

// GetOAuthClientCredentialAttrsWithClient returns the settings Google Cloud holds for the given OAuth
// client credential using the supplied *iam.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see iam_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetOAuthClientCredentialAttrsWithClient(ctx context.Context, service *iam.Service, projectID string, clientID string, credentialID string) (*iam.OauthClientCredential, error) {
	name := fmt.Sprintf("projects/%s/locations/global/oauthClients/%s/credentials/%s", projectID, clientID, credentialID)

	credential, err := service.Projects.Locations.OauthClients.Credentials.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the credential %s does not exist on OAuth client %s in project %s", credentialID, clientID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for credential %s on OAuth client %s in project %s: %w", credentialID, clientID, projectID, err)
	}

	// Google may answer with the secret, and nothing a test asserts needs it, so it does not leave
	// this function.
	credential.ClientSecret = ""

	return credential, nil
}

// NewIAMServiceE creates an IAM service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewIAMServiceE(t testing.TestingT, ctx context.Context) (*iam.Service, error) {
	return iam.NewService(ctx, append(withOptions(), option.WithScopes(iam.CloudPlatformScope))...)
}
