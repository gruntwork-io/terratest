package aws_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	awsSDK "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/stretchr/testify/require"

	terraaws "github.com/gruntwork-io/terratest/modules/aws"
)

// fakeSsmClient is an inline test double for aws.SsmAPI. It records how many times each method is
// invoked so tests can assert the retry loop exited promptly after ctx cancellation.
type fakeSsmClient struct {
	// When non-nil, SendCommand returns this response (CommandId is used by the subsequent
	// GetCommandInvocation calls). If nil, SendCommand returns a synthetic response.
	sendCommandOutput *ssm.SendCommandOutput

	// When non-nil, SendCommand errors. Otherwise sendCommandOutput is used.
	sendCommandErr error

	// getInventoryErr / getCommandInvocationErr are the errors (retryable) returned by each call.
	// When nil, each call returns a retryable condition so the loop keeps spinning until ctx is canceled.
	getInventoryErr         error
	getCommandInvocationErr error

	getInventoryCalls         int32
	sendCommandCalls          int32
	getCommandInvocationCalls int32
}

func (f *fakeSsmClient) GetInventory(_ context.Context, _ *ssm.GetInventoryInput, _ ...func(*ssm.Options)) (*ssm.GetInventoryOutput, error) {
	atomic.AddInt32(&f.getInventoryCalls, 1)

	if f.getInventoryErr != nil {
		return nil, f.getInventoryErr
	}

	// Return an empty inventory so the retry loop treats this as a retryable condition and keeps spinning.
	return &ssm.GetInventoryOutput{Entities: nil}, nil
}

func (f *fakeSsmClient) SendCommand(_ context.Context, _ *ssm.SendCommandInput, _ ...func(*ssm.Options)) (*ssm.SendCommandOutput, error) {
	atomic.AddInt32(&f.sendCommandCalls, 1)

	if f.sendCommandErr != nil {
		return nil, f.sendCommandErr
	}

	if f.sendCommandOutput != nil {
		return f.sendCommandOutput, nil
	}

	return &ssm.SendCommandOutput{Command: &types.Command{CommandId: awsSDK.String("cmd-id")}}, nil
}

func (f *fakeSsmClient) GetCommandInvocation(_ context.Context, _ *ssm.GetCommandInvocationInput, _ ...func(*ssm.Options)) (*ssm.GetCommandInvocationOutput, error) {
	atomic.AddInt32(&f.getCommandInvocationCalls, 1)

	if f.getCommandInvocationErr != nil {
		return nil, f.getCommandInvocationErr
	}

	// Return "Pending" — one of the retryable statuses — so the retry loop keeps going until ctx is canceled.
	return &ssm.GetCommandInvocationOutput{
		Status:        types.CommandInvocationStatusPending,
		StatusDetails: awsSDK.String("bad status: Pending"),
	}, nil
}

// TestWaitForSsmInstanceWithClientContextE_HonorsCtxCancellation asserts that when the caller's ctx is
// canceled, the retry loop exits within roughly one retry interval, not the full timeout.
func TestWaitForSsmInstanceWithClientContextE_HonorsCtxCancellation(t *testing.T) {
	t.Parallel()

	client := &fakeSsmClient{}

	// Configure a 10-minute timeout but cancel ctx after 100ms. Without the fix the call would
	// block for up to the full timeout. With the fix it returns within < one retry interval.
	const timeout = 10 * time.Minute

	const cancelAfter = 100 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(cancelAfter, cancel)

	defer cancel()

	start := time.Now()
	err := terraaws.WaitForSsmInstanceWithClientContextE(t, ctx, client, "i-123", timeout)
	elapsed := time.Since(start)

	require.ErrorIs(t, err, context.Canceled)

	// Allow generous slack — we expect well under one retry interval (2s) + cancelAfter.
	maxExpected := cancelAfter + 3*time.Second
	require.Less(t, elapsed, maxExpected, "retry loop did not exit promptly after ctx cancel: elapsed=%s", elapsed)
}

// TestWaitForSsmInstanceWithClientContextE_AlreadyCancelled asserts behavior when ctx is already
// canceled at entry — the function must return immediately with ctx.Err().
func TestWaitForSsmInstanceWithClientContextE_AlreadyCancelled(t *testing.T) {
	t.Parallel()

	client := &fakeSsmClient{}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	err := terraaws.WaitForSsmInstanceWithClientContextE(t, ctx, client, "i-123", time.Minute)
	elapsed := time.Since(start)

	require.ErrorIs(t, err, context.Canceled)
	require.Less(t, elapsed, time.Second)
}

// TestWaitForSsmInstanceWithClientContextE_DeadlineExceeded asserts that a ctx.DeadlineExceeded
// is propagated unwrapped (i.e. errors.Is sees it).
func TestWaitForSsmInstanceWithClientContextE_DeadlineExceeded(t *testing.T) {
	t.Parallel()

	client := &fakeSsmClient{}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)

	defer cancel()

	err := terraaws.WaitForSsmInstanceWithClientContextE(t, ctx, client, "i-123", 10*time.Minute)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

// TestCheckSSMCommandWithClientWithDocumentContextE_HonorsCtxCancellation asserts that when the caller's
// ctx is canceled mid-retry the loop exits promptly rather than waiting out the full timeout.
func TestCheckSSMCommandWithClientWithDocumentContextE_HonorsCtxCancellation(t *testing.T) {
	t.Parallel()

	client := &fakeSsmClient{}

	const timeout = 10 * time.Minute

	const cancelAfter = 100 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(cancelAfter, cancel)

	defer cancel()

	start := time.Now()
	_, err := terraaws.CheckSSMCommandWithClientWithDocumentContextE(t, ctx, client, "i-123", "echo hi", "AWS-RunShellScript", timeout)
	elapsed := time.Since(start)

	require.ErrorIs(t, err, context.Canceled)

	maxExpected := cancelAfter + 3*time.Second
	require.Less(t, elapsed, maxExpected, "retry loop did not exit promptly after ctx cancel: elapsed=%s", elapsed)
}
