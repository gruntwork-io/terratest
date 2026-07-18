package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
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
