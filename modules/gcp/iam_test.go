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
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

// newFakeIAMService points a real *iam.Service at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeIAMService(t *testing.T, handler http.Handler) *iam.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := iam.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetServiceAccountAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-management service account module sets, because
	// the point of reading settings back is asserting a module configured the account it was asked
	// for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/serviceAccounts/gw-library-test@gw-library-test-project.iam.gserviceaccount.com"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/serviceAccounts/gw-library-test@gw-library-test-project.iam.gserviceaccount.com",
			"projectId":"gw-library-test-project",
			"uniqueId":"104000000000000000000",
			"email":"gw-library-test@gw-library-test-project.iam.gserviceaccount.com",
			"displayName":"terratest account",
			"description":"created by terratest",
			"disabled":true
		}`))
	})

	account, err := gcp.GetServiceAccountAttrsWithClient(context.Background(), newFakeIAMService(t, handler),
		"gw-library-test-project", "gw-library-test@gw-library-test-project.iam.gserviceaccount.com")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test@gw-library-test-project.iam.gserviceaccount.com", account.Email)
	assert.Equal(t, "terratest account", account.DisplayName)
	assert.Equal(t, "created by terratest", account.Description)
	assert.Equal(t, "gw-library-test-project", account.ProjectId)
	assert.True(t, account.Disabled)
}

func TestGetWorkloadIdentityPoolAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-identity pool module sets, because the point of
	// reading settings back is asserting a module configured the pool it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/locations/global/workloadIdentityPools/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/global/workloadIdentityPools/gw-library-test",
			"displayName":"terratest pool",
			"description":"created by terratest",
			"state":"ACTIVE",
			"disabled":true,
			"mode":"FEDERATION_ONLY"
		}`))
	})

	pool, err := gcp.GetWorkloadIdentityPoolAttrsWithClient(context.Background(), newFakeIAMService(t, handler),
		"gw-library-test-project", "global", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest pool", pool.DisplayName)
	assert.Equal(t, "created by terratest", pool.Description)
	assert.Equal(t, "ACTIVE", pool.State)
	assert.Equal(t, "FEDERATION_ONLY", pool.Mode)
	assert.True(t, pool.Disabled)
}

func TestGetWorkloadIdentityPoolAttrsWithClientMissingPool(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the pool and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetWorkloadIdentityPoolAttrsWithClient(context.Background(), newFakeIAMService(t, handler),
		"gw-library-test-project", "global", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetServiceAccountAttrsWithClientMissingAccount(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the account and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetServiceAccountAttrsWithClient(context.Background(), newFakeIAMService(t, handler),
		"gw-library-test-project", "gone@gw-library-test-project.iam.gserviceaccount.com")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone@")
	require.ErrorContains(t, err, "gw-library-test-project")
}
