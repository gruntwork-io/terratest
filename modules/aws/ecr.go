package aws

import (
	goerrors "errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/gruntwork-io/go-commons/errors"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// CreateECRRepo creates a new ECR Repository. This will fail the test and stop execution if there is an error.
func CreateECRRepo(t testing.TestingT, region string, name string) *ecr.Repository {
	repo, err := CreateECRRepoE(t, region, name)
	require.NoError(t, err)
	return repo
}

// CreateECRRepoE creates a new ECR Repository.
func CreateECRRepoE(t testing.TestingT, region string, name string) (*ecr.Repository, error) {
	client := NewECRClient(t, region)
	resp, err := client.CreateRepository(&ecr.CreateRepositoryInput{RepositoryName: aws.String(name)})
	if err != nil {
		return nil, err
	}
	return resp.Repository, nil
}

// GetECRRepo gets an ECR repository by name. This will fail the test and stop execution if there is an error.
// An error occurs if a repository with the given name does not exist in the given region.
func GetECRRepo(t testing.TestingT, region string, name string) *ecr.Repository {
	repo, err := GetECRRepoE(t, region, name)
	require.NoError(t, err)
	return repo
}

// GetECRRepoE gets an ECR Repository by name.
// An error occurs if a repository with the given name does not exist in the given region.
func GetECRRepoE(t testing.TestingT, region string, name string) (*ecr.Repository, error) {
	client := NewECRClient(t, region)
	repositoryNames := []*string{aws.String(name)}
	resp, err := client.DescribeRepositories(&ecr.DescribeRepositoriesInput{RepositoryNames: repositoryNames})
	if err != nil {
		return nil, err
	}
	if len(resp.Repositories) != 1 {
		return nil, errors.WithStackTrace(goerrors.New(("An unexpected condition occurred. Please file an issue at github.com/gruntwork-io/terratest")))
	}
	return resp.Repositories[0], nil
}

// DeleteECRRepo will force delete the ECR repo by deleting all images prior to deleting the ECR repository.
// This will fail the test and stop execution if there is an error.
func DeleteECRRepo(t testing.TestingT, region string, repo *ecr.Repository) {
	err := DeleteECRRepoE(t, region, repo)
	require.NoError(t, err)
}

// DeleteECRRepoE will force delete the ECR repo by deleting all images prior to deleting the ECR repository.
func DeleteECRRepoE(t testing.TestingT, region string, repo *ecr.Repository) error {
	client := NewECRClient(t, region)
	resp, err := client.ListImages(&ecr.ListImagesInput{RepositoryName: repo.RepositoryName})
	if err != nil {
		return err
	}
	if len(resp.ImageIds) > 0 {
		_, err = client.BatchDeleteImage(&ecr.BatchDeleteImageInput{
			RepositoryName: repo.RepositoryName,
			ImageIds:       resp.ImageIds,
		})
		if err != nil {
			return err
		}
	}

	_, err = client.DeleteRepository(&ecr.DeleteRepositoryInput{RepositoryName: repo.RepositoryName})
	if err != nil {
		return err
	}
	return nil
}

// NewECRClient returns a client for the Elastic Container Registry. This will fail the test and
// stop execution if there is an error.
func NewECRClient(t testing.TestingT, region string) *ecr.ECR {
	sess, err := NewECRClientE(t, region)
	require.NoError(t, err)
	return sess
}

// NewECRClient returns a client for the Elastic Container Registry.
func NewECRClientE(t testing.TestingT, region string) (*ecr.ECR, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return ecr.New(sess), nil
}
