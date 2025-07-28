package terragrunt

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// TgOutput calls terragrunt stack output for the given variable and returns its value as a string
func TgOutput(t testing.TestingT, options *Options, key string) string {
	out, err := TgOutputE(t, options, key)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// TgOutputE calls terragrunt stack output for the given variable and returns its value as a string
func TgOutputE(t testing.TestingT, options *Options, key string) (string, error) {
	// Build args for stack output command
	args := []string{"-no-color"}
	args = append(args, options.ExtraArgs...)
	if key != "" {
		args = append(args, key)
	}
	
	// Use runTerragruntStackCommandE but pass empty additionalArgs since we've already
	// built the complete args list (output doesn't use -- separator)
	rawOutput, err := runTerragruntStackCommandE(t, options, "output", args...)
	if err != nil {
		return "", err
	}
	
	// Clean the output to extract the actual value
	cleaned, err := cleanTerragruntOutput(rawOutput)
	if err != nil {
		return "", err
	}
	return cleaned, nil
}

// TgOutputJson calls terragrunt stack output for the given variable and returns the
// result as the json string.
// If key is an empty string, it will return all the output variables.
func TgOutputJson(t testing.TestingT, options *Options, key string) string {
	str, err := TgOutputJsonE(t, options, key)
	if err != nil {
		t.Fatal(err)
	}
	return str
}

// TgOutputJsonE calls terragrunt stack output for the given variable and returns the
// result as the json string.
// If key is an empty string, it will return all the output variables.
func TgOutputJsonE(t testing.TestingT, options *Options, key string) (string, error) {
	// Build args for stack output command with JSON flag
	args := []string{"-no-color", "-json"}
	args = append(args, options.ExtraArgs...)
	if key != "" {
		args = append(args, key)
	}
	
	// Use runTerragruntStackCommandE but pass empty additionalArgs since we've already
	// built the complete args list (output doesn't use -- separator)
	rawOutput, err := runTerragruntStackCommandE(t, options, "output", args...)
	if err != nil {
		return "", err
	}
	
	// Clean and format the JSON output
	return cleanTerragruntJson(rawOutput)
}


var (
	// tgLogLevel matches log lines containing fields for time, level, prefix, binary, and message
	tgLogLevel = regexp.MustCompile(`.*time=\S+ level=\S+ prefix=\S+ binary=\S+ msg=.*`)
)

// cleanTerragruntOutput extracts the actual output value from terragrunt stack's verbose output
//
// Example input (raw terragrunt output):
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
	// Debug: Log raw output
	// fmt.Printf("DEBUG: Raw output: %q\n", rawOutput)
	// Remove terragrunt log lines
	cleaned := tgLogLevel.ReplaceAllString(rawOutput, "")

	lines := strings.Split(cleaned, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip empty lines and lines that are clearly log lines (containing msg= with log context)
		if trimmed != "" && !strings.Contains(line, " msg=") {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return "", nil
	}

	// Join all result lines
	finalOutput := strings.Join(result, "\n")

	// Check if it's JSON (starts with { or [)
	finalOutput = strings.TrimSpace(finalOutput)
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
// Example input (raw terragrunt JSON output):
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
	cleaned := tgLogLevel.ReplaceAllString(input, "")

	lines := strings.Split(cleaned, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip empty lines and lines that are clearly log lines (containing msg= with log context)
		if trimmed != "" && !strings.Contains(line, " msg=") {
			result = append(result, trimmed)
		}
	}
	ansiClean := strings.Join(result, "\n")

	var jsonObj interface{}
	if err := json.Unmarshal([]byte(ansiClean), &jsonObj); err != nil {
		return "", err
	}

	// Format JSON output with indentation
	normalized, err := json.MarshalIndent(jsonObj, "", "  ")
	if err != nil {
		return "", err
	}

	return string(normalized), nil
}
