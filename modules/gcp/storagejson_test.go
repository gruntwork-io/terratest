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

// newFakeStorageJSONService points a real Cloud Storage JSON client at a local test server, so the
// Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeStorageJSONService(t *testing.T, handler http.Handler) *storagev1.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := storagev1.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetStorageFolderAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a folder the terraform-google-data-storage folder module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/b/gw-library-test/folders/terratest/"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bucket":"gw-library-test","name":"terratest/","metageneration":"1","createTime":"2026-09-23T12:00:00.000Z"}`))
	})

	folder, err := gcp.GetStorageFolderAttrsWithClient(context.Background(), newFakeStorageJSONService(t, handler), "gw-library-test", "terratest/")
	require.NoError(t, err)

	assert.Equal(t, "terratest/", folder.Name)
	assert.Equal(t, "gw-library-test", folder.Bucket)
}

func TestGetStorageManagedFolderAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a managed folder the terraform-google-data-storage managed folder module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/b/gw-library-test/managedFolders/terratest/"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bucket":"gw-library-test","name":"terratest/","metageneration":"1","createTime":"2026-09-23T12:00:00.000Z"}`))
	})

	managedFolder, err := gcp.GetStorageManagedFolderAttrsWithClient(context.Background(), newFakeStorageJSONService(t, handler), "gw-library-test", "terratest/")
	require.NoError(t, err)

	assert.Equal(t, "terratest/", managedFolder.Name)
	assert.Equal(t, "gw-library-test", managedFolder.Bucket)
}

func TestGetStorageHmacKeyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an HMAC key the terraform-google-data-storage HMAC key module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/hmacKeys/GOOG1ETERRATEST"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"accessId":"GOOG1ETERRATEST","projectId":"gw-library-test-project","serviceAccountEmail":"gw-library-test@gw-library-test-project.iam.gserviceaccount.com","state":"INACTIVE"}`))
	})

	key, err := gcp.GetStorageHmacKeyAttrsWithClient(context.Background(), newFakeStorageJSONService(t, handler), "gw-library-test-project", "GOOG1ETERRATEST")
	require.NoError(t, err)

	assert.Equal(t, "INACTIVE", key.State)
	assert.Equal(t, "gw-library-test@gw-library-test-project.iam.gserviceaccount.com", key.ServiceAccountEmail)
	assert.Equal(t, "GOOG1ETERRATEST", key.AccessId)
}

func TestGetStorageNotificationAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a notification the terraform-google-data-storage notification module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/b/gw-library-test/notificationConfigs/7"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"7","topic":"//pubsub.googleapis.com/projects/gw-library-test-project/topics/gw-library-test","payload_format":"JSON_API_V1","event_types":["OBJECT_FINALIZE"],"object_name_prefix":"terratest/","custom_attributes":{"purpose":"terratest"}}`))
	})

	notification, err := gcp.GetStorageNotificationAttrsWithClient(context.Background(), newFakeStorageJSONService(t, handler), "gw-library-test", "7")
	require.NoError(t, err)

	assert.Equal(t, "JSON_API_V1", notification.PayloadFormat)
	assert.Equal(t, []string{"OBJECT_FINALIZE"}, notification.EventTypes)
	assert.Equal(t, "terratest/", notification.ObjectNamePrefix)
	assert.Equal(t, "terratest", notification.CustomAttributes["purpose"])
}

func TestGetStorageObjectAccessControlAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an object access control the terraform-google-data-storage object access control module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/b/gw-library-test/o/terratest.txt/acl/allUsers"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bucket":"gw-library-test","object":"terratest.txt","entity":"allUsers","role":"READER","generation":"1"}`))
	})

	accessControl, err := gcp.GetStorageObjectAccessControlAttrsWithClient(context.Background(), newFakeStorageJSONService(t, handler), "gw-library-test", "terratest.txt", "allUsers")
	require.NoError(t, err)

	assert.Equal(t, "READER", accessControl.Role)
	assert.Equal(t, "allUsers", accessControl.Entity)
	assert.Equal(t, "terratest.txt", accessControl.Object)
}

func TestGetStorageDefaultObjectAccessControlAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a default object access control the terraform-google-data-storage default object access control module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/b/gw-library-test/defaultObjectAcl/allAuthenticatedUsers"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"bucket":"gw-library-test","entity":"allAuthenticatedUsers","role":"OWNER"}`))
	})

	accessControl, err := gcp.GetStorageDefaultObjectAccessControlAttrsWithClient(context.Background(), newFakeStorageJSONService(t, handler), "gw-library-test", "allAuthenticatedUsers")
	require.NoError(t, err)

	assert.Equal(t, "OWNER", accessControl.Role)
	assert.Equal(t, "allAuthenticatedUsers", accessControl.Entity)
}

func TestGetStorageManagedFolderIamPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a managed folder the
	// terraform-google-data-storage managed folder module granted access on, not a copy of any one
	// fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/b/gw-library-test/managedFolders/terratest//iam"), "unexpected path %s", r.URL.Path)

		// A conditional binding only comes back at version 3, so the read has to ask for it.
		assert.Equal(t, "3", r.URL.Query().Get("optionsRequestedPolicyVersion"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"resourceId":"projects/_/buckets/gw-library-test/managedFolders/terratest/","version":3,"bindings":[{"role":"roles/storage.objectViewer","members":["serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"],"condition":{"title":"until 2030","expression":"request.time < timestamp(\"2030-01-01T00:00:00Z\")"}}]}`))
	})

	policy, err := gcp.GetStorageManagedFolderIamPolicyAttrsWithClient(context.Background(), newFakeStorageJSONService(t, handler), "gw-library-test", "terratest/")
	require.NoError(t, err)

	require.Len(t, policy.Bindings, 1)
	assert.Equal(t, int64(3), policy.Version)
	assert.Equal(t, "roles/storage.objectViewer", policy.Bindings[0].Role)
	require.NotNil(t, policy.Bindings[0].Condition, "a conditional binding should keep its condition")
	assert.Equal(t, `request.time < timestamp("2030-01-01T00:00:00Z")`, policy.Bindings[0].Condition.Expression)
	assert.Equal(t, []string{"serviceAccount:gw-library-test@gw-library-test-project.iam.gserviceaccount.com"}, policy.Bindings[0].Members)
}
