package terraform

import (
	"context"
	"errors"
	"os"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// InitAndPlanContext runs terraform init and plan with the given options and returns stdout/stderr from the plan command.
// The context argument can be used for cancellation or timeout control.
// This will fail the test if there is an error in the command.
func InitAndPlanContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := InitAndPlanContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// InitAndPlanContextE runs terraform init and plan with the given options and returns stdout/stderr from the plan
// command. The context argument can be used for cancellation or timeout control.
func InitAndPlanContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	if _, err := InitContextE(t, ctx, options); err != nil {
		return "", err
	}

	return PlanContextE(t, ctx, options)
}

// PlanContext runs terraform plan with the given options and returns stdout/stderr.
// The context argument can be used for cancellation or timeout control.
// This will fail the test if there is an error in the command.
func PlanContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := PlanContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// PlanContextE runs terraform plan with the given options and returns stdout/stderr.
// The context argument can be used for cancellation or timeout control.
func PlanContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return RunTerraformCommandContextE(t, ctx, options, FormatArgs(options, prepend(options.ExtraArgs.Plan, "plan", "-input=false", "-lock=false")...)...)
}

// InitAndPlanAndShowContext runs terraform init, then terraform plan, and then terraform show with the given options,
// and returns the json output of the plan file. The context argument can be used for cancellation or timeout control.
// This will fail the test if there is an error in the command.
func InitAndPlanAndShowContext(t testing.TestingT, ctx context.Context, options *Options) string {
	jsonOut, err := InitAndPlanAndShowContextE(t, ctx, options)
	require.NoError(t, err)

	return jsonOut
}

// InitAndPlanAndShowContextE runs terraform init, then terraform plan, and then terraform show with the given options,
// and returns the json output of the plan file. The context argument can be used for cancellation or timeout control.
func InitAndPlanAndShowContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	if options.PlanFilePath == "" {
		return "", PlanFilePathRequired
	}

	_, err := InitAndPlanContextE(t, ctx, options)
	if err != nil {
		return "", err
	}

	return ShowContextE(t, ctx, options)
}

// InitAndPlanAndShowWithStructNoLogTempPlanFileContext runs InitAndPlanAndShowWithStructContext without logging and also
// by allocating a temporary plan file destination that is discarded before returning the struct. The context argument
// can be used for cancellation or timeout control.
func InitAndPlanAndShowWithStructNoLogTempPlanFileContext(t testing.TestingT, ctx context.Context, options *Options) *PlanStruct {
	oldLogger := options.Logger
	options.Logger = logger.Discard

	defer func() { options.Logger = oldLogger }()

	tmpFile, err := os.CreateTemp("", "terratest-plan-file-")
	require.NoError(t, err)

	require.NoError(t, tmpFile.Close())
	defer require.NoError(t, os.Remove(tmpFile.Name()))

	options.PlanFilePath = tmpFile.Name()

	return InitAndPlanAndShowWithStructContext(t, ctx, options)
}

// InitAndPlanAndShowWithStructContext runs terraform init, then terraform plan, and then terraform show with the given
// options, and parses the json result into a go struct. The context argument can be used for cancellation or timeout
// control. This will fail the test if there is an error in the command.
func InitAndPlanAndShowWithStructContext(t testing.TestingT, ctx context.Context, options *Options) *PlanStruct {
	plan, err := InitAndPlanAndShowWithStructContextE(t, ctx, options)
	require.NoError(t, err)

	return plan
}

// InitAndPlanAndShowWithStructContextE runs terraform init, then terraform plan, and then terraform show with the
// given options, and parses the json result into a go struct. The context argument can be used for cancellation or
// timeout control.
func InitAndPlanAndShowWithStructContextE(t testing.TestingT, ctx context.Context, options *Options) (*PlanStruct, error) {
	jsonOut, err := InitAndPlanAndShowContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	return ParsePlanJSON(jsonOut)
}

// InitAndPlanWithExitCodeContext runs terraform init and plan with the given options and returns the exitcode for the
// plan command. The context argument can be used for cancellation or timeout control.
// This will fail the test if there is an error in the command.
func InitAndPlanWithExitCodeContext(t testing.TestingT, ctx context.Context, options *Options) int {
	exitCode, err := InitAndPlanWithExitCodeContextE(t, ctx, options)
	require.NoError(t, err)

	return exitCode
}

// InitAndPlanWithExitCodeContextE runs terraform init and plan with the given options and returns the exitcode for the
// plan command. The context argument can be used for cancellation or timeout control.
func InitAndPlanWithExitCodeContextE(t testing.TestingT, ctx context.Context, options *Options) (int, error) {
	if _, err := InitContextE(t, ctx, options); err != nil {
		return DefaultErrorExitCode, err
	}

	return PlanExitCodeContextE(t, ctx, options)
}

// PlanExitCodeContext runs terraform plan with the given options and returns the detailed exitcode.
// The context argument can be used for cancellation or timeout control.
// This will fail the test if there is an error in the command.
func PlanExitCodeContext(t testing.TestingT, ctx context.Context, options *Options) int {
	exitCode, err := PlanExitCodeContextE(t, ctx, options)
	require.NoError(t, err)

	return exitCode
}

// PlanExitCodeContextE runs terraform plan with the given options and returns the detailed exitcode.
// The context argument can be used for cancellation or timeout control.
func PlanExitCodeContextE(t testing.TestingT, ctx context.Context, options *Options) (int, error) {
	return GetExitCodeForTerraformCommandContextE(t, ctx, options, FormatArgs(options, prepend(options.ExtraArgs.Plan, "plan", "-input=false", "-detailed-exitcode")...)...)
}

// InitAndPlan runs terraform init and plan with the given options and returns stdout/stderr from the plan command.
// This will fail the test if there is an error in the command.
//
// Deprecated: Use [InitAndPlanContext] instead.
func InitAndPlan(t testing.TestingT, options *Options) string {
	return InitAndPlanContext(t, context.Background(), options)
}

// InitAndPlanE runs terraform init and plan with the given options and returns stdout/stderr from the plan command.
//
// Deprecated: Use [InitAndPlanContextE] instead.
func InitAndPlanE(t testing.TestingT, options *Options) (string, error) {
	return InitAndPlanContextE(t, context.Background(), options)
}

// Plan runs terraform plan with the given options and returns stdout/stderr.
// This will fail the test if there is an error in the command.
//
// Deprecated: Use [PlanContext] instead.
func Plan(t testing.TestingT, options *Options) string {
	return PlanContext(t, context.Background(), options)
}

// PlanE runs terraform plan with the given options and returns stdout/stderr.
//
// Deprecated: Use [PlanContextE] instead.
func PlanE(t testing.TestingT, options *Options) (string, error) {
	return PlanContextE(t, context.Background(), options)
}

// InitAndPlanAndShow runs terraform init, then terraform plan, and then terraform show with the given options, and
// returns the json output of the plan file. This will fail the test if there is an error in the command.
//
// Deprecated: Use [InitAndPlanAndShowContext] instead.
func InitAndPlanAndShow(t testing.TestingT, options *Options) string {
	return InitAndPlanAndShowContext(t, context.Background(), options)
}

// InitAndPlanAndShowE runs terraform init, then terraform plan, and then terraform show with the given options, and
// returns the json output of the plan file.
//
// Deprecated: Use [InitAndPlanAndShowContextE] instead.
func InitAndPlanAndShowE(t testing.TestingT, options *Options) (string, error) {
	return InitAndPlanAndShowContextE(t, context.Background(), options)
}

// InitAndPlanAndShowWithStructNoLogTempPlanFile runs InitAndPlanAndShowWithStruct without logging and also by allocating
// a temporary plan file destination that is discarded before returning the struct.
//
// Deprecated: Use [InitAndPlanAndShowWithStructNoLogTempPlanFileContext] instead.
func InitAndPlanAndShowWithStructNoLogTempPlanFile(t testing.TestingT, options *Options) *PlanStruct {
	return InitAndPlanAndShowWithStructNoLogTempPlanFileContext(t, context.Background(), options)
}

// InitAndPlanAndShowWithStruct runs terraform init, then terraform plan, and then terraform show with the given
// options, and parses the json result into a go struct. This will fail the test if there is an error in the command.
//
// Deprecated: Use [InitAndPlanAndShowWithStructContext] instead.
func InitAndPlanAndShowWithStruct(t testing.TestingT, options *Options) *PlanStruct {
	return InitAndPlanAndShowWithStructContext(t, context.Background(), options)
}

// InitAndPlanAndShowWithStructE runs terraform init, then terraform plan, and then terraform show with the given options, and
// parses the json result into a go struct.
//
// Deprecated: Use [InitAndPlanAndShowWithStructContextE] instead.
func InitAndPlanAndShowWithStructE(t testing.TestingT, options *Options) (*PlanStruct, error) {
	return InitAndPlanAndShowWithStructContextE(t, context.Background(), options)
}

// InitAndPlanWithExitCode runs terraform init and plan with the given options and returns exitcode for the plan command.
// This will fail the test if there is an error in the command.
//
// Deprecated: Use [InitAndPlanWithExitCodeContext] instead.
func InitAndPlanWithExitCode(t testing.TestingT, options *Options) int {
	return InitAndPlanWithExitCodeContext(t, context.Background(), options)
}

// InitAndPlanWithExitCodeE runs terraform init and plan with the given options and returns exitcode for the plan command.
//
// Deprecated: Use [InitAndPlanWithExitCodeContextE] instead.
func InitAndPlanWithExitCodeE(t testing.TestingT, options *Options) (int, error) {
	return InitAndPlanWithExitCodeContextE(t, context.Background(), options)
}

// PlanExitCode runs terraform plan with the given options and returns the detailed exitcode.
// This will fail the test if there is an error in the command.
//
// Deprecated: Use [PlanExitCodeContext] instead.
func PlanExitCode(t testing.TestingT, options *Options) int {
	return PlanExitCodeContext(t, context.Background(), options)
}

// PlanExitCodeE runs terraform plan with the given options and returns the detailed exitcode.
//
// Deprecated: Use [PlanExitCodeContextE] instead.
func PlanExitCodeE(t testing.TestingT, options *Options) (int, error) {
	return PlanExitCodeContextE(t, context.Background(), options)
}

// Custom errors

var (
	// PlanFilePathRequired is returned when PlanFilePath is not set on Options but is required by the called function.
	PlanFilePathRequired = errors.New("you must set PlanFilePath on options struct to use this function") //nolint:staticcheck // preserving existing public variable name
)
