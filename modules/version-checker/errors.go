package version_checker //nolint:staticcheck // package name matches directory convention

import "fmt"

// VersionMismatchErr is returned when a binary's version does not satisfy the
// given version constraint.
//
// Deprecated: scheduled for removal in Terratest v2 along with the version_checker package.
type VersionMismatchErr struct {
	Actual     string
	Constraint string
}

func (e *VersionMismatchErr) Error() string {
	return fmt.Sprintf("actual version {%s} failed the version constraint {%s}", e.Actual, e.Constraint)
}

// MissingParamErr is returned when a required field in [CheckVersionParams] is
// empty.
//
// Deprecated: scheduled for removal in Terratest v2 along with the version_checker package.
type MissingParamErr struct {
	Param string
}

func (e *MissingParamErr) Error() string {
	return fmt.Sprintf("set %s in params", e.Param)
}

// InvalidVersionConstraintErr is returned when the VersionConstraint string
// cannot be parsed.
//
// Deprecated: scheduled for removal in Terratest v2 along with the version_checker package.
type InvalidVersionConstraintErr struct {
	Underlying error
	Constraint string
}

func (e *InvalidVersionConstraintErr) Error() string {
	return fmt.Sprintf("invalid version constraint format found {%s}: %s", e.Constraint, e.Underlying)
}

func (e *InvalidVersionConstraintErr) Unwrap() error {
	return e.Underlying
}

// UnsupportedBinaryErr is returned when [CheckVersionParams.Binary] is set to
// an unknown [VersionCheckerBinary] value.
//
// Deprecated: scheduled for removal in Terratest v2 along with the version_checker package.
type UnsupportedBinaryErr struct {
	Binary VersionCheckerBinary
}

func (e *UnsupportedBinaryErr) Error() string {
	return fmt.Sprintf("unsupported Binary for checking versions {%d}", e.Binary)
}

// InvalidVersionFormatErr is returned when a version string cannot be parsed.
//
// Deprecated: scheduled for removal in Terratest v2 along with the version_checker package.
type InvalidVersionFormatErr struct {
	Underlying error
	Field      string
	Value      string
}

func (e *InvalidVersionFormatErr) Error() string {
	return fmt.Sprintf("invalid version format found for %s %s: %s", e.Field, e.Value, e.Underlying)
}

func (e *InvalidVersionFormatErr) Unwrap() error {
	return e.Underlying
}

// VersionExtractionErr is returned when the version string cannot be parsed
// from the output of running a binary's --version command.
//
// Deprecated: scheduled for removal in Terratest v2 along with the version_checker package.
type VersionExtractionErr struct {
	Output string
}

func (e *VersionExtractionErr) Error() string {
	return fmt.Sprintf("failed to find version using regex matcher in output {%s}", e.Output)
}
