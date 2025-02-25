package valuer

import (
	"database/sql"
	"geektime-go-orm/orm/db/register"
)

// Value 是对结构体实例的内部抽象
type Value interface {
	SetColumns(rows *sql.Rows) error
	Field(name string) (any, error)
}

type ValueCreator func(val any, model *register.Model) Value
