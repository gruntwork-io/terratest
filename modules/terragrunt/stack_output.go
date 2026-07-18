package terragrunt

import (
	"context"
	"encoding/json"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// StackOutputContext calls terragrunt stack output for the given variable and returns its value as a string.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func StackOutputContext(t testing.TestingT, ctx context.Context, options *Options, key string) string {
	out, err := StackOutputContextE(t, ctx, options, key)
	require.NoError(t, err)

	return out
}

// StackOutputContextE calls terragrunt stack output for the given variable and returns its value as a string.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func StackOutputContextE(t testing.TestingT, ctx context.Context, options *Options, key string) (string, error) {

	optsCopy := *options
	optsCopy.TerragruntArgs = append([]string{"--no-color"}, options.TerragruntArgs...)

	var args []string
	if key != "" {
		args = append(args, key)
	}

	if len(options.TerraformArgs) > 0 {
		args = append(args, options.TerraformArgs...)
	}

	rawOutput, err := runTerragruntStackCommandE(
		t, ctx, &optsCopy, "output", args...)
	if err != nil {
		return "", err
	}

	return CleanTerragruntOutput(rawOutput), nil
}

// StackOutputJSONContext calls terragrunt stack output for the given variable and returns the result as a JSON string.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. If key is an empty string, it will return all the output variables.
func StackOutputJSONContext(t testing.TestingT, ctx context.Context, options *Options, key string) string {
	str, err := StackOutputJSONContextE(t, ctx, options, key)
	require.NoError(t, err)

	return str
}

// StackOutputJSONContextE calls terragrunt stack output for the given variable and returns the result as a JSON string.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control. If key is an empty string, it will return all the output variables.
func StackOutputJSONContextE(t testing.TestingT, ctx context.Context, options *Options, key string) (string, error) {

	optsCopy := *options
	optsCopy.TerragruntArgs = append([]string{"--no-color"}, options.TerragruntArgs...)

	args := []string{"-json"}
	if key != "" {
		args = append(args, key)
	}

	if len(options.TerraformArgs) > 0 {
		args = append(args, options.TerraformArgs...)
	}

	rawOutput, err := runTerragruntStackCommandE(
		t, ctx, &optsCopy, "output", args...)
	if err != nil {
		return "", err
	}

	return CleanTerragruntJSON(rawOutput)
}

// StackOutputAllContext gets all stack outputs and returns them as a map[string]any.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func StackOutputAllContext(t testing.TestingT, ctx context.Context, options *Options) map[string]any {
	outputs, err := StackOutputAllContextE(t, ctx, options)
	require.NoError(t, err)

	return outputs
}

// StackOutputAllContextE gets all stack outputs and returns them as a map[string]any.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func StackOutputAllContextE(t testing.TestingT, ctx context.Context, options *Options) (map[string]any, error) {
	jsonOutput, err := StackOutputJSONContextE(t, ctx, options, "")
	if err != nil {
		return nil, err
	}

	var outputs map[string]any
	if err := json.Unmarshal([]byte(jsonOutput), &outputs); err != nil {
		return nil, err
	}

	return outputs, nil
}

// StackOutputListAllContext gets all stack output variable names and returns them as a slice.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func StackOutputListAllContext(t testing.TestingT, ctx context.Context, options *Options) []string {
	keys, err := StackOutputListAllContextE(t, ctx, options)
	require.NoError(t, err)

	return keys
}

// StackOutputListAllContextE gets all stack output variable names and returns them as a slice.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func StackOutputListAllContextE(t testing.TestingT, ctx context.Context, options *Options) ([]string, error) {
	outputs, err := StackOutputAllContextE(t, ctx, options)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(outputs))
	for key := range outputs {
		keys = append(keys, key)
	}

	return keys, nil
}
