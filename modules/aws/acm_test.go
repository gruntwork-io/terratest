package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/stretchr/testify/require"
)

// mockAcmClient is a test double for AcmAPI that returns canned responses.
type mockAcmClient struct {
	ListCertificatesOutput *acm.ListCertificatesOutput
	ListCertificatesErr    error
}

func (m *mockAcmClient) ListCertificates(_ context.Context, _ *acm.ListCertificatesInput, _ ...func(*acm.Options)) (*acm.ListCertificatesOutput, error) {
	if m.ListCertificatesErr != nil {
		return nil, m.ListCertificatesErr
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
			{DomainName: aws.String(domain1), CertificateArn: aws.String(arn1)},
			{DomainName: aws.String(domain2), CertificateArn: aws.String(arn2)},
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			arn, err := GetAcmCertificateArnWithClientContextE(t, context.Background(), tc.client, tc.query)
			if tc.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedArn, arn)
		})
	}
}

