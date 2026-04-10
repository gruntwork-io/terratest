package aws

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// sqsMaxWaitTimeSeconds is the maximum number of seconds to wait for a message on a single SQS receive call.
const sqsMaxWaitTimeSeconds = 20

// CreateRandomQueueContextE creates a new SQS queue with a random name that starts with the given prefix and return the queue URL.
// The ctx parameter supports cancellation and timeouts.
func CreateRandomQueueContextE(t testing.TestingT, ctx context.Context, awsRegion string, prefix string) (string, error) {
	logger.Default.Logf(t, "Creating randomly named SQS queue with prefix %s", prefix)

	sqsClient, err := NewSqsClientContextE(t, ctx, awsRegion)
	if err != nil {
		return "", err
	}

	channel, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	channelName := fmt.Sprintf("%s-%s", prefix, channel.String())

	queue, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(channelName),
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(queue.QueueUrl), nil
}

// CreateRandomQueueContext creates a new SQS queue with a random name that starts with the given prefix and return the queue URL.
// The ctx parameter supports cancellation and timeouts.
func CreateRandomQueueContext(t testing.TestingT, ctx context.Context, awsRegion string, prefix string) string {
	t.Helper()

	url, err := CreateRandomQueueContextE(t, ctx, awsRegion, prefix)
	require.NoError(t, err)

	return url
}

// CreateRandomQueue creates a new SQS queue with a random name that starts with the given prefix and return the queue URL.
//
// Deprecated: Use [CreateRandomQueueContext] instead.
func CreateRandomQueue(t testing.TestingT, awsRegion string, prefix string) string {
	t.Helper()

	return CreateRandomQueueContext(t, context.Background(), awsRegion, prefix)
}

// CreateRandomQueueE creates a new SQS queue with a random name that starts with the given prefix and return the queue URL.
//
// Deprecated: Use [CreateRandomQueueContextE] instead.
func CreateRandomQueueE(t testing.TestingT, awsRegion string, prefix string) (string, error) {
	return CreateRandomQueueContextE(t, context.Background(), awsRegion, prefix)
}

// CreateRandomFifoQueueContextE creates a new FIFO SQS queue with a random name that starts with the given prefix and return the queue URL.
// The ctx parameter supports cancellation and timeouts.
func CreateRandomFifoQueueContextE(t testing.TestingT, ctx context.Context, awsRegion string, prefix string) (string, error) {
	logger.Default.Logf(t, "Creating randomly named FIFO SQS queue with prefix %s", prefix)

	sqsClient, err := NewSqsClientContextE(t, ctx, awsRegion)
	if err != nil {
		return "", err
	}

	channel, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	channelName := fmt.Sprintf("%s-%s.fifo", prefix, channel.String())

	queue, err := sqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(channelName),
		Attributes: map[string]string{
			"ContentBasedDeduplication": "true",
			"FifoQueue":                 "true",
		},
	})
	if err != nil {
		return "", err
	}

	return aws.ToString(queue.QueueUrl), nil
}

// CreateRandomFifoQueueContext creates a new FIFO SQS queue with a random name that starts with the given prefix and return the queue URL.
// The ctx parameter supports cancellation and timeouts.
func CreateRandomFifoQueueContext(t testing.TestingT, ctx context.Context, awsRegion string, prefix string) string {
	t.Helper()

	url, err := CreateRandomFifoQueueContextE(t, ctx, awsRegion, prefix)
	require.NoError(t, err)

	return url
}

// CreateRandomFifoQueue creates a new FIFO SQS queue with a random name that starts with the given prefix and return the queue URL.
//
// Deprecated: Use [CreateRandomFifoQueueContext] instead.
func CreateRandomFifoQueue(t testing.TestingT, awsRegion string, prefix string) string {
	t.Helper()

	return CreateRandomFifoQueueContext(t, context.Background(), awsRegion, prefix)
}

// CreateRandomFifoQueueE creates a new FIFO SQS queue with a random name that starts with the given prefix and return the queue URL.
//
// Deprecated: Use [CreateRandomFifoQueueContextE] instead.
func CreateRandomFifoQueueE(t testing.TestingT, awsRegion string, prefix string) (string, error) {
	return CreateRandomFifoQueueContextE(t, context.Background(), awsRegion, prefix)
}

// DeleteQueueContextE deletes the SQS queue with the given URL.
// The ctx parameter supports cancellation and timeouts.
func DeleteQueueContextE(t testing.TestingT, ctx context.Context, awsRegion string, queueURL string) error {
	logger.Default.Logf(t, "Deleting SQS Queue %s", queueURL)

	sqsClient, err := NewSqsClientContextE(t, ctx, awsRegion)
	if err != nil {
		return err
	}

	_, err = sqsClient.DeleteQueue(ctx, &sqs.DeleteQueueInput{
		QueueUrl: aws.String(queueURL),
	})

	return err
}

// DeleteQueueContext deletes the SQS queue with the given URL.
// The ctx parameter supports cancellation and timeouts.
func DeleteQueueContext(t testing.TestingT, ctx context.Context, awsRegion string, queueURL string) {
	t.Helper()

	err := DeleteQueueContextE(t, ctx, awsRegion, queueURL)
	require.NoError(t, err)
}

// DeleteQueue deletes the SQS queue with the given URL.
//
// Deprecated: Use [DeleteQueueContext] instead.
func DeleteQueue(t testing.TestingT, awsRegion string, queueURL string) {
	t.Helper()

	DeleteQueueContext(t, context.Background(), awsRegion, queueURL)
}

// DeleteQueueE deletes the SQS queue with the given URL.
//
// Deprecated: Use [DeleteQueueContextE] instead.
func DeleteQueueE(t testing.TestingT, awsRegion string, queueURL string) error {
	return DeleteQueueContextE(t, context.Background(), awsRegion, queueURL)
}

// DeleteMessageFromQueueContextE deletes the message with the given receipt from the SQS queue with the given URL.
// The ctx parameter supports cancellation and timeouts.
func DeleteMessageFromQueueContextE(t testing.TestingT, ctx context.Context, awsRegion string, queueURL string, receipt string) error {
	logger.Default.Logf(t, "Deleting message from queue %s (%s)", queueURL, receipt)

	sqsClient, err := NewSqsClientContextE(t, ctx, awsRegion)
	if err != nil {
		return err
	}

	_, err = sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		ReceiptHandle: &receipt,
		QueueUrl:      &queueURL,
	})

	return err
}

// DeleteMessageFromQueueContext deletes the message with the given receipt from the SQS queue with the given URL.
// The ctx parameter supports cancellation and timeouts.
func DeleteMessageFromQueueContext(t testing.TestingT, ctx context.Context, awsRegion string, queueURL string, receipt string) {
	t.Helper()

	err := DeleteMessageFromQueueContextE(t, ctx, awsRegion, queueURL, receipt)
	require.NoError(t, err)
}

// DeleteMessageFromQueue deletes the message with the given receipt from the SQS queue with the given URL.
//
// Deprecated: Use [DeleteMessageFromQueueContext] instead.
func DeleteMessageFromQueue(t testing.TestingT, awsRegion string, queueURL string, receipt string) {
	t.Helper()

	DeleteMessageFromQueueContext(t, context.Background(), awsRegion, queueURL, receipt)
}

// DeleteMessageFromQueueE deletes the message with the given receipt from the SQS queue with the given URL.
//
// Deprecated: Use [DeleteMessageFromQueueContextE] instead.
func DeleteMessageFromQueueE(t testing.TestingT, awsRegion string, queueURL string, receipt string) error {
	return DeleteMessageFromQueueContextE(t, context.Background(), awsRegion, queueURL, receipt)
}

// SendMessageToQueueContextE sends the given message to the SQS queue with the given URL.
// The ctx parameter supports cancellation and timeouts.
func SendMessageToQueueContextE(t testing.TestingT, ctx context.Context, awsRegion string, queueURL string, message string) error {
	logger.Default.Logf(t, "Sending message %s to queue %s", message, queueURL)

	sqsClient, err := NewSqsClientContextE(t, ctx, awsRegion)
	if err != nil {
		return err
	}

	res, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: &message,
		QueueUrl:    &queueURL,
	})
	if err != nil {
		if strings.Contains(err.Error(), "AWS.SimpleQueueService.NonExistentQueue") {
			logger.Default.Logf(t, "WARN: Client has stopped listening on queue %s", queueURL)

			return nil
		}

		return err
	}

	logger.Default.Logf(t, "Message id %s sent to queue %s", aws.ToString(res.MessageId), queueURL)

	return nil
}

// SendMessageToQueueContext sends the given message to the SQS queue with the given URL.
// The ctx parameter supports cancellation and timeouts.
func SendMessageToQueueContext(t testing.TestingT, ctx context.Context, awsRegion string, queueURL string, message string) {
	t.Helper()

	err := SendMessageToQueueContextE(t, ctx, awsRegion, queueURL, message)
	require.NoError(t, err)
}

// SendMessageToQueue sends the given message to the SQS queue with the given URL.
//
// Deprecated: Use [SendMessageToQueueContext] instead.
func SendMessageToQueue(t testing.TestingT, awsRegion string, queueURL string, message string) {
	t.Helper()

	SendMessageToQueueContext(t, context.Background(), awsRegion, queueURL, message)
}

// SendMessageToQueueE sends the given message to the SQS queue with the given URL.
//
// Deprecated: Use [SendMessageToQueueContextE] instead.
func SendMessageToQueueE(t testing.TestingT, awsRegion string, queueURL string, message string) error {
	return SendMessageToQueueContextE(t, context.Background(), awsRegion, queueURL, message)
}

// SendMessageFifoToQueueContextE sends the given message to the FIFO SQS queue with the given URL.
// The ctx parameter supports cancellation and timeouts.
func SendMessageFifoToQueueContextE(t testing.TestingT, ctx context.Context, awsRegion string, queueURL string, message string, messageGroupID string) error {
	logger.Default.Logf(t, "Sending message %s to queue %s", message, queueURL)

	sqsClient, err := NewSqsClientContextE(t, ctx, awsRegion)
	if err != nil {
		return err
	}

	res, err := sqsClient.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody:    &message,
		QueueUrl:       &queueURL,
		MessageGroupId: &messageGroupID,
	})
	if err != nil {
		if strings.Contains(err.Error(), "AWS.SimpleQueueService.NonExistentQueue") {
			logger.Default.Logf(t, "WARN: Client has stopped listening on queue %s", queueURL)

			return nil
		}

		return err
	}

	logger.Default.Logf(t, "Message id %s sent to queue %s", aws.ToString(res.MessageId), queueURL)

	return nil
}

// SendMessageFifoToQueueContext sends the given message to the FIFO SQS queue with the given URL.
// The ctx parameter supports cancellation and timeouts.
func SendMessageFifoToQueueContext(t testing.TestingT, ctx context.Context, awsRegion string, queueURL string, message string, messageGroupID string) {
	t.Helper()

	err := SendMessageFifoToQueueContextE(t, ctx, awsRegion, queueURL, message, messageGroupID)
	require.NoError(t, err)
}

// SendMessageFifoToQueue sends the given message to the FIFO SQS queue with the given URL.
//
// Deprecated: Use [SendMessageFifoToQueueContext] instead.
func SendMessageFifoToQueue(t testing.TestingT, awsRegion string, queueURL string, message string, messageGroupID string) {
	t.Helper()

	SendMessageFifoToQueueContext(t, context.Background(), awsRegion, queueURL, message, messageGroupID)
}

// SendMessageToFifoQueueE sends the given message to the FIFO SQS queue with the given URL.
//
// Deprecated: Use [SendMessageFifoToQueueContextE] instead.
func SendMessageToFifoQueueE(t testing.TestingT, awsRegion string, queueURL string, message string, messageGroupID string) error {
	return SendMessageFifoToQueueContextE(t, context.Background(), awsRegion, queueURL, message, messageGroupID)
}

// QueueMessageResponse contains a queue message.
type QueueMessageResponse struct {
	Error         error
	ReceiptHandle string
	MessageBody   string
}

// WaitForQueueMessageContext waits to receive a message from on the queueURL. Since the API only allows us to wait a max 20 seconds for a new
// message to arrive, we must loop TIMEOUT/20 number of times to be able to wait for a total of TIMEOUT seconds.
// The ctx parameter supports cancellation and timeouts.
func WaitForQueueMessageContext(t testing.TestingT, ctx context.Context, awsRegion string, queueURL string, timeout int) QueueMessageResponse {
	sqsClient, err := NewSqsClientContextE(t, ctx, awsRegion)
	if err != nil {
		return QueueMessageResponse{Error: err}
	}

	cycles := timeout
	cycleLength := 1

	if timeout >= sqsMaxWaitTimeSeconds {
		cycleLength = sqsMaxWaitTimeSeconds
		cycles = timeout / cycleLength
	}

	for i := 0; i < cycles; i++ {
		logger.Default.Logf(t, "Waiting for message on %s (%ss)", queueURL, strconv.Itoa(i*cycleLength))

		result, err := sqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
			QueueUrl:                    aws.String(queueURL),
			MessageSystemAttributeNames: []types.MessageSystemAttributeName{types.MessageSystemAttributeNameSentTimestamp},
			MaxNumberOfMessages:         int32(1),
			MessageAttributeNames:       []string{"All"},
			WaitTimeSeconds:             int32(cycleLength),
		})
		if err != nil {
			return QueueMessageResponse{Error: err}
		}

		if len(result.Messages) > 0 {
			logger.Default.Logf(t, "Message %s received on %s", *result.Messages[0].MessageId, queueURL)

			return QueueMessageResponse{ReceiptHandle: *result.Messages[0].ReceiptHandle, MessageBody: *result.Messages[0].Body}
		}
	}

	return QueueMessageResponse{Error: ReceiveMessageTimeout{QueueUrl: queueURL, TimeoutSec: timeout}}
}

// WaitForQueueMessage waits to receive a message from on the queueURL. Since the API only allows us to wait a max 20 seconds for a new
// message to arrive, we must loop TIMEOUT/20 number of times to be able to wait for a total of TIMEOUT seconds
//
// Deprecated: Use [WaitForQueueMessageContext] instead.
func WaitForQueueMessage(t testing.TestingT, awsRegion string, queueURL string, timeout int) QueueMessageResponse {
	return WaitForQueueMessageContext(t, context.Background(), awsRegion, queueURL, timeout)
}

// NewSqsClientContextE creates a new SQS client.
// The ctx parameter supports cancellation and timeouts.
func NewSqsClientContextE(t testing.TestingT, ctx context.Context, region string) (*sqs.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return sqs.NewFromConfig(*sess), nil
}

// NewSqsClientContext creates a new SQS client.
// The ctx parameter supports cancellation and timeouts.
func NewSqsClientContext(t testing.TestingT, ctx context.Context, region string) *sqs.Client {
	t.Helper()

	client, err := NewSqsClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// NewSqsClient creates a new SQS client.
//
// Deprecated: Use [NewSqsClientContext] instead.
func NewSqsClient(t testing.TestingT, region string) *sqs.Client {
	t.Helper()

	return NewSqsClientContext(t, context.Background(), region)
}

// NewSqsClientE creates a new SQS client.
//
// Deprecated: Use [NewSqsClientContextE] instead.
func NewSqsClientE(t testing.TestingT, region string) (*sqs.Client, error) {
	return NewSqsClientContextE(t, context.Background(), region)
}

// ReceiveMessageTimeout is an error that occurs if receiving a message times out.
type ReceiveMessageTimeout struct {
	QueueUrl   string //nolint:staticcheck,revive // preserving existing field name
	TimeoutSec int
}

func (err ReceiveMessageTimeout) Error() string {
	return fmt.Sprintf("Failed to receive messages on %s within %s seconds", err.QueueUrl, strconv.Itoa(err.TimeoutSec))
}
