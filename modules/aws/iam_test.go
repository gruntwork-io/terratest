package aws

import (
	"strings"
	"testing"

	"github.com/Briansbum/terratest/modules/random"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/stretchr/testify/assert"
)

func TestGetIamCurrentUserName(t *testing.T) {
	t.Parallel()

	username := GetIamCurrentUserName(t)
	assert.NotEmpty(t, username)
}

func TestGetIamCurrentUserArn(t *testing.T) {
	t.Parallel()

	username := GetIamCurrentUserArn(t)
	assert.Regexp(t, "^arn:aws:iam::[0-9]{12}:user/.+$", username)
}

func TestAssertIAMPolicyIsCorrect(t *testing.T) {
	t.Parallel()

	region := GetRandomRegion(t, nil, nil)

	policyDocument := "{\n  \"Version\": \"2012-10-17\",\n  \"Statement\": [\n    {\n      \"Sid\": \"Stmt1530709892083\",\n      \"Action\": \"*\",\n      \"Effect\": \"Allow\",\n      \"Resource\": \"*\"\n    }\n  ]\n}"

	iamClient, err := NewIamClientE(t, region)
	if err != nil {
		t.Fatal(err)
	}

	input := &iam.CreatePolicyInput{
		PolicyName:     aws.String(strings.ToLower(random.UniqueId())),
		PolicyDocument: aws.String(policyDocument),
	}

	policy, err := iamClient.CreatePolicy(input)
	if err != nil {
		t.Fatal(err)
	}

	AssertIAMPolicyIsEqual(t, region, *policy.Policy.Arn, policyDocument)
}
