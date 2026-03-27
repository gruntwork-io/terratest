package terragrunt

import (
	"context"
	"encoding/json"

	"github.com/gruntwork-io/terratest/modules/testing"
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
	// Prepare options with no-color flag for parsing
	optsCopy := *options
	optsCopy.TerragruntArgs = append([]string{"--no-color"}, options.TerragruntArgs...)

	var args []string
	if key != "" {
		args = append(args, key)
	}
	// Append any user-provided TerraformArgs
	if len(options.TerraformArgs) > 0 {
		args = append(args, options.TerraformArgs...)
	}

	// Output command for stack
	rawOutput, err := runTerragruntStackCommandE(
		t, ctx, &optsCopy, "output", args...)
	if err != nil {
		return "", err
	}

	// Extract the actual value from output
	return CleanTerragruntOutput(rawOutput), nil
}

// StackOutput calls terragrunt stack output for the given variable and returns its value as a string.
//
// Deprecated: Use [StackOutputContext] instead.
func StackOutput(t testing.TestingT, options *Options, key string) string {
	return StackOutputContext(t, context.Background(), options, key)
}

// StackOutputE calls terragrunt stack output for the given variable and returns its value as a string.
//
// Deprecated: Use [StackOutputContextE] instead.
func StackOutputE(t testing.TestingT, options *Options, key string) (string, error) {
	return StackOutputContextE(t, context.Background(), options, key)
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
	// Prepare options with no-color flag
	optsCopy := *options
	optsCopy.TerragruntArgs = append([]string{"--no-color"}, options.TerragruntArgs...)

	// -json is an OpenTofu/Terraform flag that should go after the output command
	args := []string{"-json"}
	if key != "" {
		args = append(args, key)
	}
	// Append any user-provided TerraformArgs
	if len(options.TerraformArgs) > 0 {
		args = append(args, options.TerraformArgs...)
	}

	// Output command for stack
	rawOutput, err := runTerragruntStackCommandE(
		t, ctx, &optsCopy, "output", args...)
	if err != nil {
		return "", err
	}

	// Parse and format JSON output
	return CleanTerragruntJSON(rawOutput)
}

// StackOutputJSON calls terragrunt stack output for the given variable and returns the result as a JSON string.
// If key is an empty string, it will return all the output variables.
//
// Deprecated: Use [StackOutputJSONContext] instead.
func StackOutputJSON(t testing.TestingT, options *Options, key string) string {
	return StackOutputJSONContext(t, context.Background(), options, key)
}

// StackOutputJSONE calls terragrunt stack output for the given variable and returns the result as a JSON string.
// If key is an empty string, it will return all the output variables.
//
// Deprecated: Use [StackOutputJSONContextE] instead.
func StackOutputJSONE(t testing.TestingT, options *Options, key string) (string, error) {
	return StackOutputJSONContextE(t, context.Background(), options, key)
}

// StackOutputJson calls terragrunt stack output for the given variable and returns the result as a JSON string.
// If key is an empty string, it will return all the output variables.
//
// Deprecated: Use [StackOutputJSONContext] instead.
func StackOutputJson(t testing.TestingT, options *Options, key string) string { //nolint:staticcheck // Deprecated wrapper kept for backward compatibility.
	return StackOutputJSONContext(t, context.Background(), options, key)
}

// StackOutputJsonE calls terragrunt stack output for the given variable and returns the result as a JSON string.
// If key is an empty string, it will return all the output variables.
//
// Deprecated: Use [StackOutputJSONContextE] instead.
func StackOutputJsonE(t testing.TestingT, options *Options, key string) (string, error) { //nolint:staticcheck // Deprecated wrapper kept for backward compatibility.
	return StackOutputJSONContextE(t, context.Background(), options, key)
}

// StackOutputAllContext gets all stack outputs and returns them as a map[string]interface{}.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func StackOutputAllContext(t testing.TestingT, ctx context.Context, options *Options) map[string]interface{} {
	outputs, err := StackOutputAllContextE(t, ctx, options)
	require.NoError(t, err)

	return outputs
}

// StackOutputAllContextE gets all stack outputs and returns them as a map[string]interface{}.
// The provided context is passed through to the underlying command execution, allowing for timeout
// and cancellation control.
func StackOutputAllContextE(t testing.TestingT, ctx context.Context, options *Options) (map[string]interface{}, error) {
	jsonOutput, err := StackOutputJSONContextE(t, ctx, options, "")
	if err != nil {
		return nil, err
	}

	var outputs map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOutput), &outputs); err != nil {
		return nil, err
	}

	return outputs, nil
}

// StackOutputAll gets all stack outputs and returns them as a map[string]interface{}.
//
// Deprecated: Use [StackOutputAllContext] instead.
func StackOutputAll(t testing.TestingT, options *Options) map[string]interface{} {
	return StackOutputAllContext(t, context.Background(), options)
}

// StackOutputAllE gets all stack outputs and returns them as a map[string]interface{}.
//
// Deprecated: Use [StackOutputAllContextE] instead.
func StackOutputAllE(t testing.TestingT, options *Options) (map[string]interface{}, error) {
	return StackOutputAllContextE(t, context.Background(), options)
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

// StackOutputListAll gets all stack output variable names and returns them as a slice.
//
// Deprecated: Use [StackOutputListAllContext] instead.
func StackOutputListAll(t testing.TestingT, options *Options) []string {
	return StackOutputListAllContext(t, context.Background(), options)
}

// StackOutputListAllE gets all stack output variable names and returns them as a slice.
//
// Deprecated: Use [StackOutputListAllContextE] instead.
func StackOutputListAllE(t testing.TestingT, options *Options) ([]string, error) {
	return StackOutputListAllContextE(t, context.Background(), options)
}
