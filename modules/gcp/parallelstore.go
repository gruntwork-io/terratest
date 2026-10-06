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
	"google.golang.org/api/parallelstore/v1"
)

// GetParallelstoreInstanceAttrs returns the settings Google Cloud holds for the given Parallelstore instance, so a test can assert on what was
// actually created rather than only that it exists.
// An instance is the scratch filesystem itself, so its capacity, its tier and the network it is reachable on are what a job gets.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetParallelstoreInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) *parallelstore.Instance {
	attrs, err := GetParallelstoreInstanceAttrsE(t, ctx, projectID, location, instanceID)
	require.NoError(t, err)

	return attrs
}

// GetParallelstoreInstanceAttrsE returns the settings Google Cloud holds for the given Parallelstore instance.
// The ctx parameter supports cancellation and timeouts.
func GetParallelstoreInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) (*parallelstore.Instance, error) {
	logger.Default.Logf(t, "Getting settings for Parallelstore instance %s in %s in project %s", instanceID, location, projectID)

	service, err := NewParallelstoreServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetParallelstoreInstanceAttrsWithClient(ctx, service, projectID, location, instanceID)
}

// GetParallelstoreInstanceAttrsWithClient returns the settings Google Cloud holds for the given Parallelstore instance using the supplied
// *parallelstore.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see parallelstore_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetParallelstoreInstanceAttrsWithClient(ctx context.Context, service *parallelstore.Service, projectID string, location string, instanceID string) (*parallelstore.Instance, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/instances/%s", projectID, location, instanceID)

	attrs, err := service.Projects.Locations.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Parallelstore instance %s in %s in project %s does not exist", instanceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Parallelstore instance %s in %s in project %s: %w", instanceID, location, projectID, err)
	}

	return attrs, nil
}

// NewParallelstoreServiceE creates a Parallelstore service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewParallelstoreServiceE(t testing.TestingT, ctx context.Context) (*parallelstore.Service, error) {
	return parallelstore.NewService(ctx, append(withOptions(), option.WithScopes(parallelstore.CloudPlatformScope))...)
}
