package aws_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	terraaws "github.com/gruntwork-io/terratest/modules/aws"
)

// TestSleepWithContext_CancelReturnsPromptly asserts that sleepWithContext exits within a small
// slack of ctx cancellation rather than waiting out the full duration. This is the core fix that
// lets EnableMfaDeviceContextE honor ctx cancellation during its MFA-propagation waits.
func TestSleepWithContext_CancelReturnsPromptly(t *testing.T) {
	t.Parallel()

	const sleepDuration = 30 * time.Second

	const cancelAfter = 50 * time.Millisecond

	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(cancelAfter, cancel)

	defer cancel()

	start := time.Now()
	err := terraaws.SleepWithContextForTest(ctx, sleepDuration)
	elapsed := time.Since(start)

	require.ErrorIs(t, err, context.Canceled)
	require.Less(t, elapsed, cancelAfter+time.Second)
}

// TestSleepWithContext_DeadlineExceeded propagates ctx.DeadlineExceeded cleanly.
func TestSleepWithContext_DeadlineExceeded(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)

	defer cancel()

	start := time.Now()
	err := terraaws.SleepWithContextForTest(ctx, time.Minute)
	elapsed := time.Since(start)

	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Less(t, elapsed, time.Second)
}

// TestSleepWithContext_AlreadyCancelled: if ctx is already canceled at entry, return immediately.
func TestSleepWithContext_AlreadyCancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	start := time.Now()
	err := terraaws.SleepWithContextForTest(ctx, time.Hour)
	elapsed := time.Since(start)

	require.ErrorIs(t, err, context.Canceled)
	require.Less(t, elapsed, 100*time.Millisecond)
}

// TestSleepWithContext_NormalSleepCompletes: when ctx is not canceled, the function returns
// nil after the duration elapses.
func TestSleepWithContext_NormalSleepCompletes(t *testing.T) {
	t.Parallel()

	const sleepDuration = 30 * time.Millisecond

	ctx := context.Background()

	start := time.Now()
	err := terraaws.SleepWithContextForTest(ctx, sleepDuration)
	elapsed := time.Since(start)

	require.NoError(t, err)
	require.GreaterOrEqual(t, elapsed, sleepDuration)
}
