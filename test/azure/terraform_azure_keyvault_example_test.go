//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureKeyVaultExample(t *testing.T) {
	t.Parallel()

	uniquePostfix := random.UniqueId()

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-keyvault-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}

	// website::tag::6:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	keyVaultName := terraform.Output(t, terraformOptions, "key_vault_name")
	expectedSecretName := terraform.Output(t, terraformOptions, "secret_name")
	expectedKeyName := terraform.Output(t, terraformOptions, "key_name")
	expectedCertificateName := terraform.Output(t, terraformOptions, "certificate_name")

	// website::tag::4:: Determine whether the keyvault exists
	keyVault := azure.GetKeyVault(t, resourceGroupName, keyVaultName, "")
	assert.Equal(t, keyVaultName, *keyVault.Name)

	// website::tag::5:: Determine whether the secret, key, and certificate exists
	secretExists := azure.KeyVaultSecretExists(t, keyVaultName, expectedSecretName)
	assert.True(t, secretExists, "kv-secret does not exist")

	keyExists := azure.KeyVaultKeyExists(t, keyVaultName, expectedKeyName)
	assert.True(t, keyExists, "kv-key does not exist")

	certificateExists := azure.KeyVaultCertificateExists(t, keyVaultName, expectedCertificateName)
	assert.True(t, certificateExists, "kv-cert does not exist")
}
