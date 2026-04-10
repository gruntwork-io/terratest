//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure_test

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
)

func TestGetNameFromResourceID(t *testing.T) {
	t.Parallel()

	// set slice variables
	sliceSource := "this/is/a/long/slash/separated/string/ResourceID"
	sliceResult := "ResourceID"
	sliceNotFound := "noresourcepresent"

	// verify success
	resultSuccess := azure.GetNameFromResourceID(sliceSource)
	assert.Equal(t, sliceResult, resultSuccess)

	// verify error when separator not found
	resultBadSeparator := azure.GetNameFromResourceID(sliceNotFound)
	assert.Empty(t, resultBadSeparator)
}
