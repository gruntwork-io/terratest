package test

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestPostgreSQLDatabase(t *testing.T) {
	t.Parallel()

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../examples/azure/terraform-azure-postgresql-example",
		NoColor:      true,
	})
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)
	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")

	// get the actual data via terraform output
	expectedServername := "terratestdevpostgresqlsrvfoo" // must match fixture
	actualServername := terraform.Output(t, terraformOptions, "servername")
	rgName := terraform.Output(t, terraformOptions, "rgname")
	expectedSkuName := terraform.Output(t, terraformOptions, "sku_name")
	actualServer := GetPostgresqlServer(t, rgName, actualServername, subscriptionID)
	// Verify
	assert.Equal(t, expectedServername, actualServername)
	assert.Equal(t, expectedSkuName, *actualServer.Sku.Name)

}
