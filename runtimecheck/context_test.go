package runtimecheck_test

import (
	"context"
	"seal-go/runtimecheck"
	"testing"
	"time"
)

// Test WithCancel
func TestWithCancel(t *testing.T) {
	tracker := runtimecheck.NewTracker()

	_, cancel := runtimecheck.WithCancel(
		tracker,
		context.Background(),
	)

	if got := len(tracker.OpenResources()); got != 1 {
		t.Fatalf("before cancel: got %d open resources, want 1", got)
	}

	cancel()

	if got := len(tracker.OpenResources()); got != 0 {
		t.Fatalf("after cancel: got %d open resources, want 0", got)
	}
}

// Test WithTimeOut

func TestWithTimeout(t *testing.T) {
	tracker := runtimecheck.NewTracker()

	_, cancel := runtimecheck.WithTimeout(
		tracker,
		context.Background(),
		time.Second,
	)

	if got := len(tracker.OpenResources()); got != 1 {
		t.Fatalf("before timeout: got %d open resources, want 1", got)
	}

	cancel()

	if got := len(tracker.OpenResources()); got != 0 {
		t.Fatalf("after timeout: got %d open resources, want 0", got)
	}
}

// Test WithDeadline

func TestWithDeadline(t *testing.T) {
	tracker := runtimecheck.NewTracker()

	_, cancel := runtimecheck.WithDeadline(
		tracker,
		context.Background(),
		time.Now().Add(time.Second),
	)

	if got := len(tracker.OpenResources()); got != 1 {
		t.Fatalf("before deadline: got %d open resources, want 1", got)
	}

	cancel()

	if got := len(tracker.OpenResources()); got != 0 {
		t.Fatalf("after deadline: got %d open resources, want 0", got)
	}
}

func TestUncancelledContextRemainsOpen(t *testing.T) {
	tracker := runtimecheck.NewTracker()

	runtimecheck.WithCancel(tracker, context.Background())

	if got := len(tracker.OpenResources()); got != 1 {
		t.Fatalf("got %d open resources, want 1", got)
	}
}
