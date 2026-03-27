package aws_test

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	terraaws "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSqsQueueMethods(t *testing.T) {
	t.Parallel()

	region := terraaws.GetRandomStableRegion(t, nil, nil)
	uniqueID := random.UniqueID()
	namePrefix := "sqs-queue-test-" + uniqueID

	url := terraaws.CreateRandomQueue(t, region, namePrefix)
	defer deleteQueue(t, region, url)

	assert.True(t, queueExists(t, region, url))

	message := "test-message-" + uniqueID
	timeoutSec := 20

	terraaws.SendMessageToQueue(t, region, url, message)

	firstResponse := terraaws.WaitForQueueMessage(t, region, url, timeoutSec)
	require.NoError(t, firstResponse.Error)
	assert.Equal(t, message, firstResponse.MessageBody)

	terraaws.DeleteMessageFromQueue(t, region, url, firstResponse.ReceiptHandle)

	secondResponse := terraaws.WaitForQueueMessage(t, region, url, timeoutSec)
	assert.Error(t, secondResponse.Error, terraaws.ReceiveMessageTimeout{QueueUrl: url, TimeoutSec: timeoutSec})
}

func TestFifoSqsQueueMethods(t *testing.T) {
	t.Parallel()

	region := terraaws.GetRandomStableRegion(t, nil, nil)
	uniqueID := random.UniqueID()
	namePrefix := "sqs-queue-test-" + uniqueID
	fifoMessageGroupID := "g1"

	url := terraaws.CreateRandomFifoQueue(t, region, namePrefix)
	defer deleteQueue(t, region, url)

	assert.True(t, queueExists(t, region, url))

	message := "test-message-" + uniqueID
	timeoutSec := 20

	terraaws.SendMessageFifoToQueue(t, region, url, message, fifoMessageGroupID)

	firstResponse := terraaws.WaitForQueueMessage(t, region, url, timeoutSec)
	require.NoError(t, firstResponse.Error)
	assert.Equal(t, message, firstResponse.MessageBody)

	terraaws.DeleteMessageFromQueue(t, region, url, firstResponse.ReceiptHandle)

	secondResponse := terraaws.WaitForQueueMessage(t, region, url, timeoutSec)
	assert.Error(t, secondResponse.Error, terraaws.ReceiveMessageTimeout{QueueUrl: url, TimeoutSec: timeoutSec})
}

func queueExists(t *testing.T, region string, url string) bool {
	t.Helper()

	sqsClient := terraaws.NewSqsClient(t, region)

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

	terraaws.DeleteQueue(t, region, url)
	assert.False(t, queueExists(t, region, url))
}
