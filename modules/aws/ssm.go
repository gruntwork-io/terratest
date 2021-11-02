package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/testing"
	"github.com/stretchr/testify/require"
)

// GetParameter retrieves the latest version of SSM Parameter at keyName with decryption.
func GetParameter(t testing.TestingT, awsRegion string, keyName string) string {
	keyValue, err := GetParameterE(t, awsRegion, keyName)
	require.NoError(t, err)
	return keyValue
}

// GetParameterE retrieves the latest version of SSM Parameter at keyName with decryption.
func GetParameterE(t testing.TestingT, awsRegion string, keyName string) (string, error) {
	ssmClient, err := NewSsmClientE(t, awsRegion)
	if err != nil {
		return "", err
	}

	return GetParameterWithClientE(t, ssmClient, keyName)
}

// GetParameterE retrieves the latest version of SSM Parameter at keyName with decryption with the ability to provide the SSM client.
func GetParameterWithClientE(t testing.TestingT, client *ssm.SSM, keyName string) (string, error) {
	resp, err := client.GetParameter(&ssm.GetParameterInput{Name: aws.String(keyName), WithDecryption: aws.Bool(true)})
	if err != nil {
		return "", err
	}

	parameter := *resp.Parameter
	return *parameter.Value, nil
}

// PutParameter creates new version of SSM Parameter at keyName with keyValue as SecureString.
func PutParameter(t testing.TestingT, awsRegion string, keyName string, keyDescription string, keyValue string) int64 {
	version, err := PutParameterE(t, awsRegion, keyName, keyDescription, keyValue)
	require.NoError(t, err)
	return version
}

// PutParameterE creates new version of SSM Parameter at keyName with keyValue as SecureString.
func PutParameterE(t testing.TestingT, awsRegion string, keyName string, keyDescription string, keyValue string) (int64, error) {
	ssmClient, err := NewSsmClientE(t, awsRegion)
	if err != nil {
		return 0, err
	}
	return PutParameterWithClientE(t, ssmClient, keyName, keyDescription, keyValue)
}

// PutParameterE creates new version of SSM Parameter at keyName with keyValue as SecureString with the ability to provide the SSM client.
func PutParameterWithClientE(t testing.TestingT, client *ssm.SSM, keyName string, keyDescription string, keyValue string) (int64, error) {
	resp, err := client.PutParameter(&ssm.PutParameterInput{Name: aws.String(keyName), Description: aws.String(keyDescription), Value: aws.String(keyValue), Type: aws.String("SecureString")})
	if err != nil {
		return 0, err
	}

	return *resp.Version, nil
}

// DeleteParameter deletes all versions of SSM Parameter at keyName.
func DeleteParameter(t testing.TestingT, awsRegion string, keyName string) {
	err := DeleteParameterE(t, awsRegion, keyName)
	require.NoError(t, err)
}

// DeleteParameterE deletes all versions of SSM Parameter at keyName.
func DeleteParameterE(t testing.TestingT, awsRegion string, keyName string) error {
	ssmClient, err := NewSsmClientE(t, awsRegion)
	if err != nil {
		return err
	}
	return DeleteParameterWithClientE(t, ssmClient, keyName)
}

// DeleteParameterE deletes all versions of SSM Parameter at keyName with the ability to provide the SSM client.
func DeleteParameterWithClientE(t testing.TestingT, client *ssm.SSM, keyName string) error {
	_, err := client.DeleteParameter(&ssm.DeleteParameterInput{Name: aws.String(keyName)})
	if err != nil {
		return err
	}

	return nil
}

// NewSsmClient creates a SSM client.
func NewSsmClient(t testing.TestingT, region string) *ssm.SSM {
	client, err := NewSsmClientE(t, region)
	require.NoError(t, err)
	return client
}

// NewSsmClientE creates an SSM client.
func NewSsmClientE(t testing.TestingT, region string) (*ssm.SSM, error) {
	sess, err := NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}

	return ssm.New(sess), nil
}

// WaitForSsmInstanceE waits until the instance get registered to the SSM inventory.
func WaitForSsmInstanceE(t testing.TestingT, awsRegion, instanceID string, timeout time.Duration) error {
	client, err := NewSsmClientE(t, awsRegion)
	if err != nil {
		return err
	}
	return WaitForSsmInstanceWithClientE(t, client, instanceID, timeout)
}

// WaitForSsmInstanceE waits until the instance get registered to the SSM inventory with the ability to provide the SSM client.
func WaitForSsmInstanceWithClientE(t testing.TestingT, client *ssm.SSM, instanceID string, timeout time.Duration) error {
	timeBetweenRetries := 2 * time.Second
	maxRetries := int(timeout.Seconds() / timeBetweenRetries.Seconds())
	description := fmt.Sprintf("Waiting for %s to appear in the SSM inventory", instanceID)

	input := &ssm.GetInventoryInput{
		Filters: []*ssm.InventoryFilter{
			{
				Key:    aws.String("AWS:InstanceInformation.InstanceId"),
				Type:   aws.String("Equal"),
				Values: aws.StringSlice([]string{instanceID}),
			},
		},
	}
	_, err := retry.DoWithRetryE(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		resp, err := client.GetInventory(input)

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

// WaitForSsmInstance waits until the instance get registered to the SSM inventory.
func WaitForSsmInstance(t testing.TestingT, awsRegion, instanceID string, timeout time.Duration) {
	err := WaitForSsmInstanceE(t, awsRegion, instanceID, timeout)
	require.NoError(t, err)
}

// CheckSsmCommand checks that you can run the given command on the given instance through AWS SSM.
func CheckSsmCommand(t testing.TestingT, awsRegion, instanceID, command string, timeout time.Duration) *CommandOutput {
	return CheckSsmCommandWithDocument(t, awsRegion, instanceID, command, "AWS-RunShellScript", timeout)
}

// CommandOutput contains the result of the SSM command.
type CommandOutput struct {
	Stdout   string
	Stderr   string
	ExitCode int64
}

// CheckSsmCommandE checks that you can run the given command on the given instance through AWS SSM. Returns the result and an error if one occurs.
func CheckSsmCommandE(t testing.TestingT, awsRegion, instanceID, command string, timeout time.Duration) (*CommandOutput, error) {
	return CheckSsmCommandWithDocumentE(t, awsRegion, instanceID, command, "AWS-RunShellScript", timeout)
}

// CheckSSMCommandWithClientE checks that you can run the given command on the given instance through AWS SSM with the ability to provide the SSM client. Returns the result and an error if one occurs.
func CheckSSMCommandWithClientE(t testing.TestingT, client *ssm.SSM, instanceID, command string, timeout time.Duration) (*CommandOutput, error) {
	return CheckSSMCommandWithClientWithDocumentE(t, client, instanceID, command, "AWS-RunShellScript", timeout)
}

// CheckSsmCommandWithDocument checks that you can run the given command on the given instance through AWS SSM with specified Command Doc type.
func CheckSsmCommandWithDocument(t testing.TestingT, awsRegion, instanceID, command string, commandDocName string, timeout time.Duration) *CommandOutput {
	result, err := CheckSsmCommandWithDocumentE(t, awsRegion, instanceID, command, commandDocName, timeout)
	require.NoErrorf(t, err, "failed to execute '%s' on %s (%v):]\n  stdout: %#v\n  stderr: %#v", command, instanceID, err, result.Stdout, result.Stderr)
	return result
}

// CheckSsmCommandWithDocumentE checks that you can run the given command on the given instance through AWS SSM with specified Command Doc type. Returns the result and an error if one occurs.
func CheckSsmCommandWithDocumentE(t testing.TestingT, awsRegion, instanceID, command string, commandDocName string, timeout time.Duration) (*CommandOutput, error) {
	logger.Logf(t, "Running command '%s' on EC2 instance with ID '%s'", command, instanceID)

	// Now that we know the instance in the SSM inventory, we can send the command
	client, err := NewSsmClientE(t, awsRegion)
	if err != nil {
		return nil, err
	}
	return CheckSSMCommandWithClientWithDocumentE(t, client, instanceID, command, commandDocName, timeout)
}

// CheckSSMCommandWithClientWithDocumentE checks that you can run the given command on the given instance through AWS SSM with the ability to provide the SSM client with specified Command Doc type. Returns the result and an error if one occurs.
func CheckSSMCommandWithClientWithDocumentE(t testing.TestingT, client *ssm.SSM, instanceID, command string, commandDocName string, timeout time.Duration) (*CommandOutput, error) {

	timeBetweenRetries := 2 * time.Second
	maxRetries := int(timeout.Seconds() / timeBetweenRetries.Seconds())

	resp, err := client.SendCommand(&ssm.SendCommandInput{
		Comment:      aws.String("Terratest SSM"),
		DocumentName: aws.String(commandDocName),
		InstanceIds:  aws.StringSlice([]string{instanceID}),
		Parameters: map[string][]*string{
			"commands": aws.StringSlice([]string{command}),
		},
	})
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
	_, err = retry.DoWithRetryableErrorsE(t, description, retryableErrors, maxRetries, timeBetweenRetries, func() (string, error) {
		resp, err := client.GetCommandInvocation(&ssm.GetCommandInvocationInput{
			CommandId:  resp.Command.CommandId,
			InstanceId: &instanceID,
		})

		if err != nil {
			return "", err
		}

		result.Stderr = aws.StringValue(resp.StandardErrorContent)
		result.Stdout = aws.StringValue(resp.StandardOutputContent)
		result.ExitCode = aws.Int64Value(resp.ResponseCode)

		status := aws.StringValue(resp.Status)

		if status == ssm.CommandInvocationStatusSuccess {
			return "", nil
		}

		if status == ssm.CommandInvocationStatusFailed {
			return "", fmt.Errorf(aws.StringValue(resp.StatusDetails))
		}

		return "", fmt.Errorf("bad status: %s", status)
	})

	if err != nil {
		if actualErr, ok := err.(retry.FatalError); ok {
			return result, actualErr.Underlying
		}
		return result, fmt.Errorf("unexpected error: %v", err)
	}

	return result, nil
}
