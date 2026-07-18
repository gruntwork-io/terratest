package terraform

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// ValidateContext calls terraform validate and returns stdout/stderr. The provided context is passed through to the
// underlying command execution, allowing for timeout and cancellation control.
func ValidateContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := ValidateContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// ValidateContextE calls terraform validate and returns stdout/stderr. The provided context is passed through to the
// underlying command execution, allowing for timeout and cancellation control.
func ValidateContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return RunTerraformCommandContextE(t, ctx, options, FormatArgs(options, prepend(options.ExtraArgs.Validate, "validate")...)...)
}

// InitAndValidateContext runs terraform init and validate with the given options and returns stdout/stderr from the
// validate command. The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This will fail the test if there is an error in the command.
func InitAndValidateContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := InitAndValidateContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// InitAndValidateContextE runs terraform init and validate with the given options and returns stdout/stderr from the
// validate command. The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func InitAndValidateContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	if _, err := InitContextE(t, ctx, options); err != nil {
		return "", err
	}

	return ValidateContextE(t, ctx, options)
}
