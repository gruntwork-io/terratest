package collections_test

import (
	"slices"
	"testing"

	"github.com/gruntwork-io/terratest/modules/core/v2/collections"
	"github.com/stretchr/testify/assert"
)

func TestListContains(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		element     string
		list        []string
		expected    bool
	}{
		{"empty list, empty element", "", []string{}, false},
		{"empty list, non-empty element", "foo", []string{}, false},
		{"list with 1 item, element matches", "foo", []string{"foo"}, true},
		{"list with 1 item, element doesn't match", "foo", []string{"bar"}, false},
		{"list with 3 items, element matches", "foo", []string{"bar", "foo", "baz"}, true},
		{"list with 3 items, element doesn't match", "nope", []string{"bar", "foo", "baz"}, false},
		{"list with 3 items, empty element", "", []string{"bar", "foo", "baz"}, false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actual := slices.Contains(testCase.list, testCase.element)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestSubtract(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		list1       []string
		list2       []string
		expected    []string
	}{
		{"empty list, empty list", []string{}, []string{}, []string{}},
		{"empty list, non-empty list", []string{}, []string{"foo"}, []string{}},
		{"non-empty list, empty list", []string{"foo"}, []string{}, []string{"foo"}},
		{"list with 1 item, list with no matches", []string{"foo"}, []string{"bar"}, []string{"foo"}},
		{"list with 1 item, list with 1 match", []string{"foo"}, []string{"foo"}, []string{}},
		{"list with 1 item, list with multiple matches and non-matches", []string{"foo"}, []string{"foo", "bar", "foo"}, []string{}},
		{"list with multiple items, list with no matches", []string{"foo", "bar", "baz"}, []string{"abc", "def"}, []string{"foo", "bar", "baz"}},
		{"list with multiple items, list with 1 match", []string{"foo", "bar", "baz"}, []string{"abc", "foo", "def"}, []string{"bar", "baz"}},
		{"list with multiple items, list with multiple matches", []string{"foo", "bar", "baz", "foo", "bar", "baz"}, []string{"abc", "foo", "baz"}, []string{"bar", "bar"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actual := collections.ListSubtract(testCase.list1, testCase.list2)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}

func TestIntersection(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		description string
		list1       []string
		list2       []string
		expected    []string
	}{
		{"empty list, empty list", []string{}, []string{}, []string{}},
		{"empty list, non-empty list", []string{}, []string{"foo"}, []string{}},
		{"non-empty list, empty list", []string{"foo"}, []string{}, []string{}},
		{"list with 1 item, list with no matches", []string{"foo"}, []string{"bar"}, []string{}},
		{"list with 1 item, list with 1 match", []string{"foo"}, []string{"foo"}, []string{"foo"}},
		{"list with 1 item, list with multiple matches and non-matches", []string{"foo"}, []string{"foo", "bar", "foo"}, []string{"foo"}},
		{"list with multiple items, list with no matches", []string{"foo", "bar", "baz"}, []string{"abc", "def"}, []string{}},
		{"list with multiple items, list with 1 match", []string{"foo", "bar", "baz"}, []string{"abc", "foo", "def"}, []string{"foo"}},
		{"list with multiple items, list with multiple matches", []string{"foo", "bar", "baz", "foo", "bar", "baz"}, []string{"abc", "foo", "baz"}, []string{"foo", "baz"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			t.Parallel()

			actual := collections.ListIntersection(testCase.list1, testCase.list2)
			assert.Equal(t, testCase.expected, actual)
		})
	}
}
