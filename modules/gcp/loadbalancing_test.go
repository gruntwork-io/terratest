package gcp_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Each response below is shaped like the one Google returns for a resource the
// terraform-google-networking load balancing modules created, not a copy of any one fixture's values.

func TestFetchHealthCheckWithClient(t *testing.T) {
	t.Parallel()

	handler := respond(t, http.MethodGet, "/projects/gw-library-test-project/global/healthChecks/gw-library-test", http.StatusOK, `{
		"kind":"compute#healthCheck",
		"name":"gw-library-test",
		"type":"HTTP",
		"checkIntervalSec":10,
		"timeoutSec":5,
		"healthyThreshold":2,
		"unhealthyThreshold":3,
		"httpHealthCheck":{"port":8080,"requestPath":"/healthz"}
	}`)

	healthCheck, err := gcp.FetchHealthCheckWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "HTTP", healthCheck.Type)
	assert.Equal(t, int64(10), healthCheck.CheckIntervalSec)
	assert.Equal(t, int64(5), healthCheck.TimeoutSec)
	require.NotNil(t, healthCheck.HttpHealthCheck)
	assert.Equal(t, int64(8080), healthCheck.HttpHealthCheck.Port)
	assert.Equal(t, "/healthz", healthCheck.HttpHealthCheck.RequestPath)
}

func TestFetchHealthCheckWithClientMissingHealthCheck(t *testing.T) {
	t.Parallel()

	_, err := gcp.FetchHealthCheckWithClient(context.Background(), newFakeComputeService(t, respond(t, "", "", http.StatusNotFound, "")), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "HealthChecks.Get(gw-library-test-project, gone)")
}

func TestFetchBackendServiceWithClient(t *testing.T) {
	t.Parallel()

	handler := respond(t, http.MethodGet, "/projects/gw-library-test-project/global/backendServices/gw-library-test", http.StatusOK, `{
		"kind":"compute#backendService",
		"name":"gw-library-test",
		"protocol":"HTTP",
		"loadBalancingScheme":"EXTERNAL_MANAGED",
		"timeoutSec":30,
		"healthChecks":["https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/healthChecks/gw-library-test"],
		"logConfig":{"enable":true,"sampleRate":0.5}
	}`)

	backendService, err := gcp.FetchBackendServiceWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "HTTP", backendService.Protocol)
	assert.Equal(t, "EXTERNAL_MANAGED", backendService.LoadBalancingScheme)
	assert.Equal(t, int64(30), backendService.TimeoutSec)
	require.Len(t, backendService.HealthChecks, 1)
	assert.True(t, strings.HasSuffix(backendService.HealthChecks[0], "/global/healthChecks/gw-library-test"))
	require.NotNil(t, backendService.LogConfig)
	assert.InDelta(t, 0.5, backendService.LogConfig.SampleRate, 1e-9)
}

func TestFetchBackendServiceWithClientMissingBackendService(t *testing.T) {
	t.Parallel()

	_, err := gcp.FetchBackendServiceWithClient(context.Background(), newFakeComputeService(t, respond(t, "", "", http.StatusNotFound, "")), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "BackendServices.Get(gw-library-test-project, gone)")
}

func TestFetchURLMapWithClient(t *testing.T) {
	t.Parallel()

	handler := respond(t, http.MethodGet, "/projects/gw-library-test-project/global/urlMaps/gw-library-test", http.StatusOK, `{
		"kind":"compute#urlMap",
		"name":"gw-library-test",
		"defaultService":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/backendServices/gw-library-test",
		"hostRules":[{"hosts":["app.example.com"],"pathMatcher":"app"}],
		"pathMatchers":[{"name":"app","defaultService":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/backendServices/gw-library-test"}]
	}`)

	urlMap, err := gcp.FetchURLMapWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.True(t, strings.HasSuffix(urlMap.DefaultService, "/global/backendServices/gw-library-test"))
	require.Len(t, urlMap.HostRules, 1)
	assert.Equal(t, []string{"app.example.com"}, urlMap.HostRules[0].Hosts)
	require.Len(t, urlMap.PathMatchers, 1)
	assert.Equal(t, "app", urlMap.PathMatchers[0].Name)
}

func TestFetchURLMapWithClientMissingURLMap(t *testing.T) {
	t.Parallel()

	_, err := gcp.FetchURLMapWithClient(context.Background(), newFakeComputeService(t, respond(t, "", "", http.StatusNotFound, "")), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "UrlMaps.Get(gw-library-test-project, gone)")
}

func TestFetchTargetHTTPSProxyWithClient(t *testing.T) {
	t.Parallel()

	handler := respond(t, http.MethodGet, "/projects/gw-library-test-project/global/targetHttpsProxies/gw-library-test", http.StatusOK, `{
		"kind":"compute#targetHttpsProxy",
		"name":"gw-library-test",
		"urlMap":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/urlMaps/gw-library-test",
		"sslCertificates":["https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/sslCertificates/gw-library-test"],
		"quicOverride":"DISABLE"
	}`)

	proxy, err := gcp.FetchTargetHTTPSProxyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.True(t, strings.HasSuffix(proxy.UrlMap, "/global/urlMaps/gw-library-test"))
	require.Len(t, proxy.SslCertificates, 1)
	assert.True(t, strings.HasSuffix(proxy.SslCertificates[0], "/global/sslCertificates/gw-library-test"))
	assert.Equal(t, "DISABLE", proxy.QuicOverride)
}

func TestFetchTargetHTTPSProxyWithClientMissingProxy(t *testing.T) {
	t.Parallel()

	_, err := gcp.FetchTargetHTTPSProxyWithClient(context.Background(), newFakeComputeService(t, respond(t, "", "", http.StatusNotFound, "")), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "TargetHttpsProxies.Get(gw-library-test-project, gone)")
}

func TestFetchSSLCertificateWithClient(t *testing.T) {
	t.Parallel()

	// A managed certificate is readable while it is still provisioning, which is the state a test
	// usually sees, so that is the state this case returns.
	handler := respond(t, http.MethodGet, "/projects/gw-library-test-project/global/sslCertificates/gw-library-test", http.StatusOK, `{
		"kind":"compute#sslCertificate",
		"name":"gw-library-test",
		"type":"MANAGED",
		"managed":{"domains":["app.example.com"],"status":"PROVISIONING","domainStatus":{"app.example.com":"PROVISIONING"}}
	}`)

	certificate, err := gcp.FetchSSLCertificateWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "MANAGED", certificate.Type)
	require.NotNil(t, certificate.Managed)
	assert.Equal(t, []string{"app.example.com"}, certificate.Managed.Domains)
	assert.Equal(t, "PROVISIONING", certificate.Managed.Status)
}

func TestFetchSSLCertificateWithClientMissingCertificate(t *testing.T) {
	t.Parallel()

	_, err := gcp.FetchSSLCertificateWithClient(context.Background(), newFakeComputeService(t, respond(t, "", "", http.StatusNotFound, "")), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "SslCertificates.Get(gw-library-test-project, gone)")
}

func TestFetchGlobalForwardingRuleWithClient(t *testing.T) {
	t.Parallel()

	handler := respond(t, http.MethodGet, "/projects/gw-library-test-project/global/forwardingRules/gw-library-test", http.StatusOK, `{
		"kind":"compute#forwardingRule",
		"name":"gw-library-test",
		"IPProtocol":"TCP",
		"portRange":"443-443",
		"loadBalancingScheme":"EXTERNAL_MANAGED",
		"target":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/targetHttpsProxies/gw-library-test",
		"labels":{"purpose":"terratest"}
	}`)

	forwardingRule, err := gcp.FetchGlobalForwardingRuleWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "TCP", forwardingRule.IPProtocol)
	assert.Equal(t, "443-443", forwardingRule.PortRange)
	assert.Equal(t, "EXTERNAL_MANAGED", forwardingRule.LoadBalancingScheme)
	assert.True(t, strings.HasSuffix(forwardingRule.Target, "/global/targetHttpsProxies/gw-library-test"))
	assert.Equal(t, "terratest", forwardingRule.Labels["purpose"])
}

func TestFetchGlobalForwardingRuleWithClientMissingRule(t *testing.T) {
	t.Parallel()

	_, err := gcp.FetchGlobalForwardingRuleWithClient(context.Background(), newFakeComputeService(t, respond(t, "", "", http.StatusNotFound, "")), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "GlobalForwardingRules.Get(gw-library-test-project, gone)")
}
