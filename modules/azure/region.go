package azure

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// Reference for region list: https://azure.microsoft.com/en-us/global-infrastructure/locations/
var stableRegions = []string{
	// Americas
	"centralus",
	"eastus",
	"eastus2",
	"northcentralus",
	"southcentralus",
	"westcentralus",
	"westus",
	"westus2",
	"canadacentral",
	"canadaeast",
	"brazilsouth",

	// Europe
	"northeurope",
	"westeurope",
	"francecentral",
	"francesouth",
	"uksouth",
	"ukwest",
	// "germanycentral", // Shows as active on Azure website, but not from API
	// "germanynortheast", // Shows as active on Azure website, but not from API

	// Asia Pacific
	"eastasia",
	"southeastasia",
	"australiacentral",
	"australiacentral2",
	"australiaeast",
	"australiasoutheast",
	"chinaeast",
	"chinaeast2",
	"chinanorth",
	"chinanorth2",
	"centralindia",
	"southindia",
	"westindia",
	"japaneast",
	"japanwest",
	"koreacentral",
	"koreasouth",

	// Middle East and Africa
	"southafricanorth",
	"southafricawest",
	"uaecentral",
	"uaenorth",
}

// GetRandomStableRegion gets a randomly chosen Azure region that is considered stable. Like GetRandomRegion, you can
// further restrict the stable region list using approvedRegions and forbiddenRegions. We consider stable regions to be
// those that have been around for at least 1 year.
// Note that regions in the approvedRegions list that are not considered stable are ignored.
func GetRandomStableRegion(t testing.TestingT, approvedRegions []string, forbiddenRegions []string, subscriptionID string) string {
	regionsToPickFrom := stableRegions
	if len(approvedRegions) > 0 {
		regionsToPickFrom = collections.ListIntersection(regionsToPickFrom, approvedRegions)
	}
	if len(forbiddenRegions) > 0 {
		regionsToPickFrom = collections.ListSubtract(regionsToPickFrom, forbiddenRegions)
	}
	return GetRandomRegion(t, regionsToPickFrom, nil, subscriptionID)
}

// GetRandomRegion gets a randomly chosen Azure region. If approvedRegions is not empty, this will be a region from the approvedRegions
// list; otherwise, this method will fetch the latest list of regions from the Azure APIs and pick one of those. If
// forbiddenRegions is not empty, this method will make sure the returned region is not in the forbiddenRegions list.
func GetRandomRegion(t testing.TestingT, approvedRegions []string, forbiddenRegions []string, subscriptionID string) string {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		t.Fatal(err)
	}

	region, err := GetRandomRegionE(t, approvedRegions, forbiddenRegions, subscriptionID)
	if err != nil {
		t.Fatal(err)
	}
	return region
}

// GetRandomRegionE gets a randomly chosen Azure region. If approvedRegions is not empty, this will be a region from the approvedRegions
// list; otherwise, this method will fetch the latest list of regions from the Azure APIs and pick one of those. If
// forbiddenRegions is not empty, this method will make sure the returned region is not in the forbiddenRegions list
func GetRandomRegionE(t testing.TestingT, approvedRegions []string, forbiddenRegions []string, subscriptionID string) (string, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return "", err
	}

	regionsToPickFrom := approvedRegions

	if len(regionsToPickFrom) == 0 {
		allRegions, err := GetAllAzureRegionsE(t, subscriptionID)
		if err != nil {
			return "", err
		}
		regionsToPickFrom = allRegions
	}

	regionsToPickFrom = collections.ListSubtract(regionsToPickFrom, forbiddenRegions)
	region := random.RandomString(regionsToPickFrom)

	return region, nil
}

// GetAllAzureRegions gets the list of Azure regions available in this subscription.
func GetAllAzureRegions(t testing.TestingT, subscriptionID string) []string {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		t.Fatal(err)
	}

	// Get list of Azure locations
	out, err := GetAllAzureRegionsE(t, subscriptionID)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// GetAllAzureRegionsE gets the list of Azure regions available in this subscription
func GetAllAzureRegionsE(t testing.TestingT, subscriptionID string) ([]string, error) {

	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Setup Subscription client
	subscriptionClient, err := GetSubscriptionClientE()
	if err != nil {
		return nil, err
	}

	// Get list of Azure locations
	out, err := subscriptionClient.ListLocations(context.Background(), subscriptionID)
	if err != nil {
		return nil, err
	}

	// Populate a return slice
	regions := []string{}
	for _, region := range *out.Value {
		regions = append(regions, *region.Name)
	}

	return regions, nil
}
