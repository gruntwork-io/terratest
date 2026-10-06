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
	dialogflowv2 "google.golang.org/api/dialogflow/v2"
	"google.golang.org/api/option"
)

// newFakeDialogflowV2Service points a real Dialogflow ES client at an httptest server, so a read can be
// exercised against a response we control without reaching Google.
func newFakeDialogflowV2Service(t *testing.T, handler http.Handler) *dialogflowv2.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := dialogflowv2.NewService(context.Background(),
		option.WithoutAuthentication(), option.WithEndpoint(server.URL), option.WithHTTPClient(server.Client()))
	require.NoError(t, err)

	return service
}

func TestGetDialogflowV2AgentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dialogflow ES agent the the library suite Dialogflow ES agent module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/agent"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"parent":"projects/gw-library-test-project","displayName":"terratest agent","defaultLanguageCode":"en","timeZone":"America/New_York","matchMode":"MATCH_MODE_ML_ONLY","enableLogging":true}`))
	})

	attrs, err := gcp.GetDialogflowV2AgentAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project")
	require.NoError(t, err)

	assert.Equal(t, "terratest agent", attrs.DisplayName)
	assert.Equal(t, "America/New_York", attrs.TimeZone)
	assert.Equal(t, "MATCH_MODE_ML_ONLY", attrs.MatchMode)
}

func TestGetDialogflowV2AgentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dialogflow ES agent that is not there should read a sentence about that Dialogflow ES agent, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDialogflowV2AgentAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDialogflowV2ConversationProfileAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dialogflow ES conversation profile the the library suite Dialogflow ES conversation profile module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/conversationProfiles/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/conversationProfiles/gw-library-test","displayName":"terratest profile","languageCode":"en-US"}`))
	})

	attrs, err := gcp.GetDialogflowV2ConversationProfileAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest profile", attrs.DisplayName)
	assert.Equal(t, "en-US", attrs.LanguageCode)
}

func TestGetDialogflowV2ConversationProfileAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dialogflow ES conversation profile that is not there should read a sentence about that Dialogflow ES conversation profile, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDialogflowV2ConversationProfileAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDialogflowV2EncryptionSpecAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dialogflow ES encryption spec the the library suite Dialogflow ES encryption spec module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/encryptionSpec"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/encryptionSpec","kmsKey":"projects/gw-library-test-project/locations/us-central1/keyRings/gw-library-test/cryptoKeys/gw-library-test"}`))
	})

	attrs, err := gcp.GetDialogflowV2EncryptionSpecAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "us-central1")
	require.NoError(t, err)

	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1/keyRings/gw-library-test/cryptoKeys/gw-library-test", attrs.KmsKey)
}

func TestGetDialogflowV2EncryptionSpecAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dialogflow ES encryption spec that is not there should read a sentence about that Dialogflow ES encryption spec, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDialogflowV2EncryptionSpecAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDialogflowV2EntityTypeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dialogflow ES entity type the the library suite Dialogflow ES entity type module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/agent/entityTypes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/agent/entityTypes/gw-library-test","displayName":"terratest entity","kind":"KIND_MAP","autoExpansionMode":"AUTO_EXPANSION_MODE_DEFAULT"}`))
	})

	attrs, err := gcp.GetDialogflowV2EntityTypeAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest entity", attrs.DisplayName)
	assert.Equal(t, "KIND_MAP", attrs.Kind)
	assert.Equal(t, "AUTO_EXPANSION_MODE_DEFAULT", attrs.AutoExpansionMode)
}

func TestGetDialogflowV2EntityTypeAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dialogflow ES entity type that is not there should read a sentence about that Dialogflow ES entity type, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDialogflowV2EntityTypeAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDialogflowV2EnvironmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dialogflow ES environment the the library suite Dialogflow ES environment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/agent/environments/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/agent/environments/gw-library-test","description":"created by terratest","agentVersion":"projects/gw-library-test-project/agent/versions/1","state":"RUNNING"}`))
	})

	attrs, err := gcp.GetDialogflowV2EnvironmentAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "RUNNING", attrs.State)
}

func TestGetDialogflowV2EnvironmentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dialogflow ES environment that is not there should read a sentence about that Dialogflow ES environment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDialogflowV2EnvironmentAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDialogflowV2FulfillmentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dialogflow ES fulfillment the the library suite Dialogflow ES fulfillment module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/agent/fulfillment"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/agent/fulfillment","displayName":"terratest fulfillment","enabled":true,"genericWebService":{"uri":"https://terratest.example.com/hook"}}`))
	})

	attrs, err := gcp.GetDialogflowV2FulfillmentAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project")
	require.NoError(t, err)

	assert.Equal(t, "terratest fulfillment", attrs.DisplayName)
	assert.True(t, attrs.Enabled)
}

func TestGetDialogflowV2FulfillmentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dialogflow ES fulfillment that is not there should read a sentence about that Dialogflow ES fulfillment, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDialogflowV2FulfillmentAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDialogflowV2GeneratorAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dialogflow ES generator the the library suite Dialogflow ES generator module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/generators/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/generators/gw-library-test","description":"created by terratest","triggerEvent":"END_OF_HUMAN_AGENT_TURN"}`))
	})

	attrs, err := gcp.GetDialogflowV2GeneratorAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "END_OF_HUMAN_AGENT_TURN", attrs.TriggerEvent)
}

func TestGetDialogflowV2GeneratorAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dialogflow ES generator that is not there should read a sentence about that Dialogflow ES generator, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDialogflowV2GeneratorAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDialogflowV2IntentAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dialogflow ES intent the the library suite Dialogflow ES intent module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/agent/intents/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/agent/intents/gw-library-test","displayName":"terratest intent","priority":500000,"webhookState":"WEBHOOK_STATE_ENABLED","endInteraction":true}`))
	})

	attrs, err := gcp.GetDialogflowV2IntentAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest intent", attrs.DisplayName)
	assert.Equal(t, int64(500000), attrs.Priority)
	assert.Equal(t, "WEBHOOK_STATE_ENABLED", attrs.WebhookState)
}

func TestGetDialogflowV2IntentAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dialogflow ES intent that is not there should read a sentence about that Dialogflow ES intent, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDialogflowV2IntentAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDialogflowV2SipTrunkAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dialogflow ES SIP trunk the the library suite Dialogflow ES SIP trunk module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/sipTrunks/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/sipTrunks/gw-library-test","displayName":"terratest trunk","expectedHostname":["terratest.example.com"]}`))
	})

	attrs, err := gcp.GetDialogflowV2SipTrunkAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest trunk", attrs.DisplayName)
}

func TestGetDialogflowV2SipTrunkAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dialogflow ES SIP trunk that is not there should read a sentence about that Dialogflow ES SIP trunk, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDialogflowV2SipTrunkAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "us-central1", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetDialogflowV2VersionAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a Dialogflow ES agent version the the library suite Dialogflow ES agent version module created,
	// not a copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/agent/versions/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/agent/versions/gw-library-test","description":"created by terratest","versionNumber":1,"status":"READY"}`))
	})

	attrs, err := gcp.GetDialogflowV2VersionAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "created by terratest", attrs.Description)
	assert.Equal(t, "READY", attrs.Status)
}

func TestGetDialogflowV2VersionAttrsWithClientMissingResource(t *testing.T) {
	t.Parallel()

	// A caller who names a Dialogflow ES agent version that is not there should read a sentence about that Dialogflow ES agent version, not a
	// status code.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	_, err := gcp.GetDialogflowV2VersionAttrsWithClient(context.Background(), newFakeDialogflowV2Service(t, handler), "gw-library-test-project", "gw-library-missing")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "gw-library-missing")
	assert.Contains(t, err.Error(), "does not exist")
}
