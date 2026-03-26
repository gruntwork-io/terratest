package terraform

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Destroy runs terraform destroy with the given options and return stdout/stderr.
//
// Deprecated: Use [DestroyContext] instead.
func Destroy(t testing.TestingT, options *Options) string {
	return DestroyContext(t, context.Background(), options)
}

// DestroyContext runs terraform destroy with the given options and returns stdout/stderr. The provided context is
// passed through to the underlying command execution, allowing for timeout and cancellation control.
func DestroyContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := DestroyContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// DestroyE runs terraform destroy with the given options and return stdout/stderr.
//
// Deprecated: Use [DestroyContextE] instead.
func DestroyE(t testing.TestingT, options *Options) (string, error) {
	return DestroyContextE(t, context.Background(), options)
}

// DestroyContextE runs terraform destroy with the given options and returns stdout/stderr. The provided context is
// passed through to the underlying command execution, allowing for timeout and cancellation control.
func DestroyContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return RunTerraformCommandContextE(t, ctx, options, FormatArgs(options, prepend(options.ExtraArgs.Destroy, "destroy", "-auto-approve", "-input=false")...)...)
}
