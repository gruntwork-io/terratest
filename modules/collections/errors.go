package collections

// SliceValueNotFoundError is returned when a provided values file input is not found on the host path.
//
// Deprecated: scheduled for removal in Terratest v2 along with the collections package.
type SliceValueNotFoundError struct {
	sourceString string
}

func (err SliceValueNotFoundError) Error() string {
	return "Could not resolve requested slice value from string " + err.sourceString
}

// NewSliceValueNotFoundError creates a new slice found error
//
// Deprecated: scheduled for removal in Terratest v2 along with the collections package.
func NewSliceValueNotFoundError(sourceString string) SliceValueNotFoundError {
	return SliceValueNotFoundError{sourceString}
}
