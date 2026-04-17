package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// DeleteEbsSnapshotContextE deletes the given EBS snapshot.
// The ctx parameter supports cancellation and timeouts.
func DeleteEbsSnapshotContextE(t testing.TestingT, ctx context.Context, region string, snapshot string) error {
	logger.Default.Logf(t, "Deleting EBS snapshot %s", snapshot)

	sess, err := NewAuthConfigContextE(t, ctx, region)
	if err != nil {
		return err
	}

	ec2Client := ec2.NewFromConfig(*sess)

	_, err = ec2Client.DeleteSnapshot(ctx, &ec2.DeleteSnapshotInput{
		SnapshotId: aws.String(snapshot),
	})

	return err
}

// DeleteEbsSnapshotContext deletes the given EBS snapshot.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func DeleteEbsSnapshotContext(t testing.TestingT, ctx context.Context, region string, snapshot string) {
	t.Helper()

	err := DeleteEbsSnapshotContextE(t, ctx, region, snapshot)
	require.NoError(t, err)
}

// DeleteEbsSnapshot deletes the given EBS snapshot.
//
// Deprecated: Use [DeleteEbsSnapshotContext] instead.
func DeleteEbsSnapshot(t testing.TestingT, region string, snapshot string) {
	t.Helper()

	DeleteEbsSnapshotContext(t, context.Background(), region, snapshot)
}

// DeleteEbsSnapshotE deletes the given EBS snapshot.
//
// Deprecated: Use [DeleteEbsSnapshotContextE] instead.
func DeleteEbsSnapshotE(t testing.TestingT, region string, snapshot string) error {
	return DeleteEbsSnapshotContextE(t, context.Background(), region, snapshot)
}
