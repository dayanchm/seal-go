package runtimecheck

import (
	"sort"
	"sync"
	"time"
)

// Resource describes a registered resource whose cleanup has not been recorded.
// An open record is not necessarily a leak: the resource may still be in use.
type Resource struct {
	ID        uint64
	Kind      string
	CreatedAt time.Time
}

// Tracker records outstanding resources. Its zero value is ready to use.
// A Tracker is safe for concurrent use and must not be copied after first use.
type Tracker struct {
	mu        sync.Mutex
	nextID    uint64
	resources map[uint64]Resource
}

func (t *Tracker) Unregister(id uint64) {
	panic("unimplemented")
}

func (t *Tracker) Track(s string) any {
	panic("unimplemented")
}

func NewTracker() *Tracker {
	return &Tracker{}
}

// Register records a resource and returns an ID unique within this tracker.
func (t *Tracker) Register(kind string) uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.resources == nil {
		t.resources = make(map[uint64]Resource)
	}
	t.nextID++
	id := t.nextID
	t.resources[id] = Resource{
		ID:        id,
		Kind:      kind,
		CreatedAt: time.Now(),
	}
	return id
}

// Release removes a record; it does not close or cancel the actual resource.
// Releasing an unknown or already released ID has no effect.
func (t *Tracker) Release(id uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.resources, id)
}

// OpenResources returns an independent snapshot ordered by registration ID.
func (t *Tracker) OpenResources() []Resource {
	t.mu.Lock()
	resources := make([]Resource, 0, len(t.resources))
	for _, resource := range t.resources {
		resources = append(resources, resource)
	}
	t.mu.Unlock()

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].ID < resources[j].ID
	})
	return resources
}
