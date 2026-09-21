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

// GetUptimeCheckConfigAttrs returns the settings Google Cloud holds for the given uptime check, so
// a test can assert on what was actually created rather than only that it exists. The id is the
// one Google assigns, which the check's resource name carries as its last segment.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetUptimeCheckConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, checkID string) *monitoring.UptimeCheckConfig {
	check, err := GetUptimeCheckConfigAttrsE(t, ctx, projectID, checkID)
	require.NoError(t, err)

	return check
}

// GetUptimeCheckConfigAttrsE returns the settings Google Cloud holds for the given uptime check.
// The ctx parameter supports cancellation and timeouts.
func GetUptimeCheckConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, checkID string) (*monitoring.UptimeCheckConfig, error) {
	logger.Default.Logf(t, "Getting settings for uptime check %s in project %s", checkID, projectID)

	service, err := NewMonitoringServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetUptimeCheckConfigAttrsWithClient(ctx, service, projectID, checkID)
}

// GetUptimeCheckConfigAttrsWithClient returns the settings Google Cloud holds for the given uptime
// check using the supplied *monitoring.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see monitoring_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetUptimeCheckConfigAttrsWithClient(ctx context.Context, service *monitoring.Service, projectID string, checkID string) (*monitoring.UptimeCheckConfig, error) {
	name := fmt.Sprintf("projects/%s/uptimeCheckConfigs/%s", projectID, checkID)

	check, err := service.Projects.UptimeCheckConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("uptime check %s does not exist in project %s", checkID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for uptime check %s in project %s: %w", checkID, projectID, err)
	}

	return check, nil
}

// GetAlertPolicyAttrs returns the settings Google Cloud holds for the given alert policy, so a test
// can assert on what was actually created rather than only that it exists. Google assigns a
// policy's id when it is created, so the caller passes the id it got back, the last segment of the
// policy's name.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAlertPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, policyID string) *monitoring.AlertPolicy {
	policy, err := GetAlertPolicyAttrsE(t, ctx, projectID, policyID)
	require.NoError(t, err)

	return policy
}

// GetAlertPolicyAttrsE returns the settings Google Cloud holds for the given alert policy.
// The ctx parameter supports cancellation and timeouts.
func GetAlertPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, policyID string) (*monitoring.AlertPolicy, error) {
	logger.Default.Logf(t, "Getting settings for alert policy %s in project %s", policyID, projectID)

	service, err := NewMonitoringServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAlertPolicyAttrsWithClient(ctx, service, projectID, policyID)
}

// GetAlertPolicyAttrsWithClient returns the settings Google Cloud holds for the given alert policy
// using the supplied *monitoring.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see monitoring_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAlertPolicyAttrsWithClient(ctx context.Context, service *monitoring.Service, projectID string, policyID string) (*monitoring.AlertPolicy, error) {
	name := fmt.Sprintf("projects/%s/alertPolicies/%s", projectID, policyID)

	policy, err := service.Projects.AlertPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("alert policy %s does not exist in project %s", policyID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for alert policy %s in project %s: %w", policyID, projectID, err)
	}

	return policy, nil
}

// GetServiceLevelObjectiveAttrs returns the settings Google Cloud holds for the given service level
// objective, so a test can assert on what was actually created rather than only that it exists. An
// SLO is named by the Monitoring service it measures as well as its own id.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetServiceLevelObjectiveAttrs(t testing.TestingT, ctx context.Context, projectID string, serviceID string, sloID string) *monitoring.ServiceLevelObjective {
	slo, err := GetServiceLevelObjectiveAttrsE(t, ctx, projectID, serviceID, sloID)
	require.NoError(t, err)

	return slo
}

// GetServiceLevelObjectiveAttrsE returns the settings Google Cloud holds for the given service
// level objective.
// The ctx parameter supports cancellation and timeouts.
func GetServiceLevelObjectiveAttrsE(t testing.TestingT, ctx context.Context, projectID string, serviceID string, sloID string) (*monitoring.ServiceLevelObjective, error) {
	logger.Default.Logf(t, "Getting settings for SLO %s on service %s in project %s", sloID, serviceID, projectID)

	service, err := NewMonitoringServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetServiceLevelObjectiveAttrsWithClient(ctx, service, projectID, serviceID, sloID)
}

// GetServiceLevelObjectiveAttrsWithClient returns the settings Google Cloud holds for the given
// service level objective using the supplied *monitoring.Service. Prefer this variant in unit tests
// where the service is backed by an httptest fake server (see monitoring_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetServiceLevelObjectiveAttrsWithClient(ctx context.Context, service *monitoring.Service, projectID string, serviceID string, sloID string) (*monitoring.ServiceLevelObjective, error) {
	name := fmt.Sprintf("projects/%s/services/%s/serviceLevelObjectives/%s", projectID, serviceID, sloID)

	slo, err := service.Services.ServiceLevelObjectives.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("SLO %s does not exist on service %s in project %s", sloID, serviceID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for SLO %s on service %s in project %s: %w", sloID, serviceID, projectID, err)
	}

	return slo, nil
}

// NewMonitoringServiceE creates a Cloud Monitoring service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewMonitoringServiceE(t testing.TestingT, ctx context.Context) (*monitoring.Service, error) {
	return monitoring.NewService(ctx, append(withOptions(), option.WithScopes(monitoring.CloudPlatformScope))...)
}
