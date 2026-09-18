package runtimecheck

import (
	"time"
)

type TrackedTicker struct {
	*time.Ticker
	tracker *Tracker
	id      uint64
}

func NewTicker(tracker *Tracker, duration time.Duration) *TrackedTicker {
	ticker := time.NewTicker(duration)

	id := tracker.Register("ticket")

	return &TrackedTicker{
		Ticker:  ticker,
		tracker: tracker,
		id:      id}
}

func (ticker *TrackedTicker) Stop() {
	ticker.Ticker.Stop()
	ticker.tracker.Release(ticker.id)
}

func trackTicker(tracker *Tracker, ticker *time.Ticker) *TrackedTicker {
	id := tracker.Register("ticker")

	return &TrackedTicker{
		Ticker:  ticker,
		tracker: tracker,
		id:      id,
	}
}
