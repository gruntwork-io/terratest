package gcp

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/compute/v1"
)

// FetchHealthCheck queries GCP to return the settings it holds for the given global health
// check, so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchHealthCheck(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.HealthCheck {
	healthCheck, err := FetchHealthCheckE(t, ctx, projectID, name)
	require.NoError(t, err)

	return healthCheck
}

// FetchHealthCheckE queries GCP to return the settings it holds for the given global health
// check.
// The ctx parameter supports cancellation and timeouts.
func FetchHealthCheckE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.HealthCheck, error) {
	logger.Default.Logf(t, "Getting global health check %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchHealthCheckWithClient(ctx, service, projectID, name)
}

// FetchHealthCheckWithClient queries GCP to return the settings it holds for the given global
// health check using the supplied *compute.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchHealthCheckWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.HealthCheck, error) {
	healthCheck, err := service.HealthChecks.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("HealthChecks.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return healthCheck, nil
}

// FetchBackendService queries GCP to return the settings it holds for the given global
// backend service, so a test can assert on what was actually created rather than only that it
// exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchBackendService(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.BackendService {
	backendService, err := FetchBackendServiceE(t, ctx, projectID, name)
	require.NoError(t, err)

	return backendService
}

// FetchBackendServiceE queries GCP to return the settings it holds for the given global
// backend service.
// The ctx parameter supports cancellation and timeouts.
func FetchBackendServiceE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.BackendService, error) {
	logger.Default.Logf(t, "Getting global backend service %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchBackendServiceWithClient(ctx, service, projectID, name)
}

// FetchBackendServiceWithClient queries GCP to return the settings it holds for the given global
// backend service using the supplied *compute.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchBackendServiceWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.BackendService, error) {
	backendService, err := service.BackendServices.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("BackendServices.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return backendService, nil
}

// FetchURLMap queries GCP to return the settings it holds for the given global URL map, so a
// test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchURLMap(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.UrlMap {
	urlMap, err := FetchURLMapE(t, ctx, projectID, name)
	require.NoError(t, err)

	return urlMap
}

// FetchURLMapE queries GCP to return the settings it holds for the given global URL map.
// The ctx parameter supports cancellation and timeouts.
func FetchURLMapE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.UrlMap, error) {
	logger.Default.Logf(t, "Getting global URL map %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchURLMapWithClient(ctx, service, projectID, name)
}

// FetchURLMapWithClient queries GCP to return the settings it holds for the given global URL map
// using the supplied *compute.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchURLMapWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.UrlMap, error) {
	urlMap, err := service.UrlMaps.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("UrlMaps.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return urlMap, nil
}

// FetchTargetHTTPSProxy queries GCP to return the settings it holds for the given global
// target HTTPS proxy, so a test can assert on what was actually created rather than only that it
// exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetHTTPSProxy(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.TargetHttpsProxy {
	proxy, err := FetchTargetHTTPSProxyE(t, ctx, projectID, name)
	require.NoError(t, err)

	return proxy
}

// FetchTargetHTTPSProxyE queries GCP to return the settings it holds for the given global
// target HTTPS proxy.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetHTTPSProxyE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.TargetHttpsProxy, error) {
	logger.Default.Logf(t, "Getting global target HTTPS proxy %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchTargetHTTPSProxyWithClient(ctx, service, projectID, name)
}

// FetchTargetHTTPSProxyWithClient queries GCP to return the settings it holds for the given global
// target HTTPS proxy using the supplied *compute.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchTargetHTTPSProxyWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.TargetHttpsProxy, error) {
	proxy, err := service.TargetHttpsProxies.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("TargetHttpsProxies.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return proxy, nil
}

// FetchSSLCertificate queries GCP to return the settings it holds for the given global SSL
// certificate, so a test can assert on what was actually created rather than only that it exists. A
// Google-managed certificate is returned while it is still provisioning, with its domains' states
// in the Managed field.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchSSLCertificate(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.SslCertificate {
	certificate, err := FetchSSLCertificateE(t, ctx, projectID, name)
	require.NoError(t, err)

	return certificate
}

// FetchSSLCertificateE queries GCP to return the settings it holds for the given global SSL
// certificate.
// The ctx parameter supports cancellation and timeouts.
func FetchSSLCertificateE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.SslCertificate, error) {
	logger.Default.Logf(t, "Getting global SSL certificate %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchSSLCertificateWithClient(ctx, service, projectID, name)
}

// FetchSSLCertificateWithClient queries GCP to return the settings it holds for the given global
// SSL certificate using the supplied *compute.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchSSLCertificateWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.SslCertificate, error) {
	certificate, err := service.SslCertificates.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("SslCertificates.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return certificate, nil
}

// FetchGlobalForwardingRule queries GCP to return the settings it holds for the given global
// forwarding rule, so a test can assert on what was actually created rather than only that it
// exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchGlobalForwardingRule(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.ForwardingRule {
	forwardingRule, err := FetchGlobalForwardingRuleE(t, ctx, projectID, name)
	require.NoError(t, err)

	return forwardingRule
}

// FetchGlobalForwardingRuleE queries GCP to return the settings it holds for the given
// global forwarding rule.
// The ctx parameter supports cancellation and timeouts.
func FetchGlobalForwardingRuleE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.ForwardingRule, error) {
	logger.Default.Logf(t, "Getting global forwarding rule %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchGlobalForwardingRuleWithClient(ctx, service, projectID, name)
}

// FetchGlobalForwardingRuleWithClient queries GCP to return the settings it holds for the given
// global forwarding rule using the supplied *compute.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see loadbalancing_test.go for the
// pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchGlobalForwardingRuleWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.ForwardingRule, error) {
	forwardingRule, err := service.GlobalForwardingRules.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("GlobalForwardingRules.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return forwardingRule, nil
}

// FetchRegionSSLPolicy queries GCP to return the settings it holds for the given regional SSL policy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionSSLPolicy(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.SslPolicy {
	policy, err := FetchRegionSSLPolicyE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return policy
}

// FetchRegionSSLPolicyE queries GCP to return the settings it holds for the given regional SSL policy.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionSSLPolicyE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.SslPolicy, error) {
	logger.Default.Logf(t, "Getting regional SSL policy %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionSSLPolicyWithClient(ctx, service, projectID, region, name)
}

// FetchRegionSSLPolicyWithClient queries GCP to return the settings it holds for the given regional SSL policy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionSSLPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.SslPolicy, error) {
	policy, err := service.RegionSslPolicies.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionSslPolicies.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return policy, nil
}

// FetchRegionSSLCertificate queries GCP to return the settings it holds for the given regional SSL certificate, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionSSLCertificate(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.SslCertificate {
	certificate, err := FetchRegionSSLCertificateE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return certificate
}

// FetchRegionSSLCertificateE queries GCP to return the settings it holds for the given regional SSL certificate.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionSSLCertificateE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.SslCertificate, error) {
	logger.Default.Logf(t, "Getting regional SSL certificate %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionSSLCertificateWithClient(ctx, service, projectID, region, name)
}

// FetchRegionSSLCertificateWithClient queries GCP to return the settings it holds for the given regional SSL certificate using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionSSLCertificateWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.SslCertificate, error) {
	certificate, err := service.RegionSslCertificates.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionSslCertificates.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return certificate, nil
}

// FetchRegionURLMap queries GCP to return the settings it holds for the given regional URL map, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionURLMap(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.UrlMap {
	urlMap, err := FetchRegionURLMapE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return urlMap
}

// FetchRegionURLMapE queries GCP to return the settings it holds for the given regional URL map.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionURLMapE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.UrlMap, error) {
	logger.Default.Logf(t, "Getting regional URL map %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionURLMapWithClient(ctx, service, projectID, region, name)
}

// FetchRegionURLMapWithClient queries GCP to return the settings it holds for the given regional URL map using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionURLMapWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.UrlMap, error) {
	urlMap, err := service.RegionUrlMaps.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionUrlMaps.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return urlMap, nil
}

// FetchRegionTargetHTTPProxy queries GCP to return the settings it holds for the given regional target HTTP proxy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionTargetHTTPProxy(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.TargetHttpProxy {
	proxy, err := FetchRegionTargetHTTPProxyE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return proxy
}

// FetchRegionTargetHTTPProxyE queries GCP to return the settings it holds for the given regional target HTTP proxy.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionTargetHTTPProxyE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.TargetHttpProxy, error) {
	logger.Default.Logf(t, "Getting regional target HTTP proxy %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionTargetHTTPProxyWithClient(ctx, service, projectID, region, name)
}

// FetchRegionTargetHTTPProxyWithClient queries GCP to return the settings it holds for the given regional target HTTP proxy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionTargetHTTPProxyWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.TargetHttpProxy, error) {
	proxy, err := service.RegionTargetHttpProxies.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionTargetHttpProxies.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return proxy, nil
}

// FetchGlobalNetworkEndpointGroup queries GCP to return the settings it holds for the given global network endpoint group, so a test can
// assert on what was actually created rather than only that it exists. A global group holds endpoints outside Google Cloud, so it has no zone of its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchGlobalNetworkEndpointGroup(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.NetworkEndpointGroup {
	group, err := FetchGlobalNetworkEndpointGroupE(t, ctx, projectID, name)
	require.NoError(t, err)

	return group
}

// FetchGlobalNetworkEndpointGroupE queries GCP to return the settings it holds for the given global network endpoint group.
// The ctx parameter supports cancellation and timeouts.
func FetchGlobalNetworkEndpointGroupE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.NetworkEndpointGroup, error) {
	logger.Default.Logf(t, "Getting global network endpoint group %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchGlobalNetworkEndpointGroupWithClient(ctx, service, projectID, name)
}

// FetchGlobalNetworkEndpointGroupWithClient queries GCP to return the settings it holds for the given global network endpoint group using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchGlobalNetworkEndpointGroupWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.NetworkEndpointGroup, error) {
	group, err := service.GlobalNetworkEndpointGroups.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("GlobalNetworkEndpointGroups.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return group, nil
}

// FetchRegionNetworkEndpointGroup queries GCP to return the settings it holds for the given regional network endpoint group, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionNetworkEndpointGroup(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.NetworkEndpointGroup {
	group, err := FetchRegionNetworkEndpointGroupE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return group
}

// FetchRegionNetworkEndpointGroupE queries GCP to return the settings it holds for the given regional network endpoint group.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionNetworkEndpointGroupE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.NetworkEndpointGroup, error) {
	logger.Default.Logf(t, "Getting regional network endpoint group %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionNetworkEndpointGroupWithClient(ctx, service, projectID, region, name)
}

// FetchRegionNetworkEndpointGroupWithClient queries GCP to return the settings it holds for the given regional network endpoint group using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionNetworkEndpointGroupWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.NetworkEndpointGroup, error) {
	group, err := service.RegionNetworkEndpointGroups.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionNetworkEndpointGroups.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return group, nil
}

// FetchTargetTCPProxy queries GCP to return the settings it holds for the given target TCP proxy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetTCPProxy(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.TargetTcpProxy {
	proxy, err := FetchTargetTCPProxyE(t, ctx, projectID, name)
	require.NoError(t, err)

	return proxy
}

// FetchTargetTCPProxyE queries GCP to return the settings it holds for the given target TCP proxy.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetTCPProxyE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.TargetTcpProxy, error) {
	logger.Default.Logf(t, "Getting target TCP proxy %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchTargetTCPProxyWithClient(ctx, service, projectID, name)
}

// FetchTargetTCPProxyWithClient queries GCP to return the settings it holds for the given target TCP proxy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchTargetTCPProxyWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.TargetTcpProxy, error) {
	proxy, err := service.TargetTcpProxies.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("TargetTcpProxies.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return proxy, nil
}

// FetchTargetSSLProxy queries GCP to return the settings it holds for the given target SSL proxy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetSSLProxy(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.TargetSslProxy {
	proxy, err := FetchTargetSSLProxyE(t, ctx, projectID, name)
	require.NoError(t, err)

	return proxy
}

// FetchTargetSSLProxyE queries GCP to return the settings it holds for the given target SSL proxy.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetSSLProxyE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.TargetSslProxy, error) {
	logger.Default.Logf(t, "Getting target SSL proxy %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchTargetSSLProxyWithClient(ctx, service, projectID, name)
}

// FetchTargetSSLProxyWithClient queries GCP to return the settings it holds for the given target SSL proxy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchTargetSSLProxyWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.TargetSslProxy, error) {
	proxy, err := service.TargetSslProxies.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("TargetSslProxies.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return proxy, nil
}

// FetchTargetGRPCProxy queries GCP to return the settings it holds for the given target gRPC proxy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetGRPCProxy(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.TargetGrpcProxy {
	proxy, err := FetchTargetGRPCProxyE(t, ctx, projectID, name)
	require.NoError(t, err)

	return proxy
}

// FetchTargetGRPCProxyE queries GCP to return the settings it holds for the given target gRPC proxy.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetGRPCProxyE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.TargetGrpcProxy, error) {
	logger.Default.Logf(t, "Getting target gRPC proxy %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchTargetGRPCProxyWithClient(ctx, service, projectID, name)
}

// FetchTargetGRPCProxyWithClient queries GCP to return the settings it holds for the given target gRPC proxy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchTargetGRPCProxyWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.TargetGrpcProxy, error) {
	proxy, err := service.TargetGrpcProxies.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("TargetGrpcProxies.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return proxy, nil
}

// FetchTargetInstance queries GCP to return the settings it holds for the given target instance, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetInstance(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *compute.TargetInstance {
	target, err := FetchTargetInstanceE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return target
}

// FetchTargetInstanceE queries GCP to return the settings it holds for the given target instance.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetInstanceE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*compute.TargetInstance, error) {
	logger.Default.Logf(t, "Getting target instance %s in zone %s", name, zone)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchTargetInstanceWithClient(ctx, service, projectID, zone, name)
}

// FetchTargetInstanceWithClient queries GCP to return the settings it holds for the given target instance using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchTargetInstanceWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*compute.TargetInstance, error) {
	target, err := service.TargetInstances.Get(projectID, zone, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("TargetInstances.Get(%s, %s, %s) got error: %w", projectID, zone, name, err)
	}

	return target, nil
}
