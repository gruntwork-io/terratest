//go:build aws
// +build aws

package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAccountId(t *testing.T) {
	accountID := GetAccountId(t)
	assert.Regexp(t, "^[0-9]{12}$", accountID)
}

func TestExtractAccountIdFromValidArn(t *testing.T) {
	t.Parallel()

	expectedAccountID := "123456789012"
	arn := "arn:aws:iam::" + expectedAccountID + ":user/test"

	actualAccountID, err := extractAccountIDFromARN(arn)
	if err != nil {
		t.Fatalf("Unexpected error while extracting account id from arn %s: %s", arn, err)
	}

	if actualAccountID != expectedAccountID {
		t.Fatalf("Did not get expected account id. Expected: %s. Actual: %s.", expectedAccountID, actualAccountID)
	}
}

func TestExtractAccountIdFromInvalidArn(t *testing.T) {
	t.Parallel()

	_, err := extractAccountIDFromARN("invalidArn")
	if err == nil {
		t.Fatalf("Expected an error when extracting an account id from an invalid ARN, but got nil")
	}
}
