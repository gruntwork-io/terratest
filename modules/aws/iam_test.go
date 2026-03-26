package aws_test

import (
	"context"
	"strings"
	"testing"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	aws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIamCurrentUserName(t *testing.T) {
	t.Parallel()

	username := aws.GetIamCurrentUserName(t)
	assert.NotEmpty(t, username)
}

func TestGetIamCurrentUserArn(t *testing.T) {
	t.Parallel()

	username := aws.GetIamCurrentUserArn(t)
	assert.Regexp(t, "^arn:aws:iam::[0-9]{12}:user/.+$", username)
}

func TestGetIAMPolicyDocument(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomRegion(t, nil, nil)

	t.Run("Exists", func(t *testing.T) {
		t.Parallel()

		iamClient, err := aws.NewIamClientE(t, region)
		require.NoError(t, err)

		policyDocument := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Sid": "Stmt1530709892083",
					"Action": "*",
					"Effect": "Allow",
					"Resource": "*"
				}
			]
		}`
		input := &iam.CreatePolicyInput{
			PolicyName:     awsSDK.String(strings.ToLower(random.UniqueID())),
			PolicyDocument: awsSDK.String(policyDocument),
		}

		policy, err := iamClient.CreatePolicy(context.Background(), input)
		require.NoError(t, err)

		t.Cleanup(func() {
			t.Log("Deleting IAM Policy Document")

			_, err := iamClient.DeletePolicy(context.Background(), &iam.DeletePolicyInput{
				PolicyArn: policy.Policy.Arn,
			})
			require.NoError(t, err)
		})

		p := aws.GetIamPolicyDocument(t, region, *policy.Policy.Arn)
		t.Log("Retrieved Policy Document:", p)
		assert.JSONEq(t, policyDocument, p)
	})

	t.Run("DoesNotExist", func(t *testing.T) {
		t.Parallel()

		_, err := aws.GetIamPolicyDocumentE(t, region, "arn:aws:iam::1234567890:policy/does-not-exist")
		require.Error(t, err)
	})
}
