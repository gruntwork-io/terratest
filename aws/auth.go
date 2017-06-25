package aws

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pquerna/otp/totp"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gruntwork-io/terratest/util"
	"github.com/gruntwork-io/gruntwork-cli/logging"
)

// Create a new AWS session using the system environment to authenticate to AWS. For info on how to configure the system
// environment, see https://docs.aws.amazon.com/sdk-for-go/v1/developerguide/configuring-sdk.html.
func CreateAwsSession(awsRegion string) (*session.Session, error) {
	awsConfig := defaults.Get().Config.WithRegion(awsRegion)

	_, err := awsConfig.Credentials.Get()
	if err != nil {
		return nil, err
	}

	return session.New(awsConfig), nil
}

// Create a new AWS session using explicit credentials. This is useful if you want to create an IAM User dynamically and
// create an AWS session authenticated as the new IAM User.
func CreateAwsSessionWithCreds(awsRegion string, accessKeyId string, secretAccessKey string) *session.Session {
	creds := CreateAwsCredentials(accessKeyId, secretAccessKey)
	awsConfig := defaults.Get().Config.WithRegion(awsRegion).WithCredentials(creds)
	return session.New(awsConfig)
}

func CreateAwsSessionWithCredsAndMfa(awsRegion string, stsClient *sts.STS, iamClient *iam.IAM, mfaDevice *iam.VirtualMFADevice) (*session.Session, error) {
	// Get a one-time password
	tokenCode, err := GetTimeBasedOneTimePassword(mfaDevice)
	if err != nil {
		return nil, err
	}

	// Now get temp credentials from STS
	output, err := stsClient.GetSessionToken(&sts.GetSessionTokenInput{
		SerialNumber: mfaDevice.SerialNumber,
		TokenCode: aws.String(tokenCode),
	})
	if err != nil {
		return nil, err
	}

	accessKeyId := *output.Credentials.AccessKeyId
	secretAccessKey := *output.Credentials.SecretAccessKey
	sessionToken := *output.Credentials.SessionToken

	// Now authenticate a session with MFA
	creds := CreateAwsCredentialsWithSessionToken(accessKeyId, secretAccessKey, sessionToken)
	awsConfig := defaults.Get().Config.WithRegion(awsRegion).WithCredentials(creds)
	return session.New(awsConfig), nil
}

// Create an AWS configuration with specific AWS credentials.
func CreateAwsCredentials(accessKeyId string, secretAccessKey string) *credentials.Credentials {
	creds := credentials.Value{AccessKeyID: accessKeyId, SecretAccessKey: secretAccessKey }
	return credentials.NewStaticCredentialsFromCreds(creds)
}

// Create an AWS configuration with temporary AWS credentials by including a session token (used for authenticating with MFA)
func CreateAwsCredentialsWithSessionToken(accessKeyId, secretAccessKey, sessionToken string) *credentials.Credentials {
	creds := credentials.Value{
		AccessKeyID: accessKeyId,
		SecretAccessKey: secretAccessKey,
		SessionToken: sessionToken,
	}
	return credentials.NewStaticCredentialsFromCreds(creds)
}

func CreateMfaDevice(iamClient *iam.IAM, deviceName string) (*iam.VirtualMFADevice, error) {
	output, err := iamClient.CreateVirtualMFADevice(&iam.CreateVirtualMFADeviceInput{
		VirtualMFADeviceName: aws.String(deviceName),
	})
	if err != nil {
		return nil, err
	}

	mfaDevice := output.VirtualMFADevice

	EnableMfaDevice(iamClient, mfaDevice)

	return mfaDevice, nil
}

// Enable a newly created MFA Device (by supplying the first two one-time passwords) so that it can be used for future
// logins by the given IAM User
func EnableMfaDevice(iamClient *iam.IAM, mfaDevice *iam.VirtualMFADevice) (error) {
	iamUserName, err := GetUserName(iamClient)
	if err != nil {
		return err
	}

	authCode1, err := GetTimeBasedOneTimePassword(mfaDevice)
	if err != nil {
		return err
	}

	logger := logging.GetLogger("EnableMfaDevice")
	logger.Debug("Waiting 30 seconds for a new MFA Token to be generated...")
	time.Sleep(30 * time.Second)

	authCode2, err := GetTimeBasedOneTimePassword(mfaDevice)
	if err != nil {
		return err
	}

	_, err = iamClient.EnableMFADevice(&iam.EnableMFADeviceInput{
		AuthenticationCode1: aws.String(authCode1),
		AuthenticationCode2: aws.String(authCode2),
		SerialNumber: mfaDevice.SerialNumber,
		UserName: aws.String(iamUserName),
	})

	if err != nil {
		return err
	}

	util.SleepWithMessage(logger, 10 * time.Second, "Waiting for MFA Device enablement to propagate.")
	return nil
}

// Get the user name of the given IAM Client session
func GetUserName(iamClient *iam.IAM) (string, error) {
	output, err := iamClient.GetUser(&iam.GetUserInput{})
	if err != nil {
		return "", err
	}

	return *output.User.UserName, nil
}

// Get a One-Time Password from the given mfaDevice. Per the RFC 6238 standard, this value will be different every 30 seconds.
func GetTimeBasedOneTimePassword(mfaDevice *iam.VirtualMFADevice) (string, error) {
	base32StringSeed := string(mfaDevice.Base32StringSeed)

	otp, err := totp.GenerateCode(base32StringSeed, time.Now())
	if err != nil {
		return "", err
	}

	return otp, nil
}

func ReadPasswordPolicyMinPasswordLength(iamClient *iam.IAM) (int, error) {
	output, err := iamClient.GetAccountPasswordPolicy(&iam.GetAccountPasswordPolicyInput{})
	if err != nil {
		return -1, err
	}

	return int(*output.PasswordPolicy.MinimumPasswordLength), nil
}
