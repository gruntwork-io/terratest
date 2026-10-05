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
	"google.golang.org/api/discoveryengine/v1"
	"google.golang.org/api/option"
)

// newFakeDiscoveryEngineService points a real Discovery Engine client at an httptest server, so a read
// can be exercised against a response we control without reaching Google.
func newFakeDiscoveryEngineService(t *testing.T, handler http.Handler) *discoveryengine.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := discoveryengine.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetDiscoveryEngineACLConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine ACL config the terraform-google-ml Discovery Engine ACL config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/aclConfig"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/aclConfig","idpConfig":{"idpType":"GSUITE"}}`))
	})

	attrs, err := gcp.GetDiscoveryEngineACLConfigAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global")
	require.NoError(t, err)

	assert.Equal(t, "GSUITE", attrs.IdpConfig.IdpType)
}

func TestGetDiscoveryEngineACLConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine ACL config that is not there should read a sentence about that Discovery Engine ACL config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineACLConfigAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineCmekConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine CMEK config the terraform-google-ml Discovery Engine CMEK config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/cmekConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/cmekConfigs/gw-library-test","kmsKey":"projects/gw-library-test-project/locations/us-central1/keyRings/gw-library-test/cryptoKeys/gw-library-test","state":"ACTIVE","isDefault":true}`))
	})

	attrs, err := gcp.GetDiscoveryEngineCmekConfigAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1/keyRings/gw-library-test/cryptoKeys/gw-library-test", attrs.KmsKey)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetDiscoveryEngineCmekConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine CMEK config that is not there should read a sentence about that Discovery Engine CMEK config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineCmekConfigAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineDataStoreAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine data store the terraform-google-ml Discovery Engine data store module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/collections/default_collection/dataStores/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/collections/default_collection/dataStores/gw-library-test","displayName":"terratest data store","industryVertical":"GENERIC","contentConfig":"NO_CONTENT","solutionTypes":["SOLUTION_TYPE_SEARCH"]}`))
	})

	attrs, err := gcp.GetDiscoveryEngineDataStoreAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest data store", attrs.DisplayName)
	assert.Equal(t, "GENERIC", attrs.IndustryVertical)
	assert.Equal(t, "NO_CONTENT", attrs.ContentConfig)
}

func TestGetDiscoveryEngineDataStoreAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine data store that is not there should read a sentence about that Discovery Engine data store, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineDataStoreAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineEngineAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine engine the terraform-google-ml Discovery Engine engine module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/collections/default_collection/engines/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/collections/default_collection/engines/gw-library-test","displayName":"terratest engine","solutionType":"SOLUTION_TYPE_SEARCH","industryVertical":"GENERIC","dataStoreIds":["gw-library-parent"]}`))
	})

	attrs, err := gcp.GetDiscoveryEngineEngineAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest engine", attrs.DisplayName)
	assert.Equal(t, "SOLUTION_TYPE_SEARCH", attrs.SolutionType)
	assert.Equal(t, "gw-library-parent", attrs.DataStoreIds[0])
}

func TestGetDiscoveryEngineEngineAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine engine that is not there should read a sentence about that Discovery Engine engine, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineEngineAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineEngineIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine engine the terraform-google-ml Discovery Engine engine module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/collections/default_collection/engines/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/discoveryengine.viewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetDiscoveryEngineEngineIamPolicyAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/discoveryengine.viewer", policy.Bindings[0].Role)
}

func TestGetDiscoveryEngineEngineIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine engine that is not there should read a sentence about that Discovery Engine engine, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineEngineIamPolicyAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineAssistantAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine assistant the terraform-google-ml Discovery Engine assistant module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/collections/default_collection/engines/gw-library-parent/assistants/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/collections/default_collection/engines/gw-library-parent/assistants/gw-library-test","displayName":"terratest assistant"}`))
	})

	attrs, err := gcp.GetDiscoveryEngineAssistantAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest assistant", attrs.DisplayName)
}

func TestGetDiscoveryEngineAssistantAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine assistant that is not there should read a sentence about that Discovery Engine assistant, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineAssistantAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineDataConnectorAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine data connector the terraform-google-ml Discovery Engine data connector module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/collections/default_collection/dataConnector"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/collections/default_collection/dataConnector","dataSource":"terratest","refreshInterval":"86400s","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetDiscoveryEngineDataConnectorAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection")
	require.NoError(t, err)

	assert.Equal(t, "terratest", attrs.DataSource)
	assert.Equal(t, "86400s", attrs.RefreshInterval)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetDiscoveryEngineDataConnectorAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine data connector that is not there should read a sentence about that Discovery Engine data connector, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineDataConnectorAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineSchemaAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine schema the terraform-google-ml Discovery Engine schema module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/collections/default_collection/dataStores/gw-library-parent/schemas/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/collections/default_collection/dataStores/gw-library-parent/schemas/gw-library-test","jsonSchema":"{\"type\":\"object\"}"}`))
	})

	attrs, err := gcp.GetDiscoveryEngineSchemaAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.JSONEq(t, `{"type":"object"}`, attrs.JsonSchema)
}

func TestGetDiscoveryEngineSchemaAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine schema that is not there should read a sentence about that Discovery Engine schema, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineSchemaAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineServingConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine serving config the terraform-google-ml Discovery Engine serving config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/collections/default_collection/dataStores/gw-library-parent/servingConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/collections/default_collection/dataStores/gw-library-parent/servingConfigs/gw-library-test","displayName":"terratest serving config","solutionType":"SOLUTION_TYPE_SEARCH"}`))
	})

	attrs, err := gcp.GetDiscoveryEngineServingConfigAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest serving config", attrs.DisplayName)
	assert.Equal(t, "SOLUTION_TYPE_SEARCH", attrs.SolutionType)
}

func TestGetDiscoveryEngineServingConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine serving config that is not there should read a sentence about that Discovery Engine serving config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineServingConfigAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineControlAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine control the terraform-google-ml Discovery Engine control module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/collections/default_collection/dataStores/gw-library-parent/controls/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/collections/default_collection/dataStores/gw-library-parent/controls/gw-library-test","displayName":"terratest control","solutionType":"SOLUTION_TYPE_SEARCH"}`))
	})

	attrs, err := gcp.GetDiscoveryEngineControlAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest control", attrs.DisplayName)
	assert.Equal(t, "SOLUTION_TYPE_SEARCH", attrs.SolutionType)
}

func TestGetDiscoveryEngineControlAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine control that is not there should read a sentence about that Discovery Engine control, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineControlAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineTargetSiteAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine target site the terraform-google-ml Discovery Engine target site module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/collections/default_collection/dataStores/gw-library-parent/siteSearchEngine/targetSites/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/collections/default_collection/dataStores/gw-library-parent/siteSearchEngine/targetSites/gw-library-test","providedUriPattern":"terratest.example.com/*","type":"INCLUDE","exactMatch":false}`))
	})

	attrs, err := gcp.GetDiscoveryEngineTargetSiteAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest.example.com/*", attrs.ProvidedUriPattern)
	assert.Equal(t, "INCLUDE", attrs.Type)
}

func TestGetDiscoveryEngineTargetSiteAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine target site that is not there should read a sentence about that Discovery Engine target site, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineTargetSiteAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineWidgetConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine widget config the terraform-google-ml Discovery Engine widget config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/collections/default_collection/dataStores/gw-library-parent/widgetConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/collections/default_collection/dataStores/gw-library-parent/widgetConfigs/gw-library-test","displayName":"terratest widget","industryVertical":"GENERIC"}`))
	})

	attrs, err := gcp.GetDiscoveryEngineWidgetConfigAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest widget", attrs.DisplayName)
	assert.Equal(t, "GENERIC", attrs.IndustryVertical)
}

func TestGetDiscoveryEngineWidgetConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine widget config that is not there should read a sentence about that Discovery Engine widget config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineWidgetConfigAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "default_collection", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineLicenseConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine license config the terraform-google-ml Discovery Engine license config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/licenseConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/licenseConfigs/gw-library-test","licenseCount":"5","subscriptionTier":"SUBSCRIPTION_TIER_SEARCH","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetDiscoveryEngineLicenseConfigAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(5), attrs.LicenseCount)
	assert.Equal(t, "SUBSCRIPTION_TIER_SEARCH", attrs.SubscriptionTier)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetDiscoveryEngineLicenseConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine license config that is not there should read a sentence about that Discovery Engine license config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineLicenseConfigAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDiscoveryEngineUserStoreAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Discovery Engine user store the terraform-google-ml Discovery Engine user store module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/global/userStores/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/global/userStores/gw-library-test","displayName":"terratest user store","defaultLicenseConfig":"projects/gw-library-test-project/locations/global/licenseConfigs/gw-library-parent"}`))
	})

	attrs, err := gcp.GetDiscoveryEngineUserStoreAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest user store", attrs.DisplayName)
}

func TestGetDiscoveryEngineUserStoreAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Discovery Engine user store that is not there should read a sentence about that Discovery Engine user store, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDiscoveryEngineUserStoreAttrsWithClient(context.Background(), newFakeDiscoveryEngineService(t, handler), "gw-library-test-project", "global", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
