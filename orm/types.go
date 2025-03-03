package sql

type QueryContext struct {
	// Type 声明查询类型。即Select, Update, Delete 和 Insert
	Type    string
	Builder QueryBuilder
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
