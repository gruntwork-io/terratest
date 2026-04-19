package aws

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
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

// GetParameter retrieves the latest version of SSM Parameter at keyName with decryption.
//
// Deprecated: Use [GetParameterContext] instead.
func GetParameter(t testing.TestingT, awsRegion string, keyName string) string {
	t.Helper()
	return GetParameterContext(t, context.Background(), awsRegion, keyName)
}

// GetParameterE retrieves the latest version of SSM Parameter at keyName with decryption.
//
// Deprecated: Use [GetParameterContextE] instead.
func GetParameterE(t testing.TestingT, awsRegion string, keyName string) (string, error) {
	return GetParameterContextE(t, context.Background(), awsRegion, keyName)
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

// GetParameterWithClientE retrieves the latest version of SSM Parameter at keyName with decryption with the ability to provide the SSM client.
//
// Deprecated: Use [GetParameterWithClientContextE] instead.
func GetParameterWithClientE(t testing.TestingT, client *ssm.Client, keyName string) (string, error) {
	return GetParameterWithClientContextE(t, context.Background(), client, keyName)
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

// PutParameter creates new version of SSM Parameter at keyName with keyValue as SecureString.
//
// Deprecated: Use [PutParameterContext] instead.
func PutParameter(t testing.TestingT, awsRegion string, keyName string, keyDescription string, keyValue string) int64 {
	t.Helper()
	return PutParameterContext(t, context.Background(), awsRegion, keyName, keyDescription, keyValue)
}

// PutParameterE creates new version of SSM Parameter at keyName with keyValue as SecureString.
//
// Deprecated: Use [PutParameterContextE] instead.
func PutParameterE(t testing.TestingT, awsRegion string, keyName string, keyDescription string, keyValue string) (int64, error) {
	return PutParameterContextE(t, context.Background(), awsRegion, keyName, keyDescription, keyValue)
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

// PutParameterWithClientE creates new version of SSM Parameter at keyName with keyValue as SecureString with the ability to provide the SSM client.
//
// Deprecated: Use [PutParameterWithClientContextE] instead.
func PutParameterWithClientE(t testing.TestingT, client *ssm.Client, keyName string, keyDescription string, keyValue string) (int64, error) {
	return PutParameterWithClientContextE(t, context.Background(), client, keyName, keyDescription, keyValue)
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

// DeleteParameter deletes all versions of SSM Parameter at keyName.
//
// Deprecated: Use [DeleteParameterContext] instead.
func DeleteParameter(t testing.TestingT, awsRegion string, keyName string) {
	t.Helper()
	DeleteParameterContext(t, context.Background(), awsRegion, keyName)
}

// DeleteParameterE deletes all versions of SSM Parameter at keyName.
//
// Deprecated: Use [DeleteParameterContextE] instead.
func DeleteParameterE(t testing.TestingT, awsRegion string, keyName string) error {
	return DeleteParameterContextE(t, context.Background(), awsRegion, keyName)
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

// DeleteParameterWithClientE deletes all versions of SSM Parameter at keyName with the ability to provide the SSM client.
//
// Deprecated: Use [DeleteParameterWithClientContextE] instead.
func DeleteParameterWithClientE(t testing.TestingT, client *ssm.Client, keyName string) error {
	return DeleteParameterWithClientContextE(t, context.Background(), client, keyName)
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

// NewSsmClient creates an SSM client.
//
// Deprecated: Use [NewSsmClientContext] instead.
func NewSsmClient(t testing.TestingT, region string) *ssm.Client {
	t.Helper()
	return NewSsmClientContext(t, context.Background(), region)
}

// NewSsmClientE creates an SSM client.
//
// Deprecated: Use [NewSsmClientContextE] instead.
func NewSsmClientE(t testing.TestingT, region string) (*ssm.Client, error) {
	return NewSsmClientContextE(t, context.Background(), region)
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

// WaitForSsmInstance waits until the instance get registered to the SSM inventory.
//
// Deprecated: Use [WaitForSsmInstanceContext] instead.
func WaitForSsmInstance(t testing.TestingT, awsRegion, instanceID string, timeout time.Duration) {
	t.Helper()
	WaitForSsmInstanceContext(t, context.Background(), awsRegion, instanceID, timeout)
}

// WaitForSsmInstanceE waits until the instance get registered to the SSM inventory.
//
// Deprecated: Use [WaitForSsmInstanceContextE] instead.
func WaitForSsmInstanceE(t testing.TestingT, awsRegion, instanceID string, timeout time.Duration) error {
	return WaitForSsmInstanceContextE(t, context.Background(), awsRegion, instanceID, timeout)
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

// WaitForSsmInstanceWithClientE waits until the instance get registered to the SSM inventory with the ability to provide the SSM client.
//
// Deprecated: Use [WaitForSsmInstanceWithClientContextE] instead.
func WaitForSsmInstanceWithClientE(t testing.TestingT, client *ssm.Client, instanceID string, timeout time.Duration) error {
	return WaitForSsmInstanceWithClientContextE(t, context.Background(), client, instanceID, timeout)
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

// CheckSsmCommand checks that you can run the given command on the given instance through AWS SSM.
//
// Deprecated: Use [CheckSsmCommandContext] instead.
func CheckSsmCommand(t testing.TestingT, awsRegion, instanceID, command string, timeout time.Duration) *CommandOutput {
	t.Helper()
	return CheckSsmCommandContext(t, context.Background(), awsRegion, instanceID, command, timeout)
}

// CommandOutput contains the result of the SSM command.
type CommandOutput struct {
	Stdout   string
	Stderr   string
	ExitCode int64
}

// CheckSsmCommandE checks that you can run the given command on the given instance through AWS SSM. Returns the result and an error if one occurs.
//
// Deprecated: Use [CheckSsmCommandContextE] instead.
func CheckSsmCommandE(t testing.TestingT, awsRegion, instanceID, command string, timeout time.Duration) (*CommandOutput, error) {
	return CheckSsmCommandContextE(t, context.Background(), awsRegion, instanceID, command, timeout)
}

// CheckSSMCommandWithClientContextE checks that you can run the given command on the given instance through AWS SSM with the ability to provide the SSM client. Returns the result and an error if one occurs.
// The ctx parameter supports cancellation and timeouts.
func CheckSSMCommandWithClientContextE(t testing.TestingT, ctx context.Context, client *ssm.Client, instanceID, command string, timeout time.Duration) (*CommandOutput, error) {
	return CheckSSMCommandWithClientWithDocumentContextE(t, ctx, client, instanceID, command, "AWS-RunShellScript", timeout)
}

// CheckSSMCommandWithClientE checks that you can run the given command on the given instance through AWS SSM with the ability to provide the SSM client. Returns the result and an error if one occurs.
//
// Deprecated: Use [CheckSSMCommandWithClientContextE] instead.
func CheckSSMCommandWithClientE(t testing.TestingT, client *ssm.Client, instanceID, command string, timeout time.Duration) (*CommandOutput, error) {
	return CheckSSMCommandWithClientContextE(t, context.Background(), client, instanceID, command, timeout)
}

// CheckSsmCommandWithDocumentContextE checks that you can run the given command on the given instance through AWS SSM with specified Command Doc type. Returns the result and an error if one occurs.
// The ctx parameter supports cancellation and timeouts.
func CheckSsmCommandWithDocumentContextE(t testing.TestingT, ctx context.Context, awsRegion, instanceID, command string, commandDocName string, timeout time.Duration) (*CommandOutput, error) {
	logger.Default.Logf(t, "Running command '%s' on EC2 instance with ID '%s'", command, instanceID)

	// Now that we know the instance in the SSM inventory, we can send the command
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

// CheckSsmCommandWithDocument checks that you can run the given command on the given instance through AWS SSM with specified Command Doc type.
//
// Deprecated: Use [CheckSsmCommandWithDocumentContext] instead.
func CheckSsmCommandWithDocument(t testing.TestingT, awsRegion, instanceID, command string, commandDocName string, timeout time.Duration) *CommandOutput {
	t.Helper()
	return CheckSsmCommandWithDocumentContext(t, context.Background(), awsRegion, instanceID, command, commandDocName, timeout)
}

// CheckSsmCommandWithDocumentE checks that you can run the given command on the given instance through AWS SSM with specified Command Doc type. Returns the result and an error if one occurs.
//
// Deprecated: Use [CheckSsmCommandWithDocumentContextE] instead.
func CheckSsmCommandWithDocumentE(t testing.TestingT, awsRegion, instanceID, command string, commandDocName string, timeout time.Duration) (*CommandOutput, error) {
	return CheckSsmCommandWithDocumentContextE(t, context.Background(), awsRegion, instanceID, command, commandDocName, timeout)
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

	// Wait for the result
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

// CheckSSMCommandWithClientWithDocumentE checks that you can run the given command on the given instance through AWS SSM with the ability to provide the SSM client with specified Command Doc type. Returns the result and an error if one occurs.
//
// Deprecated: Use [CheckSSMCommandWithClientWithDocumentContextE] instead.
func CheckSSMCommandWithClientWithDocumentE(t testing.TestingT, client *ssm.Client, instanceID, command string, commandDocName string, timeout time.Duration) (*CommandOutput, error) {
	return CheckSSMCommandWithClientWithDocumentContextE(t, context.Background(), client, instanceID, command, commandDocName, timeout)
}
