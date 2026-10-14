package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/gkehub/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetGKEHubScopeAttrs returns the settings Google Cloud holds for the given fleet scope, so a test can
// assert on the labels it pushes to the namespaces inside it. A scope groups the namespaces a fleet
// shares; it needs no cluster to exist.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEHubScopeAttrs(t testing.TestingT, ctx context.Context, projectID string, scopeID string) *gkehub.Scope {
	scope, err := GetGKEHubScopeAttrsE(t, ctx, projectID, scopeID)
	require.NoError(t, err)

	return scope
}

// GetGKEHubScopeAttrsE returns the settings Google Cloud holds for the given fleet scope.
// The ctx parameter supports cancellation and timeouts.
func GetGKEHubScopeAttrsE(t testing.TestingT, ctx context.Context, projectID string, scopeID string) (*gkehub.Scope, error) {
	logger.Default.Logf(t, "Getting settings for fleet scope %s in project %s", scopeID, projectID)

	service, err := NewGKEHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEHubScopeAttrsWithClient(ctx, service, projectID, scopeID)
}

// GetGKEHubScopeAttrsWithClient returns the settings Google Cloud holds for the given fleet scope using
// the supplied *gkehub.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see gkehub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEHubScopeAttrsWithClient(ctx context.Context, service *gkehub.Service, projectID string, scopeID string) (*gkehub.Scope, error) {
	// A fleet scope lives in the global collection: a fleet spans regions, so its scopes belong to none.
	name := fmt.Sprintf("projects/%s/locations/global/scopes/%s", projectID, scopeID)

	scope, err := service.Projects.Locations.Scopes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the fleet scope %s in project %s does not exist", scopeID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for fleet scope %s in project %s: %w", scopeID, projectID, err)
	}

	return scope, nil
}

// GetGKEHubNamespaceAttrs returns the settings Google Cloud holds for the given fleet namespace, so a
// test can assert on the scope it belongs to and the labels it carries.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEHubNamespaceAttrs(t testing.TestingT, ctx context.Context, projectID string, scopeID string, namespaceID string) *gkehub.Namespace {
	namespace, err := GetGKEHubNamespaceAttrsE(t, ctx, projectID, scopeID, namespaceID)
	require.NoError(t, err)

	return namespace
}

// GetGKEHubNamespaceAttrsE returns the settings Google Cloud holds for the given fleet namespace.
// The ctx parameter supports cancellation and timeouts.
func GetGKEHubNamespaceAttrsE(t testing.TestingT, ctx context.Context, projectID string, scopeID string, namespaceID string) (*gkehub.Namespace, error) {
	logger.Default.Logf(t, "Getting settings for fleet namespace %s in scope %s in project %s", namespaceID, scopeID, projectID)

	service, err := NewGKEHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEHubNamespaceAttrsWithClient(ctx, service, projectID, scopeID, namespaceID)
}

// GetGKEHubNamespaceAttrsWithClient returns the settings Google Cloud holds for the given fleet
// namespace using the supplied *gkehub.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see gkehub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEHubNamespaceAttrsWithClient(ctx context.Context, service *gkehub.Service, projectID string, scopeID string, namespaceID string) (*gkehub.Namespace, error) {
	name := fmt.Sprintf("projects/%s/locations/global/scopes/%s/namespaces/%s", projectID, scopeID, namespaceID)

	namespace, err := service.Projects.Locations.Scopes.Namespaces.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the fleet namespace %s in scope %s in project %s does not exist", namespaceID, scopeID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for fleet namespace %s in scope %s in project %s: %w", namespaceID, scopeID, projectID, err)
	}

	return namespace, nil
}

// GetGKEHubMembershipAttrs returns the settings Google Cloud holds for the fleet membership, so a test can assert on what was
// actually created rather than only that it exists.
// A membership is the fleet's record of one registered cluster, and the endpoint it reports is how a
// test tells which cluster was actually registered.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEHubMembershipAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, membershipID string) *gkehub.Membership {
	result, err := GetGKEHubMembershipAttrsE(t, ctx, projectID, location, membershipID)
	require.NoError(t, err)

	return result
}

// GetGKEHubMembershipAttrsE returns the settings Google Cloud holds for the fleet membership.
// The ctx parameter supports cancellation and timeouts.
func GetGKEHubMembershipAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, membershipID string) (*gkehub.Membership, error) {
	logger.Default.Logf(t, "Getting settings for fleet membership %s in %s in project %s", membershipID, location, projectID)

	service, err := NewGKEHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEHubMembershipAttrsWithClient(ctx, service, projectID, location, membershipID)
}

// GetGKEHubMembershipAttrsWithClient returns the settings Google Cloud holds for the fleet membership using the supplied
// *gkehub.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see gkehub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEHubMembershipAttrsWithClient(ctx context.Context, service *gkehub.Service, projectID string, location string, membershipID string) (*gkehub.Membership, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/memberships/%s", projectID, location, membershipID)

	result, err := service.Projects.Locations.Memberships.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the fleet membership %s in %s in project %s does not exist", membershipID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for fleet membership %s in %s in project %s: %w", membershipID, location, projectID, err)
	}

	return result, nil
}

// GetGKEHubMembershipBindingAttrs returns the settings Google Cloud holds for the fleet membership binding, so a test can assert on what was
// actually created rather than only that it exists.
// The binding is what puts a membership into a fleet scope, so the scope it names is the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetGKEHubMembershipBindingAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, membershipID string, bindingID string) *gkehub.MembershipBinding {
	result, err := GetGKEHubMembershipBindingAttrsE(t, ctx, projectID, location, membershipID, bindingID)
	require.NoError(t, err)

	return result
}

// GetGKEHubMembershipBindingAttrsE returns the settings Google Cloud holds for the fleet membership binding.
// The ctx parameter supports cancellation and timeouts.
func GetGKEHubMembershipBindingAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, membershipID string, bindingID string) (*gkehub.MembershipBinding, error) {
	logger.Default.Logf(t, "Getting settings for fleet membership binding %s of membership %s in %s in project %s", bindingID, membershipID, location, projectID)

	service, err := NewGKEHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetGKEHubMembershipBindingAttrsWithClient(ctx, service, projectID, location, membershipID, bindingID)
}

// GetGKEHubMembershipBindingAttrsWithClient returns the settings Google Cloud holds for the fleet membership binding using the supplied
// *gkehub.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see gkehub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetGKEHubMembershipBindingAttrsWithClient(ctx context.Context, service *gkehub.Service, projectID string, location string, membershipID string, bindingID string) (*gkehub.MembershipBinding, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/memberships/%s/bindings/%s", projectID, location, membershipID, bindingID)

	result, err := service.Projects.Locations.Memberships.Bindings.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the fleet membership binding %s of membership %s in %s in project %s does not exist", bindingID, membershipID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for fleet membership binding %s of membership %s in %s in project %s: %w", bindingID, membershipID, location, projectID, err)
	}

	return result, nil
}

// NewGKEHubServiceE creates a GKE Hub service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewGKEHubServiceE(t testing.TestingT, ctx context.Context) (*gkehub.Service, error) {
	return gkehub.NewService(ctx, append(withOptions(), option.WithScopes(gkehub.CloudPlatformScope))...)
}
