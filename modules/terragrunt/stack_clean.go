package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
)

// TgStackClean calls terragrunt stack clean and returns stdout/stderr
func TgStackClean(t testing.TestingT, options *Options) string {
	out, err := TgStackCleanE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TgStackCleanE calls terragrunt stack clean and returns stdout/stderr
func TgStackCleanE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, options, "clean", cleanStackArgs(options)...)
}

// cleanStackArgs builds the argument list for terragrunt stack clean command.
// All terragrunt command-line flags are now passed via ExtraArgs.
func cleanStackArgs(options *Options) []string {
	// Return all user-specified terragrunt command-line arguments
	// The user passes the specific args they need for their stack clean operation
	return options.ExtraArgs
}