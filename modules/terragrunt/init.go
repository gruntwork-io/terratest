package terragrunt

import (
	"github.com/gruntwork-io/terratest/internal/lib/formatting"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Init calls terragrunt init and return stdout/stderr
func Init(t testing.TestingT, options *Options) string {
	out, err := InitE(t, options)
	require.NoError(t, err)
	return out
}

// InitE calls terragrunt init and return stdout/stderr
func InitE(t testing.TestingT, options *Options) (string, error) {
	// Use regular terragrunt init command (not terragrunt stack init)
	return runTerragruntCommandE(t, options, "init", initArgs(options)...)
}

// initArgs builds the argument list for terragrunt init command.
// This function handles complex configuration that requires special formatting.
func initArgs(options *Options) []string {
	var args []string

	// Add complex configuration that requires special formatting
	// These are terraform-specific arguments that need special formatting
	args = append(args, formatting.FormatBackendConfigAsArgs(options.BackendConfig)...)
	args = append(args, formatting.FormatPluginDirAsArgs(options.PluginDir)...)
	return args
}
