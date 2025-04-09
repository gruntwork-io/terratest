//go:build gcp
// +build gcp

// NOTE: We use build tags to differentiate GCP testing for better isolation and parallelism when executing our tests.

package test

import (
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
	projectId := gcp.GetGoogleProjectIDFromEnvVar(t)

	region := gcp.GetRandomRegion(t, projectId, RegionsThatSupportF1Micro, nil)

	randomValidGcpName := gcp.RandomValidGcpName()
	clusterSize := 3

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code instances located
		TerraformDir: exampleDir,

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"gcp_project_id": projectId,
			"gcp_region":     region,
			"cluster_name":   randomValidGcpName,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	instanceGroupName := terraform.Output(t, terraformOptions, "instance_group_name")

	instanceGroup := gcp.FetchRegionalInstanceGroup(t, projectId, region, instanceGroupName)

	// Validate that GetInstances() returns a non-zero number of Instances
	maxRetries := 40
	sleepBetweenRetries := 3 * time.Second

	retry.DoWithRetry(t, "Attempting to fetch Instances from Instance Group", maxRetries, sleepBetweenRetries, func() (string, error) {
		instances, err := instanceGroup.GetInstancesE(t, projectId)
		if err != nil {
			return "", fmt.Errorf("Failed to get Instances: %s", err)
		}

		if len(instances) != clusterSize {
			return "", fmt.Errorf("Expected to find exactly %d Compute Instances in Instance Group but found %d.", clusterSize, len(instances))
		}

		return "", nil
	})

	// Validate that we get the right number of IP addresses
	retry.DoWithRetry(t, "Attempting to fetch Public IP addresses from Instance Group", maxRetries, sleepBetweenRetries, func() (string, error) {
		ips, err := instanceGroup.GetPublicIpsE(t, projectId)
		if err != nil {
			return "", fmt.Errorf("Failed to get public IPs from Instance Group")
		}

		if len(ips) != clusterSize {
			return "", fmt.Errorf("Expected to get exactly %d public IP addresses but found %d.", clusterSize, len(ips))
		}

		return "", nil
	})
}
