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
	"google.golang.org/api/option"
	"google.golang.org/api/servicemanagement/v1"
)

// newFakeServiceManagementService points a real Service Management client at a local test server, so the
// Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeServiceManagementService(t *testing.T, handler http.Handler) *servicemanagement.APIService {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := servicemanagement.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetEndpointsServiceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// An Endpoints service is named by its own DNS name rather than by a path under the project, so the
	// name goes into the path whole.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/services/gw-library-test.endpoints.gw-library-test-project.cloud.goog"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"serviceName":"gw-library-test.endpoints.gw-library-test-project.cloud.goog","producerProjectId":"gw-library-test-project"}`))
	})

	managed, err := gcp.GetEndpointsServiceAttrsWithClient(context.Background(), newFakeServiceManagementService(t, handler),
		"gw-library-test.endpoints.gw-library-test-project.cloud.goog")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test-project", managed.ProducerProjectId)
}

func TestGetEndpointsServiceIamPolicyAttrsWithClientAsksForConditionalBindings(t *testing.T) {
	t.Parallel()

	// A policy is returned at version 1 unless asked otherwise, and a version 1 answer has no room for
	// a condition, so a read that did not ask would silently drop one.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The policy lives under the service's resource path, so a read that passed the bare DNS name
		// would ask for a resource Google does not have.
		assert.True(t, strings.HasSuffix(r.URL.Path, "/services/gw-library-test.endpoints.gw-library-test-project.cloud.goog:getIamPolicy"),
			"unexpected path %s", r.URL.Path)

		body, err := io.ReadAll(r.Body)
		if assert.NoError(t, err) {
			var request servicemanagement.GetIamPolicyRequest
			if assert.NoError(t, json.Unmarshal(body, &request)) && assert.NotNil(t, request.Options) {
				assert.Equal(t, int64(3), request.Options.RequestedPolicyVersion)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"bindings":[{"role":"roles/servicemanagement.serviceController","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.GetEndpointsServiceIamPolicyAttrsWithClient(context.Background(), newFakeServiceManagementService(t, handler),
		"gw-library-test.endpoints.gw-library-test-project.cloud.goog")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/servicemanagement.serviceController", policy.Bindings[0].Role)
}

func TestGetEndpointsServiceConsumerIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The consumer policy is a separate resource from the service's own, and its name carries the
	// consumer project rather than the producer.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// The consumer is named by its bare project id, which is the only spelling the API serves.
		assert.Contains(t, r.URL.Path, "/consumers/gw-library-consumer")
		assert.NotContains(t, r.URL.Path, "project:")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"bindings":[{"role":"roles/servicemanagement.serviceConsumer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"]}]}`))
	})

	policy, err := gcp.GetEndpointsServiceConsumerIamPolicyAttrsWithClient(context.Background(), newFakeServiceManagementService(t, handler),
		"gw-library-test.endpoints.gw-library-test-project.cloud.goog", "gw-library-consumer")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, "roles/servicemanagement.serviceConsumer", policy.Bindings[0].Role)
}
