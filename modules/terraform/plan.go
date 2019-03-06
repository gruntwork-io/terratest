package terraform

import (
	"testing"
)

// InitAndPlan runs terraform init and plan with the given options and return stdout/stderr from the apply command.
func InitAndPlan(t *testing.T, options *Options) int {
	exitCode, err := InitAndPlanE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return exitCode
}

// InitAndPlanE runs terraform init and plan with the given options and return stdout/stderr from the apply command.
func InitAndPlanE(t *testing.T, options *Options) (int, error) {
	if _, err := InitE(t, options); err != nil {
		return DefaultErrorExitCode, err
	}

	return PlanExitCodeE(t, options)
}

// PlanExitCode runs terraform apply with the given options and returns the detailed exitcode.
func PlanExitCode(t *testing.T, options *Options) int {
	exitCode, err := PlanExitCodeE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return exitCode
}

// PlanExitCodeE runs terraform apply with the given options and returns the detailed exitcode.
func PlanExitCodeE(t *testing.T, options *Options) (int, error) {
	return GetExitCodeForTerraformCommandE(t, options, FormatArgs(options, "plan", "-input=false", "-detailed-exitcode")...)
}
