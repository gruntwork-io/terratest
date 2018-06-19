package aws

import (
	"testing"

	"github.com/Briansbum/terratest/modules/logger"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// DeleteEbsSnapshot deletes the given EBS snapshot
func DeleteEbsSnapshot(t *testing.T, region string, snapshot string, sessExists ...*session.Session) {
	err := DeleteEbsSnapshotE(t, region, snapshot, sessExists[0])
	if err != nil {
		t.Fatal(err)
	}
}

// DeleteEbsSnapshot deletes the given EBS snapshot
func DeleteEbsSnapshotE(t *testing.T, region string, snapshot string, sessExists ...*session.Session) error {
	logger.Logf(t, "Deleting EBS snapshot %s", snapshot)
	ec2Client, err := NewEc2ClientE(t, region, sessExists[0])
	if err != nil {
		return err
	}

	_, err = ec2Client.DeleteSnapshot(&ec2.DeleteSnapshotInput{
		SnapshotId: aws.String(snapshot),
	})
	return err
}
