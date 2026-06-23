package collections

import (
	"strings"
)

// GetSliceLastValueE will take a source string and returns the last value when split by the separator char.
//
// Deprecated: scheduled for removal in Terratest v2. Use strings.Split at the call
// site, e.g.:
//
//	parts := strings.Split(source, separator)
//	last := parts[len(parts)-1]
//
// Note: strings.Split returns the whole string when the separator is absent, so add a
// strings.Contains check first if you need the not-found error this helper returned.
func GetSliceLastValueE(source string, separator string) (string, error) {
	if len(source) > 0 && len(separator) > 0 && strings.Contains(source, separator) {
		tmp := strings.Split(source, separator)

		return tmp[len(tmp)-1], nil
	}

	return "", NewSliceValueNotFoundError(source)
}

// GetSliceIndexValueE will take a source string and returns the value at the given index when split by
// the separator char.
//
// Deprecated: scheduled for removal in Terratest v2. Use strings.Split at the call
// site, e.g.:
//
//	parts := strings.Split(source, separator)
//	val := parts[index]
//
// Note: bounds-check index first; this helper returned the not-found error for an
// out-of-range or negative index rather than panicking.
func GetSliceIndexValueE(source string, separator string, index int) (string, error) {
	if len(source) > 0 && len(separator) > 0 && strings.Contains(source, separator) && index >= 0 {
		tmp := strings.Split(source, separator)
		if index >= len(tmp) {
			return "", NewSliceValueNotFoundError(source)
		}

		return tmp[index], nil
	}

	return "", NewSliceValueNotFoundError(source)
}
