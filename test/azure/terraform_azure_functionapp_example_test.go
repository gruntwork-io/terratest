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

func TestTerraformAzureFunctionAppExample(t *testing.T) {
	t.Parallel()

	// _random := strings.ToLower(random.UniqueID())
	uniquePostfix := strings.ToLower(random.UniqueID())

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/azure/terraform-azure-functionapp-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}
	// website::tag::5:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	appName := terraform.OutputContext(t, t.Context(), terraformOptions, "function_app_name")

	appID := terraform.OutputContext(t, t.Context(), terraformOptions, "function_app_id")
	appDefaultHostName := terraform.OutputContext(t, t.Context(), terraformOptions, "default_hostname")
	appKind := terraform.OutputContext(t, t.Context(), terraformOptions, "function_app_kind")

	// website::tag::4:: Assert
	assert.True(t, azure.AppExistsContext(t, t.Context(), appName, resourceGroupName, ""))
	site := azure.GetAppServiceContext(t, t.Context(), appName, resourceGroupName, "")

	assert.Equal(t, appID, *site.ID)
	assert.Equal(t, appDefaultHostName, *site.Properties.DefaultHostName)
	assert.Equal(t, appKind, *site.Kind)

	assert.NotEmpty(t, *site.Properties.OutboundIPAddresses)
	assert.Equal(t, "Running", *site.Properties.State)
}
