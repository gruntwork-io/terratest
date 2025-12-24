package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing/v2"
)

// StackGenerate calls terragrunt stack generate and returns stdout/stderr
func StackGenerate(t testing.TestingT, options *Options) string {
	out, err := StackGenerateE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// StackGenerateE calls terragrunt stack generate and returns stdout/stderr
func StackGenerateE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, options, "generate")
}
