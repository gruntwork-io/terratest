package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing/v2"
)

// StackRun calls terragrunt stack run and returns stdout/stderr
func StackRun(t testing.TestingT, options *Options) string {
	out, err := StackRunE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// StackRunE calls terragrunt stack run and returns stdout/stderr
func StackRunE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, options, "run")
}
