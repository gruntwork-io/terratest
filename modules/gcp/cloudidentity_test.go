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
	"google.golang.org/api/cloudidentity/v1"
	"google.golang.org/api/option"
)

// newFakeCloudIdentityService points a real Cloud Identity client at a local test server, so the
// Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeCloudIdentityService(t *testing.T, handler http.Handler) *cloudidentity.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := cloudidentity.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetCloudIdentityGroupAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a group the Cloud Identity group module
	// created, not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/groups/01234abcd5678ef"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"groups/01234abcd5678ef","groupKey":{"id":"terratest@example.com"},"parent":"customers/C01abc234","displayName":"terratest group","description":"created by terratest","labels":{"cloudidentity.googleapis.com/groups.discussion_forum":""}}`))
	})

	group, err := gcp.GetCloudIdentityGroupAttrsWithClient(context.Background(), newFakeCloudIdentityService(t, handler), "groups/01234abcd5678ef")
	require.NoError(t, err)

	require.NotNil(t, group.GroupKey)
	assert.Equal(t, "terratest@example.com", group.GroupKey.Id)
	assert.Equal(t, "customers/C01abc234", group.Parent)
	assert.Equal(t, "terratest group", group.DisplayName)
}

func TestGetCloudIdentityGroupAttrsWithClientReportsAMissingGroup(t *testing.T) {
	t.Parallel()

	// A caller who asks for a group that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Not found."}}`))
	})

	_, err := gcp.GetCloudIdentityGroupAttrsWithClient(context.Background(), newFakeCloudIdentityService(t, handler), "groups/00000missing000")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
