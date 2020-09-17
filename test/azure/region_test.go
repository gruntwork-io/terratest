// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomRegion(t *testing.T) {
	t.Parallel()

	randomRegion := azure.GetRandomRegion(t, nil, nil, "")
	assertLooksLikeRegionName(t, randomRegion)
}

func TestGetRandomRegionExcludesForbiddenRegions(t *testing.T) {
	t.Parallel()

	approvedRegions := []string{"canadacentral", "eastus", "eastus2", "westus", "westus2", "westeurope", "northeurope", "uksouth", "southeastasia", "eastasia", "japaneast", "australiacentral"}
	forbiddenRegions := []string{"westus2", "japaneast"}

	for i := 0; i < 1000; i++ {
		randomRegion := azure.GetRandomRegion(t, approvedRegions, forbiddenRegions, "")
		assert.NotContains(t, forbiddenRegions, randomRegion)
	}
}

func TestGetAllAzureRegions(t *testing.T) {
	t.Parallel()

	regions := azure.GetAllAzureRegions(t, "")

	// The typical subscription had access to 30+ live regions as of July 2019: https://azure.microsoft.com/en-us/global-infrastructure/regions/
	assert.True(t, len(regions) >= 30, "Number of regions: %d", len(regions))
	for _, region := range regions {
		assertLooksLikeRegionName(t, region)
	}
}

func assertLooksLikeRegionName(t *testing.T, regionName string) {
	assert.Regexp(t, "[a-z]", regionName)
}
