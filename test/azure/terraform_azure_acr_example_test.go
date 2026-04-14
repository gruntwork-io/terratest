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

func TestTerraformAzureACRExample(t *testing.T) {
	t.Parallel()

	uniquePostfix := strings.ToLower(random.UniqueID())
	acrSKU := "Premium"

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/azure/terraform-azure-acr-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
			"sku":     acrSKU,
		},
	}

	// website::tag::5:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	acrName := terraform.OutputContext(t, t.Context(), terraformOptions, "container_registry_name")
	loginServer := terraform.OutputContext(t, t.Context(), terraformOptions, "login_server")

	// website::tag::4:: Assert
	assert.True(t, azure.ContainerRegistryExistsContext(t, t.Context(), acrName, resourceGroupName, ""))

	actualACR := azure.GetContainerRegistryContext(t, t.Context(), acrName, resourceGroupName, "")

	assert.Equal(t, loginServer, *actualACR.Properties.LoginServer)
	assert.True(t, *actualACR.Properties.AdminUserEnabled)
	assert.Equal(t, acrSKU, string(*actualACR.SKU.Name))
}
