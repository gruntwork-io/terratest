package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// GetCmkArnContextE gets the ARN of a KMS Customer Master Key (CMK) in the given region with the given ID. The ID can be an alias, such
// as "alias/my-cmk".
// The ctx parameter supports cancellation and timeouts.
func GetCmkArnContextE(t testing.TestingT, ctx context.Context, region string, cmkID string) (string, error) {
	kmsClient, err := NewKmsClientContextE(t, ctx, region)
	if err != nil {
		return "", err
	}

	result, err := kmsClient.DescribeKey(ctx, &kms.DescribeKeyInput{
		KeyId: aws.String(cmkID),
	})
	if err != nil {
		return "", err
	}

	return *result.KeyMetadata.Arn, nil
}

// GetCmkArnContext gets the ARN of a KMS Customer Master Key (CMK) in the given region with the given ID. The ID can be an alias, such
// as "alias/my-cmk".
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetCmkArnContext(t testing.TestingT, ctx context.Context, region string, cmkID string) string {
	t.Helper()

	out, err := GetCmkArnContextE(t, ctx, region, cmkID)
	if err != nil {
		t.Fatal(err)
	}

	return out
}

// GetCmkArn gets the ARN of a KMS Customer Master Key (CMK) in the given region with the given ID. The ID can be an alias, such
// as "alias/my-cmk".
//
// Deprecated: Use [GetCmkArnContext] instead.
func GetCmkArn(t testing.TestingT, region string, cmkID string) string {
	t.Helper()

	return GetCmkArnContext(t, context.Background(), region, cmkID)
}

// GetCmkArnE gets the ARN of a KMS Customer Master Key (CMK) in the given region with the given ID. The ID can be an alias, such
// as "alias/my-cmk".
//
// Deprecated: Use [GetCmkArnContextE] instead.
func GetCmkArnE(t testing.TestingT, region string, cmkID string) (string, error) {
	return GetCmkArnContextE(t, context.Background(), region, cmkID)
}

// NewKmsClientContextE creates a KMS client.
// The ctx parameter supports cancellation and timeouts.
func NewKmsClientContextE(t testing.TestingT, ctx context.Context, region string) (*kms.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return kms.NewFromConfig(*sess), nil
}

// NewKmsClientContext creates a KMS client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewKmsClientContext(t testing.TestingT, ctx context.Context, region string) *kms.Client {
	t.Helper()

	client, err := NewKmsClientContextE(t, ctx, region)
	if err != nil {
		t.Fatal(err)
	}

	return client
}

// NewKmsClient creates a KMS client.
//
// Deprecated: Use [NewKmsClientContext] instead.
func NewKmsClient(t testing.TestingT, region string) *kms.Client {
	t.Helper()

	return NewKmsClientContext(t, context.Background(), region)
}

// NewKmsClientE creates a KMS client.
//
// Deprecated: Use [NewKmsClientContextE] instead.
func NewKmsClientE(t testing.TestingT, region string) (*kms.Client, error) {
	return NewKmsClientContextE(t, context.Background(), region)
}
