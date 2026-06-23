package collections

import "slices"

// ListIntersection returns all the items in both list1 and list2. Note that this will dedup the items so that the
// output is more predictable. Otherwise, the end list depends on which list was used as the base.
//
// Deprecated: scheduled for removal in Terratest v2. The collections package is being
// dropped, so there is no drop-in public replacement; build it inline with the slices
// package (Go 1.21+) at the call site, e.g.:
//
//	out := []T{}
//	for _, x := range list1 {
//		if slices.Contains(list2, x) && !slices.Contains(out, x) {
//			out = append(out, x)
//		}
//	}
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
//
// Deprecated: scheduled for removal in Terratest v2. The collections package is being
// dropped, so there is no drop-in public replacement; build it inline with the slices
// package (Go 1.21+) at the call site, e.g.:
//
//	out := []T{}
//	for _, x := range list1 {
//		if !slices.Contains(list2, x) {
//			out = append(out, x)
//		}
//	}
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
// Deprecated: scheduled for removal in Terratest v2. Replace at the call site with
// slices.Contains(haystack, needle) (Go 1.21+).
func ListContains(haystack []string, needle string) bool {
	return slices.Contains(haystack, needle)
}
