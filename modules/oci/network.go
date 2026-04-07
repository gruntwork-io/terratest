package oci

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/oracle/oci-go-sdk/core"
)

// GetAllVcnIDsContextE gets the list of VCNs available in the given compartment.
// The ctx parameter supports cancellation and timeouts.
func GetAllVcnIDsContextE(t testing.TestingT, ctx context.Context, compartmentID string) ([]string, error) {
	configProvider := common.DefaultConfigProvider()

	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(configProvider)
	if err != nil {
		return nil, err
	}

	request := core.ListVcnsRequest{CompartmentId: &compartmentID}

	response, err := client.ListVcns(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, NoVCNsFoundError{CompartmentID: compartmentID}
	}

	return vcnsIDs(response.Items), nil
}

// GetAllVcnIDsContext gets the list of VCNs available in the given compartment.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAllVcnIDsContext(t testing.TestingT, ctx context.Context, compartmentID string) []string {
	t.Helper()

	vcnIDs, err := GetAllVcnIDsContextE(t, ctx, compartmentID)
	if err != nil {
		t.Fatal(err)
	}

	return vcnIDs
}

// GetAllVcnIDs gets the list of VCNs available in the given compartment.
//
// Deprecated: Use [GetAllVcnIDsContext] instead.
func GetAllVcnIDs(t testing.TestingT, compartmentID string) []string {
	t.Helper()

	return GetAllVcnIDsContext(t, context.Background(), compartmentID)
}

// GetAllVcnIDsE gets the list of VCNs available in the given compartment.
//
// Deprecated: Use [GetAllVcnIDsContextE] instead.
func GetAllVcnIDsE(t testing.TestingT, compartmentID string) ([]string, error) {
	return GetAllVcnIDsContextE(t, context.Background(), compartmentID)
}

// GetRandomSubnetIDContextE gets a randomly chosen subnet OCID in the given availability domain.
// The returned value can be overridden by of the environment variable TF_VAR_subnet_ocid.
// The ctx parameter supports cancellation and timeouts.
func GetRandomSubnetIDContextE(t testing.TestingT, ctx context.Context, compartmentID string, availabilityDomain string) (string, error) {
	configProvider := common.DefaultConfigProvider()

	client, err := core.NewVirtualNetworkClientWithConfigurationProvider(configProvider)
	if err != nil {
		return "", err
	}

	vcnIDs, err := GetAllVcnIDsContextE(t, ctx, compartmentID)
	if err != nil {
		return "", err
	}

	allSubnetIDs := map[string][]string{}

	for _, vcnID := range vcnIDs {
		request := core.ListSubnetsRequest{
			CompartmentId: &compartmentID,
			VcnId:         &vcnID,
		}

		response, err := client.ListSubnets(ctx, request)
		if err != nil {
			return "", err
		}

		mapSubnetsByAvailabilityDomain(allSubnetIDs, response.Items)
	}

	subnetID := random.RandomString(allSubnetIDs[availabilityDomain])

	logger.Default.Logf(t, "Using subnet with OCID %s", subnetID)

	return subnetID, nil
}

// GetRandomSubnetIDContext gets a randomly chosen subnet OCID in the given availability domain.
// The returned value can be overridden by of the environment variable TF_VAR_subnet_ocid.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRandomSubnetIDContext(t testing.TestingT, ctx context.Context, compartmentID string, availabilityDomain string) string {
	t.Helper()

	ocid, err := GetRandomSubnetIDContextE(t, ctx, compartmentID, availabilityDomain)
	if err != nil {
		t.Fatal(err)
	}

	return ocid
}

// GetRandomSubnetID gets a randomly chosen subnet OCID in the given availability domain.
// The returned value can be overridden by of the environment variable TF_VAR_subnet_ocid.
//
// Deprecated: Use [GetRandomSubnetIDContext] instead.
func GetRandomSubnetID(t testing.TestingT, compartmentID string, availabilityDomain string) string {
	t.Helper()

	return GetRandomSubnetIDContext(t, context.Background(), compartmentID, availabilityDomain)
}

// GetRandomSubnetIDE gets a randomly chosen subnet OCID in the given availability domain.
// The returned value can be overridden by of the environment variable TF_VAR_subnet_ocid.
//
// Deprecated: Use [GetRandomSubnetIDContextE] instead.
func GetRandomSubnetIDE(t testing.TestingT, compartmentID string, availabilityDomain string) (string, error) {
	return GetRandomSubnetIDContextE(t, context.Background(), compartmentID, availabilityDomain)
}

func mapSubnetsByAvailabilityDomain(allSubnets map[string][]string, subnets []core.Subnet) map[string][]string {
	for i := range subnets {
		allSubnets[*subnets[i].AvailabilityDomain] = append(allSubnets[*subnets[i].AvailabilityDomain], *subnets[i].Id)
	}

	return allSubnets
}

func vcnsIDs(vcns []core.Vcn) []string {
	ids := make([]string, 0, len(vcns))

	for i := range vcns {
		ids = append(ids, *vcns[i].Id)
	}

	return ids
}
