package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// EbsAPI is the subset of *ec2.Client operations used by the EBS helpers in this file. Declared
// as an interface so tests can substitute a mock; a real *ec2.Client satisfies it automatically.
type EbsAPI interface {
	DeleteSnapshot(ctx context.Context, params *ec2.DeleteSnapshotInput, optFns ...func(*ec2.Options)) (*ec2.DeleteSnapshotOutput, error)
}

// DeleteEbsSnapshotContextE deletes the given EBS snapshot.
// The ctx parameter supports cancellation and timeouts.
func DeleteEbsSnapshotContextE(t testing.TestingT, ctx context.Context, region string, snapshot string) error {
	logger.Default.Logf(t, "Deleting EBS snapshot %s", snapshot)

	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return err
	}

	return DeleteEbsSnapshotWithClientContextE(t, ctx, ec2.NewFromConfig(*sess), snapshot)
}

// DeleteEbsSnapshotWithClientContextE deletes the given EBS snapshot using the provided EC2 client.
// The ctx parameter supports cancellation and timeouts.
func DeleteEbsSnapshotWithClientContextE(t testing.TestingT, ctx context.Context, client EbsAPI, snapshot string) error {
	_, err := client.DeleteSnapshot(ctx, &ec2.DeleteSnapshotInput{
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
