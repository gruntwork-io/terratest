package aws_test

import (
	"context"
	"errors"
	"testing"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/stretchr/testify/require"

	aws "github.com/gruntwork-io/terratest/modules/aws"
)

// mockKmsClient is a test double for aws.KmsAPI that returns canned responses.
type mockKmsClient struct {
	DescribeKeyOutput *kms.DescribeKeyOutput
	DescribeKeyErr    error
	lastKeyID         string
}

func (m *mockKmsClient) DescribeKey(_ context.Context, params *kms.DescribeKeyInput, _ ...func(*kms.Options)) (*kms.DescribeKeyOutput, error) {
	m.lastKeyID = awsSDK.ToString(params.KeyId)

	if m.DescribeKeyErr != nil {
		return nil, m.DescribeKeyErr
	}

	return m.DescribeKeyOutput, nil
}

func TestGetCmkArnWithClientContextE(t *testing.T) {
	t.Parallel()

	const (
		keyArn = "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
		alias  = "alias/my-cmk"
	)

	t.Run("returns arn for key id", func(t *testing.T) {
		t.Parallel()

		client := &mockKmsClient{
			DescribeKeyOutput: &kms.DescribeKeyOutput{
				KeyMetadata: &types.KeyMetadata{Arn: awsSDK.String(keyArn)},
			},
		}

		got, err := aws.GetCmkArnWithClientContextE(t, context.Background(), client, alias)
		require.NoError(t, err)
		require.Equal(t, keyArn, got)
		require.Equal(t, alias, client.lastKeyID)
	})

	t.Run("propagates api error", func(t *testing.T) {
		t.Parallel()

		client := &mockKmsClient{DescribeKeyErr: errors.New("NotFoundException")}

		_, err := aws.GetCmkArnWithClientContextE(t, context.Background(), client, "alias/missing")
		require.Error(t, err)
	})
}
