// Package version_checker provides utilities for checking binary versions against constraints.
package version_checker //nolint:staticcheck // package name matches directory convention

import (
	"context"
	"fmt"
	"regexp"

	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/hashicorp/go-version"
	"github.com/stretchr/testify/require"
)

// VersionCheckerBinary is an enum for supported version checking.
type VersionCheckerBinary int

// List of binaries supported for version checking.
const (
	Docker VersionCheckerBinary = iota
	Terraform
	Packer
)

const (
	// versionRegexMatcher is a regex used to extract version string from shell command output.
	versionRegexMatcher = `\d+(\.\d+)+`
	// defaultVersionArg is a default arg to pass in to get version output from shell command.
	defaultVersionArg = "--version"
)

// CheckVersionParams contains the parameters for checking a binary's version.
type CheckVersionParams struct {
	// BinaryPath is a path to the binary you want to check the version for.
	BinaryPath string
	// VersionConstraint is a string literal containing one or more conditions, which are separated by commas.
	// More information here:https://www.terraform.io/language/expressions/version-constraints
	VersionConstraint string
	// WorkingDir is a directory you want to run the shell command.
	WorkingDir string
	// Binary is the name of the binary you want to check the version for.
	Binary VersionCheckerBinary
}

// CheckVersionE checks whether the given Binary version is greater than or equal
// to the given expected version.
//
// Deprecated: Use CheckVersionContextE instead.
func CheckVersionE(
	t testing.TestingT,
	params CheckVersionParams) error {
	return CheckVersionContextE(t, context.Background(), params)
}

// CheckVersionContextE is like CheckVersionE but includes a context.
func CheckVersionContextE(
	t testing.TestingT,
	ctx context.Context,
	params CheckVersionParams) error {
	if err := validateParams(params); err != nil {
		return err
	}

	binaryVersion, err := getVersionWithShellCommand(t, ctx, params)
	if err != nil {
		return err
	}

	return CheckVersionConstraint(binaryVersion, params.VersionConstraint)
}

// CheckVersion checks whether the given Binary version is greater than or equal to the
// given expected version and fails if it's not.
//
// Deprecated: Use CheckVersionContext instead.
func CheckVersion(
	t testing.TestingT,
	params CheckVersionParams) {
	CheckVersionContext(t, context.Background(), params)
}

// CheckVersionContext is like CheckVersion but includes a context.
func CheckVersionContext(
	t testing.TestingT,
	ctx context.Context,
	params CheckVersionParams) {
	require.NoError(t, CheckVersionContextE(t, ctx, params))
}

// validateParams checks whether the given params contain valid data.
func validateParams(params CheckVersionParams) error {
	if params.WorkingDir == "" {
		return &MissingParamErr{Param: "WorkingDir"}
	} else if params.VersionConstraint == "" {
		return &MissingParamErr{Param: "VersionConstraint"}
	}

	if _, err := version.NewConstraint(params.VersionConstraint); err != nil {
		return &InvalidVersionConstraintErr{Constraint: params.VersionConstraint, Underlying: err}
	}

	return nil
}

// getVersionWithShellCommand get version by running a shell command.
func getVersionWithShellCommand(t testing.TestingT, ctx context.Context, params CheckVersionParams) (string, error) {
	var versionArg = defaultVersionArg

	binary, err := getBinary(params)
	if err != nil {
		return "", err
	}

	// Run a shell command to get the version string.
	output, err := shell.RunCommandContextAndGetOutputE(t, ctx, &shell.Command{
		Command:    binary,
		Args:       []string{versionArg},
		WorkingDir: params.WorkingDir,
		Env:        map[string]string{},
	})
	if err != nil {
		return "", fmt.Errorf("failed to run shell command for Binary {%s} "+
			"w/ version args {%s}: %w", binary, versionArg, err)
	}

	versionStr, err := ExtractVersionFromShellCommandOutput(output)
	if err != nil {
		return "", fmt.Errorf("failed to extract version from shell "+
			"command output {%s}: %w", output, err)
	}

	return versionStr, nil
}

// getBinary retrieves the binary to use from the given params.
func getBinary(params CheckVersionParams) (string, error) {
	// Use BinaryPath if it is set, otherwise use the binary enum.
	if params.BinaryPath != "" {
		return params.BinaryPath, nil
	}

	switch params.Binary {
	case Docker:
		return "docker", nil
	case Packer:
		return "packer", nil
	case Terraform:
		return terraform.DefaultExecutable, nil
	default:
		return "", &UnsupportedBinaryErr{Binary: params.Binary}
	}
}

// ExtractVersionFromShellCommandOutput extracts version with regex string matching
// from the given shell command output string.
func ExtractVersionFromShellCommandOutput(output string) (string, error) {
	regexMatcher := regexp.MustCompile(versionRegexMatcher)

	versionStr := regexMatcher.FindString(output)
	if versionStr == "" {
		return "", &VersionExtractionErr{Output: output}
	}

	return versionStr, nil
}

// CheckVersionConstraint checks whether the given version passes the version constraint.
//
// It returns [InvalidVersionFormatErr] for ill-formatted version strings and
// [VersionMismatchErr] when the version does not satisfy the constraint.
//
//	CheckVersionConstraint("1.2.31",  ">= 1.2.0, < 2.0.0") - no error
//	CheckVersionConstraint("1.0.31",  ">= 1.2.0, < 2.0.0") - error
func CheckVersionConstraint(actualVersionStr string, versionConstraintStr string) error {
	actualVersion, err := version.NewVersion(actualVersionStr)
	if err != nil {
		return &InvalidVersionFormatErr{Field: "actualVersionStr", Value: actualVersionStr, Underlying: err}
	}

	versionConstraint, err := version.NewConstraint(versionConstraintStr)
	if err != nil {
		return &InvalidVersionFormatErr{Field: "versionConstraint", Value: versionConstraintStr, Underlying: err}
	}

	if !versionConstraint.Check(actualVersion) {
		return &VersionMismatchErr{
			Actual:     actualVersionStr,
			Constraint: versionConstraintStr,
		}
	}

	return nil
}
