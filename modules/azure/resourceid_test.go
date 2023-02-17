//go:build azure
// +build azure

// NOTE: We use build tags to differentiate azure testing because we currently do not have azure access setup for
// CircleCI.

package azure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetNameFromResourceID(t *testing.T) {
	t.Parallel()

	// set slice variables
	sliceSource := "this/is/a/long/slash/separated/string/ResourceID"
	sliceResult := "ResourceID"
	sliceNotFound := "noresourcepresent"

	// verify success
	resultSuccess := GetNameFromResourceID(sliceSource)
	assert.Equal(t, sliceResult, resultSuccess)

	// verify error when seperator not found
	resultBadSeperator := GetNameFromResourceID(sliceNotFound)
	assert.Equal(t, "", resultBadSeperator)
}
