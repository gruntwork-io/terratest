package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

//PutSecret creates new secret of SecretManager.
func PutSecret(t testing.TestingT, awsRegion string, secretID string, secretDescription string, keyValue string) string {
	version, err := PutSecretE(t, awsRegion, secretID, secretDescription, keyValue)
	require.NoError(t, err)
	return version
}

//PutSecretE creates new secret of SecretManager.
func PutSecretE(t testing.TestingT, awsRegion string, secretID string, secretDescription string, secretValue string) (string, error) {
	secretManagerClient, err := NewSecretManagerClientE(t, awsRegion)
	if err != nil {
		return "", err
	}
	resp, err := secretManagerClient.CreateSecret(&secretsmanager.CreateSecretInput{Name: aws.String(secretID), Description: aws.String(secretDescription), SecretString: aws.String(secretValue)})
	if err != nil {
		return "", err
	}

	return *resp.VersionId, nil
}

//GetSecret creates new secret of SecretManager.
func GetSecret(t testing.TestingT, awsRegion string, secretID string) string {
	keyName, err := GetSecretE(t, awsRegion, secretID)
	require.NoError(t, err)
	return keyName
}

//GetSecretE creates new secret of SecretManager.
func GetSecretE(t testing.TestingT, awsRegion string, secretID string) (string, error) {
	secretManagerClient, err := NewSecretManagerClientE(t, awsRegion)
	if err != nil {
		return "", err
	}

	resp, err := secretManagerClient.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: aws.String(secretID)})

	if err != nil {
		return "", err
	}

	return *resp.Name, nil
}

//NewSecretManagerClient creates a new SQS client.
func NewSecretManagerClient(t testing.TestingT, region string) *secretsmanager.SecretsManager {
	client, err := NewSecretManagerClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

//NewSecretManagerClientE creates a new SQS client.
func NewSecretManagerClientE(t testing.TestingT, region string) (*secretsmanager.SecretsManager, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return secretsmanager.New(sess), nil
}
