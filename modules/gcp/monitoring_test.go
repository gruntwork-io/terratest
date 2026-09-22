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
	"google.golang.org/api/monitoring/v3"
	"google.golang.org/api/option"
)

// newFakeMonitoringService points a real *monitoring.Service at a local test server, so the Google
// transport is exercised rather than a hand-written stand-in for a type we do not own.
func newFakeMonitoringService(t *testing.T, handler http.Handler) *monitoring.Service {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	service, err := monitoring.NewService(context.Background(),
		option.WithEndpoint(server.URL), option.WithoutAuthentication())
	require.NoError(t, err)

	return service
}

func TestGetNotificationChannelAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-observability notification channel module sets,
	// because the point of reading settings back is asserting a module configured the channel it
	// was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/notificationChannels/1234567890"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/notificationChannels/1234567890",
			"type":"email",
			"displayName":"terratest channel",
			"description":"created by terratest",
			"labels":{"email_address":"terratest@example.com"},
			"userLabels":{"purpose":"terratest"},
			"enabled":false
		}`))
	})

	channel, err := gcp.GetNotificationChannelAttrsWithClient(context.Background(), newFakeMonitoringService(t, handler), "gw-library-test-project", "1234567890")
	require.NoError(t, err)

	assert.Equal(t, "email", channel.Type)
	assert.Equal(t, "terratest channel", channel.DisplayName)
	assert.Equal(t, "created by terratest", channel.Description)
	assert.Equal(t, "terratest@example.com", channel.Labels["email_address"])
	assert.Equal(t, map[string]string{"purpose": "terratest"}, channel.UserLabels)
	assert.False(t, channel.Enabled)
}

func TestGetNotificationChannelAttrsWithClientMissingChannel(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the channel and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetNotificationChannelAttrsWithClient(context.Background(), newFakeMonitoringService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetUptimeCheckConfigAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The values are the ones the terraform-google-observability uptime check module sets, because
	// the point of reading settings back is asserting a module configured the check it was asked for.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/uptimeCheckConfigs/abc123"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/uptimeCheckConfigs/abc123",
			"displayName":"terratest check",
			"period":"300s",
			"timeout":"10s",
			"checkerType":"STATIC_IP_CHECKERS",
			"httpCheck":{"path":"/healthz","port":443,"useSsl":true,"validateSsl":true,"requestMethod":"GET"},
			"monitoredResource":{"type":"uptime_url","labels":{"host":"example.com","project_id":"gw-library-test-project"}},
			"userLabels":{"managed-by":"terratest"}
		}`))
	})

	check, err := gcp.GetUptimeCheckConfigAttrsWithClient(context.Background(), newFakeMonitoringService(t, handler), "gw-library-test-project", "abc123")
	require.NoError(t, err)

	assert.Equal(t, "terratest check", check.DisplayName)
	assert.Equal(t, "300s", check.Period)
	assert.Equal(t, "10s", check.Timeout)
	require.NotNil(t, check.HttpCheck)
	assert.Equal(t, "/healthz", check.HttpCheck.Path)
	assert.Equal(t, int64(443), check.HttpCheck.Port)
	assert.True(t, check.HttpCheck.UseSsl)
	require.NotNil(t, check.MonitoredResource)
	assert.Equal(t, "example.com", check.MonitoredResource.Labels["host"])
}

func TestGetUptimeCheckConfigAttrsWithClientMissingCheck(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the check and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetUptimeCheckConfigAttrsWithClient(context.Background(), newFakeMonitoringService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetAlertPolicyAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an alert policy the
	// terraform-google-observability module created, not a copy of any one fixture's values. Google
	// assigns the policy's id, so the caller passes the id it got back.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/alertPolicies/1234567890"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/gw-library-test-project/alertPolicies/1234567890",
			"displayName":"terratest policy",
			"combiner":"OR",
			"enabled":false,
			"severity":"WARNING",
			"userLabels":{"purpose":"terratest"},
			"conditions":[{
				"displayName":"uptime check failing",
				"conditionThreshold":{"comparison":"COMPARISON_GT","thresholdValue":1,"duration":"300s"}
			}]
		}`))
	})

	policy, err := gcp.GetAlertPolicyAttrsWithClient(context.Background(), newFakeMonitoringService(t, handler), "gw-library-test-project", "1234567890")
	require.NoError(t, err)

	assert.Equal(t, "terratest policy", policy.DisplayName)
	assert.Equal(t, "OR", policy.Combiner)
	assert.False(t, policy.Enabled)
	assert.Equal(t, "WARNING", policy.Severity)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, policy.UserLabels)
	require.Len(t, policy.Conditions, 1)
	require.NotNil(t, policy.Conditions[0].ConditionThreshold)
	assert.Equal(t, "COMPARISON_GT", policy.Conditions[0].ConditionThreshold.Comparison)
	assert.Equal(t, "300s", policy.Conditions[0].ConditionThreshold.Duration)
}

func TestGetAlertPolicyAttrsWithClientMissingPolicy(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the policy and the project as well as saying it is absent, so all three are
	// asserted rather than only the phrase.
	_, err := gcp.GetAlertPolicyAttrsWithClient(context.Background(), newFakeMonitoringService(t, handler), "gw-library-test-project", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "gone")
	require.ErrorContains(t, err, "gw-library-test-project")
}

func TestGetServiceLevelObjectiveAttrsWithClient(t *testing.T) {
	t.Parallel()

	// The response is shaped like the one Google returns for an SLO the terraform-google-observability
	// module created, not a copy of any one fixture's values. An SLO is named by the service it
	// measures as well as its own id, so both have to reach the request.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.True(t, strings.HasSuffix(r.URL.Path, "/projects/gw-library-test-project/services/gw-library-test/serviceLevelObjectives/gw-library-test-slo"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"projects/123/services/gw-library-test/serviceLevelObjectives/gw-library-test-slo",
			"displayName":"terratest slo",
			"goal":0.99,
			"rollingPeriod":"604800s",
			"userLabels":{"purpose":"terratest"}
		}`))
	})

	slo, err := gcp.GetServiceLevelObjectiveAttrsWithClient(context.Background(), newFakeMonitoringService(t, handler), "gw-library-test-project", "gw-library-test", "gw-library-test-slo")
	require.NoError(t, err)

	assert.Equal(t, "terratest slo", slo.DisplayName)
	assert.InDelta(t, 0.99, slo.Goal, 1e-9)
	assert.Equal(t, "604800s", slo.RollingPeriod)
	assert.Equal(t, map[string]string{"purpose": "terratest"}, slo.UserLabels)
}

func TestGetServiceLevelObjectiveAttrsWithClientMissingSLO(t *testing.T) {
	t.Parallel()

	handler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// The error names the SLO, its service and the project as well as saying it is absent, since an
	// SLO is only identified by all three together. The service is named so that no other value in
	// the error contains it, or its check could not fail.
	_, err := gcp.GetServiceLevelObjectiveAttrsWithClient(context.Background(), newFakeMonitoringService(t, handler), "gw-library-test-project", "gw-service", "gone")
	require.ErrorContains(t, err, "does not exist")
	require.ErrorContains(t, err, "SLO gone ")
	require.ErrorContains(t, err, "service gw-service ")
	require.ErrorContains(t, err, "gw-library-test-project")
}
