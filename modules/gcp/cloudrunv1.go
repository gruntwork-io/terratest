package gcp

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	runv1 "google.golang.org/api/run/v1"
)

// GetCloudRunV1ServiceAttrs returns the settings Google Cloud holds for the given Cloud Run service
// as the v1 API describes it, so a test can assert on what the older module created. The v1 API is a
// separate surface from the v2 one the rest of this package reads, and it answers in Knative's shape
// rather than Google's.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunV1ServiceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string) *runv1.Service {
	service, err := GetCloudRunV1ServiceAttrsE(t, ctx, projectID, location, serviceID)
	require.NoError(t, err)

	return service
}

// GetCloudRunV1ServiceAttrsE returns the settings Google Cloud holds for the given Cloud Run service
// as the v1 API describes it.
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunV1ServiceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string) (*runv1.Service, error) {
	logger.Default.Logf(t, "Getting v1 settings for Cloud Run service %s in %s in project %s", serviceID, location, projectID)

	// The v1 API is served per region rather than globally, and a call to the global endpoint finds
	// nothing, so the client is pointed at the region the service lives in.
	client, err := NewCloudRunV1ServiceE(t, ctx, location)
	if err != nil {
		return nil, err
	}

	return GetCloudRunV1ServiceAttrsWithClient(ctx, client, projectID, location, serviceID)
}

// GetCloudRunV1ServiceAttrsWithClient returns the settings Google Cloud holds for the given Cloud
// Run service using the supplied *runv1.APIService. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see cloudrunv1_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetCloudRunV1ServiceAttrsWithClient(ctx context.Context, client *runv1.APIService, projectID string, location string, serviceID string) (*runv1.Service, error) {
	// The v1 API names a service by Knative namespace, and the namespace is the project.
	name := fmt.Sprintf("namespaces/%s/services/%s", projectID, serviceID)

	service, err := client.Namespaces.Services.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Cloud Run service %s in %s in project %s does not exist", serviceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get v1 settings for Cloud Run service %s in %s in project %s: %w", serviceID, location, projectID, err)
	}

	return service, nil
}

// NewCloudRunV1ServiceE creates a Cloud Run v1 client pointed at one region, authenticated the same
// way every other client in this module is. The v1 API serves each region from its own endpoint.
// The ctx parameter supports cancellation and timeouts.
func NewCloudRunV1ServiceE(t testing.TestingT, ctx context.Context, location string) (*runv1.APIService, error) {
	if !cloudRunV1LocationPattern.MatchString(location) {
		return nil, fmt.Errorf("%q is not a valid location: a location may hold only lowercase letters, digits and hyphens", location)
	}

	endpoint := fmt.Sprintf("https://%s-run.googleapis.com/", location)

	return runv1.NewService(ctx, append(withOptions(),
		option.WithScopes(runv1.CloudPlatformScope), option.WithEndpoint(endpoint))...)
}

// cloudRunV1LocationPattern is what a Cloud Run location may look like. The location goes into the
// endpoint host, so anything else would send the request, and the caller's credentials with it,
// somewhere the caller did not name.
var cloudRunV1LocationPattern = regexp.MustCompile(`^[a-z0-9-]+$`)
