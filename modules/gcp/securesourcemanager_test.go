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
	"google.golang.org/api/securesourcemanager/v1"
)

// newFakeSecureSourceManagerService points a real Secure Source Manager client at an httptest server, so a read can be exercised
// against a response we control without reaching Google.
func newFakeSecureSourceManagerService(t *testing.T, handler http.Handler) *securesourcemanager.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := securesourcemanager.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetSecureSourceManagerInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Secure Source Manager instance the terraform-google-devtools Secure Source Manager instance module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/instances/gw-library-test","state":"ACTIVE","kmsKey":"projects/gw-library-test-project/locations/us-central1/keyRings/gw-library-test/cryptoKeys/gw-library-test","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetSecureSourceManagerInstanceAttrsWithClient(context.Background(), newFakeSecureSourceManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "ACTIVE", attrs.State)
	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1/keyRings/gw-library-test/cryptoKeys/gw-library-test", attrs.KmsKey)
}

func TestGetSecureSourceManagerInstanceAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Secure Source Manager instance that is not there should read a sentence about that Secure Source Manager instance, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetSecureSourceManagerInstanceAttrsWithClient(context.Background(), newFakeSecureSourceManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetSecureSourceManagerRepositoryAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Secure Source Manager repository the terraform-google-devtools Secure Source Manager repository module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/repositories/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/repositories/gw-library-test","description":"created by terratest","instance":"projects/gw-library-test-project/locations/us-central1/instances/gw-library-parent"}`))
	})

	attrs, err := gcp.GetSecureSourceManagerRepositoryAttrsWithClient(context.Background(), newFakeSecureSourceManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1/instances/gw-library-parent", attrs.Instance)
}

func TestGetSecureSourceManagerRepositoryAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Secure Source Manager repository that is not there should read a sentence about that Secure Source Manager repository, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetSecureSourceManagerRepositoryAttrsWithClient(context.Background(), newFakeSecureSourceManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetSecureSourceManagerBranchRuleAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Secure Source Manager branch rule the terraform-google-devtools Secure Source Manager branch rule module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/repositories/gw-library-parent/branchRules/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/repositories/gw-library-parent/branchRules/gw-library-test","includePattern":"main","requirePullRequest":true,"minimumApprovalsCount":2}`))
	})

	attrs, err := gcp.GetSecureSourceManagerBranchRuleAttrsWithClient(context.Background(), newFakeSecureSourceManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "main", attrs.IncludePattern)
	assert.Equal(t, int64(2), attrs.MinimumApprovalsCount)
}

func TestGetSecureSourceManagerBranchRuleAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Secure Source Manager branch rule that is not there should read a sentence about that Secure Source Manager branch rule, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetSecureSourceManagerBranchRuleAttrsWithClient(context.Background(), newFakeSecureSourceManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetSecureSourceManagerHookAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Secure Source Manager hook the terraform-google-devtools Secure Source Manager hook module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/repositories/gw-library-parent/hooks/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/repositories/gw-library-parent/hooks/gw-library-test","targetUri":"https://terratest.example.com/hook","disabled":false,"events":["PUSH"]}`))
	})

	attrs, err := gcp.GetSecureSourceManagerHookAttrsWithClient(context.Background(), newFakeSecureSourceManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "https://terratest.example.com/hook", attrs.TargetUri)
	assert.Equal(t, "PUSH", attrs.Events[0])
}

func TestGetSecureSourceManagerHookAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Secure Source Manager hook that is not there should read a sentence about that Secure Source Manager hook, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetSecureSourceManagerHookAttrsWithClient(context.Background(), newFakeSecureSourceManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetSecureSourceManagerInstanceIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Secure Source Manager instance the terraform-google-devtools Secure Source Manager instance module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/instances/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/securesourcemanager.instanceOwner","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetSecureSourceManagerInstanceIamPolicyAttrsWithClient(context.Background(), newFakeSecureSourceManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/securesourcemanager.instanceOwner", policy.Bindings[0].Role)
}

func TestGetSecureSourceManagerInstanceIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Secure Source Manager instance that is not there should read a sentence about that Secure Source Manager instance, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetSecureSourceManagerInstanceIamPolicyAttrsWithClient(context.Background(), newFakeSecureSourceManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetSecureSourceManagerRepositoryIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Secure Source Manager repository the terraform-google-devtools Secure Source Manager repository module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/repositories/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/securesourcemanager.repoReader","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetSecureSourceManagerRepositoryIamPolicyAttrsWithClient(context.Background(), newFakeSecureSourceManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/securesourcemanager.repoReader", policy.Bindings[0].Role)
}

func TestGetSecureSourceManagerRepositoryIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Secure Source Manager repository that is not there should read a sentence about that Secure Source Manager repository, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetSecureSourceManagerRepositoryIamPolicyAttrsWithClient(context.Background(), newFakeSecureSourceManagerService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
