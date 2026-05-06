package gcp_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
)

// newFakeStorageClient points a *storage.Client at a local httptest server. The SDK uses a
// JSON/REST transport by default, so option.WithEndpoint routes calls through the fake.
func newFakeStorageClient(t *testing.T, handler http.Handler) *storage.Client {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	client, err := storage.NewClient(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)
	t.Cleanup(func() { _ = client.Close() })

	return client
}

func TestCreateStorageBucketWithClient(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "p", r.URL.Query().Get("project"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"kind":"storage#bucket","id":"b","name":"b"}`))
	})

	client := newFakeStorageClient(t, handler)

	require.NoError(t, gcp.CreateStorageBucketWithClient(context.Background(), client, "p", "b", nil))
}

func TestDeleteStorageBucketWithClient(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Contains(t, r.URL.Path, "/b/b")
		w.WriteHeader(http.StatusNoContent)
	})

	client := newFakeStorageClient(t, handler)

	require.NoError(t, gcp.DeleteStorageBucketWithClient(context.Background(), client, "b"))
}

func TestAssertStorageBucketExistsWithClient(t *testing.T) {
	t.Parallel()

	t.Run("exists", func(t *testing.T) {
		t.Parallel()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")

			if strings.HasSuffix(r.URL.Path, "/b/b") {
				_, _ = w.Write([]byte(`{"kind":"storage#bucket","id":"b","name":"b"}`))

				return
			}

			_, _ = w.Write([]byte(`{"kind":"storage#objects","items":[]}`))
		})

		client := newFakeStorageClient(t, handler)

		require.NoError(t, gcp.AssertStorageBucketExistsWithClient(context.Background(), client, "b"))
	})

	t.Run("missing", func(t *testing.T) {
		t.Parallel()

		handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		})

		client := newFakeStorageClient(t, handler)

		require.Error(t, gcp.AssertStorageBucketExistsWithClient(context.Background(), client, "b"))
	})
}

func TestReadBucketObjectWithClient(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = w.Write([]byte("hello"))
	})

	client := newFakeStorageClient(t, handler)

	r, err := gcp.ReadBucketObjectWithClient(context.Background(), client, "b", "o.txt")
	require.NoError(t, err)

	got, err := io.ReadAll(r)
	require.NoError(t, err)
	assert.Equal(t, "hello", string(got))
}

func TestEmptyStorageBucketWithClient(t *testing.T) {
	t.Parallel()

	var deleted []string

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"kind":"storage#objects","items":[{"name":"a"},{"name":"b"}]}`))
		case http.MethodDelete:
			if parts := strings.Split(r.URL.Path, "/o/"); len(parts) == 2 {
				deleted = append(deleted, parts[1])
			}

			w.WriteHeader(http.StatusNoContent)
		}
	})

	client := newFakeStorageClient(t, handler)

	require.NoError(t, gcp.EmptyStorageBucketWithClient(context.Background(), client, "b"))
	assert.ElementsMatch(t, []string{"a", "b"}, deleted)
}
