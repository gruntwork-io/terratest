package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

const rsaKeyBits = 2048

// Ec2Keypair is an EC2 key pair.
type Ec2Keypair struct {
	*ssh.KeyPair
	Name   string // The name assigned in AWS to the EC2 Key Pair
	Region string // The AWS region where the EC2 Key Pair lives
}

// CreateAndImportEC2KeyPairContextE generates a public/private KeyPair and import it into EC2 in the given region under the given name.
// The ctx parameter supports cancellation and timeouts.
func CreateAndImportEC2KeyPairContextE(t testing.TestingT, ctx context.Context, region string, name string) (*Ec2Keypair, error) {
	keyPair, err := ssh.GenerateRSAKeyPairE(t, rsaKeyBits)
	if err != nil {
		return nil, err
	}

	return ImportEC2KeyPairContextE(t, ctx, region, name, keyPair)
}

// CreateAndImportEC2KeyPairContext generates a public/private KeyPair and import it into EC2 in the given region under the given name.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CreateAndImportEC2KeyPairContext(t testing.TestingT, ctx context.Context, region string, name string) *Ec2Keypair {
	t.Helper()

	keyPair, err := CreateAndImportEC2KeyPairContextE(t, ctx, region, name)
	require.NoError(t, err)

	return keyPair
}

// CreateAndImportEC2KeyPair generates a public/private KeyPair and import it into EC2 in the given region under the given name.
//
// Deprecated: Use [CreateAndImportEC2KeyPairContext] instead.
func CreateAndImportEC2KeyPair(t testing.TestingT, region string, name string) *Ec2Keypair {
	t.Helper()

	return CreateAndImportEC2KeyPairContext(t, context.Background(), region, name)
}

// CreateAndImportEC2KeyPairE generates a public/private KeyPair and import it into EC2 in the given region under the given name.
//
// Deprecated: Use [CreateAndImportEC2KeyPairContextE] instead.
func CreateAndImportEC2KeyPairE(t testing.TestingT, region string, name string) (*Ec2Keypair, error) {
	return CreateAndImportEC2KeyPairContextE(t, context.Background(), region, name)
}

// ImportEC2KeyPairContextE creates a Key Pair in EC2 by importing an existing public key.
// The ctx parameter supports cancellation and timeouts.
func ImportEC2KeyPairContextE(t testing.TestingT, ctx context.Context, region string, name string, keyPair *ssh.KeyPair) (*Ec2Keypair, error) {
	logger.Default.Logf(t, "Creating new Key Pair in EC2 region %s named %s", region, name)

	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	client := ec2.NewFromConfig(*sess)

	params := &ec2.ImportKeyPairInput{
		KeyName:           aws.String(name),
		PublicKeyMaterial: []byte(keyPair.PublicKey),
	}

	_, err = client.ImportKeyPair(ctx, params)
	if err != nil {
		return nil, err
	}

	return &Ec2Keypair{Name: name, Region: region, KeyPair: keyPair}, nil
}

// ImportEC2KeyPairContext creates a Key Pair in EC2 by importing an existing public key.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func ImportEC2KeyPairContext(t testing.TestingT, ctx context.Context, region string, name string, keyPair *ssh.KeyPair) *Ec2Keypair {
	t.Helper()

	ec2KeyPair, err := ImportEC2KeyPairContextE(t, ctx, region, name, keyPair)
	require.NoError(t, err)

	return ec2KeyPair
}

// ImportEC2KeyPair creates a Key Pair in EC2 by importing an existing public key.
//
// Deprecated: Use [ImportEC2KeyPairContext] instead.
func ImportEC2KeyPair(t testing.TestingT, region string, name string, keyPair *ssh.KeyPair) *Ec2Keypair {
	t.Helper()

	return ImportEC2KeyPairContext(t, context.Background(), region, name, keyPair)
}

// ImportEC2KeyPairE creates a Key Pair in EC2 by importing an existing public key.
//
// Deprecated: Use [ImportEC2KeyPairContextE] instead.
func ImportEC2KeyPairE(t testing.TestingT, region string, name string, keyPair *ssh.KeyPair) (*Ec2Keypair, error) {
	return ImportEC2KeyPairContextE(t, context.Background(), region, name, keyPair)
}

// DeleteEC2KeyPairContextE deletes an EC2 key pair.
// The ctx parameter supports cancellation and timeouts.
func DeleteEC2KeyPairContextE(t testing.TestingT, ctx context.Context, keyPair *Ec2Keypair) error {
	logger.Default.Logf(t, "Deleting Key Pair in EC2 region %s named %s", keyPair.Region, keyPair.Name)

	sess, err := NewAuthenticatedSessionContext(ctx, keyPair.Region)
	if err != nil {
		return err
	}

	client := ec2.NewFromConfig(*sess)

	params := &ec2.DeleteKeyPairInput{
		KeyName: aws.String(keyPair.Name),
	}

	_, err = client.DeleteKeyPair(ctx, params)

	return err
}

// DeleteEC2KeyPairContext deletes an EC2 key pair.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteEC2KeyPairContext(t testing.TestingT, ctx context.Context, keyPair *Ec2Keypair) {
	t.Helper()

	err := DeleteEC2KeyPairContextE(t, ctx, keyPair)
	require.NoError(t, err)
}

// DeleteEC2KeyPair deletes an EC2 key pair.
//
// Deprecated: Use [DeleteEC2KeyPairContext] instead.
func DeleteEC2KeyPair(t testing.TestingT, keyPair *Ec2Keypair) {
	t.Helper()

	DeleteEC2KeyPairContext(t, context.Background(), keyPair)
}

// DeleteEC2KeyPairE deletes an EC2 key pair.
//
// Deprecated: Use [DeleteEC2KeyPairContextE] instead.
func DeleteEC2KeyPairE(t testing.TestingT, keyPair *Ec2Keypair) error {
	return DeleteEC2KeyPairContextE(t, context.Background(), keyPair)
}
