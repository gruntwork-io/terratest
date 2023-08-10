//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureSynapseExample(t *testing.T) {
	t.Parallel()

	uniquePostfix := strings.ToLower(random.UniqueId())
	expectedSynapseSqlUser := "sqladminuser"
	expectedSynapseProvisioningState := "Succeeded"
	expectedLocation := "westus2"
	expectedSyPoolSkuName := "DW100c"

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-synapse-example",
		Vars: map[string]interface{}{
			"postfix":                  uniquePostfix,
			"synapse_sql_user":         expectedSynapseSqlUser,
			"location":                 expectedLocation,
			"synapse_sqlpool_sku_name": expectedSyPoolSkuName,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)
	// terraform.InitE(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	expectedResourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	expectedSyDLgen2Name := terraform.Output(t, terraformOptions, "synapse_dlgen2_name")
	expectedSyWorkspaceName := terraform.Output(t, terraformOptions, "synapse_workspace_name")
	expectedSqlPoolName := terraform.Output(t, terraformOptions, "synapse_sqlpool_name")

	// website::tag::4:: Get synapse details and assert them against the terraform output
	actualSynapseWorkspace := azure.GetSynapseWorkspace(t, expectedResourceGroupName, expectedSyWorkspaceName, "")
	actualSynapseSqlPool := azure.GetSynapseSqlPool(t, expectedResourceGroupName, expectedSyWorkspaceName, expectedSqlPoolName, "")

	assert.Equal(t, expectedSyWorkspaceName, *actualSynapseWorkspace.Name)
	assert.Equal(t, expectedSynapseSqlUser, *actualSynapseWorkspace.WorkspaceProperties.SQLAdministratorLogin)
	assert.Equal(t, expectedSynapseProvisioningState, *actualSynapseWorkspace.WorkspaceProperties.ProvisioningState)
	assert.Equal(t, expectedLocation, *actualSynapseWorkspace.Location)
	assert.Equal(t, expectedSyDLgen2Name, *actualSynapseWorkspace.WorkspaceProperties.DefaultDataLakeStorage.Filesystem)

	assert.Equal(t, expectedSqlPoolName, *actualSynapseSqlPool.Name)
	assert.Equal(t, expectedSyPoolSkuName, *actualSynapseSqlPool.Sku.Name)
}
