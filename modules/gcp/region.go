package gcp

import (
	"context"
	"os"
	"strings"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/testing"
	"google.golang.org/api/compute/v1"
)

// You can set this environment variable to force Terratest to use a specific Region rather than a random one. This is
// convenient when iterating locally.
const regionOverrideEnvVarName = "TERRATEST_GCP_REGION"

// You can set this environment variable to force Terratest to use a specific Zone rather than a random one. This is
// convenient when iterating locally.
const zoneOverrideEnvVarName = "TERRATEST_GCP_ZONE"

// Some GCP API calls require a GCP Region. We typically require the user to set one explicitly, but in some
// cases, this doesn't make sense (e.g., for fetching the list of regions in an account), so for those cases, we use
// this Region as a default.
const defaultRegion = "us-west1"

// Some GCP API calls require a GCP Zone. We typically require the user to set one explicitly, but in some
// cases, this doesn't make sense (e.g., for fetching the list of regions in an account), so for those cases, we use
// this Zone as a default.
const defaultZone = "us-west1-b"

// GetRandomRegion gets a randomly chosen GCP Region. If approvedRegions is not empty, this will be a Region from the approvedRegions
// list; otherwise, this method will fetch the latest list of regions from the GCP APIs and pick one of those. If
// forbiddenRegions is not empty, this method will make sure the returned Region is not in the forbiddenRegions list.
func GetRandomRegion(t testing.TestingT, projectID string, approvedRegions []string, forbiddenRegions []string) string {
	region, err := GetRandomRegionE(t, projectID, approvedRegions, forbiddenRegions)
	if err != nil {
		t.Fatal(err)
	}
	return region
}

// GetRandomRegionE gets a randomly chosen GCP Region. If approvedRegions is not empty, this will be a Region from the approvedRegions
// list; otherwise, this method will fetch the latest list of regions from the GCP APIs and pick one of those. If
// forbiddenRegions is not empty, this method will make sure the returned Region is not in the forbiddenRegions list.
func GetRandomRegionE(t testing.TestingT, projectID string, approvedRegions []string, forbiddenRegions []string) (string, error) {
	regionFromEnvVar := os.Getenv(regionOverrideEnvVarName)
	if regionFromEnvVar != "" {
		logger.Logf(t, "Using GCP Region %s from environment variable %s", regionFromEnvVar, regionOverrideEnvVarName)
		return regionFromEnvVar, nil
	}

	regionsToPickFrom := approvedRegions

	if len(regionsToPickFrom) == 0 {
		allRegions, err := GetAllGcpRegionsE(t, projectID)
		if err != nil {
			return "", err
		}
		regionsToPickFrom = allRegions
	}

	regionsToPickFrom = collections.ListSubtract(regionsToPickFrom, forbiddenRegions)
	region := random.RandomString(regionsToPickFrom)

	logger.Logf(t, "Using Region %s", region)
	return region, nil
}

// GetRandomZone gets a randomly chosen GCP Zone. If approvedRegions is not empty, this will be a Zone from the approvedZones
// list; otherwise, this method will fetch the latest list of Zones from the GCP APIs and pick one of those. If
// forbiddenZones is not empty, this method will make sure the returned Region is not in the forbiddenZones list.
func GetRandomZone(t testing.TestingT, projectID string, approvedZones []string, forbiddenZones []string, forbiddenRegions []string) string {
	zone, err := GetRandomZoneE(t, projectID, approvedZones, forbiddenZones, forbiddenRegions)
	if err != nil {
		t.Fatal(err)
	}
	return zone
}

// GetRandomZoneE gets a randomly chosen GCP Zone. If approvedRegions is not empty, this will be a Zone from the approvedZones
// list; otherwise, this method will fetch the latest list of Zones from the GCP APIs and pick one of those. If
// forbiddenZones is not empty, this method will make sure the returned Region is not in the forbiddenZones list.
func GetRandomZoneE(t testing.TestingT, projectID string, approvedZones []string, forbiddenZones []string, forbiddenRegions []string) (string, error) {
	zoneFromEnvVar := os.Getenv(zoneOverrideEnvVarName)
	if zoneFromEnvVar != "" {
		logger.Logf(t, "Using GCP Zone %s from environment variable %s", zoneFromEnvVar, zoneOverrideEnvVarName)
		return zoneFromEnvVar, nil
	}

	zonesToPickFrom := approvedZones

	if len(zonesToPickFrom) == 0 {
		allZones, err := GetAllGcpZonesE(t, projectID)
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
func GetRandomZoneForRegion(t testing.TestingT, projectID string, region string) string {
	zone, err := GetRandomZoneForRegionE(t, projectID, region)
	if err != nil {
		t.Fatal(err)
	}
	return zone
}

// GetRandomZoneForRegionE gets a randomly chosen GCP Zone in the given Region.
func GetRandomZoneForRegionE(t testing.TestingT, projectID string, region string) (string, error) {
	zoneFromEnvVar := os.Getenv(zoneOverrideEnvVarName)
	if zoneFromEnvVar != "" {
		logger.Logf(t, "Using GCP Zone %s from environment variable %s", zoneFromEnvVar, zoneOverrideEnvVarName)
		return zoneFromEnvVar, nil
	}

	allZones, err := GetAllGcpZonesE(t, projectID)
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

	logger.Logf(t, "Using Zone %s", zone)
	return zone, nil
}

// GetAllGcpRegions gets the list of GCP regions available in this account.
func GetAllGcpRegions(t testing.TestingT, projectID string) []string {
	out, err := GetAllGcpRegionsE(t, projectID)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetAllGcpRegionsE gets the list of GCP regions available in this account.
func GetAllGcpRegionsE(t testing.TestingT, projectID string) ([]string, error) {
	logger.Log(t, "Looking up all GCP regions available in this account")

	// Note that NewComputeServiceE creates a context, but it appears to be empty so we keep the code simpler by
	// creating a new one here
	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	req := service.Regions.List(projectID)

	regions := []string{}
	err = req.Pages(ctx, func(page *compute.RegionList) error {
		for _, region := range page.Items {
			regions = append(regions, region.Name)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return regions, nil
}

// GetAllGcpZones gets the list of GCP Zones available in this account.
func GetAllGcpZones(t testing.TestingT, projectID string) []string {
	out, err := GetAllGcpZonesE(t, projectID)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetAllGcpZonesE gets the list of GCP Zones available in this account.
func GetAllGcpZonesE(t testing.TestingT, projectID string) ([]string, error) {
	// Note that NewComputeServiceE creates a context, but it appears to be empty so we keep the code simpler by
	// creating a new one here
	ctx := context.Background()

	service, err := NewComputeServiceE(t)
	if err != nil {
		return nil, err
	}

	req := service.Zones.List(projectID)

	zones := []string{}
	err = req.Pages(ctx, func(page *compute.ZoneList) error {
		for _, zone := range page.Items {
			zones = append(zones, zone.Name)
		}
		return err
	})
	if err != nil {
		return nil, err
	}

	return zones, nil
}

// Given a GCP Zone URL formatted like https://www.googleapis.com/compute/v1/projects/project-123456/zones/asia-east1-b,
// return "asia-east1-b".
// Todo: Improve sanity checking on this function by using a RegEx with capture groups
func ZoneUrlToZone(zoneUrl string) string {
	tokens := strings.Split(zoneUrl, "/")
	return tokens[len(tokens)-1]
}

// Given a GCP Zone URL formatted like https://www.googleapis.com/compute/v1/projects/project-123456/regions/southamerica-east1,
// return "southamerica-east1".
// Todo: Improve sanity checking on this function by using a RegEx with capture groups
func RegionUrlToRegion(zoneUrl string) string {
	tokens := strings.Split(zoneUrl, "/")
	return tokens[len(tokens)-1]
}

// Returns true if the given zone is in any of the given regions
func isInRegions(zone string, regions []string) bool {
	for _, region := range regions {
		if isInRegion(zone, region) {
			return true
		}
	}

	return false
}

// Returns true if the given zone is in the given region
func isInRegion(zone string, region string) bool {
	return strings.Contains(zone, region)
}
