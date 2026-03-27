package aws_test

import (
	"strings"
	"testing"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
)

func TestEcrRepo(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegion(t, nil, nil)
	ecrRepoName := "terratest" + strings.ToLower(random.UniqueID())

	repo1, err := aws.CreateECRRepoE(t, region, ecrRepoName)
	defer aws.DeleteECRRepo(t, region, repo1)

	require.NoError(t, err)

	assert.Equal(t, ecrRepoName, awsSDK.ToString(repo1.RepositoryName))

	repo2, err := aws.GetECRRepoE(t, region, ecrRepoName)
	require.NoError(t, err)
	assert.Equal(t, ecrRepoName, awsSDK.ToString(repo2.RepositoryName))
}

func TestGetEcrRepoLifecyclePolicyError(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegion(t, nil, nil)
	ecrRepoName := "terratest" + strings.ToLower(random.UniqueID())

	repo1, err := aws.CreateECRRepoE(t, region, ecrRepoName)
	defer aws.DeleteECRRepo(t, region, repo1)

	require.NoError(t, err)

	assert.Equal(t, ecrRepoName, awsSDK.ToString(repo1.RepositoryName))

	_, err = aws.GetECRRepoLifecyclePolicyE(t, region, repo1)
	require.Error(t, err)
}

func TestCanSetECRRepoLifecyclePolicyWithSingleRule(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegion(t, nil, nil)
	ecrRepoName := "terratest" + strings.ToLower(random.UniqueID())

	repo1, err := aws.CreateECRRepoE(t, region, ecrRepoName)
	defer aws.DeleteECRRepo(t, region, repo1)

	require.NoError(t, err)

	lifecyclePolicy := `{
		"rules": [
			{
				"rulePriority": 1,
				"description": "Expire images older than 14 days",
				"selection": {
					"tagStatus": "untagged",
					"countType": "sinceImagePushed",
					"countUnit": "days",
					"countNumber": 14
				},
				"action": {
					"type": "expire"
				}
			}
		]
	}`

	err = aws.PutECRRepoLifecyclePolicyE(t, region, repo1, lifecyclePolicy)
	require.NoError(t, err)

	policy := aws.GetECRRepoLifecyclePolicy(t, region, repo1)
	assert.JSONEq(t, lifecyclePolicy, policy)
}

func TestCanSetRepositoryPolicyWithSimplePolicy(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegion(t, nil, nil)
	ecrRepoName := "terratest" + strings.ToLower(random.UniqueID())

	repo, err := aws.CreateECRRepoE(t, region, ecrRepoName)
	defer aws.DeleteECRRepo(t, region, repo)

	require.NoError(t, err)

	repositoryPolicy := `
		{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Sid": "AllowPushPull",
				"Effect": "Allow",
				"Principal": {
					"AWS": "*"
				},
				"Action": "ecr:*"
			}
		]
	}`

	err = aws.PutECRRepoPolicyE(t, region, repo, repositoryPolicy)
	require.NoError(t, err)

	policy := aws.GetECRRepoPolicy(t, region, repo)
	assert.JSONEq(t, repositoryPolicy, policy)
}
