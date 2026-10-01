package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/firebaseapphosting/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetFirebaseAppHostingBackendAttrs returns the settings Google Cloud holds for the given Firebase App Hosting backend, so a test can assert on what was
// actually created rather than only that it exists.
// A backend is the deployed web app, so the service account it runs as and the URI it answers on are the point.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppHostingBackendAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, backendID string) *firebaseapphosting.Backend {
	attrs, err := GetFirebaseAppHostingBackendAttrsE(t, ctx, projectID, location, backendID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseAppHostingBackendAttrsE returns the settings Google Cloud holds for the given Firebase App Hosting backend.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppHostingBackendAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, backendID string) (*firebaseapphosting.Backend, error) {
	logger.Default.Logf(t, "Getting settings for Firebase App Hosting backend %s in %s in project %s", backendID, location, projectID)

	service, err := NewFirebaseAppHostingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseAppHostingBackendAttrsWithClient(ctx, service, projectID, location, backendID)
}

// GetFirebaseAppHostingBackendAttrsWithClient returns the settings Google Cloud holds for the given Firebase App Hosting backend using the supplied
// *firebaseapphosting.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebaseapphosting_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppHostingBackendAttrsWithClient(ctx context.Context, service *firebaseapphosting.Service, projectID string, location string, backendID string) (*firebaseapphosting.Backend, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backends/%s", projectID, location, backendID)

	attrs, err := service.Projects.Locations.Backends.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase App Hosting backend %s in %s in project %s does not exist", backendID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase App Hosting backend %s in %s in project %s: %w", backendID, location, projectID, err)
	}

	return attrs, nil
}

// GetFirebaseAppHostingBuildAttrs returns the settings Google Cloud holds for the given Firebase App Hosting build, so a test can assert on what was
// actually created rather than only that it exists.
// A build is one attempt to turn a commit into something servable, so its source and its state are what it records.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppHostingBuildAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, backendID string, buildID string) *firebaseapphosting.Build {
	attrs, err := GetFirebaseAppHostingBuildAttrsE(t, ctx, projectID, location, backendID, buildID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseAppHostingBuildAttrsE returns the settings Google Cloud holds for the given Firebase App Hosting build.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppHostingBuildAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, backendID string, buildID string) (*firebaseapphosting.Build, error) {
	logger.Default.Logf(t, "Getting settings for Firebase App Hosting build %s %s in %s in project %s", buildID, backendID, location, projectID)

	service, err := NewFirebaseAppHostingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseAppHostingBuildAttrsWithClient(ctx, service, projectID, location, backendID, buildID)
}

// GetFirebaseAppHostingBuildAttrsWithClient returns the settings Google Cloud holds for the given Firebase App Hosting build using the supplied
// *firebaseapphosting.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebaseapphosting_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppHostingBuildAttrsWithClient(ctx context.Context, service *firebaseapphosting.Service, projectID string, location string, backendID string, buildID string) (*firebaseapphosting.Build, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backends/%s/builds/%s", projectID, location, backendID, buildID)

	attrs, err := service.Projects.Locations.Backends.Builds.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase App Hosting build %s %s in %s in project %s does not exist", buildID, backendID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase App Hosting build %s %s in %s in project %s: %w", buildID, backendID, location, projectID, err)
	}

	return attrs, nil
}

// GetFirebaseAppHostingDomainAttrs returns the settings Google Cloud holds for the given Firebase App Hosting domain, so a test can assert on what was
// actually created rather than only that it exists.
// A domain is a hostname the backend answers on, so its serving state says whether it actually does.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppHostingDomainAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, backendID string, domainID string) *firebaseapphosting.Domain {
	attrs, err := GetFirebaseAppHostingDomainAttrsE(t, ctx, projectID, location, backendID, domainID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseAppHostingDomainAttrsE returns the settings Google Cloud holds for the given Firebase App Hosting domain.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppHostingDomainAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, backendID string, domainID string) (*firebaseapphosting.Domain, error) {
	logger.Default.Logf(t, "Getting settings for Firebase App Hosting domain %s %s in %s in project %s", domainID, backendID, location, projectID)

	service, err := NewFirebaseAppHostingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseAppHostingDomainAttrsWithClient(ctx, service, projectID, location, backendID, domainID)
}

// GetFirebaseAppHostingDomainAttrsWithClient returns the settings Google Cloud holds for the given Firebase App Hosting domain using the supplied
// *firebaseapphosting.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebaseapphosting_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppHostingDomainAttrsWithClient(ctx context.Context, service *firebaseapphosting.Service, projectID string, location string, backendID string, domainID string) (*firebaseapphosting.Domain, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backends/%s/domains/%s", projectID, location, backendID, domainID)

	attrs, err := service.Projects.Locations.Backends.Domains.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase App Hosting domain %s %s in %s in project %s does not exist", domainID, backendID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase App Hosting domain %s %s in %s in project %s: %w", domainID, backendID, location, projectID, err)
	}

	return attrs, nil
}

// GetFirebaseAppHostingTrafficAttrs returns the settings Google Cloud holds for the given Firebase App Hosting traffic, so a test can assert on what was
// actually created rather than only that it exists.
// Traffic is what decides which build serves requests, so the target it names is the live version.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppHostingTrafficAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, backendID string) *firebaseapphosting.Traffic {
	attrs, err := GetFirebaseAppHostingTrafficAttrsE(t, ctx, projectID, location, backendID)
	require.NoError(t, err)

	return attrs
}

// GetFirebaseAppHostingTrafficAttrsE returns the settings Google Cloud holds for the given Firebase App Hosting traffic.
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppHostingTrafficAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, backendID string) (*firebaseapphosting.Traffic, error) {
	logger.Default.Logf(t, "Getting settings for Firebase App Hosting traffic %s in %s in project %s", backendID, location, projectID)

	service, err := NewFirebaseAppHostingServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetFirebaseAppHostingTrafficAttrsWithClient(ctx, service, projectID, location, backendID)
}

// GetFirebaseAppHostingTrafficAttrsWithClient returns the settings Google Cloud holds for the given Firebase App Hosting traffic using the supplied
// *firebaseapphosting.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see firebaseapphosting_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetFirebaseAppHostingTrafficAttrsWithClient(ctx context.Context, service *firebaseapphosting.Service, projectID string, location string, backendID string) (*firebaseapphosting.Traffic, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/backends/%s/traffic", projectID, location, backendID)

	attrs, err := service.Projects.Locations.Backends.Traffic.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Firebase App Hosting traffic %s in %s in project %s does not exist", backendID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Firebase App Hosting traffic %s in %s in project %s: %w", backendID, location, projectID, err)
	}

	return attrs, nil
}

// NewFirebaseAppHostingServiceE creates a Firebase App Hosting service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewFirebaseAppHostingServiceE(t testing.TestingT, ctx context.Context) (*firebaseapphosting.Service, error) {
	return firebaseapphosting.NewService(ctx, append(withOptions(), option.WithScopes(firebaseapphosting.CloudPlatformScope))...)
}
