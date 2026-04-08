package aws

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	ttesting "github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetRoute53RecordContextE returns a Route 53 Record.
// The ctx parameter supports cancellation and timeouts.
func GetRoute53RecordContextE(t ttesting.TestingT, ctx context.Context, hostedZoneID, recordName, recordType, awsRegion string) (*types.ResourceRecordSet, error) {
	t.Helper()

	route53Client, err := NewRoute53ClientContextE(t, ctx, awsRegion)
	if err != nil {
		return nil, err
	}

	o, err := route53Client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
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

// GetRoute53RecordContext returns a Route 53 Record.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func GetRoute53RecordContext(t ttesting.TestingT, ctx context.Context, hostedZoneID, recordName, recordType, awsRegion string) *types.ResourceRecordSet {
	t.Helper()
	r, err := GetRoute53RecordContextE(t, ctx, hostedZoneID, recordName, recordType, awsRegion)
	require.NoError(t, err)

	return r
}

// GetRoute53Record returns a Route 53 Record.
//
// Deprecated: Use [GetRoute53RecordContext] instead.
func GetRoute53Record(t ttesting.TestingT, hostedZoneID, recordName, recordType, awsRegion string) *types.ResourceRecordSet {
	t.Helper()
	return GetRoute53RecordContext(t, context.Background(), hostedZoneID, recordName, recordType, awsRegion)
}

// GetRoute53RecordE returns a Route 53 Record.
//
// Deprecated: Use [GetRoute53RecordContextE] instead.
func GetRoute53RecordE(t ttesting.TestingT, hostedZoneID, recordName, recordType, awsRegion string) (*types.ResourceRecordSet, error) {
	t.Helper()
	return GetRoute53RecordContextE(t, context.Background(), hostedZoneID, recordName, recordType, awsRegion)
}

// NewRoute53ClientContextE creates a Route 53 client.
// The ctx parameter supports cancellation and timeouts.
func NewRoute53ClientContextE(t ttesting.TestingT, ctx context.Context, region string) (*route53.Client, error) {
	t.Helper()

	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return route53.NewFromConfig(*sess), nil
}

// NewRoute53ClientContext creates a Route 53 client.
// This function will fail the test if there is an error.
// The ctx parameter supports cancellation and timeouts.
func NewRoute53ClientContext(t ttesting.TestingT, ctx context.Context, region string) *route53.Client {
	t.Helper()
	c, err := NewRoute53ClientContextE(t, ctx, region)
	require.NoError(t, err)

	return c
}

// NewRoute53Client creates a Route 53 client.
//
// Deprecated: Use [NewRoute53ClientContext] instead.
func NewRoute53Client(t ttesting.TestingT, region string) *route53.Client {
	t.Helper()
	return NewRoute53ClientContext(t, context.Background(), region)
}

// NewRoute53ClientE creates a Route 53 client.
//
// Deprecated: Use [NewRoute53ClientContextE] instead.
func NewRoute53ClientE(t ttesting.TestingT, region string) (*route53.Client, error) {
	t.Helper()
	return NewRoute53ClientContextE(t, context.Background(), region)
}
