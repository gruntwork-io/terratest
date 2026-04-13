//go:build aws

package test_test

import (
	"testing"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-aws-dynamodb-example using Terratest.
func TestTerraformAwsDynamoDBExample(t *testing.T) {
	t.Parallel()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegionContext(t, t.Context(), nil, nil)

	// Set up expected values to be checked later
	expectedTableName := "terratest-aws-dynamodb-example-table-" + random.UniqueID()
	expectedKmsKeyArn := aws.GetCmkArnContext(t, t.Context(), awsRegion, "alias/aws/dynamodb")
	expectedKeySchema := []types.KeySchemaElement{
		{AttributeName: awsSDK.String("userId"), KeyType: types.KeyTypeHash},
		{AttributeName: awsSDK.String("department"), KeyType: types.KeyTypeRange},
	}
	expectedTags := []types.Tag{
		{Key: awsSDK.String("Environment"), Value: awsSDK.String("production")},
	}

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-aws-dynamodb-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"table_name": expectedTableName,
			"region":     awsRegion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// Look up the DynamoDB table by name
	table := aws.GetDynamoDBTableContext(t, t.Context(), awsRegion, expectedTableName)

	assert.Equal(t, "ACTIVE", string(table.TableStatus))
	assert.ElementsMatch(t, expectedKeySchema, table.KeySchema)

	// Verify server-side encryption configuration
	assert.Equal(t, expectedKmsKeyArn, awsSDK.ToString(table.SSEDescription.KMSMasterKeyArn))
	assert.Equal(t, "ENABLED", string(table.SSEDescription.Status))
	assert.Equal(t, "KMS", string(table.SSEDescription.SSEType))

	// Verify TTL configuration
	ttl := aws.GetDynamoDBTableTimeToLiveContext(t, t.Context(), awsRegion, expectedTableName)
	assert.Equal(t, "expires", awsSDK.ToString(ttl.AttributeName))
	assert.Equal(t, "ENABLED", string(ttl.TimeToLiveStatus))

	// Verify resource tags
	tags := aws.GetDynamoDBTableTagsContext(t, t.Context(), awsRegion, expectedTableName)
	assert.ElementsMatch(t, expectedTags, tags)
}
