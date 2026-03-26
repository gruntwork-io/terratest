package terragrunt //nolint:dupl // structural pattern for terragrunt command wrappers

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// StackCleanContext calls terragrunt stack clean to remove the .terragrunt-stack directory.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This command cleans up the generated stack files created by stack generate
// or stack run.
func StackCleanContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := StackCleanContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// StackCleanContextE calls terragrunt stack clean to remove the .terragrunt-stack directory.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This command cleans up the generated stack files created by stack generate
// or stack run.
func StackCleanContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, ctx, options, "clean")
}

// StackClean calls terragrunt stack clean to remove the .terragrunt-stack directory.
// This command cleans up the generated stack files created by stack generate or stack run.
//
// Deprecated: Use [StackCleanContext] instead.
func StackClean(t testing.TestingT, options *Options) string {
	return StackCleanContext(t, context.Background(), options)
}

// StackCleanE calls terragrunt stack clean to remove the .terragrunt-stack directory.
// This command cleans up the generated stack files created by stack generate or stack run.
//
// Deprecated: Use [StackCleanContextE] instead.
func StackCleanE(t testing.TestingT, options *Options) (string, error) {
	return StackCleanContextE(t, context.Background(), options)
}
