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
	"google.golang.org/api/file/v1"
	"google.golang.org/api/option"
)

// newFakeFilestoreService points a real Filestore client at a local test server, so the Google transport is
// exercised rather than a hand-written stand-in for a type we do not own.
func newFakeFilestoreService(t *testing.T, handler http.Handler) *file.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := file.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetFilestoreInstanceAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a instance the terraform-google-data-
	// storage module created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1-a/instances/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1-a/instances/gw-library-test","description":"created by terratest","tier":"BASIC_HDD","state":"READY","labels":{"purpose":"terratest"},"fileShares":[{"name":"share","capacityGb":"1024"}],"networks":[{"network":"gw-library-test","modes":["MODE_IPV4"]}]}`))
	})

	// Google sends a share's size as a JSON string and the Go client decodes it to an int64.
	instance, err := gcp.GetFilestoreInstanceAttrsWithClient(context.Background(), newFakeFilestoreService(t, handler), "gw-library-test-project", "us-central1-a", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "BASIC_HDD", instance.Tier)
	assert.Equal(t, "READY", instance.State)
	require.Len(t, instance.FileShares, 1)
	assert.Equal(t, int64(1024), instance.FileShares[0].CapacityGb)
	require.Len(t, instance.Networks, 1)
	assert.Equal(t, "gw-library-test", instance.Networks[0].Network)
}

func TestGetFilestoreInstanceAttrsWithClientMissingInstance(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the instance and everything that identifies it, and each of those is a value
	// no other part of the message contains, or its check could not fail.
	_, err := gcp.GetFilestoreInstanceAttrsWithClient(context.Background(), newFakeFilestoreService(t, handler), "gw-library-test-project", "us-central1-a", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "us-central1-a")
	require.ErrorContains(t, err, "gw-library-test-project")
}
