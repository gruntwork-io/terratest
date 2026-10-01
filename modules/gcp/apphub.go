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

// NewAppHubServiceE creates an App Hub service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewAppHubServiceE(t testing.TestingT, ctx context.Context) (*apphub.APIService, error) {
	return apphub.NewService(ctx, append(withOptions(), option.WithScopes(apphub.CloudPlatformScope))...)
}
