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
	"google.golang.org/api/firebasedataconnect/v1"
	"google.golang.org/api/option"
)

// newFakeFirebaseDataConnectService points a real Firebase Data Connect client at an httptest server, so a read can be exercised
// against a response we control without reaching Google.
func newFakeFirebaseDataConnectService(t *testing.T, handler http.Handler) *firebasedataconnect.APIService {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := firebasedataconnect.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetFirebaseDataConnectServiceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Firebase Data Connect service the terraform-google-firebase Data Connect service module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/services/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/services/gw-library-test","displayName":"terratest service","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetFirebaseDataConnectServiceAttrsWithClient(context.Background(), newFakeFirebaseDataConnectService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest service", attrs.DisplayName)
}

func TestGetFirebaseDataConnectServiceAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Firebase Data Connect service that is not there should read a sentence about that Firebase Data Connect service, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetFirebaseDataConnectServiceAttrsWithClient(context.Background(), newFakeFirebaseDataConnectService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
