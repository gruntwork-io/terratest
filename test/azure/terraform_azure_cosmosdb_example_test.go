//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureCosmosDBExample(t *testing.T) {
	t.Parallel()

	subscriptionID := ""
	uniquePostfix := random.Random(10000, 99999)
	throughput := 400

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-cosmosdb-example",
		Vars: map[string]interface{}{
			"postfix":    uniquePostfix,
			"throughput": throughput,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	resourceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "resource_group_name")
	accountName := terraform.OutputContext(t, t.Context(), terraformOptions, "account_name")

	// website::tag::4:: Get CosmosDB details and assert them against the terraform output
	// NOTE: the value of subscriptionID can be left blank, it will be replaced by the value
	//       of the environment variable ARM_SUBSCRIPTION_ID

	// Database Account properties
	actualCosmosDBAccount := azure.GetCosmosDBAccountContext(t, t.Context(), subscriptionID, resourceGroupName, accountName)
	assert.Equal(t, accountName, *actualCosmosDBAccount.Name)
	assert.Equal(t, armcosmos.DatabaseAccountKindGlobalDocumentDB, *actualCosmosDBAccount.Kind)
	assert.Equal(t, armcosmos.DefaultConsistencyLevelSession, *actualCosmosDBAccount.Properties.ConsistencyPolicy.DefaultConsistencyLevel)

	// SQL Database properties
	cosmosSQLDB := azure.GetCosmosDBSQLDatabaseContext(t, t.Context(), subscriptionID, resourceGroupName, accountName, "testdb")
	assert.Equal(t, "testdb", *cosmosSQLDB.Name)

	// SQL Database throughput
	cosmosSQLDBThroughput := azure.GetCosmosDBSQLDatabaseThroughputContext(t, t.Context(), subscriptionID, resourceGroupName, accountName, "testdb")
	assert.Equal(t, int32(throughput), *cosmosSQLDBThroughput.Properties.Resource.Throughput)

	// SQL Container properties
	cosmosSQLContainer1 := azure.GetCosmosDBSQLContainerContext(t, t.Context(), subscriptionID, resourceGroupName, accountName, "testdb", "test-container-1")
	cosmosSQLContainer2 := azure.GetCosmosDBSQLContainerContext(t, t.Context(), subscriptionID, resourceGroupName, accountName, "testdb", "test-container-2")
	cosmosSQLContainer3 := azure.GetCosmosDBSQLContainerContext(t, t.Context(), subscriptionID, resourceGroupName, accountName, "testdb", "test-container-3")
	assert.Equal(t, "test-container-1", *cosmosSQLContainer1.Name)
	assert.Equal(t, "/key1", *cosmosSQLContainer1.Properties.Resource.PartitionKey.Paths[0])
	assert.Equal(t, "test-container-2", *cosmosSQLContainer2.Name)
	assert.Equal(t, "/key2", *cosmosSQLContainer2.Properties.Resource.PartitionKey.Paths[0])
	assert.Equal(t, "test-container-3", *cosmosSQLContainer3.Name)
	assert.Equal(t, "/key3", *cosmosSQLContainer3.Properties.Resource.PartitionKey.Paths[0])

	// SQL Container throughput
	cosmosSQLContainer1Throughput := azure.GetCosmosDBSQLContainerThroughputContext(t, t.Context(), subscriptionID, resourceGroupName, accountName, "testdb", "test-container-1")
	assert.Equal(t, int32(400), *cosmosSQLContainer1Throughput.Properties.Resource.Throughput)
}
