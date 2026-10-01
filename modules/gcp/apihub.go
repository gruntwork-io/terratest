package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/apihub/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetAPIHubInstanceAttrs returns the settings Google Cloud holds for the given API hub instance, so a test can assert on what was
// actually created rather than only that it exists.
// An instance is the catalogue itself, so its state and the key it encrypts entries with are what it is.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) *apihub.GoogleCloudApihubV1ApiHubInstance {
	attrs, err := GetAPIHubInstanceAttrsE(t, ctx, projectID, location, instanceID)
	require.NoError(t, err)

	return attrs
}

// GetAPIHubInstanceAttrsE returns the settings Google Cloud holds for the given API hub instance.
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) (*apihub.GoogleCloudApihubV1ApiHubInstance, error) {
	logger.Default.Logf(t, "Getting settings for API hub instance %s in %s in project %s", instanceID, location, projectID)

	service, err := NewAPIHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAPIHubInstanceAttrsWithClient(ctx, service, projectID, location, instanceID)
}

// GetAPIHubInstanceAttrsWithClient returns the settings Google Cloud holds for the given API hub instance using the supplied
// *apihub.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apihub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubInstanceAttrsWithClient(ctx context.Context, service *apihub.Service, projectID string, location string, instanceID string) (*apihub.GoogleCloudApihubV1ApiHubInstance, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/apiHubInstances/%s", projectID, location, instanceID)

	attrs, err := service.Projects.Locations.ApiHubInstances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the API hub instance %s in %s in project %s does not exist", instanceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for API hub instance %s in %s in project %s: %w", instanceID, location, projectID, err)
	}

	return attrs, nil
}

// GetAPIHubCurationAttrs returns the settings Google Cloud holds for the given API hub curation, so a test can assert on what was
// actually created rather than only that it exists.
// A curation is the function that normalises an API before it is catalogued, so the endpoint it calls is the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubCurationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, curationID string) *apihub.GoogleCloudApihubV1Curation {
	attrs, err := GetAPIHubCurationAttrsE(t, ctx, projectID, location, curationID)
	require.NoError(t, err)

	return attrs
}

// GetAPIHubCurationAttrsE returns the settings Google Cloud holds for the given API hub curation.
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubCurationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, curationID string) (*apihub.GoogleCloudApihubV1Curation, error) {
	logger.Default.Logf(t, "Getting settings for API hub curation %s in %s in project %s", curationID, location, projectID)

	service, err := NewAPIHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAPIHubCurationAttrsWithClient(ctx, service, projectID, location, curationID)
}

// GetAPIHubCurationAttrsWithClient returns the settings Google Cloud holds for the given API hub curation using the supplied
// *apihub.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apihub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubCurationAttrsWithClient(ctx context.Context, service *apihub.Service, projectID string, location string, curationID string) (*apihub.GoogleCloudApihubV1Curation, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/curations/%s", projectID, location, curationID)

	attrs, err := service.Projects.Locations.Curations.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the API hub curation %s in %s in project %s does not exist", curationID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for API hub curation %s in %s in project %s: %w", curationID, location, projectID, err)
	}

	return attrs, nil
}

// GetAPIHubHostProjectRegistrationAttrs returns the settings Google Cloud holds for the given API hub host project registration, so a test can assert on what was
// actually created rather than only that it exists.
// The registration is what lets one project host the catalogue for others, so the project it names is the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubHostProjectRegistrationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, registrationID string) *apihub.GoogleCloudApihubV1HostProjectRegistration {
	attrs, err := GetAPIHubHostProjectRegistrationAttrsE(t, ctx, projectID, location, registrationID)
	require.NoError(t, err)

	return attrs
}

// GetAPIHubHostProjectRegistrationAttrsE returns the settings Google Cloud holds for the given API hub host project registration.
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubHostProjectRegistrationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, registrationID string) (*apihub.GoogleCloudApihubV1HostProjectRegistration, error) {
	logger.Default.Logf(t, "Getting settings for API hub host project registration %s in %s in project %s", registrationID, location, projectID)

	service, err := NewAPIHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAPIHubHostProjectRegistrationAttrsWithClient(ctx, service, projectID, location, registrationID)
}

// GetAPIHubHostProjectRegistrationAttrsWithClient returns the settings Google Cloud holds for the given API hub host project registration using the supplied
// *apihub.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apihub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubHostProjectRegistrationAttrsWithClient(ctx context.Context, service *apihub.Service, projectID string, location string, registrationID string) (*apihub.GoogleCloudApihubV1HostProjectRegistration, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/hostProjectRegistrations/%s", projectID, location, registrationID)

	attrs, err := service.Projects.Locations.HostProjectRegistrations.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the API hub host project registration %s in %s in project %s does not exist", registrationID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for API hub host project registration %s in %s in project %s: %w", registrationID, location, projectID, err)
	}

	return attrs, nil
}

// GetAPIHubPluginAttrs returns the settings Google Cloud holds for the given API hub plugin, so a test can assert on what was
// actually created rather than only that it exists.
// A plugin is a source the catalogue can import from, so its type and state decide what it brings in.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubPluginAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, pluginID string) *apihub.GoogleCloudApihubV1Plugin {
	attrs, err := GetAPIHubPluginAttrsE(t, ctx, projectID, location, pluginID)
	require.NoError(t, err)

	return attrs
}

// GetAPIHubPluginAttrsE returns the settings Google Cloud holds for the given API hub plugin.
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubPluginAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, pluginID string) (*apihub.GoogleCloudApihubV1Plugin, error) {
	logger.Default.Logf(t, "Getting settings for API hub plugin %s in %s in project %s", pluginID, location, projectID)

	service, err := NewAPIHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAPIHubPluginAttrsWithClient(ctx, service, projectID, location, pluginID)
}

// GetAPIHubPluginAttrsWithClient returns the settings Google Cloud holds for the given API hub plugin using the supplied
// *apihub.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apihub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubPluginAttrsWithClient(ctx context.Context, service *apihub.Service, projectID string, location string, pluginID string) (*apihub.GoogleCloudApihubV1Plugin, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/plugins/%s", projectID, location, pluginID)

	attrs, err := service.Projects.Locations.Plugins.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the API hub plugin %s in %s in project %s does not exist", pluginID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for API hub plugin %s in %s in project %s: %w", pluginID, location, projectID, err)
	}

	return attrs, nil
}

// GetAPIHubPluginInstanceAttrs returns the settings Google Cloud holds for the given API hub plugin instance, so a test can assert on what was
// actually created rather than only that it exists.
// An instance is one configured copy of a plugin, so its state and its schedule decide when it actually runs.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubPluginInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, pluginID string, instanceID string) *apihub.GoogleCloudApihubV1PluginInstance {
	attrs, err := GetAPIHubPluginInstanceAttrsE(t, ctx, projectID, location, pluginID, instanceID)
	require.NoError(t, err)

	return attrs
}

// GetAPIHubPluginInstanceAttrsE returns the settings Google Cloud holds for the given API hub plugin instance.
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubPluginInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, pluginID string, instanceID string) (*apihub.GoogleCloudApihubV1PluginInstance, error) {
	logger.Default.Logf(t, "Getting settings for API hub plugin instance %s %s in %s in project %s", instanceID, pluginID, location, projectID)

	service, err := NewAPIHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAPIHubPluginInstanceAttrsWithClient(ctx, service, projectID, location, pluginID, instanceID)
}

// GetAPIHubPluginInstanceAttrsWithClient returns the settings Google Cloud holds for the given API hub plugin instance using the supplied
// *apihub.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apihub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubPluginInstanceAttrsWithClient(ctx context.Context, service *apihub.Service, projectID string, location string, pluginID string, instanceID string) (*apihub.GoogleCloudApihubV1PluginInstance, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/plugins/%s/instances/%s", projectID, location, pluginID, instanceID)

	attrs, err := service.Projects.Locations.Plugins.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the API hub plugin instance %s %s in %s in project %s does not exist", instanceID, pluginID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for API hub plugin instance %s %s in %s in project %s: %w", instanceID, pluginID, location, projectID, err)
	}

	return attrs, nil
}

// GetAPIHubRuntimeProjectAttachmentAttrs returns the settings Google Cloud holds for the given API hub runtime project attachment, so a test can assert on what was
// actually created rather than only that it exists.
// The attachment is what lets the catalogue see a runtime project's APIs, so the project it names is the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubRuntimeProjectAttachmentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, attachmentID string) *apihub.GoogleCloudApihubV1RuntimeProjectAttachment {
	attrs, err := GetAPIHubRuntimeProjectAttachmentAttrsE(t, ctx, projectID, location, attachmentID)
	require.NoError(t, err)

	return attrs
}

// GetAPIHubRuntimeProjectAttachmentAttrsE returns the settings Google Cloud holds for the given API hub runtime project attachment.
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubRuntimeProjectAttachmentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, attachmentID string) (*apihub.GoogleCloudApihubV1RuntimeProjectAttachment, error) {
	logger.Default.Logf(t, "Getting settings for API hub runtime project attachment %s in %s in project %s", attachmentID, location, projectID)

	service, err := NewAPIHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAPIHubRuntimeProjectAttachmentAttrsWithClient(ctx, service, projectID, location, attachmentID)
}

// GetAPIHubRuntimeProjectAttachmentAttrsWithClient returns the settings Google Cloud holds for the given API hub runtime project attachment using the supplied
// *apihub.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apihub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAPIHubRuntimeProjectAttachmentAttrsWithClient(ctx context.Context, service *apihub.Service, projectID string, location string, attachmentID string) (*apihub.GoogleCloudApihubV1RuntimeProjectAttachment, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/runtimeProjectAttachments/%s", projectID, location, attachmentID)

	attrs, err := service.Projects.Locations.RuntimeProjectAttachments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the API hub runtime project attachment %s in %s in project %s does not exist", attachmentID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for API hub runtime project attachment %s in %s in project %s: %w", attachmentID, location, projectID, err)
	}

	return attrs, nil
}

// NewAPIHubServiceE creates a API hub service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewAPIHubServiceE(t testing.TestingT, ctx context.Context) (*apihub.Service, error) {
	return apihub.NewService(ctx, append(withOptions(), option.WithScopes(apihub.CloudPlatformScope))...)
}
