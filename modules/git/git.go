// Package git allows to interact with Git.
//
// Deprecated: The git package is scheduled for removal in Terratest v2. Each
// helper here wraps a single git command (for example, git rev-parse or
// git describe); call git directly with os/exec instead. See the Terratest v2
// migration notes for details.
package git

import (
	"context"
	"os"
	"os/exec"
	"strings"

	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetCurrentBranchName retrieves the current branch name or an empty string
// in case of detached state. Fails the test if an error occurs.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetCurrentBranchName(t testing.TestingT) string {
	return GetCurrentBranchNameContext(t, context.Background(), "")
}

// GetCurrentBranchNameContext retrieves the current branch name or an empty
// string in case of detached state. The dir parameter specifies the working
// directory for the git command; if empty, the process working directory is
// used. Fails the test if an error occurs.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetCurrentBranchNameContext(t testing.TestingT, ctx context.Context, dir string) string {
	out, err := GetCurrentBranchNameContextE(t, ctx, dir)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// GetCurrentBranchNameE retrieves the current branch name or an empty string
// in case of detached state. Uses git branch --show-current, which was
// introduced in git v2.22. Falls back to git rev-parse for older versions.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetCurrentBranchNameE(t testing.TestingT) (string, error) {
	return GetCurrentBranchNameContextE(t, context.Background(), "")
}

// GetCurrentBranchNameContextE retrieves the current branch name or an empty
// string in case of detached state. Uses git branch --show-current, which was
// introduced in git v2.22. Falls back to git rev-parse for older versions.
// The dir parameter specifies the working directory for the git command; if
// empty, the process working directory is used.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetCurrentBranchNameContextE(t testing.TestingT, ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "branch", "--show-current")
	cmd.Dir = dir

	bytes, err := cmd.Output()
	if err != nil {
		return GetCurrentBranchNameOldContextE(t, ctx, dir)
	}

	name := strings.TrimSpace(string(bytes))
	if name == "HEAD" {
		return "", nil
	}

	return name, nil
}

// GetCurrentBranchNameOldE retrieves the current branch name or an empty
// string in case of detached state using git rev-parse --abbrev-ref HEAD.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetCurrentBranchNameOldE(t testing.TestingT) (string, error) {
	return GetCurrentBranchNameOldContextE(t, context.Background(), "")
}

// GetCurrentBranchNameOldContextE retrieves the current branch name or an
// empty string in case of detached state using git rev-parse --abbrev-ref HEAD.
// This is a fallback for git versions older than v2.22 that lack
// git branch --show-current. The dir parameter specifies the working directory
// for the git command; if empty, the process working directory is used.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetCurrentBranchNameOldContextE(t testing.TestingT, ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir

	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	name := strings.TrimSpace(string(bytes))
	if name == "HEAD" {
		return "", nil
	}

	return name, nil
}

// GetCurrentGitRef retrieves the current branch name, lightweight
// (non-annotated) tag, or exact tag value if the tag points to the current
// commit. Fails the test if an error occurs.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetCurrentGitRef(t testing.TestingT) string {
	return GetCurrentGitRefContext(t, context.Background(), "")
}

// GetCurrentGitRefContext retrieves the current branch name, lightweight
// (non-annotated) tag, or exact tag value if the tag points to the current
// commit. The dir parameter specifies the working directory for the git
// command; if empty, the process working directory is used. Fails the test if
// an error occurs.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetCurrentGitRefContext(t testing.TestingT, ctx context.Context, dir string) string {
	out, err := GetCurrentGitRefContextE(t, ctx, dir)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// GetCurrentGitRefE retrieves the current branch name, lightweight
// (non-annotated) tag, or exact tag value if the tag points to the current
// commit.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetCurrentGitRefE(t testing.TestingT) (string, error) {
	return GetCurrentGitRefContextE(t, context.Background(), "")
}

// GetCurrentGitRefContextE retrieves the current branch name, lightweight
// (non-annotated) tag, or exact tag value if the tag points to the current
// commit. The dir parameter specifies the working directory for the git
// command; if empty, the process working directory is used.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetCurrentGitRefContextE(t testing.TestingT, ctx context.Context, dir string) (string, error) {
	out, err := GetCurrentBranchNameContextE(t, ctx, dir)
	if err != nil {
		return "", err
	}

	if out != "" {
		return out, nil
	}

	out, err = GetTagContextE(t, ctx, dir)
	if err != nil {
		return "", err
	}

	return out, nil
}

// GetTagE retrieves the lightweight (non-annotated) tag or exact tag value if
// the tag points to the current commit.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetTagE(t testing.TestingT) (string, error) {
	return GetTagContextE(t, context.Background(), "")
}

// GetTagContextE retrieves the lightweight (non-annotated) tag or exact tag
// value if the tag points to the current commit. The dir parameter specifies
// the working directory for the git command; if empty, the process working
// directory is used.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetTagContextE(t testing.TestingT, ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "describe", "--tags")
	cmd.Dir = dir

	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytes)), nil
}

// GetRepoRoot retrieves the path to the root directory of the repo. Fails the
// test if there is an error.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetRepoRoot(t testing.TestingT) string {
	return GetRepoRootContext(t, context.Background(), "")
}

// GetRepoRootContext retrieves the path to the root directory of the repo. The
// dir parameter specifies the working directory for the git command; if empty,
// the process working directory is used. Fails the test if there is an error.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetRepoRootContext(t testing.TestingT, ctx context.Context, dir string) string {
	out, err := GetRepoRootContextE(t, ctx, dir)
	require.NoError(t, err)

	return out
}

// GetRepoRootE retrieves the path to the root directory of the repo.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetRepoRootE(t testing.TestingT) (string, error) {
	return GetRepoRootContextE(t, context.Background(), "")
}

// GetRepoRootContextE retrieves the path to the root directory of the repo.
// The dir parameter specifies the working directory for the git command; if
// empty, the process working directory is used.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetRepoRootContextE(t testing.TestingT, ctx context.Context, dir string) (string, error) {
	if dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}

		dir = cwd
	}

	return GetRepoRootForDirContextE(t, ctx, dir)
}

// GetRepoRootForDir retrieves the path to the root directory of the repo in
// which dir resides. Fails the test if there is an error.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetRepoRootForDir(t testing.TestingT, dir string) string {
	return GetRepoRootForDirContext(t, context.Background(), dir)
}

// GetRepoRootForDirContext retrieves the path to the root directory of the
// repo in which dir resides. Fails the test if there is an error.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetRepoRootForDirContext(t testing.TestingT, ctx context.Context, dir string) string {
	out, err := GetRepoRootForDirContextE(t, ctx, dir)
	require.NoError(t, err)

	return out
}

// GetRepoRootForDirE retrieves the path to the root directory of the repo in
// which dir resides.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetRepoRootForDirE(t testing.TestingT, dir string) (string, error) {
	return GetRepoRootForDirContextE(t, context.Background(), dir)
}

// GetRepoRootForDirContextE retrieves the path to the root directory of the
// repo in which dir resides.
//
// Deprecated: scheduled for removal in Terratest v2. Shell out to git directly
// with os/exec instead.
func GetRepoRootForDirContextE(t testing.TestingT, ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir

	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytes)), nil
}
