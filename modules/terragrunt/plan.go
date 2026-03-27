package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// PlanAllExitCodeContext runs terragrunt run --all plan with the given options and returns the detailed exit code.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This will fail the test if there is an error in the command.
func PlanAllExitCodeContext(t testing.TestingT, ctx context.Context, options *Options) int {
	exitCode, err := PlanAllExitCodeContextE(t, ctx, options)
	require.NoError(t, err)

	return exitCode
}

// PlanAllExitCodeContextE runs terragrunt run --all -- plan with the given options and returns the detailed exit code.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func PlanAllExitCodeContextE(t testing.TestingT, ctx context.Context, options *Options) (int, error) {
	args := BuildRunArgs([]string{"--all"}, []string{"plan", "-input=false", "-lock=true", "-detailed-exitcode"})

	return getExitCodeForTerragruntCommandE(t, ctx, options, append([]string{"run"}, args...)...)
}

// PlanAllExitCode runs terragrunt run --all plan with the given options and returns the detailed exit code.
// This will fail the test if there is an error in the command.
//
// Deprecated: Use [PlanAllExitCodeContext] instead.
func PlanAllExitCode(t testing.TestingT, options *Options) int {
	return PlanAllExitCodeContext(t, context.Background(), options)
}

// PlanAllExitCodeE runs terragrunt run --all -- plan with the given options and returns the detailed exit code.
//
// Deprecated: Use [PlanAllExitCodeContextE] instead.
func PlanAllExitCodeE(t testing.TestingT, options *Options) (int, error) {
	return PlanAllExitCodeContextE(t, context.Background(), options)
}

// PlanContext runs terragrunt run plan for a single unit and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func PlanContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := PlanContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// PlanContextE runs terragrunt run -- plan for a single unit and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. Uses -lock=false since plan is a read-only operation that does not need state locking.
func PlanContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	args := BuildRunArgs([]string{}, []string{"plan", "-input=false", "-lock=false"})

	return runTerragruntCommandE(t, ctx, options, "run", args...)
}

// Plan runs terragrunt run plan for a single unit and returns stdout/stderr.
//
// Deprecated: Use [PlanContext] instead.
func Plan(t testing.TestingT, options *Options) string {
	return PlanContext(t, context.Background(), options)
}

// PlanE runs terragrunt run -- plan for a single unit and returns stdout/stderr.
// Uses -lock=false since plan is a read-only operation that does not need state locking.
//
// Deprecated: Use [PlanContextE] instead.
func PlanE(t testing.TestingT, options *Options) (string, error) {
	return PlanContextE(t, context.Background(), options)
}

// PlanExitCodeContext runs terragrunt run plan for a single unit and returns the detailed exit code.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This will fail the test if there is an error in the command.
func PlanExitCodeContext(t testing.TestingT, ctx context.Context, options *Options) int {
	exitCode, err := PlanExitCodeContextE(t, ctx, options)
	require.NoError(t, err)

	return exitCode
}

// PlanExitCodeContextE runs terragrunt run -- plan for a single unit and returns the detailed exit code.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func PlanExitCodeContextE(t testing.TestingT, ctx context.Context, options *Options) (int, error) {
	args := BuildRunArgs([]string{}, []string{"plan", "-input=false", "-lock=true", "-detailed-exitcode"})

	return getExitCodeForTerragruntCommandE(t, ctx, options, append([]string{"run"}, args...)...)
}

// PlanExitCode runs terragrunt run plan for a single unit and returns the detailed exit code.
// This will fail the test if there is an error in the command.
//
// Deprecated: Use [PlanExitCodeContext] instead.
func PlanExitCode(t testing.TestingT, options *Options) int {
	return PlanExitCodeContext(t, context.Background(), options)
}

// PlanExitCodeE runs terragrunt run -- plan for a single unit and returns the detailed exit code.
//
// Deprecated: Use [PlanExitCodeContextE] instead.
func PlanExitCodeE(t testing.TestingT, options *Options) (int, error) {
	return PlanExitCodeContextE(t, context.Background(), options)
}

// InitAndPlanContext runs terragrunt init followed by plan for a single unit and returns the plan stdout/stderr.
// The provided context is passed through to both the init and plan command executions.
func InitAndPlanContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := InitAndPlanContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// InitAndPlanContextE runs terragrunt init followed by plan for a single unit and returns the plan stdout/stderr.
// The provided context is passed through to both the init and plan command executions.
func InitAndPlanContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	if _, err := InitContextE(t, ctx, options); err != nil {
		return "", err
	}

	return PlanContextE(t, ctx, options)
}

// InitAndPlan runs terragrunt init followed by plan for a single unit and returns the plan stdout/stderr.
//
// Deprecated: Use [InitAndPlanContext] instead.
func InitAndPlan(t testing.TestingT, options *Options) string {
	return InitAndPlanContext(t, context.Background(), options)
}

// InitAndPlanE runs terragrunt init followed by plan for a single unit and returns the plan stdout/stderr.
//
// Deprecated: Use [InitAndPlanContextE] instead.
func InitAndPlanE(t testing.TestingT, options *Options) (string, error) {
	return InitAndPlanContextE(t, context.Background(), options)
}
