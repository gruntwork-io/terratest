// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/cosmos-db/mgmt/documentdb"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformAzureCosmosDBExample(t *testing.T) {
	t.Parallel()

	//expectedResourceGroupName := fmt.Sprintf("terratest-cosmosdb-rg-%s", random.UniqueId())
	//expectedAccountName := fmt.Sprintf("terratest-%d", random.Random(10000, 99999))
	expectedResourceGroupName := "terratest-cosmosdb-rg"
	expectedAccountName := "terratest-12345"

	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../../examples/azure/terraform-azure-cosmosdb-example",
		Vars: map[string]interface{}{
			"resource_group_name":   expectedResourceGroupName,
			"cosmosdb_account_name": expectedAccountName,
		},
	}

	// website::tag::4:: At the end of the test, run `terraform destroy` to clean up any resources that were created
	//defer terraform.Destroy(t, terraformOptions)

	// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// website::tag::3:: Run `terraform output` to get the values of output variables
	actualAccountName := terraform.Output(t, terraformOptions, "name")
	assert.Equal(t, expectedAccountName, actualAccountName)

	// website::tag::4:: Get CosmosDB details and assert them against the terraform output
	// NOTE: the value of subscriptionID can be left blank, it will be replaced by the value
	//       of the environment variable ARM_SUBSCRIPTION_ID

	// Database Account properties
	actualCosmosDBAccount := azure.GetCosmosDBAccount(t, "", expectedResourceGroupName, expectedAccountName)
	assert.Equal(t, expectedAccountName, *actualCosmosDBAccount.Name)
	assert.Equal(t, documentdb.GlobalDocumentDB, actualCosmosDBAccount.Kind)
	assert.Equal(t, documentdb.Session, actualCosmosDBAccount.DatabaseAccountGetProperties.ConsistencyPolicy.DefaultConsistencyLevel)

	// SQL Database properties
	cosmosSQLDB := azure.GetCosmosDBSQLDatabase(t, "", expectedResourceGroupName, expectedAccountName, "testdb")
	assert.Equal(t, "testdb", *cosmosSQLDB.Name)

	// SQL Container properties
	cosmosSQLContainer1 := azure.GetCosmosDBSQLContainer(t, "", expectedResourceGroupName, expectedAccountName, "testdb", "test-container-1")
	cosmosSQLContainer2 := azure.GetCosmosDBSQLContainer(t, "", expectedResourceGroupName, expectedAccountName, "testdb", "test-container-2")
	cosmosSQLContainer3 := azure.GetCosmosDBSQLContainer(t, "", expectedResourceGroupName, expectedAccountName, "testdb", "test-container-3")
	assert.Equal(t, "test-container-1", *cosmosSQLContainer1.Name)
	assert.Equal(t, "/key1", (*cosmosSQLContainer1.SQLContainerGetProperties.Resource.PartitionKey.Paths)[0])
	assert.Equal(t, "test-container-2", *cosmosSQLContainer2.Name)
	assert.Equal(t, "/key2", (*cosmosSQLContainer2.SQLContainerGetProperties.Resource.PartitionKey.Paths)[0])
	assert.Equal(t, "test-container-3", *cosmosSQLContainer3.Name)
	assert.Equal(t, "/key3", (*cosmosSQLContainer3.SQLContainerGetProperties.Resource.PartitionKey.Paths)[0])
}
