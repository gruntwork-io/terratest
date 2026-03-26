package terragrunt_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terragrunt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsLogLine(t *testing.T) {
	t.Parallel()

	// Old format (time=... level=... msg=...)
	assert.True(t, terragrunt.IsLogLine("time=2026 level=info prefix=foo tf-path=terraform msg=Running"))

	// New format (HH:MM:SS.mmm LEVEL ...)
	assert.True(t, terragrunt.IsLogLine("20:41:53.564 INFO   Generating unit father"))
	assert.True(t, terragrunt.IsLogLine("20:41:53.564 WARN   Something is off"))
	assert.True(t, terragrunt.IsLogLine("20:41:53.564 DEBUG  Detailed info"))
	assert.True(t, terragrunt.IsLogLine("20:41:53.564 STDOUT [.terragrunt-stack/mother] terraform: output"))
	assert.True(t, terragrunt.IsLogLine("20:41:53.564 STDERR [foo] error message"))
	assert.True(t, terragrunt.IsLogLine("20:41:53.564 ERROR  Something went wrong"))
	assert.True(t, terragrunt.IsLogLine("20:41:53.564 TRACE  Very detailed"))

	// Not log lines
	assert.False(t, terragrunt.IsLogLine(`{"key": "value"}`))
	assert.False(t, terragrunt.IsLogLine(`{"message": "error msg=bad"}`))
	assert.False(t, terragrunt.IsLogLine("Group 1"))
	assert.False(t, terragrunt.IsLogLine("- Unit ./foo"))
}

func TestIsMetadataLine(t *testing.T) {
	t.Parallel()

	// Metadata lines
	assert.True(t, terragrunt.IsMetadataLine("Group 1"))
	assert.True(t, terragrunt.IsMetadataLine("Group 42"))
	assert.True(t, terragrunt.IsMetadataLine("- Unit ./foo"))
	assert.True(t, terragrunt.IsMetadataLine("- Unit ./.terragrunt-stack/mother"))

	// Not metadata lines
	assert.False(t, terragrunt.IsMetadataLine(`{"key": "value"}`))
	assert.False(t, terragrunt.IsMetadataLine("mother = { output = \"./test.txt\" }"))
	assert.False(t, terragrunt.IsMetadataLine("20:41:53.564 INFO   Running"))
}

func TestRemoveLogLines(t *testing.T) {
	t.Parallel()

	// Removes old format log lines, keeps JSON
	result := terragrunt.RemoveLogLines("time=2026 level=info msg=Start\n{\"key\": \"value\"}")
	assert.JSONEq(t, `{"key": "value"}`, result)

	// Removes new format log lines
	result = terragrunt.RemoveLogLines("20:41:53.564 INFO   Running\n{\"key\": \"value\"}")
	assert.JSONEq(t, `{"key": "value"}`, result)

	// Removes metadata lines (Group, Unit)
	result = terragrunt.RemoveLogLines("Group 1\n- Unit ./foo\n{\"key\": \"value\"}")
	assert.JSONEq(t, `{"key": "value"}`, result)

	// Preserves JSON with msg= in value
	result = terragrunt.RemoveLogLines("time=2026 level=info msg=Start\n{\"message\": \"error msg=bad\"}")
	assert.Contains(t, result, "error msg=bad")
}

func TestExtractJsonContent(t *testing.T) {
	t.Parallel()

	// Extracts JSON with old format, filters non-JSON
	input := "time=2026 level=info msg=Running\nGroup 1\n- Unit ./foo\n{\"a\": 1}\n{\"b\": 2}"
	result, err := terragrunt.ExtractJSONContent(input)
	require.NoError(t, err)
	assert.Contains(t, result, `"a": 1`)
	assert.Contains(t, result, `"b": 2`)
	assert.NotContains(t, result, "Group")
	assert.NotContains(t, result, "Unit")

	// Extracts JSON with new format logs
	input = "20:41:53.564 INFO   Running\n20:41:53.564 STDOUT terraform: done\n{\"key\": \"value\"}"
	result, err = terragrunt.ExtractJSONContent(input)
	require.NoError(t, err)
	assert.JSONEq(t, `{"key": "value"}`, result)

	// Handles nested JSON
	input = "time=2026 level=info msg=Running\n{\n  \"outer\": {\n    \"inner\": true\n  }\n}"
	result, err = terragrunt.ExtractJSONContent(input)
	require.NoError(t, err)
	assert.Contains(t, result, `"inner": true`)

	// Empty when only logs/metadata
	input = "20:41:53.564 INFO   Running\nGroup 1\n- Unit ./foo"
	result, err = terragrunt.ExtractJSONContent(input)
	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestCleanTerragruntOutput(t *testing.T) {
	t.Parallel()

	// Simple quoted string value
	input := "time=2026 level=info msg=Running\n\"my-bucket-name\""
	result := terragrunt.CleanTerragruntOutput(input)
	assert.Equal(t, "my-bucket-name", result)

	// JSON output preserved
	input = "20:41:53.564 INFO   Running\n{\"key\": \"value\"}"
	result = terragrunt.CleanTerragruntOutput(input)
	assert.JSONEq(t, `{"key": "value"}`, result)

	// Filters metadata lines
	input = "Group 1\n- Unit ./foo\n\"result\""
	result = terragrunt.CleanTerragruntOutput(input)
	assert.Equal(t, "result", result)

	// Empty input returns empty
	input = "20:41:53.564 INFO   Running"
	result = terragrunt.CleanTerragruntOutput(input)
	assert.Empty(t, result)
}

func TestCleanTerragruntJson(t *testing.T) {
	t.Parallel()

	// Valid single JSON with old format logs
	input := "time=2026 level=info msg=Running\n{\"mother\":{\"output\":\"test\"}}"
	result, err := terragrunt.CleanTerragruntJSON(input)
	require.NoError(t, err)
	assert.Contains(t, result, "mother")

	// Valid single JSON with new format logs (terragrunt 0.88+)
	input = "{\"a\": 1}\n20:41:53.564 INFO   Generating unit\n20:41:53.564 STDOUT terraform: done"
	result, err = terragrunt.CleanTerragruntJSON(input)
	require.NoError(t, err)
	assert.Contains(t, result, `"a": 1`)

	// Multiple JSON objects should error
	_, err = terragrunt.CleanTerragruntJSON("{\"a\": 1}\n{\"b\": 2}")
	require.Error(t, err)

	// Empty/no-JSON input should error (documents expected behavior)
	_, err = terragrunt.CleanTerragruntJSON("20:41:53.564 INFO   Running\nGroup 1\n- Unit ./foo")
	require.Error(t, err, "cleanTerragruntJson should error when input contains no JSON")
}

func TestCleanTerragruntOutputEdgeCases(t *testing.T) {
	t.Parallel()

	// Empty string value (terraform outputs "" for empty strings)
	input := "time=2026 level=info msg=Running\n\"\""
	result := terragrunt.CleanTerragruntOutput(input)
	assert.Empty(t, result, "Empty quoted string should become empty string")

	// Value with quotes inside (terraform outputs "\"quoted\"")
	input = "20:41:53.564 INFO   Running\n\"\\\"quoted\\\"\""
	result = terragrunt.CleanTerragruntOutput(input)
	assert.Equal(t, "\\\"quoted\\\"", result, "Escaped quotes should be preserved")

	// Multiple lines of non-JSON content after filtering logs
	input = "20:41:53.564 INFO   Running\nline1\nline2"
	result = terragrunt.CleanTerragruntOutput(input)
	assert.Equal(t, "line1\nline2", result)

	// Mismatched quotes: opening quote without closing quote should be left as-is
	input = "20:41:53.564 INFO   Running\n\"no-closing-quote"
	result = terragrunt.CleanTerragruntOutput(input)
	assert.Equal(t, "\"no-closing-quote", result, "Mismatched quotes should be preserved verbatim")

	// Closing quote without opening quote should be left as-is
	input = "20:41:53.564 INFO   Running\nno-opening-quote\""
	result = terragrunt.CleanTerragruntOutput(input)
	assert.Equal(t, "no-opening-quote\"", result, "Mismatched quotes should be preserved verbatim")

	// Array JSON output preserved
	input = "20:41:53.564 INFO   Running\n[\"a\", \"b\"]"
	result = terragrunt.CleanTerragruntOutput(input)
	assert.Equal(t, `["a", "b"]`, result, "JSON array should be preserved")
}

func TestExtractJsonContentMalformedJson(t *testing.T) {
	t.Parallel()

	// Valid JSON followed by malformed JSON: returns error
	input := "{\"valid\": true}\n{broken json"
	_, err := terragrunt.ExtractJSONContent(input)
	require.Error(t, err)

	// Malformed JSON only: returns error
	input = "{not valid json at all"
	_, err = terragrunt.ExtractJSONContent(input)
	require.Error(t, err)

	// Valid JSON with log lines before and after (realistic scenario)
	input = "time=2026 level=info msg=Before\n{\"key\": 1}\ntime=2026 level=info msg=After"
	result, err := terragrunt.ExtractJSONContent(input)
	require.NoError(t, err)
	assert.Contains(t, result, `"key"`)
	assert.NotContains(t, result, "Before")
	assert.NotContains(t, result, "After")

	// Whitespace-only after filtering logs
	input = "20:41:53.564 INFO   Running\n   \n  "
	result, err = terragrunt.ExtractJSONContent(input)
	require.NoError(t, err)
	assert.Empty(t, result)

	// Two valid JSON objects separated by log lines
	input = "20:41:53.564 INFO   Start\n{\"a\": 1}\n20:41:53.564 INFO   Middle\n{\"b\": 2}\n20:41:53.564 INFO   End"
	result, err = terragrunt.ExtractJSONContent(input)
	require.NoError(t, err)
	assert.Contains(t, result, `"a"`)
	assert.Contains(t, result, `"b"`)
}
