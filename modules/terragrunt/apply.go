package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// ApplyAll runs terragrunt run-all apply and returns stdout/stderr.
// Note: This method does NOT call destroy; caller is responsible for cleanup.
func ApplyAll(t testing.TestingT, options *Options) string {
	out, err := ApplyAllE(t, options)
	require.NoError(t, err)
	return out
}

// ApplyAllE runs terragrunt run-all apply and returns stdout/stderr.
// Note: This method does NOT call destroy; caller is responsible for cleanup.
func ApplyAllE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntCommandE(t, options, "run-all", formatArgs(options, prepend(options.ExtraArgs.Apply, "apply", "-input=false", "-auto-approve")...)...)
}
