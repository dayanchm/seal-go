package runtimecheck

import (
	"context"
	"time"
)

func WithCancel(tracker *Tracker, parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)
	id := tracker.Register("context.WithCancel")

	trackedCancel := func() {
		cancel()
		tracker.Release(id)
	}

	return ctx, trackedCancel
}

func WithTimeout(tracker *Tracker, parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(parent, timeout)
	id := tracker.Register("context.WithTimeout")

	trackerCancel := func() {
		cancel()
		tracker.Release(id)
	}

	return ctx, trackerCancel
}

func WithDeadline(tracker *Tracker, parent context.Context, deadline time.Time) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithDeadline(parent, deadline)
	id := tracker.Register("context.WithDeadline")

	trackedCancel := func() {
		cancel()
		tracker.Release(id)
	}

	return ctx, trackedCancel
}

func (t *Tracker) ContextOpenResources() []Resource {
	t.mu.Lock()

	defer t.mu.Unlock()

	resources := make([]Resource, 0, len(t.resources))

	for _, resource := range t.resources {
		resources = append(resources, resource)
	}

	return resources
}
