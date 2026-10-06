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
	"google.golang.org/api/workstations/v1"
)

// newFakeWorkstationsService points a real Cloud Workstations client at an httptest server, so a read can be exercised
// against a response we control without reaching Google.
func newFakeWorkstationsService(t *testing.T, handler http.Handler) *workstations.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := workstations.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetWorkstationClusterAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a workstation cluster the terraform-google-devtools workstation cluster module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/workstationClusters/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/workstationClusters/gw-library-test","displayName":"terratest cluster","network":"projects/gw-library-test-project/global/networks/gw-library-test","subnetwork":"projects/gw-library-test-project/regions/us-central1/subnetworks/gw-library-test","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetWorkstationClusterAttrsWithClient(context.Background(), newFakeWorkstationsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest cluster", attrs.DisplayName)
	assert.Equal(t, "projects/gw-library-test-project/global/networks/gw-library-test", attrs.Network)
}

func TestGetWorkstationClusterAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a workstation cluster that is not there should read a sentence about that workstation cluster, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetWorkstationClusterAttrsWithClient(context.Background(), newFakeWorkstationsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetWorkstationConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a workstation config the terraform-google-devtools workstation config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/workstationClusters/gw-library-parent/workstationConfigs/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/workstationClusters/gw-library-parent/workstationConfigs/gw-library-test","displayName":"terratest config","idleTimeout":"1200s","runningTimeout":"7200s","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetWorkstationConfigAttrsWithClient(context.Background(), newFakeWorkstationsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest config", attrs.DisplayName)
	assert.Equal(t, "1200s", attrs.IdleTimeout)
	assert.Equal(t, "7200s", attrs.RunningTimeout)
}

func TestGetWorkstationConfigAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a workstation config that is not there should read a sentence about that workstation config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetWorkstationConfigAttrsWithClient(context.Background(), newFakeWorkstationsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetWorkstationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a workstation the terraform-google-devtools workstation module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/workstationClusters/gw-library-parent/workstationConfigs/gw-library-config/workstations/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/workstationClusters/gw-library-parent/workstationConfigs/gw-library-config/workstations/gw-library-test","displayName":"terratest workstation","state":"STATE_STOPPED","labels":{"purpose":"terratest"}}`))
	})

	attrs, err := gcp.GetWorkstationAttrsWithClient(context.Background(), newFakeWorkstationsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-config", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest workstation", attrs.DisplayName)
	assert.Equal(t, "STATE_STOPPED", attrs.State)
}

func TestGetWorkstationAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a workstation that is not there should read a sentence about that workstation, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetWorkstationAttrsWithClient(context.Background(), newFakeWorkstationsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-config", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetWorkstationConfigIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a workstation config the terraform-google-devtools workstation config module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/workstationClusters/gw-library-parent/workstationConfigs/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/workstations.workstationCreator","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetWorkstationConfigIamPolicyAttrsWithClient(context.Background(), newFakeWorkstationsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/workstations.workstationCreator", policy.Bindings[0].Role)
}

func TestGetWorkstationConfigIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a workstation config that is not there should read a sentence about that workstation config, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetWorkstationConfigIamPolicyAttrsWithClient(context.Background(), newFakeWorkstationsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetWorkstationIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a workstation the terraform-google-devtools workstation module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/workstationClusters/gw-library-parent/workstationConfigs/gw-library-config/workstations/gw-library-test:getIamPolicy"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"version":3,"etag":"BwXhqw==","bindings":[{"role":"roles/workstations.user","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetWorkstationIamPolicyAttrsWithClient(context.Background(), newFakeWorkstationsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-config", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/workstations.user", policy.Bindings[0].Role)
}

func TestGetWorkstationIamPolicyAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a workstation that is not there should read a sentence about that workstation, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetWorkstationIamPolicyAttrsWithClient(context.Background(), newFakeWorkstationsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-config", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
