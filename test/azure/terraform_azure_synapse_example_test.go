//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test_test

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

	uniquePostfix := strings.ToLower(random.UniqueID())
	expectedSynapseSQLUser := "sqladminuser"
	expectedSynapseProvisioningState := "Succeeded"
	expectedLocation := "westus2"
	expectedSyPoolSkuName := "DW100c"

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-synapse-example",
		Vars: map[string]interface{}{
			"postfix":                  uniquePostfix,
			"synapse_sql_user":         expectedSynapseSQLUser,
			"location":                 expectedLocation,
			"synapse_sqlpool_sku_name": expectedSyPoolSkuName,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)
	// terraform.InitE(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	expectedResourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	expectedSyDLgen2Name := terraform.OutputContext(t, t.Context(), terraformOptions, "synapse_dlgen2_name")
	expectedSyWorkspaceName := terraform.OutputContext(t, t.Context(), terraformOptions, "synapse_workspace_name")
	expectedSQLPoolName := terraform.OutputContext(t, t.Context(), terraformOptions, "synapse_sqlpool_name")

	// website::tag::4:: Get synapse details and assert them against the terraform output
	actualSynapseWorkspace := azure.GetSynapseWorkspaceContext(t, t.Context(), "", expectedResourceGroupName, expectedSyWorkspaceName)
	actualSynapseSQLPool := azure.GetSynapseSQLPoolContext(t, t.Context(), "", expectedResourceGroupName, expectedSyWorkspaceName, expectedSQLPoolName)

	assert.Equal(t, expectedSyWorkspaceName, *actualSynapseWorkspace.Name)
	assert.Equal(t, expectedSynapseSQLUser, *actualSynapseWorkspace.Properties.SQLAdministratorLogin)
	assert.Equal(t, expectedSynapseProvisioningState, *actualSynapseWorkspace.Properties.ProvisioningState)
	assert.Equal(t, expectedLocation, *actualSynapseWorkspace.Location)
	assert.Equal(t, expectedSyDLgen2Name, *actualSynapseWorkspace.Properties.DefaultDataLakeStorage.Filesystem)

	assert.Equal(t, expectedSQLPoolName, *actualSynapseSQLPool.Name)
	assert.Equal(t, expectedSyPoolSkuName, *actualSynapseSQLPool.SKU.Name)
}
