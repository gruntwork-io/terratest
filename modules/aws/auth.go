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

// NewAuthConfigWithCredsE creates a new AWS Config using explicit static credentials. Useful for
// authenticating as a dynamically-created IAM User. No context variant is provided because this
// function performs no I/O — it only constructs an aws.Config.
func NewAuthConfigWithCredsE(t testing.TestingT, region string, accessKeyID string, secretAccessKey string) (*aws.Config, error) {
	return &aws.Config{
		Region:      region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	}, nil
}

// NewAuthConfigWithCreds creates a new AWS Config using explicit static credentials.
// This function will fail the test if there is an error.
func NewAuthConfigWithCreds(t testing.TestingT, region string, accessKeyID string, secretAccessKey string) *aws.Config {
	t.Helper()
	cfg, err := NewAuthConfigWithCredsE(t, region, accessKeyID, secretAccessKey)
	require.NoError(t, err)
	return cfg
}

// NewAuthConfigWithMfaContextE creates a new AWS Config authenticated with an MFA session token
// obtained via the given STS client and MFA device.
// The ctx parameter supports cancellation and timeouts.
func NewAuthConfigWithMfaContextE(t testing.TestingT, ctx context.Context, region string, stsClient *sts.Client, mfaDevice *types.VirtualMFADevice) (*aws.Config, error) {
	tokenCode, err := GenerateMfaTokenE(t, mfaDevice)
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

// NewAuthConfigWithMfaContext creates a new AWS Config authenticated with an MFA session token.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewAuthConfigWithMfaContext(t testing.TestingT, ctx context.Context, region string, stsClient *sts.Client, mfaDevice *types.VirtualMFADevice) *aws.Config {
	t.Helper()
	cfg, err := NewAuthConfigWithMfaContextE(t, ctx, region, stsClient, mfaDevice)
	require.NoError(t, err)
	return cfg
}

// NewAuthConfigWithMfaE creates a new AWS Config authenticated with an MFA session token.
func NewAuthConfigWithMfaE(t testing.TestingT, region string, stsClient *sts.Client, mfaDevice *types.VirtualMFADevice) (*aws.Config, error) {
	return NewAuthConfigWithMfaContextE(t, context.Background(), region, stsClient, mfaDevice)
}

// NewAuthConfigWithMfa creates a new AWS Config authenticated with an MFA session token.
// This function will fail the test if there is an error.
func NewAuthConfigWithMfa(t testing.TestingT, region string, stsClient *sts.Client, mfaDevice *types.VirtualMFADevice) *aws.Config {
	t.Helper()
	return NewAuthConfigWithMfaContext(t, context.Background(), region, stsClient, mfaDevice)
}

// GenerateMfaTokenE returns a time-based one-time password for the given MFA device, per RFC 6238.
// The returned value changes every 30 seconds.
func GenerateMfaTokenE(t testing.TestingT, mfaDevice *types.VirtualMFADevice) (string, error) {
	base32StringSeed := string(mfaDevice.Base32StringSeed)

	otp, err := totp.GenerateCode(base32StringSeed, time.Now())
	if err != nil {
		return "", err
	}

	return otp, nil
}

// GenerateMfaToken returns a time-based one-time password for the given MFA device, per RFC 6238.
// This function will fail the test if there is an error.
func GenerateMfaToken(t testing.TestingT, mfaDevice *types.VirtualMFADevice) string {
	t.Helper()
	token, err := GenerateMfaTokenE(t, mfaDevice)
	require.NoError(t, err)
	return token
}

// GetPasswordPolicyMinLengthContextE returns the minimum password length from the account's
// IAM password policy.
// The ctx parameter supports cancellation and timeouts.
func GetPasswordPolicyMinLengthContextE(t testing.TestingT, ctx context.Context, iamClient *iam.Client) (int, error) {
	output, err := iamClient.GetAccountPasswordPolicy(ctx, &iam.GetAccountPasswordPolicyInput{})
	if err != nil {
		return -1, err
	}

	return int(*output.PasswordPolicy.MinimumPasswordLength), nil
}

// GetPasswordPolicyMinLengthContext returns the minimum password length from the account's IAM password policy.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetPasswordPolicyMinLengthContext(t testing.TestingT, ctx context.Context, iamClient *iam.Client) int {
	t.Helper()
	n, err := GetPasswordPolicyMinLengthContextE(t, ctx, iamClient)
	require.NoError(t, err)
	return n
}

// GetPasswordPolicyMinLengthE returns the minimum password length from the account's IAM password policy.
func GetPasswordPolicyMinLengthE(t testing.TestingT, iamClient *iam.Client) (int, error) {
	return GetPasswordPolicyMinLengthContextE(t, context.Background(), iamClient)
}

// GetPasswordPolicyMinLength returns the minimum password length from the account's IAM password policy.
// This function will fail the test if there is an error.
func GetPasswordPolicyMinLength(t testing.TestingT, iamClient *iam.Client) int {
	t.Helper()
	return GetPasswordPolicyMinLengthContext(t, context.Background(), iamClient)
}

// CredentialsError is an error that occurs because AWS credentials can't be found.
type CredentialsError struct {
	UnderlyingErr error
}

func (err CredentialsError) Error() string {
	return fmt.Sprintf("Error finding AWS credentials. Did you set the AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY environment variables or configure an AWS profile? Underlying error: %v", err.UnderlyingErr)
}

// -----------------------------------------------------------------------------
// Deprecated: legacy surface preserved for backwards compatibility.
//
// These functions predate the package's testing.TestingT + Context[E] convention.
// New callers should use the NewAuthConfig* / GenerateMfaToken / GetPasswordPolicyMinLength
// families above. The shims below will be removed in a future major release.
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
// Deprecated: Use [NewAuthConfigWithCredsE] instead. The Context variant was a placeholder
// for API symmetry — the underlying operation performs no I/O, so no context is needed.
func CreateAwsSessionWithCredsContext(ctx context.Context, region string, accessKeyID string, secretAccessKey string) (*aws.Config, error) {
	return &aws.Config{
		Region:      region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	}, nil
}

// CreateAwsSessionWithCreds creates a new AWS Config using explicit credentials.
//
// Deprecated: Use [NewAuthConfigWithCredsE] instead.
func CreateAwsSessionWithCreds(region string, accessKeyID string, secretAccessKey string) (*aws.Config, error) {
	return CreateAwsSessionWithCredsContext(context.Background(), region, accessKeyID, secretAccessKey)
}

// CreateAwsSessionWithMfaContext creates a new AWS Config authenticated using an MFA token retrieved using the given STS client and MFA Device.
//
// Deprecated: Use [NewAuthConfigWithMfaContextE] instead.
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
// Deprecated: Use [NewAuthConfigWithMfaE] instead.
func CreateAwsSessionWithMfa(region string, stsClient *sts.Client, mfaDevice *types.VirtualMFADevice) (*aws.Config, error) {
	return CreateAwsSessionWithMfaContext(context.Background(), region, stsClient, mfaDevice)
}

// GetTimeBasedOneTimePassword gets a One-Time Password from the given mfaDevice. Per the RFC 6238 standard, this value will be different every 30 seconds.
//
// Deprecated: Use [GenerateMfaTokenE] instead.
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
// Deprecated: Use [GetPasswordPolicyMinLengthContextE] instead.
func ReadPasswordPolicyMinPasswordLengthContext(ctx context.Context, iamClient *iam.Client) (int, error) {
	output, err := iamClient.GetAccountPasswordPolicy(ctx, &iam.GetAccountPasswordPolicyInput{})
	if err != nil {
		return -1, err
	}

	return int(*output.PasswordPolicy.MinimumPasswordLength), nil
}

// ReadPasswordPolicyMinPasswordLength returns the minimal password length.
//
// Deprecated: Use [GetPasswordPolicyMinLengthE] instead.
func ReadPasswordPolicyMinPasswordLength(iamClient *iam.Client) (int, error) {
	return ReadPasswordPolicyMinPasswordLengthContext(context.Background(), iamClient)
}
