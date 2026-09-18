package runtimecheck_test

import (
	"testing"

	"seal-go/runtimecheck"
)

func TestListenTracksAndClosesListener(t *testing.T) {
	tracker := runtimecheck.NewTracker()

	listener, err := runtimecheck.Listen(tracker, "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}

	resources := tracker.OpenResources()
	if len(resources) != 1 {
		listener.Close()
		t.Fatalf("after listen: got %d open resources, want 1", len(resources))
	}
	if resources[0].Kind != "net.Listener" {
		listener.Close()
		t.Fatalf("got resource kind %q, want %q", resources[0].Kind, "net.Listener")
	}

	if err := listener.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}
	if got := len(tracker.OpenResources()); got != 0 {
		t.Fatalf("after close: got %d open resources, want 0", got)
	}
}

func TestFailedListenDoesNotRegisterListener(t *testing.T) {
	tracker := runtimecheck.NewTracker()

	listener, err := runtimecheck.Listen(tracker, "invalid-network", "127.0.0.1:0")
	if err == nil {
		listener.Close()
		t.Fatal("listen succeeded, want an error")
	}
	if listener != nil {
		t.Fatal("listener should be nil after listen error")
	}
	if got := len(tracker.OpenResources()); got != 0 {
		t.Fatalf("got %d open resources, want 0", got)
	}
}
