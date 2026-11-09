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
	"google.golang.org/api/option"
	"google.golang.org/api/securitycenter/v1"
)

// newFakeSecurityCenterService points a real Security Command Center client at a local test server,
// so the Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeSecurityCenterService(t *testing.T, handler http.Handler) *securitycenter.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := securitycenter.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetSecurityCenterMuteConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a config the mute config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/muteConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/muteConfigs/gw-library-test","description":"created by terratest","filter":"category=\"OPEN_FIREWALL\"","type":"STATIC","mostRecentEditor":"terratest@gw-library-test-project.iam.gserviceaccount.com"}`))
	})

	config, err := gcp.GetSecurityCenterMuteConfigAttrsWithClient(context.Background(), newFakeSecurityCenterService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, `category="OPEN_FIREWALL"`, config.Filter)
	assert.Equal(t, "STATIC", config.Type)
}

func TestGetSecurityCenterMuteConfigAttrsWithClientReportsAMissingConfig(t *testing.T) {
	t.Parallel()

	// A caller who asks for a config that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetSecurityCenterMuteConfigAttrsWithClient(context.Background(), newFakeSecurityCenterService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetSecurityCenterNotificationConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a config the notification config module
	// created, not a copy of any one fixture's values. The service account is filled in by the API.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/notificationConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/notificationConfigs/gw-library-test","description":"created by terratest","pubsubTopic":"projects/gw-library-test-project/topics/gw-library-test","serviceAccount":"service-424242424242@gcp-sa-scc-notification.iam.gserviceaccount.com","streamingConfig":{"filter":"state=\"ACTIVE\""}}`))
	})

	config, err := gcp.GetSecurityCenterNotificationConfigAttrsWithClient(context.Background(), newFakeSecurityCenterService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/topics/gw-library-test", config.PubsubTopic)
	assert.Contains(t, config.ServiceAccount, "gcp-sa-scc-notification")
	require.NotNil(t, config.StreamingConfig)
	assert.Equal(t, `state="ACTIVE"`, config.StreamingConfig.Filter)
}

func TestGetSecurityCenterNotificationConfigAttrsWithClientReportsAMissingConfig(t *testing.T) {
	t.Parallel()

	// A caller who asks for a config that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetSecurityCenterNotificationConfigAttrsWithClient(context.Background(), newFakeSecurityCenterService(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetSecurityCenterSecurityHealthAnalyticsCustomModuleAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a module the project custom module
	// module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/securityHealthAnalyticsSettings/customModules/424242424242"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/securityHealthAnalyticsSettings/customModules/424242424242","displayName":"gw_library_test","enablementState":"ENABLED","cloudProvider":"GOOGLE_CLOUD_PLATFORM","customConfig":{"severity":"LOW","predicate":{"expression":"resource.name.contains(\"terratest\")"},"recommendation":"created by terratest"}}`))
	})

	module, err := gcp.GetSecurityCenterSecurityHealthAnalyticsCustomModuleAttrsWithClient(context.Background(), newFakeSecurityCenterService(t, handler), "gw-library-test-project", "424242424242")
	require.NoError(t, err)

	assert.Equal(t, "ENABLED", module.EnablementState)
	require.NotNil(t, module.CustomConfig)
	assert.Equal(t, "LOW", module.CustomConfig.Severity)
}

func TestGetSecurityCenterSecurityHealthAnalyticsCustomModuleAttrsWithClientReportsAMissingModule(t *testing.T) {
	t.Parallel()

	// A caller who asks for a module that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetSecurityCenterSecurityHealthAnalyticsCustomModuleAttrsWithClient(context.Background(), newFakeSecurityCenterService(t, handler), "gw-library-test-project", "515151515151")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetSecurityCenterEventThreatDetectionCustomModuleAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a module the Event Threat Detection
	// custom module module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/eventThreatDetectionSettings/customModules/424242424242"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/eventThreatDetectionSettings/customModules/424242424242","displayName":"gw_library_test","enablementState":"ENABLED","type":"CONFIGURABLE_BAD_IP","config":{"metadata":{"severity":"LOW","description":"created by terratest"},"ips":["192.0.2.1"]}}`))
	})

	module, err := gcp.GetSecurityCenterEventThreatDetectionCustomModuleAttrsWithClient(context.Background(), newFakeSecurityCenterService(t, handler), "gw-library-test-project", "424242424242")
	require.NoError(t, err)

	assert.Equal(t, "CONFIGURABLE_BAD_IP", module.Type)
	assert.Equal(t, "ENABLED", module.EnablementState)
	assert.Contains(t, string(module.Config), "192.0.2.1")
}

func TestGetSecurityCenterEventThreatDetectionCustomModuleAttrsWithClientReportsAMissingModule(t *testing.T) {
	t.Parallel()

	// A caller who asks for a module that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetSecurityCenterEventThreatDetectionCustomModuleAttrsWithClient(context.Background(), newFakeSecurityCenterService(t, handler), "gw-library-test-project", "515151515151")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
