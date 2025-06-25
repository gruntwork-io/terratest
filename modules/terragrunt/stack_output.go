package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
)

// TgStackOutput calls terragrunt output and return stdout/stderr
func TgStackOutput(t testing.TestingT, options *Options) string {
	out, err := TgStackOutputE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TgStackOutputE calls terragrunt output and return stdout/stderr
func TgStackOutputE(t testing.TestingT, options *Options) (string, error) {
	return terragruntStackCommandE(t, options, outputArgs(options)...)
}

func outputArgs(options *Options) []string {
	args := []string{"output"}

	// Append no-color option if needed
	if options.NoColor {
		args = append(args, "-no-color")
	}

	// Use Apply extra args for output command as it's a similar operation
	if len(options.ExtraArgs.Apply) > 0 {
		args = append(args, options.ExtraArgs.Apply...)
	}
	return args
}
