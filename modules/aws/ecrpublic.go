package aws

import (
	goerrors "errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecrpublic"
	"github.com/gruntwork-io/go-commons/errors"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// CreateECRRepo creates a new ECR Repository. This will fail the test and stop execution if there is an error.
func CreateECRPublicRepo(t testing.TestingT, region string, name string) *ecrpublic.Repository {
	repo, err := CreateECRPublicRepoE(t, region, name)
	require.NoError(t, err)
	return repo
}

// CreateECRRepoE creates a new ECR Repository.
func CreateECRPublicRepoE(t testing.TestingT, region string, name string) (*ecrpublic.Repository, error) {
	client := NewECRPublicClient(t, region)
	resp, err := client.CreateRepository(&ecrpublic.CreateRepositoryInput{RepositoryName: aws.String(name)})
	if err != nil {
		return nil, err
	}
	return resp.Repository, nil
}

// GetECRRepo gets an ECR repository by name. This will fail the test and stop execution if there is an error.
// An error occurs if a repository with the given name does not exist in the given region.
func GetECRPublicRepo(t testing.TestingT, region string, name string) *ecrpublic.Repository {
	repo, err := GetECRPublicRepoE(t, region, name)
	require.NoError(t, err)
	return repo
}

// GetECRRepoE gets an ECR Repository by name.
// An error occurs if a repository with the given name does not exist in the given region.
func GetECRPublicRepoE(t testing.TestingT, region string, name string) (*ecrpublic.Repository, error) {
	client := NewECRPublicClient(t, region)
	repositoryNames := []*string{aws.String(name)}
	resp, err := client.DescribeRepositories(&ecrpublic.DescribeRepositoriesInput{RepositoryNames: repositoryNames})
	if err != nil {
		return nil, err
	}
	if len(resp.Repositories) != 1 {
		return nil, errors.WithStackTrace(goerrors.New(("An unexpected condition occurred. Please file an issue at github.com/gruntwork-io/terratest")))
	}
	return resp.Repositories[0], nil
}

// DeleteECRPublicRepo will force delete the ECR repo by deleting all images prior to deleting the ECR repository.
// This will fail the test and stop execution if there is an error.
func DeleteECRPublicRepo(t testing.TestingT, region string, repo *ecrpublic.Repository) {
	err := DeleteECRPublicRepoE(t, region, repo)
	require.NoError(t, err)
}

// DeleteECRPublicRepoE will force delete the ECR repo by deleting all images prior to deleting the ECR repository.
func DeleteECRPublicRepoE(t testing.TestingT, region string, repo *ecrpublic.Repository) error {
	client := NewECRPublicClient(t, region)
	_, err := client.DeleteRepository(&ecrpublic.DeleteRepositoryInput{
		RepositoryName: repo.RepositoryName,
	})
	if err != nil {
		return err
	}
	return nil
}

// NewECRPublicClient returns a client for the Elastic Container Registry. This will fail the test and
// stop execution if there is an error.
func NewECRPublicClient(t testing.TestingT, region string) *ecrpublic.ECRPublic {
	sess, err := NewECRPublicClientE(t, region)
	require.NoError(t, err)
	return sess
}

// NewECRPublicClient returns a client for the Elastic Container Registry.
func NewECRPublicClientE(t testing.TestingT, region string) (*ecrpublic.ECRPublic, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return ecrpublic.New(sess), nil
}

// GetECRRepoLifecyclePolicy gets the policies for the given ECR repository.
// This will fail the test and stop execution if there is an error.
func GetECRPublicRepoPolicy(t testing.TestingT, region string, repo *ecrpublic.Repository) string {
	policy, err := GetECRPublicRepoPolicyE(t, region, repo)
	require.NoError(t, err)
	return policy
}

// GetECRRepoLifecyclePolicyE gets the policies for the given ECR repository.
func GetECRPublicRepoPolicyE(t testing.TestingT, region string, repo *ecrpublic.Repository) (string, error) {
	client := NewECRPublicClient(t, region)
	resp, err := client.GetRepositoryPolicy(&ecrpublic.GetRepositoryPolicyInput{RepositoryName: repo.RepositoryName})
	if err != nil {
		return "", err
	}
	return *resp.PolicyText, nil
}

// PutECRRepoLifecyclePolicy puts the given policy for the given ECR repository.
// This will fail the test and stop execution if there is an error.
func SetECRPublicRepoPolicy(t testing.TestingT, region string, repo *ecrpublic.Repository, policy string) {
	err := SetECRPublicRepoPolicyE(t, region, repo, policy)
	require.NoError(t, err)
}

// PutEcrRepoLifecyclePolicy puts the given policy for the given ECR repository.
func SetECRPublicRepoPolicyE(t testing.TestingT, region string, repo *ecrpublic.Repository, policy string) error {
	logger.Logf(t, "Applying policy for repository %s in %s", *repo.RepositoryName, region)

	client, err := NewECRPublicClientE(t, region)
	if err != nil {
		return err
	}

	input := &ecrpublic.SetRepositoryPolicyInput{
		RepositoryName: repo.RepositoryName,
		PolicyText:     aws.String(policy),
	}

	_, err = client.SetRepositoryPolicy(input)
	return err
}
