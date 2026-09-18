package runtimecheck_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"seal-go/runtimecheck"
)

func TestQueryRegistersAndCloseReleasesRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock database: %v", err)
	}
	defer db.Close()

	query := "SELECT id FROM users"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1),
		)

	tracker := runtimecheck.NewTracker()

	rows, err := runtimecheck.Query(tracker, db, query)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if got := len(tracker.OpenResources()); got != 1 {
		t.Fatalf("before close: got %d open resources, want 1", got)
	}

	if err := rows.Close(); err != nil {
		t.Fatalf("close rows: %v", err)
	}

	if got := len(tracker.OpenResources()); got != 0 {
		t.Fatalf("after close: got %d open resources, want 0", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestQueryContextRegistersRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock database: %v", err)
	}
	defer db.Close()

	query := "SELECT name FROM users"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(
			sqlmock.NewRows([]string{"name"}).AddRow("Ada"),
		)

	tracker := runtimecheck.NewTracker()

	rows, err := runtimecheck.QueryContext(
		tracker,
		db,
		context.Background(),
		query,
	)
	if err != nil {
		t.Fatalf("query context failed: %v", err)
	}
	defer rows.Close()

	if got := len(tracker.OpenResources()); got != 1 {
		t.Fatalf("got %d open resources, want 1", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestNextReleasesRowsAfterIteration(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock database: %v", err)
	}
	defer db.Close()

	query := "SELECT id FROM users"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).
				AddRow(1).
				AddRow(2),
		)

	tracker := runtimecheck.NewTracker()

	rows, err := runtimecheck.Query(tracker, db, query)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	var count int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			t.Fatalf("scan row: %v", err)
		}
		count++
	}

	if err := rows.Err(); err != nil {
		t.Fatalf("iterate rows: %v", err)
	}

	if count != 2 {
		t.Fatalf("got %d rows, want 2", count)
	}

	if got := len(tracker.OpenResources()); got != 0 {
		t.Fatalf("after iteration: got %d open resources, want 0", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestQueryErrorDoesNotRegisterRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock database: %v", err)
	}
	defer db.Close()

	query := "SELECT id FROM missing_table"
	queryErr := errors.New("query failed")

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnError(queryErr)

	tracker := runtimecheck.NewTracker()

	rows, err := runtimecheck.Query(tracker, db, query)
	if !errors.Is(err, queryErr) {
		t.Fatalf("got error %v, want %v", err, queryErr)
	}

	if rows != nil {
		t.Fatal("rows should be nil after query error")
	}

	if got := len(tracker.OpenResources()); got != 0 {
		t.Fatalf("got %d open resources, want 0", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestCloseCanBeCalledTwice(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock database: %v", err)
	}
	defer db.Close()

	query := "SELECT id FROM users"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id"}).AddRow(1),
		)

	tracker := runtimecheck.NewTracker()

	rows, err := runtimecheck.Query(tracker, db, query)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}

	if err := rows.Close(); err != nil {
		t.Fatalf("first close: %v", err)
	}

	if err := rows.Close(); err != nil {
		t.Fatalf("second close: %v", err)
	}

	if got := len(tracker.OpenResources()); got != 0 {
		t.Fatalf("got %d open resources, want 0", got)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}

func TestQueryForwardsArguments(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create mock database: %v", err)
	}
	defer db.Close()

	query := "SELECT name FROM users WHERE id = ? AND role = ?"

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(42, "admin").
		WillReturnRows(
			sqlmock.NewRows([]string{"name"}).AddRow("Ada"),
		)

	tracker := runtimecheck.NewTracker()

	rows, err := runtimecheck.Query(
		tracker,
		db,
		query,
		42,
		"admin",
	)
	if err != nil {
		t.Fatalf("query failed: %v", err)
	}
	defer rows.Close()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("database expectations: %v", err)
	}
}
