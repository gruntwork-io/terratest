package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/memcache/v1"
	"google.golang.org/api/option"
)

// GetMemcacheInstanceAttrs returns the settings Google Cloud holds for the Memorystore for Memcached instance, so a test can assert on what was
// actually created rather than only that it exists.
// The node count and the memcache version Google settled on are reported here, and a create alone
// does not say either.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMemcacheInstanceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) *memcache.Instance {
	result, err := GetMemcacheInstanceAttrsE(t, ctx, projectID, location, instanceID)
	require.NoError(t, err)

	return result
}

// GetMemcacheInstanceAttrsE returns the settings Google Cloud holds for the Memorystore for Memcached instance.
// The ctx parameter supports cancellation and timeouts.
func GetMemcacheInstanceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, instanceID string) (*memcache.Instance, error) {
	logger.Default.Logf(t, "Getting settings for Memorystore for Memcached instance %s in %s in project %s", instanceID, location, projectID)

	service, err := NewMemcacheServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMemcacheInstanceAttrsWithClient(ctx, service, projectID, location, instanceID)
}

// GetMemcacheInstanceAttrsWithClient returns the settings Google Cloud holds for the Memorystore for Memcached instance using the supplied
// *memcache.Service. Prefer this variant in unit tests where the service is backed by an httptest
// fake server (see memcache_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMemcacheInstanceAttrsWithClient(ctx context.Context, service *memcache.Service, projectID string, location string, instanceID string) (*memcache.Instance, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/instances/%s", projectID, location, instanceID)

	result, err := service.Projects.Locations.Instances.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Memorystore for Memcached instance %s in %s in project %s does not exist", instanceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Memorystore for Memcached instance %s in %s in project %s: %w", instanceID, location, projectID, err)
	}

	return result, nil
}

// NewMemcacheServiceE creates a Memorystore for Memcached service authenticated the same way every other client in this
// module is.
// The API is named memcache rather than memorystore, and it is a different service from the Redis one
// redis.go reads.
// The ctx parameter supports cancellation and timeouts.
func NewMemcacheServiceE(t testing.TestingT, ctx context.Context) (*memcache.Service, error) {
	return memcache.NewService(ctx, append(withOptions(), option.WithScopes(memcache.CloudPlatformScope))...)
}
