package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	bigquerydatapolicyv1 "google.golang.org/api/bigquerydatapolicy/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// The v2 data policy surface is a separate API with its own reads, in bigquerydatapolicy.go. A policy
// created on one is not visible to the other, so the version is part of every name here.

// GetBigQueryDataPolicyV1Attrs returns the settings Google Cloud holds for the given BigQuery data policy, so a test can assert on
// what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryDataPolicyV1Attrs(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) *bigquerydatapolicyv1.DataPolicy {
	policy, err := GetBigQueryDataPolicyV1AttrsE(t, ctx, projectID, location, policyID)
	require.NoError(t, err)

	return policy
}

// GetBigQueryDataPolicyV1AttrsE returns the settings Google Cloud holds for the given BigQuery data policy.
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryDataPolicyV1AttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, policyID string) (*bigquerydatapolicyv1.DataPolicy, error) {
	logger.Default.Logf(t, "Getting settings for BigQuery data policy %s in %s in project %s", policyID, location, projectID)

	service, err := NewBigQueryDataPolicyV1ServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetBigQueryDataPolicyV1AttrsWithClient(ctx, service, projectID, location, policyID)
}

// GetBigQueryDataPolicyV1AttrsWithClient returns the settings Google Cloud holds for the given BigQuery data policy using the supplied
// *bigquerydatapolicyv1.Service. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see bigquerydatapolicyv1_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetBigQueryDataPolicyV1AttrsWithClient(ctx context.Context, service *bigquerydatapolicyv1.Service, projectID string, location string, policyID string) (*bigquerydatapolicyv1.DataPolicy, error) {
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

// NewBigQueryDataPolicyV1ServiceE creates a BigQuery Data Policy service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewBigQueryDataPolicyV1ServiceE(t testing.TestingT, ctx context.Context) (*bigquerydatapolicyv1.Service, error) {
	return bigquerydatapolicyv1.NewService(ctx, append(withOptions(), option.WithScopes(bigquerydatapolicyv1.CloudPlatformScope))...)
}
