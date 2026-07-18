package terragrunt

import (
	"context"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

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

	return ExtractJSONContent(rawOutput)
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
