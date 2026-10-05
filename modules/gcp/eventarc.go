package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/eventarc/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetEventarcTriggerAttrs returns the settings Google Cloud holds for the given Eventarc trigger,
// so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcTriggerAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, triggerID string) *eventarc.Trigger {
	trigger, err := GetEventarcTriggerAttrsE(t, ctx, projectID, region, triggerID)
	require.NoError(t, err)

	return trigger
}

// GetEventarcTriggerAttrsE returns the settings Google Cloud holds for the given Eventarc trigger.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcTriggerAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, triggerID string) (*eventarc.Trigger, error) {
	logger.Default.Logf(t, "Getting settings for Eventarc trigger %s in region %s in project %s", triggerID, region, projectID)

	service, err := NewEventarcServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetEventarcTriggerAttrsWithClient(ctx, service, projectID, region, triggerID)
}

// GetEventarcTriggerAttrsWithClient returns the settings Google Cloud holds for the given Eventarc
// trigger using the supplied *eventarc.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see eventarc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetEventarcTriggerAttrsWithClient(ctx context.Context, service *eventarc.Service, projectID string, region string, triggerID string) (*eventarc.Trigger, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/triggers/%s", projectID, region, triggerID)

	trigger, err := service.Projects.Locations.Triggers.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Eventarc trigger %s does not exist in region %s in project %s", triggerID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Eventarc trigger %s in region %s in project %s: %w", triggerID, region, projectID, err)
	}

	return trigger, nil
}

// GetEventarcEnrollmentAttrs returns the settings Google Cloud holds for the given Eventarc enrollment, so a test can assert on what was
// actually created rather than only that it exists.
// An enrollment is the subscription that decides which messages on a bus reach one pipeline, so its filter is the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcEnrollmentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *eventarc.Enrollment {
	attrs, err := GetEventarcEnrollmentAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetEventarcEnrollmentAttrsE returns the settings Google Cloud holds for the given Eventarc enrollment.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcEnrollmentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*eventarc.Enrollment, error) {
	logger.Default.Logf(t, "Getting settings for Eventarc enrollment %s in %s in project %s", id, location, projectID)

	service, err := NewEventarcServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetEventarcEnrollmentAttrsWithClient(ctx, service, projectID, location, id)
}

// GetEventarcEnrollmentAttrsWithClient returns the settings Google Cloud holds for the given Eventarc enrollment using the supplied
// *eventarc.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see eventarc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetEventarcEnrollmentAttrsWithClient(ctx context.Context, service *eventarc.Service, projectID string, location string, id string) (*eventarc.Enrollment, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/enrollments/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.Enrollments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Eventarc enrollment %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Eventarc enrollment %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetEventarcGoogleAPISourceAttrs returns the settings Google Cloud holds for the given Eventarc Google API source, so a test can assert on what was
// actually created rather than only that it exists.
// This is what puts a project's own Google API events onto a bus, so the bus it feeds is the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcGoogleAPISourceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *eventarc.GoogleApiSource {
	attrs, err := GetEventarcGoogleAPISourceAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetEventarcGoogleAPISourceAttrsE returns the settings Google Cloud holds for the given Eventarc Google API source.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcGoogleAPISourceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*eventarc.GoogleApiSource, error) {
	logger.Default.Logf(t, "Getting settings for Eventarc Google API source %s in %s in project %s", id, location, projectID)

	service, err := NewEventarcServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetEventarcGoogleAPISourceAttrsWithClient(ctx, service, projectID, location, id)
}

// GetEventarcGoogleAPISourceAttrsWithClient returns the settings Google Cloud holds for the given Eventarc Google API source using the supplied
// *eventarc.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see eventarc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetEventarcGoogleAPISourceAttrsWithClient(ctx context.Context, service *eventarc.Service, projectID string, location string, id string) (*eventarc.GoogleApiSource, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/googleApiSources/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.GoogleApiSources.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Eventarc Google API source %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Eventarc Google API source %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetEventarcMessageBusAttrs returns the settings Google Cloud holds for the given Eventarc message bus, so a test can assert on what was
// actually created rather than only that it exists.
// A bus is where publishers put events and enrollments read them, so its own existence and encryption are what it offers.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcMessageBusAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *eventarc.MessageBus {
	attrs, err := GetEventarcMessageBusAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetEventarcMessageBusAttrsE returns the settings Google Cloud holds for the given Eventarc message bus.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcMessageBusAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*eventarc.MessageBus, error) {
	logger.Default.Logf(t, "Getting settings for Eventarc message bus %s in %s in project %s", id, location, projectID)

	service, err := NewEventarcServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetEventarcMessageBusAttrsWithClient(ctx, service, projectID, location, id)
}

// GetEventarcMessageBusAttrsWithClient returns the settings Google Cloud holds for the given Eventarc message bus using the supplied
// *eventarc.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see eventarc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetEventarcMessageBusAttrsWithClient(ctx context.Context, service *eventarc.Service, projectID string, location string, id string) (*eventarc.MessageBus, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/messageBuses/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.MessageBuses.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Eventarc message bus %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Eventarc message bus %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// GetEventarcChannelAttrs returns the settings Google Cloud holds for the given Eventarc channel, so a test can assert on
// what was actually created rather than only that it exists. A channel is how events from a third party provider reach
// this project.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcChannelAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, channelID string) *eventarc.Channel {
	channel, err := GetEventarcChannelAttrsE(t, ctx, projectID, location, channelID)
	require.NoError(t, err)

	return channel
}

// GetEventarcChannelAttrsE returns the settings Google Cloud holds for the given Eventarc channel.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcChannelAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, channelID string) (*eventarc.Channel, error) {
	logger.Default.Logf(t, "Getting settings for Eventarc channel %s in %s in project %s", channelID, location, projectID)

	service, err := NewEventarcServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetEventarcChannelAttrsWithClient(ctx, service, projectID, location, channelID)
}

// GetEventarcChannelAttrsWithClient returns the settings Google Cloud holds for the given Eventarc channel using the supplied
// *eventarc.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see eventarc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetEventarcChannelAttrsWithClient(ctx context.Context, service *eventarc.Service, projectID string, location string, channelID string) (*eventarc.Channel, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/channels/%s", projectID, location, channelID)

	channel, err := service.Projects.Locations.Channels.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Eventarc channel %s does not exist in %s in project %s", channelID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Eventarc channel %s in %s in project %s: %w", channelID, location, projectID, err)
	}

	return channel, nil
}

// GetEventarcPipelineAttrs returns the settings Google Cloud holds for the given pipeline, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcPipelineAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, pipelineID string) *eventarc.Pipeline {
	pipeline, err := GetEventarcPipelineAttrsE(t, ctx, projectID, location, pipelineID)
	require.NoError(t, err)

	return pipeline
}

// GetEventarcPipelineAttrsE returns the settings Google Cloud holds for the given pipeline.
// The ctx parameter supports cancellation and timeouts.
func GetEventarcPipelineAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, pipelineID string) (*eventarc.Pipeline, error) {
	logger.Default.Logf(t, "Getting settings for pipeline %s in %s in project %s", pipelineID, location, projectID)

	service, err := NewEventarcServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetEventarcPipelineAttrsWithClient(ctx, service, projectID, location, pipelineID)
}

// GetEventarcPipelineAttrsWithClient returns the settings Google Cloud holds for the given pipeline using the supplied
// *eventarc.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see eventarc_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetEventarcPipelineAttrsWithClient(ctx context.Context, service *eventarc.Service, projectID string, location string, pipelineID string) (*eventarc.Pipeline, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/pipelines/%s", projectID, location, pipelineID)

	pipeline, err := service.Projects.Locations.Pipelines.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the pipeline %s does not exist in %s in project %s", pipelineID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for pipeline %s in %s in project %s: %w", pipelineID, location, projectID, err)
	}

	return pipeline, nil
}

// NewEventarcServiceE creates a Eventarc service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewEventarcServiceE(t testing.TestingT, ctx context.Context) (*eventarc.Service, error) {
	return eventarc.NewService(ctx, append(withOptions(), option.WithScopes(eventarc.CloudPlatformScope))...)
}
