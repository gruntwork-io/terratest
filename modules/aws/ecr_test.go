package aws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEcrRepo(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	ecrRepoName := fmt.Sprintf("terratest%s", strings.ToLower(random.UniqueId()))
	repo1, err := CreateECRRepoE(t, region, ecrRepoName)
	defer DeleteECRRepo(t, region, repo1)
	require.NoError(t, err)

	assert.Equal(t, ecrRepoName, aws.ToString(repo1.RepositoryName))

	repo2, err := GetECRRepoE(t, region, ecrRepoName)
	require.NoError(t, err)
	assert.Equal(t, ecrRepoName, aws.ToString(repo2.RepositoryName))
}

func TestGetEcrRepoLifecyclePolicyError(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	ecrRepoName := fmt.Sprintf("terratest%s", strings.ToLower(random.UniqueId()))
	repo1, err := CreateECRRepoE(t, region, ecrRepoName)
	defer DeleteECRRepo(t, region, repo1)
	require.NoError(t, err)

	assert.Equal(t, ecrRepoName, aws.ToString(repo1.RepositoryName))

	_, err = GetECRRepoLifecyclePolicyE(t, region, repo1)
	require.Error(t, err)
}

func TestCanSetECRRepoLifecyclePolicyWithSingleRule(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	ecrRepoName := fmt.Sprintf("terratest%s", strings.ToLower(random.UniqueId()))
	repo1, err := CreateECRRepoE(t, region, ecrRepoName)
	defer DeleteECRRepo(t, region, repo1)
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

	err = PutECRRepoLifecyclePolicyE(t, region, repo1, lifecyclePolicy)
	require.NoError(t, err)

	policy := GetECRRepoLifecyclePolicy(t, region, repo1)
	assert.JSONEq(t, lifecyclePolicy, policy)
}

func TestCanSetRepositoryPolicyWithSimplePolicy(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	ecrRepoName := fmt.Sprintf("terratest%s", strings.ToLower(random.UniqueId()))
	repo, err := CreateECRRepoE(t, region, ecrRepoName)
	defer DeleteECRRepo(t, region, repo)
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

	err = PutECRRepoPolicyE(t, region, repo, repositoryPolicy)
	require.NoError(t, err)

	policy := GetECRRepoPolicy(t, region, repo)
	assert.JSONEq(t, repositoryPolicy, policy)
}
