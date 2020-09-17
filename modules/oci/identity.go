package oci

import (
	"context"
	"fmt"
	"os"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/identity"
)

// GetRandomAvailabilityDomain gets a randomly chosen availability domain for given compartment.
// The returned value can be overridden by of the environment variable TF_VAR_availability_domain.
func GetRandomAvailabilityDomain(t testing.TestingT, compartmentID string) string {
	ad, err := GetRandomAvailabilityDomainE(t, compartmentID)
	if err != nil {
		t.Fatal(err)
	}
	return ad
}

// GetRandomAvailabilityDomainE gets a randomly chosen availability domain for given compartment.
// The returned value can be overridden by of the environment variable TF_VAR_availability_domain.
func GetRandomAvailabilityDomainE(t testing.TestingT, compartmentID string) (string, error) {
	adFromEnvVar := os.Getenv(availabilityDomainEnvVar)
	if adFromEnvVar != "" {
		logger.Logf(t, "Using availability domain %s from environment variable %s", adFromEnvVar, availabilityDomainEnvVar)
		return adFromEnvVar, nil
	}

	allADs, err := GetAllAvailabilityDomainsE(t, compartmentID)
	if err != nil {
		return "", err
	}

	ad := random.RandomString(allADs)

	logger.Logf(t, "Using availability domain %s", ad)
	return ad, nil
}

// GetAllAvailabilityDomains gets the list of availability domains available in the given compartment.
func GetAllAvailabilityDomains(t testing.TestingT, compartmentID string) []string {
	ads, err := GetAllAvailabilityDomainsE(t, compartmentID)
	if err != nil {
		t.Fatal(err)
	}
	return ads
}

// GetAllAvailabilityDomainsE gets the list of availability domains available in the given compartment.
func GetAllAvailabilityDomainsE(t testing.TestingT, compartmentID string) ([]string, error) {
	configProvider := common.DefaultConfigProvider()
	client, err := identity.NewIdentityClientWithConfigurationProvider(configProvider)
	if err != nil {
		return nil, err
	}

	request := identity.ListAvailabilityDomainsRequest{CompartmentId: &compartmentID}
	response, err := client.ListAvailabilityDomains(context.Background(), request)
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("No availability domains found in the %s compartment", compartmentID)
	}

	return availabilityDomainsNames(response.Items), nil
}

func availabilityDomainsNames(ads []identity.AvailabilityDomain) []string {
	names := []string{}
	for _, ad := range ads {
		names = append(names, *ad.Name)
	}
	return names
}
