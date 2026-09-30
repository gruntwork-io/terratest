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
	storagev1 "google.golang.org/api/storage/v1"
)

// newFakeStorageRESTService points a real Cloud Storage REST client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeStorageRESTService(t *testing.T, handler http.Handler) *storagev1.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := storagev1.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDefaultObjectACLAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a default ACL entry the terraform-google-data-storage default object ACL module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/b/gw-library-test/defaultObjectAcl/allUsers"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"kind":"storage#objectAccessControl","bucket":"gw-library-test","entity":"allUsers","role":"READER"}`))
	})

	acl, err := gcp.GetDefaultObjectACLAttrsWithClient(context.Background(), newFakeStorageRESTService(t, handler), "gw-library-test", "allUsers")
	require.NoError(t, err)

	assert.Equal(t, "allUsers", acl.Entity)
	assert.Equal(t, "READER", acl.Role)
}

func TestGetObjectACLAttrsWithClient(t *testing.T) {
	t.Parallel()

	// An object ACL entry names the object as well as the bucket and the entity, so the path is one
	// segment deeper than the default one.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/b/gw-library-test/o/terratest.txt/acl/allUsers"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"kind":"storage#objectAccessControl","bucket":"gw-library-test","object":"terratest.txt","entity":"allUsers","role":"READER"}`))
	})

	acl, err := gcp.GetObjectACLAttrsWithClient(context.Background(), newFakeStorageRESTService(t, handler), "gw-library-test", "terratest.txt", "allUsers")
	require.NoError(t, err)

	assert.Equal(t, "terratest.txt", acl.Object)
	assert.Equal(t, "allUsers", acl.Entity)
	assert.Equal(t, "READER", acl.Role)
}

func TestGetDefaultObjectACLAttrsWithClientReportsAMissingEntry(t *testing.T) {
	t.Parallel()

	// A caller who asks for an entry that is not there should be told that, rather than handed the
	// transport's own wording for a 404. A bucket with uniform access is a different case: Google
	// refuses the request outright with a 400, so it is not what this fixture covers.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not Found"}}`))
	})

	_, err := gcp.GetDefaultObjectACLAttrsWithClient(context.Background(), newFakeStorageRESTService(t, handler), "gw-library-test", "allUsers")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no default object ACL entry")
}
