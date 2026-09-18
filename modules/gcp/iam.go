package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/option"
)

// GetServiceAccountAttrs returns the settings Google Cloud holds for the given service account, so
// a test can assert on what was actually created rather than only that it exists. The email is the
// full one, `<account-id>@<project>.iam.gserviceaccount.com`.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetServiceAccountAttrs(t testing.TestingT, ctx context.Context, projectID string, email string) *iam.ServiceAccount {
	account, err := GetServiceAccountAttrsE(t, ctx, projectID, email)
	require.NoError(t, err)

	return account
}

// GetServiceAccountAttrsE returns the settings Google Cloud holds for the given service account.
// The ctx parameter supports cancellation and timeouts.
func GetServiceAccountAttrsE(t testing.TestingT, ctx context.Context, projectID string, email string) (*iam.ServiceAccount, error) {
	logger.Default.Logf(t, "Getting settings for service account %s in project %s", email, projectID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetServiceAccountAttrsWithClient(ctx, service, projectID, email)
}

// GetServiceAccountAttrsWithClient returns the settings Google Cloud holds for the given service
// account using the supplied *iam.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see iam_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetServiceAccountAttrsWithClient(ctx context.Context, service *iam.Service, projectID string, email string) (*iam.ServiceAccount, error) {
	name := fmt.Sprintf("projects/%s/serviceAccounts/%s", projectID, email)

	account, err := service.Projects.ServiceAccounts.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("service account %s does not exist in project %s", email, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for service account %s in project %s: %w", email, projectID, err)
	}

	return account, nil
}

// GetWorkloadIdentityPoolAttrs returns the settings Google Cloud holds for the given workload
// identity pool, so a test can assert on what was actually created rather than only that it exists.
// The pool lives in a location, which is `global` for every pool the console creates.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetWorkloadIdentityPoolAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string) *iam.WorkloadIdentityPool {
	pool, err := GetWorkloadIdentityPoolAttrsE(t, ctx, projectID, location, poolID)
	require.NoError(t, err)

	return pool
}

// GetWorkloadIdentityPoolAttrsE returns the settings Google Cloud holds for the given workload
// identity pool.
// The ctx parameter supports cancellation and timeouts.
func GetWorkloadIdentityPoolAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, poolID string) (*iam.WorkloadIdentityPool, error) {
	logger.Default.Logf(t, "Getting settings for workload identity pool %s in project %s", poolID, projectID)

	service, err := NewIAMServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetWorkloadIdentityPoolAttrsWithClient(ctx, service, projectID, location, poolID)
}

// GetWorkloadIdentityPoolAttrsWithClient returns the settings Google Cloud holds for the given
// workload identity pool using the supplied *iam.Service. Prefer this variant in unit tests where
// the service is backed by an httptest fake server (see iam_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetWorkloadIdentityPoolAttrsWithClient(ctx context.Context, service *iam.Service, projectID string, location string, poolID string) (*iam.WorkloadIdentityPool, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/workloadIdentityPools/%s", projectID, location, poolID)

	pool, err := service.Projects.Locations.WorkloadIdentityPools.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("workload identity pool %s does not exist in project %s", poolID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for workload identity pool %s in project %s: %w", poolID, projectID, err)
	}

	return pool, nil
}

// NewIAMServiceE creates an IAM service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewIAMServiceE(t testing.TestingT, ctx context.Context) (*iam.Service, error) {
	return iam.NewService(ctx, append(withOptions(), option.WithScopes(iam.CloudPlatformScope))...)
}
