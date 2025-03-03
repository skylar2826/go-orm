package session

import (
	"context"
	"database/sql"
	sql2 "geektime-go-orm/orm"
)

type Tx struct {
	tx *sql.Tx
	db *DB
}

func (tx *Tx) GetCore() Core {
	return tx.db.Core
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...any) *sql2.QueryResult {
	res, err := tx.tx.ExecContext(ctx, query, args...)
	return &sql2.QueryResult{
		Result: res,
		Err:    err,
	}
}

func (tx *Tx) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.tx.QueryContext(ctx, query, args...)
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}
