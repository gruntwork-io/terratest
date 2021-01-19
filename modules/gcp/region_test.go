// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package gcp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomRegion(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)

	randomRegion := GetRandomRegion(t, projectID, nil, nil)
	assertLooksLikeRegionName(t, randomRegion)
}

func TestGetRandomZone(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)

	randomZone := GetRandomZone(t, projectID, nil, nil, nil)
	assertLooksLikeZoneName(t, randomZone)
}

func TestGetRandomRegionExcludesForbiddenRegions(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)

	approvedRegions := []string{"asia-east1", "asia-northeast1", "asia-south1", "asia-southeast1", "australia-southeast1", "europe-north1", "europe-west1", "europe-west2", "europe-west3", "northamerica-northeast1", "southamerica-east1", "us-central1", "us-east1", "us-east4", "us-west2"}
	forbiddenRegions := []string{"europe-west4", "us-west1"}

	for i := 0; i < 1000; i++ {
		randomRegion := GetRandomRegion(t, projectID, approvedRegions, forbiddenRegions)
		assert.NotContains(t, forbiddenRegions, randomRegion)
	}
}

func TestGetRandomZoneExcludesForbiddenZones(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)

	approvedZones := []string{"us-east1-b", "us-east1-c", "us-east1-d", "us-east4-a", "us-east4-b", "us-east4-c", "us-west2-a", "us-west2-b", "us-west2-c", "us-central1-f", "europe-west2-b"}
	forbiddenZones := []string{"us-east1-a", "europe-west1-a", "europe-west2-a", "europe-west2-c"}

	for i := 0; i < 1000; i++ {
		randomZone := GetRandomZone(t, projectID, approvedZones, forbiddenZones, nil)
		assert.NotContains(t, forbiddenZones, randomZone)
	}
}

func TestGetRandomZoneExcludesForbiddenRegions(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)

	approvedZones := []string{"us-east1-b", "us-east1-c", "us-east1-d", "us-east4-a", "us-east4-b", "us-east4-c", "us-west2-a", "us-west2-b", "us-west2-c", "us-central1-f", "europe-west2-b"}
	forbiddenRegions := []string{"europe-west2"}

	for i := 0; i < 1000; i++ {
		randomZone := GetRandomZone(t, projectID, approvedZones, nil, forbiddenRegions)

		for _, forbiddenRegion := range forbiddenRegions {
			assert.True(t, !isInRegion(randomZone, forbiddenRegion), "Expected that selected zone %s would not be in region %s, but it is.", randomZone, forbiddenRegion)
		}
	}
}

func TestGetAllGcpRegions(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)

	regions := GetAllGcpRegions(t, projectID)

	// The typical account had access to 17 regions as of August, 2018: https://cloud.google.com/compute/docs/regions-zones/
	assert.True(t, len(regions) >= 17, "Number of regions: %d", len(regions))
	for _, region := range regions {
		assertLooksLikeRegionName(t, region)
	}
}

func TestGetAllGcpZones(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)

	zones := GetAllGcpZones(t, projectID)

	// The typical account had access to 52 zones as of August, 2018: https://cloud.google.com/compute/docs/regions-zones/
	assert.True(t, len(zones) >= 52, "Number of zones: %d", len(zones))
	for _, zone := range zones {
		assertLooksLikeZoneName(t, zone)
	}
}

func TestGetRandomZoneForRegion(t *testing.T) {
	t.Parallel()

	projectID := GetGoogleProjectIDFromEnvVar(t)

	regions := []string{
		"us-west1",
		"us-west2",
		"us-central1",
	}

	for _, region := range regions {
		zone, err := GetRandomZoneForRegionE(t, projectID, region)
		if err != nil {
			t.Fatal(err)
		}

		assert.True(t, strings.Contains(zone, region), "Expected zone %s to be in region %s", zone, region)
	}
}

func TestGetInRegion(t *testing.T) {
	t.Parallel()

	testData := []struct {
		zone     string
		region   string
		expected bool
	}{
		{"us-west2a", "us-west2", true},
		{"us-west2b", "us-west2", true},
		{"us-west2a", "us-east1", false},
	}

	for _, td := range testData {
		actual := isInRegion(td.zone, td.region)
		assert.Equal(t, td.expected, actual, "Expected %t for isInRegion(%s, %s) but got %t", td.expected, td.zone, td.region, actual)
	}
}

func TestGetInRegions(t *testing.T) {
	t.Parallel()

	testData := []struct {
		zone     string
		regions  []string
		expected bool
	}{
		{"us-west2a", []string{"us-west2", "us-east1"}, true},
		{"us-west2b", []string{"us-west2", "us-east1"}, true},
		{"us-west2a", []string{"us-west2", "us-east1"}, true},
		{"us-west2a", []string{"us-east1", "europe-west1"}, false},
	}

	for _, td := range testData {
		actual := isInRegions(td.zone, td.regions)
		assert.Equal(t, td.expected, actual, "Expected %t for isInRegions(%s, %v) but got %t", td.expected, td.zone, td.regions, actual)
	}
}

func assertLooksLikeRegionName(t *testing.T, regionName string) {
	assert.Regexp(t, "[a-z]+-[a-z]+[[:digit:]]+", regionName)
}

func assertLooksLikeZoneName(t *testing.T, zoneName string) {
	assert.Regexp(t, "[a-z]+-[a-z]+[[:digit:]]+-[a-z]{1}", zoneName)
}
