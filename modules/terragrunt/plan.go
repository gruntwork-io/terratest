package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
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

// Plan runs terragrunt plan for a single unit and returns stdout/stderr.
func Plan(t testing.TestingT, options *Options) string {
	out, err := PlanE(t, options)
	require.NoError(t, err)
	return out
}

// PlanE runs terragrunt plan for a single unit and returns stdout/stderr.
func PlanE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntCommandE(t, options, "plan", "-input=false", "-lock=false")
}

// PlanExitCode runs terragrunt plan for a single unit and returns the detailed exit code.
// This will fail the test if there is an error in the command.
func PlanExitCode(t testing.TestingT, options *Options) int {
	exitCode, err := PlanExitCodeE(t, options)
	require.NoError(t, err)
	return exitCode
}

// PlanExitCodeE runs terragrunt plan for a single unit and returns the detailed exit code.
func PlanExitCodeE(t testing.TestingT, options *Options) (int, error) {
	return getExitCodeForTerragruntCommandE(t, options, "plan", "-input=false",
		"-lock=true", "-detailed-exitcode")
}

// InitAndPlan runs terragrunt init followed by plan for a single unit and returns the plan stdout/stderr.
func InitAndPlan(t testing.TestingT, options *Options) string {
	out, err := InitAndPlanE(t, options)
	require.NoError(t, err)
	return out
}

// InitAndPlanE runs terragrunt init followed by plan for a single unit and returns the plan stdout/stderr.
func InitAndPlanE(t testing.TestingT, options *Options) (string, error) {
	if _, err := InitE(t, options); err != nil {
		return "", err
	}
	return PlanE(t, options)
}
