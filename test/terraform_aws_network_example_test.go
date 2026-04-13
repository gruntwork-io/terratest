//go:build aws

package test_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// An example of how to test the Terraform module in examples/terraform-aws-network-example using Terratest.
func TestTerraformAwsNetworkExample(t *testing.T) {
	t.Parallel()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegionContext(t, t.Context(), nil, nil)

	// Give the VPC and the subnets correct CIDRs
	vpcCidr := "10.10.0.0/16"
	privateSubnetCidr := "10.10.1.0/24"
	publicSubnetCidr := "10.10.2.0/24"

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-aws-network-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"main_vpc_cidr":       vpcCidr,
			"private_subnet_cidr": privateSubnetCidr,
			"public_subnet_cidr":  publicSubnetCidr,
			"aws_region":          awsRegion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// Run `terraform output` to get the value of an output variable
	publicSubnetID := terraform.OutputContext(t, t.Context(), terraformOptions, "public_subnet_id")
	privateSubnetID := terraform.OutputContext(t, t.Context(), terraformOptions, "private_subnet_id")
	vpcID := terraform.OutputContext(t, t.Context(), terraformOptions, "main_vpc_id")

	subnets := aws.GetSubnetsForVpcContext(t, t.Context(), vpcID, awsRegion)

	require.Len(t, subnets, 2)
	// Verify if the network that is supposed to be public is really public
	assert.True(t, aws.IsPublicSubnetContext(t, t.Context(), publicSubnetID, awsRegion))
	// Verify if the network that is supposed to be private is really private
	assert.False(t, aws.IsPublicSubnetContext(t, t.Context(), privateSubnetID, awsRegion))
}
