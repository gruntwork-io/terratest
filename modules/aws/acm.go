package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/acm"
)

// GetAcmCertificateArn gets the ACM certificate for the given domain name in the given region.
func GetAcmCertificateArn(t *testing.T, awsRegion string, certDomainName string, sessExists ...*session.Session) string {
	arn, err := GetAcmCertificateArnE(t, awsRegion, certDomainName, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return arn
}

// GetAcmCertificateArnE gets the ACM certificate for the given domain name in the given region.
func GetAcmCertificateArnE(t *testing.T, awsRegion string, certDomainName string, sessExists ...*session.Session) (string, error) {
	acmClient, err := NewAcmClientE(t, awsRegion, sessExists[0])
	if err != nil {
		return "", err
	}

	result, err := acmClient.ListCertificates(&acm.ListCertificatesInput{})
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
func NewAcmClient(t *testing.T, region string, sessExists ...*session.Session) *acm.ACM {
	client, err := NewAcmClientE(t, region, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
	return client
}

// NewAcmClientE creates a new ACM client.
func NewAcmClientE(t *testing.T, awsRegion string, sessExists ...*session.Session) (*acm.ACM, error) {
	sess, err := NewAuthenticatedSession(awsRegion, sessExists[0])
	if err != nil {
		return nil, err
	}

	return acm.New(sess), nil
}
