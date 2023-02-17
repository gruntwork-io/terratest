//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureStorageExample(t *testing.T) {
	t.Parallel()

	// subscriptionID is overridden by the environment variable "ARM_SUBSCRIPTION_ID"
	subscriptionID := ""
	uniquePostfix := random.UniqueId()

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-storage-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"postfix": strings.ToLower(uniquePostfix),
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.Output(t, terraformOptions, "resource_group_name")
	storageAccountName := terraform.Output(t, terraformOptions, "storage_account_name")
	storageAccountTier := terraform.Output(t, terraformOptions, "storage_account_account_tier")
	storageAccountKind := terraform.Output(t, terraformOptions, "storage_account_account_kind")
	storageBlobContainerName := terraform.Output(t, terraformOptions, "storage_container_name")

	// website::tag::4:: Verify storage account properties and ensure it matches the output.
	storageAccountExists := azure.StorageAccountExists(t, storageAccountName, resourceGroupName, subscriptionID)
	assert.True(t, storageAccountExists, "storage account does not exist")

	containerExists := azure.StorageBlobContainerExists(t, storageBlobContainerName, storageAccountName, resourceGroupName, subscriptionID)
	assert.True(t, containerExists, "storage container does not exist")

	publicAccess := azure.GetStorageBlobContainerPublicAccess(t, storageBlobContainerName, storageAccountName, resourceGroupName, subscriptionID)
	assert.False(t, publicAccess, "storage container has public access")

	accountKind := azure.GetStorageAccountKind(t, storageAccountName, resourceGroupName, subscriptionID)
	assert.Equal(t, storageAccountKind, accountKind, "storage account kind mismatch")

	skuTier := azure.GetStorageAccountSkuTier(t, storageAccountName, resourceGroupName, subscriptionID)
	assert.Equal(t, storageAccountTier, skuTier, "sku tier mismatch")

	actualDNSString := azure.GetStorageDNSString(t, storageAccountName, resourceGroupName, subscriptionID)
	storageSuffix, _ := azure.GetStorageURISuffixE()
	expectedDNS := fmt.Sprintf("https://%s.blob.%s/", storageAccountName, storageSuffix)
	assert.Equal(t, expectedDNS, actualDNSString, "Storage DNS string mismatch")
}
