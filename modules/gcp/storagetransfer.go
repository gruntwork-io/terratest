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
	"google.golang.org/api/storagetransfer/v1"
)

// GetAgentPoolAttrs returns the settings Google Cloud holds for the given Storage Transfer agent
// pool, so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAgentPoolAttrs(t testing.TestingT, ctx context.Context, projectID string, poolName string) *storagetransfer.AgentPool {
	pool, err := GetAgentPoolAttrsE(t, ctx, projectID, poolName)
	require.NoError(t, err)

	return pool
}

// GetAgentPoolAttrsE returns the settings Google Cloud holds for the given Storage Transfer agent
// pool.
// The ctx parameter supports cancellation and timeouts.
func GetAgentPoolAttrsE(t testing.TestingT, ctx context.Context, projectID string, poolName string) (*storagetransfer.AgentPool, error) {
	logger.Default.Logf(t, "Getting settings for agent pool %s in project %s", poolName, projectID)

	service, err := NewStorageTransferServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAgentPoolAttrsWithClient(ctx, service, projectID, poolName)
}

// GetAgentPoolAttrsWithClient returns the settings Google Cloud holds for the given Storage
// Transfer agent pool using the supplied *storagetransfer.Service. Prefer this variant in unit
// tests where the service is backed by an httptest fake server (see storagetransfer_test.go for
// the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAgentPoolAttrsWithClient(ctx context.Context, service *storagetransfer.Service, projectID string, poolName string) (*storagetransfer.AgentPool, error) {
	name := fmt.Sprintf("projects/%s/agentPools/%s", projectID, poolName)

	pool, err := service.Projects.AgentPools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("agent pool %s does not exist in project %s", poolName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for agent pool %s in project %s: %w", poolName, projectID, err)
	}

	return pool, nil
}

// NewStorageTransferServiceE creates a Storage Transfer service authenticated the same way every
// other client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewStorageTransferServiceE(t testing.TestingT, ctx context.Context) (*storagetransfer.Service, error) {
	return storagetransfer.NewService(ctx, append(withOptions(), option.WithScopes(storagetransfer.CloudPlatformScope))...)
}
