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

func TestEcrRepo(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	ecrRepoName := fmt.Sprintf("terratest%s", strings.ToLower(random.UniqueId()))
	repo1, err := CreateECRRepoE(t, region, ecrRepoName)
	defer DeleteECRRepo(t, region, repo1)
	require.NoError(t, err)

	assert.Equal(t, ecrRepoName, aws.StringValue(repo1.RepositoryName))

	repo2, err := GetECRRepoE(t, region, ecrRepoName)
	require.NoError(t, err)
	assert.Equal(t, ecrRepoName, aws.StringValue(repo2.RepositoryName))
}
