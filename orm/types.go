package orm

import "geektime-go-orm/orm/db/register"

type QueryContext struct {
	// Type 声明查询类型。即Select, Update, Delete 和 Insert
	Type    string
	Builder QueryBuilder
	Model   *register.Model
}

type QueryResult struct {
	Result any
	Err    error
}

type Query struct {
	SQL  string
	Args []any
}

type QueryBuilder interface {
	Build() (*Query, error)
}
