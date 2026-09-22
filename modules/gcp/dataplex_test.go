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
	"google.golang.org/api/dataplex/v1"
	"google.golang.org/api/option"
)

// newFakeDataplexService points a real *dataplex.Service at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDataplexService(t *testing.T, handler http.Handler) *dataplex.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := dataplex.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDataplexLakeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-data-analytics lake module sets, because the
	// point of reading settings back is asserting a module configured the lake it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/lakes/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/us-central1/lakes/gw-library-test",
			"displayName":"terratest lake",
			"description":"created by terratest",
			"state":"ACTIVE",
			"labels":{"managed-by":"terratest"}
		}`))
	})

	lake, err := gcp.GetDataplexLakeAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest lake", lake.DisplayName)
	assert.Equal(t, "created by terratest", lake.Description)
	assert.Equal(t, "ACTIVE", lake.State)
	assert.Equal(t, "terratest", lake.Labels["managed-by"])
}

func TestGetDataplexLakeAttrsWithClientMissingLake(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the lake and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetDataplexLakeAttrsWithClient(context.Background(), newFakeDataplexService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
