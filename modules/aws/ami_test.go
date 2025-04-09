package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetUbuntu1404AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := GetUbuntu1404Ami(t, "us-east-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}

func TestGetUbuntu1604AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := GetUbuntu1604Ami(t, "us-west-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}

func TestGetUbuntu2004AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := GetUbuntu2004Ami(t, "us-west-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}

func TestGetUbuntu2204AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := GetUbuntu2204Ami(t, "us-west-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}

func TestGetCentos7AmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := GetCentos7Ami(t, "eu-west-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}

func TestGetAmazonLinuxAmiReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := GetAmazonLinuxAmi(t, "ap-southeast-1")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}

func TestGetEcsOptimizedAmazonLinuxAmiEReturnsSomeAmi(t *testing.T) {
	t.Parallel()

	amiID := GetEcsOptimizedAmazonLinuxAmi(t, "us-east-2")
	assert.Regexp(t, "^ami-[[:alnum:]]+$", amiID)
}
