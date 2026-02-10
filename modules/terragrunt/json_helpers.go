package terragrunt

import (
	"encoding/json"
	"regexp"
	"strings"
)

// removeLogLines removes terragrunt log lines and metadata from output
func removeLogLines(rawOutput string) string {
	lines := strings.Split(rawOutput, "\n")
	var result []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip empty lines, terragrunt log lines, and metadata lines
		if trimmed == "" {
			continue
		}
		if isLogLine(trimmed) || isMetadataLine(trimmed) {
			continue
		}
		result = append(result, trimmed)
	}

	return strings.Join(result, "\n")
}

// isMetadataLine checks if a line is terragrunt metadata (e.g., "Group 1", "- Unit ./foo")
func isMetadataLine(line string) bool {
	return tgMetadataPattern.MatchString(line)
}

// newLogLinePattern matches the new terragrunt log format: "HH:MM:SS.mmm LEVEL ..."
// Example: "20:41:53.564 INFO   Generating unit father..."
var newLogLinePattern = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}\.\d{3}\s+(INFO|WARN|ERROR|DEBUG|TRACE|STDOUT|STDERR)\s+`)

// tgMetadataPattern matches terragrunt metadata lines like "Group 1" or "- Unit ./foo"
var tgMetadataPattern = regexp.MustCompile(`^(Group \d+|- Unit )`)

// isLogLine checks if a line is a terragrunt log line
func isLogLine(line string) bool {
	// Old format: time=... level=... msg=...
	if strings.HasPrefix(line, "time=") && strings.Contains(line, "level=") && strings.Contains(line, "msg=") {
		return true
	}
	// New format (terragrunt 0.88+): HH:MM:SS.mmm LEVEL message
	return newLogLinePattern.MatchString(line)
}

// extractJsonContent extracts only JSON objects from terragrunt output,
// filtering out log lines and other non-JSON content like "Group 1" or "- Unit ./foo"
func extractJsonContent(rawOutput string) string {
	lines := strings.Split(rawOutput, "\n")
	var result []string
	braceDepth := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Skip empty lines, log lines, and metadata lines
		if trimmed == "" || isLogLine(trimmed) || isMetadataLine(trimmed) {
			continue
		}

		// Count braces to track JSON depth
		openBraces := strings.Count(trimmed, "{")
		closeBraces := strings.Count(trimmed, "}")

		// If we're starting a new JSON object or inside one, include the line
		if openBraces > 0 || braceDepth > 0 {
			result = append(result, line)
		}

		braceDepth += openBraces - closeBraces
		if braceDepth < 0 {
			braceDepth = 0
		}
	}

	return strings.Join(result, "\n")
}

// cleanTerragruntOutput extracts the actual output value from terragrunt stack's verbose output
//
// Example input (raw tg output):
//
//	time=2023-07-11T10:30:45Z level=info prefix=foo tf-path=terraform msg=Initializing...
//	time=2023-07-11T10:30:46Z level=info prefix=foo tf-path=terraform msg=Running command...
//	"my-bucket-name"
//
// Example output (cleaned):
//
//	my-bucket-name
//
// For JSON values, it preserves the structure:
// Input:
//
//	time=2023-07-11T10:30:45Z level=info prefix=foo tf-path=terraform msg=Running...
//	{"vpc_id": "vpc-12345", "subnet_ids": ["subnet-1", "subnet-2"]}
//
// Output:
//
//	{"vpc_id": "vpc-12345", "subnet_ids": ["subnet-1", "subnet-2"]}
func cleanTerragruntOutput(rawOutput string) (string, error) {
	// Remove terragrunt log lines and metadata
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
	// Use TrimPrefix/TrimSuffix to remove exactly one quote from each end
	if strings.HasPrefix(finalOutput, "\"") && strings.HasSuffix(finalOutput, "\"") {
		finalOutput = strings.TrimPrefix(finalOutput, "\"")
		finalOutput = strings.TrimSuffix(finalOutput, "\"")
	}

	return finalOutput, nil
}

// cleanTerragruntJson cleans the JSON output from terragrunt stack command
//
// Example input (raw tg JSON output):
//
//	time=2023-07-11T10:30:45Z level=info prefix=mother tf-path=terraform msg=Initializing...
//	time=2023-07-11T10:30:46Z level=info prefix=mother tf-path=terraform msg=Running command...
//	{"mother":{"output":"./test.txt"},"father":{"output":"./test.txt"}}
//
// Example output (cleaned and formatted):
//
//	{
//	  "mother": {
//	    "output": "./test.txt"
//	  },
//	  "father": {
//	    "output": "./test.txt"
//	  }
//	}
func cleanTerragruntJson(input string) (string, error) {
	// Extract only JSON content, filtering out log lines and other non-JSON content
	cleaned := extractJsonContent(input)

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
