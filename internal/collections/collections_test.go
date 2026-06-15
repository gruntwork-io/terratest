package collections_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/internal/collections"
	"github.com/stretchr/testify/assert"
)

func TestIntersection(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []string{"b", "c"}, collections.Intersection([]string{"a", "b", "c"}, []string{"b", "c", "d"}))
	assert.Equal(t, []string{"a"}, collections.Intersection([]string{"a", "a"}, []string{"a"}), "dedups output")
	assert.Equal(t, []string{}, collections.Intersection([]string{"a"}, []string{"b"}), "returns empty, not nil")
}

func TestSubtract(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []string{"a"}, collections.Subtract([]string{"a", "b", "c"}, []string{"b", "c"}))
	assert.Equal(t, []string{}, collections.Subtract([]string{"a", "b"}, []string{"a", "b"}), "returns empty, not nil")

	in := []string{"a", "b"}
	collections.Subtract(in, []string{"a"})
	assert.Equal(t, []string{"a", "b"}, in, "does not mutate the input slice")
}
