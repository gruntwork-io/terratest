// Package collections allows to interact with lists of things.
//
// Deprecated: The collections package is scheduled for removal in Terratest v2.
// Go's standard library covers these helpers as of Go 1.21+. Replace at the call
// site:
//
//	ListContains(haystack, needle)        -> slices.Contains(haystack, needle)
//	ListIntersection / ListSubtract       -> a short slices.Contains loop (see each function)
//	GetSliceLastValueE / GetSliceIndexValueE -> strings.Split, then index the result
package collections
