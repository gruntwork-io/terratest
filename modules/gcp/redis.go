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
	"google.golang.org/api/redis/v1"
)

// GetRedisInstanceAttrs returns the settings Google Cloud holds for the given Memorystore for Redis
// instance, so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRedisInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, region string, instanceID string) *redis.Instance {
	instance, err := GetRedisInstanceAttrsE(t, ctx, projectID, region, instanceID)
	require.NoError(t, err)

	return instance
}

// GetRedisInstanceAttrsE returns the settings Google Cloud holds for the given Memorystore for
// Redis instance.
// The ctx parameter supports cancellation and timeouts.
func GetRedisInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, region string, instanceID string) (*redis.Instance, error) {
	logger.Default.Logf(t, "Getting settings for Memorystore for Redis instance %s in region %s in project %s", instanceID, region, projectID)

	service, err := NewRedisServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetRedisInstanceAttrsWithClient(ctx, service, projectID, region, instanceID)
}

// GetRedisInstanceAttrsWithClient returns the settings Google Cloud holds for the given Memorystore
// for Redis instance using the supplied *redis.Service. Prefer this variant in unit tests where the
// service is backed by an httptest fake server (see redis_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetRedisInstanceAttrsWithClient(ctx context.Context, service *redis.Service, projectID string, region string, instanceID string) (*redis.Instance, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/instances/%s", projectID, region, instanceID)

	instance, err := service.Projects.Locations.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Memorystore for Redis instance %s does not exist in region %s in project %s", instanceID, region, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Memorystore for Redis instance %s in region %s in project %s: %w", instanceID, region, projectID, err)
	}

	return instance, nil
}

// NewRedisServiceE creates a Memorystore for Redis service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewRedisServiceE(t testing.TestingT, ctx context.Context) (*redis.Service, error) {
	return redis.NewService(ctx, append(withOptions(), option.WithScopes(redis.CloudPlatformScope))...)
}
