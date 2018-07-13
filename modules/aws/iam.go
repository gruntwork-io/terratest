package aws

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/stretchr/testify/assert"
)

// NoPolicyError is used to alert that a valid policy was found from an ARN but a document doesn't exist against it.
type NoPolicyError struct {
	error
}

func (e *NoPolicyError) Error() string {
	return "the policy exists but no policy document was returned"
}

// AssertJsonEqual takes two json strings and compares them
func AssertJsonEqual(t *testing.T, correct string, comparison string) {
	err := AssertJsonEqualE(t, correct, comparison)
	if err != nil {
		t.Fatal(err)
	}
	return
}

// AssertJsonEqualE takes two json strings and compares them, it returns an error if they're not identical
func AssertJsonEqualE(t *testing.T, correct string, comparison string) error {
	dst := new(bytes.Buffer)

	src := []byte(correct)

	err := json.Compact(dst, src)
	if err != nil {
		return err
	}

	equal := assert.JSONEq(
		t,
		dst.String(),
		comparison,
	)
	if equal == false {
		return errors.New("text")
	}

	return nil
}

// GetIamCurrentUserName gets the username for the current IAM user.
func GetIamCurrentUserName(t *testing.T) string {
	out, err := GetIamCurrentUserNameE(t)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetIamCurrentUserNameE gets the username for the current IAM user.
func GetIamCurrentUserNameE(t *testing.T) (string, error) {
	iamClient, err := NewIamClientE(t, defaultRegion)
	if err != nil {
		return "", err
	}

	resp, err := iamClient.GetUser(&iam.GetUserInput{})
	if err != nil {
		return "", err
	}

	return *resp.User.UserName, nil
}

// GetIamCurrentUserArn gets the ARN for the current IAM user.
func GetIamCurrentUserArn(t *testing.T) string {
	out, err := GetIamCurrentUserArnE(t)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// GetIamCurrentUserArnE gets the ARN for the current IAM user.
func GetIamCurrentUserArnE(t *testing.T) (string, error) {
	iamClient, err := NewIamClientE(t, defaultRegion)
	if err != nil {
		return "", err
	}

	resp, err := iamClient.GetUser(&iam.GetUserInput{})
	if err != nil {
		return "", err
	}

	return *resp.User.Arn, nil
}

// CreateMfaDevice creates an MFA device using the given IAM client.
func CreateMfaDevice(t *testing.T, iamClient *iam.IAM, deviceName string) *iam.VirtualMFADevice {
	mfaDevice, err := CreateMfaDeviceE(t, iamClient, deviceName)
	if err != nil {
		t.Fatal(err)
	}
	return mfaDevice
}

// CreateMfaDeviceE creates an MFA device using the given IAM client.
func CreateMfaDeviceE(t *testing.T, iamClient *iam.IAM, deviceName string) (*iam.VirtualMFADevice, error) {
	logger.Logf(t, "Creating an MFA device called %s", deviceName)

	output, err := iamClient.CreateVirtualMFADevice(&iam.CreateVirtualMFADeviceInput{
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
func EnableMfaDevice(t *testing.T, iamClient *iam.IAM, mfaDevice *iam.VirtualMFADevice) {
	err := EnableMfaDeviceE(t, iamClient, mfaDevice)
	if err != nil {
		t.Fatal(err)
	}
}

// EnableMfaDeviceE enables a newly created MFA Device by supplying the first two one-time passwords, so that it can be used for future
// logins by the given IAM User.
func EnableMfaDeviceE(t *testing.T, iamClient *iam.IAM, mfaDevice *iam.VirtualMFADevice) error {
	logger.Logf(t, "Enabling MFA device %s", aws.StringValue(mfaDevice.SerialNumber))

	iamUserName, err := GetIamCurrentUserArnE(t)
	if err != nil {
		return err
	}

	authCode1, err := GetTimeBasedOneTimePassword(mfaDevice)
	if err != nil {
		return err
	}

	logger.Logf(t, "Waiting 30 seconds for a new MFA Token to be generated...")
	time.Sleep(30 * time.Second)

	authCode2, err := GetTimeBasedOneTimePassword(mfaDevice)
	if err != nil {
		return err
	}

	_, err = iamClient.EnableMFADevice(&iam.EnableMFADeviceInput{
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

// AssertIamPolicyExists checks if the given IAM policy exists in the given region and fail the test if it does not.
func AssertIamPolicyExists(t *testing.T, region string, policyARN string) {
	_, err := GetIamPolicyDocumentE(t, region, policyARN)
	if err != nil {
		t.Fatal(err)
	}
}

// AssertIAMPolicyIsEqual fetches the contents of the IAM policy and matches it against the provided policy document.
func AssertIAMPolicyIsEqual(t *testing.T, region string, policyARN string, policyDocument string) {
	err := AssertIAMPolicyIsEqualE(t, region, policyARN, policyDocument)
	if err != nil {
		t.Fatal(err)
	}
}

// AssertIAMPolicyIsEqualE fetches the contents of the IAM policy and matches it against the provided policy document returning an error if they don't match.
func AssertIAMPolicyIsEqualE(t *testing.T, region string, policyARN string, policyDocument string) error {
	documentFromPolicy, err := GetIamPolicyDocumentE(t, region, policyARN)
	if err != nil {
		return err
	}

	err = AssertJsonEqualE(t, policyDocument, documentFromPolicy)
	if err != nil {
		return err
	}

	return nil
}

// GetIamPolicyDocument gets the most recent policy document for an IAM policy and fatals if it can't get one.
func GetIamPolicyDocument(t *testing.T, region string, policyARN string) string {
	content, err := GetIamPolicyDocumentE(t, region, policyARN)
	if err != nil {
		t.Fatal(err)
	}

	return content
}

// GetIamPolicyDocumentE gets the most recent policy document for an IAM policy and errors if it can't get one.
func GetIamPolicyDocumentE(t *testing.T, region string, policyARN string) (string, error) {
	iamClient, err := NewIamClientE(t, region)
	if err != nil {
		return "", err
	}

	versions, err := iamClient.ListPolicyVersions(&iam.ListPolicyVersionsInput{
		PolicyArn: &policyARN,
	})
	if err != nil {
		return "", err
	}

	var defaultVersion string

	for _, version := range versions.Versions {
		if *version.IsDefaultVersion == true {
			defaultVersion = *version.VersionId
		}
	}

	document, err := iamClient.GetPolicyVersion(&iam.GetPolicyVersionInput{
		PolicyArn: aws.String(policyARN),
		VersionId: aws.String(defaultVersion),
	})

	var unescapedDocument *string

	unescapedDocument = document.PolicyVersion.Document

	if unescapedDocument == nil {
		return "", new(NoPolicyError)
	}

	escapedDocument, err := url.QueryUnescape(*unescapedDocument)
	if err != nil {
		return "", err
	}

	return escapedDocument, nil
}

// NewIamClient creates a new IAM client.
func NewIamClient(t *testing.T, region string) *iam.IAM {
	client, err := NewIamClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewIamClientE creates a new IAM client.
func NewIamClientE(t *testing.T, region string) (*iam.IAM, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return iam.New(sess), nil
}
