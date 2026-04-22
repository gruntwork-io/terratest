package gcp

import (
	"context"
	"os"
	"strings"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/compute/v1"
)

// You can set this environment variable to force Terratest to use a specific Region rather than a random one. This is
// convenient when iterating locally.
const regionOverrideEnvVarName = "TERRATEST_GCP_REGION"

// You can set this environment variable to force Terratest to use a specific Zone rather than a random one. This is
// convenient when iterating locally.
const zoneOverrideEnvVarName = "TERRATEST_GCP_ZONE"

// GetRandomRegion gets a randomly chosen GCP Region. If approvedRegions is not empty, this will be a Region from the
// approvedRegions list; otherwise, this method will fetch the latest list of regions from the GCP APIs and pick one of
// those. If forbiddenRegions is not empty, this method will make sure the returned Region is not in the forbiddenRegions
// list. This will fail the test if there is an error.
//
// Deprecated: Use [GetRandomRegionContext] instead.
func GetRandomRegion(t testing.TestingT, projectID string, approvedRegions []string, forbiddenRegions []string) string {
	return GetRandomRegionContext(t, context.Background(), projectID, approvedRegions, forbiddenRegions)
}

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

// GetRandomRegionE gets a randomly chosen GCP Region. If approvedRegions is not empty, this will be a Region from the
// approvedRegions list; otherwise, this method will fetch the latest list of regions from the GCP APIs and pick one of
// those. If forbiddenRegions is not empty, this method will make sure the returned Region is not in the forbiddenRegions
// list.
//
// Deprecated: Use [GetRandomRegionContextE] instead.
func GetRandomRegionE(t testing.TestingT, projectID string, approvedRegions []string, forbiddenRegions []string) (string, error) {
	return GetRandomRegionContextE(t, context.Background(), projectID, approvedRegions, forbiddenRegions)
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

	regionsToPickFrom = collections.ListSubtract(regionsToPickFrom, forbiddenRegions)
	region := random.RandomString(regionsToPickFrom)

	logger.Default.Logf(t, "Using Region %s", region)

	return region, nil
}

// GetRandomZone gets a randomly chosen GCP Zone. If approvedZones is not empty, this will be a Zone from the
// approvedZones list; otherwise, this method will fetch the latest list of Zones from the GCP APIs and pick one of
// those. If forbiddenZones is not empty, this method will make sure the returned Zone is not in the forbiddenZones list.
// This will fail the test if there is an error.
//
// Deprecated: Use [GetRandomZoneContext] instead.
func GetRandomZone(t testing.TestingT, projectID string, approvedZones []string, forbiddenZones []string, forbiddenRegions []string) string {
	return GetRandomZoneContext(t, context.Background(), projectID, approvedZones, forbiddenZones, forbiddenRegions)
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

// GetRandomZoneE gets a randomly chosen GCP Zone. If approvedZones is not empty, this will be a Zone from the
// approvedZones list; otherwise, this method will fetch the latest list of Zones from the GCP APIs and pick one of
// those. If forbiddenZones is not empty, this method will make sure the returned Zone is not in the forbiddenZones list.
//
// Deprecated: Use [GetRandomZoneContextE] instead.
func GetRandomZoneE(t testing.TestingT, projectID string, approvedZones []string, forbiddenZones []string, forbiddenRegions []string) (string, error) {
	return GetRandomZoneContextE(t, context.Background(), projectID, approvedZones, forbiddenZones, forbiddenRegions)
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

	zonesToPickFrom = collections.ListSubtract(zonesToPickFrom, forbiddenZones)

	var zonesToPickFromFiltered []string

	for _, zone := range zonesToPickFrom {
		if !isInRegions(zone, forbiddenRegions) {
			zonesToPickFromFiltered = append(zonesToPickFromFiltered, zone)
		}
	}

	zone := random.RandomString(zonesToPickFromFiltered)

	return zone, nil
}

// GetRandomZoneForRegion gets a randomly chosen GCP Zone in the given Region.
// This will fail the test if there is an error.
//
// Deprecated: Use [GetRandomZoneForRegionContext] instead.
func GetRandomZoneForRegion(t testing.TestingT, projectID string, region string) string {
	return GetRandomZoneForRegionContext(t, context.Background(), projectID, region)
}

// GetRandomZoneForRegionContext gets a randomly chosen GCP Zone in the given Region.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRandomZoneForRegionContext(t testing.TestingT, ctx context.Context, projectID string, region string) string {
	zone, err := GetRandomZoneForRegionContextE(t, ctx, projectID, region)
	require.NoError(t, err)

	return zone
}

// GetRandomZoneForRegionE gets a randomly chosen GCP Zone in the given Region.
//
// Deprecated: Use [GetRandomZoneForRegionContextE] instead.
func GetRandomZoneForRegionE(t testing.TestingT, projectID string, region string) (string, error) {
	return GetRandomZoneForRegionContextE(t, context.Background(), projectID, region)
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

// GetAllGCPRegions gets the list of GCP regions available in this account.
// This will fail the test if there is an error.
//
// Deprecated: Use [GetAllGCPRegionsContext] instead.
func GetAllGCPRegions(t testing.TestingT, projectID string) []string {
	return GetAllGCPRegionsContext(t, context.Background(), projectID)
}

// GetAllGCPRegionsContext gets the list of GCP regions available in this account.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAllGCPRegionsContext(t testing.TestingT, ctx context.Context, projectID string) []string {
	out, err := GetAllGCPRegionsContextE(t, ctx, projectID)
	require.NoError(t, err)

	return out
}

// GetAllGCPRegionsE gets the list of GCP regions available in this account.
//
// Deprecated: Use [GetAllGCPRegionsContextE] instead.
func GetAllGCPRegionsE(t testing.TestingT, projectID string) ([]string, error) {
	return GetAllGCPRegionsContextE(t, context.Background(), projectID)
}

// GetAllGcpRegions gets the list of GCP regions available in this account.
// This will fail the test if there is an error.
//
// Deprecated: Use [GetAllGCPRegionsContext] instead.
func GetAllGcpRegions(t testing.TestingT, projectID string) []string { //nolint:staticcheck,revive // preserving deprecated function name
	return GetAllGCPRegionsContext(t, context.Background(), projectID)
}

// GetAllGcpRegionsContext gets the list of GCP regions available in this account.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
//
// Deprecated: Use [GetAllGCPRegionsContext] instead.
func GetAllGcpRegionsContext(t testing.TestingT, ctx context.Context, projectID string) []string { //nolint:staticcheck,revive // preserving deprecated function name
	return GetAllGCPRegionsContext(t, ctx, projectID)
}

// GetAllGcpRegionsE gets the list of GCP regions available in this account.
//
// Deprecated: Use [GetAllGCPRegionsContextE] instead.
func GetAllGcpRegionsE(t testing.TestingT, projectID string) ([]string, error) { //nolint:staticcheck,revive // preserving deprecated function name
	return GetAllGCPRegionsContextE(t, context.Background(), projectID)
}

// GetAllGcpRegionsContextE gets the list of GCP regions available in this account.
// The ctx parameter supports cancellation and timeouts.
//
// Deprecated: Use [GetAllGCPRegionsContextE] instead.
func GetAllGcpRegionsContextE(t testing.TestingT, ctx context.Context, projectID string) ([]string, error) { //nolint:staticcheck,revive // preserving deprecated function name
	return GetAllGCPRegionsContextE(t, ctx, projectID)
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

// GetAllGCPZones gets the list of GCP Zones available in this account.
// This will fail the test if there is an error.
//
// Deprecated: Use [GetAllGCPZonesContext] instead.
func GetAllGCPZones(t testing.TestingT, projectID string) []string {
	return GetAllGCPZonesContext(t, context.Background(), projectID)
}

// GetAllGCPZonesContext gets the list of GCP Zones available in this account.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAllGCPZonesContext(t testing.TestingT, ctx context.Context, projectID string) []string {
	out, err := GetAllGCPZonesContextE(t, ctx, projectID)
	require.NoError(t, err)

	return out
}

// GetAllGCPZonesE gets the list of GCP Zones available in this account.
//
// Deprecated: Use [GetAllGCPZonesContextE] instead.
func GetAllGCPZonesE(t testing.TestingT, projectID string) ([]string, error) {
	return GetAllGCPZonesContextE(t, context.Background(), projectID)
}

// GetAllGcpZones gets the list of GCP Zones available in this account.
// This will fail the test if there is an error.
//
// Deprecated: Use [GetAllGCPZonesContext] instead.
func GetAllGcpZones(t testing.TestingT, projectID string) []string { //nolint:staticcheck,revive // preserving deprecated function name
	return GetAllGCPZonesContext(t, context.Background(), projectID)
}

// GetAllGcpZonesContext gets the list of GCP Zones available in this account.
// This will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
//
// Deprecated: Use [GetAllGCPZonesContext] instead.
func GetAllGcpZonesContext(t testing.TestingT, ctx context.Context, projectID string) []string { //nolint:staticcheck,revive // preserving deprecated function name
	return GetAllGCPZonesContext(t, ctx, projectID)
}

// GetAllGcpZonesE gets the list of GCP Zones available in this account.
//
// Deprecated: Use [GetAllGCPZonesContextE] instead.
func GetAllGcpZonesE(t testing.TestingT, projectID string) ([]string, error) { //nolint:staticcheck,revive // preserving deprecated function name
	return GetAllGCPZonesContextE(t, context.Background(), projectID)
}

// GetAllGcpZonesContextE gets the list of GCP Zones available in this account.
// The ctx parameter supports cancellation and timeouts.
//
// Deprecated: Use [GetAllGCPZonesContextE] instead.
func GetAllGcpZonesContextE(t testing.TestingT, ctx context.Context, projectID string) ([]string, error) { //nolint:staticcheck,revive // preserving deprecated function name
	return GetAllGCPZonesContextE(t, ctx, projectID)
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

// ZoneUrlToZone extracts the zone name from a GCP Zone URL.
//
// Deprecated: Use [ZoneURLToZone] instead.
func ZoneUrlToZone(zoneURL string) string { //nolint:staticcheck,revive // preserving deprecated function name
	return ZoneURLToZone(zoneURL)
}

// RegionURLToRegion extracts the region name from a GCP Region URL formatted like
// https://www.googleapis.com/compute/v1/projects/project-123456/regions/southamerica-east1 and returns
// "southamerica-east1".
func RegionURLToRegion(regionURL string) string {
	tokens := strings.Split(regionURL, "/")

	return tokens[len(tokens)-1]
}

// RegionUrlToRegion extracts the region name from a GCP Region URL.
//
// Deprecated: Use [RegionURLToRegion] instead.
func RegionUrlToRegion(regionURL string) string { //nolint:staticcheck,revive // preserving deprecated function name
	return RegionURLToRegion(regionURL)
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
