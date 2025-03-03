package db

import "database/sql"

type Tx struct {
	tx *sql.Tx
	db *sql.DB
}

func (tx *Tx) Commit() error {
	
}
