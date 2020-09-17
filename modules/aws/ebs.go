package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
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

// DeleteEbsSnapshot deletes the given EBS snapshot
func DeleteEbsSnapshotE(t testing.TestingT, region string, snapshot string) error {
	logger.Logf(t, "Deleting EBS snapshot %s", snapshot)
	ec2Client, err := NewEc2ClientE(t, region)
	if err != nil {
		return err
	}

	_, err = ec2Client.DeleteSnapshot(&ec2.DeleteSnapshotInput{
		SnapshotId: aws.String(snapshot),
	})
	return err
}
