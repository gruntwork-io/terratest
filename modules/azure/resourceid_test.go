package azure_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNameFromResourceID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		resourceID string
		want       string
	}{
		{"normal resource ID", "this/is/a/long/slash/separated/string/ResourceID", "ResourceID"},
		{"no separator", "noresourcepresent", ""},
		{"trailing slash", "this/is/a/ResourceID/", ""},
		{"empty", "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, azure.GetNameFromResourceID(tc.resourceID))
		})
	}
}

func TestGetNameFromResourceIDE(t *testing.T) {
	t.Parallel()

	name, err := azure.GetNameFromResourceIDE("this/is/a/long/slash/separated/string/ResourceID")
	require.NoError(t, err)
	assert.Equal(t, "ResourceID", name)

	tests := []struct {
		name       string
		resourceID string
	}{
		{"no separator", "noresourcepresent"},
		{"trailing slash", "this/is/a/ResourceID/"},
		{"empty", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := azure.GetNameFromResourceIDE(tc.resourceID)
			require.Error(t, err)

			var notFound azure.ResourceIDNameNotFoundError
			require.ErrorAs(t, err, &notFound)
			assert.Equal(t, tc.resourceID, notFound.ResourceID)
		})
	}
}
