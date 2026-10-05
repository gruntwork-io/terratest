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
	"google.golang.org/api/apphub/v1"
	"google.golang.org/api/option"
)

// newFakeAppHubService points a real App Hub client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeAppHubService(t *testing.T, handler http.Handler) *apphub.APIService {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := apphub.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetAppHubApplicationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an application the terraform-google-management App Hub application module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/applications/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/applications/gw-library-test","displayName":"terratest application","description":"created by terratest","scope":{"type":"REGIONAL"},"attributes":{"criticality":{"type":"LOW"},"environment":{"type":"TEST"}}}`))
	})

	application, err := gcp.GetAppHubApplicationAttrsWithClient(context.Background(), newFakeAppHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest application", application.DisplayName)
	require.NotNil(t, application.Scope)
	assert.Equal(t, "REGIONAL", application.Scope.Type)
	require.NotNil(t, application.Attributes)
	require.NotNil(t, application.Attributes.Criticality)
	assert.Equal(t, "LOW", application.Attributes.Criticality.Type)
}

func TestGetAppHubApplicationAttrsWithClientReportsAMissingApplication(t *testing.T) {
	t.Parallel()

	// A caller who asks for an application that is not there should be told that, rather than be
	// handed the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetAppHubApplicationAttrsWithClient(context.Background(), newFakeAppHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetAppHubServiceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a App Hub service the the library suite App Hub service module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/applications/gw-library-parent/services/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/applications/gw-library-parent/services/gw-library-test","displayName":"terratest service","description":"created by terratest","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetAppHubServiceAttrsWithClient(context.Background(), newFakeAppHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest service", attrs.DisplayName)
	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetAppHubServiceAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a App Hub service that is not there should read a sentence about that App Hub service, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAppHubServiceAttrsWithClient(context.Background(), newFakeAppHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetAppHubWorkloadAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a App Hub workload the the library suite App Hub workload module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/applications/gw-library-parent/workloads/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/applications/gw-library-parent/workloads/gw-library-test","displayName":"terratest workload","description":"created by terratest","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetAppHubWorkloadAttrsWithClient(context.Background(), newFakeAppHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest workload", attrs.DisplayName)
	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetAppHubWorkloadAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a App Hub workload that is not there should read a sentence about that App Hub workload, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAppHubWorkloadAttrsWithClient(context.Background(), newFakeAppHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetAppHubServiceProjectAttachmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a App Hub service project attachment the the library suite App Hub service project attachment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/serviceProjectAttachments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/serviceProjectAttachments/gw-library-test","serviceProject":"projects/gw-library-other-project","state":"ACTIVE"}`))
	})

	attrs, err := gcp.GetAppHubServiceProjectAttachmentAttrsWithClient(context.Background(), newFakeAppHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-other-project", attrs.ServiceProject)
	assert.Equal(t, "ACTIVE", attrs.State)
}

func TestGetAppHubServiceProjectAttachmentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a App Hub service project attachment that is not there should read a sentence about that App Hub service project attachment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetAppHubServiceProjectAttachmentAttrsWithClient(context.Background(), newFakeAppHubService(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
