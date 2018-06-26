package terragrunt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/collections"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/shell"
)

// RunTerragruntCommand runs terragrunt with the given arguments and options and return stdout/stderr.
func RunTerragruntCommand(t *testing.T, options *Options, args ...string) string {
	out, err := RunTerragruntCommandE(t, options, args...)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// RunTerragruntCommandE runs terragrunt with the given arguments and options and return stdout/stderr.
func RunTerragruntCommandE(t *testing.T, options *Options, args ...string) (string, error) {
	if options.NoColor && !collections.ListContains(args, "-no-color") {
		args = append(args, "-no-color")
	}

	description := fmt.Sprintf("Running terragrunt %v", args)
	return retry.DoWithRetryE(t, description, options.MaxRetries, options.TimeBetweenRetries, func() (string, error) {
		cmd := shell.Command{
			Command:    "terragrunt",
			Args:       args,
			WorkingDir: options.TerragruntDir,
			Env:        options.EnvVars,
		}

		out, err := shell.RunCommandAndGetOutputE(t, cmd)
		if err == nil {
			return out, nil
		}

		for errorText, errorMessage := range options.RetryableTerragruntErrors {
			if strings.Contains(out, errorText) {
				logger.Logf(t, "terragrunt failed with the error '%s' but this error was expected and warrants a retry. Further details: %s\n", errorText, errorMessage)
				return out, err
			}
		}

		return out, retry.FatalError{Underlying: err}
	})
}
