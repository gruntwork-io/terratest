package collections_test

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/internal/collections"
	"github.com/stretchr/testify/assert"
)

func TestIntersection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		list1 []string
		list2 []string
		want  []string
	}{
		{"common items, ordered by list1", []string{"a", "b", "c"}, []string{"b", "c", "d"}, []string{"b", "c"}},
		{"dedups output", []string{"a", "a"}, []string{"a"}, []string{"a"}},
		{"no overlap returns empty, not nil", []string{"a"}, []string{"b"}, []string{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, collections.Intersection(tc.list1, tc.list2))
		})
	}
}

func TestSubtract(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		list1 []string
		list2 []string
		want  []string
	}{
		{"removes list2 items", []string{"a", "b", "c"}, []string{"b", "c"}, []string{"a"}},
		{"everything removed returns empty, not nil", []string{"a", "b"}, []string{"a", "b"}, []string{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, collections.Subtract(tc.list1, tc.list2))
		})
	}
}

func TestSubtractDoesNotMutateInput(t *testing.T) {
	t.Parallel()

	in := []string{"a", "b"}
	collections.Subtract(in, []string{"a"})
	assert.Equal(t, []string{"a", "b"}, in, "does not mutate the input slice")
}

func ExampleIntersection() {
	fmt.Println(collections.Intersection([]int{1, 2, 3}, []int{2, 3, 4}))
	// Output: [2 3]
}

func ExampleSubtract() {
	fmt.Println(collections.Subtract([]int{1, 2, 3}, []int{2, 3}))
	// Output: [1]
}
