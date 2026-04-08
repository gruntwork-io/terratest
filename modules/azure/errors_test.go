package azure //nolint:testpackage // tests access unexported functions

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/assert"
)

func TestResourceNotFoundErrorExists(t *testing.T) {
	t.Parallel()

	tests := []struct { //nolint:govet // fieldalignment not worth optimizing in test structs
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"non-azcore error", errors.New("something failed"), false},
		{"ResourceNotFound", &azcore.ResponseError{ErrorCode: "ResourceNotFound", StatusCode: http.StatusNotFound}, true},
		{"ResourceGroupNotFound", &azcore.ResponseError{ErrorCode: "ResourceGroupNotFound", StatusCode: http.StatusNotFound}, true},
		{"AuthorizationFailed", &azcore.ResponseError{ErrorCode: "AuthorizationFailed", StatusCode: http.StatusForbidden}, false},
		{"wrapped ResourceNotFound", fmt.Errorf("outer: %w", &azcore.ResponseError{ErrorCode: "ResourceNotFound", StatusCode: http.StatusNotFound}), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.expected, ResourceNotFoundErrorExists(tt.err))
		})
	}
}
