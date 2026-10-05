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

func TestNewDialogflowServiceERefusesABadLocation(t *testing.T) {
	t.Parallel()

	// Each of these would build a URL pointing somewhere other than Google, so the constructor has
	// to refuse them before the endpoint is built.
	for _, location := range []string{"us/../evil.com", "evil.com", "us:8080", "user@evil.com", "US", ""} {
		_, err := gcp.NewDialogflowServiceE(t, context.Background(), location)
		require.ErrorContains(t, err, "not a valid location", "location %q should be refused", location)
	}

	// A real one is accepted, and so is the global one that builds no endpoint.
	for _, location := range []string{"us-central1", "global"} {
		_, err := gcp.NewDialogflowServiceE(t, context.Background(), location)
		require.NoError(t, err, "location %q should be accepted", location)
	}
}

func TestGetDialogflowCXGenerativeSettingsAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dialogflow CX generative settings the the library suite Dialogflow CX generative settings module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-parent/generativeSettings"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/agents/gw-library-parent/generativeSettings","languageCode":"en","fallbackSettings":{"selectedPrompt":"terratest"}}`))
	})

	attrs, err := gcp.GetDialogflowCXGenerativeSettingsAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "en")
	require.NoError(t, err)

	assert.Equal(t, "en", attrs.LanguageCode)
}

func TestGetDialogflowCXGenerativeSettingsAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dialogflow CX generative settings that is not there should read a sentence about that Dialogflow CX generative settings, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDialogflowCXGenerativeSettingsAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDialogflowCXToolVersionAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dialogflow CX tool version the the library suite Dialogflow CX tool version module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-parent/tools/gw-library-tool/versions/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/agents/gw-library-parent/tools/gw-library-tool/versions/gw-library-test","displayName":"terratest tool version"}`))
	})

	attrs, err := gcp.GetDialogflowCXToolVersionAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-tool", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest tool version", attrs.DisplayName)
}

func TestGetDialogflowCXToolVersionAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dialogflow CX tool version that is not there should read a sentence about that Dialogflow CX tool version, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDialogflowCXToolVersionAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-parent", "gw-library-tool", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
