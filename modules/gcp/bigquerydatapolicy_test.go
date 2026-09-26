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
	"google.golang.org/api/bigquerydatapolicy/v1"
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

	// The response is shaped like the one Google returns for a policy the terraform-google-data-analytics data policy module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/dataPolicies/gw_library_test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/37950160017/locations/us-central1/dataPolicies/gw_library_test","dataPolicyId":"gw_library_test","dataPolicyType":"DATA_MASKING_POLICY","policyTag":"projects/gw-library-test-project/locations/us-central1/taxonomies/1234567890/policyTags/987654321","dataMaskingPolicy":{"predefinedExpression":"SHA256"}}`))
	})

	policy, err := gcp.GetBigQueryDataPolicyAttrsWithClient(context.Background(), newFakeBigQueryDataPolicyService(t, handler), "gw-library-test-project", "us-central1", "gw_library_test")
	require.NoError(t, err)

	assert.Equal(t, "DATA_MASKING_POLICY", policy.DataPolicyType)
	require.NotNil(t, policy.DataMaskingPolicy)
	assert.Equal(t, "SHA256", policy.DataMaskingPolicy.PredefinedExpression)
}
