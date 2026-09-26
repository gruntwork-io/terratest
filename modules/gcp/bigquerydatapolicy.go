package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/bigquerydatapolicy/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetBigQueryDataPolicyAttrs returns the settings Google Cloud holds for the given BigQuery data policy, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryDataPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) *bigquerydatapolicy.DataPolicy {
	policy, err := GetBigQueryDataPolicyAttrsE(t, ctx, projectID, location, policyID)
	require.NoError(t, err)

	return policy
}

// GetBigQueryDataPolicyAttrsE returns the settings Google Cloud holds for the given BigQuery data policy.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryDataPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) (*bigquerydatapolicy.DataPolicy, error) {
	logger.Default.Logf(t, "Getting settings for BigQuery data policy %s in %s in project %s", policyID, location, projectID)

	service, err := NewBigQueryDataPolicyServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigQueryDataPolicyAttrsWithClient(ctx, service, projectID, location, policyID)
}

// GetBigQueryDataPolicyAttrsWithClient returns the settings Google Cloud holds for the given BigQuery data policy using the supplied
// *bigquerydatapolicy.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see bigquerydatapolicy_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryDataPolicyAttrsWithClient(ctx context.Context, service *bigquerydatapolicy.Service, projectID string, location string, policyID string) (*bigquerydatapolicy.DataPolicy, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/dataPolicies/%s", projectID, location, policyID)

	policy, err := service.Projects.Locations.DataPolicies.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the BigQuery data policy %s does not exist in %s in project %s", policyID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for BigQuery data policy %s in %s in project %s: %w", policyID, location, projectID, err)
	}

	return policy, nil
}

// NewBigQueryDataPolicyServiceE creates a BigQuery Data Policy service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewBigQueryDataPolicyServiceE(t testing.TestingT, ctx context.Context) (*bigquerydatapolicy.Service, error) {
	return bigquerydatapolicy.NewService(ctx, append(withOptions(), option.WithScopes(bigquerydatapolicy.CloudPlatformScope))...)
}
