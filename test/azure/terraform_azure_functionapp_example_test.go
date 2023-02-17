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

func TestTerraformAzureFunctionAppExample(t *testing.T) {
	t.Parallel()

	//_random := strings.ToLower(random.UniqueId())
	uniquePostfix := strings.ToLower(random.UniqueId())

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/azure/terraform-azure-functionapp-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}
	// website::tag::5:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	appName := terraform.Output(t, terraformOptions, "function_app_name")

	appId := terraform.Output(t, terraformOptions, "function_app_id")
	appDefaultHostName := terraform.Output(t, terraformOptions, "default_hostname")
	appKind := terraform.Output(t, terraformOptions, "function_app_kind")

	// website::tag::4:: Assert
	assert.True(t, azure.AppExists(t, appName, resourceGroupName, ""))
	site := azure.GetAppService(t, appName, resourceGroupName, "")

	assert.Equal(t, appId, *site.ID)
	assert.Equal(t, appDefaultHostName, *site.DefaultHostName)
	assert.Equal(t, appKind, *site.Kind)

	assert.NotEmpty(t, *site.OutboundIPAddresses)
	assert.Equal(t, "Running", *site.State)
}
