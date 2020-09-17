package oci

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
)

// GetRandomSubnetID gets a randomly chosen subnet OCID in the given availability domain.
// The returned value can be overridden by of the environment variable TF_VAR_subnet_ocid.
func GetRandomSubnetID(t testing.TestingT, compartmentID string, availabilityDomain string) string {
	ocid, err := GetRandomSubnetIDE(t, compartmentID, availabilityDomain)
	if err != nil {
		t.Fatal(err)
	}
	return ocid
}

// GetRandomSubnetIDE gets a randomly chosen subnet OCID in the given availability domain.
// The returned value can be overridden by of the environment variable TF_VAR_subnet_ocid.
func GetRandomSubnetIDE(t testing.TestingT, compartmentID string, availabilityDomain string) (string, error) {
	configProvider := common.DefaultConfigProvider()
	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(configProvider)
	if err != nil {
		return "", err
	}

	vcnIDs, err := GetAllVcnIDsE(t, compartmentID)
	if err != nil {
		return "", err
	}

	allSubnetIDs := map[string][]string{}
	for _, vcnID := range vcnIDs {
		request := core.ListSubnetsRequest{
			CompartmentId: &compartmentID,
			VcnId:         &vcnID,
		}
		response, err := client.ListSubnets(context.Background(), request)
		if err != nil {
			return "", err
		}

		mapSubnetsByAvailabilityDomain(allSubnetIDs, response.Items)
	}

	subnetID := random.RandomString(allSubnetIDs[availabilityDomain])

	logger.Logf(t, "Using subnet with OCID %s", subnetID)
	return subnetID, nil
}

// GetAllVcnIDs gets the list of VCNs available in the given compartment.
func GetAllVcnIDs(t testing.TestingT, compartmentID string) []string {
	vcnIDS, err := GetAllVcnIDsE(t, compartmentID)
	if err != nil {
		t.Fatal(err)
	}
	return vcnIDS
}

// GetAllVcnIDsE gets the list of VCNs available in the given compartment.
func GetAllVcnIDsE(t testing.TestingT, compartmentID string) ([]string, error) {
	configProvider := common.DefaultConfigProvider()
	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(configProvider)
	if err != nil {
		return nil, err
	}

	request := core.ListVcnsRequest{CompartmentId: &compartmentID}
	response, err := client.ListVcns(context.Background(), request)
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("No VCNs found in the %s compartment", compartmentID)
	}

	return vcnsIDs(response.Items), nil
}

func mapSubnetsByAvailabilityDomain(allSubnets map[string][]string, subnets []core.Subnet) map[string][]string {
	for _, subnet := range subnets {
		allSubnets[*subnet.AvailabilityDomain] = append(allSubnets[*subnet.AvailabilityDomain], *subnet.Id)
	}
	return allSubnets
}

func vcnsIDs(vcns []core.Vcn) []string {
	ids := []string{}
	for _, vcn := range vcns {
		ids = append(ids, *vcn.Id)
	}
	return ids
}
