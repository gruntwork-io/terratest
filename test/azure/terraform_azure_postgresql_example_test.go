// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestPostgreSQLDatabase(t *testing.T) {
	t.Parallel()

	uniquePostfix := strings.ToLower(random.UniqueId())

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../../examples/azure/terraform-azure-postgresql-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
		NoColor: true,
	})
	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")

	// website::tag::3:: Run `terraform output` to get the values of output variables
	expectedServername := "postgresqlserver-" + uniquePostfix // see fixture
	actualServername := terraform.Output(t, terraformOptions, "servername")
	rgName := terraform.Output(t, terraformOptions, "rgname")
	expectedSkuName := terraform.Output(t, terraformOptions, "sku_name")

	// website::tag::4:: Get the Server details and assert them against the terraform output
	actualServer := azure.GetPostgreSQLServer(t, rgName, actualServername, subscriptionID)
	// Verify
	assert.NotNil(t, actualServer)
	assert.Equal(t, expectedServername, actualServername)
	assert.Equal(t, expectedSkuName, *actualServer.Sku.Name)

}
