// Package random contains different random generators.
package random

import (
	"bytes"
	"math/rand"
	"time"
)

// Random generates a random int between min and max, inclusive.
func Random(min int, max int) int {
	return newRand().Intn(max-min+1) + min
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

// UniqueId returns a unique (ish) id we can attach to resources and tfstate files so they don't conflict with each other
// Uses base 62 to generate a 6 character string that's unlikely to collide with the handful of tests we run in
// parallel. Based on code here: http://stackoverflow.com/a/9543797/483528
func UniqueId() string {
	var out bytes.Buffer

	generator := newRand()
	for i := 0; i < uniqueIDLength; i++ {
		out.WriteByte(base62chars[generator.Intn(len(base62chars))])
	}

	return out.String()
}

const base36chars = "0123456789abcdefghijklmnopqrstuvwxyz"
const uniqueLowercaseIDLength = 6 // Should be good for 36^6 = 2+ billion combinations

// UniqueLowercaseId returns a unique (ish) id we can attach to resources and tfstate files so they don't conflict with each other
// Uses base 36 to generate a 6 character string that's unlikely to collide with the handful of tests we run in
// parallel. Based on code here: http://stackoverflow.com/a/9543797/483528
func UniqueLowercaseId() string {
	var out bytes.Buffer

	generator := newRand()
	for i := 0; i < uniqueLowercaseIDLength; i++ {
		out.WriteByte(base36chars[generator.Intn(len(base36chars))])
	}

	return out.String()
}

// newRand creates a new random number generator, seeding it with the current system time.
func newRand() *rand.Rand {
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
