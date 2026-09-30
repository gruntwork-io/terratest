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

func TestGetDialogflowCXIntentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an intent a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-agent/intents/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest intent","description":"created by terratest","priority":250000,"isFallback":false,"labels":{"purpose":"terratest"}}`))
	})

	result, err := gcp.GetDialogflowCXIntentAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-agent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest intent", result.DisplayName)
	assert.Equal(t, int64(250000), result.Priority)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, result.Labels)
}

func TestGetDialogflowCXEntityTypeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an entity type a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-agent/entityTypes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest entity type","kind":"KIND_MAP","autoExpansionMode":"AUTO_EXPANSION_MODE_DEFAULT","enableFuzzyExtraction":true,"entities":[{"value":"terratest","synonyms":["terratest","tt"]}]}`))
	})

	result, err := gcp.GetDialogflowCXEntityTypeAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-agent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "KIND_MAP", result.Kind)
	assert.True(t, result.EnableFuzzyExtraction)
	require.Len(t, result.Entities, 1)
	assert.Equal(t, "terratest", result.Entities[0].Value)
}

func TestGetDialogflowCXFlowAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a flow a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-agent/flows/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest flow","description":"created by terratest"}`))
	})

	result, err := gcp.GetDialogflowCXFlowAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-agent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest flow", result.DisplayName)
	assert.Equal(t, "created by terratest", result.Description)
}

func TestGetDialogflowCXWebhookAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a webhook a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-agent/webhooks/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest webhook","disabled":true,"timeout":"7s","genericWebService":{"uri":"https://terratest.example.com/hook"}}`))
	})

	result, err := gcp.GetDialogflowCXWebhookAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-agent", "gw-library-test")
	require.NoError(t, err)

	assert.True(t, result.Disabled)
	assert.Equal(t, "7s", result.Timeout)
	require.NotNil(t, result.GenericWebService)
	assert.Equal(t, "https://terratest.example.com/hook", result.GenericWebService.Uri)
}

func TestGetDialogflowCXEnvironmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an environment a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-agent/environments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest environment","description":"created by terratest","versionConfigs":[{"version":"projects/p/locations/us-central1/agents/a/flows/f/versions/1"}]}`))
	})

	result, err := gcp.GetDialogflowCXEnvironmentAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-agent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest environment", result.DisplayName)
	require.Len(t, result.VersionConfigs, 1)
}

func TestGetDialogflowCXTestCaseAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a test case a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-agent/testCases/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest test case","notes":"created by terratest","tags":["#terratest"]}`))
	})

	result, err := gcp.GetDialogflowCXTestCaseAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-agent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest test case", result.DisplayName)
	assert.Equal(t, []string{"#terratest"}, result.Tags)
}

func TestGetDialogflowCXGeneratorAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a generator a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-agent/generators/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest generator","promptText":{"text":"Summarise the conversation for terratest"}}`))
	})

	result, err := gcp.GetDialogflowCXGeneratorAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-agent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest generator", result.DisplayName)
	require.NotNil(t, result.PromptText)
	assert.Contains(t, result.PromptText.Text, "terratest")
}

func TestGetDialogflowCXToolAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a tool a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-agent/tools/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest tool","description":"created by terratest","toolType":"CUSTOMIZED_TOOL","openApiSpec":{"textSchema":"openapi: 3.0.0"}}`))
	})

	result, err := gcp.GetDialogflowCXToolAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-agent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest tool", result.DisplayName)
	assert.Equal(t, "CUSTOMIZED_TOOL", result.ToolType)
	require.NotNil(t, result.OpenApiSpec)
}

func TestGetDialogflowCXPlaybookAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a playbook a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-agent/playbooks/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest playbook","goal":"Answer terratest questions"}`))
	})

	result, err := gcp.GetDialogflowCXPlaybookAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-agent", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest playbook", result.DisplayName)
	assert.Contains(t, result.Goal, "terratest")
}

func TestGetDialogflowCXPageAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a page a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-agent/flows/gw-library-flow/pages/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest page","description":"created by terratest"}`))
	})

	result, err := gcp.GetDialogflowCXPageAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-agent", "gw-library-flow", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest page", result.DisplayName)
}

func TestGetDialogflowCXVersionAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a version a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/agents/gw-library-agent/flows/gw-library-flow/versions/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest version","description":"created by terratest","state":"RUNNING"}`))
	})

	result, err := gcp.GetDialogflowCXVersionAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-agent", "gw-library-flow", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest version", result.DisplayName)
	assert.Equal(t, "RUNNING", result.State)
}

func TestGetDialogflowCXSecuritySettingsAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for security settings a Gruntwork module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/securitySettings/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"...","displayName":"terratest security settings","retentionWindowDays":7,"purgeDataTypes":["DIALOGFLOW_HISTORY"],"redactionStrategy":"REDACT_WITH_SERVICE"}`))
	})

	result, err := gcp.GetDialogflowCXSecuritySettingsAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest security settings", result.DisplayName)
	assert.Equal(t, int64(7), result.RetentionWindowDays)
	assert.Equal(t, []string{"DIALOGFLOW_HISTORY"}, result.PurgeDataTypes)
}

func TestGetDialogflowCXIntentAttrsWithClientReportsAMissingIntent(t *testing.T) {
	t.Parallel()

	// A caller who asks for an intent that is not there should be told that, rather than be handed the
	// transport's own wording for a 404.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":{"code":404,"message":"Resource not found."}}`))
	})

	_, err := gcp.GetDialogflowCXIntentAttrsWithClient(context.Background(), newFakeDialogflowService(t, handler),
		"gw-library-test-project", "us-central1", "gw-library-agent", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not exist")
}
