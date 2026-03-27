package aws_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	aws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoute53Record(t *testing.T) {
	t.Parallel()
	region := aws.GetRandomStableRegion(t, nil, nil)
	c, err := aws.NewRoute53ClientE(t, region)
	require.NoError(t, err)

	domain := "terratest" + strconv.FormatInt(time.Now().UnixNano(), 10) + "example.com"
	hostedZone, err := c.CreateHostedZone(context.Background(), &route53.CreateHostedZoneInput{
		Name:            awsSDK.String(domain),
		CallerReference: awsSDK.String(strconv.FormatInt(time.Now().UnixNano(), 10)),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := c.DeleteHostedZone(context.Background(), &route53.DeleteHostedZoneInput{
			Id: hostedZone.HostedZone.Id,
		})
		require.NoError(t, err)
	})

	recordName := "record." + domain
	resourceRecordSet := &types.ResourceRecordSet{
		Name: &recordName,
		Type: types.RRTypeA,
		TTL:  awsSDK.Int64(60),
		ResourceRecords: []types.ResourceRecord{
			{
				Value: awsSDK.String("127.0.0.1"),
			},
		},
	}
	_, err = c.ChangeResourceRecordSets(context.Background(), &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: hostedZone.HostedZone.Id,
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action:            types.ChangeActionCreate,
					ResourceRecordSet: resourceRecordSet,
				},
			},
		},
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		_, err := c.ChangeResourceRecordSets(context.Background(), &route53.ChangeResourceRecordSetsInput{
			HostedZoneId: hostedZone.HostedZone.Id,
			ChangeBatch: &types.ChangeBatch{
				Changes: []types.Change{
					{
						Action:            types.ChangeActionDelete,
						ResourceRecordSet: resourceRecordSet,
					},
				},
			},
		})
		require.NoError(t, err)
	})

	t.Run("ExistingRecord", func(t *testing.T) {
		t.Parallel()

		route53Record := aws.GetRoute53Record(t, *hostedZone.HostedZone.Id, recordName, string(resourceRecordSet.Type), region)
		require.NotNil(t, route53Record)
		assert.Equal(t, recordName+".", *route53Record.Name)
		assert.Equal(t, resourceRecordSet.Type, route53Record.Type)
		assert.Equal(t, "127.0.0.1", *route53Record.ResourceRecords[0].Value)
	})

	t.Run("NotExistRecord", func(t *testing.T) {
		t.Parallel()

		route53Record, err := aws.GetRoute53RecordE(t, *hostedZone.HostedZone.Id, "ne"+recordName, "A", region)
		require.Error(t, err)
		assert.Nil(t, route53Record)
	})
}
