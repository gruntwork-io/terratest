package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/networksecurity/v1"
	"google.golang.org/api/option"
)

// GetNetworkSecurityAddressGroupAttrs returns the settings Google Cloud holds for the given address group, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityAddressGroupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) *networksecurity.AddressGroup {
	group, err := GetNetworkSecurityAddressGroupAttrsE(t, ctx, projectID, location, groupID)
	require.NoError(t, err)

	return group
}

// GetNetworkSecurityAddressGroupAttrsE returns the settings Google Cloud holds for the given address group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityAddressGroupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) (*networksecurity.AddressGroup, error) {
	logger.Default.Logf(t, "Getting settings for address group %s in %s in project %s", groupID, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityAddressGroupAttrsWithClient(ctx, service, projectID, location, groupID)
}

// GetNetworkSecurityAddressGroupAttrsWithClient returns the settings Google Cloud holds for the given address group using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityAddressGroupAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, groupID string) (*networksecurity.AddressGroup, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/addressGroups/%s", projectID, location, groupID)

	group, err := service.Projects.Locations.AddressGroups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the address group %s does not exist in %s in project %s", groupID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for address group %s in %s in project %s: %w", groupID, location, projectID, err)
	}

	return group, nil
}

// GetNetworkSecurityClientTLSPolicyAttrs returns the settings Google Cloud holds for the given client TLS policy, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityClientTLSPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) *networksecurity.ClientTlsPolicy {
	policy, err := GetNetworkSecurityClientTLSPolicyAttrsE(t, ctx, projectID, location, policyID)
	require.NoError(t, err)

	return policy
}

// GetNetworkSecurityClientTLSPolicyAttrsE returns the settings Google Cloud holds for the given client TLS policy.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityClientTLSPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) (*networksecurity.ClientTlsPolicy, error) {
	logger.Default.Logf(t, "Getting settings for client TLS policy %s in %s in project %s", policyID, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityClientTLSPolicyAttrsWithClient(ctx, service, projectID, location, policyID)
}

// GetNetworkSecurityClientTLSPolicyAttrsWithClient returns the settings Google Cloud holds for the given client TLS policy using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityClientTLSPolicyAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, policyID string) (*networksecurity.ClientTlsPolicy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/clientTlsPolicies/%s", projectID, location, policyID)

	policy, err := service.Projects.Locations.ClientTlsPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the client TLS policy %s does not exist in %s in project %s", policyID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for client TLS policy %s in %s in project %s: %w", policyID, location, projectID, err)
	}

	return policy, nil
}

// GetNetworkSecurityServerTLSPolicyAttrs returns the settings Google Cloud holds for the given server TLS policy, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityServerTLSPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) *networksecurity.ServerTlsPolicy {
	policy, err := GetNetworkSecurityServerTLSPolicyAttrsE(t, ctx, projectID, location, policyID)
	require.NoError(t, err)

	return policy
}

// GetNetworkSecurityServerTLSPolicyAttrsE returns the settings Google Cloud holds for the given server TLS policy.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityServerTLSPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) (*networksecurity.ServerTlsPolicy, error) {
	logger.Default.Logf(t, "Getting settings for server TLS policy %s in %s in project %s", policyID, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityServerTLSPolicyAttrsWithClient(ctx, service, projectID, location, policyID)
}

// GetNetworkSecurityServerTLSPolicyAttrsWithClient returns the settings Google Cloud holds for the given server TLS policy using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityServerTLSPolicyAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, policyID string) (*networksecurity.ServerTlsPolicy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/serverTlsPolicies/%s", projectID, location, policyID)

	policy, err := service.Projects.Locations.ServerTlsPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the server TLS policy %s does not exist in %s in project %s", policyID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for server TLS policy %s in %s in project %s: %w", policyID, location, projectID, err)
	}

	return policy, nil
}

// GetNetworkSecurityURLListAttrs returns the settings Google Cloud holds for the given URL list, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityURLListAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, listID string) *networksecurity.UrlList {
	list, err := GetNetworkSecurityURLListAttrsE(t, ctx, projectID, location, listID)
	require.NoError(t, err)

	return list
}

// GetNetworkSecurityURLListAttrsE returns the settings Google Cloud holds for the given URL list.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityURLListAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, listID string) (*networksecurity.UrlList, error) {
	logger.Default.Logf(t, "Getting settings for URL list %s in %s in project %s", listID, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityURLListAttrsWithClient(ctx, service, projectID, location, listID)
}

// GetNetworkSecurityURLListAttrsWithClient returns the settings Google Cloud holds for the given URL list using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityURLListAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, listID string) (*networksecurity.UrlList, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/urlLists/%s", projectID, location, listID)

	list, err := service.Projects.Locations.UrlLists.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the URL list %s does not exist in %s in project %s", listID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for URL list %s in %s in project %s: %w", listID, location, projectID, err)
	}

	return list, nil
}

// GetNetworkSecurityGatewaySecurityPolicyAttrs returns the settings Google Cloud holds for the given gateway security policy, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityGatewaySecurityPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) *networksecurity.GatewaySecurityPolicy {
	policy, err := GetNetworkSecurityGatewaySecurityPolicyAttrsE(t, ctx, projectID, location, policyID)
	require.NoError(t, err)

	return policy
}

// GetNetworkSecurityGatewaySecurityPolicyAttrsE returns the settings Google Cloud holds for the given gateway security policy.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityGatewaySecurityPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) (*networksecurity.GatewaySecurityPolicy, error) {
	logger.Default.Logf(t, "Getting settings for gateway security policy %s in %s in project %s", policyID, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityGatewaySecurityPolicyAttrsWithClient(ctx, service, projectID, location, policyID)
}

// GetNetworkSecurityGatewaySecurityPolicyAttrsWithClient returns the settings Google Cloud holds for the given gateway security policy using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityGatewaySecurityPolicyAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, policyID string) (*networksecurity.GatewaySecurityPolicy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/gatewaySecurityPolicies/%s", projectID, location, policyID)

	policy, err := service.Projects.Locations.GatewaySecurityPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the gateway security policy %s does not exist in %s in project %s", policyID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for gateway security policy %s in %s in project %s: %w", policyID, location, projectID, err)
	}

	return policy, nil
}

// NewNetworkSecurityServiceE creates a Network Security service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewNetworkSecurityServiceE(t testing.TestingT, ctx context.Context) (*networksecurity.Service, error) {
	return networksecurity.NewService(ctx, append(withOptions(), option.WithScopes(networksecurity.CloudPlatformScope))...)
}

// GetNetworkSecurityAddressGroupIamPolicyAttrs returns the IAM policy Google Cloud holds for the given
// address group, so a test can assert on who was actually granted access to it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityAddressGroupIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) *networksecurity.GoogleIamV1Policy {
	policy, err := GetNetworkSecurityAddressGroupIamPolicyAttrsE(t, ctx, projectID, location, groupID)
	require.NoError(t, err)

	return policy
}

// GetNetworkSecurityAddressGroupIamPolicyAttrsE returns the IAM policy Google Cloud holds for the
// given address group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityAddressGroupIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, groupID string) (*networksecurity.GoogleIamV1Policy, error) {
	logger.Default.Logf(t, "Getting the IAM policy for address group %s in %s in project %s", groupID, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityAddressGroupIamPolicyAttrsWithClient(ctx, service, projectID, location, groupID)
}

// GetNetworkSecurityAddressGroupIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for
// the given address group using the supplied *networksecurity.Service. Prefer this variant in unit
// tests where the service is backed by an httptest fake server (see networksecurity_test.go).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityAddressGroupIamPolicyAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, groupID string) (*networksecurity.GoogleIamV1Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/addressGroups/%s", projectID, location, groupID)

	// A policy carrying a conditional binding is only returned in full at version 3, so that is what is
	// asked for: at a lower version Google drops the condition or refuses the call outright.
	policy, err := service.Projects.Locations.AddressGroups.GetIamPolicy(resource).
		OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the address group %s does not exist in %s in project %s", groupID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get the IAM policy for address group %s in %s in project %s: %w", groupID, location, projectID, err)
	}

	return policy, nil
}
