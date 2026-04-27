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
			name:     "embedded single quote",
			input:    "it's",
			expected: `'it'\''s'`,
		},
		{
			name:     "spaces",
			input:    "/path with spaces/file.txt",
			expected: "'/path with spaces/file.txt'",
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
