// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"strings"

	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureAppServicePlanExample(t *testing.T) {
	t.Parallel()
	_random := strings.ToLower(random.UniqueId())

	expectedResourceGroupName := fmt.Sprintf("tmp-rg-%s", _random)
	expectedAppName := fmt.Sprintf("tmp-asp-%s", _random)
	expectedEnvironment := "dev"
	expectedLocation := "westus2"
	expectedSkuSize := "S1"
	expectedSkuTier := "Standard"
	var expectedSkuCapacity int32
	expectedSkuCapacity = 1
	expectedKind := "Windows"
	expectedReserved := false

	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/azure/terraform-azure-appserviceplan-example",
		Vars: map[string]interface{}{
			"resourceGroupName": expectedResourceGroupName,
			"appName":           expectedAppName,
			"environment":       expectedEnvironment,
			"location":          expectedLocation,
		},
	}
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	assert := assert.New(t)

	outputValue := terraform.Output(t, terraformOptions, "planids")
	assert.NotNil(outputValue)
	assert.Contains(outputValue, expectedAppName)

	plan := azure.GetAppServicePlan(t, expectedAppName, expectedResourceGroupName, "")

	assert.Equal(expectedSkuSize, *plan.Sku.Size)
	assert.Equal(expectedSkuTier, *plan.Sku.Tier)
	assert.Equal(expectedSkuCapacity, *plan.Sku.Capacity)
	assert.Equal(expectedKind, *plan.Kind)
	assert.Equal(expectedReserved, *plan.Reserved)
}
