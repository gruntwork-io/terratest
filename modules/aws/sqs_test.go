package aws_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	terraaws "github.com/gruntwork-io/terratest/modules/aws/v2"
	"github.com/gruntwork-io/terratest/modules/core/v2/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSqsQueueMethods(t *testing.T) {
	t.Parallel()

	region := terraaws.GetRandomStableRegionContext(t, context.Background(), nil, nil)
	uniqueID := random.UniqueID()
	namePrefix := "sqs-queue-test-" + uniqueID

	url := terraaws.CreateRandomQueueContext(t, context.Background(), region, namePrefix)
	defer deleteQueue(t, region, url)

	assert.True(t, queueExists(t, region, url))

	message := "test-message-" + uniqueID
	timeoutSec := 20

	terraaws.SendMessageToQueueContext(t, context.Background(), region, url, message)

	firstResponse := terraaws.WaitForQueueMessageContext(t, context.Background(), region, url, timeoutSec)
	require.NoError(t, firstResponse.Error)
	assert.Equal(t, message, firstResponse.MessageBody)

	terraaws.DeleteMessageFromQueueContext(t, context.Background(), region, url, firstResponse.ReceiptHandle)

	secondResponse := terraaws.WaitForQueueMessageContext(t, context.Background(), region, url, timeoutSec)
	assert.Error(t, secondResponse.Error, terraaws.ReceiveMessageTimeout{QueueUrl: url, TimeoutSec: timeoutSec})
}

func TestFifoSqsQueueMethods(t *testing.T) {
	t.Parallel()

	region := terraaws.GetRandomStableRegionContext(t, context.Background(), nil, nil)
	uniqueID := random.UniqueID()
	namePrefix := "sqs-queue-test-" + uniqueID
	fifoMessageGroupID := "g1"

	url := terraaws.CreateRandomFifoQueueContext(t, context.Background(), region, namePrefix)
	defer deleteQueue(t, region, url)

	assert.True(t, queueExists(t, region, url))

	message := "test-message-" + uniqueID
	timeoutSec := 20

	terraaws.SendMessageFifoToQueueContext(t, context.Background(), region, url, message, fifoMessageGroupID)

	firstResponse := terraaws.WaitForQueueMessageContext(t, context.Background(), region, url, timeoutSec)
	require.NoError(t, firstResponse.Error)
	assert.Equal(t, message, firstResponse.MessageBody)

	terraaws.DeleteMessageFromQueueContext(t, context.Background(), region, url, firstResponse.ReceiptHandle)

	secondResponse := terraaws.WaitForQueueMessageContext(t, context.Background(), region, url, timeoutSec)
	assert.Error(t, secondResponse.Error, terraaws.ReceiveMessageTimeout{QueueUrl: url, TimeoutSec: timeoutSec})
}

func queueExists(t *testing.T, region string, url string) bool {
	t.Helper()

	sqsClient := terraaws.NewSqsClientContext(t, context.Background(), region)

	input := sqs.GetQueueAttributesInput{QueueUrl: aws.String(url)}

	if _, err := sqsClient.GetQueueAttributes(context.Background(), &input); err != nil {
		if strings.Contains(err.Error(), "NonExistentQueue") {
			return false
		}

		t.Fatal(err)
	}

	return true
}

func deleteQueue(t *testing.T, region string, url string) {
	t.Helper()

	terraaws.DeleteQueueContext(t, context.Background(), region, url)
	assert.False(t, queueExists(t, region, url))
}
