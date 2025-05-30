package terraform

import (
	"fmt"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// Init calls terraform init and return stdout/stderr.
func Init(t testing.TestingT, options *Options) string {
	out, err := InitE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// InitE calls terraform init and return stdout/stderr.
func InitE(t testing.TestingT, options *Options) (string, error) {
	return RunTerraformCommandE(t, options, initArgs(options)...)
}

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
	if options.TerraformBinary != "terragrunt" {
		return "", TgInvalidBinary(options.TerraformBinary)
	}
	return runTerragruntStackCommandE(t, options, initArgs(options)...)
}

func initArgs(options *Options) []string {
	args := []string{"init", fmt.Sprintf("-upgrade=%t", options.Upgrade)}

	// Append reconfigure option if specified
	if options.Reconfigure {
		args = append(args, "-reconfigure")
	}
	// Append combination of migrate-state and force-copy to suppress answer prompt
	if options.MigrateState {
		args = append(args, "-migrate-state", "-force-copy")
	}
	// Append no-color option if needed
	if options.NoColor {
		args = append(args, "-no-color")
	}

	args = append(args, FormatTerraformBackendConfigAsArgs(options.BackendConfig)...)
	args = append(args, FormatTerraformPluginDirAsArgs(options.PluginDir)...)
	args = append(args, options.ExtraArgs.Init...)
	return args
}
