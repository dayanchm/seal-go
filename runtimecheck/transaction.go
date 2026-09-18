package runtimecheck

import (
	"context"
	"database/sql"
)

type TrackedTx struct {
	*sql.Tx
	tracker *Tracker
	id      uint64
}

func Begin(tracker *Tracker, db *sql.DB) (*TrackedTx, error) {
	tx, err := db.Begin()

	if err != nil {
		return nil, err
	}

	return trackTx(
		tx,
		tracker,
	), nil
}

func (tx *TrackedTx) Commit() error {
	err := tx.Tx.Commit()

	tx.tracker.Release(tx.id)

	return err
}

func (tx *TrackedTx) Rollback() error {
	err := tx.Tx.Rollback()
	tx.tracker.Release(tx.id)

	return err
}

func BeginTx(tracker *Tracker, db *sql.DB, ctx context.Context, opts *sql.TxOptions) (*TrackedTx, error) {
	tx, err := db.BeginTx(ctx, opts)

	if err != nil {
		return nil, err
	}

	return trackTx(tx, tracker), nil
}

func trackTx(tx *sql.Tx, tracker *Tracker) *TrackedTx {
	id := tracker.Register("sql.Tx")

	return &TrackedTx{
		Tx:      tx,
		tracker: tracker,
		id:      id,
	}
}
