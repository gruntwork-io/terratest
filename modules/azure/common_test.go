package azure //nolint:testpackage // tests access unexported functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSafePtrToString(t *testing.T) {
	t.Parallel()

	var nilPtr *string

	nilResult := safePtrToString(nilPtr)
	assert.Empty(t, nilResult)

	stringPtr := "Test"
	stringResult := safePtrToString(&stringPtr)
	assert.Equal(t, "Test", stringResult)
}

func TestSafePtrToInt32(t *testing.T) {
	t.Parallel()

	var nilPtr *int32

	nilResult := safePtrToInt32(nilPtr)
	assert.Equal(t, int32(0), nilResult)

	intPtr := int32(42)
	intResult := safePtrToInt32(&intPtr)
	assert.Equal(t, int32(42), intResult)
}
