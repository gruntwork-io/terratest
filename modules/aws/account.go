package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"

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

	identity, err := stsClient.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	return aws.ToString(identity.Account), nil
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
func NewStsClientE(t testing.TestingT, region string) (*sts.Client, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return sts.NewFromConfig(*sess), nil
}
