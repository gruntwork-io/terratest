package terragrunt_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terragrunt"
	"github.com/stretchr/testify/require"
)

// TestBuildRunArgs verifies the argument construction logic for the run command.
func TestBuildRunArgs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tgArgs   []string
		tfArgs   []string
		expected []string
	}{
		{
			name:     "only tf args",
			tgArgs:   []string{},
			tfArgs:   []string{"apply", "-auto-approve"},
			expected: []string{"--", "apply", "-auto-approve"},
		},
		{
			name:     "nil tg args with tf args",
			tgArgs:   nil,
			tfArgs:   []string{"plan"},
			expected: []string{"--", "plan"},
		},
		{
			name:     "both tg and tf args",
			tgArgs:   []string{"--all"},
			tfArgs:   []string{"apply", "-input=false", "-auto-approve"},
			expected: []string{"--all", "--", "apply", "-input=false", "-auto-approve"},
		},
		{
			name:     "multiple tg args",
			tgArgs:   []string{"--all", "--exclude-dir", "staging"},
			tfArgs:   []string{"plan"},
			expected: []string{"--all", "--exclude-dir", "staging", "--", "plan"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual := terragrunt.BuildRunArgs(tt.tgArgs, tt.tfArgs)
			require.Equal(t, tt.expected, actual)
		})
	}
}

// TestRunE_EmptyTfArgs verifies that RunE returns an error when tfArgs is empty.
func TestRunE_EmptyTfArgs(t *testing.T) {
	t.Parallel()

	options := &terragrunt.Options{
		TerragruntDir: "/some/path",
	}

	_, err := terragrunt.RunE(t, options, []string{}, []string{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "tfArgs cannot be empty")

	_, err = terragrunt.RunE(t, options, []string{"--all"}, nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "tfArgs cannot be empty")
}

// TestRun verifies that Run executes terragrunt run -- apply successfully.
func TestRun(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	defer terragrunt.Run(t, options, []string{}, []string{"destroy", "-auto-approve"})

	out := terragrunt.Run(t, options, []string{}, []string{"apply", "-input=false", "-auto-approve"})
	require.Contains(t, out, "Hello, World")
}

// TestRunE verifies that RunE returns an error on failure rather than calling t.Fatal.
func TestRunE(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	// Run an invalid tf command to trigger an error
	_, err = terragrunt.RunE(t, options, []string{}, []string{"not-a-real-command"})
	require.Error(t, err)
}

// TestRunWithTgArgs verifies that terragrunt-specific args are passed before the -- separator.
func TestRunWithTgArgs(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerragruntFolderToTemp("testdata/terragrunt-no-error", t.Name())
	require.NoError(t, err)

	options := &terragrunt.Options{
		TerragruntDir:    testFolder,
		TerragruntBinary: "terragrunt",
	}

	defer terragrunt.Run(t, options, []string{}, []string{"destroy", "-auto-approve"})

	// Use --log-level error as a tg arg to verify it's respected
	out := terragrunt.Run(t, options, []string{"--log-level", "error"}, []string{"apply", "-input=false", "-auto-approve"})
	require.Contains(t, out, "Hello, World")
	require.NotContains(t, out, "level=info",
		"With --log-level error, info logs should not appear")
}

// TestRunE_ValidationError verifies that RunE returns an error for invalid options.
func TestRunE_ValidationError(t *testing.T) {
	t.Parallel()

	// Missing TerragruntDir
	options := &terragrunt.Options{}
	_, err := terragrunt.RunE(t, options, []string{}, []string{"apply"})
	require.Error(t, err)
	require.Contains(t, err.Error(), "TerragruntDir is required")
}
