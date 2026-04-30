//go:build azure
// +build azure

package test_test

import (
	"strconv"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mysql/armmysqlflexibleservers"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureMySQLDBExample(t *testing.T) {
	t.Parallel()

	uniquePostfix := strings.ToLower(random.UniqueID())
	expectedServerSkuName := "B_Standard_B1ms"
	expectedServerStorageSizeGB := "20"
	expectedDatabaseCharSet := "utf8"
	expectedDatabaseCollation := "utf8_unicode_ci"

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-mysqldb-example",
		Vars: map[string]interface{}{
			"postfix":                     uniquePostfix,
			"mysqlserver_sku_name":        expectedServerSkuName,
			"mysqlserver_storage_size_gb": expectedServerStorageSizeGB,
			"mysqldb_charset":             expectedDatabaseCharSet,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	expectedResourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	expectedMYSQLServerName := terraform.OutputContext(t, t.Context(), terraformOptions, "mysql_server_name")

	expectedMYSQLDBName := terraform.OutputContext(t, t.Context(), terraformOptions, "mysql_database_name")

	// website::tag::4:: Get mySQL flexible server details and assert them against the terraform output
	actualMYSQLServer := azure.GetMYSQLServerContext(t, t.Context(), "", expectedResourceGroupName, expectedMYSQLServerName)

	assert.Equal(t, expectedServerSkuName, *actualMYSQLServer.SKU.Name)
	// Flexible servers expose storage in GB (StorageSizeGB) rather than the legacy StorageMB.
	assert.Equal(t, expectedServerStorageSizeGB, strconv.Itoa(int(*actualMYSQLServer.Properties.Storage.StorageSizeGB)))

	// Flexible servers use Properties.State (typed as *ServerState) rather than the legacy UserVisibleState.
	assert.Equal(t, armmysqlflexibleservers.ServerStateReady, *actualMYSQLServer.Properties.State)

	// website::tag::5:: Get  mySQL flexible server DB details and assert them against the terraform output
	actualDatabase := azure.GetMYSQLDBContext(t, t.Context(), "", expectedResourceGroupName, expectedMYSQLServerName, expectedMYSQLDBName)

	assert.Equal(t, expectedDatabaseCharSet, *actualDatabase.Properties.Charset)
	assert.Equal(t, expectedDatabaseCollation, *actualDatabase.Properties.Collation)
}
