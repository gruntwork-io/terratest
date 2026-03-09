//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation

package test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"github.com/stretchr/testify/assert"
)

func TestTerraformGcpCloudSQLExample(t *testing.T) {
	ttable := []struct {
		name            string
		databaseVersion string
	}{
		{
			name:            "mysql",
			databaseVersion: "MYSQL_8_0",
		},
		{
			name:            "postgres",
			databaseVersion: "POSTGRES_14",
		},
	}

	for _, tt := range ttable {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Get the Project ID from the environment variable.
			projectID := gcp.GetGoogleProjectIDFromEnvVar(t)

			// Create a random unique name for the Cloud SQL instance
			// so multiple tests running simultaneously don't collide.
			expectedInstanceName := fmt.Sprintf("terratest-cloudsql-%s", random.UniqueId())

			exampleDir := test_structure.CopyTerraformFolderToTemp(t, "../../", "examples/terraform-gcp-cloudsql-example")

			terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
				// The path to where our Terraform code is located
				TerraformDir: exampleDir,

				// Variables to pass to our Terraform code using -var options
				Vars: map[string]interface{}{
					"gcp_project_id":   projectID,
					"instance_name":    expectedInstanceName,
					"database_version": tt.databaseVersion,
				},
			})

			// At the end of the test, run `terraform destroy` to clean up any resources that were created
			defer terraform.Destroy(t, terraformOptions)

			// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
			terraform.InitAndApply(t, terraformOptions)

			// Pull out the outputs from the Terraform configuration
			actualInstanceName := terraform.Output(t, terraformOptions, "instance_name")
			actualDatabaseVersion := terraform.Output(t, terraformOptions, "database_version")

			// Verify the instance name matches what we expected
			assert.Equal(t, expectedInstanceName, actualInstanceName)

			// Verify the instance exists in GCP
			gcp.AssertCloudSQLInstanceExists(t, projectID, actualInstanceName)

			// Verify the database version from Terraform output matches what we deployed
			assert.Equal(t, tt.databaseVersion, actualDatabaseVersion)

			// Verify the database version from the GCP API matches what we deployed
			actualVersionFromAPI := gcp.GetCloudSQLInstanceDatabaseVersion(t, projectID, actualInstanceName)
			assert.Equal(t, tt.databaseVersion, actualVersionFromAPI)
		})
	}
}
