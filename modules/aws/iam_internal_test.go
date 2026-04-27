package aws

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSleepContextHonorsCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	err := sleepContext(ctx, 30*time.Second)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if time.Since(start) > 5*time.Second {
		t.Fatalf("sleepContext did not return promptly after cancellation")
	}
}
