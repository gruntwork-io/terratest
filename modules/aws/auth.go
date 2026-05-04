package aws

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/pquerna/otp/totp"
)

const (
	// AuthAssumeRoleEnvVar is the OS environment variable name through which an
	// Assume Role ARN may be passed for authentication.
	AuthAssumeRoleEnvVar = "TERRATEST_IAM_ROLE"
)

// NewAuthenticatedSessionContextE creates an AWS Config following to standard AWS authentication workflow.
// If AuthAssumeIamRoleEnvVar environment variable is set, assumes IAM role specified in it.
// The ctx parameter supports cancellation and timeouts.
func NewAuthenticatedSessionContextE(ctx context.Context, region string) (*aws.Config, error) {
	if assumeRoleArn, ok := os.LookupEnv(AuthAssumeRoleEnvVar); ok {
		return NewAuthenticatedSessionFromRoleContextE(ctx, region, assumeRoleArn)
	}

	return NewAuthenticatedSessionFromDefaultCredentialsContextE(ctx, region)
}

// Deprecated: Use [NewAuthenticatedSessionContextE] instead.
func NewAuthenticatedSessionContext(ctx context.Context, region string) (*aws.Config, error) {
	return NewAuthenticatedSessionContextE(ctx, region)
}

// NewAuthenticatedSession creates an AWS Config following to standard AWS authentication workflow.
// If AuthAssumeIamRoleEnvVar environment variable is set, assumes IAM role specified in it.
//
// Deprecated: Use [NewAuthenticatedSessionContextE] instead.
func NewAuthenticatedSession(region string) (*aws.Config, error) {
	return NewAuthenticatedSessionContextE(context.Background(), region)
}

// NewAuthenticatedSessionFromDefaultCredentialsContextE gets an AWS Config, checking that the user has credentials properly configured in their environment.
// The ctx parameter supports cancellation and timeouts.
func NewAuthenticatedSessionFromDefaultCredentialsContextE(ctx context.Context, region string) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, CredentialsError{UnderlyingErr: err}
	}

	return &cfg, nil
}

// Deprecated: Use [NewAuthenticatedSessionFromDefaultCredentialsContextE] instead.
func NewAuthenticatedSessionFromDefaultCredentialsContext(ctx context.Context, region string) (*aws.Config, error) {
	return NewAuthenticatedSessionFromDefaultCredentialsContextE(ctx, region)
}

// NewAuthenticatedSessionFromDefaultCredentials gets an AWS Config, checking that the user has credentials properly configured in their environment.
//
// Deprecated: Use [NewAuthenticatedSessionFromDefaultCredentialsContextE] instead.
func NewAuthenticatedSessionFromDefaultCredentials(region string) (*aws.Config, error) {
	return NewAuthenticatedSessionFromDefaultCredentialsContextE(context.Background(), region)
}

// NewAuthenticatedSessionFromRoleContextE returns a new AWS Config after assuming the
// role whose ARN is provided in roleARN. If the credentials are not properly
// configured in the underlying environment, an error is returned.
// The ctx parameter supports cancellation and timeouts.
func NewAuthenticatedSessionFromRoleContextE(ctx context.Context, region string, roleARN string) (*aws.Config, error) {
	cfg, err := NewAuthenticatedSessionFromDefaultCredentialsContextE(ctx, region)
	if err != nil {
		return nil, err
	}

	client := sts.NewFromConfig(*cfg)

	roleProvider := stscreds.NewAssumeRoleProvider(client, roleARN)

	retrieve, err := roleProvider.Retrieve(ctx)
	if err != nil {
		return nil, CredentialsError{UnderlyingErr: err}
	}

	return &aws.Config{
		Region: region,
		Credentials: aws.NewCredentialsCache(credentials.StaticCredentialsProvider{
			Value: retrieve,
		}),
	}, nil
}

// Deprecated: Use [NewAuthenticatedSessionFromRoleContextE] instead.
func NewAuthenticatedSessionFromRoleContext(ctx context.Context, region string, roleARN string) (*aws.Config, error) {
	return NewAuthenticatedSessionFromRoleContextE(ctx, region, roleARN)
}

// NewAuthenticatedSessionFromRole returns a new AWS Config after assuming the
// role whose ARN is provided in roleARN. If the credentials are not properly
// configured in the underlying environment, an error is returned.
//
// Deprecated: Use [NewAuthenticatedSessionFromRoleContextE] instead.
func NewAuthenticatedSessionFromRole(region string, roleARN string) (*aws.Config, error) {
	return NewAuthenticatedSessionFromRoleContextE(context.Background(), region, roleARN)
}

// CreateAwsSessionWithCredsContextE creates a new AWS Config using explicit credentials. This is useful if you want to create an IAM User dynamically and
// create an AWS Config authenticated as the new IAM User.
// The ctx parameter is accepted for API consistency but not currently used.
func CreateAwsSessionWithCredsContextE(_ context.Context, region string, accessKeyID string, secretAccessKey string) (*aws.Config, error) {
	return &aws.Config{
		Region:      region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	}, nil
}

// Deprecated: Use [CreateAwsSessionWithCredsContextE] instead.
func CreateAwsSessionWithCredsContext(ctx context.Context, region string, accessKeyID string, secretAccessKey string) (*aws.Config, error) {
	return CreateAwsSessionWithCredsContextE(ctx, region, accessKeyID, secretAccessKey)
}

// CreateAwsSessionWithCreds creates a new AWS Config using explicit credentials. This is useful if you want to create an IAM User dynamically and
// create an AWS Config authenticated as the new IAM User.
//
// Deprecated: Use [CreateAwsSessionWithCredsContextE] instead.
func CreateAwsSessionWithCreds(region string, accessKeyID string, secretAccessKey string) (*aws.Config, error) {
	return CreateAwsSessionWithCredsContextE(context.Background(), region, accessKeyID, secretAccessKey)
}

// CreateAwsSessionWithMfaContextE creates a new AWS Config authenticated using an MFA token retrieved using the given STS client and MFA Device.
// The ctx parameter supports cancellation and timeouts.
func CreateAwsSessionWithMfaContextE(ctx context.Context, region string, stsClient *sts.Client, mfaDevice *types.VirtualMFADevice) (*aws.Config, error) {
	tokenCode, err := GetTimeBasedOneTimePassword(mfaDevice)
	if err != nil {
		return nil, err
	}

	output, err := stsClient.GetSessionToken(ctx, &sts.GetSessionTokenInput{
		SerialNumber: mfaDevice.SerialNumber,
		TokenCode:    aws.String(tokenCode),
	})
	if err != nil {
		return nil, err
	}

	accessKeyID := *output.Credentials.AccessKeyId
	secretAccessKey := *output.Credentials.SecretAccessKey
	sessionToken := *output.Credentials.SessionToken

	return &aws.Config{
		Region:      region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, sessionToken)),
	}, nil
}

// Deprecated: Use [CreateAwsSessionWithMfaContextE] instead.
func CreateAwsSessionWithMfaContext(ctx context.Context, region string, stsClient *sts.Client, mfaDevice *types.VirtualMFADevice) (*aws.Config, error) {
	return CreateAwsSessionWithMfaContextE(ctx, region, stsClient, mfaDevice)
}

// CreateAwsSessionWithMfa creates a new AWS Config authenticated using an MFA token retrieved using the given STS client and MFA Device.
//
// Deprecated: Use [CreateAwsSessionWithMfaContextE] instead.
func CreateAwsSessionWithMfa(region string, stsClient *sts.Client, mfaDevice *types.VirtualMFADevice) (*aws.Config, error) {
	return CreateAwsSessionWithMfaContextE(context.Background(), region, stsClient, mfaDevice)
}

// GetTimeBasedOneTimePassword gets a One-Time Password from the given mfaDevice. Per the RFC 6238 standard, this value will be different every 30 seconds.
func GetTimeBasedOneTimePassword(mfaDevice *types.VirtualMFADevice) (string, error) {
	base32StringSeed := string(mfaDevice.Base32StringSeed)

	otp, err := totp.GenerateCode(base32StringSeed, time.Now())
	if err != nil {
		return "", err
	}

	return otp, nil
}

// ReadPasswordPolicyMinPasswordLengthContextE returns the minimal password length.
// The ctx parameter supports cancellation and timeouts.
func ReadPasswordPolicyMinPasswordLengthContextE(ctx context.Context, iamClient *iam.Client) (int, error) {
	output, err := iamClient.GetAccountPasswordPolicy(ctx, &iam.GetAccountPasswordPolicyInput{})
	if err != nil {
		return -1, err
	}

	return int(*output.PasswordPolicy.MinimumPasswordLength), nil
}

// Deprecated: Use [ReadPasswordPolicyMinPasswordLengthContextE] instead.
func ReadPasswordPolicyMinPasswordLengthContext(ctx context.Context, iamClient *iam.Client) (int, error) {
	return ReadPasswordPolicyMinPasswordLengthContextE(ctx, iamClient)
}

// ReadPasswordPolicyMinPasswordLength returns the minimal password length.
//
// Deprecated: Use [ReadPasswordPolicyMinPasswordLengthContextE] instead.
func ReadPasswordPolicyMinPasswordLength(iamClient *iam.Client) (int, error) {
	return ReadPasswordPolicyMinPasswordLengthContextE(context.Background(), iamClient)
}

// CredentialsError is an error that occurs because AWS credentials can't be found.
type CredentialsError struct {
	UnderlyingErr error
}

func (err CredentialsError) Error() string {
	return fmt.Sprintf("Error finding AWS credentials. Did you set the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables or configure an AWS profile? Underlying error: %v", err.UnderlyingErr)
}
