package terragrunt

import (
	"context"
	"strings"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// RenderContext runs terragrunt render to output the resolved terragrunt configuration as HCL.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This is useful for verifying merged includes, resolved dependencies,
// and executed functions without actually applying any changes.
func RenderContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := RenderContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// RenderContextE runs terragrunt render to output the resolved terragrunt configuration as HCL.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This is useful for verifying merged includes, resolved dependencies,
// and executed functions without actually applying any changes. Log lines are stripped from the output.
func RenderContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	rawOutput, err := runTerragruntCommandE(t, ctx, options, "render")
	if err != nil {
		return "", err
	}

	return FilterLogLines(rawOutput), nil
}

// FilterLogLines removes terragrunt log lines while preserving original indentation.
// Unlike [RemoveLogLines] (which trims whitespace for JSON extraction), this keeps
// leading whitespace intact so HCL output structure is preserved.
func FilterLogLines(rawOutput string) string {
	lines := strings.Split(rawOutput, "\n")

	var result []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || IsLogLine(trimmed) || IsMetadataLine(trimmed) {
			continue
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// RenderJSONContext runs terragrunt render --format json and returns the cleaned JSON output.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This is useful for programmatic assertions on the resolved terragrunt
// configuration.
func RenderJSONContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := RenderJSONContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// RenderJSONContextE runs terragrunt render --format json and returns the cleaned JSON output.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. This is useful for programmatic assertions on the resolved terragrunt
// configuration.
func RenderJSONContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	optsCopy := *options
	optsCopy.TerragruntArgs = append([]string{"--no-color"}, options.TerragruntArgs...)

	rawOutput, err := runTerragruntCommandE(t, ctx, &optsCopy, "render", "--format", "json")
	if err != nil {
		return "", err
	}

	return CleanTerragruntJSON(rawOutput)
}
