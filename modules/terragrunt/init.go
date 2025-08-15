package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// TgStackInit calls terragrunt init and return stdout/stderr
func TgStackInit(t testing.TestingT, options *Options) string {
	out, err := TgStackInitE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TgStackInitE calls terragrunt init and return stdout/stderr
func TgStackInitE(t testing.TestingT, options *Options) (string, error) {
	// Use regular terragrunt init command (not terragrunt stack init)
	return runTerragruntCommandE(t, options, "init", initStackArgs(options)...)
}

// initStackArgs builds the argument list for terragrunt init command.
// This function handles complex configuration that requires special formatting.
func initStackArgs(options *Options) []string {
	var args []string

	// Add complex configuration that requires special formatting
	// These are terraform-specific arguments that need special formatting
	args = append(args, terraform.FormatTerraformBackendConfigAsArgs(options.BackendConfig)...)
	args = append(args, terraform.FormatTerraformPluginDirAsArgs(options.PluginDir)...)

	// User-specified arguments are now handled by GetArgsForCommand
	// which properly separates TerragruntArgs and TerraformArgs

	return args
}
