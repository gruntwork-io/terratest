package terragrunt

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing/v2"
)

// Output calls terragrunt stack output for the given variable and returns its value as a string
func Output(t testing.TestingT, options *Options, key string) string {
	out, err := OutputE(t, options, key)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// OutputE calls terragrunt stack output for the given variable and returns its value as a string
func OutputE(t testing.TestingT, options *Options, key string) (string, error) {
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

// OutputJson calls terragrunt stack output for the given variable and returns the result as the json string.
// If key is an empty string, it will return all the output variables.
func OutputJson(t testing.TestingT, options *Options, key string) string {
	str, err := OutputJsonE(t, options, key)
	if err != nil {
		t.Fatal(err)
	}
	return str
}

// OutputJsonE calls terragrunt stack output for the given variable and returns the
// result as the json string.
// If key is an empty string, it will return all the output variables.
func OutputJsonE(t testing.TestingT, options *Options, key string) (string, error) {
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

// OutputAll gets all stack outputs and returns them as a map[string]interface{}
func OutputAll(t testing.TestingT, options *Options) map[string]interface{} {
	outputs, err := OutputAllE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return outputs
}

// OutputAllE gets all stack outputs and returns them as a map[string]interface{}
func OutputAllE(t testing.TestingT, options *Options) (map[string]interface{}, error) {
	jsonOutput, err := OutputJsonE(t, options, "")
	if err != nil {
		return nil, err
	}

	var outputs map[string]interface{}
	if err := json.Unmarshal([]byte(jsonOutput), &outputs); err != nil {
		return nil, err
	}

	return outputs, nil
}

// OutputListAll gets all stack output variable names and returns them as a slice
func OutputListAll(t testing.TestingT, options *Options) []string {
	keys, err := OutputListAllE(t, options)
	if err != nil {
		t.Fatal(err)
	}
	return keys
}

// OutputListAllE gets all stack output variable names and returns them as a slice
func OutputListAllE(t testing.TestingT, options *Options) ([]string, error) {
	outputs, err := OutputAllE(t, options)
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(outputs))
	for key := range outputs {
		keys = append(keys, key)
	}

	return keys, nil
}

var (
	// tgLogLevel matches log lines containing fields for time, level, prefix, binary, and message
	tgLogLevel = regexp.MustCompile(`.*time=\S+ level=\S+ prefix=\S+ binary=\S+ msg=.*`)
)

// removeLogLines removes terragrunt log lines from output
func removeLogLines(rawOutput string) string {
	// Remove lines matching the log pattern
	cleaned := tgLogLevel.ReplaceAllString(rawOutput, "")

	// Split into lines and filter
	lines := strings.Split(cleaned, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip empty lines and lines that are clearly log lines (containing msg= with log context)
		if trimmed != "" && !strings.Contains(line, " msg=") {
			result = append(result, trimmed)
		}
	}

	return strings.Join(result, "\n")
}

// cleanTerragruntOutput extracts the actual output value from terragrunt stack's verbose output
//
// Example input (raw tg output):
//
//	time=2023-07-11T10:30:45Z level=info prefix=terragrunt binary=terragrunt msg="Initializing..."
//	time=2023-07-11T10:30:46Z level=info prefix=terragrunt binary=terragrunt msg="Running command..."
//	"my-bucket-name"
//
// Example output (cleaned):
//
//	my-bucket-name
//
// For JSON values, it preserves the structure:
// Input:
//
//	time=2023-07-11T10:30:45Z level=info prefix=terragrunt binary=terragrunt msg="Running..."
//	{"vpc_id": "vpc-12345", "subnet_ids": ["subnet-1", "subnet-2"]}
//
// Output:
//
//	{"vpc_id": "vpc-12345", "subnet_ids": ["subnet-1", "subnet-2"]}
func cleanTerragruntOutput(rawOutput string) (string, error) {
	// Remove terragrunt log lines
	finalOutput := removeLogLines(rawOutput)
	if finalOutput == "" {
		return "", nil
	}

	// Check if it's JSON (starts with { or [)
	if strings.HasPrefix(finalOutput, "{") || strings.HasPrefix(finalOutput, "[") {
		// For JSON output, return as-is
		return finalOutput, nil
	}

	// For simple values, remove surrounding quotes if present
	if strings.HasPrefix(finalOutput, "\"") && strings.HasSuffix(finalOutput, "\"") {
		finalOutput = strings.Trim(finalOutput, "\"")
	}

	return finalOutput, nil
}

// cleanTerragruntJson cleans the JSON output from terragrunt stack command
//
// Example input (raw tg JSON output):
//
//	time=2023-07-11T10:30:45Z level=info prefix=terragrunt binary=terragrunt msg="Initializing..."
//	time=2023-07-11T10:30:46Z level=info prefix=terragrunt binary=terragrunt msg="Running command..."
//	{"mother.output":{"sensitive":false,"type":"string","value":"mother/test.txt"},"father.output":{"sensitive":false,"type":"string","value":"father/test.txt"}}
//
// Example output (cleaned and formatted):
//
//	{
//	  "mother.output": {
//	    "sensitive": false,
//	    "type": "string",
//	    "value": "mother/test.txt"
//	  },
//	  "father.output": {
//	    "sensitive": false,
//	    "type": "string",
//	    "value": "father/test.txt"
//	  }
//	}
func cleanTerragruntJson(input string) (string, error) {
	// Remove terragrunt log lines
	cleaned := removeLogLines(input)

	// Parse JSON
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(cleaned), &jsonObj); err != nil {
		return "", err
	}

	// Format JSON output with indentation
	normalized, err := json.MarshalIndent(jsonObj, "", "  ")
	if err != nil {
		return "", err
	}

	return string(normalized), nil
}
