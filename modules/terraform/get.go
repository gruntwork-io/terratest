package terraform

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// Get calls terraform get and return stdout/stderr.
//
// Deprecated: Use [GetContext] instead.
func Get(t testing.TestingT, options *Options) string {
	return GetContext(t, context.Background(), options)
}

// GetContext calls terraform get and returns stdout/stderr. The provided context is passed through to the underlying
// command execution, allowing for timeout and cancellation control.
func GetContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := GetContextE(t, ctx, options)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// GetE calls terraform get and return stdout/stderr.
//
// Deprecated: Use [GetContextE] instead.
func GetE(t testing.TestingT, options *Options) (string, error) {
	return GetContextE(t, context.Background(), options)
}

// GetContextE calls terraform get and returns stdout/stderr. The provided context is passed through to the underlying
// command execution, allowing for timeout and cancellation control.
func GetContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	return RunTerraformCommandContextE(t, ctx, options, prepend(options.ExtraArgs.Get, "get", "-update")...)
}
