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
	"google.golang.org/api/dialogflow/v3"
	"google.golang.org/api/option"
)

// newFakeDialogflowService points a real *dialogflow.Service at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDialogflowService(t *testing.T, handler http.Handler) *dialogflow.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := dialogflow.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDialogflowAgentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-ml agent module sets, because the point of
	// reading settings back is asserting a module configured the agent it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/locations/global/agents/abc123"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/global/agents/abc123",
			"displayName":"terratest agent",
			"description":"created by terratest",
			"defaultLanguageCode":"en",
			"timeZone":"America/New_York",
			"enableSpellCorrection":true
		}`))
	})

	agent, err := gcp.GetDialogflowAgentAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "global", "abc123")
	require.NoError(t, err)

	assert.Equal(t, "terratest agent", agent.DisplayName)
	assert.Equal(t, "created by terratest", agent.Description)
	assert.Equal(t, "en", agent.DefaultLanguageCode)
	assert.Equal(t, "America/New_York", agent.TimeZone)
	assert.True(t, agent.EnableSpellCorrection)
}

func TestGetDialogflowAgentAttrsWithClientMissingAgent(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the agent and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetDialogflowAgentAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "global", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
