package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/observability/v1"
	"google.golang.org/api/option"
)

// GetTraceScopeAttrs returns the settings Google Cloud holds for the given trace scope, so a test
// can assert on which projects' traces it gathers rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetTraceScopeAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, scopeID string) *observability.TraceScope {
	scope, err := GetTraceScopeAttrsE(t, ctx, projectID, location, scopeID)
	require.NoError(t, err)

	return scope
}

// GetTraceScopeAttrsE returns the settings Google Cloud holds for the given trace scope.
// The ctx parameter supports cancellation and timeouts.
func GetTraceScopeAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, scopeID string) (*observability.TraceScope, error) {
	logger.Default.Logf(t, "Getting settings for trace scope %s in %s in project %s", scopeID, location, projectID)

	service, err := NewObservabilityServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetTraceScopeAttrsWithClient(ctx, service, projectID, location, scopeID)
}

// GetTraceScopeAttrsWithClient returns the settings Google Cloud holds for the given trace scope
// using the supplied *observability.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see observability_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetTraceScopeAttrsWithClient(ctx context.Context, service *observability.Service, projectID string, location string, scopeID string) (*observability.TraceScope, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/traceScopes/%s", projectID, location, scopeID)

	scope, err := service.Projects.Locations.TraceScopes.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the trace scope %s in %s in project %s does not exist", scopeID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for trace scope %s in %s in project %s: %w", scopeID, location, projectID, err)
	}

	return scope, nil
}

// NewObservabilityServiceE creates an Observability service authenticated the same way every other
// client in this module is.
// The ctx parameter supports cancellation and timeouts.
func NewObservabilityServiceE(t testing.TestingT, ctx context.Context) (*observability.Service, error) {
	return observability.NewService(ctx, append(withOptions(), option.WithScopes(observability.CloudPlatformScope))...)
}
