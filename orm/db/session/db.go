package session

import (
	"context"
	"database/sql"
	"fmt"
	sql2 "geektime-go-orm/orm"
	"geektime-go-orm/orm/db/dialect"
	"geektime-go-orm/orm/db/register"
	"geektime-go-orm/orm/db/valuer"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	DB *sql.DB
	Core
}

func (db *DB) GetCore() Core {
	return db.Core
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...any) *sql2.QueryResult {
	res, err := db.DB.ExecContext(ctx, query, args...)
	return &sql2.QueryResult{
		Result: res,
		Err:    err,
	}
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.DB.QueryContext(ctx, query, args...)
}

func (db *DB) BeginTx(ctx context.Context, txOpts *sql.TxOptions) (*Tx, error) {
	tx, err := db.DB.BeginTx(ctx, txOpts)
	if err != nil {
		return nil, err
	}
	return &Tx{
		tx: tx,
		db: db,
	}, nil
}

type TxKey struct {
}

func (db *DB) BeginTxV2(ctx context.Context, txOpts *sql.TxOptions) (context.Context, *Tx, error) {
	val := ctx.Value(TxKey{})
	if val != nil {
		return ctx, val.(*Tx), nil
	}
	tx, err := db.BeginTx(ctx, txOpts)
	if err != nil {
		return ctx, nil, err
	}
	ctx = context.WithValue(ctx, TxKey{}, tx)
	return ctx, tx, nil
}

func (db *DB) DoTx(ctx context.Context, fn func(ctx context.Context, tx *Tx) error, txOpts *sql.TxOptions) error {
	tx, err := db.BeginTx(ctx, txOpts)
	if err != nil {
		return err
	}

	panicked := true
	err = fn(ctx, tx)
	panicked = false

	defer func() {
		if err != nil || panicked {
			e := tx.Rollback()

			err = fmt.Errorf("业务错误： %w, 回滚错误： %s, 是否panic: %t\n", err, e, panicked)
		} else {
			err = tx.Commit()
		}
	}()
	return nil
}

type DBOption func(db *DB)

func WithReflectValue() DBOption {
	return func(db *DB) {
		db.ValueCreator = valuer.NewReflectValue
	}
}

func WithUnsafeValue() DBOption {
	return func(db *DB) {
		db.ValueCreator = valuer.NewUnsafeValue
	}
}

func WithDialect(dialect2 dialect.Dialect) DBOption {
	return func(db *DB) {
		db.Dialect = dialect2
	}
}

func Open(driver string, datasourceName string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, datasourceName)
	if err != nil {
		return nil, err
	}

	return OpenDB(db, opts...)
}

func OpenDB(sqlDB *sql.DB, opts ...DBOption) (*DB, error) {
	db := &DB{
		Core: Core{
			R: &register.Register{
				Models: make(map[string]*register.Model, 1),
			},
			ValueCreator: valuer.NewUnsafeValue,
			Dialect:      dialect.NewStandardSQL(),
		},
		DB: sqlDB,
	}

	for _, opt := range opts {
		opt(db)
	}

	return db, nil
}
