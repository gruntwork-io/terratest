package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-aws-rds-example using Terratest.
func TestTerraformAwsRdsExample(t *testing.T) {
	t.Parallel()

	// Give this RDS Instance a unique ID for a name tag so we can distinguish it from any other RDS Instance running
	// in your AWS account
	expectedName := fmt.Sprintf("terratest-aws-rds-example-%s", strings.ToLower(random.UniqueId()))
	expectedPort := int64(3306)
	expectedDatabaseName := "terratest"
	username := "username"
	password := "password"
	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)
	instanceType := aws.GetRecommendedRdsInstanceType(t, awsRegion, "mysql", "5.7.21", []string{"db.t2.micro", "db.t3.micro"})

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-aws-rds-example",

		// Variables to pass to our Terraform code using -var options
		// "username" and "password" should not be passed from here in a production scenario.
		Vars: map[string]interface{}{
			"name":                 expectedName,
			"engine_name":          "mysql",
			"major_engine_version": "5.7",
			"family":               "mysql5.7",
			"instance_class":       instanceType,
			"username":             username,
			"password":             password,
			"allocated_storage":    5,
			"license_model":        "general-public-license",
			"engine_version":       "5.7.21",
			"port":                 expectedPort,
			"database_name":        expectedDatabaseName,
			"region":               awsRegion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	dbInstanceID := terraform.Output(t, terraformOptions, "db_instance_id")

	// Look up the endpoint address and port of the RDS instance
	address := aws.GetAddressOfRdsInstance(t, dbInstanceID, awsRegion)
	port := aws.GetPortOfRdsInstance(t, dbInstanceID, awsRegion)
	schemaExistsInRdsInstance := aws.GetWhetherSchemaExistsInRdsMySqlInstance(t, address, port, username, password, expectedDatabaseName)
	// Lookup parameter values. All defined values are strings in the API call response
	generalLogParameterValue := aws.GetParameterValueForParameterOfRdsInstance(t, "general_log", dbInstanceID, awsRegion)
	allowSuspiciousUdfsParameterValue := aws.GetParameterValueForParameterOfRdsInstance(t, "allow-suspicious-udfs", dbInstanceID, awsRegion)

	// Lookup option values. All defined values are strings in the API call response
	mariadbAuditPluginServerAuditEventsOptionValue := aws.GetOptionSettingForOfRdsInstance(t, "MARIADB_AUDIT_PLUGIN", "SERVER_AUDIT_EVENTS", dbInstanceID, awsRegion)

	// Verify that the address is not null
	assert.NotNil(t, address)
	// Verify that the DB instance is listening on the port mentioned
	assert.Equal(t, expectedPort, port)
	// Verify that the table/schema requested for creation is actually present in the database
	assert.True(t, schemaExistsInRdsInstance)
	// Booleans are (string) "0", "1"
	assert.Equal(t, "0", generalLogParameterValue)
	// Values not set are "". This is custom behavior defined.
	assert.Equal(t, "", allowSuspiciousUdfsParameterValue)
	// assert.Equal(t, "", mariadbAuditPluginServerAuditEventsOptionValue)
	assert.Equal(t, "CONNECT", mariadbAuditPluginServerAuditEventsOptionValue)
}
