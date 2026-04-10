//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureKeyVaultExample(t *testing.T) {
	t.Parallel()

	uniquePostfix := random.UniqueID()

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-keyvault-example",
		Vars: map[string]interface{}{
			"postfix": uniquePostfix,
		},
	}

	// website::tag::6:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	keyVaultName := terraform.OutputContext(t, t.Context(), terraformOptions, "key_vault_name")
	expectedSecretName := terraform.OutputContext(t, t.Context(), terraformOptions, "secret_name")
	expectedKeyName := terraform.OutputContext(t, t.Context(), terraformOptions, "key_name")
	expectedCertificateName := terraform.OutputContext(t, t.Context(), terraformOptions, "certificate_name")

	// website::tag::4:: Determine whether the keyvault exists
	keyVault := azure.GetKeyVaultContext(t, t.Context(), resourceGroupName, keyVaultName, "")
	assert.Equal(t, keyVaultName, *keyVault.Name)

	// website::tag::5:: Determine whether the secret, key, and certificate exists
	secretExists := azure.KeyVaultSecretExistsContext(t, t.Context(), keyVaultName, expectedSecretName)
	assert.True(t, secretExists, "kv-secret does not exist")

	keyExists := azure.KeyVaultKeyExistsContext(t, t.Context(), keyVaultName, expectedKeyName)
	assert.True(t, keyExists, "kv-key does not exist")

	certificateExists := azure.KeyVaultCertificateExistsContext(t, t.Context(), keyVaultName, expectedCertificateName)
	assert.True(t, certificateExists, "kv-cert does not exist")
}
