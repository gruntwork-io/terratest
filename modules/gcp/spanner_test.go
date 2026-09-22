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
	"google.golang.org/api/spanner/v1"
)

// newFakeSpannerService points a real Spanner client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeSpannerService(t *testing.T, handler http.Handler) *spanner.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := spanner.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetSpannerInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a instance the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/instances/gw-library-test","displayName":"terratest instance","config":"projects/gw-library-test-project/instanceConfigs/regional-us-central1","processingUnits":100,"state":"READY","labels":{"purpose":"terratest"}}`))
	})

	instance, err := gcp.GetSpannerInstanceAttrsWithClient(context.Background(), newFakeSpannerService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest instance", instance.DisplayName)
	assert.Equal(t, int64(100), instance.ProcessingUnits)
	assert.Equal(t, "READY", instance.State)
	assert.True(t, strings.HasSuffix(instance.Config, "/instanceConfigs/regional-us-central1"))
}

func TestGetSpannerInstanceAttrsWithClientMissingInstance(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the instance and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetSpannerInstanceAttrsWithClient(context.Background(), newFakeSpannerService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetSpannerDatabaseAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a database the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/instances/gw-instance/databases/gw-database"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/instances/gw-instance/databases/gw-database","state":"READY","databaseDialect":"GOOGLE_STANDARD_SQL","versionRetentionPeriod":"1h"}`))
	})

	database, err := gcp.GetSpannerDatabaseAttrsWithClient(context.Background(), newFakeSpannerService(t, handler), "gw-library-test-project", "gw-instance", "gw-database")
	require.NoError(t, err)

	assert.Equal(t, "READY", database.State)
	assert.Equal(t, "GOOGLE_STANDARD_SQL", database.DatabaseDialect)
	assert.Equal(t, "1h", database.VersionRetentionPeriod)
}

func TestGetSpannerDatabaseAttrsWithClientMissingDatabase(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the database and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetSpannerDatabaseAttrsWithClient(context.Background(), newFakeSpannerService(t, handler), "gw-library-test-project", "gw-instance", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-instance")
	require.ErrorContains(t, err, "gw-library-test-project")
}
