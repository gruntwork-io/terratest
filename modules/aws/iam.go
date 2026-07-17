package aws

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// GetIamCurrentUserNameContextE gets the username for the current IAM user.
// The ctx parameter supports cancellation and timeouts.
func GetIamCurrentUserNameContextE(t testing.TestingT, ctx context.Context) (string, error) {
	iamClient, err := NewIamClientContextE(t, ctx, defaultRegion)
	if err != nil {
		return "", err
	}

	resp, err := iamClient.GetUser(ctx, &iam.GetUserInput{})
	if err != nil {
		return "", err
	}

	return *resp.User.UserName, nil
}

// GetIamCurrentUserNameContext gets the username for the current IAM user.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIamCurrentUserNameContext(t testing.TestingT, ctx context.Context) string {
	t.Helper()

	out, err := GetIamCurrentUserNameContextE(t, ctx)
	require.NoError(t, err)

	return out
}

// GetIamCurrentUserArnContextE gets the ARN for the current IAM user.
// The ctx parameter supports cancellation and timeouts.
func GetIamCurrentUserArnContextE(t testing.TestingT, ctx context.Context) (string, error) {
	iamClient, err := NewIamClientContextE(t, ctx, defaultRegion)
	if err != nil {
		return "", err
	}

	resp, err := iamClient.GetUser(ctx, &iam.GetUserInput{})
	if err != nil {
		return "", err
	}

	return *resp.User.Arn, nil
}

// GetIamCurrentUserArnContext gets the ARN for the current IAM user.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIamCurrentUserArnContext(t testing.TestingT, ctx context.Context) string {
	t.Helper()

	out, err := GetIamCurrentUserArnContextE(t, ctx)
	require.NoError(t, err)

	return out
}

// GetIamPolicyDocumentContextE gets the most recent policy (JSON) document for an IAM policy.
// The ctx parameter supports cancellation and timeouts.
func GetIamPolicyDocumentContextE(t testing.TestingT, ctx context.Context, region string, policyARN string) (string, error) {
	iamClient, err := NewIamClientContextE(t, ctx, region)
	if err != nil {
		return "", err
	}

	var defaultVersion string

	paginator := iam.NewListPolicyVersionsPaginator(iamClient, &iam.ListPolicyVersionsInput{
		PolicyArn: &policyARN,
	})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return "", err
		}

		for _, version := range page.Versions {
			if version.IsDefaultVersion && version.VersionId != nil {
				defaultVersion = *version.VersionId
			}
		}
	}

	if defaultVersion == "" {
		return "", fmt.Errorf("no default version found for IAM policy %s", policyARN)
	}

	document, err := iamClient.GetPolicyVersion(ctx, &iam.GetPolicyVersionInput{
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

// GetIamPolicyDocumentContext gets the most recent policy (JSON) document for an IAM policy.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetIamPolicyDocumentContext(t testing.TestingT, ctx context.Context, region string, policyARN string) string {
	t.Helper()

	out, err := GetIamPolicyDocumentContextE(t, ctx, region, policyARN)
	require.NoError(t, err)

	return out
}

// CreateMfaDeviceContextE creates an MFA device using the given IAM client.
// The ctx parameter supports cancellation and timeouts.
func CreateMfaDeviceContextE(t testing.TestingT, ctx context.Context, iamClient *iam.Client, deviceName string) (*types.VirtualMFADevice, error) {
	logger.Default.Logf(t, "Creating an MFA device called %s", deviceName)

	output, err := iamClient.CreateVirtualMFADevice(ctx, &iam.CreateVirtualMFADeviceInput{
		VirtualMFADeviceName: aws.String(deviceName),
	})
	if err != nil {
		return nil, err
	}

	if err := EnableMfaDeviceContextE(t, ctx, iamClient, output.VirtualMFADevice); err != nil {
		return nil, err
	}

	return output.VirtualMFADevice, nil
}

// CreateMfaDeviceContext creates an MFA device using the given IAM client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CreateMfaDeviceContext(t testing.TestingT, ctx context.Context, iamClient *iam.Client, deviceName string) *types.VirtualMFADevice {
	t.Helper()

	mfaDevice, err := CreateMfaDeviceContextE(t, ctx, iamClient, deviceName)
	require.NoError(t, err)

	return mfaDevice
}

// EnableMfaDeviceContextE enables a newly created MFA Device by supplying the first two one-time passwords, so that it can be used for future
// logins by the given IAM User.
// The ctx parameter supports cancellation and timeouts.
func EnableMfaDeviceContextE(t testing.TestingT, ctx context.Context, iamClient *iam.Client, mfaDevice *types.VirtualMFADevice) error {
	logger.Default.Logf(t, "Enabling MFA device %s", aws.ToString(mfaDevice.SerialNumber))

	iamUserArn, err := GetIamCurrentUserArnContextE(t, ctx)
	if err != nil {
		return err
	}

	authCode1, err := GetTimeBasedOneTimePassword(mfaDevice)
	if err != nil {
		return err
	}

	const mfaEnableWait = 30 * time.Second

	logger.Default.Logf(t, "Waiting 30 seconds for a new MFA Token to be generated...")
	time.Sleep(mfaEnableWait)

	authCode2, err := GetTimeBasedOneTimePassword(mfaDevice)
	if err != nil {
		return err
	}

	_, err = iamClient.EnableMFADevice(ctx, &iam.EnableMFADeviceInput{
		AuthenticationCode1: aws.String(authCode1),
		AuthenticationCode2: aws.String(authCode2),
		SerialNumber:        mfaDevice.SerialNumber,
		UserName:            aws.String(iamUserArn),
	})
	if err != nil {
		return err
	}

	const mfaTokenWait = 10 * time.Second

	logger.Log(t, "Waiting for MFA Device enablement to propagate.")
	time.Sleep(mfaTokenWait)

	return nil
}

// EnableMfaDeviceContext enables a newly created MFA Device by supplying the first two one-time passwords, so that it can be used for future
// logins by the given IAM User.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func EnableMfaDeviceContext(t testing.TestingT, ctx context.Context, iamClient *iam.Client, mfaDevice *types.VirtualMFADevice) {
	t.Helper()

	err := EnableMfaDeviceContextE(t, ctx, iamClient, mfaDevice)
	require.NoError(t, err)
}

// NewIamClientContextE creates a new IAM client.
// The ctx parameter supports cancellation and timeouts.
func NewIamClientContextE(t testing.TestingT, ctx context.Context, region string) (*iam.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return iam.NewFromConfig(*sess), nil
}

// NewIamClientContext creates a new IAM client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewIamClientContext(t testing.TestingT, ctx context.Context, region string) *iam.Client {
	t.Helper()

	client, err := NewIamClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}
