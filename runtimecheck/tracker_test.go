package runtimecheck_test

import (
	"sync"
	"testing"
	"time"

	"seal-go/runtimecheck"
)

func TestTrackerLifecycle(t *testing.T) {
	var tracker runtimecheck.Tracker
	tracker.Release(0)
	if got := tracker.OpenResources(); len(got) != 0 {
		t.Fatalf("new tracker has resources: %v", got)
	}

	before := time.Now()
	first := tracker.Register("context.WithCancel")
	second := tracker.Register("os.File")
	after := time.Now()
	resources := tracker.OpenResources()
	if len(resources) != 2 || resources[0].ID != first || resources[1].ID != second {
		t.Fatalf("unexpected registration order: %v", resources)
	}
	if resources[0].Kind != "context.WithCancel" || resources[1].Kind != "os.File" {
		t.Fatalf("unexpected resource kinds: %v", resources)
	}
	for _, resource := range resources {
		if resource.CreatedAt.Before(before) || resource.CreatedAt.After(after) {
			t.Errorf("timestamp outside registration interval: %v", resource)
		}
	}

	resources[0].Kind = "changed snapshot"
	if tracker.OpenResources()[0].Kind != "context.WithCancel" {
		t.Fatal("editing snapshot changed tracker state")
	}
	tracker.Release(first)
	tracker.Release(first)
	tracker.Release(second + 100)
	if got := tracker.OpenResources(); len(got) != 1 || got[0].ID != second {
		t.Fatalf("release affected the wrong resources: %v", got)
	}
	tracker.Release(second)
	third := tracker.Register("context.WithDeadline")
	if third <= second {
		t.Fatalf("registration ID was reused: %d", third)
	}
	tracker.Release(third)
	if got := tracker.OpenResources(); len(got) != 0 {
		t.Fatalf("released resources remain: %v", got)
	}
}

func TestTrackerConcurrentUse(t *testing.T) {
	tracker := runtimecheck.NewTracker()
	const workers = 100
	ids := make(chan uint64, workers)
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			id := tracker.Register("context.WithCancel")
			ids <- id
			tracker.OpenResources()
			tracker.Release(id)
		}()
	}
	wg.Wait()
	close(ids)

	seen := make(map[uint64]bool)
	for id := range ids {
		if id == 0 || seen[id] {
			t.Errorf("invalid or duplicate resource ID: %d", id)
		}
		seen[id] = true
	}
	if got := tracker.OpenResources(); len(got) != 0 {
		t.Fatalf("resources remain after concurrent cleanup: %v", got)
	}
}
