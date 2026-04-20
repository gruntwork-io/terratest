package azure

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
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

// GetRandomStableRegionContext gets a randomly chosen Azure region that is considered stable.
// Like GetRandomRegionContext, you can further restrict the stable region list using approvedRegions
// and forbiddenRegions. We consider stable regions to be those that have been around for at least 1 year.
// Note that regions in the approvedRegions list that are not considered stable are ignored.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRandomStableRegionContext(t testing.TestingT, ctx context.Context, approvedRegions []string, forbiddenRegions []string, subscriptionID string) string {
	t.Helper()

	regionsToPickFrom := stableRegions

	if len(approvedRegions) > 0 {
		regionsToPickFrom = collections.ListIntersection(regionsToPickFrom, approvedRegions)
	}

	if len(forbiddenRegions) > 0 {
		regionsToPickFrom = collections.ListSubtract(regionsToPickFrom, forbiddenRegions)
	}

	return GetRandomRegionContext(t, ctx, regionsToPickFrom, nil, subscriptionID) //nolint:staticcheck
}

// GetRandomStableRegion gets a randomly chosen Azure region that is considered stable. Like GetRandomRegion, you can
// further restrict the stable region list using approvedRegions and forbiddenRegions. We consider stable regions to be
// those that have been around for at least 1 year.
// Note that regions in the approvedRegions list that are not considered stable are ignored.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetRandomStableRegionContext] instead.
func GetRandomStableRegion(t testing.TestingT, approvedRegions []string, forbiddenRegions []string, subscriptionID string) string {
	t.Helper()

	return GetRandomStableRegionContext(t, context.Background(), approvedRegions, forbiddenRegions, subscriptionID)
}

// GetRandomRegionContext gets a randomly chosen Azure region.
// If approvedRegions is not empty, this will be a region from the approvedRegions list; otherwise,
// this method will fetch the latest list of regions from the Azure APIs and pick one of those.
// If forbiddenRegions is not empty, this method will make sure the returned region is not in the forbiddenRegions list.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRandomRegionContext(t testing.TestingT, ctx context.Context, approvedRegions []string, forbiddenRegions []string, subscriptionID string) string {
	t.Helper()

	region, err := GetRandomRegionContextE(t, ctx, approvedRegions, forbiddenRegions, subscriptionID)
	require.NoError(t, err)

	return region
}

// GetRandomRegion gets a randomly chosen Azure region. If approvedRegions is not empty, this will be a region from the approvedRegions
// list; otherwise, this method will fetch the latest list of regions from the Azure APIs and pick one of those. If
// forbiddenRegions is not empty, this method will make sure the returned region is not in the forbiddenRegions list.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetRandomRegionContext] instead.
func GetRandomRegion(t testing.TestingT, approvedRegions []string, forbiddenRegions []string, subscriptionID string) string {
	t.Helper()

	return GetRandomRegionContext(t, context.Background(), approvedRegions, forbiddenRegions, subscriptionID) //nolint:staticcheck
}

// GetRandomRegionContextE gets a randomly chosen Azure region.
// If approvedRegions is not empty, this will be a region from the approvedRegions list; otherwise,
// this method will fetch the latest list of regions from the Azure APIs and pick one of those.
// If forbiddenRegions is not empty, this method will make sure the returned region is not in the forbiddenRegions list.
// The ctx parameter supports cancellation and timeouts.
func GetRandomRegionContextE(t testing.TestingT, ctx context.Context, approvedRegions []string, forbiddenRegions []string, subscriptionID string) (string, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return "", err
	}

	regionsToPickFrom := approvedRegions

	if len(regionsToPickFrom) == 0 {
		allRegions, err := GetAllAzureRegionsContextE(t, ctx, subscriptionID)
		if err != nil {
			return "", err
		}

		regionsToPickFrom = allRegions
	}

	regionsToPickFrom = collections.ListSubtract(regionsToPickFrom, forbiddenRegions)
	region := random.RandomString(regionsToPickFrom)

	return region, nil
}

// GetRandomRegionE gets a randomly chosen Azure region. If approvedRegions is not empty, this will be a region from the approvedRegions
// list; otherwise, this method will fetch the latest list of regions from the Azure APIs and pick one of those. If
// forbiddenRegions is not empty, this method will make sure the returned region is not in the forbiddenRegions list.
//
// Deprecated: Use [GetRandomRegionContextE] instead.
func GetRandomRegionE(t testing.TestingT, approvedRegions []string, forbiddenRegions []string, subscriptionID string) (string, error) {
	return GetRandomRegionContextE(t, context.Background(), approvedRegions, forbiddenRegions, subscriptionID)
}

// GetAllAzureRegionsContext gets the list of Azure regions available in this subscription.
// This function would fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAllAzureRegionsContext(t testing.TestingT, ctx context.Context, subscriptionID string) []string {
	t.Helper()

	out, err := GetAllAzureRegionsContextE(t, ctx, subscriptionID)
	require.NoError(t, err)

	return out
}

// GetAllAzureRegions gets the list of Azure regions available in this subscription.
// This function would fail the test if there is an error.
//
// Deprecated: Use [GetAllAzureRegionsContext] instead.
func GetAllAzureRegions(t testing.TestingT, subscriptionID string) []string {
	t.Helper()

	return GetAllAzureRegionsContext(t, context.Background(), subscriptionID)
}

// GetAllAzureRegionsContextE gets the list of Azure regions available in this subscription.
// The ctx parameter supports cancellation and timeouts.
func GetAllAzureRegionsContextE(t testing.TestingT, ctx context.Context, subscriptionID string) ([]string, error) {
	// Validate Azure subscription ID
	subscriptionID, err := getTargetAzureSubscription(subscriptionID)
	if err != nil {
		return nil, err
	}

	// Setup Subscription client
	subscriptionClient, err := CreateSubscriptionsClientContextE(ctx)
	if err != nil {
		return nil, err
	}

	// Get list of Azure locations via pager
	var regions []string

	pager := subscriptionClient.NewListLocationsPager(subscriptionID, nil)
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, region := range page.Value {
			if region == nil || region.Name == nil {
				continue
			}

			regions = append(regions, *region.Name)
		}
	}

	return regions, nil
}

// GetAllAzureRegionsE gets the list of Azure regions available in this subscription.
//
// Deprecated: Use [GetAllAzureRegionsContextE] instead.
func GetAllAzureRegionsE(t testing.TestingT, subscriptionID string) ([]string, error) {
	return GetAllAzureRegionsContextE(t, context.Background(), subscriptionID)
}
