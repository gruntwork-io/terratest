package aws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEcrPublicRepo(t *testing.T) {
	t.Parallel()

	region := "us-east-1"
	ecrRepoName := fmt.Sprintf("terratest%s", strings.ToLower(random.UniqueId()))
	repo1, err := CreateECRPublicRepoE(t, region, ecrRepoName)
	defer DeleteECRPublicRepo(t, region, repo1)
	require.NoError(t, err)

	assert.Equal(t, ecrRepoName, aws.StringValue(repo1.RepositoryName))

	repo2, err := GetECRPublicRepoE(t, region, ecrRepoName)
	require.NoError(t, err)
	assert.Equal(t, ecrRepoName, aws.StringValue(repo2.RepositoryName))
}
