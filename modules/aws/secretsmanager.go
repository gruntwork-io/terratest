package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// CreateSecretStringWithDefaultKey creates a new secret in Secrets Manager using the default "aws/secretsmanager" KMS key and returns the secret ARN
func CreateSecretStringWithDefaultKey(t testing.TestingT, awsRegion, description, name, secretString string) string {
	arn, err := CreateSecretStringWithDefaultKeyE(t, awsRegion, description, name, secretString)
	require.NoError(t, err)
	return arn
}

// CreateSecretStringWithDefaultKeyE creates a new secret in Secrets Manager using the default "aws/secretsmanager" KMS key and returns the secret ARN
func CreateSecretStringWithDefaultKeyE(t testing.TestingT, awsRegion, description, name, secretString string) (string, error) {
	logger.Logf(t, "Creating new secret in secrets manager named %s", name)

	client := NewSecretsManagerClient(t, awsRegion)

	secret, err := client.CreateSecret(&secretsmanager.CreateSecretInput{
		Description:  aws.String(description),
		Name:         aws.String(name),
		SecretString: aws.String(secretString),
	})

	if err != nil {
		return "", err
	}

	return aws.StringValue(secret.ARN), nil
}

// GetSecretValue takes the friendly name or ARN of a secret and returns the plaintext value
func GetSecretValue(t testing.TestingT, awsRegion, id string) string {
	secret, err := GetSecretValueE(t, awsRegion, id)
	require.NoError(t, err)
	return secret
}

// GetSecretValueE takes the friendly name or ARN of a secret and returns the plaintext value
func GetSecretValueE(t testing.TestingT, awsRegion, id string) (string, error) {
	logger.Logf(t, "Getting value of secret with ID %s", id)

	client := NewSecretsManagerClient(t, awsRegion)

	secret, err := client.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(id),
	})
	if err != nil {
		return "", err
	}

	return aws.StringValue(secret.SecretString), nil
}

// DeleteSecret deletes a secret. If forceDelete is true, the secret will be deleted after a short delay. If forceDelete is false, the secret will be deleted after a 30 day recovery window.
func DeleteSecret(t testing.TestingT, awsRegion, id string, forceDelete bool) {
	err := DeleteSecretE(t, awsRegion, id, forceDelete)
	require.NoError(t, err)
}

// DeleteSecretE deletes a secret. If forceDelete is true, the secret will be deleted after a short delay. If forceDelete is false, the secret will be deleted after a 30 day recovery window.
func DeleteSecretE(t testing.TestingT, awsRegion, id string, forceDelete bool) error {
	logger.Logf(t, "Deleting secret with ID %s", id)

	client := NewSecretsManagerClient(t, awsRegion)

	_, err := client.DeleteSecret(&secretsmanager.DeleteSecretInput{
		ForceDeleteWithoutRecovery: aws.Bool(forceDelete),
		SecretId:                   aws.String(id),
	})

	return err
}

// NewSecretsManagerClient creates a new SecretsManager client.
func NewSecretsManagerClient(t testing.TestingT, region string) *secretsmanager.SecretsManager {
	client, err := NewSecretsManagerClientE(t, region)
	require.NoError(t, err)
	return client
}

// NewSecretsManagerClientE creates a new SecretsManager client.
func NewSecretsManagerClientE(t testing.TestingT, region string) (*secretsmanager.SecretsManager, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return secretsmanager.New(sess), nil
}
