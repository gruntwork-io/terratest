package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
)

// TgStackRun calls tg stack run and returns stdout/stderr
func TgStackRun(t testing.TestingT, options *Options) string {
	out, err := TgStackRunE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TgStackRunE calls tg stack run and returns stdout/stderr
func TgStackRunE(t testing.TestingT, options *Options) (string, error) {
	return runTerragruntStackCommandE(t, options, "run")
}
