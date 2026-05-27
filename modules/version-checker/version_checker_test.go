package version_checker_test //nolint:staticcheck // package name matches directory convention

import (
	"errors"
	"testing"

	checker "github.com/gruntwork-io/terratest/modules/version-checker/v2"
	"github.com/stretchr/testify/require"
)

func TestParamValidation(t *testing.T) {
	t.Parallel()

	t.Run("Empty Params", func(t *testing.T) {
		t.Parallel()
		err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{})

		var missingErr *checker.MissingParamErr

		require.ErrorAs(t, err, &missingErr)
	})

	t.Run("Missing VersionConstraint", func(t *testing.T) {
		t.Parallel()
		err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{
			Binary:            checker.Docker,
			VersionConstraint: "",
			WorkingDir:        ".",
		})

		var missingErr *checker.MissingParamErr

		require.ErrorAs(t, err, &missingErr)
	})

	t.Run("Invalid Version Constraint Format", func(t *testing.T) {
		t.Parallel()
		err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{
			Binary:            checker.Docker,
			VersionConstraint: "abc",
			WorkingDir:        ".",
		})

		var constraintErr *checker.InvalidVersionConstraintErr

		require.ErrorAs(t, err, &constraintErr)
	})

	t.Run("Valid Params", func(t *testing.T) {
		t.Parallel()
		err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{
			Binary:            checker.Docker,
			VersionConstraint: ">= 0.0.1",
			WorkingDir:        ".",
		})
		require.NoError(t, err)
	})
}

func TestCheckVersionConstraintMismatch(t *testing.T) {
	t.Parallel()

	err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{
		Binary:            checker.Docker,
		VersionConstraint: ">= 999.999.999",
		WorkingDir:        ".",
	})

	var mismatchErr *checker.VersionMismatchErr

	require.ErrorAs(t, err, &mismatchErr)
	require.Equal(t, ">= 999.999.999", mismatchErr.Constraint)
}

func TestMissingParamErrFields(t *testing.T) {
	t.Parallel()

	err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{})

	var missingErr *checker.MissingParamErr

	require.ErrorAs(t, err, &missingErr)
	require.Equal(t, "WorkingDir", missingErr.Param)
}

func TestInvalidVersionConstraintErrUnwrap(t *testing.T) {
	t.Parallel()

	err := checker.CheckVersionContextE(t, t.Context(), checker.CheckVersionParams{
		Binary:            checker.Docker,
		VersionConstraint: "abc",
		WorkingDir:        ".",
	})

	var constraintErr *checker.InvalidVersionConstraintErr

	require.ErrorAs(t, err, &constraintErr)
	require.Error(t, errors.Unwrap(constraintErr), "underlying parse error should be wrapped")
}

func TestExtractVersionFromShellCommandOutput(t *testing.T) {
	t.Parallel()

	t.Run("Stand-alone version string", func(t *testing.T) {
		t.Parallel()

		versionStr, err := checker.ExtractVersionFromShellCommandOutput("version is 1.2.3")
		require.NoError(t, err)
		require.Equal(t, "1.2.3", versionStr)
	})

	t.Run("version string with v prefix", func(t *testing.T) {
		t.Parallel()

		versionStr, err := checker.ExtractVersionFromShellCommandOutput("version is v1.0.0")
		require.NoError(t, err)
		require.Equal(t, "1.0.0", versionStr)
	})

	t.Run("2 digit version string", func(t *testing.T) {
		t.Parallel()

		versionStr, err := checker.ExtractVersionFromShellCommandOutput("version is v1.0")
		require.NoError(t, err)
		require.Equal(t, "1.0", versionStr)
	})

	t.Run("invalid output string", func(t *testing.T) {
		t.Parallel()

		_, err := checker.ExtractVersionFromShellCommandOutput("version is vabc")

		var extractionErr *checker.VersionExtractionErr

		require.ErrorAs(t, err, &extractionErr)
	})

	t.Run("empty output string", func(t *testing.T) {
		t.Parallel()

		_, err := checker.ExtractVersionFromShellCommandOutput("")

		var extractionErr *checker.VersionExtractionErr

		require.ErrorAs(t, err, &extractionErr)
	})
}

func TestCheckVersionConstraint(t *testing.T) {
	t.Parallel()

	t.Run("invalid actualVersionStr", func(t *testing.T) {
		t.Parallel()

		err := checker.CheckVersionConstraint("", "1.2.3")

		var formatErr *checker.InvalidVersionFormatErr

		require.ErrorAs(t, err, &formatErr)
		require.Equal(t, "actualVersionStr", formatErr.Field)
	})

	t.Run("invalid versionConstraint", func(t *testing.T) {
		t.Parallel()

		err := checker.CheckVersionConstraint("1.2.3", "")

		var formatErr *checker.InvalidVersionFormatErr

		require.ErrorAs(t, err, &formatErr)
		require.Equal(t, "versionConstraint", formatErr.Field)
	})

	t.Run("pass version constraint", func(t *testing.T) {
		t.Parallel()

		err := checker.CheckVersionConstraint("1.2.3", "1.2.3")
		require.NoError(t, err)
	})

	t.Run("fail version constraint", func(t *testing.T) {
		t.Parallel()

		err := checker.CheckVersionConstraint("1.2.3", "1.2.4")

		var mismatchErr *checker.VersionMismatchErr

		require.ErrorAs(t, err, &mismatchErr)
		require.Equal(t, "1.2.3", mismatchErr.Actual)
		require.Equal(t, "1.2.4", mismatchErr.Constraint)
	})

	t.Run("pessimistic constraint operator", func(t *testing.T) {
		t.Parallel()

		err := checker.CheckVersionConstraint("1.0.5", "~> 1.0.4")
		require.NoError(t, err)
	})

	t.Run("compound constraint operators", func(t *testing.T) {
		t.Parallel()

		err := checker.CheckVersionConstraint("1.2.7", ">= 1.2.0, < 2.0.0")
		require.NoError(t, err)
	})
}

// TestCheckVersionEndToEnd assumes Docker, Terraform, and Packer are installed
// and their versions are greater than 0.0.1.
func TestCheckVersionEndToEnd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		param checker.CheckVersionParams
	}{
		{
			name: "Docker",
			param: checker.CheckVersionParams{
				Binary:            checker.Docker,
				VersionConstraint: ">= 0.0.1",
				WorkingDir:        ".",
			},
		},
		{
			name: "Terraform",
			param: checker.CheckVersionParams{
				Binary:            checker.Terraform,
				VersionConstraint: ">= 0.0.1",
				WorkingDir:        ".",
			},
		},
		{
			name: "Packer",
			param: checker.CheckVersionParams{
				BinaryPath:        "/usr/local/bin/packer",
				Binary:            checker.Packer,
				VersionConstraint: ">= 0.0.1",
				WorkingDir:        ".",
			},
		},
	}

	for _, tc := range tests {
		err := checker.CheckVersionContextE(t, t.Context(), tc.param)
		require.NoError(t, err, tc.name)
	}
}
