package terragrunt

import (
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Render runs terragrunt render to output the resolved terragrunt configuration as HCL.
// This is useful for verifying merged includes, resolved dependencies, and executed functions
// without actually applying any changes.
func Render(t testing.TestingT, options *Options) string {
	out, err := RenderE(t, options)
	require.NoError(t, err)
	return out
}

// RenderE runs terragrunt render to output the resolved terragrunt configuration as HCL.
// This is useful for verifying merged includes, resolved dependencies, and executed functions
// without actually applying any changes. Log lines are stripped from the output.
func RenderE(t testing.TestingT, options *Options) (string, error) {
	rawOutput, err := runTerragruntCommandE(t, options, "render")
	if err != nil {
		return "", err
	}

	return filterLogLines(rawOutput), nil
}

// filterLogLines removes terragrunt log lines while preserving original indentation.
// Unlike removeLogLines (which trims whitespace for JSON extraction), this keeps
// leading whitespace intact so HCL output structure is preserved.
func filterLogLines(rawOutput string) string {
	lines := strings.Split(rawOutput, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || isLogLine(trimmed) || isMetadataLine(trimmed) {
			continue
		}
		result = append(result, line)
	}
	return strings.Join(result, "\n")
}

// RenderJson runs terragrunt render --format json and returns the cleaned JSON output.
// This is useful for programmatic assertions on the resolved terragrunt configuration.
func RenderJson(t testing.TestingT, options *Options) string {
	out, err := RenderJsonE(t, options)
	require.NoError(t, err)
	return out
}

// RenderJsonE runs terragrunt render --format json and returns the cleaned JSON output.
// This is useful for programmatic assertions on the resolved terragrunt configuration.
func RenderJsonE(t testing.TestingT, options *Options) (string, error) {
	optsCopy := *options
	optsCopy.TerragruntArgs = append([]string{"--no-color"}, options.TerragruntArgs...)

	rawOutput, err := runTerragruntCommandE(t, &optsCopy, "render", "--format", "json")
	if err != nil {
		return "", err
	}

	return cleanTerragruntJson(rawOutput)
}
