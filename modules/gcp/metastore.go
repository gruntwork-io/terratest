package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/metastore/v1"
	"google.golang.org/api/option"
)

// GetMetastoreServiceAttrs returns the settings Google Cloud holds for the given Dataproc Metastore service, so a test can assert on what was
// actually created rather than only that it exists.
// The service is the Hive metastore itself, so its version, its tier and the network it is reachable on decide what can query it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreServiceAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string) *metastore.Service {
	attrs, err := GetMetastoreServiceAttrsE(t, ctx, projectID, location, serviceID)
	require.NoError(t, err)

	return attrs
}

// GetMetastoreServiceAttrsE returns the settings Google Cloud holds for the given Dataproc Metastore service.
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreServiceAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string) (*metastore.Service, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc Metastore service %s in %s in project %s", serviceID, location, projectID)

	service, err := NewMetastoreServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMetastoreServiceAttrsWithClient(ctx, service, projectID, location, serviceID)
}

// GetMetastoreServiceAttrsWithClient returns the settings Google Cloud holds for the given Dataproc Metastore service using the supplied
// *metastore.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see metastore_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreServiceAttrsWithClient(ctx context.Context, service *metastore.APIService, projectID string, location string, serviceID string) (*metastore.Service, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/services/%s", projectID, location, serviceID)

	attrs, err := service.Projects.Locations.Services.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc Metastore service %s in %s in project %s does not exist", serviceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc Metastore service %s in %s in project %s: %w", serviceID, location, projectID, err)
	}

	return attrs, nil
}

// GetMetastoreFederationAttrs returns the settings Google Cloud holds for the given Dataproc Metastore federation, so a test can assert on what was
// actually created rather than only that it exists.
// A federation presents several metastores as one endpoint, so which backends it names is the whole point of it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreFederationAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, federationID string) *metastore.Federation {
	attrs, err := GetMetastoreFederationAttrsE(t, ctx, projectID, location, federationID)
	require.NoError(t, err)

	return attrs
}

// GetMetastoreFederationAttrsE returns the settings Google Cloud holds for the given Dataproc Metastore federation.
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreFederationAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, federationID string) (*metastore.Federation, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc Metastore federation %s in %s in project %s", federationID, location, projectID)

	service, err := NewMetastoreServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMetastoreFederationAttrsWithClient(ctx, service, projectID, location, federationID)
}

// GetMetastoreFederationAttrsWithClient returns the settings Google Cloud holds for the given Dataproc Metastore federation using the supplied
// *metastore.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see metastore_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreFederationAttrsWithClient(ctx context.Context, service *metastore.APIService, projectID string, location string, federationID string) (*metastore.Federation, error) {
	name := fmt.Sprintf("projects/%s/locations/%s/federations/%s", projectID, location, federationID)

	attrs, err := service.Projects.Locations.Federations.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc Metastore federation %s in %s in project %s does not exist", federationID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc Metastore federation %s in %s in project %s: %w", federationID, location, projectID, err)
	}

	return attrs, nil
}

// GetMetastoreServiceIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Dataproc Metastore service, so a test can assert on what was
// actually created rather than only that it exists.
// Who may query or administer the metastore is what the policy decides, and a service Google created on its own carries no binding at all.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreServiceIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string) *metastore.Policy {
	policy, err := GetMetastoreServiceIamPolicyAttrsE(t, ctx, projectID, location, serviceID)
	require.NoError(t, err)

	return policy
}

// GetMetastoreServiceIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given Dataproc Metastore service.
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreServiceIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string) (*metastore.Policy, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc Metastore service %s in %s in project %s", serviceID, location, projectID)

	service, err := NewMetastoreServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMetastoreServiceIamPolicyAttrsWithClient(ctx, service, projectID, location, serviceID)
}

// GetMetastoreServiceIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given Dataproc Metastore service using the supplied
// *metastore.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see metastore_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreServiceIamPolicyAttrsWithClient(ctx context.Context, service *metastore.APIService, projectID string, location string, serviceID string) (*metastore.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/services/%s", projectID, location, serviceID)

	policy, err := service.Projects.Locations.Services.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc Metastore service %s in %s in project %s does not exist", serviceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc Metastore service %s in %s in project %s: %w", serviceID, location, projectID, err)
	}

	return policy, nil
}

// GetMetastoreFederationIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Dataproc Metastore federation, so a test can assert on what was
// actually created rather than only that it exists.
// Who may query through the federation is what the policy decides, and it is separate from the policy on any metastore behind it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreFederationIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, federationID string) *metastore.Policy {
	policy, err := GetMetastoreFederationIamPolicyAttrsE(t, ctx, projectID, location, federationID)
	require.NoError(t, err)

	return policy
}

// GetMetastoreFederationIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given Dataproc Metastore federation.
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreFederationIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, federationID string) (*metastore.Policy, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc Metastore federation %s in %s in project %s", federationID, location, projectID)

	service, err := NewMetastoreServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMetastoreFederationIamPolicyAttrsWithClient(ctx, service, projectID, location, federationID)
}

// GetMetastoreFederationIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given Dataproc Metastore federation using the supplied
// *metastore.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see metastore_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreFederationIamPolicyAttrsWithClient(ctx context.Context, service *metastore.APIService, projectID string, location string, federationID string) (*metastore.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/federations/%s", projectID, location, federationID)

	policy, err := service.Projects.Locations.Federations.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc Metastore federation %s in %s in project %s does not exist", federationID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc Metastore federation %s in %s in project %s: %w", federationID, location, projectID, err)
	}

	return policy, nil
}

// GetMetastoreDatabaseIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Dataproc Metastore database, so a test can assert on what was
// actually created rather than only that it exists.
// A database inside a metastore carries its own policy, so access can be granted to one database rather than to the whole service.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreDatabaseIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string, databaseID string) *metastore.Policy {
	policy, err := GetMetastoreDatabaseIamPolicyAttrsE(t, ctx, projectID, location, serviceID, databaseID)
	require.NoError(t, err)

	return policy
}

// GetMetastoreDatabaseIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given Dataproc Metastore database.
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreDatabaseIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string, databaseID string) (*metastore.Policy, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc Metastore database %s %s in %s in project %s", databaseID, serviceID, location, projectID)

	service, err := NewMetastoreServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMetastoreDatabaseIamPolicyAttrsWithClient(ctx, service, projectID, location, serviceID, databaseID)
}

// GetMetastoreDatabaseIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given Dataproc Metastore database using the supplied
// *metastore.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see metastore_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreDatabaseIamPolicyAttrsWithClient(ctx context.Context, service *metastore.APIService, projectID string, location string, serviceID string, databaseID string) (*metastore.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/services/%s/databases/%s", projectID, location, serviceID, databaseID)

	policy, err := service.Projects.Locations.Services.Databases.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc Metastore database %s %s in %s in project %s does not exist", databaseID, serviceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc Metastore database %s %s in %s in project %s: %w", databaseID, serviceID, location, projectID, err)
	}

	return policy, nil
}

// GetMetastoreTableIamPolicyAttrs returns the IAM policy Google Cloud holds for the given Dataproc Metastore table, so a test can assert on what was
// actually created rather than only that it exists.
// A table carries its own policy, which is the finest grain Metastore grants access at.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreTableIamPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string, databaseID string, tableID string) *metastore.Policy {
	policy, err := GetMetastoreTableIamPolicyAttrsE(t, ctx, projectID, location, serviceID, databaseID, tableID)
	require.NoError(t, err)

	return policy
}

// GetMetastoreTableIamPolicyAttrsE returns the IAM policy Google Cloud holds for the given Dataproc Metastore table.
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreTableIamPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, location string, serviceID string, databaseID string, tableID string) (*metastore.Policy, error) {
	logger.Default.Logf(t, "Getting settings for Dataproc Metastore table %s %s %s in %s in project %s", tableID, databaseID, serviceID, location, projectID)

	service, err := NewMetastoreServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetMetastoreTableIamPolicyAttrsWithClient(ctx, service, projectID, location, serviceID, databaseID, tableID)
}

// GetMetastoreTableIamPolicyAttrsWithClient returns the IAM policy Google Cloud holds for the given Dataproc Metastore table using the supplied
// *metastore.APIService. Prefer this variant in unit tests where the service is backed by an httptest fake
// server (see metastore_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetMetastoreTableIamPolicyAttrsWithClient(ctx context.Context, service *metastore.APIService, projectID string, location string, serviceID string, databaseID string, tableID string) (*metastore.Policy, error) {
	resource := fmt.Sprintf("projects/%s/locations/%s/services/%s/databases/%s/tables/%s", projectID, location, serviceID, databaseID, tableID)

	policy, err := service.Projects.Locations.Services.Databases.Tables.GetIamPolicy(resource).OptionsRequestedPolicyVersion(iamPolicyVersionWithConditions).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Dataproc Metastore table %s %s %s in %s in project %s does not exist", tableID, databaseID, serviceID, location, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Dataproc Metastore table %s %s %s in %s in project %s: %w", tableID, databaseID, serviceID, location, projectID, err)
	}

	return policy, nil
}

// NewMetastoreServiceE creates a Dataproc Metastore service authenticated the same way every other client
// in this module is. The Go client calls this type APIService, because Metastore already spends the name
// Service on a resource: a metastore service is the Hive metastore itself.
// The ctx parameter supports cancellation and timeouts.
func NewMetastoreServiceE(t testing.TestingT, ctx context.Context) (*metastore.APIService, error) {
	return metastore.NewService(ctx, append(withOptions(), option.WithScopes(metastore.CloudPlatformScope))...)
}
