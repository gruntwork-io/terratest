//go:build aws

package test_test

import (
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-aws-rds-example using Terratest.
func TestTerraformAwsRdsExample(t *testing.T) {
	t.Parallel()

	ttable := []struct {
		schemaCheck    func(t *testing.T, dbURL string, dbPort int32, dbUsername string, dbPassword string, expectedSchemaName string) bool
		expectedOptins map[struct {
			opName  string
			setName string
		}]string
		expectedParameter  map[string]string
		name               string
		engineName         string
		majorEngineVersion string
		engineFamily       string
		licenseModel       string
	}{
		{
			name:               "mysql",
			engineName:         "mysql",
			majorEngineVersion: "5.7",
			engineFamily:       "mysql5.7",
			licenseModel:       "general-public-license",
			schemaCheck: func(t *testing.T, dbURL string, dbPort int32, dbUsername, dbPassword, expectedSchemaName string) bool {
				t.Helper()

				return aws.GetWhetherSchemaExistsInRdsMySQLInstanceContext(t, t.Context(), dbURL, dbPort, dbUsername, dbPassword, expectedSchemaName)
			},
			expectedOptins: map[struct {
				opName  string
				setName string
			}]string{
				{opName: "MARIADB_AUDIT_PLUGIN", setName: "SERVER_AUDIT_EVENTS"}: "CONNECT",
			},
			expectedParameter: map[string]string{
				"general_log":           "0",
				"allow-suspicious-udfs": "",
			},
		},
		{
			name:               "postgres",
			engineName:         "postgres",
			majorEngineVersion: "13",
			engineFamily:       "postgres13",
			licenseModel:       "postgresql-license",
			schemaCheck: func(t *testing.T, dbURL string, dbPort int32, dbUsername, dbPassword, expectedSchemaName string) bool {
				t.Helper()

				return aws.GetWhetherSchemaExistsInRdsPostgresInstanceContext(t, t.Context(), dbURL, dbPort, dbUsername, dbPassword, expectedSchemaName)
			},
		},
	}

	for _, tt := range ttable {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Give this RDS Instance a unique ID for a name tag so we can distinguish it from any other RDS Instance running
			// in your AWS account
			expectedName := "terratest-aws-rds-example-" + strings.ToLower(random.UniqueID())
			expectedPort := int32(3306)
			expectedDatabaseName := "terratest"
			username := "username"
			password := "password"
			// Pick a random AWS region to test in. This helps ensure your code works in all regions.
			awsRegion := aws.GetRandomStableRegionContext(t, t.Context(), nil, nil)
			engineVersion := aws.GetValidEngineVersionContext(t, t.Context(), awsRegion, tt.engineName, tt.majorEngineVersion)
			instanceType := aws.GetRecommendedRdsInstanceTypeContext(t, t.Context(), awsRegion, tt.engineName, engineVersion, []string{"db.t2.micro", "db.t3.micro", "db.t3.small"})
			moduleFolder := test_structure.CopyTerraformFolderToTemp(t, "../", "examples/terraform-aws-rds-example")

			// Construct the terraform options with default retryable errors to handle the most common retryable errors in
			// terraform testing.
			terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
				// The path to where our Terraform code is located
				TerraformDir: moduleFolder,

				// Variables to pass to our Terraform code using -var options
				// "username" and "password" should not be passed from here in a production scenario.
				Vars: map[string]interface{}{
					"name":                 expectedName,
					"engine_name":          tt.engineName,
					"major_engine_version": tt.majorEngineVersion,
					"family":               tt.engineFamily,
					"instance_class":       instanceType,
					"username":             username,
					"password":             password,
					"allocated_storage":    5,
					"license_model":        tt.licenseModel,
					"engine_version":       engineVersion,
					"port":                 expectedPort,
					"database_name":        expectedDatabaseName,
					"region":               awsRegion,
				},
			})

			// At the end of the test, run `terraform destroy` to clean up any resources that were created
			defer terraform.DestroyContext(t, t.Context(), terraformOptions)

			// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
			terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

			// Run `terraform output` to get the value of an output variable
			dbInstanceID := terraform.OutputContext(t, t.Context(), terraformOptions, "db_instance_id")

			// Look up the endpoint address and port of the RDS instance
			address := aws.GetAddressOfRdsInstanceContext(t, t.Context(), dbInstanceID, awsRegion)
			port := aws.GetPortOfRdsInstanceContext(t, t.Context(), dbInstanceID, awsRegion)
			schemaExistsInRdsInstance := tt.schemaCheck(t, address, port, username, password, expectedDatabaseName)
			// Lookup parameter values. All defined values are strings in the API call response

			// Verify that the address is not null
			assert.NotNil(t, address)
			// Verify that the DB instance is listening on the port mentioned
			assert.Equal(t, expectedPort, port)
			// Verify that the table/schema requested for creation is actually present in the database
			assert.True(t, schemaExistsInRdsInstance)

			// assert expected parameters
			for k, v := range tt.expectedParameter {
				assert.Equal(t, v, aws.GetParameterValueForParameterOfRdsInstanceContext(t, t.Context(), k, dbInstanceID, awsRegion))
			}

			// assert all parameters
			params := aws.GetAllParametersOfRdsInstanceContext(t, t.Context(), dbInstanceID, awsRegion)

			paramNames := map[string]struct{}{}
			for _, param := range params {
				paramNames[*param.ParameterName] = struct{}{}
			}

			assert.Len(t, paramNames, len(params), "should return no duplicate parameters")
			assert.Greater(t, len(paramNames), 100)

			// assert expected options
			for k, v := range tt.expectedOptins {
				// Lookup option values. All defined values are strings in the API call response
				assert.Equal(t, v, aws.GetOptionSettingForOfRdsInstanceContext(t, t.Context(), k.opName, k.setName, dbInstanceID, awsRegion))
			}
		})
	}
}
