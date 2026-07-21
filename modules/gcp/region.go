package gcp

import (
	"context"
	"os"
	"strings"

	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/compute/v1"
)

// You can set this environment variable to force Terratest to use a specific Region rather than a random one. This is
// convenient when iterating locally.
const regionOverrideEnvVarName = "TERRATEST_GCP_REGION"

// You can set this environment variable to force Terratest to use a specific Zone rather than a random one. This is
// convenient when iterating locally.
const zoneOverrideEnvVarName = "TERRATEST_GCP_ZONE"

// GetRandomRegionContext gets a randomly chosen GCP Region. If approvedRegions is not empty, this will be a Region from
// the approvedRegions list; otherwise, this method will fetch the latest list of regions from the GCP APIs and pick one
// of those. If forbiddenRegions is not empty, this method will make sure the returned Region is not in the
// forbiddenRegions list. This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRandomRegionContext(t testing.TestingT, ctx context.Context, projectID string, approvedRegions []string, forbiddenRegions []string) string {
	region, err := GetRandomRegionContextE(t, ctx, projectID, approvedRegions, forbiddenRegions)
	require.NoError(t, err)

	return region
}

// GetRandomRegionContextE gets a randomly chosen GCP Region. If approvedRegions is not empty, this will be a Region
// from the approvedRegions list; otherwise, this method will fetch the latest list of regions from the GCP APIs and pick
// one of those. If forbiddenRegions is not empty, this method will make sure the returned Region is not in the
// forbiddenRegions list. The ctx parameter supports cancellation and timeouts.
func GetRandomRegionContextE(t testing.TestingT, ctx context.Context, projectID string, approvedRegions []string, forbiddenRegions []string) (string, error) {
	regionFromEnvVar := os.Getenv(regionOverrideEnvVarName)
	if regionFromEnvVar != "" {
		logger.Default.Logf(t, "Using GCP Region %s from environment variable %s", regionFromEnvVar, regionOverrideEnvVarName)

		return regionFromEnvVar, nil
	}

	regionsToPickFrom := approvedRegions

	if len(regionsToPickFrom) == 0 {
		allRegions, err := GetAllGCPRegionsContextE(t, ctx, projectID)
		if err != nil {
			return "", err
		}

		regionsToPickFrom = allRegions
	}

	regionsToPickFrom = subtract(regionsToPickFrom, forbiddenRegions)
	region := random.RandomString(regionsToPickFrom)

	logger.Default.Logf(t, "Using Region %s", region)

	return region, nil
}

// GetRandomZoneContext gets a randomly chosen GCP Zone. If approvedZones is not empty, this will be a Zone from the
// approvedZones list; otherwise, this method will fetch the latest list of Zones from the GCP APIs and pick one of
// those. If forbiddenZones is not empty, this method will make sure the returned Zone is not in the forbiddenZones list.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRandomZoneContext(t testing.TestingT, ctx context.Context, projectID string, approvedZones []string, forbiddenZones []string, forbiddenRegions []string) string {
	zone, err := GetRandomZoneContextE(t, ctx, projectID, approvedZones, forbiddenZones, forbiddenRegions)
	require.NoError(t, err)

	return zone
}

// GetRandomZoneContextE gets a randomly chosen GCP Zone. If approvedZones is not empty, this will be a Zone from the
// approvedZones list; otherwise, this method will fetch the latest list of Zones from the GCP APIs and pick one of
// those. If forbiddenZones is not empty, this method will make sure the returned Zone is not in the forbiddenZones list.
// The ctx parameter supports cancellation and timeouts.
func GetRandomZoneContextE(t testing.TestingT, ctx context.Context, projectID string, approvedZones []string, forbiddenZones []string, forbiddenRegions []string) (string, error) {
	zoneFromEnvVar := os.Getenv(zoneOverrideEnvVarName)
	if zoneFromEnvVar != "" {
		logger.Default.Logf(t, "Using GCP Zone %s from environment variable %s", zoneFromEnvVar, zoneOverrideEnvVarName)

		return zoneFromEnvVar, nil
	}

	zonesToPickFrom := approvedZones

	if len(zonesToPickFrom) == 0 {
		allZones, err := GetAllGCPZonesContextE(t, ctx, projectID)
		if err != nil {
			return "", err
		}

		zonesToPickFrom = allZones
	}

	zonesToPickFrom = subtract(zonesToPickFrom, forbiddenZones)

	var zonesToPickFromFiltered []string

	for _, zone := range zonesToPickFrom {
		if !isInRegions(zone, forbiddenRegions) {
			zonesToPickFromFiltered = append(zonesToPickFromFiltered, zone)
		}
	}

	zone := random.RandomString(zonesToPickFromFiltered)

	return zone, nil
}

// GetRandomZoneForRegionContext gets a randomly chosen GCP Zone in the given Region.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRandomZoneForRegionContext(t testing.TestingT, ctx context.Context, projectID string, region string) string {
	zone, err := GetRandomZoneForRegionContextE(t, ctx, projectID, region)
	require.NoError(t, err)

	return zone
}

// GetRandomZoneForRegionContextE gets a randomly chosen GCP Zone in the given Region.
// The ctx parameter supports cancellation and timeouts.
func GetRandomZoneForRegionContextE(t testing.TestingT, ctx context.Context, projectID string, region string) (string, error) {
	zoneFromEnvVar := os.Getenv(zoneOverrideEnvVarName)
	if zoneFromEnvVar != "" {
		logger.Default.Logf(t, "Using GCP Zone %s from environment variable %s", zoneFromEnvVar, zoneOverrideEnvVarName)

		return zoneFromEnvVar, nil
	}

	allZones, err := GetAllGCPZonesContextE(t, ctx, projectID)
	if err != nil {
		return "", err
	}

	zonesToPickFrom := []string{}

	for _, zone := range allZones {
		if strings.Contains(zone, region) {
			zonesToPickFrom = append(zonesToPickFrom, zone)
		}
	}

	zone := random.RandomString(zonesToPickFrom)

	logger.Default.Logf(t, "Using Zone %s", zone)

	return zone, nil
}

// GetAllGCPRegionsContext gets the list of GCP regions available in this account.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAllGCPRegionsContext(t testing.TestingT, ctx context.Context, projectID string) []string {
	out, err := GetAllGCPRegionsContextE(t, ctx, projectID)
	require.NoError(t, err)

	return out
}

// GetAllGCPRegionsContextE gets the list of GCP regions available in this account.
// The ctx parameter supports cancellation and timeouts.
func GetAllGCPRegionsContextE(t testing.TestingT, ctx context.Context, projectID string) ([]string, error) {
	logger.Default.Logf(t, "Looking up all GCP regions available in this account")

	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAllGCPRegionsWithClient(ctx, service, projectID)
}

// GetAllGCPRegionsWithClient gets the list of GCP regions available in this account using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see region_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAllGCPRegionsWithClient(ctx context.Context, service *compute.Service, projectID string) ([]string, error) {
	req := service.Regions.List(projectID)
	regions := []string{}

	err := req.Pages(ctx, func(page *compute.RegionList) error {
		for _, region := range page.Items {
			regions = append(regions, region.Name)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return regions, nil
}

// GetAllGCPZonesContext gets the list of GCP Zones available in this account.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAllGCPZonesContext(t testing.TestingT, ctx context.Context, projectID string) []string {
	out, err := GetAllGCPZonesContextE(t, ctx, projectID)
	require.NoError(t, err)

	return out
}

// GetAllGCPZonesContextE gets the list of GCP Zones available in this account.
// The ctx parameter supports cancellation and timeouts.
func GetAllGCPZonesContextE(t testing.TestingT, ctx context.Context, projectID string) ([]string, error) {
	service, err := NewComputeServiceContextE(t, ctx)
	if err != nil {
		return nil, err
	}

	return GetAllGCPZonesWithClient(ctx, service, projectID)
}

// GetAllGCPZonesWithClient gets the list of GCP Zones available in this account using the supplied
// *compute.Service. Prefer this variant in unit tests where the service is backed by an httptest fake server
// (see region_unit_test.go for the pattern).
// The ctx parameter supports cancellation and timeouts.
func GetAllGCPZonesWithClient(ctx context.Context, service *compute.Service, projectID string) ([]string, error) {
	req := service.Zones.List(projectID)
	zones := []string{}

	err := req.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			zones = append(zones, zone.Name)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return zones, nil
}

// ZoneURLToZone extracts the zone name from a GCP Zone URL formatted like
// https://www.googleapis.com/compute/v1/projects/project-123456/zones/asia-east1-b and returns "asia-east1-b".
func ZoneURLToZone(zoneURL string) string {
	tokens := strings.Split(zoneURL, "/")

	return tokens[len(tokens)-1]
}

// RegionURLToRegion extracts the region name from a GCP Region URL formatted like
// https://www.googleapis.com/compute/v1/projects/project-123456/regions/southamerica-east1 and returns
// "southamerica-east1".
func RegionURLToRegion(regionURL string) string {
	tokens := strings.Split(regionURL, "/")

	return tokens[len(tokens)-1]
}

// isInRegions returns true if the given zone is in any of the given regions.
func isInRegions(zone string, regions []string) bool {
	for _, region := range regions {
		if isInRegion(zone, region) {
			return true
		}
	}

	return false
}

// isInRegion returns true if the given zone is in the given region.
func isInRegion(zone string, region string) bool {
	return strings.Contains(zone, region)
}

// subtract returns the items in list1 that are not in list2.
func subtract[T comparable](list1, list2 []T) []T {
	lookups := make(map[T]struct{}, len(list2))
	for _, item := range list2 {
		lookups[item] = struct{}{}
	}

	out := make([]T, 0, len(list1))
	for _, item := range list1 {
		if _, found := lookups[item]; !found {
			out = append(out, item)
		}
	}

	return out
}
