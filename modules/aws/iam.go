package aws

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// GetIamCurrentUserName gets the username for the current IAM user.
func GetIamCurrentUserName(t testing.TestingT) string {
	out, err := GetIamCurrentUserNameE(t)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetIamCurrentUserNameE gets the username for the current IAM user.
func GetIamCurrentUserNameE(t testing.TestingT) (string, error) {
	iamClient, err := NewIamClientE(t, defaultRegion)
	if err != nil {
		return "", err
	}

	resp, err := iamClient.GetUser(context.Background(), &iam.GetUserInput{})
	if err != nil {
		return "", err
	}

	return *resp.User.UserName, nil
}

// GetIamCurrentUserArn gets the ARN for the current IAM user.
func GetIamCurrentUserArn(t testing.TestingT) string {
	out, err := GetIamCurrentUserArnE(t)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetIamCurrentUserArnE gets the ARN for the current IAM user.
func GetIamCurrentUserArnE(t testing.TestingT) (string, error) {
	iamClient, err := NewIamClientE(t, defaultRegion)
	if err != nil {
		return "", err
	}

	resp, err := iamClient.GetUser(context.Background(), &iam.GetUserInput{})
	if err != nil {
		return "", err
	}

	return *resp.User.Arn, nil
}

// GetIamPolicyDocument gets the most recent policy (JSON) document for an IAM policy.
func GetIamPolicyDocument(t testing.TestingT, region string, policyARN string) string {
	out, err := GetIamPolicyDocumentE(t, region, policyARN)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetIamPolicyDocumentE gets the most recent policy (JSON) document for an IAM policy.
func GetIamPolicyDocumentE(t testing.TestingT, region string, policyARN string) (string, error) {
	iamClient, err := NewIamClientE(t, region)
	if err != nil {
		return "", err
	}

	versions, err := iamClient.ListPolicyVersions(context.Background(), &iam.ListPolicyVersionsInput{
		PolicyArn: &policyARN,
	})
	if err != nil {
		return "", err
	}

	var defaultVersion string
	for _, version := range versions.Versions {
		if version.IsDefaultVersion == true {
			defaultVersion = *version.VersionId
		}
	}

	document, err := iamClient.GetPolicyVersion(context.Background(), &iam.GetPolicyVersionInput{
		PolicyArn: aws.String(policyARN),
		VersionId: aws.String(defaultVersion),
	})
	if err != nil {
		return "", err
	}

	unescapedDocument := document.PolicyVersion.Document
	if unescapedDocument == nil {
		return "", fmt.Errorf("no policy document found for policy %s", policyARN)
	}

	escapedDocument, err := url.QueryUnescape(*unescapedDocument)
	if err != nil {
		return "", err
	}

	return escapedDocument, nil
}

// CreateMfaDevice creates an MFA device using the given IAM client.
func CreateMfaDevice(t testing.TestingT, iamClient *iam.Client, deviceName string) *types.VirtualMFADevice {
	mfaDevice, err := CreateMfaDeviceE(t, iamClient, deviceName)
	if err != nil {
		t.Fatal(err)
	}
	return mfaDevice
}

// CreateMfaDeviceE creates an MFA device using the given IAM client.
func CreateMfaDeviceE(t testing.TestingT, iamClient *iam.Client, deviceName string) (*types.VirtualMFADevice, error) {
	logger.Default.Logf(t, "Creating an MFA device called %s", deviceName)

	output, err := iamClient.CreateVirtualMFADevice(context.Background(), &iam.CreateVirtualMFADeviceInput{
		VirtualMFADeviceName: aws.String(deviceName),
	})
	if err != nil {
		return nil, err
	}

	if err := EnableMfaDeviceE(t, iamClient, output.VirtualMFADevice); err != nil {
		return nil, err
	}

	return output.VirtualMFADevice, nil
}

// EnableMfaDevice enables a newly created MFA Device by supplying the first two one-time passwords, so that it can be used for future
// logins by the given IAM User.
func EnableMfaDevice(t testing.TestingT, iamClient *iam.Client, mfaDevice *types.VirtualMFADevice) {
	err := EnableMfaDeviceE(t, iamClient, mfaDevice)
	if err != nil {
		t.Fatal(err)
	}
}

// EnableMfaDeviceE enables a newly created MFA Device by supplying the first two one-time passwords, so that it can be used for future
// logins by the given IAM User.
func EnableMfaDeviceE(t testing.TestingT, iamClient *iam.Client, mfaDevice *types.VirtualMFADevice) error {
	logger.Default.Logf(t, "Enabling MFA device %s", aws.ToString(mfaDevice.SerialNumber))

	iamUserName, err := GetIamCurrentUserArnE(t)
	if err != nil {
		return err
	}

	authCode1, err := GetTimeBasedOneTimePassword(mfaDevice)
	if err != nil {
		return err
	}

	logger.Default.Logf(t, "Waiting 30 seconds for a new MFA Token to be generated...")
	time.Sleep(30 * time.Second)

	authCode2, err := GetTimeBasedOneTimePassword(mfaDevice)
	if err != nil {
		return err
	}

	_, err = iamClient.EnableMFADevice(context.Background(), &iam.EnableMFADeviceInput{
		AuthenticationCode1: aws.String(authCode1),
		AuthenticationCode2: aws.String(authCode2),
		SerialNumber:        mfaDevice.SerialNumber,
		UserName:            aws.String(iamUserName),
	})

	if err != nil {
		return err
	}

	logger.Log(t, "Waiting for MFA Device enablement to propagate.")
	time.Sleep(10 * time.Second)

	return nil
}

// NewIamClient creates a new IAM client.
func NewIamClient(t testing.TestingT, region string) *iam.Client {
	client, err := NewIamClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewIamClientE creates a new IAM client.
func NewIamClientE(t testing.TestingT, region string) (*iam.Client, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return iam.NewFromConfig(*sess), nil
}
