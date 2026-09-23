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
