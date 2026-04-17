package aws

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/require"
)

// mockEbsClient is a test double for EbsAPI that captures the snapshot ID passed to DeleteSnapshot
// and returns a canned error.
type mockEbsClient struct {
	DeleteSnapshotErr error
	lastSnapshotID    string
	callCount         int
}

func (m *mockEbsClient) DeleteSnapshot(_ context.Context, params *ec2.DeleteSnapshotInput, _ ...func(*ec2.Options)) (*ec2.DeleteSnapshotOutput, error) {
	m.callCount++
	m.lastSnapshotID = aws.ToString(params.SnapshotId)

	if m.DeleteSnapshotErr != nil {
		return nil, m.DeleteSnapshotErr
	}
	return &ec2.DeleteSnapshotOutput{}, nil
}

func TestDeleteEbsSnapshotWithClientContextE(t *testing.T) {
	t.Parallel()

	t.Run("forwards snapshot id and succeeds", func(t *testing.T) {
		t.Parallel()

		client := &mockEbsClient{}
		err := DeleteEbsSnapshotWithClientContextE(t, context.Background(), client, "snap-0123456789abcdef0")
		require.NoError(t, err)
		require.Equal(t, 1, client.callCount)
		require.Equal(t, "snap-0123456789abcdef0", client.lastSnapshotID)
	})

	t.Run("propagates api error", func(t *testing.T) {
		t.Parallel()

		client := &mockEbsClient{DeleteSnapshotErr: errors.New("InvalidSnapshot.NotFound")}
		err := DeleteEbsSnapshotWithClientContextE(t, context.Background(), client, "snap-missing")
		require.Error(t, err)
	})
}
