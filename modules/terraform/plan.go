package terraform

import (
	"fmt"
	"os"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/assert"
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

// InitAndPlanAndShow runs terraform init, then terraform plan, and then terraform show with the given options, and
// returns the json output of the plan file. This will fail the test if there is an error in the command.
func InitAndPlanAndShow(t testing.TestingT, options *Options) string {
	jsonOut, err := InitAndPlanAndShowE(t, options)
	require.NoError(t, err)
	return jsonOut
}

// InitAndPlanAndShowE runs terraform init, then terraform plan, and then terraform show with the given options, and
// returns the json output of the plan file.
func InitAndPlanAndShowE(t testing.TestingT, options *Options) (string, error) {
	if options.PlanFilePath == "" {
		return "", PlanFilePathRequired
	}

	_, err := InitAndPlanE(t, options)
	if err != nil {
		return "", err
	}
	return ShowE(t, options)
}

// InitAndPlanAndShowWithStructNoLog runs InitAndPlanAndShowWithStruct without logging and also by allocating a
// temporary plan file destination that is discarded before returning the struct.
func InitAndPlanAndShowWithStructNoLogTempPlanFile(t testing.TestingT, options *Options) *PlanStruct {
	oldLogger := options.Logger
	options.Logger = logger.Discard
	defer func() { options.Logger = oldLogger }()

	tmpFile, err := os.CreateTemp("", "terratest-plan-file-")
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())
	defer require.NoError(t, os.Remove(tmpFile.Name()))

	options.PlanFilePath = tmpFile.Name()
	return InitAndPlanAndShowWithStruct(t, options)
}

// InitAndPlanAndShowWithStruct runs terraform init, then terraform plan, and then terraform show with the given
// options, and parses the json result into a go struct. This will fail the test if there is an error in the command.
func InitAndPlanAndShowWithStruct(t testing.TestingT, options *Options) *PlanStruct {
	plan, err := InitAndPlanAndShowWithStructE(t, options)
	require.NoError(t, err)
	return plan
}

// InitAndPlanAndShowWithStructE runs terraform init, then terraform plan, and then terraform show with the given options, and
// parses the json result into a go struct.
func InitAndPlanAndShowWithStructE(t testing.TestingT, options *Options) (*PlanStruct, error) {
	jsonOut, err := InitAndPlanAndShowE(t, options)
	if err != nil {
		return nil, err
	}
	return ParsePlanJSON(jsonOut)
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

	return GetExitCodeForTerraformCommandE(t, options, FormatArgs(options, "run-all", "plan", "--input=false",
		"--lock=true", "--detailed-exitcode")...)
}

// AssertTgPlanAllExitCode asserts the succuess (or failure) of a terragrunt run-all plan.
// On success, terragrunt will exit 0 on a plan that has previously been applied (has state)
// and exit with 2 for plans that have never been applied when ran with `-detailed-exitcode`.
func AssertTgPlanAllExitCode(t testing.TestingT, exitCode int, assertTrue bool) {

	validExitCodes := map[int]bool{
		0: true,
		2: true,
	}

	_, hasKey := validExitCodes[exitCode]
	if assertTrue {
		assert.True(t, hasKey)
	} else {
		assert.False(t, hasKey)
	}
}

// Custom errors

var (
	PlanFilePathRequired = fmt.Errorf("You must set PlanFilePath on options struct to use this function.")
)
