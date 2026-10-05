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
	"google.golang.org/api/apihub/v1"
	"google.golang.org/api/option"
)

// newFakeAPIHubService points a real API hub client at an httptest server, so a read can be exercised
// against a response we control without reaching Google.
func newFakeAPIHubService(t *testing.T, handler http.Handler) *apihub.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := apihub.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetAPIHubInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a API hub instance the terraform-google-api-management API hub instance module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/apiHubInstances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/apiHubInstances/gw-library-test","description":"created by terratest","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetAPIHubInstanceAttrsWithClient(context.Background(), newFakeAPIHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetAPIHubInstanceAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a API hub instance that is not there should read a sentence about that API hub instance, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAPIHubInstanceAttrsWithClient(context.Background(), newFakeAPIHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetAPIHubCurationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a API hub curation the terraform-google-api-management API hub curation module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/curations/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/curations/gw-library-test","displayName":"terratest curation","description":"created by terratest","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetAPIHubCurationAttrsWithClient(context.Background(), newFakeAPIHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest curation", attrs.DisplayName)
	assert.Equal(t, "created by terratest", attrs.Description)
}

func TestGetAPIHubCurationAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a API hub curation that is not there should read a sentence about that API hub curation, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAPIHubCurationAttrsWithClient(context.Background(), newFakeAPIHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetAPIHubHostProjectRegistrationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a API hub host project registration the terraform-google-api-management API hub host project registration module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/hostProjectRegistrations/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/hostProjectRegistrations/gw-library-test","gcpProject":"projects/gw-library-other-project"}`))
	})

	attrs, err := gcp.GetAPIHubHostProjectRegistrationAttrsWithClient(context.Background(), newFakeAPIHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-other-project", attrs.GcpProject)
}

func TestGetAPIHubHostProjectRegistrationAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a API hub host project registration that is not there should read a sentence about that API hub host project registration, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAPIHubHostProjectRegistrationAttrsWithClient(context.Background(), newFakeAPIHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetAPIHubPluginAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a API hub plugin the terraform-google-api-management API hub plugin module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/plugins/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/plugins/gw-library-test","displayName":"terratest plugin","description":"created by terratest","state":"ENABLED"}`))
	})

	attrs, err := gcp.GetAPIHubPluginAttrsWithClient(context.Background(), newFakeAPIHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest plugin", attrs.DisplayName)
	assert.Equal(t, "ENABLED", attrs.State)
}

func TestGetAPIHubPluginAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a API hub plugin that is not there should read a sentence about that API hub plugin, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAPIHubPluginAttrsWithClient(context.Background(), newFakeAPIHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetAPIHubPluginInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a API hub plugin instance the terraform-google-api-management API hub plugin instance module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/plugins/gw-library-parent/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/plugins/gw-library-parent/instances/gw-library-test","displayName":"terratest plugin instance","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetAPIHubPluginInstanceAttrsWithClient(context.Background(), newFakeAPIHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest plugin instance", attrs.DisplayName)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetAPIHubPluginInstanceAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a API hub plugin instance that is not there should read a sentence about that API hub plugin instance, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAPIHubPluginInstanceAttrsWithClient(context.Background(), newFakeAPIHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetAPIHubRuntimeProjectAttachmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a API hub runtime project attachment the terraform-google-api-management API hub runtime project attachment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/runtimeProjectAttachments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/runtimeProjectAttachments/gw-library-test","runtimeProject":"projects/gw-library-other-project"}`))
	})

	attrs, err := gcp.GetAPIHubRuntimeProjectAttachmentAttrsWithClient(context.Background(), newFakeAPIHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-other-project", attrs.RuntimeProject)
}

func TestGetAPIHubRuntimeProjectAttachmentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a API hub runtime project attachment that is not there should read a sentence about that API hub runtime project attachment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAPIHubRuntimeProjectAttachmentAttrsWithClient(context.Background(), newFakeAPIHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
