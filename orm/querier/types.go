package querier

import (
	"context"
	"database/sql"
)

//res, err = db.ExecContext(ctx, "INSERT INTO test_model(`id`, `first_name`, `age`, `last_name`) VALUES (?, ?, ?, ?)", 1, "Tom", 18, "xi")

type Querier[T any] interface {
	// Get 列、聚合函数、子查询
	Get(ctx context.Context) (any, error)
	GetMutli(ctx context.Context) ([]any, error)
}

type Executer interface {
	Execute(ctx *context.Context) (sql.Result, error)
}

type Query struct {
	SQL  string
	Args []any
}

type QueryBuilder interface {
	Build() (*Query, error)
}
