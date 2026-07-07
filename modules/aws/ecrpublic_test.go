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

// ecrPublicRegion is fixed because the ECR Public API is only available in us-east-1.
const ecrPublicRegion = "us-east-1"

func TestEcrPublicRepo(t *testing.T) {
	t.Parallel()

	ecrRepoName := "terratest" + strings.ToLower(random.UniqueID())

	repo1, err := aws.CreateECRPublicRepoE(t, ecrPublicRegion, ecrRepoName)
	defer aws.DeleteECRPublicRepo(t, ecrPublicRegion, repo1)

	require.NoError(t, err)

	assert.Equal(t, ecrRepoName, awsSDK.ToString(repo1.RepositoryName))

	repo2, err := aws.GetECRPublicRepoE(t, ecrPublicRegion, ecrRepoName)
	require.NoError(t, err)
	assert.Equal(t, ecrRepoName, awsSDK.ToString(repo2.RepositoryName))
}

func TestGetEcrPublicRepoError(t *testing.T) {
	t.Parallel()

	_, err := aws.GetECRPublicRepoE(t, ecrPublicRegion, "terratest"+strings.ToLower(random.UniqueID()))
	require.Error(t, err)
}
