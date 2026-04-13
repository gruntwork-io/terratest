//go:build aws

package test_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/terraform"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
)

// An example of how to test the Terraform module in examples/terraform-aws-ecs-example using Terratest.
func TestTerraformAwsEcsExample(t *testing.T) {
	t.Parallel()

	expectedClusterName := "terratest-aws-ecs-example-cluster-" + random.UniqueID()
	expectedServiceName := "terratest-aws-ecs-example-service-" + random.UniqueID()

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegionContext(t, t.Context(), []string{"us-east-1", "eu-west-1"}, nil)

	// Construct the terraform options with default retryable errors to handle the most common retryable errors in
	// terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../examples/terraform-aws-ecs-example",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"cluster_name": expectedClusterName,
			"service_name": expectedServiceName,
			"region":       awsRegion,
		},
	})

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	// Run `terraform output` to get the value of an output variable
	taskDefinition := terraform.OutputContext(t, t.Context(), terraformOptions, "task_definition")

	// Look up the ECS cluster by name
	cluster := aws.GetEcsClusterContext(t, t.Context(), awsRegion, expectedClusterName)

	assert.Equal(t, int32(1), cluster.ActiveServicesCount)

	// Look up the ECS service by name
	service := aws.GetEcsServiceContext(t, t.Context(), awsRegion, expectedClusterName, expectedServiceName)

	assert.Equal(t, int32(0), service.DesiredCount)
	assert.Equal(t, types.LaunchTypeFargate, service.LaunchType)

	// Look up the ECS task definition by ARN
	task := aws.GetEcsTaskDefinitionContext(t, t.Context(), awsRegion, taskDefinition)

	assert.Equal(t, "256", awsSDK.ToString(task.Cpu))
	assert.Equal(t, "512", awsSDK.ToString(task.Memory))
	assert.Equal(t, types.NetworkModeAwsvpc, task.NetworkMode)
}
