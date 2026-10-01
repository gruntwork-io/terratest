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
	"google.golang.org/api/bigquerydatapolicy/v2"
	"google.golang.org/api/option"
)

// newFakeBigQueryDataPolicyService points a real BigQuery Data Policy client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeBigQueryDataPolicyService(t *testing.T, handler http.Handler) *bigquerydatapolicy.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := bigquerydatapolicy.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetBigQueryDataPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a data policy a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/dataPolicies/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/dataPolicies/gw-library-test","dataPolicyId":"gw-library-test","dataPolicyType":"DATA_MASKING_POLICY","dataMaskingPolicy":{"predefinedExpression":"SHA256"}}`))
	})

	result, err := gcp.GetBigQueryDataPolicyAttrsWithClient(context.Background(), newFakeBigQueryDataPolicyService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "DATA_MASKING_POLICY", result.DataPolicyType)
	require.NotNil(t, result.DataMaskingPolicy)
	assert.Equal(t, "SHA256", result.DataMaskingPolicy.PredefinedExpression)
}

func TestGetBigQueryDataPolicyAttrsWithClientReportsAMissingOne(t *testing.T) {
	t.Parallel()

	// A caller who asks for something that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetBigQueryDataPolicyAttrsWithClient(context.Background(), newFakeBigQueryDataPolicyService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
