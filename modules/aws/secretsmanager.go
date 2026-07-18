package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// CreateSecretStringWithDefaultKeyContextE creates a new secret in Secrets Manager using the default "aws/secretsmanager" KMS key and returns the secret ARN.
// The ctx parameter supports cancellation and timeouts.
func CreateSecretStringWithDefaultKeyContextE(t testing.TestingT, ctx context.Context, awsRegion, description, name, secretString string) (string, error) {
	logger.Default.Logf(t, "Creating new secret in secrets manager named %s", name)

	client, err := NewSecretsManagerClientContextE(t, ctx, awsRegion)
	if err != nil {
		return "", err
	}

	secret, err := client.CreateSecret(ctx, &secretsmanager.CreateSecretInput{
		Description:  aws.String(description),
		Name:         aws.String(name),
		SecretString: aws.String(secretString),
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(secret.ARN), nil
}

// CreateSecretStringWithDefaultKeyContext creates a new secret in Secrets Manager using the default "aws/secretsmanager" KMS key and returns the secret ARN.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CreateSecretStringWithDefaultKeyContext(t testing.TestingT, ctx context.Context, awsRegion, description, name, secretString string) string {
	t.Helper()
	arn, err := CreateSecretStringWithDefaultKeyContextE(t, ctx, awsRegion, description, name, secretString)
	require.NoError(t, err)

	return arn
}

// GetSecretValueContextE takes the friendly name or ARN of a secret and returns the plaintext value.
// The ctx parameter supports cancellation and timeouts.
func GetSecretValueContextE(t testing.TestingT, ctx context.Context, awsRegion, id string) (string, error) {
	logger.Default.Logf(t, "Getting value of secret with ID %s", id)

	client, err := NewSecretsManagerClientContextE(t, ctx, awsRegion)
	if err != nil {
		return "", err
	}

	secret, err := client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(id),
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(secret.SecretString), nil
}

// GetSecretValueContext takes the friendly name or ARN of a secret and returns the plaintext value.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetSecretValueContext(t testing.TestingT, ctx context.Context, awsRegion, id string) string {
	t.Helper()
	secret, err := GetSecretValueContextE(t, ctx, awsRegion, id)
	require.NoError(t, err)

	return secret
}

// PutSecretStringContextE updates a secret in Secrets Manager to a new string value.
// The ctx parameter supports cancellation and timeouts.
func PutSecretStringContextE(t testing.TestingT, ctx context.Context, awsRegion, id string, secretString string) error {
	logger.Default.Logf(t, "Updating secret with ID %s", id)

	client, err := NewSecretsManagerClientContextE(t, ctx, awsRegion)
	if err != nil {
		return err
	}

	_, err = client.PutSecretValue(ctx, &secretsmanager.PutSecretValueInput{
		SecretId:     aws.String(id),
		SecretString: aws.String(secretString),
	})

	return err
}

// PutSecretStringContext updates a secret in Secrets Manager to a new string value.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func PutSecretStringContext(t testing.TestingT, ctx context.Context, awsRegion, id string, secretString string) {
	t.Helper()
	err := PutSecretStringContextE(t, ctx, awsRegion, id, secretString)
	require.NoError(t, err)
}

// DeleteSecretContextE deletes a secret. If forceDelete is true, the secret will be deleted after a short delay. If forceDelete is false, the secret will be deleted after a 30-day recovery window.
// The ctx parameter supports cancellation and timeouts.
func DeleteSecretContextE(t testing.TestingT, ctx context.Context, awsRegion, id string, forceDelete bool) error {
	logger.Default.Logf(t, "Deleting secret with ID %s", id)

	client, err := NewSecretsManagerClientContextE(t, ctx, awsRegion)
	if err != nil {
		return err
	}

	_, err = client.DeleteSecret(ctx, &secretsmanager.DeleteSecretInput{
		ForceDeleteWithoutRecovery: aws.Bool(forceDelete),
		SecretId:                   aws.String(id),
	})

	return err
}

// DeleteSecretContext deletes a secret. If forceDelete is true, the secret will be deleted after a short delay. If forceDelete is false, the secret will be deleted after a 30-day recovery window.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteSecretContext(t testing.TestingT, ctx context.Context, awsRegion, id string, forceDelete bool) {
	t.Helper()
	err := DeleteSecretContextE(t, ctx, awsRegion, id, forceDelete)
	require.NoError(t, err)
}

// NewSecretsManagerClientContextE creates a new SecretsManager client.
// The ctx parameter supports cancellation and timeouts.
func NewSecretsManagerClientContextE(t testing.TestingT, ctx context.Context, region string) (*secretsmanager.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return secretsmanager.NewFromConfig(*sess), nil
}

// NewSecretsManagerClientContext creates a new SecretsManager client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewSecretsManagerClientContext(t testing.TestingT, ctx context.Context, region string) *secretsmanager.Client {
	t.Helper()
	client, err := NewSecretsManagerClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}
