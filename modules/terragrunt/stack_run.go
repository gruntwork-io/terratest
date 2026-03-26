package terragrunt //nolint:dupl // structural pattern for terragrunt command wrappers

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// StackRunContext calls terragrunt stack run and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func StackRunContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := StackRunContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// StackRunContextE calls terragrunt stack run and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func StackRunContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, ctx, options, "run")
}

// StackRun calls terragrunt stack run and returns stdout/stderr.
//
// Deprecated: Use [StackRunContext] instead.
func StackRun(t testing.TestingT, options *Options) string {
	return StackRunContext(t, context.Background(), options)
}

// StackRunE calls terragrunt stack run and returns stdout/stderr.
//
// Deprecated: Use [StackRunContextE] instead.
func StackRunE(t testing.TestingT, options *Options) (string, error) {
	return StackRunContextE(t, context.Background(), options)
}
