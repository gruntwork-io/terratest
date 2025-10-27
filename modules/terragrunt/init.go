package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/formatting"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// TgInit calls tg init and return stdout/stderr
func TgInit(t testing.TestingT, options *Options) string {
	out, err := TgInitE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TgInitE calls tg init and return stdout/stderr
func TgInitE(t testing.TestingT, options *Options) (string, error) {
	// Use regular tg init command (not tg stack init)
	return runTerragruntCommandE(t, options, "init", initArgs(options)...)
}

// initArgs builds the argument list for tg init command.
// This function handles complex configuration that requires special formatting.
func initArgs(options *Options) []string {
	var args []string

	// Add backend configuration and plugin directory arguments
	// These arguments are passed through to the underlying terraform init command
	args = append(args, formatting.FormatBackendConfigAsArgs(options.BackendConfig)...)
	args = append(args, formatting.FormatPluginDirAsArgs(options.PluginDir)...)
	return args
}
