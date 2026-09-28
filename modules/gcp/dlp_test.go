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
	"google.golang.org/api/dlp/v2"
	"google.golang.org/api/option"
)

// newFakeDLPService points a real Sensitive Data Protection client at a local test server, so the Google transport
// is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeDLPService(t *testing.T, handler http.Handler) *dlp.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := dlp.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetDLPInspectTemplateAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a template the terraform-google-security inspect template module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/inspectTemplates/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/inspectTemplates/gw-library-test","displayName":"terratest template","description":"created by terratest","inspectConfig":{"minLikelihood":"LIKELY","includeQuote":true,"infoTypes":[{"name":"EMAIL_ADDRESS"}],"limits":{"maxFindingsPerRequest":10}}}`))
	})

	template, err := gcp.GetDLPInspectTemplateAttrsWithClient(context.Background(), newFakeDLPService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest template", template.DisplayName)
	require.NotNil(t, template.InspectConfig)
	assert.Equal(t, "LIKELY", template.InspectConfig.MinLikelihood)
	assert.True(t, template.InspectConfig.IncludeQuote)
	require.Len(t, template.InspectConfig.InfoTypes, 1)
	assert.Equal(t, "EMAIL_ADDRESS", template.InspectConfig.InfoTypes[0].Name)
}

func TestGetDLPDeidentifyTemplateAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a template the terraform-google-security deidentify template module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/deidentifyTemplates/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/deidentifyTemplates/gw-library-test","displayName":"terratest template","description":"created by terratest","deidentifyConfig":{"infoTypeTransformations":{"transformations":[{"primitiveTransformation":{"replaceWithInfoTypeConfig":{}}}]}}}`))
	})

	template, err := gcp.GetDLPDeidentifyTemplateAttrsWithClient(context.Background(), newFakeDLPService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest template", template.DisplayName)
	require.NotNil(t, template.DeidentifyConfig)
	require.NotNil(t, template.DeidentifyConfig.InfoTypeTransformations)
	require.Len(t, template.DeidentifyConfig.InfoTypeTransformations.Transformations, 1)
}

func TestGetDLPStoredInfoTypeAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a type the terraform-google-security stored info type module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/storedInfoTypes/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/storedInfoTypes/gw-library-test","currentVersion":{"state":"READY","config":{"displayName":"terratest type","description":"created by terratest","regex":{"pattern":"terratest-[0-9]+"}}}}`))
	})

	infoType, err := gcp.GetDLPStoredInfoTypeAttrsWithClient(context.Background(), newFakeDLPService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	require.NotNil(t, infoType.CurrentVersion)
	require.NotNil(t, infoType.CurrentVersion.Config)
	assert.Equal(t, "terratest type", infoType.CurrentVersion.Config.DisplayName)
	require.NotNil(t, infoType.CurrentVersion.Config.Regex)
	assert.Equal(t, "terratest-[0-9]+", infoType.CurrentVersion.Config.Regex.Pattern)
}

func TestGetDLPJobTriggerAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for a trigger the terraform-google-security job trigger module created, not a
	// copy of any one fixture's values.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/locations/us-central1/jobTriggers/gw-library-test"), "unexpected path %s", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"projects/gw-library-test-project/locations/us-central1/jobTriggers/gw-library-test","displayName":"terratest trigger","description":"created by terratest","status":"PAUSED","triggers":[{"schedule":{"recurrencePeriodDuration":"86400s"}}],"inspectJob":{"inspectTemplateName":"projects/gw-library-test-project/locations/us-central1/inspectTemplates/gw-library-test"}}`))
	})

	trigger, err := gcp.GetDLPJobTriggerAttrsWithClient(context.Background(), newFakeDLPService(t, handler), "gw-library-test-project", "us-central1", "gw-library-test")
	require.NoError(t, err)

	assert.Equal(t, "terratest trigger", trigger.DisplayName)
	assert.Equal(t, "PAUSED", trigger.Status)
	require.Len(t, trigger.Triggers, 1)
	require.NotNil(t, trigger.Triggers[0].Schedule)
	assert.Equal(t, "86400s", trigger.Triggers[0].Schedule.RecurrencePeriodDuration)
}
