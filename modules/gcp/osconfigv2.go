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
	osconfigv2 "google.golang.org/api/osconfig/v2"
)

// GetOSConfigPolicyOrchestratorAttrs returns the settings Google Cloud holds for the given OS Config policy orchestrator, so a test can assert on
// what was actually created rather than only that it exists.
// An orchestrator is always global, so the location is not a parameter.
// The rollout state Google reports is what tells a test whether the orchestrator is actually acting
// on the resources it selects.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetOSConfigPolicyOrchestratorAttrs(t testing.TestingT, ctx context.Context, projectID string, orchestratorID string) *osconfigv2.GoogleCloudOsconfigV2__PolicyOrchestrator {
	result, err := GetOSConfigPolicyOrchestratorAttrsE(t, ctx, projectID, orchestratorID)
	require.NoError(t, err)

	return result
}

// GetOSConfigPolicyOrchestratorAttrsE returns the settings Google Cloud holds for the given OS Config policy orchestrator.
// The ctx parameter supports cancellation and timeouts.
func GetOSConfigPolicyOrchestratorAttrsE(t testing.TestingT, ctx context.Context, projectID string, orchestratorID string) (*osconfigv2.GoogleCloudOsconfigV2__PolicyOrchestrator, error) {
	logger.Default.Logf(t, "Getting settings for OS Config policy orchestrator %s in project %s", orchestratorID, projectID)

	service, err := NewOSConfigV2ServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetOSConfigPolicyOrchestratorAttrsWithClient(ctx, service, projectID, orchestratorID)
}

// GetOSConfigPolicyOrchestratorAttrsWithClient returns the settings Google Cloud holds for the given OS Config policy orchestrator using the supplied
// *osconfigv2.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see osconfigv2_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetOSConfigPolicyOrchestratorAttrsWithClient(ctx context.Context, service *osconfigv2.Service, projectID string, orchestratorID string) (*osconfigv2.GoogleCloudOsconfigV2__PolicyOrchestrator, error) {
	name := fmt.Sprintf("projects/%s/locations/global/policyOrchestrators/%s", projectID, orchestratorID)

	result, err := service.Projects.Locations.Global.PolicyOrchestrators.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the OS Config policy orchestrator %s does not exist in project %s", orchestratorID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for OS Config policy orchestrator %s in project %s: %w", orchestratorID, projectID, err)
	}

	return result, nil
}

// NewOSConfigV2ServiceE creates a OS Config v2 service authenticated the same way every other client in
// this module is.
// Policy orchestrators are served only by the second version of the OS Config API, which is why this
// client sits beside the one osconfig.go builds.
// The ctx parameter supports cancellation and timeouts.
func NewOSConfigV2ServiceE(t testing.TestingT, ctx context.Context) (*osconfigv2.Service, error) {
	return osconfigv2.NewService(ctx, append(withOptions(), option.WithScopes(osconfigv2.CloudPlatformScope))...)
}
