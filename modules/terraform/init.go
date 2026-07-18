package terraform

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// InitContext calls terraform init with the given options and returns stdout/stderr.
// The context argument can be used for cancellation or timeout control.
func InitContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := InitContextE(t, ctx, options)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// InitContextE calls terraform init with the given options and returns stdout/stderr.
// The context argument can be used for cancellation or timeout control.
func InitContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	args := []string{"init", fmt.Sprintf("-upgrade=%t", options.Upgrade)}

	if options.Reconfigure {
		args = append(args, "-reconfigure")
	}

	if options.MigrateState {
		args = append(args, "-migrate-state", "-force-copy")
	}

	if options.NoColor {
		args = append(args, "-no-color")
	}

	args = append(args, FormatTerraformBackendConfigAsArgs(options.BackendConfig)...)
	args = append(args, FormatTerraformPluginDirAsArgs(options.PluginDir)...)

	return RunTerraformCommandContextE(t, ctx, options, prepend(options.ExtraArgs.Init, args...)...)
}
