package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// minARNParts is the minimum number of colon-separated parts in a valid IAM ARN.
const minARNParts = 5

// GetAccountIDContextE gets the Account ID for the currently logged in IAM User.
// The ctx parameter supports cancellation and timeouts.
func GetAccountIDContextE(t testing.TestingT, ctx context.Context) (string, error) {
	stsClient, err := NewStsClientContextE(t, ctx, defaultRegion)
	if err != nil {
		return "", err
	}

	identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}

	return aws.ToString(identity.Account), nil
}

// GetAccountIDContext gets the Account ID for the currently logged in IAM User.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAccountIDContext(t testing.TestingT, ctx context.Context) string {
	t.Helper()

	id, err := GetAccountIDContextE(t, ctx)
	require.NoError(t, err)

	return id
}

// GetAccountID gets the Account ID for the currently logged in IAM User.
//
// Deprecated: Use [GetAccountIDContext] instead.
func GetAccountID(t testing.TestingT) string {
	t.Helper()

	return GetAccountIDContext(t, context.Background())
}

// GetAccountIDE gets the Account ID for the currently logged in IAM User.
//
// Deprecated: Use [GetAccountIDContextE] instead.
func GetAccountIDE(t testing.TestingT) (string, error) {
	return GetAccountIDContextE(t, context.Background())
}

// GetAccountId gets the Account ID for the currently logged in IAM User.
//
// Deprecated: Use [GetAccountID] instead.
//
//nolint:staticcheck,revive // preserving deprecated function name
func GetAccountId(t testing.TestingT) string {
	return GetAccountID(t)
}

// GetAccountIdE gets the Account ID for the currently logged in IAM User.
//
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

// NewStsClientContextE creates a new STS client.
// The ctx parameter supports cancellation and timeouts.
func NewStsClientContextE(t testing.TestingT, ctx context.Context, region string) (*sts.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return sts.NewFromConfig(*sess), nil
}

// NewStsClientContext creates a new STS client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewStsClientContext(t testing.TestingT, ctx context.Context, region string) *sts.Client {
	t.Helper()

	client, err := NewStsClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// NewStsClient creates a new STS client.
//
// Deprecated: Use [NewStsClientContext] instead.
func NewStsClient(t testing.TestingT, region string) *sts.Client {
	t.Helper()

	return NewStsClientContext(t, context.Background(), region)
}

// NewStsClientE creates a new STS client.
//
// Deprecated: Use [NewStsClientContextE] instead.
func NewStsClientE(t testing.TestingT, region string) (*sts.Client, error) {
	return NewStsClientContextE(t, context.Background(), region)
}
