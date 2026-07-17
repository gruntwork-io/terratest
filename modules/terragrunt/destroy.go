package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// DestroyAllContext runs terragrunt run --all destroy with the given options and returns stdout.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func DestroyAllContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := DestroyAllContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// DestroyAllContextE runs terragrunt run --all -- destroy with the given options and returns stdout.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func DestroyAllContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	args := BuildRunArgs([]string{"--all"}, []string{"destroy", "-auto-approve", "-input=false"})

	return runTerragruntCommandE(t, ctx, options, "run", args...)
}

// DestroyContext runs terragrunt run destroy for a single unit and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func DestroyContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := DestroyContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// DestroyContextE runs terragrunt run -- destroy for a single unit and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func DestroyContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	args := BuildRunArgs([]string{}, []string{"destroy", "-auto-approve", "-input=false"})

	return runTerragruntCommandE(t, ctx, options, "run", args...)
}
