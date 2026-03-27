package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// minARNParts is the minimum number of colon-separated parts in a valid IAM ARN.
const minARNParts = 5

// GetAccountID gets the Account ID for the currently logged in IAM User.
func GetAccountID(t testing.TestingT) string {
	id, err := GetAccountIDE(t)
	if err != nil {
		t.Fatal(err)
	}

	return id
}

// GetAccountIDE gets the Account ID for the currently logged in IAM User.
func GetAccountIDE(t testing.TestingT) (string, error) {
	stsClient, err := NewStsClientE(t, defaultRegion)
	if err != nil {
		return "", err
	}

	identity, err := stsClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	return aws.ToString(identity.Account), nil
}

// Deprecated: Use [GetAccountID] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetAccountId(t testing.TestingT) string {
	return GetAccountID(t)
}

// Deprecated: Use [GetAccountIDE] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetAccountIdE(t testing.TestingT) (string, error) {
	return GetAccountIDE(t)
}

// ExtractAccountIDFromARN extracts the AWS account ID from an IAM ARN.
// An IAM ARN is of the format arn:aws:iam::123456789012:user/test. The account ID is the number after arn:aws:iam::,
// so we split on a colon and return the 5th item.
func ExtractAccountIDFromARN(arn string) (string, error) {
	arnParts := strings.Split(arn, ":")

	if len(arnParts) < minARNParts {
		return "", errors.New("Unrecognized format for IAM ARN: " + arn)
	}

	return arnParts[4], nil
}

// NewStsClientE creates a new STS client.
func NewStsClientE(t testing.TestingT, region string) (*sts.Client, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return sts.NewFromConfig(*sess), nil
}
