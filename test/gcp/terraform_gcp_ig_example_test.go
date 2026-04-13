//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package test_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/gcp"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

func TestTerraformGcpInstanceGroupExample(t *testing.T) {
	t.Parallel()

	exampleDir := test_structure.CopyTerraformFolderToTemp(t, "../../", "examples/terraform-gcp-ig-example")

	// Setup values for our Terraform apply
	projectID := gcp.GetGoogleProjectIDFromEnvVar(t)

	region := gcp.GetRandomRegionContext(t, t.Context(), projectID, RegionsThatSupportF1Micro, nil)

	randomValidGcpName := gcp.RandomValidGCPName()
	clusterSize := 3

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code instances located
		TerraformDir: exampleDir,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"gcp_project_id": projectID,
			"gcp_region":     region,
			"cluster_name":   randomValidGcpName,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	instanceGroupName := terraform.OutputContext(t, t.Context(), terraformOptions, "instance_group_name")

	instanceGroup := gcp.FetchRegionalInstanceGroupContext(t, t.Context(), projectID, region, instanceGroupName)

	// Validate that GetInstances() returns a non-zero number of Instances
	maxRetries := 100
	sleepBetweenRetries := 3 * time.Second

	retry.DoWithRetryContext(t, t.Context(), "Attempting to fetch Instances from Instance Group", maxRetries, sleepBetweenRetries, func() (string, error) {
		instances, err := instanceGroup.GetInstancesContextE(t, t.Context(), projectID)
		if err != nil {
			return "", fmt.Errorf("Failed to get Instances: %w", err)
		}

		if len(instances) != clusterSize {
			return "", fmt.Errorf("Expected to find exactly %d Compute Instances in Instance Group but found %d.", clusterSize, len(instances))
		}

		return "", nil
	})

	// Validate that we get the right number of IP addresses
	retry.DoWithRetryContext(t, t.Context(), "Attempting to fetch Public IP addresses from Instance Group", maxRetries, sleepBetweenRetries, func() (string, error) {
		ips, err := instanceGroup.GetPublicIPsContextE(t, t.Context(), projectID)
		if err != nil {
			return "", errors.New("Failed to get public IPs from Instance Group")
		}

		if len(ips) != clusterSize {
			return "", fmt.Errorf("Expected to get exactly %d public IP addresses but found %d.", clusterSize, len(ips))
		}

		return "", nil
	})
}
