package gcp_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/appengine/v1"
	"google.golang.org/api/option"
)

// newFakeAppEngineService points a real App Engine Admin client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeAppEngineService(t *testing.T, handler http.Handler) *appengine.APIService {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := appengine.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetAppEngineDomainMappingAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The records Google wants the domain's owner to publish are reported on the mapping, and a create does not list
	// them.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/apps/gw-library-test-project/domainMappings/terratest.example.com"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"apps/gw-library-test-project/domainMappings/terratest.example.com","id":"terratest.example.com","sslSettings":{"sslManagementType":"AUTOMATIC"},"resourceRecords":[{"name":"terratest","rrdata":"ghs.googlehosted.com.","type":"CNAME"}]}`))
	})

	mapping, err := gcp.GetAppEngineDomainMappingAttrsWithClient(context.Background(), newFakeAppEngineService(t, handler), "gw-library-test-project", "terratest.example.com")
	require.NoError(t, err)

	assert.Equal(t, "terratest.example.com", mapping.Id)
	require.NotNil(t, mapping.SslSettings)
	assert.Equal(t, "AUTOMATIC", mapping.SslSettings.SslManagementType)
	require.Len(t, mapping.ResourceRecords, 1)
	assert.Equal(t, "CNAME", mapping.ResourceRecords[0].Type)
}

func TestGetAppEngineDomainMappingAttrsWithClientReportsAMissingMapping(t *testing.T) {
	t.Parallel()

	// A caller who names a domain that is not mapped should read a sentence about that domain, not a status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAppEngineDomainMappingAttrsWithClient(context.Background(), newFakeAppEngineService(t, handler), "gw-library-test-project", "missing.example.com")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing.example.com")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetAppEngineFirewallRuleAttrsWithClient(t *testing.T) {
	t.Parallel()

	// A rule is identified by its priority rather than by a name, so the priority is what the path carries.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/apps/gw-library-test-project/firewall/ingressRules/1000"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"priority":1000,"action":"ALLOW","sourceRange":"192.0.2.0/24","description":"terratest firewall rule"}`))
	})

	rule, err := gcp.GetAppEngineFirewallRuleAttrsWithClient(context.Background(), newFakeAppEngineService(t, handler), "gw-library-test-project", "1000")
	require.NoError(t, err)

	assert.Equal(t, int64(1000), rule.Priority)
	assert.Equal(t, "ALLOW", rule.Action)
	assert.Equal(t, "192.0.2.0/24", rule.SourceRange)
}

func TestGetAppEngineFirewallRuleAttrsWithClientReportsAMissingRule(t *testing.T) {
	t.Parallel()

	// A caller who names a priority that holds no rule should read a sentence about that priority, not a status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAppEngineFirewallRuleAttrsWithClient(context.Background(), newFakeAppEngineService(t, handler), "gw-library-test-project", "4000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "4000")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetAppEngineServiceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// One read answers for two modules, so the response carries both the traffic split and the network settings.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/apps/gw-library-test-project/services/default"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"apps/gw-library-test-project/services/default","id":"default","split":{"shardBy":"IP","allocations":{"v1":0.6,"v2":0.4}},"networkSettings":{"ingressTrafficAllowed":"INGRESS_TRAFFIC_ALLOWED_INTERNAL_ONLY"}}`))
	})

	service, err := gcp.GetAppEngineServiceAttrsWithClient(context.Background(), newFakeAppEngineService(t, handler), "gw-library-test-project", "default")
	require.NoError(t, err)

	require.NotNil(t, service.Split)
	assert.Equal(t, "IP", service.Split.ShardBy)
	assert.Len(t, service.Split.Allocations, 2)
	require.NotNil(t, service.NetworkSettings)
	assert.Equal(t, "INGRESS_TRAFFIC_ALLOWED_INTERNAL_ONLY", service.NetworkSettings.IngressTrafficAllowed)
}

func TestGetAppEngineServiceAttrsWithClientReportsAMissingService(t *testing.T) {
	t.Parallel()

	// A caller who names a service that is not there should read a sentence about that service, not a status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAppEngineServiceAttrsWithClient(context.Background(), newFakeAppEngineService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetAppEngineVersionAttrsWithClientAsksForTheFullView(t *testing.T) {
	t.Parallel()

	// The basic view leaves out the handlers, the entrypoint and the scaling settings, which is most of what a version
	// module sets, so the full view is asked for and the query is checked here.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.True(t, strings.HasSuffix(r.URL.Path, "/apps/gw-library-test-project/services/default/versions/gw-library-test"), "unexpected path %s", r.URL.Path)
		assert.Equal(t, "FULL", r.URL.Query().Get("view"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"apps/gw-library-test-project/services/default/versions/gw-library-test","id":"gw-library-test","runtime":"python312","servingStatus":"SERVING","instanceClass":"F1","automaticScaling":{"maxConcurrentRequests":10},"entrypoint":{"shell":"gunicorn -b :8080 main:app"}}`))
	})

	version, err := gcp.GetAppEngineVersionAttrsWithClient(context.Background(), newFakeAppEngineService(t, handler), "gw-library-test-project", "default", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "python312", version.Runtime)
	assert.Equal(t, "SERVING", version.ServingStatus)
	assert.Equal(t, "F1", version.InstanceClass)
	require.NotNil(t, version.AutomaticScaling)
	assert.Equal(t, int64(10), version.AutomaticScaling.MaxConcurrentRequests)
}

func TestGetAppEngineVersionAttrsWithClientReportsAMissingVersion(t *testing.T) {
	t.Parallel()

	// The sentence has to name the version and the service it was looked for on, because a version id is only unique
	// within its service.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAppEngineVersionAttrsWithClient(context.Background(), newFakeAppEngineService(t, handler), "gw-library-test-project", "default", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "default")
	assert.Contains(t, err.Error(), "does not exist")
}
