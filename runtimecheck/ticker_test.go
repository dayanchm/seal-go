package runtimecheck_test

import (
	"testing"
	"time"

	"seal-go/runtimecheck"
)

func TestNewTickerTracksAndStopsTicker(t *testing.T) {
	tracker := runtimecheck.NewTracker()
	ticker := runtimecheck.NewTicker(tracker, time.Hour)

	resources := tracker.OpenResources()
	if len(resources) != 1 {
		ticker.Stop()
		t.Fatalf("after NewTicker: got %d open resources, want 1", len(resources))
	}
	if resources[0].Kind != "time.Ticker" {
		ticker.Stop()
		t.Fatalf("got resource kind %q, want %q", resources[0].Kind, "time.Ticker")
	}

	ticker.Stop()
	if got := len(tracker.OpenResources()); got != 0 {
		t.Fatalf("after Stop: got %d open resources, want 0", got)
	}
}

func TestTrackedTickerEmitsTick(t *testing.T) {
	tracker := runtimecheck.NewTracker()
	ticker := runtimecheck.NewTicker(tracker, 10*time.Millisecond)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		// The embedded time.Ticker still behaves like a normal ticker.
	case <-time.After(time.Second):
		t.Fatal("ticker did not emit a tick")
	}
}

func TestNewTickerPanicsForNonPositiveDuration(t *testing.T) {
	tracker := runtimecheck.NewTracker()

	defer func() {
		if recover() == nil {
			t.Fatal("NewTicker did not panic for zero duration")
		}
		if got := len(tracker.OpenResources()); got != 0 {
			t.Fatalf("got %d open resources after panic, want 0", got)
		}
	}()

	runtimecheck.NewTicker(tracker, 0)
}
