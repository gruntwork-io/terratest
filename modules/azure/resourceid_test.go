package azure_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
)

func TestGetNameFromResourceID(t *testing.T) {
	t.Parallel()

	resultSuccess := azure.GetNameFromResourceID("this/is/a/long/slash/separated/string/ResourceID")
	assert.Equal(t, "ResourceID", resultSuccess)

	resultBadSeparator := azure.GetNameFromResourceID("noresourcepresent")
	assert.Empty(t, resultBadSeparator)
}
