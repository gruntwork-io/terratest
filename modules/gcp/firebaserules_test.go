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
	"google.golang.org/api/firebaserules/v1"
	"google.golang.org/api/option"
)

// newFakeFirebaseRulesService points a real *firebaserules.Service at a local test server, so the
// Google transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeFirebaseRulesService(t *testing.T, handler http.Handler) *firebaserules.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := firebaserules.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetFirebaseRulesetAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-firebase ruleset module sets, because the point
	// of reading settings back is asserting a module configured the ruleset it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/rulesets/abc123"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/rulesets/abc123",
			"source":{"files":[{"name":"firestore.rules","content":"service cloud.firestore { match /{document=**} { allow read, write: if false; } }"}]},
			"metadata":{"services":["cloud.firestore"]},
			"createTime":"2026-09-14T00:00:00Z"
		}`))
	})

	ruleset, err := gcp.GetFirebaseRulesetAttrsWithClient(context.Background(), newFakeFirebaseRulesService(t, handler), "gw-library-test-project", "abc123")
	require.NoError(t, err)

	require.Len(t, ruleset.Source.Files, 1)
	assert.Equal(t, "firestore.rules", ruleset.Source.Files[0].Name)
	assert.Contains(t, ruleset.Source.Files[0].Content, "allow read, write: if false")
	assert.Equal(t, []string{"cloud.firestore"}, ruleset.Metadata.Services)
}

func TestGetFirebaseRulesetAttrsWithClientMissingRuleset(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the ruleset and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetFirebaseRulesetAttrsWithClient(context.Background(), newFakeFirebaseRulesService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}
