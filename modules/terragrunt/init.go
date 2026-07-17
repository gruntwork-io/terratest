package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/formatting"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// InitContext calls terragrunt run init and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func InitContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := InitContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// InitContextE calls terragrunt run -- init and returns stdout/stderr.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func InitContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	args := BuildRunArgs([]string{}, append([]string{"init"}, initArgs(options)...))

	return runTerragruntCommandE(t, ctx, options, "run", args...)
}

// initArgs builds the argument list for terragrunt init command.
// This function handles complex configuration that requires special formatting.
func initArgs(options *Options) []string {
	backendArgs := formatting.FormatBackendConfigAsArgs(options.BackendConfig)
	pluginArgs := formatting.FormatPluginDirAsArgs(options.PluginDir)

	args := make([]string, 0, len(backendArgs)+len(pluginArgs))
	args = append(args, backendArgs...)
	args = append(args, pluginArgs...)

	return args
}
