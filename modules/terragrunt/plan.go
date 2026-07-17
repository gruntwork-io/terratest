package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
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
