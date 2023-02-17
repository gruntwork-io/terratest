package test

import (
	"fmt"
	"testing"

	awsSDK "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-aws-dynamodb-example using Terratest.
func TestTerraformAwsDynamoDBExample(t *testing.T) {
	t.Parallel()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	// Set up expected values to be checked later
	expectedTableName := fmt.Sprintf("terratest-aws-dynamodb-example-table-%s", random.UniqueId())
	expectedKmsKeyArn := aws.GetCmkArn(t, awsRegion, "alias/aws/dynamodb")
	expectedKeySchema := []*dynamodb.KeySchemaElement{
		{AttributeName: awsSDK.String("userId"), KeyType: awsSDK.String("HASH")},
		{AttributeName: awsSDK.String("department"), KeyType: awsSDK.String("RANGE")},
	}
	expectedTags := []*dynamodb.Tag{
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
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Look up the DynamoDB table by name
	table := aws.GetDynamoDBTable(t, awsRegion, expectedTableName)

	assert.Equal(t, "ACTIVE", awsSDK.StringValue(table.TableStatus))
	assert.ElementsMatch(t, expectedKeySchema, table.KeySchema)

	// Verify server-side encryption configuration
	assert.Equal(t, expectedKmsKeyArn, awsSDK.StringValue(table.SSEDescription.KMSMasterKeyArn))
	assert.Equal(t, "ENABLED", awsSDK.StringValue(table.SSEDescription.Status))
	assert.Equal(t, "KMS", awsSDK.StringValue(table.SSEDescription.SSEType))

	// Verify TTL configuration
	ttl := aws.GetDynamoDBTableTimeToLive(t, awsRegion, expectedTableName)
	assert.Equal(t, "expires", awsSDK.StringValue(ttl.AttributeName))
	assert.Equal(t, "ENABLED", awsSDK.StringValue(ttl.TimeToLiveStatus))

	// Verify resource tags
	tags := aws.GetDynamoDbTableTags(t, awsRegion, expectedTableName)
	assert.ElementsMatch(t, expectedTags, tags)
}
