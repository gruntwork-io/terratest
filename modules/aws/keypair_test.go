package aws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func TestCreateImportAndDeleteEC2KeyPair(t *testing.T) {
	t.Parallel()

	region := GetRandomRegion(t, nil, nil)
	uniqueID := random.UniqueId()
	name := fmt.Sprintf("test-key-pair-%s", uniqueID)

	keyPair := CreateAndImportEC2KeyPair(t, region, name)
	defer deleteKeyPair(t, keyPair)

	assert.True(t, keyPairExists(t, keyPair))
	assert.Equal(t, name, keyPair.Name)
	assert.Equal(t, region, keyPair.Region)
	assert.Contains(t, keyPair.PublicKey, "ssh-rsa")
	assert.Contains(t, keyPair.PrivateKey, "-----BEGIN RSA PRIVATE KEY-----")
}

func keyPairExists(t *testing.T, keyPair *Ec2Keypair) bool {
	client := NewEc2Client(t, keyPair.Region)

	input := ec2.DescribeKeyPairsInput{
		KeyNames: aws.StringSlice([]string{keyPair.Name}),
	}

	out, err := client.DescribeKeyPairs(&input)
	if err != nil {
		if strings.Contains(err.Error(), "InvalidKeyPair.NotFound") {
			return false
		}
		t.Fatal(err)
	}

	return len(out.KeyPairs) == 1
}

func deleteKeyPair(t *testing.T, keyPair *Ec2Keypair) {
	DeleteEC2KeyPair(t, keyPair)
	assert.False(t, keyPairExists(t, keyPair))
}
