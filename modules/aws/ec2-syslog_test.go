package aws_test

import (
	"context"
	"encoding/base64"
	"errors"
	"testing"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/stretchr/testify/require"

	aws "github.com/gruntwork-io/terratest/modules/aws"
)

// mockEc2SyslogClient is a test double for aws.Ec2SyslogAPI.
type mockEc2SyslogClient struct {
	GetConsoleOutputOutput *ec2.GetConsoleOutputOutput
	GetConsoleOutputErr    error
	lastInstanceID         string
}

func (m *mockEc2SyslogClient) GetConsoleOutput(_ context.Context, params *ec2.GetConsoleOutputInput, _ ...func(*ec2.Options)) (*ec2.GetConsoleOutputOutput, error) {
	m.lastInstanceID = awsSDK.ToString(params.InstanceId)

	if m.GetConsoleOutputErr != nil {
		return nil, m.GetConsoleOutputErr
	}

	return m.GetConsoleOutputOutput, nil
}

func TestGetSyslogForInstanceWithClientContextE(t *testing.T) {
	t.Parallel()

	const (
		instanceID = "i-0abc1234def5678"
		plaintext  = "Welcome to Ubuntu\nKernel booted"
	)

	encoded := base64.StdEncoding.EncodeToString([]byte(plaintext))

	t.Run("returns decoded syslog and forwards instance id", func(t *testing.T) {
		t.Parallel()

		client := &mockEc2SyslogClient{
			GetConsoleOutputOutput: &ec2.GetConsoleOutputOutput{Output: awsSDK.String(encoded)},
		}

		got, err := aws.GetSyslogForInstanceWithClientContextE(t, context.Background(), client, instanceID)
		require.NoError(t, err)
		require.Equal(t, plaintext, got)
		require.Equal(t, instanceID, client.lastInstanceID)
	})

	t.Run("returns empty string without error when syslog not yet available", func(t *testing.T) {
		t.Parallel()

		client := &mockEc2SyslogClient{
			GetConsoleOutputOutput: &ec2.GetConsoleOutputOutput{Output: awsSDK.String("")},
		}

		got, err := aws.GetSyslogForInstanceWithClientContextE(t, context.Background(), client, instanceID)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("propagates api error", func(t *testing.T) {
		t.Parallel()

		client := &mockEc2SyslogClient{GetConsoleOutputErr: errors.New("InvalidInstanceID.NotFound")}

		_, err := aws.GetSyslogForInstanceWithClientContextE(t, context.Background(), client, instanceID)
		require.Error(t, err)
	})

	t.Run("returns decode error on invalid base64 payload", func(t *testing.T) {
		t.Parallel()

		client := &mockEc2SyslogClient{
			GetConsoleOutputOutput: &ec2.GetConsoleOutputOutput{Output: awsSDK.String("!!!not-base64!!!")},
		}

		_, err := aws.GetSyslogForInstanceWithClientContextE(t, context.Background(), client, instanceID)
		require.Error(t, err)
	})
}
