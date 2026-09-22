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
	"google.golang.org/api/firestore/v1"
	"google.golang.org/api/option"
)

// newFakeFirestoreService points a real Firestore client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeFirestoreService(t *testing.T, handler http.Handler) *firestore.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := firestore.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetFirestoreDatabaseAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a database the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/databases/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/databases/gw-library-test","locationId":"nam5","type":"FIRESTORE_NATIVE","concurrencyMode":"OPTIMISTIC","appEngineIntegrationMode":"DISABLED","deleteProtectionState":"DELETE_PROTECTION_DISABLED"}`))
	})

	database, err := gcp.GetFirestoreDatabaseAttrsWithClient(context.Background(), newFakeFirestoreService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "nam5", database.LocationId)
	assert.Equal(t, "FIRESTORE_NATIVE", database.Type)
	assert.Equal(t, "OPTIMISTIC", database.ConcurrencyMode)
	assert.Equal(t, "DELETE_PROTECTION_DISABLED", database.DeleteProtectionState)
}

func TestGetFirestoreDatabaseAttrsWithClientMissingDatabase(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the database and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetFirestoreDatabaseAttrsWithClient(context.Background(), newFakeFirestoreService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
