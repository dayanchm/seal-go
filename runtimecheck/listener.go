package runtimecheck

import (
	"net"
)

type TrackedListener struct {
	net.Listener
	tracker *Tracker
	id      uint64
}

func Listen(tracker *Tracker, network string, address string) (*TrackedListener, error) {
	listen, err := net.Listen(network, address)

	if err != nil {
		return nil, err
	}

	return trackListen(tracker, listen), nil

}

func (l *TrackedListener) Close() error {
	err := l.Listener.Close()
	l.tracker.Release(l.id)
	return err
}

func trackListen(tracker *Tracker, listen net.Listener) *TrackedListener {
	id := tracker.Register("net.Listener")

	return &TrackedListener{
		Listener: listen,
		tracker:  tracker,
		id:       id,
	}
}
