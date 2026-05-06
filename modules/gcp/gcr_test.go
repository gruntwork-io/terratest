package gcp_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeRegistry runs an httptest server mimicking the minimal Docker Registry v2 endpoints
// the go-containerregistry library exercises (ping, tags list, manifest delete).
type fakeRegistry struct {
	deletedPaths map[string]int
	host         string
	listBody     string
	deleteStatus int
}

func newFakeRegistry(t *testing.T, listBody string) *fakeRegistry {
	t.Helper()

	fr := &fakeRegistry{deletedPaths: map[string]int{}, listBody: listBody, deleteStatus: http.StatusOK}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/v2/":
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodGet && strings.HasSuffix(r.URL.Path, "/tags/list"):
			_, _ = w.Write([]byte(fr.listBody))
		case r.Method == http.MethodDelete && strings.Contains(r.URL.Path, "/manifests/"):
			fr.deletedPaths[r.URL.Path]++

			w.WriteHeader(fr.deleteStatus)
		default:
			http.Error(w, "unexpected", http.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)

	fr.host = u.Host

	return fr
}

func anonymousGCRClient() *gcp.GCRClient {
	return &gcp.GCRClient{Authenticator: authn.Anonymous, Transport: http.DefaultTransport}
}

func TestDeleteGCRImageRefWithClient(t *testing.T) {
	t.Parallel()

	t.Run("deletes tag reference", func(t *testing.T) {
		t.Parallel()

		fr := newFakeRegistry(t, "")

		require.NoError(t, gcp.DeleteGCRImageRefWithClient(context.Background(), anonymousGCRClient(), fr.host+"/p/r:latest"))
		assert.Equal(t, 1, fr.deletedPaths["/v2/p/r/manifests/latest"])
	})

	t.Run("bad reference errors", func(t *testing.T) {
		t.Parallel()

		err := gcp.DeleteGCRImageRefWithClient(context.Background(), anonymousGCRClient(), "::not-a-ref")
		require.ErrorContains(t, err, "failed to parse reference")
	})
}

func TestDeleteGCRRepoWithClient(t *testing.T) {
	t.Parallel()

	t.Run("deletes :latest and each listed digest", func(t *testing.T) {
		t.Parallel()

		d1 := "sha256:1111111111111111111111111111111111111111111111111111111111111111"
		d2 := "sha256:2222222222222222222222222222222222222222222222222222222222222222"

		listBody := fmt.Sprintf(`{"child":[],"manifest":{%q:{"tag":["t1"]},%q:{"tag":["t2"]}},"tags":["latest","t1","t2"]}`, d1, d2)
		fr := newFakeRegistry(t, listBody)

		require.NoError(t, gcp.DeleteGCRRepoWithClient(context.Background(), anonymousGCRClient(), fr.host+"/p/r"))
		assert.Equal(t, 1, fr.deletedPaths["/v2/p/r/manifests/latest"])
		assert.Equal(t, 1, fr.deletedPaths["/v2/p/r/manifests/"+d1])
		assert.Equal(t, 1, fr.deletedPaths["/v2/p/r/manifests/"+d2])
	})

	t.Run("list returns non-JSON", func(t *testing.T) {
		t.Parallel()

		fr := newFakeRegistry(t, "notjson")

		err := gcp.DeleteGCRRepoWithClient(context.Background(), anonymousGCRClient(), fr.host+"/p/r")
		require.ErrorContains(t, err, "failed to list tags")
	})
}
