package aws

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// GetAccountId gets the Account ID for the currently logged in IAM User.
func GetAccountId(t testing.TestingT) string {
	id, err := GetAccountIdE(t)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

// GetAccountIdE gets the Account ID for the currently logged in IAM User.
func GetAccountIdE(t testing.TestingT) (string, error) {
	stsClient, err := NewStsClientE(t, defaultRegion)
	if err != nil {
		return "", err
	}

	identity, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	return aws.StringValue(identity.Account), nil
}

// An IAM arn is of the format arn:aws:iam::123456789012:user/test. The account id is the number after arn:aws:iam::,
// so we split on a colon and return the 5th item.
func extractAccountIDFromARN(arn string) (string, error) {
	arnParts := strings.Split(arn, ":")

	if len(arnParts) < 5 {
		return "", errors.New("Unrecognized format for IAM ARN: " + arn)
	}

	return arnParts[4], nil
}

// NewStsClientE creates a new STS client.
func NewStsClientE(t testing.TestingT, region string) (*sts.STS, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return sts.New(sess), nil
}
