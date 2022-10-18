package version_checker

// VersionMismatchErr is an error to indicate version mismatch.
type VersionMismatchErr struct {
	errorMessage string
}

func (r *VersionMismatchErr) Error() string {
	return r.errorMessage
}
