package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/apphub/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetAppHubApplicationAttrs returns the settings Google Cloud holds for the given App Hub
// application, so a test can assert on the scope and attributes it was given rather than only that
// it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAppHubApplicationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, applicationID string) *apphub.Application {
	application, err := GetAppHubApplicationAttrsE(t, ctx, projectID, location, applicationID)
	require.NoError(t, err)

	return application
}

// GetAppHubApplicationAttrsE returns the settings Google Cloud holds for the given App Hub
// application.
// The ctx parameter supports cancellation and timeouts.
func GetAppHubApplicationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, applicationID string) (*apphub.Application, error) {
	logger.Default.Logf(t, "Getting settings for App Hub application %s in %s in project %s", applicationID, location, projectID)

	service, err := NewAppHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAppHubApplicationAttrsWithClient(ctx, service, projectID, location, applicationID)
}

// GetAppHubApplicationAttrsWithClient returns the settings Google Cloud holds for the given App Hub
// application using the supplied *apphub.APIService. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see apphub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAppHubApplicationAttrsWithClient(ctx context.Context, service *apphub.APIService, projectID string, location string, applicationID string) (*apphub.Application, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/applications/%s", projectID, location, applicationID)

	application, err := service.Projects.Locations.Applications.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the App Hub application %s in %s in project %s does not exist", applicationID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for App Hub application %s in %s in project %s: %w", applicationID, location, projectID, err)
	}

	return application, nil
}

// GetAppHubServiceAttrs returns the settings Google Cloud holds for the given App Hub service, so a test can assert on what was
// actually created rather than only that it exists.
// A service is one piece of a registered application, so the discovered service behind it and its criticality are what the catalogue is for.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAppHubServiceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, applicationID string, id string) *apphub.Service {
	attrs, err := GetAppHubServiceAttrsE(t, ctx, projectID, location, applicationID, id)
	require.NoError(t, err)

	return attrs
}

// GetAppHubServiceAttrsE returns the settings Google Cloud holds for the given App Hub service.
// The ctx parameter supports cancellation and timeouts.
func GetAppHubServiceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, applicationID string, id string) (*apphub.Service, error) {
	logger.Default.Logf(t, "Getting settings for App Hub service %s in application %s in %s in project %s", id, applicationID, location, projectID)

	service, err := NewAppHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAppHubServiceAttrsWithClient(ctx, service, projectID, location, applicationID, id)
}

// GetAppHubServiceAttrsWithClient returns the settings Google Cloud holds for the given App Hub service using the supplied
// *apphub.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apphub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAppHubServiceAttrsWithClient(ctx context.Context, service *apphub.APIService, projectID string, location string, applicationID string, id string) (*apphub.Service, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/applications/%s/services/%s", projectID, location, applicationID, id)

	attrs, err := service.Projects.Locations.Applications.Services.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the App Hub service %s in application %s in %s in project %s does not exist", id, applicationID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for App Hub service %s in application %s in %s in project %s: %w", id, applicationID, location, projectID, err)
	}

	return attrs, nil
}

// GetAppHubWorkloadAttrs returns the settings Google Cloud holds for the given App Hub workload, so a test can assert on what was
// actually created rather than only that it exists.
// A workload is the running thing a registered application is made of, so what it points at and its state are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAppHubWorkloadAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, applicationID string, id string) *apphub.Workload {
	attrs, err := GetAppHubWorkloadAttrsE(t, ctx, projectID, location, applicationID, id)
	require.NoError(t, err)

	return attrs
}

// GetAppHubWorkloadAttrsE returns the settings Google Cloud holds for the given App Hub workload.
// The ctx parameter supports cancellation and timeouts.
func GetAppHubWorkloadAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, applicationID string, id string) (*apphub.Workload, error) {
	logger.Default.Logf(t, "Getting settings for App Hub workload %s in application %s in %s in project %s", id, applicationID, location, projectID)

	service, err := NewAppHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAppHubWorkloadAttrsWithClient(ctx, service, projectID, location, applicationID, id)
}

// GetAppHubWorkloadAttrsWithClient returns the settings Google Cloud holds for the given App Hub workload using the supplied
// *apphub.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apphub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAppHubWorkloadAttrsWithClient(ctx context.Context, service *apphub.APIService, projectID string, location string, applicationID string, id string) (*apphub.Workload, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/applications/%s/workloads/%s", projectID, location, applicationID, id)

	attrs, err := service.Projects.Locations.Applications.Workloads.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the App Hub workload %s in application %s in %s in project %s does not exist", id, applicationID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for App Hub workload %s in application %s in %s in project %s: %w", id, applicationID, location, projectID, err)
	}

	return attrs, nil
}

// GetAppHubServiceProjectAttachmentAttrs returns the settings Google Cloud holds for the given App Hub service project attachment, so a test can assert on what was
// actually created rather than only that it exists.
// The attachment is what lets a host project see another project's resources, so the project it names is the whole point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAppHubServiceProjectAttachmentAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, id string) *apphub.ServiceProjectAttachment {
	attrs, err := GetAppHubServiceProjectAttachmentAttrsE(t, ctx, projectID, location, id)
	require.NoError(t, err)

	return attrs
}

// GetAppHubServiceProjectAttachmentAttrsE returns the settings Google Cloud holds for the given App Hub service project attachment.
// The ctx parameter supports cancellation and timeouts.
func GetAppHubServiceProjectAttachmentAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, id string) (*apphub.ServiceProjectAttachment, error) {
	logger.Default.Logf(t, "Getting settings for App Hub service project attachment %s in %s in project %s", id, location, projectID)

	service, err := NewAppHubServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAppHubServiceProjectAttachmentAttrsWithClient(ctx, service, projectID, location, id)
}

// GetAppHubServiceProjectAttachmentAttrsWithClient returns the settings Google Cloud holds for the given App Hub service project attachment using the supplied
// *apphub.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see apphub_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAppHubServiceProjectAttachmentAttrsWithClient(ctx context.Context, service *apphub.APIService, projectID string, location string, id string) (*apphub.ServiceProjectAttachment, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/serviceProjectAttachments/%s", projectID, location, id)

	attrs, err := service.Projects.Locations.ServiceProjectAttachments.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the App Hub service project attachment %s in %s in project %s does not exist", id, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for App Hub service project attachment %s in %s in project %s: %w", id, location, projectID, err)
	}

	return attrs, nil
}

// NewAppHubServiceE creates an App Hub service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewAppHubServiceE(t testing.TestingT, ctx context.Context) (*apphub.APIService, error) {
	return apphub.NewService(ctx, append(withOptions(), option.WithScopes(apphub.CloudPlatformScope))...)
}
