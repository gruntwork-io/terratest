//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/mysql/mgmt/mysqlflexibleservers"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureMySqlFlexibleServerDBExample(t *testing.T) {
	t.Parallel()

	uniquePostfix := strings.ToLower(random.UniqueId())
	expectedFlexibleServerSkuName := "Standard_B1ms"
	expectedFlexibleServerStorageGb := "32"
	expectedFlexibleServerDatabaseCharSet := "utf8"
	expectedFlexibleServerDatabaseCollation := "utf8_unicode_ci"

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-mysqldb-flexible-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
			// SKU name returned by Azure is different from
			// the SKU name that terraform use to create the resource.
			// Appending the SKU Family as workaround.
			"mysql_flexible_server_sku_name":        "B_" + expectedFlexibleServerSkuName,
			"mysql_flexible_server_storage_size_gb": expectedFlexibleServerStorageGb,
			"mysql_flexible_server_db_charset":      expectedFlexibleServerDatabaseCharSet,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	expectedResourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	expectedMySqlFlexibleServerName := terraform.Output(t, terraformOptions, "mysql_flexible_server_name")

	expectedMySqlFlexibleServerDBName := terraform.Output(t, terraformOptions, "mysql_flexible_server_db_name")

	// website::tag::4:: Get mySql server details and assert them against the terraform output
	actualMySqlFlexibleServer := azure.GetMySqlFlexibleServer(t, expectedResourceGroupName, expectedMySqlFlexibleServerName, "")

	assert.Equal(t, expectedFlexibleServerSkuName, *actualMySqlFlexibleServer.Sku.Name)
	assert.Equal(t, expectedFlexibleServerStorageGb, fmt.Sprint(*actualMySqlFlexibleServer.ServerProperties.Storage.StorageSizeGB))

	assert.Equal(t, mysqlflexibleservers.ServerStateReady, actualMySqlFlexibleServer.ServerProperties.State)

	// website::tag::5:: Get  mySql server DB details and assert them against the terraform output
	actualMySqlServerFlexibleServerDatabase := azure.GetMySqlFlexibleServerDB(t, expectedResourceGroupName, expectedMySqlFlexibleServerName, expectedMySqlFlexibleServerDBName, "")

	assert.Equal(t, expectedFlexibleServerDatabaseCharSet, *actualMySqlServerFlexibleServerDatabase.DatabaseProperties.Charset)
	assert.Equal(t, expectedFlexibleServerDatabaseCollation, *actualMySqlServerFlexibleServerDatabase.DatabaseProperties.Collation)
}
