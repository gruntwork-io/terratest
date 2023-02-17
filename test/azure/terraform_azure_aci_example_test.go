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

func TestTerraformAzureACIExample(t *testing.T) {
	t.Parallel()

	uniquePostfix := strings.ToLower(random.UniqueId())

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		TerraformDir: "../../examples/azure/terraform-azure-aci-example",
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
	aciName := terraform.Output(t, terraformOptions, "container_instance_name")
	ipAddress := terraform.Output(t, terraformOptions, "ip_address")
	fqdn := terraform.Output(t, terraformOptions, "fqdn")

	// website::tag::4:: Assert
	assert.True(t, azure.ContainerInstanceExists(t, aciName, resourceGroupName, ""))

	actualInstance := azure.GetContainerInstance(t, aciName, resourceGroupName, "")

	assert.Equal(t, ipAddress, *actualInstance.ContainerGroupProperties.IPAddress.IP)
	assert.Equal(t, fqdn, *actualInstance.ContainerGroupProperties.IPAddress.Fqdn)
}
