package aws

import (
	"context"
	goerrors "errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/gruntwork-io/go-commons/errors"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// CreateECRRepoContextE creates a new ECR Repository.
// The ctx parameter supports cancellation and timeouts.
func CreateECRRepoContextE(t testing.TestingT, ctx context.Context, region string, name string) (*types.Repository, error) {
	client, err := NewECRClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	resp, err := client.CreateRepository(ctx, &ecr.CreateRepositoryInput{RepositoryName: aws.String(name)})
	if err != nil {
		return nil, err
	}

	return resp.Repository, nil
}

// CreateECRRepoContext creates a new ECR Repository.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func CreateECRRepoContext(t testing.TestingT, ctx context.Context, region string, name string) *types.Repository {
	t.Helper()

	repo, err := CreateECRRepoContextE(t, ctx, region, name)
	require.NoError(t, err)

	return repo
}

// GetECRRepoContextE gets an ECR Repository by name.
// An error occurs if a repository with the given name does not exist in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetECRRepoContextE(t testing.TestingT, ctx context.Context, region string, name string) (*types.Repository, error) {
	client, err := NewECRClientContextE(t, ctx, region)
	if err != nil {
		return nil, err
	}

	repositoryNames := []string{name}

	resp, err := client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{RepositoryNames: repositoryNames})
	if err != nil {
		return nil, err
	}

	if len(resp.Repositories) != 1 {
		return nil, errors.WithStackTrace(goerrors.New("an unexpected condition occurred. Please file an issue at github.com/gruntwork-io/terratest"))
	}

	return &resp.Repositories[0], nil
}

// GetECRRepoContext gets an ECR repository by name.
// This function will fail the test if there is an error.
// An error occurs if a repository with the given name does not exist in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetECRRepoContext(t testing.TestingT, ctx context.Context, region string, name string) *types.Repository {
	t.Helper()

	repo, err := GetECRRepoContextE(t, ctx, region, name)
	require.NoError(t, err)

	return repo
}

// DeleteECRRepoContextE will force delete the ECR repo by deleting all images prior to deleting the ECR repository.
// The ctx parameter supports cancellation and timeouts.
func DeleteECRRepoContextE(t testing.TestingT, ctx context.Context, region string, repo *types.Repository) error {
	client, err := NewECRClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	paginator := ecr.NewListImagesPaginator(client, &ecr.ListImagesInput{RepositoryName: repo.RepositoryName})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		if len(page.ImageIds) == 0 {
			continue
		}

		if _, err := client.BatchDeleteImage(ctx, &ecr.BatchDeleteImageInput{
			RepositoryName: repo.RepositoryName,
			ImageIds:       page.ImageIds,
		}); err != nil {
			return err
		}
	}

	if _, err := client.DeleteRepository(ctx, &ecr.DeleteRepositoryInput{RepositoryName: repo.RepositoryName}); err != nil {
		return err
	}

	return nil
}

// DeleteECRRepoContext will force delete the ECR repo by deleting all images prior to deleting the ECR repository.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteECRRepoContext(t testing.TestingT, ctx context.Context, region string, repo *types.Repository) {
	t.Helper()

	err := DeleteECRRepoContextE(t, ctx, region, repo)
	require.NoError(t, err)
}

// NewECRClientContextE returns a client for the Elastic Container Registry.
// The ctx parameter supports cancellation and timeouts.
func NewECRClientContextE(t testing.TestingT, ctx context.Context, region string) (*ecr.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return ecr.NewFromConfig(*sess), nil
}

// NewECRClientContext returns a client for the Elastic Container Registry.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewECRClientContext(t testing.TestingT, ctx context.Context, region string) *ecr.Client {
	t.Helper()

	client, err := NewECRClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// GetECRRepoLifecyclePolicyContextE gets the policies for the given ECR repository.
// The ctx parameter supports cancellation and timeouts.
func GetECRRepoLifecyclePolicyContextE(t testing.TestingT, ctx context.Context, region string, repo *types.Repository) (string, error) {
	client, err := NewECRClientContextE(t, ctx, region)
	if err != nil {
		return "", err
	}

	resp, err := client.GetLifecyclePolicy(ctx, &ecr.GetLifecyclePolicyInput{RepositoryName: repo.RepositoryName})
	if err != nil {
		return "", err
	}

	return *resp.LifecyclePolicyText, nil
}

// GetECRRepoLifecyclePolicyContext gets the policies for the given ECR repository.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetECRRepoLifecyclePolicyContext(t testing.TestingT, ctx context.Context, region string, repo *types.Repository) string {
	t.Helper()

	policy, err := GetECRRepoLifecyclePolicyContextE(t, ctx, region, repo)
	require.NoError(t, err)

	return policy
}

// PutECRRepoLifecyclePolicyContextE puts the given policy for the given ECR repository.
// The ctx parameter supports cancellation and timeouts.
func PutECRRepoLifecyclePolicyContextE(t testing.TestingT, ctx context.Context, region string, repo *types.Repository, policy string) error {
	logger.Default.Logf(t, "Applying policy for repository %s in %s", *repo.RepositoryName, region)

	client, err := NewECRClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	input := &ecr.PutLifecyclePolicyInput{
		RepositoryName:      repo.RepositoryName,
		LifecyclePolicyText: aws.String(policy),
	}

	_, err = client.PutLifecyclePolicy(ctx, input)

	return err
}

// PutECRRepoLifecyclePolicyContext puts the given policy for the given ECR repository.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func PutECRRepoLifecyclePolicyContext(t testing.TestingT, ctx context.Context, region string, repo *types.Repository, policy string) {
	t.Helper()

	err := PutECRRepoLifecyclePolicyContextE(t, ctx, region, repo, policy)
	require.NoError(t, err)
}

// GetECRRepoPolicyContextE gets the policies for the given ECR repository.
// The ctx parameter supports cancellation and timeouts.
func GetECRRepoPolicyContextE(t testing.TestingT, ctx context.Context, region string, repo *types.Repository) (string, error) {
	client, err := NewECRClientContextE(t, ctx, region)
	if err != nil {
		return "", err
	}

	resp, err := client.GetRepositoryPolicy(ctx, &ecr.GetRepositoryPolicyInput{RepositoryName: repo.RepositoryName})
	if err != nil {
		return "", err
	}

	return *resp.PolicyText, nil
}

// GetECRRepoPolicyContext gets the permissions for the given ECR repository.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetECRRepoPolicyContext(t testing.TestingT, ctx context.Context, region string, repo *types.Repository) string {
	t.Helper()

	policy, err := GetECRRepoPolicyContextE(t, ctx, region, repo)
	require.NoError(t, err)

	return policy
}

// PutECRRepoPolicyContextE puts the given policy for the given ECR repository.
// The ctx parameter supports cancellation and timeouts.
func PutECRRepoPolicyContextE(t testing.TestingT, ctx context.Context, region string, repo *types.Repository, policy string) error {
	logger.Default.Logf(t, "Applying repo policy for repository %s in %s", *repo.RepositoryName, region)

	client, err := NewECRClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	input := &ecr.SetRepositoryPolicyInput{
		PolicyText:     &policy,
		RepositoryName: repo.RepositoryName,
	}

	_, err = client.SetRepositoryPolicy(ctx, input)

	return err
}

// PutECRRepoPolicyContext puts the given policy for the given ECR repository.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func PutECRRepoPolicyContext(t testing.TestingT, ctx context.Context, region string, repo *types.Repository, policy string) {
	t.Helper()

	err := PutECRRepoPolicyContextE(t, ctx, region, repo, policy)
	require.NoError(t, err)
}
