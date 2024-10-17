package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// DeleteEbsSnapshot deletes the given EBS snapshot
func DeleteEbsSnapshot(t testing.TestingT, region string, snapshot string) {
	err := DeleteEbsSnapshotE(t, region, snapshot)
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteEbsSnapshotE deletes the given EBS snapshot
func DeleteEbsSnapshotE(t testing.TestingT, region string, snapshot string) error {
	logger.Default.Logf(t, "Deleting EBS snapshot %s", snapshot)
	ec2Client, err := NewEc2ClientE(t, region)
	if err != nil {
		return err
	}

	_, err = ec2Client.DeleteSnapshot(context.Background(), &ec2.DeleteSnapshotInput{
		SnapshotId: aws.String(snapshot),
	})
	return err
}
