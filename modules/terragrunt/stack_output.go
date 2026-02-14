package terragrunt

import (
	"encoding/json"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// StackOutput calls terragrunt stack output for the given variable and returns its value as a string
func StackOutput(t testing.TestingT, options *Options, key string) string {
	out, err := StackOutputE(t, options, key)
	require.NoError(t, err)
	return out
}

// StackOutputE calls terragrunt stack output for the given variable and returns its value as a string
func StackOutputE(t testing.TestingT, options *Options, key string) (string, error) {
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
		t, &optsCopy, "output", args...)
	if err != nil {
		return "", err
	}

	// Extract the actual value from output
	cleaned, err := cleanTerragruntOutput(rawOutput)
	if err != nil {
		return "", err
	}
	return cleaned, nil
}

// StackOutputJson calls terragrunt stack output for the given variable and returns the result as the json string.
// If key is an empty string, it will return all the output variables.
func StackOutputJson(t testing.TestingT, options *Options, key string) string {
	str, err := StackOutputJsonE(t, options, key)
	require.NoError(t, err)
	return str
}

// StackOutputJsonE calls terragrunt stack output for the given variable and returns the
// result as the json string.
// If key is an empty string, it will return all the output variables.
func StackOutputJsonE(t testing.TestingT, options *Options, key string) (string, error) {
	// Prepare options with no-color flag
	optsCopy := *options
	optsCopy.TerragruntArgs = append([]string{"--no-color"}, options.TerragruntArgs...)

	// -json is a terraform flag that should go after the output command
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
		t, &optsCopy, "output", args...)
	if err != nil {
		return "", err
	}

	// Parse and format JSON output
	return cleanTerragruntJson(rawOutput)
}

// StackOutputAll gets all stack outputs and returns them as a map[string]interface{}
func StackOutputAll(t testing.TestingT, options *Options) map[string]interface{} {
	outputs, err := StackOutputAllE(t, options)
	require.NoError(t, err)
	return outputs
}

// StackOutputAllE gets all stack outputs and returns them as a map[string]interface{}
func StackOutputAllE(t testing.TestingT, options *Options) (map[string]interface{}, error) {
	jsonOutput, err := StackOutputJsonE(t, options, "")
	if err != nil {
		return nil, err
	}

	var outputs map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOutput), &outputs); err != nil {
		return nil, err
	}

	return outputs, nil
}

// StackOutputListAll gets all stack output variable names and returns them as a slice
func StackOutputListAll(t testing.TestingT, options *Options) []string {
	keys, err := StackOutputListAllE(t, options)
	require.NoError(t, err)
	return keys
}

// StackOutputListAllE gets all stack output variable names and returns them as a slice
func StackOutputListAllE(t testing.TestingT, options *Options) ([]string, error) {
	outputs, err := StackOutputAllE(t, options)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(outputs))
	for key := range outputs {
		keys = append(keys, key)
	}

	return keys, nil
}
