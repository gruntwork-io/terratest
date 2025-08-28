package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// PlanAllExitCode runs terragrunt plan-all with the given options and returns the detailed exitcode.
// This will fail the test if there is an error in the command.
func PlanAllExitCode(t testing.TestingT, options *Options) int {
	exitCode, err := PlanAllExitCodeE(t, options)
	require.NoError(t, err)
	return exitCode
}

// PlanAllExitCodeE runs terragrunt plan-all with the given options and returns the detailed exitcode.
func PlanAllExitCodeE(t testing.TestingT, options *Options) (int, error) {
	return getExitCodeForTerraformCommandE(t, options, formatArgs(options, prepend(options.ExtraArgs.Plan, "run-all", "plan", "--input=false", "--lock=true", "--detailed-exitcode")...)...)
}
