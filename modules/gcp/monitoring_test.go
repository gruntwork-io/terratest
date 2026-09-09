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
