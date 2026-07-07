package aws

import (
	"context"
	goerrors "errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic"
	"github.com/aws/aws-sdk-go-v2/service/ecrpublic/types"
	"github.com/gruntwork-io/go-commons/errors"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// Note: the ECR Public API is only available in the us-east-1 region, so pass "us-east-1" as the region to these
// functions regardless of where the rest of your infrastructure lives.

// CreateECRPublicRepoContextE creates a new ECR Public Repository.
// The ctx parameter supports cancellation and timeouts.
func CreateECRPublicRepoContextE(t testing.TestingT, ctx context.Context, region string, name string) (*types.Repository, error) {
	client, err := NewECRPublicClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	resp, err := client.CreateRepository(ctx, &ecrpublic.CreateRepositoryInput{RepositoryName: aws.String(name)})
	if err != nil {
		return nil, err
	}

	return resp.Repository, nil
}

// CreateECRPublicRepoContext creates a new ECR Public Repository.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CreateECRPublicRepoContext(t testing.TestingT, ctx context.Context, region string, name string) *types.Repository {
	t.Helper()

	repo, err := CreateECRPublicRepoContextE(t, ctx, region, name)
	require.NoError(t, err)

	return repo
}

// CreateECRPublicRepoE creates a new ECR Public Repository.
func CreateECRPublicRepoE(t testing.TestingT, region string, name string) (*types.Repository, error) {
	return CreateECRPublicRepoContextE(t, context.Background(), region, name)
}

// CreateECRPublicRepo creates a new ECR Public Repository. This will fail the test and stop execution if there is an error.
func CreateECRPublicRepo(t testing.TestingT, region string, name string) *types.Repository {
	t.Helper()

	return CreateECRPublicRepoContext(t, context.Background(), region, name)
}

// GetECRPublicRepoContextE gets an ECR Public Repository by name.
// An error occurs if a repository with the given name does not exist.
// The ctx parameter supports cancellation and timeouts.
func GetECRPublicRepoContextE(t testing.TestingT, ctx context.Context, region string, name string) (*types.Repository, error) {
	client, err := NewECRPublicClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	resp, err := client.DescribeRepositories(ctx, &ecrpublic.DescribeRepositoriesInput{RepositoryNames: []string{name}})
	if err != nil {
		return nil, err
	}

	if len(resp.Repositories) != 1 {
		return nil, errors.WithStackTrace(goerrors.New("an unexpected condition occurred. Please file an issue at github.com/gruntwork-io/terratest"))
	}

	return &resp.Repositories[0], nil
}

// GetECRPublicRepoContext gets an ECR Public Repository by name.
// This function will fail the test if there is an error.
// An error occurs if a repository with the given name does not exist.
// The ctx parameter supports cancellation and timeouts.
func GetECRPublicRepoContext(t testing.TestingT, ctx context.Context, region string, name string) *types.Repository {
	t.Helper()

	repo, err := GetECRPublicRepoContextE(t, ctx, region, name)
	require.NoError(t, err)

	return repo
}

// GetECRPublicRepoE gets an ECR Public Repository by name.
// An error occurs if a repository with the given name does not exist.
func GetECRPublicRepoE(t testing.TestingT, region string, name string) (*types.Repository, error) {
	return GetECRPublicRepoContextE(t, context.Background(), region, name)
}

// GetECRPublicRepo gets an ECR Public Repository by name. This will fail the test and stop execution if there is an error.
// An error occurs if a repository with the given name does not exist.
func GetECRPublicRepo(t testing.TestingT, region string, name string) *types.Repository {
	t.Helper()

	return GetECRPublicRepoContext(t, context.Background(), region, name)
}

// DeleteECRPublicRepoContextE force deletes the ECR Public repository, including any images it contains.
// The ctx parameter supports cancellation and timeouts.
func DeleteECRPublicRepoContextE(t testing.TestingT, ctx context.Context, region string, repo *types.Repository) error {
	client, err := NewECRPublicClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	_, err = client.DeleteRepository(ctx, &ecrpublic.DeleteRepositoryInput{
		RepositoryName: repo.RepositoryName,
		Force:          true,
	})

	return err
}

// DeleteECRPublicRepoContext force deletes the ECR Public repository, including any images it contains.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteECRPublicRepoContext(t testing.TestingT, ctx context.Context, region string, repo *types.Repository) {
	t.Helper()

	err := DeleteECRPublicRepoContextE(t, ctx, region, repo)
	require.NoError(t, err)
}

// DeleteECRPublicRepoE force deletes the ECR Public repository, including any images it contains.
func DeleteECRPublicRepoE(t testing.TestingT, region string, repo *types.Repository) error {
	return DeleteECRPublicRepoContextE(t, context.Background(), region, repo)
}

// DeleteECRPublicRepo force deletes the ECR Public repository, including any images it contains.
// This will fail the test and stop execution if there is an error.
func DeleteECRPublicRepo(t testing.TestingT, region string, repo *types.Repository) {
	t.Helper()

	DeleteECRPublicRepoContext(t, context.Background(), region, repo)
}

// NewECRPublicClientContextE returns a client for the Elastic Container Registry Public.
// The ctx parameter supports cancellation and timeouts.
func NewECRPublicClientContextE(t testing.TestingT, ctx context.Context, region string) (*ecrpublic.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return ecrpublic.NewFromConfig(*sess), nil
}

// NewECRPublicClientContext returns a client for the Elastic Container Registry Public.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewECRPublicClientContext(t testing.TestingT, ctx context.Context, region string) *ecrpublic.Client {
	t.Helper()

	client, err := NewECRPublicClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// NewECRPublicClientE returns a client for the Elastic Container Registry Public.
func NewECRPublicClientE(t testing.TestingT, region string) (*ecrpublic.Client, error) {
	return NewECRPublicClientContextE(t, context.Background(), region)
}

// NewECRPublicClient returns a client for the Elastic Container Registry Public. This will fail the test and
// stop execution if there is an error.
func NewECRPublicClient(t testing.TestingT, region string) *ecrpublic.Client {
	t.Helper()

	return NewECRPublicClientContext(t, context.Background(), region)
}
