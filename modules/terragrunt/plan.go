package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing/v2"
	"github.com/stretchr/testify/require"
)

// PlanAllExitCode runs terragrunt plan --all with the given options and returns the detailed exit code.
// This will fail the test if there is an error in the command.
func PlanAllExitCode(t testing.TestingT, options *Options) int {
	exitCode, err := PlanAllExitCodeE(t, options)
	require.NoError(t, err)
	return exitCode
}

// PlanAllExitCodeE runs terragrunt plan --all with the given options and returns the detailed exit code.
func PlanAllExitCodeE(t testing.TestingT, options *Options) (int, error) {
	return getExitCodeForTerragruntCommandE(t, options, "plan", "--all", "-input=false",
		"-lock=true", "-detailed-exitcode")
}
