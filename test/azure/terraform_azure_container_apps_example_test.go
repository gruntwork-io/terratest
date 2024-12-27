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

func TestTerraformAzureContainerAppExample(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	uniquePostfix := strings.ToLower(random.UniqueId())

	terraformOptions := &terraform.Options{
		TerraformBinary: "",
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-container-apps-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}

	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	envName := terraform.Output(t, terraformOptions, "container_app_env_name")
	containerAppName := terraform.Output(t, terraformOptions, "container_app_name")
	containerAppJobName := terraform.Output(t, terraformOptions, "container_app_job_name")

	// NOTE: the value of subscriptionID can be left blank, it will be replaced by the value
	//       of the environment variable ARM_SUBSCRIPTION_ID

	envExsists := azure.ManagedEnvironmentExists(t, envName, resourceGroupName, subscriptionID)
	assert.True(t, envExsists)

	actualEnv := azure.GetManagedEnvironment(t, envName, resourceGroupName, subscriptionID)
	assert.Equal(t, envName, *actualEnv.Name)

	containerAppExists := azure.ContainerAppExists(t, containerAppName, resourceGroupName, subscriptionID)
	assert.True(t, containerAppExists)

	actualContainerApp := azure.GetContainerApp(t, containerAppName, resourceGroupName, subscriptionID)
	assert.Equal(t, containerAppName, *actualContainerApp.Name)

	containerAppJobExists := azure.ContainerAppJobExists(t, containerAppJobName, resourceGroupName, subscriptionID)
	assert.True(t, containerAppJobExists)

	actualContainerAppJob := azure.GetContainerAppJob(t, containerAppJobName, resourceGroupName, subscriptionID)
	assert.Equal(t, containerAppJobName, *actualContainerAppJob.Name)
}
