package terragrunt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsLogLine(t *testing.T) {
	t.Parallel()

	// Old format (time=... level=... msg=...)
	assert.True(t, isLogLine("time=2026 level=info prefix=foo tf-path=terraform msg=Running"))

	// New format (HH:MM:SS.mmm LEVEL ...)
	assert.True(t, isLogLine("20:41:53.564 INFO   Generating unit father"))
	assert.True(t, isLogLine("20:41:53.564 WARN   Something is off"))
	assert.True(t, isLogLine("20:41:53.564 DEBUG  Detailed info"))
	assert.True(t, isLogLine("20:41:53.564 STDOUT [.terragrunt-stack/mother] terraform: output"))
	assert.True(t, isLogLine("20:41:53.564 STDERR [foo] error message"))
	assert.True(t, isLogLine("20:41:53.564 ERROR  Something went wrong"))
	assert.True(t, isLogLine("20:41:53.564 TRACE  Very detailed"))

	// Not log lines
	assert.False(t, isLogLine(`{"key": "value"}`))
	assert.False(t, isLogLine(`{"message": "error msg=bad"}`))
	assert.False(t, isLogLine("Group 1"))
	assert.False(t, isLogLine("- Unit ./foo"))
}

func TestIsMetadataLine(t *testing.T) {
	t.Parallel()

	// Metadata lines
	assert.True(t, isMetadataLine("Group 1"))
	assert.True(t, isMetadataLine("Group 42"))
	assert.True(t, isMetadataLine("- Unit ./foo"))
	assert.True(t, isMetadataLine("- Unit ./.terragrunt-stack/mother"))

	// Not metadata lines
	assert.False(t, isMetadataLine(`{"key": "value"}`))
	assert.False(t, isMetadataLine("mother = { output = \"./test.txt\" }"))
	assert.False(t, isMetadataLine("20:41:53.564 INFO   Running"))
}

func TestRemoveLogLines(t *testing.T) {
	t.Parallel()

	// Removes old format log lines, keeps JSON
	result := removeLogLines("time=2026 level=info msg=Start\n{\"key\": \"value\"}")
	assert.Equal(t, `{"key": "value"}`, result)

	// Removes new format log lines
	result = removeLogLines("20:41:53.564 INFO   Running\n{\"key\": \"value\"}")
	assert.Equal(t, `{"key": "value"}`, result)

	// Removes metadata lines (Group, Unit)
	result = removeLogLines("Group 1\n- Unit ./foo\n{\"key\": \"value\"}")
	assert.Equal(t, `{"key": "value"}`, result)

	// Preserves JSON with msg= in value
	result = removeLogLines("time=2026 level=info msg=Start\n{\"message\": \"error msg=bad\"}")
	assert.Contains(t, result, "error msg=bad")
}

func TestExtractJsonContent(t *testing.T) {
	t.Parallel()

	// Extracts JSON with old format, filters non-JSON
	input := "time=2026 level=info msg=Running\nGroup 1\n- Unit ./foo\n{\"a\": 1}\n{\"b\": 2}"
	result := extractJsonContent(input)
	assert.Contains(t, result, `"a": 1`)
	assert.Contains(t, result, `"b": 2`)
	assert.NotContains(t, result, "Group")
	assert.NotContains(t, result, "Unit")

	// Extracts JSON with new format logs
	input = "20:41:53.564 INFO   Running\n20:41:53.564 STDOUT terraform: done\n{\"key\": \"value\"}"
	result = extractJsonContent(input)
	assert.Equal(t, `{"key": "value"}`, result)

	// Handles nested JSON
	input = "time=2026 level=info msg=Running\n{\n  \"outer\": {\n    \"inner\": true\n  }\n}"
	result = extractJsonContent(input)
	assert.Contains(t, result, `"inner": true`)

	// Empty when only logs/metadata
	input = "20:41:53.564 INFO   Running\nGroup 1\n- Unit ./foo"
	result = extractJsonContent(input)
	assert.Equal(t, "", result)
}

func TestCleanTerragruntOutput(t *testing.T) {
	t.Parallel()

	// Simple quoted string value
	input := "time=2026 level=info msg=Running\n\"my-bucket-name\""
	result, err := cleanTerragruntOutput(input)
	require.NoError(t, err)
	assert.Equal(t, "my-bucket-name", result)

	// JSON output preserved
	input = "20:41:53.564 INFO   Running\n{\"key\": \"value\"}"
	result, err = cleanTerragruntOutput(input)
	require.NoError(t, err)
	assert.Equal(t, `{"key": "value"}`, result)

	// Filters metadata lines
	input = "Group 1\n- Unit ./foo\n\"result\""
	result, err = cleanTerragruntOutput(input)
	require.NoError(t, err)
	assert.Equal(t, "result", result)

	// Empty input returns empty
	input = "20:41:53.564 INFO   Running"
	result, err = cleanTerragruntOutput(input)
	require.NoError(t, err)
	assert.Equal(t, "", result)
}

func TestCleanTerragruntJson(t *testing.T) {
	t.Parallel()

	// Valid single JSON with old format logs
	input := "time=2026 level=info msg=Running\n{\"mother\":{\"output\":\"test\"}}"
	result, err := cleanTerragruntJson(input)
	require.NoError(t, err)
	assert.Contains(t, result, "mother")

	// Valid single JSON with new format logs (terragrunt 0.88+)
	input = "{\"a\": 1}\n20:41:53.564 INFO   Generating unit\n20:41:53.564 STDOUT terraform: done"
	result, err = cleanTerragruntJson(input)
	require.NoError(t, err)
	assert.Contains(t, result, `"a": 1`)

	// Multiple JSON objects should error
	_, err = cleanTerragruntJson("{\"a\": 1}\n{\"b\": 2}")
	require.Error(t, err)

	// Empty/no-JSON input should error (documents expected behavior)
	_, err = cleanTerragruntJson("20:41:53.564 INFO   Running\nGroup 1\n- Unit ./foo")
	require.Error(t, err, "cleanTerragruntJson should error when input contains no JSON")
}

func TestCleanTerragruntOutputEdgeCases(t *testing.T) {
	t.Parallel()

	// Empty string value (terraform outputs "" for empty strings)
	input := "time=2026 level=info msg=Running\n\"\""
	result, err := cleanTerragruntOutput(input)
	require.NoError(t, err)
	assert.Equal(t, "", result, "Empty quoted string should become empty string")

	// Value with quotes inside (terraform outputs "\"quoted\"")
	input = "20:41:53.564 INFO   Running\n\"\\\"quoted\\\"\""
	result, err = cleanTerragruntOutput(input)
	require.NoError(t, err)
	assert.Equal(t, "\\\"quoted\\\"", result, "Escaped quotes should be preserved")

	// Multiple lines of non-JSON content after filtering logs
	input = "20:41:53.564 INFO   Running\nline1\nline2"
	result, err = cleanTerragruntOutput(input)
	require.NoError(t, err)
	assert.Equal(t, "line1\nline2", result)
}
