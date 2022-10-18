//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/services/sql/mgmt/2014-04-01/sql"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureSQLDBExample(t *testing.T) {
	t.Parallel()

	uniquePostfix := strings.ToLower(random.UniqueId())

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-sqldb-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	expectedSQLServerID := terraform.Output(t, terraformOptions, "sql_server_id")
	expectedSQLServerName := terraform.Output(t, terraformOptions, "sql_server_name")

	expectedSQLServerFullDomainName := terraform.Output(t, terraformOptions, "sql_server_full_domain_name")
	expectedSQLDBName := terraform.Output(t, terraformOptions, "sql_database_name")

	expectedSQLDBID := terraform.Output(t, terraformOptions, "sql_database_id")
	expectedResourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	expectedSQLDBStatus := "Online"

	// website::tag::4:: Get the SQL server details and assert them against the terraform output
	actualSQLServer := azure.GetSQLServer(t, expectedResourceGroupName, expectedSQLServerName, "")

	assert.Equal(t, expectedSQLServerID, *actualSQLServer.ID)
	assert.Equal(t, expectedSQLServerFullDomainName, *actualSQLServer.FullyQualifiedDomainName)
	assert.Equal(t, sql.ServerStateReady, actualSQLServer.State)

	// website::tag::5:: Get the SQL server DB details and assert them against the terraform output
	actualSQLDatabase := azure.GetSQLDatabase(t, expectedResourceGroupName, expectedSQLServerName, expectedSQLDBName, "")

	assert.Equal(t, expectedSQLDBID, *actualSQLDatabase.ID)
	assert.Equal(t, expectedSQLDBStatus, *actualSQLDatabase.Status)
}
