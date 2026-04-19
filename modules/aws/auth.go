package aws

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/require"
)

const (
	AuthAssumeRoleEnvVar = "TERRATEST_IAM_ROLE" // OS environment variable name through which Assume Role ARN may be passed for authentication
)

// NewAuthConfigContextE creates an AWS Config following the standard AWS authentication workflow.
// If AuthAssumeRoleEnvVar is set, assumes the IAM role specified in it; otherwise falls back to
// the default credential chain.
// The ctx parameter supports cancellation and timeouts.
func NewAuthConfigContextE(t testing.TestingT, ctx context.Context, region string) (*aws.Config, error) {
	if assumeRoleArn, ok := os.LookupEnv(AuthAssumeRoleEnvVar); ok {
		return NewAuthConfigFromRoleContextE(t, ctx, region, assumeRoleArn)
	}

	return NewAuthConfigFromDefaultCredentialsContextE(t, ctx, region)
}

// NewAuthConfigContext creates an AWS Config following the standard AWS authentication workflow.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewAuthConfigContext(t testing.TestingT, ctx context.Context, region string) *aws.Config {
	t.Helper()
	cfg, err := NewAuthConfigContextE(t, ctx, region)
	require.NoError(t, err)
	return cfg
}

// NewAuthConfigE creates an AWS Config following the standard AWS authentication workflow.
func NewAuthConfigE(t testing.TestingT, region string) (*aws.Config, error) {
	return NewAuthConfigContextE(t, context.Background(), region)
}

// NewAuthConfig creates an AWS Config following the standard AWS authentication workflow.
// This function will fail the test if there is an error.
func NewAuthConfig(t testing.TestingT, region string) *aws.Config {
	t.Helper()
	return NewAuthConfigContext(t, context.Background(), region)
}

// NewAuthConfigFromDefaultCredentialsContextE gets an AWS Config using the default credential
// chain, checking that the user has credentials properly configured in their environment.
// The ctx parameter supports cancellation and timeouts.
func NewAuthConfigFromDefaultCredentialsContextE(t testing.TestingT, ctx context.Context, region string) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, CredentialsError{UnderlyingErr: err}
	}

	return &cfg, nil
}

// NewAuthConfigFromDefaultCredentialsContext gets an AWS Config using the default credential chain.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewAuthConfigFromDefaultCredentialsContext(t testing.TestingT, ctx context.Context, region string) *aws.Config {
	t.Helper()
	cfg, err := NewAuthConfigFromDefaultCredentialsContextE(t, ctx, region)
	require.NoError(t, err)
	return cfg
}

// NewAuthConfigFromDefaultCredentialsE gets an AWS Config using the default credential chain.
func NewAuthConfigFromDefaultCredentialsE(t testing.TestingT, region string) (*aws.Config, error) {
	return NewAuthConfigFromDefaultCredentialsContextE(t, context.Background(), region)
}

// NewAuthConfigFromDefaultCredentials gets an AWS Config using the default credential chain.
// This function will fail the test if there is an error.
func NewAuthConfigFromDefaultCredentials(t testing.TestingT, region string) *aws.Config {
	t.Helper()
	return NewAuthConfigFromDefaultCredentialsContext(t, context.Background(), region)
}

// NewAuthConfigFromRoleContextE returns a new AWS Config after assuming the role whose ARN
// is provided in roleARN. If the credentials are not properly configured in the underlying
// environment, an error is returned.
// The ctx parameter supports cancellation and timeouts.
func NewAuthConfigFromRoleContextE(t testing.TestingT, ctx context.Context, region string, roleARN string) (*aws.Config, error) {
	cfg, err := NewAuthConfigFromDefaultCredentialsContextE(t, ctx, region)
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

// NewAuthConfigFromRoleContext returns a new AWS Config after assuming the role whose ARN is provided in roleARN.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewAuthConfigFromRoleContext(t testing.TestingT, ctx context.Context, region string, roleARN string) *aws.Config {
	t.Helper()
	cfg, err := NewAuthConfigFromRoleContextE(t, ctx, region, roleARN)
	require.NoError(t, err)
	return cfg
}

// NewAuthConfigFromRoleE returns a new AWS Config after assuming the role whose ARN is provided in roleARN.
func NewAuthConfigFromRoleE(t testing.TestingT, region string, roleARN string) (*aws.Config, error) {
	return NewAuthConfigFromRoleContextE(t, context.Background(), region, roleARN)
}

// NewAuthConfigFromRole returns a new AWS Config after assuming the role whose ARN is provided in roleARN.
// This function will fail the test if there is an error.
func NewAuthConfigFromRole(t testing.TestingT, region string, roleARN string) *aws.Config {
	t.Helper()
	return NewAuthConfigFromRoleContext(t, context.Background(), region, roleARN)
}

// CreateAwsSessionWithCredsE creates a new AWS Config using explicit static credentials. Useful
// for authenticating as a dynamically-created IAM User. No context variant is provided because
// this function performs no I/O — it only constructs an aws.Config. No panic variant is provided
// because CreateAwsSessionWithCreds is already taken by the deprecated error-returning shim.
func CreateAwsSessionWithCredsE(t testing.TestingT, region string, accessKeyID string, secretAccessKey string) (*aws.Config, error) {
	return &aws.Config{
		Region:      region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	}, nil
}

// CreateAwsSessionWithMfaContextE creates a new AWS Config authenticated with an MFA session
// token obtained via the given STS client and MFA device.
// The ctx parameter supports cancellation and timeouts.
func CreateAwsSessionWithMfaContextE(t testing.TestingT, ctx context.Context, region string, stsClient *sts.Client, mfaDevice *types.VirtualMFADevice) (*aws.Config, error) {
	tokenCode, err := GetTimeBasedOneTimePasswordE(t, mfaDevice)
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

// CreateAwsSessionWithMfaE creates a new AWS Config authenticated with an MFA session token.
func CreateAwsSessionWithMfaE(t testing.TestingT, region string, stsClient *sts.Client, mfaDevice *types.VirtualMFADevice) (*aws.Config, error) {
	return CreateAwsSessionWithMfaContextE(t, context.Background(), region, stsClient, mfaDevice)
}

// GetTimeBasedOneTimePasswordE returns a time-based one-time password for the given MFA device,
// per RFC 6238. The returned value changes every 30 seconds. No context variant is provided
// because this function performs no I/O — it is a local HMAC-based computation.
func GetTimeBasedOneTimePasswordE(t testing.TestingT, mfaDevice *types.VirtualMFADevice) (string, error) {
	base32StringSeed := string(mfaDevice.Base32StringSeed)

	otp, err := totp.GenerateCode(base32StringSeed, time.Now())
	if err != nil {
		return "", err
	}

	return otp, nil
}

// ReadPasswordPolicyMinPasswordLengthContextE returns the minimum password length from the
// account's IAM password policy.
// The ctx parameter supports cancellation and timeouts.
func ReadPasswordPolicyMinPasswordLengthContextE(t testing.TestingT, ctx context.Context, iamClient *iam.Client) (int, error) {
	output, err := iamClient.GetAccountPasswordPolicy(ctx, &iam.GetAccountPasswordPolicyInput{})
	if err != nil {
		return -1, err
	}

	return int(*output.PasswordPolicy.MinimumPasswordLength), nil
}

// ReadPasswordPolicyMinPasswordLengthE returns the minimum password length from the account's
// IAM password policy.
func ReadPasswordPolicyMinPasswordLengthE(t testing.TestingT, iamClient *iam.Client) (int, error) {
	return ReadPasswordPolicyMinPasswordLengthContextE(t, context.Background(), iamClient)
}

// CredentialsError is an error that occurs because AWS credentials can't be found.
type CredentialsError struct {
	UnderlyingErr error
}

func (err CredentialsError) Error() string {
	return "Error finding AWS credentials. Did you set the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables or configure an AWS profile?"
}

// Unwrap returns the underlying error so callers can inspect it via errors.As / errors.Is
// without it being embedded in the user-facing message (which otherwise leaks into CI logs).
func (err CredentialsError) Unwrap() error {
	return err.UnderlyingErr
}

// -----------------------------------------------------------------------------
// Deprecated: legacy surface preserved for backwards compatibility.
//
// These functions predate the package's testing.TestingT + Context[E] convention.
// New callers should use the NewAuthConfig* family (or the *E variants of the original names
// for MFA / credentials / password-policy helpers) defined above.
// The shims below will be removed in a future major release.
// -----------------------------------------------------------------------------

// NewAuthenticatedSessionContext creates an AWS Config following to standard AWS authentication workflow.
// If AuthAssumeIamRoleEnvVar environment variable is set, assumes IAM role specified in it.
//
// Deprecated: Use [NewAuthConfigContextE] (or [NewAuthConfigContext] for panic-on-error semantics) instead.
func NewAuthenticatedSessionContext(ctx context.Context, region string) (*aws.Config, error) {
	if assumeRoleArn, ok := os.LookupEnv(AuthAssumeRoleEnvVar); ok {
		return NewAuthenticatedSessionFromRoleContext(ctx, region, assumeRoleArn)
	}

	return NewAuthenticatedSessionFromDefaultCredentialsContext(ctx, region)
}

// NewAuthenticatedSession creates an AWS Config following to standard AWS authentication workflow.
// If AuthAssumeIamRoleEnvVar environment variable is set, assumes IAM role specified in it.
//
// Deprecated: Use [NewAuthConfigE] (or [NewAuthConfig] for panic-on-error semantics) instead.
func NewAuthenticatedSession(region string) (*aws.Config, error) {
	return NewAuthenticatedSessionContext(context.Background(), region)
}

// NewAuthenticatedSessionFromDefaultCredentialsContext gets an AWS Config, checking that the user has credentials properly configured in their environment.
//
// Deprecated: Use [NewAuthConfigFromDefaultCredentialsContextE] instead.
func NewAuthenticatedSessionFromDefaultCredentialsContext(ctx context.Context, region string) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, CredentialsError{UnderlyingErr: err}
	}

	return &cfg, nil
}

// NewAuthenticatedSessionFromDefaultCredentials gets an AWS Config, checking that the user has credentials properly configured in their environment.
//
// Deprecated: Use [NewAuthConfigFromDefaultCredentialsE] instead.
func NewAuthenticatedSessionFromDefaultCredentials(region string) (*aws.Config, error) {
	return NewAuthenticatedSessionFromDefaultCredentialsContext(context.Background(), region)
}

// NewAuthenticatedSessionFromRoleContext returns a new AWS Config after assuming the
// role whose ARN is provided in roleARN. If the credentials are not properly
// configured in the underlying environment, an error is returned.
//
// Deprecated: Use [NewAuthConfigFromRoleContextE] instead.
func NewAuthenticatedSessionFromRoleContext(ctx context.Context, region string, roleARN string) (*aws.Config, error) {
	cfg, err := NewAuthenticatedSessionFromDefaultCredentialsContext(ctx, region)
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

// NewAuthenticatedSessionFromRole returns a new AWS Config after assuming the
// role whose ARN is provided in roleARN. If the credentials are not properly
// configured in the underlying environment, an error is returned.
//
// Deprecated: Use [NewAuthConfigFromRoleE] instead.
func NewAuthenticatedSessionFromRole(region string, roleARN string) (*aws.Config, error) {
	return NewAuthenticatedSessionFromRoleContext(context.Background(), region, roleARN)
}

// CreateAwsSessionWithCredsContext creates a new AWS Config using explicit credentials.
//
// Deprecated: Use [CreateAwsSessionWithCredsE] instead. The Context variant was a placeholder
// for API symmetry — the underlying operation performs no I/O, so no context is needed.
func CreateAwsSessionWithCredsContext(ctx context.Context, region string, accessKeyID string, secretAccessKey string) (*aws.Config, error) {
	return &aws.Config{
		Region:      region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	}, nil
}

// CreateAwsSessionWithCreds creates a new AWS Config using explicit credentials.
//
// Deprecated: Use [CreateAwsSessionWithCredsE] instead.
func CreateAwsSessionWithCreds(region string, accessKeyID string, secretAccessKey string) (*aws.Config, error) {
	return CreateAwsSessionWithCredsContext(context.Background(), region, accessKeyID, secretAccessKey)
}

// CreateAwsSessionWithMfaContext creates a new AWS Config authenticated using an MFA token retrieved using the given STS client and MFA Device.
//
// Deprecated: Use [CreateAwsSessionWithMfaContextE] instead.
func CreateAwsSessionWithMfaContext(ctx context.Context, region string, stsClient *sts.Client, mfaDevice *types.VirtualMFADevice) (*aws.Config, error) {
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

// CreateAwsSessionWithMfa creates a new AWS Config authenticated using an MFA token retrieved using the given STS client and MFA Device.
//
// Deprecated: Use [CreateAwsSessionWithMfaE] instead.
func CreateAwsSessionWithMfa(region string, stsClient *sts.Client, mfaDevice *types.VirtualMFADevice) (*aws.Config, error) {
	return CreateAwsSessionWithMfaContext(context.Background(), region, stsClient, mfaDevice)
}

// GetTimeBasedOneTimePassword gets a One-Time Password from the given mfaDevice. Per the RFC 6238 standard, this value will be different every 30 seconds.
//
// Deprecated: Use [GetTimeBasedOneTimePasswordE] instead.
func GetTimeBasedOneTimePassword(mfaDevice *types.VirtualMFADevice) (string, error) {
	base32StringSeed := string(mfaDevice.Base32StringSeed)

	otp, err := totp.GenerateCode(base32StringSeed, time.Now())
	if err != nil {
		return "", err
	}

	return otp, nil
}

// ReadPasswordPolicyMinPasswordLengthContext returns the minimal password length.
//
// Deprecated: Use [ReadPasswordPolicyMinPasswordLengthContextE] instead.
func ReadPasswordPolicyMinPasswordLengthContext(ctx context.Context, iamClient *iam.Client) (int, error) {
	output, err := iamClient.GetAccountPasswordPolicy(ctx, &iam.GetAccountPasswordPolicyInput{})
	if err != nil {
		return -1, err
	}

	return int(*output.PasswordPolicy.MinimumPasswordLength), nil
}

// ReadPasswordPolicyMinPasswordLength returns the minimal password length.
//
// Deprecated: Use [ReadPasswordPolicyMinPasswordLengthE] instead.
func ReadPasswordPolicyMinPasswordLength(iamClient *iam.Client) (int, error) {
	return ReadPasswordPolicyMinPasswordLengthContext(context.Background(), iamClient)
}
