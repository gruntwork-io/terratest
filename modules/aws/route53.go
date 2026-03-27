package aws

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/stretchr/testify/require"
)

// GetRoute53Record returns a Route 53 Record.
func GetRoute53Record(t *testing.T, hostedZoneID, recordName, recordType, awsRegion string) *types.ResourceRecordSet {
	t.Helper()

	r, err := GetRoute53RecordE(t, hostedZoneID, recordName, recordType, awsRegion)
	require.NoError(t, err)

	return r
}

// GetRoute53RecordE returns a Route 53 Record.
func GetRoute53RecordE(t *testing.T, hostedZoneID, recordName, recordType, awsRegion string) (*types.ResourceRecordSet, error) {
	t.Helper()

	route53Client, err := NewRoute53ClientE(t, awsRegion)
	if err != nil {
		return nil, err
	}

	o, err := route53Client.ListResourceRecordSets(context.Background(), &route53.ListResourceRecordSetsInput{
		HostedZoneId:    &hostedZoneID,
		StartRecordName: &recordName,
		StartRecordType: types.RRType(recordType),
		MaxItems:        aws.Int32(1),
	})
	if err != nil {
		return nil, err
	}

	for i := range o.ResourceRecordSets {
		if strings.EqualFold(recordName+".", *o.ResourceRecordSets[i].Name) {
			return &o.ResourceRecordSets[i], nil
		}
	}

	return nil, errors.New("record not found")
}

// NewRoute53Client creates a Route 53 client.
func NewRoute53Client(t *testing.T, region string) *route53.Client {
	t.Helper()

	c, err := NewRoute53ClientE(t, region)
	require.NoError(t, err)

	return c
}

// NewRoute53ClientE creates a Route 53 client.
func NewRoute53ClientE(t *testing.T, region string) (*route53.Client, error) {
	t.Helper()

	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return route53.NewFromConfig(*sess), nil
}
