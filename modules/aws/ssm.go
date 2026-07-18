package aws

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/gruntwork-io/terratest/modules/core/v2/logger"
	"github.com/gruntwork-io/terratest/modules/core/v2/retry"
	"github.com/gruntwork-io/terratest/modules/core/v2/testing"
	"github.com/stretchr/testify/require"
)

// ssmRetryInterval is the time between retries when waiting for SSM operations.
const ssmRetryInterval = 2 * time.Second

// GetParameterContextE retrieves the latest version of SSM Parameter at keyName with decryption.
// The ctx parameter supports cancellation and timeouts.
func GetParameterContextE(t testing.TestingT, ctx context.Context, awsRegion string, keyName string) (string, error) {
	ssmClient, err := NewSsmClientContextE(t, ctx, awsRegion)
	if err != nil {
		return "", err
	}

	return GetParameterWithClientContextE(t, ctx, ssmClient, keyName)
}

// GetParameterContext retrieves the latest version of SSM Parameter at keyName with decryption.
// The ctx parameter supports cancellation and timeouts.
func GetParameterContext(t testing.TestingT, ctx context.Context, awsRegion string, keyName string) string {
	t.Helper()
	keyValue, err := GetParameterContextE(t, ctx, awsRegion, keyName)
	require.NoError(t, err)

	return keyValue
}

// GetParameterWithClientContextE retrieves the latest version of SSM Parameter at keyName with decryption with the ability to provide the SSM client.
// The ctx parameter supports cancellation and timeouts.
func GetParameterWithClientContextE(t testing.TestingT, ctx context.Context, client *ssm.Client, keyName string) (string, error) {
	resp, err := client.GetParameter(ctx, &ssm.GetParameterInput{Name: aws.String(keyName), WithDecryption: aws.Bool(true)})
	if err != nil {
		return "", err
	}

	parameter := *resp.Parameter

	return *parameter.Value, nil
}

// PutParameterContextE creates new version of SSM Parameter at keyName with keyValue as SecureString.
// The ctx parameter supports cancellation and timeouts.
func PutParameterContextE(t testing.TestingT, ctx context.Context, awsRegion string, keyName string, keyDescription string, keyValue string) (int64, error) {
	ssmClient, err := NewSsmClientContextE(t, ctx, awsRegion)
	if err != nil {
		return 0, err
	}

	return PutParameterWithClientContextE(t, ctx, ssmClient, keyName, keyDescription, keyValue)
}

// PutParameterContext creates new version of SSM Parameter at keyName with keyValue as SecureString.
// The ctx parameter supports cancellation and timeouts.
func PutParameterContext(t testing.TestingT, ctx context.Context, awsRegion string, keyName string, keyDescription string, keyValue string) int64 {
	t.Helper()
	version, err := PutParameterContextE(t, ctx, awsRegion, keyName, keyDescription, keyValue)
	require.NoError(t, err)

	return version
}

// PutParameterWithClientContextE creates new version of SSM Parameter at keyName with keyValue as SecureString with the ability to provide the SSM client.
// The ctx parameter supports cancellation and timeouts.
func PutParameterWithClientContextE(t testing.TestingT, ctx context.Context, client *ssm.Client, keyName string, keyDescription string, keyValue string) (int64, error) {
	resp, err := client.PutParameter(ctx, &ssm.PutParameterInput{
		Name:        aws.String(keyName),
		Description: aws.String(keyDescription),
		Value:       aws.String(keyValue),
		Type:        types.ParameterTypeSecureString,
	})
	if err != nil {
		return 0, err
	}

	return resp.Version, nil
}

// DeleteParameterContextE deletes all versions of SSM Parameter at keyName.
// The ctx parameter supports cancellation and timeouts.
func DeleteParameterContextE(t testing.TestingT, ctx context.Context, awsRegion string, keyName string) error {
	ssmClient, err := NewSsmClientContextE(t, ctx, awsRegion)
	if err != nil {
		return err
	}

	return DeleteParameterWithClientContextE(t, ctx, ssmClient, keyName)
}

// DeleteParameterContext deletes all versions of SSM Parameter at keyName.
// The ctx parameter supports cancellation and timeouts.
func DeleteParameterContext(t testing.TestingT, ctx context.Context, awsRegion string, keyName string) {
	t.Helper()
	err := DeleteParameterContextE(t, ctx, awsRegion, keyName)
	require.NoError(t, err)
}

// DeleteParameterWithClientContextE deletes all versions of SSM Parameter at keyName with the ability to provide the SSM client.
// The ctx parameter supports cancellation and timeouts.
func DeleteParameterWithClientContextE(t testing.TestingT, ctx context.Context, client *ssm.Client, keyName string) error {
	_, err := client.DeleteParameter(ctx, &ssm.DeleteParameterInput{Name: aws.String(keyName)})
	if err != nil {
		return err
	}

	return nil
}

// NewSsmClientContextE creates an SSM client.
// The ctx parameter supports cancellation and timeouts.
func NewSsmClientContextE(t testing.TestingT, ctx context.Context, region string) (*ssm.Client, error) {
	sess, err := NewAuthenticatedSessionContext(ctx, region)
	if err != nil {
		return nil, err
	}

	return ssm.NewFromConfig(*sess), nil
}

// NewSsmClientContext creates an SSM client.
// The ctx parameter supports cancellation and timeouts.
func NewSsmClientContext(t testing.TestingT, ctx context.Context, region string) *ssm.Client {
	t.Helper()
	client, err := NewSsmClientContextE(t, ctx, region)
	require.NoError(t, err)

	return client
}

// WaitForSsmInstanceContextE waits until the instance get registered to the SSM inventory.
// The ctx parameter supports cancellation and timeouts.
func WaitForSsmInstanceContextE(t testing.TestingT, ctx context.Context, awsRegion, instanceID string, timeout time.Duration) error {
	client, err := NewSsmClientContextE(t, ctx, awsRegion)
	if err != nil {
		return err
	}

	return WaitForSsmInstanceWithClientContextE(t, ctx, client, instanceID, timeout)
}

// WaitForSsmInstanceContext waits until the instance get registered to the SSM inventory.
// The ctx parameter supports cancellation and timeouts.
func WaitForSsmInstanceContext(t testing.TestingT, ctx context.Context, awsRegion, instanceID string, timeout time.Duration) {
	t.Helper()
	err := WaitForSsmInstanceContextE(t, ctx, awsRegion, instanceID, timeout)
	require.NoError(t, err)
}

// WaitForSsmInstanceWithClientContextE waits until the instance get registered to the SSM inventory with the ability to provide the SSM client.
// The ctx parameter supports cancellation and timeouts.
func WaitForSsmInstanceWithClientContextE(t testing.TestingT, ctx context.Context, client *ssm.Client, instanceID string, timeout time.Duration) error {
	timeBetweenRetries := ssmRetryInterval
	maxRetries := int(timeout.Seconds() / timeBetweenRetries.Seconds())
	description := fmt.Sprintf("Waiting for %s to appear in the SSM inventory", instanceID)

	input := &ssm.GetInventoryInput{
		Filters: []types.InventoryFilter{
			{
				Key:    aws.String("AWS:InstanceInformation.InstanceId"),
				Type:   types.InventoryQueryOperatorTypeEqual,
				Values: []string{instanceID},
			},
		},
	}

	_, err := retry.DoWithRetryContextE(t, ctx, description, maxRetries, timeBetweenRetries, func() (string, error) {
		resp, err := client.GetInventory(ctx, input)
		if err != nil {
			return "", err
		}

		if len(resp.Entities) != 1 {
			return "", fmt.Errorf("%s is not in the SSM inventory", instanceID)
		}

		return "", nil
	})

	return err
}

// CheckSsmCommandContextE checks that you can run the given command on the given instance through AWS SSM. Returns the result and an error if one occurs.
// The ctx parameter supports cancellation and timeouts.
func CheckSsmCommandContextE(t testing.TestingT, ctx context.Context, awsRegion, instanceID, command string, timeout time.Duration) (*CommandOutput, error) {
	return CheckSsmCommandWithDocumentContextE(t, ctx, awsRegion, instanceID, command, "AWS-RunShellScript", timeout)
}

// CheckSsmCommandContext checks that you can run the given command on the given instance through AWS SSM.
// The ctx parameter supports cancellation and timeouts.
func CheckSsmCommandContext(t testing.TestingT, ctx context.Context, awsRegion, instanceID, command string, timeout time.Duration) *CommandOutput {
	t.Helper()
	return CheckSsmCommandWithDocumentContext(t, ctx, awsRegion, instanceID, command, "AWS-RunShellScript", timeout)
}

// CommandOutput contains the result of the SSM command.
type CommandOutput struct {
	Stdout   string
	Stderr   string
	ExitCode int64
}

// CheckSSMCommandWithClientContextE checks that you can run the given command on the given instance through AWS SSM with the ability to provide the SSM client. Returns the result and an error if one occurs.
// The ctx parameter supports cancellation and timeouts.
func CheckSSMCommandWithClientContextE(t testing.TestingT, ctx context.Context, client *ssm.Client, instanceID, command string, timeout time.Duration) (*CommandOutput, error) {
	return CheckSSMCommandWithClientWithDocumentContextE(t, ctx, client, instanceID, command, "AWS-RunShellScript", timeout)
}

// CheckSsmCommandWithDocumentContextE checks that you can run the given command on the given instance through AWS SSM with specified Command Doc type. Returns the result and an error if one occurs.
// The ctx parameter supports cancellation and timeouts.
func CheckSsmCommandWithDocumentContextE(t testing.TestingT, ctx context.Context, awsRegion, instanceID, command string, commandDocName string, timeout time.Duration) (*CommandOutput, error) {
	logger.Default.Logf(t, "Running command '%s' on EC2 instance with ID '%s'", command, instanceID)

	client, err := NewSsmClientContextE(t, ctx, awsRegion)
	if err != nil {
		return nil, err
	}

	return CheckSSMCommandWithClientWithDocumentContextE(t, ctx, client, instanceID, command, commandDocName, timeout)
}

// CheckSsmCommandWithDocumentContext checks that you can run the given command on the given instance through AWS SSM with specified Command Doc type.
// The ctx parameter supports cancellation and timeouts.
func CheckSsmCommandWithDocumentContext(t testing.TestingT, ctx context.Context, awsRegion, instanceID, command string, commandDocName string, timeout time.Duration) *CommandOutput {
	t.Helper()
	result, err := CheckSsmCommandWithDocumentContextE(t, ctx, awsRegion, instanceID, command, commandDocName, timeout)
	require.NoErrorf(t, err, "failed to execute '%s' on %s (%v):]\n  stdout: %#v\n  stderr: %#v", command, instanceID, err, result.Stdout, result.Stderr)

	return result
}

// CheckSSMCommandWithClientWithDocumentContextE checks that you can run the given command on the given instance through AWS SSM with the ability to provide the SSM client with specified Command Doc type. Returns the result and an error if one occurs.
// The ctx parameter supports cancellation and timeouts.
func CheckSSMCommandWithClientWithDocumentContextE(t testing.TestingT, ctx context.Context, client *ssm.Client, instanceID, command string, commandDocName string, timeout time.Duration) (*CommandOutput, error) {
	timeBetweenRetries := ssmRetryInterval
	maxRetries := int(timeout.Seconds() / timeBetweenRetries.Seconds())

	resp, err := client.SendCommand(
		ctx,
		&ssm.SendCommandInput{
			Comment:      aws.String("Terratest SSM"),
			DocumentName: aws.String(commandDocName),
			InstanceIds:  []string{instanceID},
			Parameters: map[string][]string{
				"commands": {command},
			},
		},
	)
	if err != nil {
		return nil, err
	}

	description := "Waiting for the result of the command"
	retryableErrors := map[string]string{
		"InvocationDoesNotExist": "InvocationDoesNotExist",
		"bad status: Pending":    "bad status: Pending",
		"bad status: InProgress": "bad status: InProgress",
		"bad status: Delayed":    "bad status: Delayed",
	}

	result := &CommandOutput{}

	_, err = retry.DoWithRetryableErrorsContextE(t, ctx, description, retryableErrors, maxRetries, timeBetweenRetries, func() (string, error) {
		resp, err := client.GetCommandInvocation(ctx, &ssm.GetCommandInvocationInput{
			CommandId:  resp.Command.CommandId,
			InstanceId: &instanceID,
		})
		if err != nil {
			return "", err
		}

		result.Stderr = aws.ToString(resp.StandardErrorContent)
		result.Stdout = aws.ToString(resp.StandardOutputContent)
		result.ExitCode = int64(resp.ResponseCode)

		status := resp.Status

		if status == types.CommandInvocationStatusSuccess {
			return "", nil
		}

		if status == types.CommandInvocationStatusFailed {
			return "", fmt.Errorf("%s", aws.ToString(resp.StatusDetails))
		}

		return "", fmt.Errorf("bad status: %s", status)
	})
	if err != nil {
		var actualErr retry.FatalError
		if errors.As(err, &actualErr) {
			return result, actualErr.Underlying
		}

		return result, fmt.Errorf("unexpected error: %w", err)
	}

	return result, nil
}
