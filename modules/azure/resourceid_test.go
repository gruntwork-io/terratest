package azure_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetNameFromResourceID(t *testing.T) {
	t.Parallel()

	resultSuccess := azure.GetNameFromResourceID("this/is/a/long/slash/separated/string/ResourceID")
	assert.Equal(t, "ResourceID", resultSuccess)

	resultBadSeparator := azure.GetNameFromResourceID("noresourcepresent")
	assert.Empty(t, resultBadSeparator)
}

func TestGetNameFromResourceIDE(t *testing.T) {
	t.Parallel()

	name, err := azure.GetNameFromResourceIDE("this/is/a/long/slash/separated/string/ResourceID")
	require.NoError(t, err)
	assert.Equal(t, "ResourceID", name)

	_, err = azure.GetNameFromResourceIDE("noresourcepresent")
	require.Error(t, err)
}
