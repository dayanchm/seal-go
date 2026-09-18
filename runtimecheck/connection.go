package runtimecheck

import (
	"context"
	"net"
	"time"
)

type TrackedConn struct {
	net.Conn
	tracker *Tracker
	id      uint64
}

func Dial(tracker *Tracker, network string, address string) (*TrackedConn, error) {
	conn, err := net.Dial(network, address)

	if err != nil {
		return nil, err
	}

	return trackConn(tracker, conn), nil
}

func DialTimeout(
	tracker *Tracker,
	network string,
	address string,
	timeout time.Duration,
) (*TrackedConn, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}
	return trackConn(tracker, conn), nil
}

func DialContext(
	tracker *Tracker,
	ctx context.Context,
	network, address string,
) (*TrackedConn, error) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}
	return trackConn(tracker, conn), nil
}

func (conn *TrackedConn) Close() error {

	err := conn.Conn.Close()
	if err != nil {
		return err
	}
	conn.tracker.Unregister(conn.id)
	return nil

}

func trackConn(tracker *Tracker, conn net.Conn) *TrackedConn {
	id := tracker.Register("net.Conn")

	return &TrackedConn{
		Conn:    conn,
		tracker: tracker,
		id:      id,
	}
}
