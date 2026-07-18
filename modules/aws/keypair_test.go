package aws_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	aws "github.com/gruntwork-io/terratest/modules/aws/v2"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/stretchr/testify/assert"
)

func TestCreateImportAndDeleteEC2KeyPair(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegionContext(t, context.Background(), nil, nil)
	uniqueID := random.UniqueID()
	name := "test-key-pair-" + uniqueID

	keyPair := aws.CreateAndImportEC2KeyPairContext(t, context.Background(), region, name)
	defer deleteKeyPair(t, keyPair)

	assert.True(t, keyPairExists(t, keyPair))
	assert.Equal(t, name, keyPair.Name)
	assert.Equal(t, region, keyPair.Region)
	assert.Contains(t, keyPair.PublicKey, "ssh-rsa")
	assert.Contains(t, keyPair.PrivateKey, "-----BEGIN RSA PRIVATE KEY-----")
}

func keyPairExists(t *testing.T, keyPair *aws.Ec2Keypair) bool {
	t.Helper()

	client := aws.NewEc2ClientContext(t, context.Background(), keyPair.Region)

	input := ec2.DescribeKeyPairsInput{
		KeyNames: []string{keyPair.Name},
	}

	out, err := client.DescribeKeyPairs(context.Background(), &input)
	if err != nil {
		if strings.Contains(err.Error(), "InvalidKeyPair.NotFound") {
			return false
		}

		t.Fatal(err)
	}

	return len(out.KeyPairs) == 1
}

func deleteKeyPair(t *testing.T, keyPair *aws.Ec2Keypair) {
	t.Helper()

	aws.DeleteEC2KeyPairContext(t, context.Background(), keyPair)
	assert.False(t, keyPairExists(t, keyPair))
}
