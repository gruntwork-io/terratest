package test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestTerraformAwsSsmExample(t *testing.T) {
	t.Parallel()
	region := aws.GetRandomStableRegion(t, nil, nil)

	// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
	instanceType := aws.GetRecommendedInstanceType(t, region, []string{"t2.micro", "t3.micro"})

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/terraform-aws-ssm-example",
		Vars: map[string]interface{}{
			"region":        region,
			"instance_type": instanceType,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	terraform.InitAndApply(t, terraformOptions)

	instanceID := terraform.Output(t, terraformOptions, "instance_id")
	timeout := 3 * time.Minute

	aws.WaitForSsmInstance(t, region, instanceID, timeout)

	result := aws.CheckSsmCommand(t, region, instanceID, "echo Hello, World", timeout)
	require.Equal(t, result.Stdout, "Hello, World\n")
	require.Equal(t, result.Stderr, "")
	require.Equal(t, int64(0), result.ExitCode)

	result, err := aws.CheckSsmCommandE(t, region, instanceID, "cat /wrong/file", timeout)
	require.Error(t, err)
	require.Equal(t, "Failed", err.Error())
	require.Equal(t, "cat: /wrong/file: No such file or directory\nfailed to run commands: exit status 1", result.Stderr)
	require.Equal(t, "", result.Stdout)
	require.Equal(t, int64(1), result.ExitCode)
}
