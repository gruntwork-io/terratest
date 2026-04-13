//go:build aws

package test_test

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/require"
)

func TestTerraformAwsSsmExample(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegionContext(t, t.Context(), nil, nil)

	// Some AWS regions are missing certain instance types, so pick an available type based on the region we picked
	instanceType := aws.GetRecommendedInstanceTypeContext(t, t.Context(), region, []string{"t2.micro, t3.micro", "t2.small", "t3.small"})

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/terraform-aws-ssm-example",
		Vars: map[string]interface{}{
			"region":        region,
			"instance_type": instanceType,
		},
	})

	defer terraform.DestroyContext(t, t.Context(), terraformOptions)

	terraform.InitAndApplyContext(t, t.Context(), terraformOptions)

	instanceID := terraform.OutputContext(t, t.Context(), terraformOptions, "instance_id")
	timeout := 3 * time.Minute

	aws.WaitForSsmInstanceContext(t, t.Context(), region, instanceID, timeout)

	result := aws.CheckSsmCommandContext(t, t.Context(), region, instanceID, "echo Hello, World", timeout)
	require.Equal(t, "Hello, World\n", result.Stdout)
	require.Empty(t, result.Stderr)
	require.Equal(t, int64(0), result.ExitCode)

	result, err := aws.CheckSsmCommandContextE(t, t.Context(), region, instanceID, "cat /wrong/file", timeout)
	require.Error(t, err)
	require.Equal(t, "Failed", err.Error())
	require.Equal(t, "cat: /wrong/file: No such file or directory\nfailed to run commands: exit status 1", result.Stderr)
	require.Empty(t, result.Stdout)
	require.Equal(t, int64(1), result.ExitCode)
}
