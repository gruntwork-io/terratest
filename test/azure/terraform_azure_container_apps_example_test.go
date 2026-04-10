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

func TestTerraformAzureContainerAppExample(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	uniquePostfix := strings.ToLower(random.UniqueID())

	terraformOptions := &terraform.Options{
		TerraformBinary: "",
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-container-apps-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	envName := terraform.OutputContext(t, t.Context(), terraformOptions, "container_app_env_name")
	containerAppName := terraform.OutputContext(t, t.Context(), terraformOptions, "container_app_name")
	containerAppJobName := terraform.OutputContext(t, t.Context(), terraformOptions, "container_app_job_name")

	// NOTE: the value of subscriptionID can be left blank, it will be replaced by the value
	//       of the environment variable ARM_SUBSCRIPTION_ID

	envExsists := azure.ManagedEnvironmentExistsContext(t, t.Context(), envName, resourceGroupName, subscriptionID)
	assert.True(t, envExsists)

	actualEnv := azure.GetManagedEnvironmentContext(t, t.Context(), envName, resourceGroupName, subscriptionID)
	assert.Equal(t, envName, *actualEnv.Name)

	containerAppExists := azure.ContainerAppExistsContext(t, t.Context(), containerAppName, resourceGroupName, subscriptionID)
	assert.True(t, containerAppExists)

	actualContainerApp := azure.GetContainerAppContext(t, t.Context(), containerAppName, resourceGroupName, subscriptionID)
	assert.Equal(t, containerAppName, *actualContainerApp.Name)

	containerAppJobExists := azure.ContainerAppJobExistsContext(t, t.Context(), containerAppJobName, resourceGroupName, subscriptionID)
	assert.True(t, containerAppJobExists)

	actualContainerAppJob := azure.GetContainerAppJobContext(t, t.Context(), containerAppJobName, resourceGroupName, subscriptionID)
	assert.Equal(t, containerAppJobName, *actualContainerAppJob.Name)
}
