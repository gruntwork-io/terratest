package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/firebasedataconnect/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetFirebaseDataConnectServiceAttrs returns the settings Google Cloud holds for the given Firebase Data Connect service, so a test can assert on what was
// actually created rather than only that it exists.
// A service is the GraphQL endpoint in front of a database, so the display name and state are what say it is there.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseDataConnectServiceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string) *firebasedataconnect.Service {
	attrs, err := GetFirebaseDataConnectServiceAttrsE(t, ctx, projectID, location, serviceID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseDataConnectServiceAttrsE returns the settings Google Cloud holds for the given Firebase Data Connect service.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseDataConnectServiceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string) (*firebasedataconnect.Service, error) {
	logger.Default.Logf(t, "Getting settings for Firebase Data Connect service %s in %s in project %s", serviceID, location, projectID)

	service, err := NewFirebaseDataConnectServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseDataConnectServiceAttrsWithClient(ctx, service, projectID, location, serviceID)
}

// GetFirebaseDataConnectServiceAttrsWithClient returns the settings Google Cloud holds for the given Firebase Data Connect service using the supplied
// *firebasedataconnect.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebasedataconnect_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseDataConnectServiceAttrsWithClient(ctx context.Context, service *firebasedataconnect.APIService, projectID string, location string, serviceID string) (*firebasedataconnect.Service, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/services/%s", projectID, location, serviceID)

	attrs, err := service.Projects.Locations.Services.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase Data Connect service %s in %s in project %s does not exist", serviceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase Data Connect service %s in %s in project %s: %w", serviceID, location, projectID, err)
	}

	return attrs, nil
}

// NewFirebaseDataConnectServiceE creates a Firebase Data Connect service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewFirebaseDataConnectServiceE(t testing.TestingT, ctx context.Context) (*firebasedataconnect.APIService, error) {
	return firebasedataconnect.NewService(ctx, append(withOptions(), option.WithScopes(firebasedataconnect.CloudPlatformScope))...)
}
