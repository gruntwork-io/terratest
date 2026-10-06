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

// GetDNSPolicyAttrs returns the settings Google Cloud holds for the given Cloud DNS policy, so a
// test can assert on what was actually created rather than only that it exists. The networks the
// policy applies to come back with it.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDNSPolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, policyName string) *dns.Policy {
	policy, err := GetDNSPolicyAttrsE(t, ctx, projectID, policyName)
	require.NoError(t, err)

	return policy
}

// GetDNSPolicyAttrsE returns the settings Google Cloud holds for the given Cloud DNS policy.
// The ctx parameter supports cancellation and timeouts.
func GetDNSPolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, policyName string) (*dns.Policy, error) {
	logger.Default.Logf(t, "Getting settings for DNS policy %s in project %s", policyName, projectID)

	service, err := NewDNSServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDNSPolicyAttrsWithClient(ctx, service, projectID, policyName)
}

// GetDNSPolicyAttrsWithClient returns the settings Google Cloud holds for the given Cloud DNS
// policy using the supplied *dns.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see dns_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDNSPolicyAttrsWithClient(ctx context.Context, service *dns.Service, projectID string, policyName string) (*dns.Policy, error) {
	policy, err := service.Policies.Get(projectID, policyName).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("DNS policy %s does not exist in project %s", policyName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for DNS policy %s in project %s: %w", policyName, projectID, err)
	}

	return policy, nil
}

// GetDNSRecordSetAttrs returns the settings Google Cloud holds for the given record set, so a test
// can assert on what was actually created rather than only that it exists. A record set is named by
// its managed zone, its fully qualified name with the trailing dot, and its type together, since one
// name can carry a record of every type.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDNSRecordSetAttrs(t testing.TestingT, ctx context.Context, projectID string, zoneName string, recordName string, recordType string) *dns.ResourceRecordSet {
	recordSet, err := GetDNSRecordSetAttrsE(t, ctx, projectID, zoneName, recordName, recordType)
	require.NoError(t, err)

	return recordSet
}

// GetDNSRecordSetAttrsE returns the settings Google Cloud holds for the given record set.
// The ctx parameter supports cancellation and timeouts.
func GetDNSRecordSetAttrsE(t testing.TestingT, ctx context.Context, projectID string, zoneName string, recordName string, recordType string) (*dns.ResourceRecordSet, error) {
	logger.Default.Logf(t, "Getting settings for %s record %s in DNS managed zone %s in project %s", recordType, recordName, zoneName, projectID)

	service, err := NewDNSServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDNSRecordSetAttrsWithClient(ctx, service, projectID, zoneName, recordName, recordType)
}

// GetDNSRecordSetAttrsWithClient returns the settings Google Cloud holds for the given record set
// using the supplied *dns.Service. Prefer this variant in unit tests where the service is backed by
// an httptest fake server (see dns_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDNSRecordSetAttrsWithClient(ctx context.Context, service *dns.Service, projectID string, zoneName string, recordName string, recordType string) (*dns.ResourceRecordSet, error) {
	recordSet, err := service.ResourceRecordSets.Get(projectID, zoneName, recordName, recordType).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("%s record %s does not exist in DNS managed zone %s in project %s", recordType, recordName, zoneName, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for %s record %s in DNS managed zone %s in project %s: %w", recordType, recordName, zoneName, projectID, err)
	}

	return recordSet, nil
}

// NewDNSServiceE creates a Cloud DNS service authenticated the same way every other client in this
// module is.
// The ctx parameter supports cancellation and timeouts.
func NewDNSServiceE(t testing.TestingT, ctx context.Context) (*dns.Service, error) {
	return dns.NewService(ctx, append(withOptions(), option.WithScopes(dns.CloudPlatformScope))...)
}

// GetDNSResponsePolicyAttrs returns the settings Google Cloud holds for the given DNS response policy,
// so a test can assert on what was actually created rather than only that it exists. A response policy
// answers a query before the zone does, which is how a private name is made to resolve differently
// inside a network.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetDNSResponsePolicyAttrs(t testing.TestingT, ctx context.Context, projectID string, policyID string) *dns.ResponsePolicy {
	policy, err := GetDNSResponsePolicyAttrsE(t, ctx, projectID, policyID)
	require.NoError(t, err)

	return policy
}

// GetDNSResponsePolicyAttrsE returns the settings Google Cloud holds for the given response policy.
// The ctx parameter supports cancellation and timeouts.
func GetDNSResponsePolicyAttrsE(t testing.TestingT, ctx context.Context, projectID string, policyID string) (*dns.ResponsePolicy, error) {
	logger.Default.Logf(t, "Getting settings for DNS response policy %s in project %s", policyID, projectID)

	service, err := NewDNSServiceE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetDNSResponsePolicyAttrsWithClient(ctx, service, projectID, policyID)
}

// GetDNSResponsePolicyAttrsWithClient returns the settings Google Cloud holds for the given response
// policy using the supplied *dns.Service. Prefer this variant in unit tests where the service is
// backed by an httptest fake server (see dns_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetDNSResponsePolicyAttrsWithClient(ctx context.Context, service *dns.Service, projectID string, policyID string) (*dns.ResponsePolicy, error) {
	policy, err := service.ResponsePolicies.Get(projectID, policyID).Context(ctx).Do()
	if err != nil {
		var apiErr *googleapi.Error
		if errors.As(err, &apiErr) && apiErr.Code == 404 {
			return nil, fmt.Errorf("the DNS response policy %s does not exist in project %s", policyID, projectID)
		}

		return nil, fmt.Errorf("failed to get settings for DNS response policy %s in project %s: %w", policyID, projectID, err)
	}

	return policy, nil
}
