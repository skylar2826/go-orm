package db

import (
	"database/sql"
	"geektime-go-orm/orm/db/register"
	"geektime-go-orm/orm/db/valuer"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	R            *register.Register
	DB           *sql.DB
	ValueCreator valuer.ValueCreator
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

func Open(driver string, datasourceName string, opts ...DBOption) (*DB, error) {
	db, err := sql.Open(driver, datasourceName)
	if err != nil {
		return nil, err
	}

	return OpenDB(db, opts...)
}

func OpenDB(sqlDB *sql.DB, opts ...DBOption) (*DB, error) {
	db := &DB{R: &register.Register{Models: make(map[string]*register.Model, 1)}, DB: sqlDB, ValueCreator: valuer.NewUnsafeValue}

	for _, opt := range opts {
		opt(db)
	}

	return db, nil
}
