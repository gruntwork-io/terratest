package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
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

// DestroyAll runs terragrunt run --all destroy with the given options and returns stdout.
//
// Deprecated: Use [DestroyAllContext] instead.
func DestroyAll(t testing.TestingT, options *Options) string {
	return DestroyAllContext(t, context.Background(), options)
}

// DestroyAllE runs terragrunt run --all -- destroy with the given options and returns stdout.
//
// Deprecated: Use [DestroyAllContextE] instead.
func DestroyAllE(t testing.TestingT, options *Options) (string, error) {
	return DestroyAllContextE(t, context.Background(), options)
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

// Destroy runs terragrunt run destroy for a single unit and returns stdout/stderr.
//
// Deprecated: Use [DestroyContext] instead.
func Destroy(t testing.TestingT, options *Options) string {
	return DestroyContext(t, context.Background(), options)
}

// DestroyE runs terragrunt run -- destroy for a single unit and returns stdout/stderr.
//
// Deprecated: Use [DestroyContextE] instead.
func DestroyE(t testing.TestingT, options *Options) (string, error) {
	return DestroyContextE(t, context.Background(), options)
}
