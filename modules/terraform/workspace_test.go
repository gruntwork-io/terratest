package terraform_test

import (
	"errors"
	"testing"

	"github.com/gruntwork-io/terratest/modules/files"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkspaceNew(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-workspace", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	out := terraform.WorkspaceSelectOrNew(t, options, "terratest")

	assert.Equal(t, "terratest", out)
}

func TestWorkspaceIllegalName(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-workspace", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	out, err := terraform.WorkspaceSelectOrNewE(t, options, "###@@@&&&")

	require.Error(t, err)
	assert.Empty(t, out, "%q should be an empty string", out)
}

func TestWorkspaceSelect(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-workspace", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	out := terraform.WorkspaceSelectOrNew(t, options, "terratest")
	assert.Equal(t, "terratest", out)

	out = terraform.WorkspaceSelectOrNew(t, options, "default")
	assert.Equal(t, "default", out)
}

func TestWorkspaceApply(t *testing.T) {
	t.Parallel()

	testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-workspace", t.Name())
	if err != nil {
		t.Fatal(err)
	}

	options := &terraform.Options{
		TerraformDir: testFolder,
	}

	terraform.WorkspaceSelectOrNew(t, options, "Terratest")
	out := terraform.InitAndApply(t, options)

	assert.Contains(t, out, "Hello, Terratest")
}

func TestIsExistingWorkspace(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		out      string
		name     string
		expected bool
	}{
		{out: "  default\n* foo\n", name: "default", expected: true},
		{out: "* default\n  foo\n", name: "default", expected: true},
		{out: "  foo\n* default\n", name: "default", expected: true},
		{out: "* foo\n  default\n", name: "default", expected: true},
		{out: "  foo\n* bar\n", name: "default", expected: false},
		{out: "* foo\n  bar\n", name: "default", expected: false},
		{out: "  default\n* foobar\n", name: "foo", expected: false},
		{out: "* default\n  foobar\n", name: "foo", expected: false},
		{out: "  default\n* foo\n", name: "foobar", expected: false},
		{out: "* default\n  foo\n", name: "foobar", expected: false},
		{out: "* default\n  foo\n", name: "foo", expected: true},
	}

	for _, testCase := range testCases {
		actual := terraform.IsExistingWorkspace(testCase.out, testCase.name)
		assert.Equal(t, testCase.expected, actual, "Out: %q, Name: %q", testCase.out, testCase.name)
	}
}

func TestWorkspaceDeleteE(t *testing.T) {
	t.Parallel()

	// state describes an expected status when a given testCase begins
	type state struct {
		current    string
		workspaces []string
	}

	// testCase describes a named test case with a state, args and expcted results
	type testCase struct {
		expectedError     error
		name              string
		toDeleteWorkspace string
		expectedCurrent   string
		initialState      state
	}

	testCases := []testCase{
		{
			name: "delete another existing workspace and stay on current",
			initialState: state{
				workspaces: []string{"staging", "production"},
				current:    "staging",
			},
			toDeleteWorkspace: "production",
			expectedCurrent:   "staging",
			expectedError:     nil,
		},
		{
			name: "delete current workspace and switch to a specified",
			initialState: state{
				workspaces: []string{"staging", "production"},
				current:    "production",
			},
			toDeleteWorkspace: "production",
			expectedCurrent:   "default",
			expectedError:     nil,
		},
		{
			name: "delete a non existing workspace should trigger an error",
			initialState: state{
				workspaces: []string{"staging", "production"},
				current:    "staging",
			},
			toDeleteWorkspace: "hellothere",
			expectedCurrent:   "staging",
			expectedError:     terraform.WorkspaceDoesNotExist("hellothere"),
		},
		{
			name: "delete the default workspace triggers an error",
			initialState: state{
				workspaces: []string{"staging", "production"},
				current:    "staging",
			},
			toDeleteWorkspace: "default",
			expectedCurrent:   "staging",
			expectedError:     &terraform.UnsupportedDefaultWorkspaceDeletion{},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			testFolder, err := files.CopyTerraformFolderToTemp("../../test/fixtures/terraform-workspace", testCase.name)
			require.NoError(t, err)

			options := &terraform.Options{
				TerraformDir: testFolder,
			}

			// Set up pre-existing environment based on test case description
			for _, existingWorkspace := range testCase.initialState.workspaces {
				_, err = terraform.RunTerraformCommandE(t, options, "workspace", "new", existingWorkspace)
				require.NoError(t, err)
			}

			// Switch to the specified workspace
			_, err = terraform.RunTerraformCommandE(t, options, "workspace", "select", testCase.initialState.current)
			require.NoError(t, err)

			// Testing time, wooohoooo
			gotResult, gotErr := terraform.WorkspaceDeleteE(t, options, testCase.toDeleteWorkspace)

			// Check for errors
			if testCase.expectedError != nil {
				if !errors.Is(gotErr, testCase.expectedError) {
					t.Errorf("expected error: %v, got error: %v", testCase.expectedError, gotErr)
				}
			} else {
				require.NoError(t, gotErr)
				// Check for results
				assert.Equal(t, testCase.expectedCurrent, gotResult)
				assert.False(t, terraform.IsExistingWorkspace(terraform.RunTerraformCommand(t, options, "workspace", "list"), testCase.toDeleteWorkspace))
			}
		})
	}
}
