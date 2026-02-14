package terragrunt

import (
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// TODO: Add OutputAll/OutputAllE when terragrunt supports combined JSON output format.
// Currently, `output --all -json` returns separate JSON objects per module without module prefixes,
// making it impossible to reliably map outputs to their source modules.

// OutputAllJson runs terragrunt output --all -json and returns the raw JSON string.
// Note: Current terragrunt versions return separate JSON objects per module, not a combined object.
func OutputAllJson(t testing.TestingT, options *Options) string {
	out, err := OutputAllJsonE(t, options)
	require.NoError(t, err)
	return out
}

// OutputAllJsonE runs terragrunt output --all -json and returns the raw JSON string.
// Note: Current terragrunt versions return separate JSON objects per module, not a combined object.
func OutputAllJsonE(t testing.TestingT, options *Options) (string, error) {
	optsCopy := *options
	optsCopy.TerragruntArgs = append([]string{"--no-color"}, options.TerragruntArgs...)

	rawOutput, err := runTerragruntCommandE(t, &optsCopy, "output", "--all", "-json")
	if err != nil {
		return "", err
	}

	// Extract only JSON content from output, filtering log lines and other terragrunt messages
	return extractJsonContent(rawOutput), nil
}
