package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/monitoring/v3"
	"google.golang.org/api/option"
)

// GetNotificationChannelAttrs returns the settings Google Cloud holds for the given notification
// channel, so a test can assert on what was actually created rather than only that it exists. The
// id is the one Google assigns, which a caller reads from whatever created the channel.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetNotificationChannelAttrs(t testing.TestingT, ctx context.Context, projectID string, channelID string) *monitoring.NotificationChannel {
	channel, err := GetNotificationChannelAttrsE(t, ctx, projectID, channelID)
	require.NoError(t, err)

	return channel
}

// GetNotificationChannelAttrsE returns the settings Google Cloud holds for the given notification
// channel.
// The ctx parameter supports cancellation and timeouts.
func GetNotificationChannelAttrsE(t testing.TestingT, ctx context.Context, projectID string, channelID string) (*monitoring.NotificationChannel, error) {
	logger.Default.Logf(t, "Getting settings for notification channel %s in project %s", channelID, projectID)

	service, err := NewMonitoringServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetNotificationChannelAttrsWithClient(ctx, service, projectID, channelID)
}

// GetNotificationChannelAttrsWithClient returns the settings Google Cloud holds for the given
// notification channel using the supplied *monitoring.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see monitoring_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetNotificationChannelAttrsWithClient(ctx context.Context, service *monitoring.Service, projectID string, channelID string) (*monitoring.NotificationChannel, error) {
	name := fmt.Sprintf("projects/%s/notificationChannels/%s", projectID, channelID)

	channel, err := service.Projects.NotificationChannels.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("notification channel %s does not exist in project %s", channelID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for notification channel %s in project %s: %w", channelID, projectID, err)
	}

	return channel, nil
}

// NewMonitoringServiceE creates a Cloud Monitoring service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewMonitoringServiceE(t testing.TestingT, ctx context.Context) (*monitoring.Service, error) {
	return monitoring.NewService(ctx, append(withOptions(), option.WithScopes(monitoring.CloudPlatformScope))...)
}
