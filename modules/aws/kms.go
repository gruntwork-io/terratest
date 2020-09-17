package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// GetCmkArn gets the ARN of a KMS Customer Master Key (CMK) in the given region with the given ID. The ID can be an alias, such
// as "alias/my-cmk".
func GetCmkArn(t testing.TestingT, region string, cmkID string) string {
	out, err := GetCmkArnE(t, region, cmkID)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetCmkArnE gets the ARN of a KMS Customer Master Key (CMK) in the given region with the given ID. The ID can be an alias, such
// as "alias/my-cmk".
func GetCmkArnE(t testing.TestingT, region string, cmkID string) (string, error) {
	kmsClient, err := NewKmsClientE(t, region)
	if err != nil {
		return "", err
	}

	result, err := kmsClient.DescribeKey(&kms.DescribeKeyInput{
		KeyId: aws.String(cmkID),
	})

	if err != nil {
		return "", err
	}

	return *result.KeyMetadata.Arn, nil
}

// NewKmsClient creates a KMS client.
func NewKmsClient(t testing.TestingT, region string) *kms.KMS {
	client, err := NewKmsClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewKmsClientE creates a KMS client.
func NewKmsClientE(t testing.TestingT, region string) (*kms.KMS, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return kms.New(sess), nil
}
