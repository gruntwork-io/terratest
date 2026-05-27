package collections

import (
	"strings"
)

// GetSliceLastValueE will take a source string and returns the last value when split by the separator char.
func GetSliceLastValueE(source string, separator string) (string, error) {
	if len(source) > 0 && len(separator) > 0 && strings.Contains(source, separator) {
		tmp := strings.Split(source, separator)

		return tmp[len(tmp)-1], nil
	}

	return "", NewSliceValueNotFoundError(source)
}

// GetSliceIndexValueE will take a source string and returns the value at the given index when split by
// the separator char.
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
