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
	"google.golang.org/api/cloudresourcemanager/v3"
	"google.golang.org/api/option"
)

// newFakeResourceManagerService points a real Cloud Resource Manager client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeResourceManagerService(t *testing.T, handler http.Handler) *cloudresourcemanager.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := cloudresourcemanager.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetTagKeyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a tag key the terraform-google-management tag key module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/v3/tagKeys/281476187767567"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"tagKeys/281476187767567","parent":"projects/37950160017","shortName":"gw-library-test","description":"created by terratest","purpose":"GCE_FIREWALL"}`))
	})

	tagKey, err := gcp.GetTagKeyAttrsWithClient(context.Background(), newFakeResourceManagerService(t, handler), "281476187767567")
	require.NoError(t, err)

	assert.Equal(t, "gw-library-test", tagKey.ShortName)
	assert.Equal(t, "created by terratest", tagKey.Description)
	assert.Equal(t, "projects/37950160017", tagKey.Parent)
}

func TestGetTagValueAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a tag value the terraform-google-management tag value module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/v3/tagValues/281478044408593"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"tagValues/281478044408593","parent":"tagKeys/281476187767567","shortName":"terratest","description":"created by terratest"}`))
	})

	tagValue, err := gcp.GetTagValueAttrsWithClient(context.Background(), newFakeResourceManagerService(t, handler), "281478044408593")
	require.NoError(t, err)

	assert.Equal(t, "terratest", tagValue.ShortName)
	assert.Equal(t, "tagKeys/281476187767567", tagValue.Parent)
}

func TestGetLienAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a lien the terraform-google-management resource manager lien module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/v3/liens/p1061677975479-l8a72e4b4-1d19-4e5c-9e19-3f1a3a4e1a2b"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"p1061677975479-l8a72e4b4-1d19-4e5c-9e19-3f1a3a4e1a2b","parent":"projects/37950160017","restrictions":["resourcemanager.projects.delete"],"origin":"terratest","reason":"created by terratest"}`))
	})

	lien, err := gcp.GetLienAttrsWithClient(context.Background(), newFakeResourceManagerService(t, handler), "p1061677975479-l8a72e4b4-1d19-4e5c-9e19-3f1a3a4e1a2b")
	require.NoError(t, err)

	assert.Equal(t, []string{"resourcemanager.projects.delete"}, lien.Restrictions)
	assert.Equal(t, "terratest", lien.Origin)
	assert.Equal(t, "created by terratest", lien.Reason)
}

func TestGetTagBindingsAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a project the terraform-google-management
	// tag binding module tagged, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/tagBindings"), "unexpected path %s", r.URL.Path)
		assert.Equal(t, "//cloudresourcemanager.googleapis.com/projects/37950160017", r.URL.Query().Get("parent"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"tagBindings":[{"name":"tagBindings/%2F%2Fcloudresourcemanager.googleapis.com%2Fprojects%2F37950160017/tagValues/281478044408593","parent":"//cloudresourcemanager.googleapis.com/projects/37950160017","tagValue":"tagValues/281478044408593"}]}`))
	})

	bindings, err := gcp.GetTagBindingsAttrsWithClient(context.Background(), newFakeResourceManagerService(t, handler), "//cloudresourcemanager.googleapis.com/projects/37950160017")
	require.NoError(t, err)

	require.Len(t, bindings, 1)
	assert.Equal(t, "tagValues/281478044408593", bindings[0].TagValue)
}

func TestGetTagBindingsAttrsWithClientReadsEveryPage(t *testing.T) {
	t.Parallel()

	// Google pages this list, and the binding a caller is looking for may not be on the first page, so
	// a read that stopped there would report it missing.
	var requests int

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++

		w.Header().Set("Content-Type", "application/json")

		if r.URL.Query().Get("pageToken") == "" {
			_, _ = w.Write([]byte(`{"tagBindings":[{"tagValue":"tagValues/111111111111"}],"nextPageToken":"second"}`))

			return
		}

		assert.Equal(t, "second", r.URL.Query().Get("pageToken"))

		_, _ = w.Write([]byte(`{"tagBindings":[{"tagValue":"tagValues/281478044408593"}]}`))
	})

	bindings, err := gcp.GetTagBindingsAttrsWithClient(context.Background(), newFakeResourceManagerService(t, handler), "//cloudresourcemanager.googleapis.com/projects/37950160017")
	require.NoError(t, err)

	assert.Equal(t, 2, requests, "both pages should have been asked for")
	require.Len(t, bindings, 2)
	assert.Equal(t, "tagValues/281478044408593", bindings[1].TagValue)
}
