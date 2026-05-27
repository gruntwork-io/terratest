package collections

import "slices"

// ListIntersection returns all the items in both list1 and list2. Note that this will dedup the items so that the
// output is more predictable. Otherwise, the end list depends on which list was used as the base.
func ListIntersection[T comparable](list1 []T, list2 []T) []T {
	out := []T{}

	// Only need to iterate list1, because we want items in both lists, not union.
	for _, item := range list1 {
		if slices.Contains(list2, item) && !slices.Contains(out, item) {
			out = append(out, item)
		}
	}

	return out
}

// ListSubtract removes all the items in list2 from list1.
func ListSubtract[T comparable](list1 []T, list2 []T) []T {
	out := []T{}

	for _, item := range list1 {
		if !slices.Contains(list2, item) {
			out = append(out, item)
		}
	}

	return out
}

// ListContains returns true if the given list of strings (haystack) contains the given string (needle).
//
// Deprecated: Use slices.Contains instead.
func ListContains(haystack []string, needle string) bool {
	return slices.Contains(haystack, needle)
}
