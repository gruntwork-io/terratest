// Package random contains different random generators.
package random

import (
	"bytes"
	"math/rand"
)

// Random generates a random int between min and max, inclusive.
func Random(min int, max int) int {
	return rand.Intn(max-min+1) + min
}

// RandomInt picks a random element in the slice of ints.
func RandomInt(elements []int) int {
	index := Random(0, len(elements)-1)
	return elements[index]
}

// RandomString picks a random element in the slice of string.
func RandomString(elements []string) string {
	index := Random(0, len(elements)-1)
	return elements[index]
}

const base62chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const uniqueIDLength = 6 // Should be good for 62^6 = 56+ billion combinations

// UniqueID returns a unique (ish) id we can attach to resources and tfstate files so they don't conflict with each other
// Uses base 62 to generate a 6 character string that's unlikely to collide with the handful of tests we run in
// parallel. Based on code here: http://stackoverflow.com/a/9543797/483528
func UniqueID() string {
	var out bytes.Buffer

	for i := 0; i < uniqueIDLength; i++ {
		out.WriteByte(base62chars[rand.Intn(len(base62chars))])
	}

	return out.String()
}
