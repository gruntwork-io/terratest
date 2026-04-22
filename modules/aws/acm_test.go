package aws_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stretchr/testify/require"

	aws "github.com/gruntwork-io/terratest/modules/aws"
)

// mockAcmClient is a test double for aws.AcmAPI that returns canned responses.
// When ListCertificatesPages is non-empty, it returns one page per call in order,
// advancing on each invocation; otherwise it returns ListCertificatesOutput on every call.
type mockAcmClient struct {
	ListCertificatesOutput *acm.ListCertificatesOutput
	ListCertificatesErr    error
	ListCertificatesPages  []*acm.ListCertificatesOutput
	callCount              int
}

func (m *mockAcmClient) ListCertificates(_ context.Context, _ *acm.ListCertificatesInput, _ ...func(*acm.Options)) (*acm.ListCertificatesOutput, error) {
	if m.ListCertificatesErr != nil {
		return nil, m.ListCertificatesErr
	}

	if len(m.ListCertificatesPages) > 0 {
		if m.callCount >= len(m.ListCertificatesPages) {
			return nil, fmt.Errorf("mockAcmClient: ListCertificates called %d times but only %d page(s) configured", m.callCount+1, len(m.ListCertificatesPages))
		}

		page := m.ListCertificatesPages[m.callCount]
		m.callCount++

		return page, nil
	}

	return m.ListCertificatesOutput, nil
}

func TestGetAcmCertificateArnWithClientContextE(t *testing.T) {
	t.Parallel()

	const (
		arn1    = "arn:aws:acm:us-east-1:123456789012:certificate/cert-1"
		arn2    = "arn:aws:acm:us-east-1:123456789012:certificate/cert-2"
		domain1 = "foo.example.com"
		domain2 = "bar.example.com"
	)

	twoCerts := &acm.ListCertificatesOutput{
		CertificateSummaryList: []types.CertificateSummary{
			{DomainName: awsSDK.String(domain1), CertificateArn: awsSDK.String(arn1)},
			{DomainName: awsSDK.String(domain2), CertificateArn: awsSDK.String(arn2)},
		},
	}

	page1 := &acm.ListCertificatesOutput{
		CertificateSummaryList: []types.CertificateSummary{
			{DomainName: awsSDK.String(domain1), CertificateArn: awsSDK.String(arn1)},
		},
		NextToken: awsSDK.String("page-2-token"),
	}
	page2 := &acm.ListCertificatesOutput{
		CertificateSummaryList: []types.CertificateSummary{
			{DomainName: awsSDK.String(domain2), CertificateArn: awsSDK.String(arn2)},
		},
	}

	tests := map[string]struct {
		client      *mockAcmClient
		query       string
		expectedArn string
		expectErr   bool
	}{
		"returns arn when domain matches": {
			client:      &mockAcmClient{ListCertificatesOutput: twoCerts},
			query:       domain2,
			expectedArn: arn2,
		},
		"returns first match when listed first": {
			client:      &mockAcmClient{ListCertificatesOutput: twoCerts},
			query:       domain1,
			expectedArn: arn1,
		},
		"returns empty string when no domain matches": {
			client:      &mockAcmClient{ListCertificatesOutput: twoCerts},
			query:       "nonexistent.example.com",
			expectedArn: "",
		},
		"returns empty string on empty list": {
			client:      &mockAcmClient{ListCertificatesOutput: &acm.ListCertificatesOutput{}},
			query:       domain1,
			expectedArn: "",
		},
		"propagates api error": {
			client:    &mockAcmClient{ListCertificatesErr: errors.New("AccessDenied")},
			query:     domain1,
			expectErr: true,
		},
		"finds arn on second page via next token": {
			client:      &mockAcmClient{ListCertificatesPages: []*acm.ListCertificatesOutput{page1, page2}},
			query:       domain2,
			expectedArn: arn2,
		},
		"stops paginating when last page has no next token": {
			// page2 has no NextToken, so the caller must stop after page 2 rather
			// than calling the mock a third time and over-running ListCertificatesPages.
			client:      &mockAcmClient{ListCertificatesPages: []*acm.ListCertificatesOutput{page1, page2}},
			query:       "nonexistent.example.com",
			expectedArn: "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			arn, err := aws.GetAcmCertificateArnWithClientContextE(t, context.Background(), tc.client, tc.query)
			if tc.expectErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedArn, arn)
		})
	}
}
