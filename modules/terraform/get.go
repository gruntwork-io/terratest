package terraform

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// GetContext calls terraform get and returns stdout/stderr. The provided context is passed through to the underlying
// command execution, allowing for timeout and cancellation control.
func GetContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := GetContextE(t, ctx, options)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// GetContextE calls terraform get and returns stdout/stderr. The provided context is passed through to the underlying
// command execution, allowing for timeout and cancellation control.
func GetContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return RunTerraformCommandContextE(t, ctx, options, prepend(options.ExtraArgs.Get, "get", "-update")...)
}
