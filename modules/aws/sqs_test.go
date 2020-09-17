package aws

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
)

func TestSqsQueueMethods(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	uniqueID := random.UniqueId()
	namePrefix := fmt.Sprintf("sqs-queue-test-%s", uniqueID)

	url := CreateRandomQueue(t, region, namePrefix)
	defer deleteQueue(t, region, url)

	assert.True(t, queueExists(t, region, url))

	message := fmt.Sprintf("test-message-%s", uniqueID)
	timeoutSec := 20

	SendMessageToQueue(t, region, url, message)

	firstResponse := WaitForQueueMessage(t, region, url, timeoutSec)
	assert.NoError(t, firstResponse.Error)
	assert.Equal(t, message, firstResponse.MessageBody)

	DeleteMessageFromQueue(t, region, url, firstResponse.ReceiptHandle)

	secondResponse := WaitForQueueMessage(t, region, url, timeoutSec)
	assert.Error(t, secondResponse.Error, ReceiveMessageTimeout{QueueUrl: url, TimeoutSec: timeoutSec})
}

func TestFifoSqsQueueMethods(t *testing.T) {
	t.Parallel()

	region := GetRandomStableRegion(t, nil, nil)
	uniqueID := random.UniqueId()
	namePrefix := fmt.Sprintf("sqs-queue-test-%s", uniqueID)
	fifoMessageGroupID := "g1"

	url := CreateRandomFifoQueue(t, region, namePrefix)
	defer deleteQueue(t, region, url)

	assert.True(t, queueExists(t, region, url))

	message := fmt.Sprintf("test-message-%s", uniqueID)
	timeoutSec := 20

	SendMessageFifoToQueue(t, region, url, message, fifoMessageGroupID)

	firstResponse := WaitForQueueMessage(t, region, url, timeoutSec)
	assert.NoError(t, firstResponse.Error)
	assert.Equal(t, message, firstResponse.MessageBody)

	DeleteMessageFromQueue(t, region, url, firstResponse.ReceiptHandle)

	secondResponse := WaitForQueueMessage(t, region, url, timeoutSec)
	assert.Error(t, secondResponse.Error, ReceiveMessageTimeout{QueueUrl: url, TimeoutSec: timeoutSec})
}

func queueExists(t *testing.T, region string, url string) bool {
	sqsClient := NewSqsClient(t, region)

	input := sqs.GetQueueAttributesInput{QueueUrl: aws.String(url)}

	if _, err := sqsClient.GetQueueAttributes(&input); err != nil {
		if strings.Contains(err.Error(), "NonExistentQueue") {
			return false
		}
		t.Fatal(err)
	}

	return true
}

func deleteQueue(t *testing.T, region string, url string) {
	DeleteQueue(t, region, url)
	assert.False(t, queueExists(t, region, url))
}
