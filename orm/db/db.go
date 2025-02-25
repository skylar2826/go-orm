package db

import (
	"database/sql"
	"geektime-go-orm/orm/db/dialect"
	"geektime-go-orm/orm/db/register"
	"geektime-go-orm/orm/db/valuer"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	R            *register.Register
	DB           *sql.DB
	ValueCreator valuer.ValueCreator
	Dialect      dialect.Dialect
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
	db := &DB{R: &register.Register{
		Models: make(map[string]*register.Model, 1)},
		DB:           sqlDB,
		ValueCreator: valuer.NewUnsafeValue,
		Dialect:      dialect.NewStandardSQL(),
	}

	for _, opt := range opts {
		opt(db)
	}

	return db, nil
}
