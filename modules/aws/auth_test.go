package aws_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	aws "github.com/gruntwork-io/terratest/modules/aws/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const authTestRegion = "us-east-1"

const assumeRoleResponse = `<AssumeRoleResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <AssumeRoleResult>
    <Credentials>
      <AccessKeyId>ASIAASSUMED</AccessKeyId>
      <SecretAccessKey>assumed-secret</SecretAccessKey>
      <SessionToken>assumed-token</SessionToken>
      <Expiration>2030-01-01T00:00:00Z</Expiration>
    </Credentials>
    <AssumedRoleUser>
      <Arn>arn:aws:sts::123456789012:assumed-role/test/session</Arn>
      <AssumedRoleId>AROATEST:session</AssumedRoleId>
    </AssumedRoleUser>
  </AssumeRoleResult>
</AssumeRoleResponse>`

const getSessionTokenResponse = `<GetSessionTokenResponse xmlns="https://sts.amazonaws.com/doc/2011-06-15/">
  <GetSessionTokenResult>
    <Credentials>
      <AccessKeyId>ASIAMFA</AccessKeyId>
      <SecretAccessKey>mfa-secret</SecretAccessKey>
      <SessionToken>mfa-token</SessionToken>
      <Expiration>2030-01-01T00:00:00Z</Expiration>
    </Credentials>
  </GetSessionTokenResult>
</GetSessionTokenResponse>`

// stubSts serves a canned STS response so the auth helpers can be exercised without
// reaching AWS, the way they would be pointed at an emulator such as LocalStack.
func stubSts(t *testing.T, body string) string {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/xml")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)

	return server.URL
}

// A caller pointing the SDK at an emulator via AWS_ENDPOINT_URL expects every helper that
// hands back a Config to keep that endpoint, otherwise the requests silently go to real AWS.
func TestCreateAwsSessionWithCredsContextKeepsEndpointOverride(t *testing.T) {
	endpoint := "http://localhost:4566"
	t.Setenv("AWS_ENDPOINT_URL", endpoint)

	cfg, err := aws.CreateAwsSessionWithCredsContext(context.Background(), authTestRegion, "key", "secret")
	require.NoError(t, err)

	require.NotNil(t, cfg.BaseEndpoint, "BaseEndpoint was dropped, so requests would go to real AWS")
	assert.Equal(t, endpoint, *cfg.BaseEndpoint)

	// The explicitly supplied credentials must still win.
	creds, err := cfg.Credentials.Retrieve(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "key", creds.AccessKeyID)
	assert.Equal(t, "secret", creds.SecretAccessKey)
}

func TestNewAuthenticatedSessionFromRoleContextKeepsEndpointOverride(t *testing.T) {
	endpoint := stubSts(t, assumeRoleResponse)

	t.Setenv("AWS_ENDPOINT_URL", endpoint)
	t.Setenv("AWS_ACCESS_KEY_ID", "test")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	cfg, err := aws.NewAuthenticatedSessionFromRoleContext(
		context.Background(), authTestRegion, "arn:aws:iam::123456789012:role/test")
	require.NoError(t, err)

	require.NotNil(t, cfg.BaseEndpoint, "BaseEndpoint was dropped, so requests would go to real AWS")
	assert.Equal(t, endpoint, *cfg.BaseEndpoint)

	// The assumed-role credentials must be the ones returned.
	creds, err := cfg.Credentials.Retrieve(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "ASIAASSUMED", creds.AccessKeyID)
}

func TestCreateAwsSessionWithMfaContextKeepsEndpointOverride(t *testing.T) {
	endpoint := stubSts(t, getSessionTokenResponse)

	t.Setenv("AWS_ENDPOINT_URL", endpoint)

	stsClient := sts.New(sts.Options{
		Region:       authTestRegion,
		BaseEndpoint: awsSDK.String(endpoint),
		Credentials:  credentials.NewStaticCredentialsProvider("key", "secret", ""),
	})

	mfaDevice := &types.VirtualMFADevice{
		SerialNumber: awsSDK.String("arn:aws:iam::123456789012:mfa/test"),
		// Valid base32, required by the TOTP generator.
		Base32StringSeed: []byte("JBSWY3DPEHPK3PXP"),
	}

	cfg, err := aws.CreateAwsSessionWithMfaContext(context.Background(), authTestRegion, stsClient, mfaDevice)
	require.NoError(t, err)

	require.NotNil(t, cfg.BaseEndpoint, "BaseEndpoint was dropped, so requests would go to real AWS")
	assert.Equal(t, endpoint, *cfg.BaseEndpoint)

	creds, err := cfg.Credentials.Retrieve(context.Background())
	require.NoError(t, err)
	assert.Equal(t, "ASIAMFA", creds.AccessKeyID)
}
