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
	"google.golang.org/api/sqladmin/v1"
)

// newFakeCloudSQLService points a real Cloud SQL client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudSQLService(t *testing.T, handler http.Handler) *sqladmin.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := sqladmin.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetCloudSQLInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a instance the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-library-test","databaseVersion":"POSTGRES_15","region":"us-central1","state":"RUNNABLE","settings":{"tier":"db-f1-micro","availabilityType":"ZONAL","backupConfiguration":{"enabled":false},"ipConfiguration":{"ipv4Enabled":false}}}`))
	})

	instance, err := gcp.GetCloudSQLInstanceAttrsWithClient(context.Background(), newFakeCloudSQLService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "POSTGRES_15", instance.DatabaseVersion)
	assert.Equal(t, "RUNNABLE", instance.State)
	require.NotNil(t, instance.Settings)
	assert.Equal(t, "db-f1-micro", instance.Settings.Tier)
	require.NotNil(t, instance.Settings.IpConfiguration)
	assert.False(t, instance.Settings.IpConfiguration.Ipv4Enabled)
}

func TestGetCloudSQLInstanceAttrsWithClientMissingInstance(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the instance and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetCloudSQLInstanceAttrsWithClient(context.Background(), newFakeCloudSQLService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetCloudSQLDatabaseAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a database the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/instances/gw-instance/databases/gw-database"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-database","instance":"gw-instance","charset":"UTF8","collation":"en_US.UTF8"}`))
	})

	database, err := gcp.GetCloudSQLDatabaseAttrsWithClient(context.Background(), newFakeCloudSQLService(t, handler), "gw-library-test-project", "gw-instance", "gw-database")
	require.NoError(t, err)

	assert.Equal(t, "gw-database", database.Name)
	assert.Equal(t, "UTF8", database.Charset)
	assert.Equal(t, "en_US.UTF8", database.Collation)
}

func TestGetCloudSQLDatabaseAttrsWithClientMissingDatabase(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the database and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetCloudSQLDatabaseAttrsWithClient(context.Background(), newFakeCloudSQLService(t, handler), "gw-library-test-project", "gw-instance", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-instance")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetCloudSQLUserAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a user the terraform-google-data-
	// storage module created, not a copy of any one fixture's values. A password is never in the
	// response, so nothing here can leak one.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Contains(t, r.URL.Path, "/instances/gw-instance/users", "unexpected path")
		assert.Equal(t, "%", r.URL.Query().Get("host"), "the host should be asked for, since a user is keyed by name and host")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"gw-user","instance":"gw-instance","project":"gw-library-test-project","type":"BUILT_IN"}`))
	})

	user, err := gcp.GetCloudSQLUserAttrsWithClient(context.Background(), newFakeCloudSQLService(t, handler), "gw-library-test-project", "gw-instance", "gw-user", "%")
	require.NoError(t, err)

	assert.Equal(t, "gw-user", user.Name)
	assert.Equal(t, "BUILT_IN", user.Type)
	assert.Empty(t, user.Password, "a password should never come back from the API")
}

func TestGetCloudSQLUserAttrsWithClientMissingUser(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the user and everything that identifies it, and each of those is a value no
	// other part of the message contains, or its check could not fail.
	_, err := gcp.GetCloudSQLUserAttrsWithClient(context.Background(), newFakeCloudSQLService(t, handler), "gw-library-test-project", "gw-instance", "gone", "%")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-instance")
	require.ErrorContains(t, err, "gw-library-test-project")
}
