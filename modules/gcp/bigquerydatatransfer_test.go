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
	"google.golang.org/api/bigquerydatatransfer/v1"
	"google.golang.org/api/option"
)

// newFakeBigQueryDataTransferService points a real BigQuery Data Transfer client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeBigQueryDataTransferService(t *testing.T, handler http.Handler) *bigquerydatatransfer.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := bigquerydatatransfer.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetBigQueryTransferConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a configuration the terraform-google-data-analytics data transfer config module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/transferConfigs/68d4a7d0-0000-2c0a-0000-000000000000"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/37950160017/locations/us-central1/transferConfigs/68d4a7d0-0000-2c0a-0000-000000000000","displayName":"terratest transfer","dataSourceId":"scheduled_query","schedule":"every 24 hours","disabled":true,"params":{"query":"SELECT 1","destination_table_name_template":"gw_library_test","write_disposition":"WRITE_TRUNCATE"}}`))
	})

	config, err := gcp.GetBigQueryTransferConfigAttrsWithClient(context.Background(), newFakeBigQueryDataTransferService(t, handler), "gw-library-test-project", "us-central1", "68d4a7d0-0000-2c0a-0000-000000000000")
	require.NoError(t, err)

	assert.Equal(t, "terratest transfer", config.DisplayName)
	assert.Equal(t, "scheduled_query", config.DataSourceId)
	assert.Equal(t, "every 24 hours", config.Schedule)
	assert.True(t, config.Disabled)
}

func TestGetBigQueryTransferConfigAttrsWithClientAcceptsWholeName(t *testing.T) {
	t.Parallel()

	// The Terraform resource's id is the whole resource name, built from the project number, so a
	// caller holding it must not have a second path prefixed onto it.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/projects/37950160017/locations/us-central1/transferConfigs/68d4a7d0-0000-2c0a-0000-000000000000", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/37950160017/locations/us-central1/transferConfigs/68d4a7d0-0000-2c0a-0000-000000000000","displayName":"terratest transfer","dataSourceId":"scheduled_query"}`))
	})

	config, err := gcp.GetBigQueryTransferConfigAttrsWithClient(context.Background(), newFakeBigQueryDataTransferService(t, handler),
		"gw-library-test-project", "us-central1", "projects/37950160017/locations/us-central1/transferConfigs/68d4a7d0-0000-2c0a-0000-000000000000")
	require.NoError(t, err)

	assert.Equal(t, "terratest transfer", config.DisplayName)
}
