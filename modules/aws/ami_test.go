package aws_test

import (
	"context"
	"testing"

	aws "github.com/gruntwork-io/terratest/modules/aws/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetUbuntu1404AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := aws.GetUbuntu1404AmiContext(t, context.Background(), "us-east-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}

func TestGetUbuntu1604AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := aws.GetUbuntu1604AmiContext(t, context.Background(), "us-west-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}

func TestGetUbuntu2004AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := aws.GetUbuntu2004AmiContext(t, context.Background(), "us-west-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}

func TestGetUbuntu2204AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := aws.GetUbuntu2204AmiContext(t, context.Background(), "us-west-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}

func TestGetCentos7AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := aws.GetCentos7AmiContext(t, context.Background(), "eu-west-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}

func TestGetAmazonLinuxAmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := aws.GetAmazonLinuxAmiContext(t, context.Background(), "ap-southeast-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}

func TestGetEcsOptimizedAmazonLinuxAmiEReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := aws.GetEcsOptimizedAmazonLinuxAmiContext(t, context.Background(), "us-east-2")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}
