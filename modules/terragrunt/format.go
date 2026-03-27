package terragrunt //nolint:dupl // structural pattern for terragrunt command wrappers

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
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

// FormatAll runs terragrunt hcl format to format all terragrunt.hcl files and returns stdout/stderr.
//
// Deprecated: Use [FormatAllContext] instead.
func FormatAll(t testing.TestingT, options *Options) string {
	return FormatAllContext(t, context.Background(), options)
}

// FormatAllE runs terragrunt hcl format to format all terragrunt.hcl files and returns stdout/stderr.
//
// Deprecated: Use [FormatAllContextE] instead.
func FormatAllE(t testing.TestingT, options *Options) (string, error) {
	return FormatAllContextE(t, context.Background(), options)
}
