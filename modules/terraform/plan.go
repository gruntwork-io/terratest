package terraform

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// InitAndPlan runs terraform init and plan with the given options and returns stdout/stderr from the plan command.
// This will fail the test if there is an error in the command.
func InitAndPlan(t testing.TestingT, options *Options) string {
	out, err := InitAndPlanE(t, options)
	require.NoError(t, err)
	return out
}

// InitAndPlanE runs terraform init and plan with the given options and returns stdout/stderr from the plan command.
func InitAndPlanE(t testing.TestingT, options *Options) (string, error) {
	if _, err := InitE(t, options); err != nil {
		return "", err
	}

	if _, err := GetE(t, options); err != nil {
		return "", err
	}

	return PlanE(t, options)
}

// Plan runs terraform plan with the given options and returns stdout/stderr.
// This will fail the test if there is an error in the command.
func Plan(t testing.TestingT, options *Options) string {
	out, err := PlanE(t, options)
	require.NoError(t, err)
	return out
}

// PlanE runs terraform plan with the given options and returns stdout/stderr.
func PlanE(t testing.TestingT, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, FormatArgs(options, "plan", "-input=false", "-lock=false")...)
}

// InitAndPlanWithExitCode runs terraform init and plan with the given options and returns exitcode for the plan command.
// This will fail the test if there is an error in the command.
func InitAndPlanWithExitCode(t testing.TestingT, options *Options) int {
	exitCode, err := InitAndPlanWithExitCodeE(t, options)
	require.NoError(t, err)
	return exitCode
}

// InitAndPlanWithExitCodeE runs terraform init and plan with the given options and returns exitcode for the plan command.
func InitAndPlanWithExitCodeE(t testing.TestingT, options *Options) (int, error) {
	if _, err := InitE(t, options); err != nil {
		return DefaultErrorExitCode, err
	}

	return PlanExitCodeE(t, options)
}

// PlanExitCode runs terraform plan with the given options and returns the detailed exitcode.
// This will fail the test if there is an error in the command.
func PlanExitCode(t testing.TestingT, options *Options) int {
	exitCode, err := PlanExitCodeE(t, options)
	require.NoError(t, err)
	return exitCode
}

// PlanExitCodeE runs terraform plan with the given options and returns the detailed exitcode.
func PlanExitCodeE(t testing.TestingT, options *Options) (int, error) {
	return GetExitCodeForTerraformCommandE(t, options, FormatArgs(options, "plan", "-input=false", "-detailed-exitcode")...)
}

// TgPlanAllExitCode runs terragrunt plan-all with the given options and returns the detailed exitcode.
// This will fail the test if there is an error in the command.
func TgPlanAllExitCode(t testing.TestingT, options *Options) int {
	exitCode, err := TgPlanAllExitCodeE(t, options)
	require.NoError(t, err)
	return exitCode
}

// TgPlanAllExitCodeE runs terragrunt plan-all with the given options and returns the detailed exitcode.
func TgPlanAllExitCodeE(t testing.TestingT, options *Options) (int, error) {
	if options.TerraformBinary != "terragrunt" {
		return 1, fmt.Errorf("terragrunt must be set as TerraformBinary to use this method")
	}

	return GetExitCodeForTerraformCommandE(t, options, FormatArgs(options, "plan-all", "--input=false", "--lock=true", "--detailed-exitcode")...)
}
