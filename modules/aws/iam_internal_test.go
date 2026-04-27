package aws

import (
	"context"
	"errors"
	"testing"
	"time"
)

// TestSleepContextHonorsCancellation verifies that sleepContext returns promptly
// with context.Canceled when the context is cancelled mid-sleep, rather than
// waiting the full duration.
func TestSleepContextHonorsCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after a short delay so we can verify the sleep is interrupted.
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := sleepContext(ctx, 30*time.Second)
	elapsed := time.Since(start)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}

	if elapsed > 5*time.Second {
		t.Fatalf("sleepContext did not return promptly after cancellation; took %v", elapsed)
	}
}

// TestSleepContextCompletesNormally verifies that sleepContext returns nil
// when the duration elapses before the context is cancelled.
func TestSleepContextCompletesNormally(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	if err := sleepContext(ctx, 10*time.Millisecond); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
