package runtimecheck_test

import (
	"seal-go/runtimecheck"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBeginRegistersTransaction(t *testing.T) {

	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("create mock database: %v", err)
	}

	defer db.Close()

	mock.ExpectBegin()

	tracker := runtimecheck.NewTracker()

	tx, err := runtimecheck.Begin(tracker, db)

	if err != nil {
		t.Fatalf("begin transaction: %v", err)
	}

	defer tx.Rollback()

	if got := len(tracker.OpenResources()); got != 1 {
		t.Fatalf("got %d open resources, want 1", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations : %v", err)
	}
}
