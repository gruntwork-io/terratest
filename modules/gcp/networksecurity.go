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

// GetNetworkSecurityAuthzPolicyAttrs returns the settings Google Cloud holds for the given authorization policy, so a test can assert on what was
// actually created rather than only that it exists.
// The policy decides which requests a load balancer lets through, so its action and the target it attaches to are the whole point of it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityAuthzPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.AuthzPolicy {
	attrs, err := GetNetworkSecurityAuthzPolicyAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityAuthzPolicyAttrsE returns the settings Google Cloud holds for the given authorization policy.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityAuthzPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.AuthzPolicy, error) {
	logger.Default.Logf(t, "Getting settings for authorization policy %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityAuthzPolicyAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityAuthzPolicyAttrsWithClient returns the settings Google Cloud holds for the given authorization policy using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityAuthzPolicyAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.AuthzPolicy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/authzPolicies/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.AuthzPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the authorization policy %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for authorization policy %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkSecurityBackendAuthenticationConfigAttrs returns the settings Google Cloud holds for the given backend authentication config, so a test can assert on what was
// actually created rather than only that it exists.
// This is how a load balancer decides whether to trust the certificate its backend presents, so which trust config it names matters.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityBackendAuthenticationConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.BackendAuthenticationConfig {
	attrs, err := GetNetworkSecurityBackendAuthenticationConfigAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityBackendAuthenticationConfigAttrsE returns the settings Google Cloud holds for the given backend authentication config.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityBackendAuthenticationConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.BackendAuthenticationConfig, error) {
	logger.Default.Logf(t, "Getting settings for backend authentication config %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityBackendAuthenticationConfigAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityBackendAuthenticationConfigAttrsWithClient returns the settings Google Cloud holds for the given backend authentication config using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityBackendAuthenticationConfigAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.BackendAuthenticationConfig, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backendAuthenticationConfigs/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.BackendAuthenticationConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the backend authentication config %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for backend authentication config %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkSecurityDNSThreatDetectorAttrs returns the settings Google Cloud holds for the given DNS threat detector, so a test can assert on what was
// actually created rather than only that it exists.
// A detector watches DNS traffic for known-bad names, so which provider it uses and which networks it watches are what it does.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityDNSThreatDetectorAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.DnsThreatDetector {
	attrs, err := GetNetworkSecurityDNSThreatDetectorAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityDNSThreatDetectorAttrsE returns the settings Google Cloud holds for the given DNS threat detector.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityDNSThreatDetectorAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.DnsThreatDetector, error) {
	logger.Default.Logf(t, "Getting settings for DNS threat detector %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityDNSThreatDetectorAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityDNSThreatDetectorAttrsWithClient returns the settings Google Cloud holds for the given DNS threat detector using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityDNSThreatDetectorAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.DnsThreatDetector, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/dnsThreatDetectors/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.DnsThreatDetectors.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the DNS threat detector %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for DNS threat detector %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkSecurityFirewallEndpointAssociationAttrs returns the settings Google Cloud holds for the given firewall endpoint association, so a test can assert on what was
// actually created rather than only that it exists.
// An association is what puts a network behind a firewall endpoint, so the network and the endpoint it names are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityFirewallEndpointAssociationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.FirewallEndpointAssociation {
	attrs, err := GetNetworkSecurityFirewallEndpointAssociationAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityFirewallEndpointAssociationAttrsE returns the settings Google Cloud holds for the given firewall endpoint association.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityFirewallEndpointAssociationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.FirewallEndpointAssociation, error) {
	logger.Default.Logf(t, "Getting settings for firewall endpoint association %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityFirewallEndpointAssociationAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityFirewallEndpointAssociationAttrsWithClient returns the settings Google Cloud holds for the given firewall endpoint association using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityFirewallEndpointAssociationAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.FirewallEndpointAssociation, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/firewallEndpointAssociations/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.FirewallEndpointAssociations.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the firewall endpoint association %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for firewall endpoint association %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkSecurityInterceptDeploymentAttrs returns the settings Google Cloud holds for the given intercept deployment, so a test can assert on what was
// actually created rather than only that it exists.
// A deployment is the zonal forwarding rule that traffic is sent to for inspection, so which rule it names decides where that traffic goes.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityInterceptDeploymentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.InterceptDeployment {
	attrs, err := GetNetworkSecurityInterceptDeploymentAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityInterceptDeploymentAttrsE returns the settings Google Cloud holds for the given intercept deployment.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityInterceptDeploymentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.InterceptDeployment, error) {
	logger.Default.Logf(t, "Getting settings for intercept deployment %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityInterceptDeploymentAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityInterceptDeploymentAttrsWithClient returns the settings Google Cloud holds for the given intercept deployment using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityInterceptDeploymentAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.InterceptDeployment, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/interceptDeployments/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.InterceptDeployments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the intercept deployment %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for intercept deployment %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkSecurityInterceptDeploymentGroupAttrs returns the settings Google Cloud holds for the given intercept deployment group, so a test can assert on what was
// actually created rather than only that it exists.
// A group is how deployments in several zones are treated as one target, so the network it serves is what it is for.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityInterceptDeploymentGroupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.InterceptDeploymentGroup {
	attrs, err := GetNetworkSecurityInterceptDeploymentGroupAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityInterceptDeploymentGroupAttrsE returns the settings Google Cloud holds for the given intercept deployment group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityInterceptDeploymentGroupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.InterceptDeploymentGroup, error) {
	logger.Default.Logf(t, "Getting settings for intercept deployment group %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityInterceptDeploymentGroupAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityInterceptDeploymentGroupAttrsWithClient returns the settings Google Cloud holds for the given intercept deployment group using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityInterceptDeploymentGroupAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.InterceptDeploymentGroup, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/interceptDeploymentGroups/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.InterceptDeploymentGroups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the intercept deployment group %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for intercept deployment group %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkSecurityInterceptEndpointGroupAttrs returns the settings Google Cloud holds for the given intercept endpoint group, so a test can assert on what was
// actually created rather than only that it exists.
// An endpoint group is what a firewall rule points at when it wants traffic inspected, so the deployment group behind it is the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityInterceptEndpointGroupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.InterceptEndpointGroup {
	attrs, err := GetNetworkSecurityInterceptEndpointGroupAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityInterceptEndpointGroupAttrsE returns the settings Google Cloud holds for the given intercept endpoint group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityInterceptEndpointGroupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.InterceptEndpointGroup, error) {
	logger.Default.Logf(t, "Getting settings for intercept endpoint group %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityInterceptEndpointGroupAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityInterceptEndpointGroupAttrsWithClient returns the settings Google Cloud holds for the given intercept endpoint group using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityInterceptEndpointGroupAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.InterceptEndpointGroup, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/interceptEndpointGroups/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.InterceptEndpointGroups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the intercept endpoint group %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for intercept endpoint group %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkSecurityInterceptEndpointGroupAssociationAttrs returns the settings Google Cloud holds for the given intercept endpoint group association, so a test can assert on what was
// actually created rather than only that it exists.
// The association is what joins one network to an endpoint group, so without it the group inspects nothing.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityInterceptEndpointGroupAssociationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.InterceptEndpointGroupAssociation {
	attrs, err := GetNetworkSecurityInterceptEndpointGroupAssociationAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityInterceptEndpointGroupAssociationAttrsE returns the settings Google Cloud holds for the given intercept endpoint group association.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityInterceptEndpointGroupAssociationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.InterceptEndpointGroupAssociation, error) {
	logger.Default.Logf(t, "Getting settings for intercept endpoint group association %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityInterceptEndpointGroupAssociationAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityInterceptEndpointGroupAssociationAttrsWithClient returns the settings Google Cloud holds for the given intercept endpoint group association using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityInterceptEndpointGroupAssociationAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.InterceptEndpointGroupAssociation, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/interceptEndpointGroupAssociations/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.InterceptEndpointGroupAssociations.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the intercept endpoint group association %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for intercept endpoint group association %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkSecurityMirroringDeploymentAttrs returns the settings Google Cloud holds for the given mirroring deployment, so a test can assert on what was
// actually created rather than only that it exists.
// A deployment is the zonal forwarding rule mirrored traffic is sent to, so which rule it names decides where the copy goes.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityMirroringDeploymentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.MirroringDeployment {
	attrs, err := GetNetworkSecurityMirroringDeploymentAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityMirroringDeploymentAttrsE returns the settings Google Cloud holds for the given mirroring deployment.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityMirroringDeploymentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.MirroringDeployment, error) {
	logger.Default.Logf(t, "Getting settings for mirroring deployment %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityMirroringDeploymentAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityMirroringDeploymentAttrsWithClient returns the settings Google Cloud holds for the given mirroring deployment using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityMirroringDeploymentAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.MirroringDeployment, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/mirroringDeployments/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.MirroringDeployments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the mirroring deployment %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for mirroring deployment %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkSecurityMirroringDeploymentGroupAttrs returns the settings Google Cloud holds for the given mirroring deployment group, so a test can assert on what was
// actually created rather than only that it exists.
// A group is how mirroring deployments in several zones are treated as one target, so the network it serves is what it is for.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityMirroringDeploymentGroupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.MirroringDeploymentGroup {
	attrs, err := GetNetworkSecurityMirroringDeploymentGroupAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityMirroringDeploymentGroupAttrsE returns the settings Google Cloud holds for the given mirroring deployment group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityMirroringDeploymentGroupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.MirroringDeploymentGroup, error) {
	logger.Default.Logf(t, "Getting settings for mirroring deployment group %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityMirroringDeploymentGroupAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityMirroringDeploymentGroupAttrsWithClient returns the settings Google Cloud holds for the given mirroring deployment group using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityMirroringDeploymentGroupAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.MirroringDeploymentGroup, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/mirroringDeploymentGroups/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.MirroringDeploymentGroups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the mirroring deployment group %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for mirroring deployment group %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkSecurityMirroringEndpointGroupAttrs returns the settings Google Cloud holds for the given mirroring endpoint group, so a test can assert on what was
// actually created rather than only that it exists.
// An endpoint group is what a firewall rule points at when it wants traffic mirrored, so the deployment group behind it is the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityMirroringEndpointGroupAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.MirroringEndpointGroup {
	attrs, err := GetNetworkSecurityMirroringEndpointGroupAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityMirroringEndpointGroupAttrsE returns the settings Google Cloud holds for the given mirroring endpoint group.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityMirroringEndpointGroupAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.MirroringEndpointGroup, error) {
	logger.Default.Logf(t, "Getting settings for mirroring endpoint group %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityMirroringEndpointGroupAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityMirroringEndpointGroupAttrsWithClient returns the settings Google Cloud holds for the given mirroring endpoint group using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityMirroringEndpointGroupAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.MirroringEndpointGroup, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/mirroringEndpointGroups/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.MirroringEndpointGroups.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the mirroring endpoint group %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for mirroring endpoint group %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkSecurityMirroringEndpointGroupAssociationAttrs returns the settings Google Cloud holds for the given mirroring endpoint group association, so a test can assert on what was
// actually created rather than only that it exists.
// The association is what joins one network to a mirroring endpoint group, so without it nothing is mirrored.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityMirroringEndpointGroupAssociationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.MirroringEndpointGroupAssociation {
	attrs, err := GetNetworkSecurityMirroringEndpointGroupAssociationAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityMirroringEndpointGroupAssociationAttrsE returns the settings Google Cloud holds for the given mirroring endpoint group association.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityMirroringEndpointGroupAssociationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.MirroringEndpointGroupAssociation, error) {
	logger.Default.Logf(t, "Getting settings for mirroring endpoint group association %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityMirroringEndpointGroupAssociationAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityMirroringEndpointGroupAssociationAttrsWithClient returns the settings Google Cloud holds for the given mirroring endpoint group association using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityMirroringEndpointGroupAssociationAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.MirroringEndpointGroupAssociation, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/mirroringEndpointGroupAssociations/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.MirroringEndpointGroupAssociations.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the mirroring endpoint group association %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for mirroring endpoint group association %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetNetworkSecurityTLSInspectionPolicyAttrs returns the settings Google Cloud holds for the given TLS inspection policy, so a test can assert on what was
// actually created rather than only that it exists.
// The policy holds the certificate authority a gateway uses to open TLS sessions, so which pool it names decides what it can read.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityTLSInspectionPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *networksecurity.TlsInspectionPolicy {
	attrs, err := GetNetworkSecurityTLSInspectionPolicyAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetNetworkSecurityTLSInspectionPolicyAttrsE returns the settings Google Cloud holds for the given TLS inspection policy.
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityTLSInspectionPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*networksecurity.TlsInspectionPolicy, error) {
	logger.Default.Logf(t, "Getting settings for TLS inspection policy %s in %s in project %s", id, location, projectID)

	service, err := NewNetworkSecurityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNetworkSecurityTLSInspectionPolicyAttrsWithClient(ctx, service, projectID, location, id)
}

// GetNetworkSecurityTLSInspectionPolicyAttrsWithClient returns the settings Google Cloud holds for the given TLS inspection policy using the supplied
// *networksecurity.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see networksecurity_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNetworkSecurityTLSInspectionPolicyAttrsWithClient(ctx context.Context, service *networksecurity.Service, projectID string, location string, id string) (*networksecurity.TlsInspectionPolicy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/tlsInspectionPolicies/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.TlsInspectionPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the TLS inspection policy %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for TLS inspection policy %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
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

// NewNetworkSecurityServiceE creates a Network Security service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewNetworkSecurityServiceE(t testing.TestingT, ctx context.Context) (*networksecurity.Service, error) {
	return networksecurity.NewService(ctx, append(withOptions(), option.WithScopes(networksecurity.CloudPlatformScope))...)
}
