package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/managedidentities/v1"
	"google.golang.org/api/option"
)

// GetActiveDirectoryDomainAttrs returns the settings Google Cloud holds for the given Managed Microsoft AD domain, so a test can assert on
// what was actually created rather than only that it exists.
// A domain is always global, so the location is not a parameter.
// One read answers for two modules: a trust is an entry in the domain's own trusts list rather than a
// resource with a read of its own.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetActiveDirectoryDomainAttrs(t testing.TestingT, ctx context.Context, projectID string, domainName string) *managedidentities.Domain {
	result, err := GetActiveDirectoryDomainAttrsE(t, ctx, projectID, domainName)
	require.NoError(t, err)

	return result
}

// GetActiveDirectoryDomainAttrsE returns the settings Google Cloud holds for the given Managed Microsoft AD domain.
// The ctx parameter supports cancellation and timeouts.
func GetActiveDirectoryDomainAttrsE(t testing.TestingT, ctx context.Context, projectID string, domainName string) (*managedidentities.Domain, error) {
	logger.Default.Logf(t, "Getting settings for Managed Microsoft AD domain %s in project %s", domainName, projectID)

	service, err := NewManagedIdentitiesServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetActiveDirectoryDomainAttrsWithClient(ctx, service, projectID, domainName)
}

// GetActiveDirectoryDomainAttrsWithClient returns the settings Google Cloud holds for the given Managed Microsoft AD domain using the supplied
// *managedidentities.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see managedidentities_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetActiveDirectoryDomainAttrsWithClient(ctx context.Context, service *managedidentities.Service, projectID string, domainName string) (*managedidentities.Domain, error) {
	name := fmt.Sprintf("projects/%s/locations/global/domains/%s", projectID, domainName)

	result, err := service.Projects.Locations.Global.Domains.Get(name).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the Managed Microsoft AD domain %s does not exist in project %s", domainName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Managed Microsoft AD domain %s in project %s: %w", domainName, projectID, err)
	}

	return result, nil
}

// NewManagedIdentitiesServiceE creates a Managed Microsoft AD service authenticated the same way every other client in
// this module is.
// The ctx parameter supports cancellation and timeouts.
func NewManagedIdentitiesServiceE(t testing.TestingT, ctx context.Context) (*managedidentities.Service, error) {
	return managedidentities.NewService(ctx, append(withOptions(), option.WithScopes(managedidentities.CloudPlatformScope))...)
}
