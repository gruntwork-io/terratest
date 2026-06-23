// Package collections holds small generic slice helpers shared across
// Terratest's own packages. It is internal and not part of the public API; the
// public modules/collections package is deprecated and scheduled for removal in
// v2.
package collections

// Intersection returns the items present in both lists, de-duplicated, in the
// order they appear in list1.
func Intersection[T comparable](list1 []T, list2 []T) []T {
	lookups := make(map[T]struct{}, len(list2))
	for _, item := range list2 {
		lookups[item] = struct{}{}
	}

	out := make([]T, 0, min(len(list1), len(list2)))

	for _, item := range list1 {
		if _, found := lookups[item]; found {
			out = append(out, item)
			delete(lookups, item) // delete so a repeated list1 item isn't emitted twice
		}
	}

	return out
}

// Subtract returns the items in list1 that are not in list2.
func Subtract[T comparable](list1 []T, list2 []T) []T {
	lookups := make(map[T]struct{}, len(list2))
	for _, item := range list2 {
		lookups[item] = struct{}{}
	}

	out := make([]T, 0, len(list1))

	for _, item := range list1 {
		if _, found := lookups[item]; !found {
			out = append(out, item)
		}
	}

	return out
}
