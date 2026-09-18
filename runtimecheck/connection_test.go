package runtimecheck_test

import (
	"context"
	"net"
	"testing"
	"time"

	"seal-go/runtimecheck"
)

func TestDialTracksAndClosesConnection(t *testing.T) {
	testTrackedConnection(t, func(tracker *runtimecheck.Tracker, address string) (*runtimecheck.TrackedConn, error) {
		return runtimecheck.Dial(tracker, "tcp", address)
	})
}

func TestDialTimeoutTracksAndClosesConnection(t *testing.T) {
	testTrackedConnection(t, func(tracker *runtimecheck.Tracker, address string) (*runtimecheck.TrackedConn, error) {
		return runtimecheck.DialTimeout(tracker, "tcp", address, time.Second)
	})
}

func TestDialContextTracksAndClosesConnection(t *testing.T) {
	testTrackedConnection(t, func(tracker *runtimecheck.Tracker, address string) (*runtimecheck.TrackedConn, error) {
		return runtimecheck.DialContext(tracker, context.Background(), "tcp", address)
	})
}

func TestFailedDialDoesNotRegisterConnection(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	address := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("close listener: %v", err)
	}

	tracker := runtimecheck.NewTracker()
	conn, err := runtimecheck.DialTimeout(tracker, "tcp", address, 100*time.Millisecond)
	if err == nil {
		conn.Close()
		t.Fatal("dial succeeded, want an error")
	}
	if conn != nil {
		t.Fatal("connection should be nil after dial error")
	}
	if got := len(tracker.OpenResources()); got != 0 {
		t.Fatalf("got %d open resources, want 0", got)
	}
}

func testTrackedConnection(
	t *testing.T,
	dial func(*runtimecheck.Tracker, string) (*runtimecheck.TrackedConn, error),
) {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer listener.Close()

	accepted := make(chan net.Conn, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr == nil {
			accepted <- conn
		}
		close(accepted)
	}()

	tracker := runtimecheck.NewTracker()
	conn, err := dial(tracker, listener.Addr().String())
	if err != nil {
		t.Fatalf("dial: %v", err)
	}

	serverConn := <-accepted
	if serverConn == nil {
		conn.Close()
		t.Fatal("server did not accept connection")
	}
	defer serverConn.Close()

	resources := tracker.OpenResources()
	if len(resources) != 1 {
		t.Fatalf("after dial: got %d open resources, want 1", len(resources))
	}
	if resources[0].Kind != "net.Conn" {
		t.Fatalf("got resource kind %q, want %q", resources[0].Kind, "net.Conn")
	}

	if err := conn.Close(); err != nil {
		t.Fatalf("close connection: %v", err)
	}
	if got := len(tracker.OpenResources()); got != 0 {
		t.Fatalf("after close: got %d open resources, want 0", got)
	}
}
