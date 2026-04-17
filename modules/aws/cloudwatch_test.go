package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/stretchr/testify/require"
)

// mockCloudWatchLogsClient is a test double for CloudWatchLogsAPI that returns canned responses.
type mockCloudWatchLogsClient struct {
	GetLogEventsOutput *cloudwatchlogs.GetLogEventsOutput
	GetLogEventsErr    error
}

func (m *mockCloudWatchLogsClient) GetLogEvents(_ context.Context, _ *cloudwatchlogs.GetLogEventsInput, _ ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetLogEventsOutput, error) {
	if m.GetLogEventsErr != nil {
		return nil, m.GetLogEventsErr
	}
	return m.GetLogEventsOutput, nil
}

func TestGetCloudWatchLogEntriesWithClientContextE(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		client   *mockCloudWatchLogsClient
		expected []string
		wantErr  bool
	}{
		"returns messages preserving order": {
			client: &mockCloudWatchLogsClient{
				GetLogEventsOutput: &cloudwatchlogs.GetLogEventsOutput{
					Events: []types.OutputLogEvent{
						{Message: aws.String("first line")},
						{Message: aws.String("second line")},
						{Message: aws.String("third line")},
					},
				},
			},
			expected: []string{"first line", "second line", "third line"},
		},
		"returns nil slice on empty events": {
			client: &mockCloudWatchLogsClient{
				GetLogEventsOutput: &cloudwatchlogs.GetLogEventsOutput{},
			},
			expected: nil,
		},
		"propagates api error": {
			client:  &mockCloudWatchLogsClient{GetLogEventsErr: errors.New("ResourceNotFoundException")},
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := GetCloudWatchLogEntriesWithClientContextE(t, context.Background(), tc.client, "stream", "group")
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expected, got)
		})
	}
}
