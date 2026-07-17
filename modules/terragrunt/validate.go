package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// ValidateAllContext runs terragrunt run --all validate with the given options and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func ValidateAllContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := ValidateAllContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// ValidateAllContextE runs terragrunt run --all -- validate with the given options and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func ValidateAllContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	args := BuildRunArgs([]string{"--all"}, []string{"validate"})

	return runTerragruntCommandE(t, ctx, options, "run", args...)
}

// ValidateContext runs terragrunt run validate for a single unit and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func ValidateContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := ValidateContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// ValidateContextE runs terragrunt run -- validate for a single unit and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func ValidateContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	args := BuildRunArgs([]string{}, []string{"validate"})

	return runTerragruntCommandE(t, ctx, options, "run", args...)
}

// InitAndValidateContext runs terragrunt init followed by validate for a single unit and returns the validate stdout/stderr.
// The provided context is passed through to both the init and validate command executions.
func InitAndValidateContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := InitAndValidateContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// InitAndValidateContextE runs terragrunt init followed by validate for a single unit and returns the validate stdout/stderr.
// The provided context is passed through to both the init and validate command executions.
func InitAndValidateContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	if _, err := InitContextE(t, ctx, options); err != nil {
		return "", err
	}

	return ValidateContextE(t, ctx, options)
}
