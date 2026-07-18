package terraform

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// DestroyContext runs terraform destroy with the given options and returns stdout/stderr. The provided context is
// passed through to the underlying command execution, allowing for timeout and cancellation control.
func DestroyContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := DestroyContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// DestroyContextE runs terraform destroy with the given options and returns stdout/stderr. The provided context is
// passed through to the underlying command execution, allowing for timeout and cancellation control.
func DestroyContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return RunTerraformCommandContextE(t, ctx, options, FormatArgs(options, prepend(options.ExtraArgs.Destroy, "destroy", "-auto-approve", "-input=false")...)...)
}
