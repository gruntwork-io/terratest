package terragrunt //nolint:dupl // structural pattern for terragrunt command wrappers

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// StackGenerateContext calls terragrunt stack generate and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func StackGenerateContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := StackGenerateContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// StackGenerateContextE calls terragrunt stack generate and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func StackGenerateContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, ctx, options, "generate")
}

// StackGenerate calls terragrunt stack generate and returns stdout/stderr.
//
// Deprecated: Use [StackGenerateContext] instead.
func StackGenerate(t testing.TestingT, options *Options) string {
	return StackGenerateContext(t, context.Background(), options)
}

// StackGenerateE calls terragrunt stack generate and returns stdout/stderr.
//
// Deprecated: Use [StackGenerateContextE] instead.
func StackGenerateE(t testing.TestingT, options *Options) (string, error) {
	return StackGenerateContextE(t, context.Background(), options)
}
