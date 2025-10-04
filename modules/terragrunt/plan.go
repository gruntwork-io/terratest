package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// PlanAllExitCode runs terragrunt run-all plan with --detailed-exitcode and returns the exit code.
// This will fail the test if there is an error in the command.
func PlanAllExitCode(t testing.TestingT, options *Options) int {
	exitCode, err := PlanAllExitCodeE(t, options)
	require.NoError(t, err)
	return exitCode
}

// PlanAllExitCodeE runs terragrunt run-all plan with --detailed-exitcode and returns the exit code.
func PlanAllExitCodeE(t testing.TestingT, options *Options) (int, error) {
	return getExitCodeForTerraformCommandE(t, options, formatArgs(options, prepend(options.ExtraArgs.Plan, "run-all", "plan", "--input=false", "--lock=true", "--detailed-exitcode")...)...)
}
