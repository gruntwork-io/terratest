package terraform

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// WorkspaceSelectOrNew runs terraform workspace with the given options and the workspace name
// and returns a name of the current workspace. It tries to select a workspace with the given
// name, or it creates a new one if it doesn't exist.
func WorkspaceSelectOrNew(t testing.TestingT, options *Options, name string) string {
	out, err := WorkspaceSelectOrNewE(t, options, name)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// WorkspaceSelectOrNewE runs terraform workspace with the given options and the workspace name
// and returns a name of the current workspace. It tries to select a workspace with the given
// name, or it creates a new one if it doesn't exist.
func WorkspaceSelectOrNewE(t testing.TestingT, options *Options, name string) (string, error) {
	out, err := RunTerraformCommandE(t, options, "workspace", "list")
	if err != nil {
		return "", err
	}

	if isExistingWorkspace(out, name) {
		_, err = RunTerraformCommandE(t, options, "workspace", "select", name)
	} else {
		_, err = RunTerraformCommandE(t, options, "workspace", "new", name)
	}
	if err != nil {
		return "", err
	}

	return RunTerraformCommandE(t, options, "workspace", "show")
}

func isExistingWorkspace(out string, name string) bool {
	workspaces := strings.Split(out, "\n")
	for _, ws := range workspaces {
		if nameMatchesWorkspace(name, ws) {
			return true
		}
	}
	return false
}

func nameMatchesWorkspace(name string, workspace string) bool {
	// Regex for matching workspace should match for strings with optional leading asterisk "*"
	// following optional white spaces following the workspace name.
	// E.g. for the given name "terratest", following strings will match:
	//
	//    "* terratest"
	//    "  terratest"
	match, _ := regexp.MatchString(fmt.Sprintf("^\\*?\\s*%s$", name), workspace)
	return match
}

// WorkspaceDelete removes the specified terraform workspace with the given options.
// It returns the name of the current workspace AFTER deletion, and the returned error (that can be nil).
// If the workspace to delete is the current one, then it tries to switch to the "default" workspace.
// Deleting the workspace "default" is not supported.
func WorkspaceDeleteE(t testing.TestingT, options *Options, name string) (string, error) {
	currentWorkspace, err := RunTerraformCommandE(t, options, "workspace", "show")
	if err != nil {
		return currentWorkspace, err
	}

	if name == "default" {
		return currentWorkspace, &UnsupportedDefaultWorkspaceDeletion{}
	}

	out, err := RunTerraformCommandE(t, options, "workspace", "list")
	if err != nil {
		return currentWorkspace, err
	}
	if !isExistingWorkspace(out, name) {
		return currentWorkspace, WorkspaceDoesNotExist(name)
	}

	// Switch workspace before deleting if it is the current
	if currentWorkspace == name {
		currentWorkspace, err = WorkspaceSelectOrNewE(t, options, "default")
		if err != nil {
			return currentWorkspace, err
		}
	}

	// delete workspace
	_, err = RunTerraformCommandE(t, options, "workspace", "delete", name)

	return currentWorkspace, err
}

// WorkspaceDelete removes the specified terraform workspace with the given options.
// It returns the name of the current workspace AFTER deletion.
// If the workspace to delete is the current one, then it tries to switch to the "default" workspace.
// Deleting the workspace "default" is not supported and only return an empty string (to avoid a fatal error).
func WorkspaceDelete(t testing.TestingT, options *Options, name string) string {
	out, err := WorkspaceDeleteE(t, options, name)
	require.NoError(t, err)
	return out
}
