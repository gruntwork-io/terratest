package oci

import (
	"context"
	"os"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/identity"
	"github.com/stretchr/testify/require"
)

// GetAllAvailabilityDomainsContextE gets the list of availability domains available in the given compartment.
// The ctx parameter supports cancellation and timeouts.
func GetAllAvailabilityDomainsContextE(t testing.TestingT, ctx context.Context, compartmentID string) ([]string, error) {
	configProvider := common.DefaultConfigProvider()

	client, err := identity.NewIdentityClientWithConfigurationProvider(configProvider)
	if err != nil {
		return nil, err
	}

	request := identity.ListAvailabilityDomainsRequest{CompartmentId: &compartmentID}

	response, err := client.ListAvailabilityDomains(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, NoAvailabilityDomainsFoundError{CompartmentID: compartmentID}
	}

	return availabilityDomainsNames(response.Items), nil
}

// GetAllAvailabilityDomainsContext gets the list of availability domains available in the given compartment.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAllAvailabilityDomainsContext(t testing.TestingT, ctx context.Context, compartmentID string) []string {
	t.Helper()

	ads, err := GetAllAvailabilityDomainsContextE(t, ctx, compartmentID)
	require.NoError(t, err)

	return ads
}

// GetAllAvailabilityDomains gets the list of availability domains available in the given compartment.
//
// Deprecated: Use [GetAllAvailabilityDomainsContext] instead.
func GetAllAvailabilityDomains(t testing.TestingT, compartmentID string) []string {
	t.Helper()

	return GetAllAvailabilityDomainsContext(t, context.Background(), compartmentID)
}

// GetAllAvailabilityDomainsE gets the list of availability domains available in the given compartment.
//
// Deprecated: Use [GetAllAvailabilityDomainsContextE] instead.
func GetAllAvailabilityDomainsE(t testing.TestingT, compartmentID string) ([]string, error) {
	return GetAllAvailabilityDomainsContextE(t, context.Background(), compartmentID)
}

// GetRandomAvailabilityDomainContextE gets a randomly chosen availability domain for given compartment.
// The returned value can be overridden by of the environment variable TF_VAR_availability_domain.
// The ctx parameter supports cancellation and timeouts.
func GetRandomAvailabilityDomainContextE(t testing.TestingT, ctx context.Context, compartmentID string) (string, error) {
	adFromEnvVar := os.Getenv(availabilityDomainEnvVar)
	if adFromEnvVar != "" {
		logger.Default.Logf(t, "Using availability domain %s from environment variable %s", adFromEnvVar, availabilityDomainEnvVar)
		return adFromEnvVar, nil
	}

	allADs, err := GetAllAvailabilityDomainsContextE(t, ctx, compartmentID)
	if err != nil {
		return "", err
	}

	ad := random.RandomString(allADs)

	logger.Default.Logf(t, "Using availability domain %s", ad)

	return ad, nil
}

// GetRandomAvailabilityDomainContext gets a randomly chosen availability domain for given compartment.
// The returned value can be overridden by of the environment variable TF_VAR_availability_domain.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRandomAvailabilityDomainContext(t testing.TestingT, ctx context.Context, compartmentID string) string {
	t.Helper()

	ad, err := GetRandomAvailabilityDomainContextE(t, ctx, compartmentID)
	require.NoError(t, err)

	return ad
}

// GetRandomAvailabilityDomain gets a randomly chosen availability domain for given compartment.
// The returned value can be overridden by of the environment variable TF_VAR_availability_domain.
//
// Deprecated: Use [GetRandomAvailabilityDomainContext] instead.
func GetRandomAvailabilityDomain(t testing.TestingT, compartmentID string) string {
	t.Helper()

	return GetRandomAvailabilityDomainContext(t, context.Background(), compartmentID)
}

// GetRandomAvailabilityDomainE gets a randomly chosen availability domain for given compartment.
// The returned value can be overridden by of the environment variable TF_VAR_availability_domain.
//
// Deprecated: Use [GetRandomAvailabilityDomainContextE] instead.
func GetRandomAvailabilityDomainE(t testing.TestingT, compartmentID string) (string, error) {
	return GetRandomAvailabilityDomainContextE(t, context.Background(), compartmentID)
}

func availabilityDomainsNames(ads []identity.AvailabilityDomain) []string {
	names := make([]string, 0, len(ads))
	for _, ad := range ads {
		names = append(names, *ad.Name)
	}

	return names
}
