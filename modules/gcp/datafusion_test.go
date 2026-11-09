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
	"google.golang.org/api/datafusion/v1"
	"google.golang.org/api/option"
)

// newFakeDataFusionService points a real Data Fusion client at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDataFusionService(t *testing.T, handler http.Handler) *datafusion.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := datafusion.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDataFusionInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an instance the Data Fusion instance
	// module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/instances/gw-library-test","type":"BASIC","state":"ACTIVE","version":"6.10.1","apiEndpoint":"https://gw-library-test.datafusion.googleusercontent.com/api","tenantProjectId":"gw-library-tenant","privateInstance":false,"enableStackdriverLogging":true}`))
	})

	instance, err := gcp.GetDataFusionInstanceAttrsWithClient(context.Background(), newFakeDataFusionService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "BASIC", instance.Type)
	assert.Equal(t, "ACTIVE", instance.State)
	assert.Equal(t, "gw-library-tenant", instance.TenantProjectId)
	assert.True(t, instance.EnableStackdriverLogging)
}

func TestGetDataFusionInstanceAttrsWithClientReportsAMissingInstance(t *testing.T) {
	t.Parallel()

	// A caller who asks for an instance that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetDataFusionInstanceAttrsWithClient(context.Background(), newFakeDataFusionService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
