package aws

import (
	"context"
	goerrors "errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/gruntwork-io/go-commons/errors"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
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

// CreateECRRepo creates a new ECR Repository. This will fail the test and stop execution if there is an error.
//
// Deprecated: Use [CreateECRRepoContext] instead.
func CreateECRRepo(t testing.TestingT, region string, name string) *types.Repository {
	t.Helper()

	return CreateECRRepoContext(t, context.Background(), region, name)
}

// CreateECRRepoE creates a new ECR Repository.
//
// Deprecated: Use [CreateECRRepoContextE] instead.
func CreateECRRepoE(t testing.TestingT, region string, name string) (*types.Repository, error) {
	return CreateECRRepoContextE(t, context.Background(), region, name)
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

// GetECRRepo gets an ECR repository by name. This will fail the test and stop execution if there is an error.
// An error occurs if a repository with the given name does not exist in the given region.
//
// Deprecated: Use [GetECRRepoContext] instead.
func GetECRRepo(t testing.TestingT, region string, name string) *types.Repository {
	t.Helper()

	return GetECRRepoContext(t, context.Background(), region, name)
}

// GetECRRepoE gets an ECR Repository by name.
// An error occurs if a repository with the given name does not exist in the given region.
//
// Deprecated: Use [GetECRRepoContextE] instead.
func GetECRRepoE(t testing.TestingT, region string, name string) (*types.Repository, error) {
	return GetECRRepoContextE(t, context.Background(), region, name)
}

// DeleteECRRepoContextE will force delete the ECR repo by deleting all images prior to deleting the ECR repository.
// The ctx parameter supports cancellation and timeouts.
func DeleteECRRepoContextE(t testing.TestingT, ctx context.Context, region string, repo *types.Repository) error {
	client, err := NewECRClientContextE(t, ctx, region)
	if err != nil {
		return err
	}

	resp, err := client.ListImages(ctx, &ecr.ListImagesInput{RepositoryName: repo.RepositoryName})
	if err != nil {
		return err
	}

	if len(resp.ImageIds) > 0 {
		_, err = client.BatchDeleteImage(ctx, &ecr.BatchDeleteImageInput{
			RepositoryName: repo.RepositoryName,
			ImageIds:       resp.ImageIds,
		})
		if err != nil {
			return err
		}
	}

	_, err = client.DeleteRepository(ctx, &ecr.DeleteRepositoryInput{RepositoryName: repo.RepositoryName})
	if err != nil {
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

// DeleteECRRepo will force delete the ECR repo by deleting all images prior to deleting the ECR repository.
// This will fail the test and stop execution if there is an error.
//
// Deprecated: Use [DeleteECRRepoContext] instead.
func DeleteECRRepo(t testing.TestingT, region string, repo *types.Repository) {
	t.Helper()

	DeleteECRRepoContext(t, context.Background(), region, repo)
}

// DeleteECRRepoE will force delete the ECR repo by deleting all images prior to deleting the ECR repository.
//
// Deprecated: Use [DeleteECRRepoContextE] instead.
func DeleteECRRepoE(t testing.TestingT, region string, repo *types.Repository) error {
	return DeleteECRRepoContextE(t, context.Background(), region, repo)
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

// NewECRClient returns a client for the Elastic Container Registry. This will fail the test and
// stop execution if there is an error.
//
// Deprecated: Use [NewECRClientContext] instead.
func NewECRClient(t testing.TestingT, region string) *ecr.Client {
	t.Helper()

	return NewECRClientContext(t, context.Background(), region)
}

// NewECRClientE returns a client for the Elastic Container Registry.
//
// Deprecated: Use [NewECRClientContextE] instead.
func NewECRClientE(t testing.TestingT, region string) (*ecr.Client, error) {
	return NewECRClientContextE(t, context.Background(), region)
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

// GetECRRepoLifecyclePolicy gets the policies for the given ECR repository.
// This will fail the test and stop execution if there is an error.
//
// Deprecated: Use [GetECRRepoLifecyclePolicyContext] instead.
func GetECRRepoLifecyclePolicy(t testing.TestingT, region string, repo *types.Repository) string {
	t.Helper()

	return GetECRRepoLifecyclePolicyContext(t, context.Background(), region, repo)
}

// GetECRRepoLifecyclePolicyE gets the policies for the given ECR repository.
//
// Deprecated: Use [GetECRRepoLifecyclePolicyContextE] instead.
func GetECRRepoLifecyclePolicyE(t testing.TestingT, region string, repo *types.Repository) (string, error) {
	return GetECRRepoLifecyclePolicyContextE(t, context.Background(), region, repo)
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

// PutECRRepoLifecyclePolicy puts the given policy for the given ECR repository.
// This will fail the test and stop execution if there is an error.
//
// Deprecated: Use [PutECRRepoLifecyclePolicyContext] instead.
func PutECRRepoLifecyclePolicy(t testing.TestingT, region string, repo *types.Repository, policy string) {
	t.Helper()

	PutECRRepoLifecyclePolicyContext(t, context.Background(), region, repo, policy)
}

// PutECRRepoLifecyclePolicyE puts the given policy for the given ECR repository.
//
// Deprecated: Use [PutECRRepoLifecyclePolicyContextE] instead.
func PutECRRepoLifecyclePolicyE(t testing.TestingT, region string, repo *types.Repository, policy string) error {
	return PutECRRepoLifecyclePolicyContextE(t, context.Background(), region, repo, policy)
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

// GetECRRepoPolicy gets the permissions for the given ECR repository.
// This will fail the test and stop execution if there is an error.
//
// Deprecated: Use [GetECRRepoPolicyContext] instead.
func GetECRRepoPolicy(t testing.TestingT, region string, repo *types.Repository) string {
	t.Helper()

	return GetECRRepoPolicyContext(t, context.Background(), region, repo)
}

// GetECRRepoPolicyE gets the policies for the given ECR repository.
//
// Deprecated: Use [GetECRRepoPolicyContextE] instead.
func GetECRRepoPolicyE(t testing.TestingT, region string, repo *types.Repository) (string, error) {
	return GetECRRepoPolicyContextE(t, context.Background(), region, repo)
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

// PutECRRepoPolicy puts the given policy for the given ECR repository.
// This will fail the test and stop execution if there is an error.
//
// Deprecated: Use [PutECRRepoPolicyContext] instead.
func PutECRRepoPolicy(t testing.TestingT, region string, repo *types.Repository, policy string) {
	t.Helper()

	PutECRRepoPolicyContext(t, context.Background(), region, repo, policy)
}

// PutECRRepoPolicyE puts the given policy for the given ECR repository.
//
// Deprecated: Use [PutECRRepoPolicyContextE] instead.
func PutECRRepoPolicyE(t testing.TestingT, region string, repo *types.Repository, policy string) error {
	return PutECRRepoPolicyContextE(t, context.Background(), region, repo, policy)
}
