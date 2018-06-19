package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

// GetCmkArn gets the ARN of a KMS Customer Master Key (CMK) in the given region with the given ID. The ID can be an alias, such
// as "alias/my-cmk".
func GetCmkArn(t *testing.T, region string, cmkID string, sessExists ...*session.Session) string {
	out, err := GetCmkArnE(t, region, cmkID, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetCmkArnE gets the ARN of a KMS Customer Master Key (CMK) in the given region with the given ID. The ID can be an alias, such
// as "alias/my-cmk".
func GetCmkArnE(t *testing.T, region string, cmkID string, sessExists ...*session.Session) (string, error) {
	kmsClient, err := NewKmsClientE(t, region, sessExists[0])
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
func NewKmsClient(t *testing.T, region string, sessExists ...*session.Session) *kms.KMS {
	client, err := NewKmsClientE(t, region, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewKmsClientE creates a KMS client.
func NewKmsClientE(t *testing.T, region string, sessExists ...*session.Session) (*kms.KMS, error) {
	sess, err := NewAuthenticatedSession(region, sessExists[0])
	if err != nil {
		return nil, err
	}

	return kms.New(sess), nil
}
