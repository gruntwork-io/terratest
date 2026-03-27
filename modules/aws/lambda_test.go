package aws_test

import (
	"testing"

	aws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
)

func TestFunctionError(t *testing.T) {
	t.Parallel()

	// assert that the error message contains all the components of the error, in a readable form
	err := &aws.FunctionError{Message: "message", StatusCode: 123, Payload: []byte("payload")}
	require.Contains(t, err.Error(), "message")
	require.Contains(t, err.Error(), "123")
	require.Contains(t, err.Error(), "payload")
}
