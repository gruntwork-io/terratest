package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/securitycenter/v1"
)

// GetSecurityCenterMuteConfigAttrs returns the settings Google Cloud holds for the given Security Command Center mute config, so a test can assert on
// what was actually created rather than only that it exists.
// The filter a config mutes findings by is the whole of what the module sets, and the type tells a
// static config from a dynamic one.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSecurityCenterMuteConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, configID string) *securitycenter.GoogleCloudSecuritycenterV1MuteConfig {
	result, err := GetSecurityCenterMuteConfigAttrsE(t, ctx, projectID, configID)
	require.NoError(t, err)

	return result
}

// GetSecurityCenterMuteConfigAttrsE returns the settings Google Cloud holds for the given Security Command Center mute config.
// The ctx parameter supports cancellation and timeouts.
func GetSecurityCenterMuteConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, configID string) (*securitycenter.GoogleCloudSecuritycenterV1MuteConfig, error) {
	logger.Default.Logf(t, "Getting settings for Security Command Center mute config %s in project %s", configID, projectID)

	service, err := NewSecurityCenterServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSecurityCenterMuteConfigAttrsWithClient(ctx, service, projectID, configID)
}

// GetSecurityCenterMuteConfigAttrsWithClient returns the settings Google Cloud holds for the given Security Command Center mute config using the supplied
// *securitycenter.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see securitycenter_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSecurityCenterMuteConfigAttrsWithClient(ctx context.Context, service *securitycenter.Service, projectID string, configID string) (*securitycenter.GoogleCloudSecuritycenterV1MuteConfig, error) {
	name := fmt.Sprintf("projects/%s/muteConfigs/%s", projectID, configID)

	result, err := service.Projects.MuteConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Security Command Center mute config %s does not exist in project %s", configID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Security Command Center mute config %s in project %s: %w", configID, projectID, err)
	}

	return result, nil
}

// GetSecurityCenterNotificationConfigAttrs returns the settings Google Cloud holds for the given Security Command Center notification config, so a test can assert on
// what was actually created rather than only that it exists.
// The service account Google will publish as is filled in by the API, so it can only be asserted on
// from a read.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSecurityCenterNotificationConfigAttrs(t testing.TestingT, ctx context.Context, projectID string, configID string) *securitycenter.NotificationConfig {
	result, err := GetSecurityCenterNotificationConfigAttrsE(t, ctx, projectID, configID)
	require.NoError(t, err)

	return result
}

// GetSecurityCenterNotificationConfigAttrsE returns the settings Google Cloud holds for the given Security Command Center notification config.
// The ctx parameter supports cancellation and timeouts.
func GetSecurityCenterNotificationConfigAttrsE(t testing.TestingT, ctx context.Context, projectID string, configID string) (*securitycenter.NotificationConfig, error) {
	logger.Default.Logf(t, "Getting settings for Security Command Center notification config %s in project %s", configID, projectID)

	service, err := NewSecurityCenterServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSecurityCenterNotificationConfigAttrsWithClient(ctx, service, projectID, configID)
}

// GetSecurityCenterNotificationConfigAttrsWithClient returns the settings Google Cloud holds for the given Security Command Center notification config using the supplied
// *securitycenter.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see securitycenter_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSecurityCenterNotificationConfigAttrsWithClient(ctx context.Context, service *securitycenter.Service, projectID string, configID string) (*securitycenter.NotificationConfig, error) {
	name := fmt.Sprintf("projects/%s/notificationConfigs/%s", projectID, configID)

	result, err := service.Projects.NotificationConfigs.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Security Command Center notification config %s does not exist in project %s", configID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Security Command Center notification config %s in project %s: %w", configID, projectID, err)
	}

	return result, nil
}

// GetSecurityCenterSecurityHealthAnalyticsCustomModuleAttrs returns the settings Google Cloud holds for the given Security Health Analytics custom module, so a test can assert on
// what was actually created rather than only that it exists.
// The ancestor module and the last editor Google records are reported here and nowhere in the create
// request.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSecurityCenterSecurityHealthAnalyticsCustomModuleAttrs(t testing.TestingT, ctx context.Context, projectID string, moduleID string) *securitycenter.GoogleCloudSecuritycenterV1SecurityHealthAnalyticsCustomModule {
	result, err := GetSecurityCenterSecurityHealthAnalyticsCustomModuleAttrsE(t, ctx, projectID, moduleID)
	require.NoError(t, err)

	return result
}

// GetSecurityCenterSecurityHealthAnalyticsCustomModuleAttrsE returns the settings Google Cloud holds for the given Security Health Analytics custom module.
// The ctx parameter supports cancellation and timeouts.
func GetSecurityCenterSecurityHealthAnalyticsCustomModuleAttrsE(t testing.TestingT, ctx context.Context, projectID string, moduleID string) (*securitycenter.GoogleCloudSecuritycenterV1SecurityHealthAnalyticsCustomModule, error) {
	logger.Default.Logf(t, "Getting settings for Security Health Analytics custom module %s in project %s", moduleID, projectID)

	service, err := NewSecurityCenterServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSecurityCenterSecurityHealthAnalyticsCustomModuleAttrsWithClient(ctx, service, projectID, moduleID)
}

// GetSecurityCenterSecurityHealthAnalyticsCustomModuleAttrsWithClient returns the settings Google Cloud holds for the given Security Health Analytics custom module using the supplied
// *securitycenter.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see securitycenter_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSecurityCenterSecurityHealthAnalyticsCustomModuleAttrsWithClient(ctx context.Context, service *securitycenter.Service, projectID string, moduleID string) (*securitycenter.GoogleCloudSecuritycenterV1SecurityHealthAnalyticsCustomModule, error) {
	name := fmt.Sprintf("projects/%s/securityHealthAnalyticsSettings/customModules/%s", projectID, moduleID)

	result, err := service.Projects.SecurityHealthAnalyticsSettings.CustomModules.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Security Health Analytics custom module %s does not exist in project %s", moduleID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Security Health Analytics custom module %s in project %s: %w", moduleID, projectID, err)
	}

	return result, nil
}

// GetSecurityCenterEventThreatDetectionCustomModuleAttrs returns the settings Google Cloud holds for the given Event Threat Detection custom module, so a test can assert on
// what was actually created rather than only that it exists.
// The config a module is created with comes back as the arbitrary JSON the detector reads, so a test
// asserts on the keys it set.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSecurityCenterEventThreatDetectionCustomModuleAttrs(t testing.TestingT, ctx context.Context, projectID string, moduleID string) *securitycenter.EventThreatDetectionCustomModule {
	result, err := GetSecurityCenterEventThreatDetectionCustomModuleAttrsE(t, ctx, projectID, moduleID)
	require.NoError(t, err)

	return result
}

// GetSecurityCenterEventThreatDetectionCustomModuleAttrsE returns the settings Google Cloud holds for the given Event Threat Detection custom module.
// The ctx parameter supports cancellation and timeouts.
func GetSecurityCenterEventThreatDetectionCustomModuleAttrsE(t testing.TestingT, ctx context.Context, projectID string, moduleID string) (*securitycenter.EventThreatDetectionCustomModule, error) {
	logger.Default.Logf(t, "Getting settings for Event Threat Detection custom module %s in project %s", moduleID, projectID)

	service, err := NewSecurityCenterServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetSecurityCenterEventThreatDetectionCustomModuleAttrsWithClient(ctx, service, projectID, moduleID)
}

// GetSecurityCenterEventThreatDetectionCustomModuleAttrsWithClient returns the settings Google Cloud holds for the given Event Threat Detection custom module using the supplied
// *securitycenter.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see securitycenter_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetSecurityCenterEventThreatDetectionCustomModuleAttrsWithClient(ctx context.Context, service *securitycenter.Service, projectID string, moduleID string) (*securitycenter.EventThreatDetectionCustomModule, error) {
	name := fmt.Sprintf("projects/%s/eventThreatDetectionSettings/customModules/%s", projectID, moduleID)

	result, err := service.Projects.EventThreatDetectionSettings.CustomModules.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Event Threat Detection custom module %s does not exist in project %s", moduleID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Event Threat Detection custom module %s in project %s: %w", moduleID, projectID, err)
	}

	return result, nil
}

// NewSecurityCenterServiceE creates a Security Command Center service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewSecurityCenterServiceE(t testing.TestingT, ctx context.Context) (*securitycenter.Service, error) {
	return securitycenter.NewService(ctx, append(withOptions(), option.WithScopes(securitycenter.CloudPlatformScope))...)
}
