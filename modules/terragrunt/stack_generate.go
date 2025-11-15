package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
)

// TgStackGenerate calls tg stack generate and returns stdout/stderr
func TgStackGenerate(t testing.TestingT, options *Options) string {
	out, err := TgStackGenerateE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TgStackGenerateE calls tg stack generate and returns stdout/stderr
func TgStackGenerateE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, options, "generate")
}
