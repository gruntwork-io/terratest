package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// TgStackOutput calls terragrunt stack output for the given variable and returns its value
func TgStackOutput(t testing.TestingT, options *Options, key string, extraArgs ...string) string {
	out, err := TgStackOutputE(t, options, key, extraArgs...)
	require.NoError(t, err)
	return out
}

// TgStackOutputE calls terragrunt stack output for the given variable and returns its value
func TgStackOutputE(t testing.TestingT, options *Options, key string, extraArgs ...string) (string, error) {
	return terragruntStackCommandE(t, options, outputStackArgs(options, key, extraArgs...)...)
}

func outputStackArgs(options *Options, key string, extraArgs ...string) []string {
	args := []string{"output"}
	args = append(args, extraArgs...)
	if key != "" {
		args = append(args, key)
	}
	if options.NoColor {
		args = append(args, "-no-color")
	}
	return args
}
