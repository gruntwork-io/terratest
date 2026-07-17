package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/stretchr/testify/require"

	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
)

// AcmAPI is the subset of *acm.Client operations used by the helpers in this file.
// It is declared as an interface so tests can substitute a mock without an AWS
// account. A real *acm.Client satisfies this interface automatically.
type AcmAPI interface {
	ListCertificates(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error)
}

// GetAcmCertificateArnContextE gets the ACM certificate for the given domain name in the given region.
// The ctx parameter supports cancellation and timeouts.
func GetAcmCertificateArnContextE(t testing.TestingT, ctx context.Context, awsRegion string, certDomainName string) (string, error) {
	acmClient, err := NewAcmClientContextE(t, ctx, awsRegion)
	if err != nil {
		return "", err
	}

	return GetAcmCertificateArnWithClientContextE(t, ctx, acmClient, certDomainName)
}

// GetAcmCertificateArnWithClientContextE gets the ACM certificate for the given domain name using
// the provided ACM client. Useful when a pre-configured client is available or in unit tests with
// a mock.
// The ctx parameter supports cancellation and timeouts.
func GetAcmCertificateArnWithClientContextE(t testing.TestingT, ctx context.Context, client AcmAPI, certDomainName string) (string, error) {
	input := &acm.ListCertificatesInput{}

	for {
		result, err := client.ListCertificates(ctx, input)
		if err != nil {
			return "", err
		}

		for i := range result.CertificateSummaryList {
			summary := &result.CertificateSummaryList[i]

			if *summary.DomainName == certDomainName {
				return *summary.CertificateArn, nil
			}
		}

		if result.NextToken == nil || *result.NextToken == "" {
			return "", nil
		}

		input.NextToken = result.NextToken
	}
}

// GetAcmCertificateArnContext gets the ACM certificate for the given domain name in the given region.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetAcmCertificateArnContext(t testing.TestingT, ctx context.Context, awsRegion string, certDomainName string) string {
	t.Helper()

	arn, err := GetAcmCertificateArnContextE(t, ctx, awsRegion, certDomainName)
	require.NoError(t, err)

	return arn
}

// NewAcmClientContextE creates a new ACM client.
// The ctx parameter supports cancellation and timeouts.
func NewAcmClientContextE(t testing.TestingT, ctx context.Context, region string) (*acm.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return acm.NewFromConfig(*sess), nil
}

// NewAcmClientContext creates a new ACM client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewAcmClientContext(t testing.TestingT, ctx context.Context, region string) *acm.Client {
	t.Helper()

	client, err := NewAcmClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}
