package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// Ec2Keypair is an EC2 key pair.
type Ec2Keypair struct {
	*ssh.KeyPair
	Name   string // The name assigned in AWS to the EC2 Key Pair
	Region string // The AWS region where the EC2 Key Pair lives
}

// CreateAndImportEC2KeyPair generates a public/private KeyPair and import it into EC2 in the given region under the given name.
func CreateAndImportEC2KeyPair(t testing.TestingT, region string, name string) *Ec2Keypair {
	keyPair, err := CreateAndImportEC2KeyPairE(t, region, name)
	if err != nil {
		t.Fatal(err)
	}
	return keyPair
}

// CreateAndImportEC2KeyPairE generates a public/private KeyPair and import it into EC2 in the given region under the given name.
func CreateAndImportEC2KeyPairE(t testing.TestingT, region string, name string) (*Ec2Keypair, error) {
	keyPair, err := ssh.GenerateRSAKeyPairE(t, 2048)
	if err != nil {
		return nil, err
	}

	return ImportEC2KeyPairE(t, region, name, keyPair)
}

// ImportEC2KeyPair creates a Key Pair in EC2 by importing an existing public key.
func ImportEC2KeyPair(t testing.TestingT, region string, name string, keyPair *ssh.KeyPair) *Ec2Keypair {
	ec2KeyPair, err := ImportEC2KeyPairE(t, region, name, keyPair)
	if err != nil {
		t.Fatal(err)
	}
	return ec2KeyPair
}

// ImportEC2KeyPairE creates a Key Pair in EC2 by importing an existing public key.
func ImportEC2KeyPairE(t testing.TestingT, region string, name string, keyPair *ssh.KeyPair) (*Ec2Keypair, error) {
	logger.Logf(t, "Creating new Key Pair in EC2 region %s named %s", region, name)

	client, err := NewEc2ClientE(t, region)
	if err != nil {
		return nil, err
	}

	params := &ec2.ImportKeyPairInput{
		KeyName:           aws.String(name),
		PublicKeyMaterial: []byte(keyPair.PublicKey),
	}

	_, err = client.ImportKeyPair(params)
	if err != nil {
		return nil, err
	}

	return &Ec2Keypair{Name: name, Region: region, KeyPair: keyPair}, nil
}

// DeleteEC2KeyPair deletes an EC2 key pair.
func DeleteEC2KeyPair(t testing.TestingT, keyPair *Ec2Keypair) {
	err := DeleteEC2KeyPairE(t, keyPair)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteEC2KeyPairE deletes an EC2 key pair.
func DeleteEC2KeyPairE(t testing.TestingT, keyPair *Ec2Keypair) error {
	logger.Logf(t, "Deleting Key Pair in EC2 region %s named %s", keyPair.Region, keyPair.Name)

	client, err := NewEc2ClientE(t, keyPair.Region)
	if err != nil {
		return err
	}

	params := &ec2.DeleteKeyPairInput{
		KeyName: aws.String(keyPair.Name),
	}

	_, err = client.DeleteKeyPair(params)
	return err
}
