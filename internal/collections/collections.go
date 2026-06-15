// Package collections holds small generic slice helpers shared across
// Terratest's own packages. It is internal and not part of the public API; the
// public modules/collections package is deprecated and scheduled for removal in
// v2.
package collections

import "slices"

// Intersection returns the items present in both lists, de-duplicated, in the
// order they appear in list1.
func Intersection[T comparable](list1 []T, list2 []T) []T {
	out := []T{}

	for _, item := range list1 {
		if slices.Contains(list2, item) && !slices.Contains(out, item) {
			out = append(out, item)
		}
	}

	return out
}

// Subtract returns the items in list1 that are not in list2.
func Subtract[T comparable](list1 []T, list2 []T) []T {
	out := []T{}

	for _, item := range list1 {
		if !slices.Contains(list2, item) {
			out = append(out, item)
		}
	}

	return out
}
