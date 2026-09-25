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
	"google.golang.org/api/bigqueryconnection/v1"
	"google.golang.org/api/option"
)

// newFakeBigQueryConnectionService points a real BigQuery Connection client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeBigQueryConnectionService(t *testing.T, handler http.Handler) *bigqueryconnection.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := bigqueryconnection.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetBigQueryConnectionAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a connection the terraform-google-data-analytics connection module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/connections/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/connections/gw-library-test","friendlyName":"terratest connection","description":"created by terratest","cloudResource":{"serviceAccountId":"bqcx-37950160017@gcp-sa-bigquery-condel.iam.gserviceaccount.com"}}`))
	})

	connection, err := gcp.GetBigQueryConnectionAttrsWithClient(context.Background(), newFakeBigQueryConnectionService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest connection", connection.FriendlyName)
	assert.Equal(t, "created by terratest", connection.Description)
	require.NotNil(t, connection.CloudResource, "a cloud resource connection carries the account BigQuery acts as")
}
