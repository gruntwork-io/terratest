package ssh //nolint:testpackage // white-box test for unexported shellQuote helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellQuote(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no special chars",
			input:    "/etc/hostname",
			expected: "'/etc/hostname'",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "''",
		},
		{
			name:     "single quote",
			input:    "it's",
			expected: `'it'\''s'`,
		},
		{
			name:     "spaces",
			input:    "/path with spaces/file.txt",
			expected: "'/path with spaces/file.txt'",
		},
		{
			name:     "path traversal attempt",
			input:    "a; rm -rf /",
			expected: "'a; rm -rf /'",
		},
		{
			name:     "command substitution attempt",
			input:    "$(rm -rf /)",
			expected: "'$(rm -rf /)'",
		},
		{
			name:     "backtick attempt",
			input:    "`rm -rf /`",
			expected: "'`rm -rf /`'",
		},
		{
			name:     "injection with embedded single quote",
			input:    "'; rm -rf / #",
			expected: `''\''; rm -rf / #'`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, shellQuote(tc.input))
		})
	}
}
