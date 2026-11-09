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
	runv1 "google.golang.org/api/run/v1"
)

// newFakeCloudRunV1Service points a real Cloud Run v1 client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudRunV1Service(t *testing.T, handler http.Handler) *runv1.APIService {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := runv1.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetCloudRunV1ServiceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The v1 API names a service by Knative namespace, and answers in Knative's shape rather than
	// Google's, so the path and the body both differ from the v2 read beside this one.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/namespaces/gw-library-test-project/services/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"apiVersion":"serving.knative.dev/v1","kind":"Service","metadata":{"name":"gw-library-test","namespace":"gw-library-test-project","labels":{"purpose":"terratest"}},"spec":{"template":{"spec":{"containerConcurrency":7,"containers":[{"image":"us-docker.pkg.dev/cloudrun/container/hello"}]}}}}`))
	})

	service, err := gcp.GetCloudRunV1ServiceAttrsWithClient(context.Background(), newFakeCloudRunV1Service(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.NotNil(t, service.Metadata)
	assert.Equal(t, "gw-library-test", service.Metadata.Name)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, service.Metadata.Labels)
	require.NotNil(t, service.Spec)
	require.NotNil(t, service.Spec.Template)
	require.NotNil(t, service.Spec.Template.Spec)
	assert.Equal(t, int64(7), service.Spec.Template.Spec.ContainerConcurrency)
}

func TestGetCloudRunV1ServiceAttrsWithClientReportsAMissingService(t *testing.T) {
	t.Parallel()

	// A caller who asks for a service that is not there should be told that, rather than be handed
	// the transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetCloudRunV1ServiceAttrsWithClient(context.Background(), newFakeCloudRunV1Service(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}

func TestNewCloudRunV1ServiceERefusesAnInvalidLocation(t *testing.T) {
	t.Parallel()

	// The location goes into the endpoint host, so anything but a plain location name would send the
	// request, and the caller's credentials with it, somewhere the caller did not name.
	_, err := gcp.NewCloudRunV1ServiceE(t, context.Background(), "us-central1.evil.example.com")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not a valid location")
}

func TestGetCloudRunDomainMappingAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The project is a namespace in this API, and the records the domain's owner has to publish come back in the
	// mapping's status.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/namespaces/gw-library-test-project/domainmappings/terratest.example.com"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"metadata":{"name":"terratest.example.com","namespace":"gw-library-test-project"},"spec":{"routeName":"gw-library-service","certificateMode":"AUTOMATIC"},"status":{"resourceRecords":[{"name":"terratest","rrdata":"ghs.googlehosted.com.","type":"CNAME"}]}}`))
	})

	mapping, err := gcp.GetCloudRunDomainMappingAttrsWithClient(context.Background(), newFakeCloudRunV1Service(t, handler), "gw-library-test-project", "us-central1", "terratest.example.com")
	require.NoError(t, err)

	require.NotNil(t, mapping.Spec)
	assert.Equal(t, "gw-library-service", mapping.Spec.RouteName)
	assert.Equal(t, "AUTOMATIC", mapping.Spec.CertificateMode)
	require.NotNil(t, mapping.Status)
	require.Len(t, mapping.Status.ResourceRecords, 1)
	assert.Equal(t, "CNAME", mapping.Status.ResourceRecords[0].Type)
}

func TestGetCloudRunDomainMappingAttrsWithClientReportsAMissingMapping(t *testing.T) {
	t.Parallel()

	// A caller who names a domain that is not mapped should read a sentence about that domain, not a status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetCloudRunDomainMappingAttrsWithClient(context.Background(), newFakeCloudRunV1Service(t, handler), "gw-library-test-project", "us-central1", "missing.example.com")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing.example.com")
	assert.Contains(t, err.Error(), "does not exist")
}
