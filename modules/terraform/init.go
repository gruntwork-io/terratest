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
	args := []string{"init", fmt.Sprintf("-upgrade=%t", options.Upgrade)}

	// Append reconfigure option if specified
	if options.Reconfigure {
		args = append(args, "-reconfigure")
	}
	// Append combination of migrate-state and force-copy to suppress answer prompt
	if options.MigrateState {
		args = append(args, "-migrate-state", "-force-copy")
	}

	args = append(args, FormatTerraformBackendConfigAsArgs(options.BackendConfig)...)
	args = append(args, FormatTerraformPluginDirAsArgs(options.PluginDir)...)
	return RunTerraformCommandE(t, options, args...)
}
