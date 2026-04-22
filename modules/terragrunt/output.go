package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// TODO: Add OutputAll/OutputAllE when terragrunt supports combined JSON output format.
// Currently, `output --all -json` returns separate JSON objects per module without module prefixes,
// making it impossible to reliably map outputs to their source modules.

// OutputAllJSONContext runs terragrunt run --all output -json and returns the raw JSON string.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. Note: Current terragrunt versions return separate JSON objects per module,
// not a combined object.
func OutputAllJSONContext(t testing.TestingT, ctx context.Context, options *Options) string {
	out, err := OutputAllJSONContextE(t, ctx, options)
	require.NoError(t, err)

	return out
}

// OutputAllJSONContextE runs terragrunt run --all output -json and returns the raw JSON string.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. Note: Current terragrunt versions return separate JSON objects per module,
// not a combined object.
func OutputAllJSONContextE(t testing.TestingT, ctx context.Context, options *Options) (string, error) {
	optsCopy := *options
	optsCopy.TerragruntArgs = append([]string{"--no-color"}, options.TerragruntArgs...)

	args := BuildRunArgs([]string{"--all"}, []string{"output", "-json"})

	rawOutput, err := runTerragruntCommandE(t, ctx, &optsCopy, "run", args...)
	if err != nil {
		return "", err
	}

	// Extract only JSON content from output, filtering log lines and other terragrunt messages
	return ExtractJSONContent(rawOutput)
}

// OutputAllJSON runs terragrunt run --all output -json and returns the raw JSON string.
// Note: Current terragrunt versions return separate JSON objects per module, not a combined object.
//
// Deprecated: Use [OutputAllJSONContext] instead.
func OutputAllJSON(t testing.TestingT, options *Options) string {
	return OutputAllJSONContext(t, context.Background(), options)
}

// OutputAllJSONE runs terragrunt run --all output -json and returns the raw JSON string.
// Note: Current terragrunt versions return separate JSON objects per module, not a combined object.
//
// Deprecated: Use [OutputAllJSONContextE] instead.
func OutputAllJSONE(t testing.TestingT, options *Options) (string, error) {
	return OutputAllJSONContextE(t, context.Background(), options)
}

// OutputJSONContext runs terragrunt run output -json for a single unit and returns clean JSON.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. If key is non-empty, returns the JSON value for that specific output.
// If key is empty, returns all outputs as JSON.
func OutputJSONContext(t testing.TestingT, ctx context.Context, options *Options, key string) string {
	out, err := OutputJSONContextE(t, ctx, options, key)
	require.NoError(t, err)

	return out
}

// OutputJSONContextE runs terragrunt run output -json for a single unit and returns clean JSON.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. If key is non-empty, returns the JSON value for that specific output.
// If key is empty, returns all outputs as JSON.
func OutputJSONContextE(t testing.TestingT, ctx context.Context, options *Options, key string) (string, error) {
	optsCopy := *options
	optsCopy.TerragruntArgs = append([]string{"--no-color"}, options.TerragruntArgs...)

	tfArgs := []string{"-json"}
	if key != "" {
		tfArgs = append(tfArgs, key)
	}

	args := BuildRunArgs([]string{}, append([]string{"output"}, tfArgs...))

	rawOutput, err := runTerragruntCommandE(t, ctx, &optsCopy, "run", args...)
	if err != nil {
		return "", err
	}

	return CleanTerragruntJSON(rawOutput)
}

// OutputJSON runs terragrunt run output -json for a single unit and returns clean JSON.
// If key is non-empty, returns the JSON value for that specific output.
// If key is empty, returns all outputs as JSON.
//
// Deprecated: Use [OutputJSONContext] instead.
func OutputJSON(t testing.TestingT, options *Options, key string) string {
	return OutputJSONContext(t, context.Background(), options, key)
}

// OutputJSONE runs terragrunt run output -json for a single unit and returns clean JSON.
// If key is non-empty, returns the JSON value for that specific output.
// If key is empty, returns all outputs as JSON.
//
// Deprecated: Use [OutputJSONContextE] instead.
func OutputJSONE(t testing.TestingT, options *Options, key string) (string, error) {
	return OutputJSONContextE(t, context.Background(), options, key)
}

// OutputAllJson runs terragrunt run --all output -json and returns the raw JSON string.
// Note: Current terragrunt versions return separate JSON objects per module, not a combined object.
//
// Deprecated: Use [OutputAllJSONContext] instead.
func OutputAllJson(t testing.TestingT, options *Options) string { //nolint:staticcheck // Deprecated wrapper kept for backward compatibility.
	return OutputAllJSONContext(t, context.Background(), options)
}

// OutputAllJsonE runs terragrunt run --all output -json and returns the raw JSON string.
// Note: Current terragrunt versions return separate JSON objects per module, not a combined object.
//
// Deprecated: Use [OutputAllJSONContextE] instead.
func OutputAllJsonE(t testing.TestingT, options *Options) (string, error) { //nolint:staticcheck // Deprecated wrapper kept for backward compatibility.
	return OutputAllJSONContextE(t, context.Background(), options)
}

// OutputJson runs terragrunt run output -json for a single unit and returns clean JSON.
// If key is non-empty, returns the JSON value for that specific output.
// If key is empty, returns all outputs as JSON.
//
// Deprecated: Use [OutputJSONContext] instead.
func OutputJson(t testing.TestingT, options *Options, key string) string { //nolint:staticcheck // Deprecated wrapper kept for backward compatibility.
	return OutputJSONContext(t, context.Background(), options, key)
}

// OutputJsonE runs terragrunt run output -json for a single unit and returns clean JSON.
// If key is non-empty, returns the JSON value for that specific output.
// If key is empty, returns all outputs as JSON.
//
// Deprecated: Use [OutputJSONContextE] instead.
func OutputJsonE(t testing.TestingT, options *Options, key string) (string, error) { //nolint:staticcheck // Deprecated wrapper kept for backward compatibility.
	return OutputJSONContextE(t, context.Background(), options, key)
}
