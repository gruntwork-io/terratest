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

func TestFetchBackendBucketWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking backend bucket module sets, because the point
	// of reading settings back is asserting a module configured the bucket it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/backendBuckets/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","bucketName":"gw-library-test-bucket","enableCdn":true}`))
	})

	bucket, err := gcp.FetchBackendBucketWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", bucket.Name)
	assert.Equal(t, "created by terratest", bucket.Description)
	assert.Equal(t, "gw-library-test-bucket", bucket.BucketName)
	assert.True(t, bucket.EnableCdn)
}

func TestFetchHTTPHealthCheckWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking HTTP health check module sets, because the point
	// of reading settings back is asserting a module configured the check it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/httpHealthChecks/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","port":8080,"requestPath":"/healthz","checkIntervalSec":15,"timeoutSec":10,"healthyThreshold":3,"unhealthyThreshold":4}`))
	})

	healthCheck, err := gcp.FetchHTTPHealthCheckWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", healthCheck.Name)
	assert.Equal(t, int64(8080), healthCheck.Port)
	assert.Equal(t, "/healthz", healthCheck.RequestPath)
	assert.Equal(t, int64(15), healthCheck.CheckIntervalSec)
	assert.Equal(t, int64(4), healthCheck.UnhealthyThreshold)
}

func TestFetchHTTPSHealthCheckWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking HTTPS health check module sets, because the point
	// of reading settings back is asserting a module configured the check it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/httpsHealthChecks/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","port":8443,"requestPath":"/healthz","checkIntervalSec":15,"timeoutSec":10}`))
	})

	healthCheck, err := gcp.FetchHTTPSHealthCheckWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", healthCheck.Name)
	assert.Equal(t, int64(8443), healthCheck.Port)
	assert.Equal(t, "/healthz", healthCheck.RequestPath)
	assert.Equal(t, int64(10), healthCheck.TimeoutSec)
}

func TestFetchSSLPolicyWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking SSL policy module sets, because the point
	// of reading settings back is asserting a module configured the policy it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/sslPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","profile":"RESTRICTED","minTlsVersion":"TLS_1_2"}`))
	})

	policy, err := gcp.FetchSSLPolicyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", policy.Name)
	assert.Equal(t, "RESTRICTED", policy.Profile)
	assert.Equal(t, "TLS_1_2", policy.MinTlsVersion)
}

func TestFetchTargetHTTPProxyWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking target HTTP proxy module sets, because the point
	// of reading settings back is asserting a module configured the proxy it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/global/targetHttpProxies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","urlMap":"https://www.googleapis.com/compute/v1/projects/gw-library-test-project/global/urlMaps/gw-library-test"}`))
	})

	proxy, err := gcp.FetchTargetHTTPProxyWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", proxy.Name)
	assert.Equal(t, "created by terratest", proxy.Description)
	assert.True(t, strings.HasSuffix(proxy.UrlMap, "/urlMaps/gw-library-test"))
}

func TestFetchTargetPoolWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking target pool module sets, because the point
	// of reading settings back is asserting a module configured the pool it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/targetPools/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","sessionAffinity":"CLIENT_IP","failoverRatio":0.5}`))
	})

	pool, err := gcp.FetchTargetPoolWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", pool.Name)
	assert.Equal(t, "CLIENT_IP", pool.SessionAffinity)
	assert.InDelta(t, 0.5, pool.FailoverRatio, 0.001)
}

func TestFetchNetworkEndpointGroupWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking network endpoint group module sets, because the point
	// of reading settings back is asserting a module configured the group it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/zones/us-central1-a/networkEndpointGroups/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","networkEndpointType":"GCE_VM_IP_PORT","defaultPort":8080,"size":0}`))
	})

	group, err := gcp.FetchNetworkEndpointGroupWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", group.Name)
	assert.Equal(t, "GCE_VM_IP_PORT", group.NetworkEndpointType)
	assert.Equal(t, int64(8080), group.DefaultPort)
}

func TestFetchRegionHealthCheckWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking region health check module sets, because the point
	// of reading settings back is asserting a module configured the check it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/healthChecks/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","type":"TCP","checkIntervalSec":15,"timeoutSec":10,"tcpHealthCheck":{"port":8080}}`))
	})

	healthCheck, err := gcp.FetchRegionHealthCheckWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", healthCheck.Name)
	assert.Equal(t, "TCP", healthCheck.Type)
	assert.Equal(t, int64(15), healthCheck.CheckIntervalSec)
	require.NotNil(t, healthCheck.TcpHealthCheck)
	assert.Equal(t, int64(8080), healthCheck.TcpHealthCheck.Port)
}

func TestFetchRegionBackendServiceWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-networking region backend service module sets, because the point
	// of reading settings back is asserting a module configured the service it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/regions/us-central1/backendServices/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","description":"created by terratest","protocol":"TCP","loadBalancingScheme":"INTERNAL","sessionAffinity":"CLIENT_IP","timeoutSec":25}`))
	})

	backendService, err := gcp.FetchRegionBackendServiceWithClient(context.Background(), newFakeComputeService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", backendService.Name)
	assert.Equal(t, "TCP", backendService.Protocol)
	assert.Equal(t, "INTERNAL", backendService.LoadBalancingScheme)
	assert.Equal(t, int64(25), backendService.TimeoutSec)
}
