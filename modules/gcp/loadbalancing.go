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

// FetchBackendBucket queries GCP to return the settings it holds for the given backend bucket, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchBackendBucket(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.BackendBucket {
	bucket, err := FetchBackendBucketE(t, ctx, projectID, name)
	require.NoError(t, err)

	return bucket
}

// FetchBackendBucketE queries GCP to return the settings it holds for the given backend bucket.
// The ctx parameter supports cancellation and timeouts.
func FetchBackendBucketE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.BackendBucket, error) {
	logger.Default.Logf(t, "Getting backend bucket %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchBackendBucketWithClient(ctx, service, projectID, name)
}

// FetchBackendBucketWithClient queries GCP to return the settings it holds for the given backend bucket using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchBackendBucketWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.BackendBucket, error) {
	bucket, err := service.BackendBuckets.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("BackendBuckets.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return bucket, nil
}

// FetchHTTPHealthCheck queries GCP to return the settings it holds for the given legacy HTTP health check, so a test can
// assert on what was actually created rather than only that it exists. A legacy health check is a different resource from the health check a modern load balancer uses, and it answers on its own endpoint.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchHTTPHealthCheck(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.HttpHealthCheck {
	healthCheck, err := FetchHTTPHealthCheckE(t, ctx, projectID, name)
	require.NoError(t, err)

	return healthCheck
}

// FetchHTTPHealthCheckE queries GCP to return the settings it holds for the given legacy HTTP health check.
// The ctx parameter supports cancellation and timeouts.
func FetchHTTPHealthCheckE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.HttpHealthCheck, error) {
	logger.Default.Logf(t, "Getting legacy HTTP health check %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchHTTPHealthCheckWithClient(ctx, service, projectID, name)
}

// FetchHTTPHealthCheckWithClient queries GCP to return the settings it holds for the given legacy HTTP health check using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchHTTPHealthCheckWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.HttpHealthCheck, error) {
	healthCheck, err := service.HttpHealthChecks.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("HttpHealthChecks.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return healthCheck, nil
}

// FetchHTTPSHealthCheck queries GCP to return the settings it holds for the given legacy HTTPS health check, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchHTTPSHealthCheck(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.HttpsHealthCheck {
	healthCheck, err := FetchHTTPSHealthCheckE(t, ctx, projectID, name)
	require.NoError(t, err)

	return healthCheck
}

// FetchHTTPSHealthCheckE queries GCP to return the settings it holds for the given legacy HTTPS health check.
// The ctx parameter supports cancellation and timeouts.
func FetchHTTPSHealthCheckE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.HttpsHealthCheck, error) {
	logger.Default.Logf(t, "Getting legacy HTTPS health check %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchHTTPSHealthCheckWithClient(ctx, service, projectID, name)
}

// FetchHTTPSHealthCheckWithClient queries GCP to return the settings it holds for the given legacy HTTPS health check using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchHTTPSHealthCheckWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.HttpsHealthCheck, error) {
	healthCheck, err := service.HttpsHealthChecks.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("HttpsHealthChecks.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return healthCheck, nil
}

// FetchSSLPolicy queries GCP to return the settings it holds for the given SSL policy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchSSLPolicy(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.SslPolicy {
	policy, err := FetchSSLPolicyE(t, ctx, projectID, name)
	require.NoError(t, err)

	return policy
}

// FetchSSLPolicyE queries GCP to return the settings it holds for the given SSL policy.
// The ctx parameter supports cancellation and timeouts.
func FetchSSLPolicyE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.SslPolicy, error) {
	logger.Default.Logf(t, "Getting SSL policy %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchSSLPolicyWithClient(ctx, service, projectID, name)
}

// FetchSSLPolicyWithClient queries GCP to return the settings it holds for the given SSL policy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchSSLPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.SslPolicy, error) {
	policy, err := service.SslPolicies.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("SslPolicies.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return policy, nil
}

// FetchTargetHTTPProxy queries GCP to return the settings it holds for the given target HTTP proxy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetHTTPProxy(t testing.TestingT, ctx context.Context, projectID string, name string) *compute.TargetHttpProxy {
	proxy, err := FetchTargetHTTPProxyE(t, ctx, projectID, name)
	require.NoError(t, err)

	return proxy
}

// FetchTargetHTTPProxyE queries GCP to return the settings it holds for the given target HTTP proxy.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetHTTPProxyE(t testing.TestingT, ctx context.Context, projectID string, name string) (*compute.TargetHttpProxy, error) {
	logger.Default.Logf(t, "Getting target HTTP proxy %s", name)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchTargetHTTPProxyWithClient(ctx, service, projectID, name)
}

// FetchTargetHTTPProxyWithClient queries GCP to return the settings it holds for the given target HTTP proxy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchTargetHTTPProxyWithClient(ctx context.Context, service *compute.Service, projectID string, name string) (*compute.TargetHttpProxy, error) {
	proxy, err := service.TargetHttpProxies.Get(projectID, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("TargetHttpProxies.Get(%s, %s) got error: %w", projectID, name, err)
	}

	return proxy, nil
}

// FetchTargetPool queries GCP to return the settings it holds for the given target pool, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetPool(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.TargetPool {
	pool, err := FetchTargetPoolE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return pool
}

// FetchTargetPoolE queries GCP to return the settings it holds for the given target pool.
// The ctx parameter supports cancellation and timeouts.
func FetchTargetPoolE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.TargetPool, error) {
	logger.Default.Logf(t, "Getting target pool %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchTargetPoolWithClient(ctx, service, projectID, region, name)
}

// FetchTargetPoolWithClient queries GCP to return the settings it holds for the given target pool using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchTargetPoolWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.TargetPool, error) {
	pool, err := service.TargetPools.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("TargetPools.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return pool, nil
}

// FetchNetworkEndpointGroup queries GCP to return the settings it holds for the given network endpoint group, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkEndpointGroup(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) *compute.NetworkEndpointGroup {
	group, err := FetchNetworkEndpointGroupE(t, ctx, projectID, zone, name)
	require.NoError(t, err)

	return group
}

// FetchNetworkEndpointGroupE queries GCP to return the settings it holds for the given network endpoint group.
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkEndpointGroupE(t testing.TestingT, ctx context.Context, projectID string, zone string, name string) (*compute.NetworkEndpointGroup, error) {
	logger.Default.Logf(t, "Getting network endpoint group %s in zone %s", name, zone)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchNetworkEndpointGroupWithClient(ctx, service, projectID, zone, name)
}

// FetchNetworkEndpointGroupWithClient queries GCP to return the settings it holds for the given network endpoint group using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkEndpointGroupWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, name string) (*compute.NetworkEndpointGroup, error) {
	group, err := service.NetworkEndpointGroups.Get(projectID, zone, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("NetworkEndpointGroups.Get(%s, %s, %s) got error: %w", projectID, zone, name, err)
	}

	return group, nil
}

// FetchRegionHealthCheck queries GCP to return the settings it holds for the given regional health check, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionHealthCheck(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.HealthCheck {
	healthCheck, err := FetchRegionHealthCheckE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return healthCheck
}

// FetchRegionHealthCheckE queries GCP to return the settings it holds for the given regional health check.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionHealthCheckE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.HealthCheck, error) {
	logger.Default.Logf(t, "Getting regional health check %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionHealthCheckWithClient(ctx, service, projectID, region, name)
}

// FetchRegionHealthCheckWithClient queries GCP to return the settings it holds for the given regional health check using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionHealthCheckWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.HealthCheck, error) {
	healthCheck, err := service.RegionHealthChecks.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionHealthChecks.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return healthCheck, nil
}

// FetchRegionBackendService queries GCP to return the settings it holds for the given regional backend service, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionBackendService(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.BackendService {
	backendService, err := FetchRegionBackendServiceE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return backendService
}

// FetchRegionBackendServiceE queries GCP to return the settings it holds for the given regional backend service.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionBackendServiceE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.BackendService, error) {
	logger.Default.Logf(t, "Getting regional backend service %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionBackendServiceWithClient(ctx, service, projectID, region, name)
}

// FetchRegionBackendServiceWithClient queries GCP to return the settings it holds for the given regional backend service using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionBackendServiceWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.BackendService, error) {
	backendService, err := service.RegionBackendServices.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionBackendServices.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return backendService, nil
}

// FetchForwardingRule queries GCP to return the settings it holds for the given regional forwarding rule, so a test can
// assert on what was actually created rather than only that it exists. A regional rule is a different resource from the global one, which FetchGlobalForwardingRule reads.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchForwardingRule(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.ForwardingRule {
	rule, err := FetchForwardingRuleE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return rule
}

// FetchForwardingRuleE queries GCP to return the settings it holds for the given regional forwarding rule.
// The ctx parameter supports cancellation and timeouts.
func FetchForwardingRuleE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.ForwardingRule, error) {
	logger.Default.Logf(t, "Getting regional forwarding rule %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchForwardingRuleWithClient(ctx, service, projectID, region, name)
}

// FetchForwardingRuleWithClient queries GCP to return the settings it holds for the given regional forwarding rule using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchForwardingRuleWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.ForwardingRule, error) {
	rule, err := service.ForwardingRules.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("ForwardingRules.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return rule, nil
}

// FetchRegionTargetHTTPSProxy queries GCP to return the settings it holds for the given regional target HTTPS proxy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionTargetHTTPSProxy(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.TargetHttpsProxy {
	proxy, err := FetchRegionTargetHTTPSProxyE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return proxy
}

// FetchRegionTargetHTTPSProxyE queries GCP to return the settings it holds for the given regional target HTTPS proxy.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionTargetHTTPSProxyE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.TargetHttpsProxy, error) {
	logger.Default.Logf(t, "Getting regional target HTTPS proxy %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionTargetHTTPSProxyWithClient(ctx, service, projectID, region, name)
}

// FetchRegionTargetHTTPSProxyWithClient queries GCP to return the settings it holds for the given regional target HTTPS proxy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionTargetHTTPSProxyWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.TargetHttpsProxy, error) {
	proxy, err := service.RegionTargetHttpsProxies.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionTargetHttpsProxies.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return proxy, nil
}

// FetchRegionTargetTCPProxy queries GCP to return the settings it holds for the given regional target TCP proxy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionTargetTCPProxy(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.TargetTcpProxy {
	proxy, err := FetchRegionTargetTCPProxyE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return proxy
}

// FetchRegionTargetTCPProxyE queries GCP to return the settings it holds for the given regional target TCP proxy.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionTargetTCPProxyE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.TargetTcpProxy, error) {
	logger.Default.Logf(t, "Getting regional target TCP proxy %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionTargetTCPProxyWithClient(ctx, service, projectID, region, name)
}

// FetchRegionTargetTCPProxyWithClient queries GCP to return the settings it holds for the given regional target TCP proxy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionTargetTCPProxyWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.TargetTcpProxy, error) {
	proxy, err := service.RegionTargetTcpProxies.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionTargetTcpProxies.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return proxy, nil
}

// FetchRegionCompositeHealthCheck queries GCP to return the settings it holds for the given regional composite health check, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionCompositeHealthCheck(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.CompositeHealthCheck {
	check, err := FetchRegionCompositeHealthCheckE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return check
}

// FetchRegionCompositeHealthCheckE queries GCP to return the settings it holds for the given regional composite health check.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionCompositeHealthCheckE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.CompositeHealthCheck, error) {
	logger.Default.Logf(t, "Getting regional composite health check %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionCompositeHealthCheckWithClient(ctx, service, projectID, region, name)
}

// FetchRegionCompositeHealthCheckWithClient queries GCP to return the settings it holds for the given regional composite health check using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionCompositeHealthCheckWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.CompositeHealthCheck, error) {
	check, err := service.RegionCompositeHealthChecks.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionCompositeHealthChecks.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return check, nil
}

// FetchRegionHealthAggregationPolicy queries GCP to return the settings it holds for the given regional health aggregation policy, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionHealthAggregationPolicy(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.HealthAggregationPolicy {
	policy, err := FetchRegionHealthAggregationPolicyE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return policy
}

// FetchRegionHealthAggregationPolicyE queries GCP to return the settings it holds for the given regional health aggregation policy.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionHealthAggregationPolicyE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.HealthAggregationPolicy, error) {
	logger.Default.Logf(t, "Getting regional health aggregation policy %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionHealthAggregationPolicyWithClient(ctx, service, projectID, region, name)
}

// FetchRegionHealthAggregationPolicyWithClient queries GCP to return the settings it holds for the given regional health aggregation policy using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionHealthAggregationPolicyWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.HealthAggregationPolicy, error) {
	policy, err := service.RegionHealthAggregationPolicies.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionHealthAggregationPolicies.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return policy, nil
}

// FetchRegionHealthSource queries GCP to return the settings it holds for the given regional health source, so a test can
// assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionHealthSource(t testing.TestingT, ctx context.Context, projectID string, region string, name string) *compute.HealthSource {
	source, err := FetchRegionHealthSourceE(t, ctx, projectID, region, name)
	require.NoError(t, err)

	return source
}

// FetchRegionHealthSourceE queries GCP to return the settings it holds for the given regional health source.
// The ctx parameter supports cancellation and timeouts.
func FetchRegionHealthSourceE(t testing.TestingT, ctx context.Context, projectID string, region string, name string) (*compute.HealthSource, error) {
	logger.Default.Logf(t, "Getting regional health source %s in region %s", name, region)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchRegionHealthSourceWithClient(ctx, service, projectID, region, name)
}

// FetchRegionHealthSourceWithClient queries GCP to return the settings it holds for the given regional health source using the
// supplied *compute.Service. Prefer this variant in unit tests where the service is backed by an
// httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchRegionHealthSourceWithClient(ctx context.Context, service *compute.Service, projectID string, region string, name string) (*compute.HealthSource, error) {
	source, err := service.RegionHealthSources.Get(projectID, region, name).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("RegionHealthSources.Get(%s, %s, %s) got error: %w", projectID, region, name, err)
	}

	return source, nil
}

// FetchNetworkEndpoints queries GCP to return the endpoints the given zonal network endpoint group
// holds, so a test can assert on what was actually attached rather than only that the group exists.
// An endpoint has no name of its own and cannot be fetched one at a time, so this lists the group's
// endpoints and the caller picks out the one it asked for.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkEndpoints(t testing.TestingT, ctx context.Context, projectID string, zone string, group string) []*compute.NetworkEndpointWithHealthStatus {
	endpoints, err := FetchNetworkEndpointsE(t, ctx, projectID, zone, group)
	require.NoError(t, err)

	return endpoints
}

// FetchNetworkEndpointsE queries GCP to return the endpoints the given zonal network endpoint group
// holds.
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkEndpointsE(t testing.TestingT, ctx context.Context, projectID string, zone string, group string) ([]*compute.NetworkEndpointWithHealthStatus, error) {
	logger.Default.Logf(t, "Getting the endpoints of network endpoint group %s in zone %s", group, zone)

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return FetchNetworkEndpointsWithClient(ctx, service, projectID, zone, group)
}

// FetchNetworkEndpointsWithClient queries GCP to return the endpoints the given zonal network
// endpoint group holds using the supplied *compute.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see loadbalancing_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func FetchNetworkEndpointsWithClient(ctx context.Context, service *compute.Service, projectID string, zone string, group string) ([]*compute.NetworkEndpointWithHealthStatus, error) {
	var endpoints []*compute.NetworkEndpointWithHealthStatus

	// The API takes a filter body even when nothing is being filtered.
	call := service.NetworkEndpointGroups.ListNetworkEndpoints(projectID, zone, group, &compute.NetworkEndpointGroupsListEndpointsRequest{})

	// A group may hold more endpoints than one page returns, and a caller asserting on a count
	// would be wrong if the rest were dropped.
	err := call.Pages(ctx, func(page *compute.NetworkEndpointGroupsListNetworkEndpoints) error {
		endpoints = append(endpoints, page.Items...)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("NetworkEndpointGroups.ListNetworkEndpoints(%s, %s, %s) got error: %w", projectID, zone, group, err)
	}

	return endpoints, nil
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
