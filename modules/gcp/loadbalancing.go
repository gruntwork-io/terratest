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
