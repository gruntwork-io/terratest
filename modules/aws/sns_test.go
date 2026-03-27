package aws_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	terraaws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndDeleteSnsTopic(t *testing.T) {
	t.Parallel()

	region := terraaws.GetRandomStableRegion(t, nil, nil)
	uniqueID := random.UniqueID()
	name := "test-sns-topic-" + uniqueID

	arn := terraaws.CreateSnsTopic(t, region, name)
	defer deleteTopic(t, region, arn)

	assert.True(t, snsTopicExists(t, region, arn))
}

func snsTopicExists(t *testing.T, region string, arn string) bool {
	t.Helper()

	snsClient := terraaws.NewSnsClient(t, region)

	input := sns.GetTopicAttributesInput{TopicArn: aws.String(arn)}

	if _, err := snsClient.GetTopicAttributes(context.Background(), &input); err != nil {
		if strings.Contains(err.Error(), "NotFound") {
			return false
		}

		t.Fatal(err)
	}

	return true
}

func deleteTopic(t *testing.T, region string, arn string) {
	t.Helper()

	terraaws.DeleteSNSTopic(t, region, arn)
	assert.False(t, snsTopicExists(t, region, arn))
}
