package aws_test

import (
	"context"
	"strings"
	"testing"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	aws "github.com/gruntwork-io/terratest/modules/aws/v2"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
)

func TestEcrRepo(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegionContext(t, context.Background(), nil, nil)
	ecrRepoName := "terratest" + strings.ToLower(random.UniqueID())

	repo1, err := aws.CreateECRRepoContextE(t, context.Background(), region, ecrRepoName)
	defer aws.DeleteECRRepoContext(t, context.Background(), region, repo1)

	require.NoError(t, err)

	assert.Equal(t, ecrRepoName, awsSDK.ToString(repo1.RepositoryName))

	repo2, err := aws.GetECRRepoContextE(t, context.Background(), region, ecrRepoName)
	require.NoError(t, err)
	assert.Equal(t, ecrRepoName, awsSDK.ToString(repo2.RepositoryName))
}

func TestGetEcrRepoLifecyclePolicyError(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegionContext(t, context.Background(), nil, nil)
	ecrRepoName := "terratest" + strings.ToLower(random.UniqueID())

	repo1, err := aws.CreateECRRepoContextE(t, context.Background(), region, ecrRepoName)
	defer aws.DeleteECRRepoContext(t, context.Background(), region, repo1)

	require.NoError(t, err)

	assert.Equal(t, ecrRepoName, awsSDK.ToString(repo1.RepositoryName))

	_, err = aws.GetECRRepoLifecyclePolicyContextE(t, context.Background(), region, repo1)
	require.Error(t, err)
}

func TestCanSetECRRepoLifecyclePolicyWithSingleRule(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegionContext(t, context.Background(), nil, nil)
	ecrRepoName := "terratest" + strings.ToLower(random.UniqueID())

	repo1, err := aws.CreateECRRepoContextE(t, context.Background(), region, ecrRepoName)
	defer aws.DeleteECRRepoContext(t, context.Background(), region, repo1)

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

	err = aws.PutECRRepoLifecyclePolicyContextE(t, context.Background(), region, repo1, lifecyclePolicy)
	require.NoError(t, err)

	policy := aws.GetECRRepoLifecyclePolicyContext(t, context.Background(), region, repo1)
	assert.JSONEq(t, lifecyclePolicy, policy)
}

func TestCanSetRepositoryPolicyWithSimplePolicy(t *testing.T) {
	t.Parallel()

	region := aws.GetRandomStableRegionContext(t, context.Background(), nil, nil)
	ecrRepoName := "terratest" + strings.ToLower(random.UniqueID())

	repo, err := aws.CreateECRRepoContextE(t, context.Background(), region, ecrRepoName)
	defer aws.DeleteECRRepoContext(t, context.Background(), region, repo)

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

	err = aws.PutECRRepoPolicyContextE(t, context.Background(), region, repo, repositoryPolicy)
	require.NoError(t, err)

	policy := aws.GetECRRepoPolicyContext(t, context.Background(), region, repo)
	assert.JSONEq(t, repositoryPolicy, policy)
}
