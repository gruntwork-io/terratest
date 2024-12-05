package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/acm"

	"github.com/gruntwork-io/terratest/modules/testing"
)

// GetAcmCertificateArn gets the ACM certificate for the given domain name in the given region.
func GetAcmCertificateArn(t testing.TestingT, awsRegion string, certDomainName string) string {
	arn, err := GetAcmCertificateArnE(t, awsRegion, certDomainName)
	if err != nil {
		t.Fatal(err)
	}
	return arn
}

// GetAcmCertificateArnE gets the ACM certificate for the given domain name in the given region.
func GetAcmCertificateArnE(t testing.TestingT, awsRegion string, certDomainName string) (string, error) {
	acmClient, err := NewAcmClientE(t, awsRegion)
	if err != nil {
		return "", err
	}

	result, err := acmClient.ListCertificates(context.Background(), &acm.ListCertificatesInput{})
	if err != nil {
		return "", err
	}

	for _, summary := range result.CertificateSummaryList {
		if *summary.DomainName == certDomainName {
			return *summary.CertificateArn, nil
		}
	}

	return "", nil
}

// NewAcmClient create a new ACM client.
func NewAcmClient(t testing.TestingT, region string) *acm.Client {
	client, err := NewAcmClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewAcmClientE creates a new ACM client.
func NewAcmClientE(t testing.TestingT, awsRegion string) (*acm.Client, error) {
	sess, err := NewAuthenticatedSession(awsRegion)
	if err != nil {
		return nil, err
	}

	return acm.NewFromConfig(*sess), nil
}
