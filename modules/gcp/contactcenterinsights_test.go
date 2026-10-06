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
	"google.golang.org/api/contactcenterinsights/v1"
	"google.golang.org/api/option"
)

// newFakeContactCenterInsightsService points a real *contactcenterinsights.Service at a local test
// server, so the Google transport is exercised rather than a hand-written stand-in for a type we
// do not own.
func newFakeContactCenterInsightsService(t *testing.T, handler http.Handler) *contactcenterinsights.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := contactcenterinsights.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetInsightsViewAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-business-apps view module sets, because the
	// point of reading settings back is asserting a module configured the view it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/views/abc123"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/locations/us-central1/views/abc123",
			"displayName":"terratest view",
			"value":"medium = \"CHAT\"",
			"createTime":"2026-09-14T00:00:00Z"
		}`))
	})

	view, err := gcp.GetInsightsViewAttrsWithClient(context.Background(), newFakeContactCenterInsightsService(t, handler), "gw-library-test-project", "us-central1", "abc123")
	require.NoError(t, err)

	assert.Equal(t, "terratest view", view.DisplayName)
	assert.Equal(t, `medium = "CHAT"`, view.Value)
	assert.Equal(t, "projects/gw-library-test-project/locations/us-central1/views/abc123", view.Name)
}

func TestGetInsightsViewAttrsWithClientMissingView(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the view and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetInsightsViewAttrsWithClient(context.Background(), newFakeContactCenterInsightsService(t, handler), "gw-library-test-project", "us-central1", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetContactCenterInsightsAnalysisRuleAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a rule the terraform-google-business-apps analysis rule module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/analysisRules/1234567890"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/37950160017/locations/us-central1/analysisRules/1234567890","displayName":"terratest rule","conversationFilter":"agent_id=\"terratest\"","analysisPercentage":0.5,"active":false}`))
	})

	rule, err := gcp.GetContactCenterInsightsAnalysisRuleAttrsWithClient(context.Background(), newFakeContactCenterInsightsService(t, handler), "gw-library-test-project", "us-central1", "1234567890")
	require.NoError(t, err)

	assert.Equal(t, "terratest rule", rule.DisplayName)
	assert.InDelta(t, 0.5, rule.AnalysisPercentage, 0.001)
	assert.False(t, rule.Active)
}

func TestGetContactCenterInsightsAssessmentRuleAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a rule the terraform-google-business-apps assessment rule module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/assessmentRules/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/assessmentRules/gw-library-test","displayName":"terratest assessment rule","active":true,"sampleRule":{"samplePercentage":5}}`))
	})

	rule, err := gcp.GetContactCenterInsightsAssessmentRuleAttrsWithClient(context.Background(), newFakeContactCenterInsightsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest assessment rule", rule.DisplayName)
	assert.True(t, rule.Active)
	require.NotNil(t, rule.SampleRule)
	assert.InDelta(t, 5.0, rule.SampleRule.SamplePercentage, 0.001)
}

func TestGetContactCenterInsightsAutoLabelingRuleAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a rule the terraform-google-business-apps auto labeling rule module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/autoLabelingRules/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/autoLabelingRules/gw-library-test","displayName":"terratest labeling rule","description":"created by terratest","labelKey":"purpose","active":true}`))
	})

	rule, err := gcp.GetContactCenterInsightsAutoLabelingRuleAttrsWithClient(context.Background(), newFakeContactCenterInsightsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest labeling rule", rule.DisplayName)
	assert.Equal(t, "purpose", rule.LabelKey)
	assert.True(t, rule.Active)
}

func TestGetContactCenterInsightsQaScorecardAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a scorecard the terraform-google-business-apps QA scorecard module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/qaScorecards/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/qaScorecards/gw-library-test","displayName":"terratest scorecard","description":"created by terratest"}`))
	})

	scorecard, err := gcp.GetContactCenterInsightsQaScorecardAttrsWithClient(context.Background(), newFakeContactCenterInsightsService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest scorecard", scorecard.DisplayName)
	assert.Equal(t, "created by terratest", scorecard.Description)
}
