package terraform

import (
	"context"
	"fmt"

	"github.com/gruntwork-io/terratest/modules/testing"
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

	return RunTerraformCommandContextE(t, ctx, options, prepend(options.ExtraArgs.Init, args...)...)
}

// Init calls terraform init and return stdout/stderr.
//
// Deprecated: Use [InitContext] instead.
func Init(t testing.TestingT, options *Options) string {
	return InitContext(t, context.Background(), options)
}

// InitE calls terraform init and return stdout/stderr.
//
// Deprecated: Use [InitContextE] instead.
func InitE(t testing.TestingT, options *Options) (string, error) {
	return InitContextE(t, context.Background(), options)
}
