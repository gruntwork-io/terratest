package terragrunt

import (
	"context"
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
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

// Render runs terragrunt render to output the resolved terragrunt configuration as HCL.
// This is useful for verifying merged includes, resolved dependencies, and executed functions
// without actually applying any changes.
//
// Deprecated: Use [RenderContext] instead.
func Render(t testing.TestingT, options *Options) string {
	return RenderContext(t, context.Background(), options)
}

// RenderE runs terragrunt render to output the resolved terragrunt configuration as HCL.
// This is useful for verifying merged includes, resolved dependencies, and executed functions
// without actually applying any changes. Log lines are stripped from the output.
//
// Deprecated: Use [RenderContextE] instead.
func RenderE(t testing.TestingT, options *Options) (string, error) {
	return RenderContextE(t, context.Background(), options)
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

// RenderJSON runs terragrunt render --format json and returns the cleaned JSON output.
// This is useful for programmatic assertions on the resolved terragrunt configuration.
//
// Deprecated: Use [RenderJSONContext] instead.
func RenderJSON(t testing.TestingT, options *Options) string {
	return RenderJSONContext(t, context.Background(), options)
}

// RenderJSONE runs terragrunt render --format json and returns the cleaned JSON output.
// This is useful for programmatic assertions on the resolved terragrunt configuration.
//
// Deprecated: Use [RenderJSONContextE] instead.
func RenderJSONE(t testing.TestingT, options *Options) (string, error) {
	return RenderJSONContextE(t, context.Background(), options)
}

// RenderJson runs terragrunt render --format json and returns the cleaned JSON output.
// This is useful for programmatic assertions on the resolved terragrunt configuration.
//
// Deprecated: Use [RenderJSONContext] instead.
func RenderJson(t testing.TestingT, options *Options) string { //nolint:staticcheck // Deprecated wrapper kept for backward compatibility.
	return RenderJSONContext(t, context.Background(), options)
}

// RenderJsonE runs terragrunt render --format json and returns the cleaned JSON output.
// This is useful for programmatic assertions on the resolved terragrunt configuration.
//
// Deprecated: Use [RenderJSONContextE] instead.
func RenderJsonE(t testing.TestingT, options *Options) (string, error) { //nolint:staticcheck // Deprecated wrapper kept for backward compatibility.
	return RenderJSONContextE(t, context.Background(), options)
}
