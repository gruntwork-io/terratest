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
	"google.golang.org/api/apikeys/v2"
	"google.golang.org/api/option"
)

// newFakeAPIKeysService points a real *apikeys.Service at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeAPIKeysService(t *testing.T, handler http.Handler) *apikeys.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := apikeys.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetAPIKeyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-api-management key module sets, because the
	// point of reading settings back is asserting a module configured the key it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/locations/global/keys/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/global/keys/gw-library-test",
			"uid":"gw-library-test",
			"displayName":"terratest key",
			"restrictions":{"apiTargets":[{"service":"translate.googleapis.com"}]}
		}`))
	})

	key, err := gcp.GetAPIKeyAttrsWithClient(context.Background(), newFakeAPIKeysService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest key", key.DisplayName)
	assert.Equal(t, "gw-library-test", key.Uid)
	require.NotNil(t, key.Restrictions)
	require.Len(t, key.Restrictions.ApiTargets, 1)
	assert.Equal(t, "translate.googleapis.com", key.Restrictions.ApiTargets[0].Service)
}

func TestGetAPIKeyAttrsWithClientMissingKey(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the key and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetAPIKeyAttrsWithClient(context.Background(), newFakeAPIKeysService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
