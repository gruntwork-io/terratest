package gcp

import (
	"context"
	"errors"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

// GetManagedZoneAttrs returns the settings Google Cloud holds for the given Cloud DNS managed zone,
// so a test can assert on what was actually created rather than only that it exists.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetManagedZoneAttrs(t testing.TestingT, ctx context.Context, projectID string, zoneName string) *dns.ManagedZone {
	zone, err := GetManagedZoneAttrsE(t, ctx, projectID, zoneName)
	require.NoError(t, err)

	return zone
}

// GetManagedZoneAttrsE returns the settings Google Cloud holds for the given Cloud DNS managed zone.
// The ctx parameter supports cancellation and timeouts.
func GetManagedZoneAttrsE(t testing.TestingT, ctx context.Context, projectID string, zoneName string) (*dns.ManagedZone, error) {
	logger.Default.Logf(t, "Getting settings for Cloud DNS managed zone %s in project %s", zoneName, projectID)

	service, err := NewDNSServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetManagedZoneAttrsWithClient(ctx, service, projectID, zoneName)
}

// GetManagedZoneAttrsWithClient returns the settings Google Cloud holds for the given Cloud DNS
// managed zone using the supplied *dns.Service. Prefer this variant in unit tests where the service
// is backed by an httptest fake server (see dns_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetManagedZoneAttrsWithClient(ctx context.Context, service *dns.Service, projectID string, zoneName string) (*dns.ManagedZone, error) {
	zone, err := service.ManagedZones.Get(projectID, zoneName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("Cloud DNS managed zone %s does not exist in project %s", zoneName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for Cloud DNS managed zone %s in project %s: %w", zoneName, projectID, err)
	}

	return zone, nil
}

// NewDNSServiceE creates a Cloud DNS service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewDNSServiceE(t testing.TestingT, ctx context.Context) (*dns.Service, error) {
	return dns.NewService(ctx, append(withOptions(), option.WithScopes(dns.CloudPlatformScope))...)
}
