package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIamCurrentUserName(t *testing.T) {
	t.Parallel()

	username := GetIamCurrentUserName(t)
	assert.NotEmpty(t, username)
}

func TestGetIamCurrentUserArn(t *testing.T) {
	t.Parallel()

	username := GetIamCurrentUserArn(t)
	assert.Regexp(t, "^arn:aws:iam::[0-9]{12}:user/.+$", username)
}
