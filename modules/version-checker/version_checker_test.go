package version_checker

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParamValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		param                CheckVersionParams
		containError         bool
		expectedErrorMessage string
	}{
		{
			name:                 "Empty Params",
			param:                CheckVersionParams{},
			containError:         true,
			expectedErrorMessage: "set WorkingDir in params",
		},
		{
			name: "Missing VersionConstraint",
			param: CheckVersionParams{
				Binary:            Docker,
				VersionConstraint: "",
				WorkingDir:        ".",
			},
			containError:         true,
			expectedErrorMessage: "set VersionConstraint in params",
		},
		{
			name: "Invalid Version Constraint Format",
			param: CheckVersionParams{
				Binary:            Docker,
				VersionConstraint: "abc",
				WorkingDir:        ".",
			},
			containError:         true,
			expectedErrorMessage: "invalid version constraint format found {abc}",
		},
		{
			name: "Invalid Version Regex Matcher Format",
			param: CheckVersionParams{
				Binary:              Docker,
				VersionConstraint:   "1.2.3",
				VersionRegexMatcher: "[notregex",
				WorkingDir:          ".",
			},
			containError:         true,
			expectedErrorMessage: "invalid version regex matcher format found {[notregex}",
		},
		{
			name: "Success",
			param: CheckVersionParams{
				Binary:            Docker,
				VersionConstraint: ">1.2.3",
				WorkingDir:        ".",
			},
			containError:         false,
			expectedErrorMessage: "",
		},
	}

	for _, tc := range tests {
		err := validateParams(tc.param)
		if tc.containError {
			require.EqualError(t, err, tc.expectedErrorMessage, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

func TestExtractVersionFromShellCommandOutput(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		outputStr            string
		expectedVersionStr   string
		containError         bool
		expectedErrorMessage string
		param                CheckVersionParams
	}{
		{
			name:                 "Stand-alone version string",
			outputStr:            "version is 1.2.3",
			expectedVersionStr:   "1.2.3",
			containError:         false,
			expectedErrorMessage: "",
			param:                CheckVersionParams{},
		},
		{
			name:                 "version string with v prefix",
			outputStr:            "version is v1.0.0",
			expectedVersionStr:   "1.0.0",
			containError:         false,
			expectedErrorMessage: "",
			param:                CheckVersionParams{},
		},
		{
			name:                 "2 digit version string",
			outputStr:            "version is v1.0",
			expectedVersionStr:   "1.0",
			containError:         false,
			expectedErrorMessage: "",
			param:                CheckVersionParams{},
		},
		{
			name:                 "invalid output string",
			outputStr:            "version is vabc",
			expectedVersionStr:   "",
			containError:         true,
			expectedErrorMessage: "failed to find version using regex matcher",
			param:                CheckVersionParams{},
		},
		{
			name:                 "empty output string",
			outputStr:            "",
			expectedVersionStr:   "",
			containError:         true,
			expectedErrorMessage: "failed to find version using regex matcher",
			param:                CheckVersionParams{},
		},
	}

	for _, tc := range tests {
		versionStr, err := extractVersionFromShellCommandOutput(tc.param, tc.outputStr)
		if tc.containError {
			require.EqualError(t, err, tc.expectedErrorMessage, tc.name)
		} else {
			require.NoError(t, err, tc.name)
			require.Equal(t, tc.expectedVersionStr, versionStr, tc.name)
		}
	}
}

func TestCheckVersionConstraint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		actualVersionStr     string
		versionConstraint    string
		containError         bool
		expectedErrorMessage string
	}{
		{
			name:                 "invalid actualVersionStr",
			actualVersionStr:     "",
			versionConstraint:    "1.2.3",
			containError:         true,
			expectedErrorMessage: "invalid version format found for actualVersionStr: ",
		},
		{
			name:                 "invalid versionConstraint",
			actualVersionStr:     "1.2.3",
			versionConstraint:    "",
			containError:         true,
			expectedErrorMessage: "invalid version format found for versionConstraint: ",
		},
		{
			name:                 "pass version constraint",
			actualVersionStr:     "1.2.3",
			versionConstraint:    "1.2.3",
			containError:         false,
			expectedErrorMessage: "",
		},
		{
			name:                 "fail version constraint",
			actualVersionStr:     "1.2.3",
			versionConstraint:    "1.2.4",
			containError:         true,
			expectedErrorMessage: "actual version {1.2.3} failed the version constraint {1.2.4}",
		},
		{
			name:                 "special syntax version constraint",
			actualVersionStr:     "1.0.5",
			versionConstraint:    "~> 1.0.4",
			containError:         false,
			expectedErrorMessage: "",
		},
		{
			name:                 "version constraint w/ operators",
			actualVersionStr:     "1.2.7",
			versionConstraint:    ">= 1.2.0, < 2.0.0",
			containError:         false,
			expectedErrorMessage: ""},
	}

	for _, tc := range tests {
		err := checkVersionConstraint(tc.actualVersionStr, tc.versionConstraint)
		if tc.containError {
			require.EqualError(t, err, tc.expectedErrorMessage, tc.name)
		} else {
			require.NoError(t, err, tc.name)
		}
	}
}

// Note: with the current implementation of running shell command, it's not easy to
// mock the output of running a shell command. So we assume a certain Binary is installed in the working
// directory and it's greater than 0.0.1 version.
func TestCheckVersionEndToEnd(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		param CheckVersionParams
	}{
		{name: "Docker", param: CheckVersionParams{
			Binary:            Docker,
			VersionConstraint: ">= 0.0.1",
			WorkingDir:        ".",
		}},
		{name: "Terraform", param: CheckVersionParams{
			BinaryPath:        "",
			Binary:            Terraform,
			VersionConstraint: ">= 0.0.1",
			WorkingDir:        ".",
		}},
		{name: "Packer", param: CheckVersionParams{
			BinaryPath:        "/usr/local/bin/packer",
			Binary:            Packer,
			VersionConstraint: ">= 0.0.1",
			WorkingDir:        ".",
		}},
		{name: "Go", param: CheckVersionParams{
			BinaryPath:        "go",
			VersionArg:        "version",
			VersionConstraint: ">= 1.17.0",
			WorkingDir:        ".",
		}},
	}

	for _, tc := range tests {
		err := CheckVersionE(t, tc.param)
		require.NoError(t, err, tc.name)
	}
}
