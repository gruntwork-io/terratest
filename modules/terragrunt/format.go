package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// FormatAllContext runs terragrunt hcl format to format all terragrunt.hcl files and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func FormatAllContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := FormatAllContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// FormatAllContextE runs terragrunt hcl format to format all terragrunt.hcl files and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func FormatAllContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return runTerragruntCommandE(t, ctx, options, "hcl", "format")
}
