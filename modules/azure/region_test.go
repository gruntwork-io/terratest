//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure //nolint:testpackage // tests access unexported functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRandomRegion(t *testing.T) {
	t.Parallel()

	randomRegion := GetRandomRegionContext(t, t.Context(), nil, nil, "")
	assertLooksLikeRegionName(t, randomRegion)
}

func TestGetRandomRegionExcludesForbiddenRegions(t *testing.T) {
	t.Parallel()

	approvedRegions := []string{"canadacentral", "eastus", "eastus2", "westus", "westus2", "westeurope", "northeurope", "uksouth", "southeastasia", "eastasia", "japaneast", "australiacentral"}
	forbiddenRegions := []string{"westus2", "japaneast"}

	for i := 0; i < 48; i++ {
		randomRegion := GetRandomRegionContext(t, t.Context(), approvedRegions, forbiddenRegions, "")
		assert.NotContains(t, forbiddenRegions, randomRegion)
	}
}

func TestGetAllAzureRegions(t *testing.T) {
	t.Parallel()

	regions := GetAllAzureRegionsContext(t, t.Context(), "")

	// The typical subscription had access to 30+ live regions as of
	// July 2019: https://azure.microsoft.com/en-us/global-infrastructure/regions/
	assert.GreaterOrEqual(t, len(regions), 30, "Number of regions: %d", len(regions))

	for _, region := range regions {
		assertLooksLikeRegionName(t, region)
	}
}

func assertLooksLikeRegionName(t *testing.T, regionName string) {
	t.Helper()

	assert.Regexp(t, "[a-z]", regionName)
}
