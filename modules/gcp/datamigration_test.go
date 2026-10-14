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
	"google.golang.org/api/datamigration/v1"
	"google.golang.org/api/option"
)

// newFakeDatabaseMigrationService points a real Database Migration client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDatabaseMigrationService(t *testing.T, handler http.Handler) *datamigration.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := datamigration.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetConnectionProfileAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a profile the terraform-google-migration connection profile module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/connectionProfiles/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/connectionProfiles/gw-library-test","displayName":"terratest profile","provider":"DATABASE_PROVIDER_UNSPECIFIED","state":"READY","labels":{"purpose":"terratest"},"mysql":{"host":"10.50.0.1","port":3306,"username":"terratest"}}`))
	})

	profile, err := gcp.GetConnectionProfileAttrsWithClient(context.Background(), newFakeDatabaseMigrationService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest profile", profile.DisplayName)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, profile.Labels)
	require.NotNil(t, profile.Mysql)
	assert.Equal(t, "10.50.0.1", profile.Mysql.Host)
	assert.Equal(t, int64(3306), profile.Mysql.Port)
}

func TestGetConnectionProfileAttrsWithClientReportsAMissingProfile(t *testing.T) {
	t.Parallel()

	// A caller who asks for a profile that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetConnectionProfileAttrsWithClient(context.Background(), newFakeDatabaseMigrationService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDatabaseMigrationJobAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a database migration job a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/migrationJobs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"displayName":"terratest migration job","state":"NOT_STARTED","phase":"PHASE_UNSPECIFIED","type":"ONE_TIME","source":"projects/gw-library-test-project/locations/us-central1/connectionProfiles/gw-library-test-source","destination":"projects/gw-library-test-project/locations/us-central1/connectionProfiles/gw-library-test-destination"}`))
	})

	result, err := gcp.GetDatabaseMigrationJobAttrsWithClient(context.Background(), newFakeDatabaseMigrationService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest migration job", result.DisplayName)
	assert.Equal(t, "NOT_STARTED", result.State)
	assert.Equal(t, "ONE_TIME", result.Type)
	assert.Contains(t, result.Source, "gw-library-test-source")
	assert.Contains(t, result.Destination, "gw-library-test-destination")
}

func TestGetDatabaseMigrationJobAttrsWithClientReportsAMissingJob(t *testing.T) {
	t.Parallel()

	// A caller who names a database migration job that is not there should read a sentence about that database migration job, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDatabaseMigrationJobAttrsWithClient(context.Background(), newFakeDatabaseMigrationService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
