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
	"google.golang.org/api/secretmanager/v1"
)

// newFakeSecretManagerService points a real *secretmanager.Service at a local test server, so the
// Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeSecretManagerService(t *testing.T, handler http.Handler) *secretmanager.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := secretmanager.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetSecretAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-security secret module sets, because the point
	// of reading settings back is asserting a module configured the secret it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/secrets/gw-library-test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/secrets/gw-library-test",
			"replication":{"userManaged":{"replicas":[{"location":"us-central1"}]}},
			"labels":{"purpose":"terratest"},
			"annotations":{"owner":"terratest"},
			"ttl":"3600s"
		}`))
	})

	secret, err := gcp.GetSecretAttrsWithClient(context.Background(), newFakeSecretManagerService(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/secrets/gw-library-test", secret.Name)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, secret.Labels)
	assert.Equal(t, map[string]string{"owner": "terratest"}, secret.Annotations)
	require.Len(t, secret.Replication.UserManaged.Replicas, 1)
	assert.Equal(t, "us-central1", secret.Replication.UserManaged.Replicas[0].Location)
}

func TestGetSecretAttrsWithClientMissingSecret(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the secret and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetSecretAttrsWithClient(context.Background(), newFakeSecretManagerService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
