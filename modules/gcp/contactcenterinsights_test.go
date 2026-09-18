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
	"google.golang.org/api/contactcenterinsights/v1"
	"google.golang.org/api/option"
)

// newFakeContactCenterInsightsService points a real *contactcenterinsights.Service at a local test
// server, so the Google transport is exercised rather than a hand-written stand-in for a type we
// do not own.
func newFakeContactCenterInsightsService(t *testing.T, handler http.Handler) *contactcenterinsights.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := contactcenterinsights.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetInsightsViewAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-business-apps view module sets, because the
	// point of reading settings back is asserting a module configured the view it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/views/abc123"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/us-central1/views/abc123",
			"displayName":"terratest view",
			"value":"medium = \"CHAT\"",
			"createTime":"2026-09-14T00:00:00Z"
		}`))
	})

	view, err := gcp.GetInsightsViewAttrsWithClient(context.Background(), newFakeContactCenterInsightsService(t, handler), "gw-library-test-project", "us-central1", "abc123")
	require.NoError(t, err)

	assert.Equal(t, "terratest view", view.DisplayName)
	assert.Equal(t, `medium = "CHAT"`, view.Value)
	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1/views/abc123", view.Name)
}

func TestGetInsightsViewAttrsWithClientMissingView(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the view and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetInsightsViewAttrsWithClient(context.Background(), newFakeContactCenterInsightsService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
