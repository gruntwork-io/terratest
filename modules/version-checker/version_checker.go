package version_checker

import (
	"fmt"
	"regexp"

	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/gruntwork-io/terratest/modules/shell"
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
	// defaultVersionRegexMatcher is a regex used to extract version string from shell command output.
	defaultVersionRegexMatcher = `\d+(\.\d+)+`
	// defaultVersionArg is a default arg to pass in to get version output from shell command.
	defaultVersionArg = "--version"
)

type CheckVersionParams struct {
	// BinaryPath is a path to the binary you want to check the version for.
	BinaryPath string
	// Binary is the name of the binary you want to check the version for.
	Binary VersionCheckerBinary
	// VersionArg is a string literal to pass in to get version output from shell command.
	VersionArg string
	// VersionConstraint is a string literal containing one or more conditions, which are separated by commas.
	// More information here:https://www.terraform.io/language/expressions/version-constraints
	VersionConstraint string
	// VersionRegexMatcher is a regex used to extract version string from shell command output.
	VersionRegexMatcher string
	// WorkingDir is a directory you want to run the shell command.
	WorkingDir string
}

// CheckVersionE checks whether the given Binary version is greater than or equal
// to the given expected version.
func CheckVersionE(
	t testing.TestingT,
	params CheckVersionParams) error {

	if err := validateParams(params); err != nil {
		return err
	}

	binaryVersion, err := getVersionWithShellCommand(t, params)
	if err != nil {
		return err
	}

	return checkVersionConstraint(binaryVersion, params.VersionConstraint)
}

// CheckVersion checks whether the given Binary version is greater than or equal to the
// given expected version and fails if it's not.
func CheckVersion(
	t testing.TestingT,
	params CheckVersionParams) {
	require.NoError(t, CheckVersionE(t, params))
}

// Validate whether the given params contains valid data to check version.
func validateParams(params CheckVersionParams) error {
	// Check for empty parameters
	if params.WorkingDir == "" {
		return fmt.Errorf("set WorkingDir in params")
	} else if params.VersionConstraint == "" {
		return fmt.Errorf("set VersionConstraint in params")
	}

	// Check the format of the version constraint if present.
	if _, err := version.NewConstraint(params.VersionConstraint); params.VersionConstraint != "" && err != nil {
		return fmt.Errorf(
			"invalid version constraint format found {%s}", params.VersionConstraint)
	}

	// Check the format of the version regex matcher if present.
	if _, err := regexp.Compile(params.VersionRegexMatcher); params.VersionRegexMatcher != "" && err != nil {
		return fmt.Errorf(
			"invalid version regex matcher format found {%s}", params.VersionRegexMatcher)
	}

	return nil
}

// getVersionWithShellCommand get version by running a shell command.
func getVersionWithShellCommand(t testing.TestingT, params CheckVersionParams) (string, error) {
	var versionArg string
	if params.VersionArg != "" {
		versionArg = params.VersionArg
	} else {
		versionArg = defaultVersionArg
	}

	binary, err := getBinary(params)
	if err != nil {
		return "", err
	}

	// Run a shell command to get the version string.
	output, err := shell.RunCommandAndGetOutputE(t, shell.Command{
		Command:    binary,
		Args:       []string{versionArg},
		WorkingDir: params.WorkingDir,
		Env:        map[string]string{},
	})
	if err != nil {
		return "", fmt.Errorf("failed to run shell command for Binary {%s} "+
			"w/ version args {%s}: %w", binary, versionArg, err)
	}

	versionStr, err := extractVersionFromShellCommandOutput(params, output)
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
		return "", fmt.Errorf("unsupported Binary for checking versions {%d}", params.Binary)
	}
}

// extractVersionFromShellCommandOutput extracts version with regex string matching
// from the given shell command output string.
func extractVersionFromShellCommandOutput(params CheckVersionParams, output string) (string, error) {
	var regexMatcher *regexp.Regexp
	if params.VersionRegexMatcher != "" {
		regexMatcher = regexp.MustCompile(params.VersionRegexMatcher)
	} else {
		regexMatcher = regexp.MustCompile(defaultVersionRegexMatcher)
	}

	versionStr := regexMatcher.FindString(output)
	if versionStr == "" {
		return "", fmt.Errorf("failed to find version using regex matcher")
	}

	return versionStr, nil
}

// checkVersionConstraint checks whether the given version pass the version constraint.
//
// It returns Error for ill-formatted version string and VersionMismatchErr for
// minimum version check failure.
//
//	checkVersionConstraint(t, "1.2.31",  ">= 1.2.0, < 2.0.0") - no error
//	checkVersionConstraint(t, "1.0.31",  ">= 1.2.0, < 2.0.0") - error
func checkVersionConstraint(actualVersionStr string, versionConstraintStr string) error {
	actualVersion, err := version.NewVersion(actualVersionStr)
	if err != nil {
		return fmt.Errorf("invalid version format found for actualVersionStr: %s", actualVersionStr)
	}

	versionConstraint, err := version.NewConstraint(versionConstraintStr)
	if err != nil {
		return fmt.Errorf("invalid version format found for versionConstraint: %s", versionConstraintStr)
	}

	if !versionConstraint.Check(actualVersion) {
		return &VersionMismatchErr{
			errorMessage: fmt.Sprintf("actual version {%s} failed "+
				"the version constraint {%s}", actualVersionStr, versionConstraint),
		}

	}

	return nil
}
