package terraform

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Validate calls terraform validate and returns stdout/stderr.
//
// Deprecated: Use [ValidateContext] instead.
func Validate(t testing.TestingT, options *Options) string {
	return ValidateContext(t, context.Background(), options)
}

// ValidateContext calls terraform validate and returns stdout/stderr. The provided context is passed through to the
// underlying command execution, allowing for timeout and cancellation control.
func ValidateContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := ValidateContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// ValidateE calls terraform validate and returns stdout/stderr.
//
// Deprecated: Use [ValidateContextE] instead.
func ValidateE(t testing.TestingT, options *Options) (string, error) {
	return ValidateContextE(t, context.Background(), options)
}

// ValidateContextE calls terraform validate and returns stdout/stderr. The provided context is passed through to the
// underlying command execution, allowing for timeout and cancellation control.
func ValidateContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return RunTerraformCommandContextE(t, ctx, options, FormatArgs(options, prepend(options.ExtraArgs.Validate, "validate")...)...)
}

// InitAndValidate runs terraform init and validate with the given options and returns stdout/stderr from the validate command.
// This will fail the test if there is an error in the command.
//
// Deprecated: Use [InitAndValidateContext] instead.
func InitAndValidate(t testing.TestingT, options *Options) string {
	return InitAndValidateContext(t, context.Background(), options)
}

// InitAndValidateContext runs terraform init and validate with the given options and returns stdout/stderr from the
// validate command. The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This will fail the test if there is an error in the command.
func InitAndValidateContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := InitAndValidateContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// InitAndValidateE runs terraform init and validate with the given options and returns stdout/stderr from the validate command.
//
// Deprecated: Use [InitAndValidateContextE] instead.
func InitAndValidateE(t testing.TestingT, options *Options) (string, error) {
	return InitAndValidateContextE(t, context.Background(), options)
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
