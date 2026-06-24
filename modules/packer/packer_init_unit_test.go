package packer //nolint:testpackage // white-box test for the unexported packerInitSupportsTemplate helper

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackerInitSupportsTemplate(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()

	hclFile := filepath.Join(dir, "build.pkr.hcl")
	require.NoError(t, os.WriteFile(hclFile, []byte("# hcl2 template"), 0644))

	jsonFile := filepath.Join(dir, "template.json")
	require.NoError(t, os.WriteFile(jsonFile, []byte("{}"), 0644))

	testCases := []struct {
		options  *Options
		name     string
		expected bool
	}{
		{
			name:     "hcl file",
			options:  &Options{Template: hclFile},
			expected: true,
		},
		{
			name:     "directory as absolute path",
			options:  &Options{Template: dir},
			expected: true,
		},
		{
			name:     "directory relative to working dir",
			options:  &Options{Template: ".", WorkingDir: dir},
			expected: true,
		},
		{
			name:     "hcl file relative to working dir",
			options:  &Options{Template: "build.pkr.hcl", WorkingDir: dir},
			expected: true,
		},
		{
			name:     "legacy json template",
			options:  &Options{Template: jsonFile},
			expected: false,
		},
		{
			name:     "nonexistent non-hcl path",
			options:  &Options{Template: filepath.Join(dir, "does-not-exist")},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, packerInitSupportsTemplate(tc.options))
		})
	}
}
