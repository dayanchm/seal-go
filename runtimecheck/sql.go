package runtimecheck

import (
	"context"
	"database/sql"
)

type TrackedRows struct {
	*sql.Rows
	tracker *Tracker
	id      uint64
}

func (rows *TrackedRows) Close() error {
	err := rows.Rows.Close()
	rows.tracker.Release(rows.id)
	return err
}
func QueryContext(tracker *Tracker, db *sql.DB, ctx context.Context, query string, args ...any) (*TrackedRows, error) {
	rows, err := db.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	return trackRows(
		tracker,
		rows,
	), nil

}
func Query(tracker *Tracker, db *sql.DB, query string, args ...any) (*TrackedRows, error) {
	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	return trackRows(tracker, rows), nil

}
func (rows *TrackedRows) Next() bool {
	ok := rows.Rows.Next()
	if !ok {
		rows.tracker.Release(rows.id)
	}
	return ok
}
func trackRows(tracker *Tracker, rows *sql.Rows) *TrackedRows {

	id := tracker.Register("sql.Rows")

	return &TrackedRows{
		Rows:    rows,
		tracker: tracker,
		id:      id,
	}
}
