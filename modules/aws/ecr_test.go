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
	repo1, err := CreateECRRepoE(t, region, fmt.Sprintf("terratest%s", strings.ToLower(random.UniqueId())))
	defer DeleteECRRepo(t, region, repo1)

	require.NoError(t, err)
	assert.Equal(t, "terratest", aws.StringValue(repo1.RepositoryName))

	repo2, err := GetECRRepoE(t, region, "terratest")

	require.NoError(t, err)
	assert.Equal(t, "terratest", aws.StringValue(repo2.RepositoryName))
}
