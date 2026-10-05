package gcp_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/iap/v1"
	"google.golang.org/api/option"
)

// newFakeIAPService points a real Identity Aware Proxy client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeIAPService(t *testing.T, handler http.Handler) *iap.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := iap.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetIAPIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// Every IAP resource answers on one method and they differ only in the resource name, so the read
	// is given the name in full and the test checks it reaches the path unchanged.
	resource := "projects/gw-library-test-project/iap_web/compute/services/gw-library-test"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, resource+":getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"bindings":[{"role":"roles/iap.httpsResourceAccessor","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.GetIAPIamPolicyAttrsWithClient(context.Background(), newFakeIAPService(t, handler), resource)
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/iap.httpsResourceAccessor", policy.Bindings[0].Role)
	assert.Equal(t, int64(3), policy.Version)
}

func TestGetIAPIamPolicyAttrsWithClientAsksForConditionalBindings(t *testing.T) {
	t.Parallel()

	// A policy is returned at version 1 unless asked otherwise, and a version 1 answer has no room
	// for a condition, so a read that did not ask would silently drop one.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if assert.NoError(t, err) {
			var request iap.GetIamPolicyRequest
			if assert.NoError(t, json.Unmarshal(body, &request)) {
				require.NotNil(t, request.Options)
				assert.Equal(t, int64(3), request.Options.RequestedPolicyVersion)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"bindings":[{"role":"roles/iap.httpsResourceAccessor","members":["allAuthenticatedUsers"],"condition":{"title":"terratest","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetIAPIamPolicyAttrsWithClient(context.Background(), newFakeIAPService(t, handler),
		"projects/gw-library-test-project/iap_web")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	require.NotNil(t, policy.Bindings[0].Condition, "the condition should survive the read")
	assert.Equal(t, "terratest", policy.Bindings[0].Condition.Title)
}

func TestGetIAPIamPolicyAttrsWithClientReportsAMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a resource that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetIAPIamPolicyAttrsWithClient(context.Background(), newFakeIAPService(t, handler),
		"projects/gw-library-test-project/iap_web/compute/services/gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
