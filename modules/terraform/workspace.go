package terraform

import (
	"context"
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// WorkspaceSelectOrNewContext runs terraform workspace with the given options and the workspace name
// and returns the name of the current workspace. It tries to select a workspace with the given
// name, or it creates a new one if it doesn't exist. The context argument can be used for
// cancellation or timeout control.
func WorkspaceSelectOrNewContext(t testing.TestingT, ctx context.Context, options *Options, name string) string {
	out, err := WorkspaceSelectOrNewContextE(t, ctx, options, name)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// WorkspaceSelectOrNewContextE runs terraform workspace with the given options and the workspace name
// and returns the name of the current workspace. It tries to select a workspace with the given
// name, or it creates a new one if it doesn't exist. The context argument can be used for
// cancellation or timeout control.
func WorkspaceSelectOrNewContextE(t testing.TestingT, ctx context.Context, options *Options, name string) (string, error) {
	out, err := RunTerraformCommandContextE(t, ctx, options, "workspace", "list")
	if err != nil {
		return "", err
	}

	if IsExistingWorkspace(out, name) {
		_, err = RunTerraformCommandContextE(t, ctx, options, prepend(options.ExtraArgs.WorkspaceSelect, "workspace", "select", name)...)
	} else {
		_, err = RunTerraformCommandContextE(t, ctx, options, prepend(options.ExtraArgs.WorkspaceNew, "workspace", "new", name)...)
	}

	if err != nil {
		return "", err
	}

	return RunTerraformCommandContextE(t, ctx, options, "workspace", "show")
}

// WorkspaceDeleteContext removes the specified terraform workspace with the given options.
// It returns the name of the current workspace AFTER deletion.
// If the workspace to delete is the current one, then it tries to switch to the "default" workspace.
// Deleting the workspace "default" is not supported and only returns an empty string (to avoid a fatal error).
// The context argument can be used for cancellation or timeout control.
func WorkspaceDeleteContext(t testing.TestingT, ctx context.Context, options *Options, name string) string {
	out, err := WorkspaceDeleteContextE(t, ctx, options, name)
	require.NoError(t, err)

	return out
}

// WorkspaceDeleteContextE removes the specified terraform workspace with the given options.
// It returns the name of the current workspace AFTER deletion, and the returned error (that can be nil).
// If the workspace to delete is the current one, then it tries to switch to the "default" workspace.
// Deleting the workspace "default" is not supported. The context argument can be used for cancellation
// or timeout control.
func WorkspaceDeleteContextE(t testing.TestingT, ctx context.Context, options *Options, name string) (string, error) {
	currentWorkspace, err := RunTerraformCommandContextE(t, ctx, options, "workspace", "show")
	if err != nil {
		return currentWorkspace, err
	}

	if name == "default" {
		return currentWorkspace, &UnsupportedDefaultWorkspaceDeletion{}
	}

	out, err := RunTerraformCommandContextE(t, ctx, options, "workspace", "list")
	if err != nil {
		return currentWorkspace, err
	}

	if !IsExistingWorkspace(out, name) {
		return currentWorkspace, WorkspaceDoesNotExist(name)
	}

	// Switch workspace before deleting if it is the current
	if currentWorkspace == name {
		currentWorkspace, err = WorkspaceSelectOrNewContextE(t, ctx, options, "default")
		if err != nil {
			return currentWorkspace, err
		}
	}

	// delete workspace
	_, err = RunTerraformCommandContextE(t, ctx, options, prepend(options.ExtraArgs.WorkspaceDelete, "workspace", "delete", name)...)

	return currentWorkspace, err
}

// IsExistingWorkspace checks if a workspace with the given name exists in the terraform workspace list output.
func IsExistingWorkspace(out string, name string) bool {
	workspaces := strings.Split(out, "\n")

	for _, ws := range workspaces {
		if strings.HasSuffix(ws, name) {
			return true
		}
	}

	return false
}

// WorkspaceSelectOrNew runs terraform workspace with the given options and the workspace name
// and returns a name of the current workspace. It tries to select a workspace with the given
// name, or it creates a new one if it doesn't exist.
//
// Deprecated: Use [WorkspaceSelectOrNewContext] instead.
func WorkspaceSelectOrNew(t testing.TestingT, options *Options, name string) string {
	return WorkspaceSelectOrNewContext(t, context.Background(), options, name)
}

// WorkspaceSelectOrNewE runs terraform workspace with the given options and the workspace name
// and returns a name of the current workspace. It tries to select a workspace with the given
// name, or it creates a new one if it doesn't exist.
//
// Deprecated: Use [WorkspaceSelectOrNewContextE] instead.
func WorkspaceSelectOrNewE(t testing.TestingT, options *Options, name string) (string, error) {
	return WorkspaceSelectOrNewContextE(t, context.Background(), options, name)
}

// WorkspaceDeleteE removes the specified terraform workspace with the given options.
// It returns the name of the current workspace AFTER deletion, and the returned error (that can be nil).
// If the workspace to delete is the current one, then it tries to switch to the "default" workspace.
// Deleting the workspace "default" is not supported.
//
// Deprecated: Use [WorkspaceDeleteContextE] instead.
func WorkspaceDeleteE(t testing.TestingT, options *Options, name string) (string, error) {
	return WorkspaceDeleteContextE(t, context.Background(), options, name)
}

// WorkspaceDelete removes the specified terraform workspace with the given options.
// It returns the name of the current workspace AFTER deletion.
// If the workspace to delete is the current one, then it tries to switch to the "default" workspace.
// Deleting the workspace "default" is not supported and only return an empty string (to avoid a fatal error).
//
// Deprecated: Use [WorkspaceDeleteContext] instead.
func WorkspaceDelete(t testing.TestingT, options *Options, name string) string {
	return WorkspaceDeleteContext(t, context.Background(), options, name)
}
