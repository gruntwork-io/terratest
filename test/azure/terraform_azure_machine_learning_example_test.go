//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning/v3"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureMachineLearningWorkspaceExample(t *testing.T) {
	t.Parallel()

	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	subscriptionID := ""
	uniquePostfix := strings.ToLower(random.UniqueId())

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-machinelearning-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}

	// website::tag::8:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	workspaceName := terraform.Output(t, terraformOptions, "workspace_name")
	computeName := terraform.Output(t, terraformOptions, "compute_name")

	//website::tag::4:: Verify the machine learning workspace exists
	workspaceExists := azure.MachinelearningWorkspaceExists(t, resourceGroupName, workspaceName, subscriptionID)
	assert.True(t, workspaceExists, "Machine Learning Workspace does not exist")

	//website::tag::5:: Verify the machine learning compute exists
	computeExists := azure.MachinelearningComputeExists(t, resourceGroupName, workspaceName, computeName, subscriptionID)
	assert.True(t, computeExists, "Machine Learning Compute does not exist")

	//website::tag::6:: Verify that the Machine Learning Workspace is provisioned successfully
	workspace := azure.GetMachinelearningWorkspace(t, resourceGroupName, workspaceName, subscriptionID)
	assert.Equal(t, armmachinelearning.ProvisioningStateSucceeded, *workspace.Properties.ProvisioningState)

	//website::tag::7:: Verify that the Machine Learning Compute is provisioned successfully
	compute := azure.GetMachinelearningCompute(t, resourceGroupName, workspaceName, computeName, subscriptionID)
	assert.Equal(t, armmachinelearning.ProvisioningStateSucceeded, *compute.Properties.GetCompute().ProvisioningState)
}
