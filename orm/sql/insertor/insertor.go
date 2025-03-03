package insertor

import (
	"context"
	"geektime-go-orm/orm"
	"geektime-go-orm/orm/aop"
	"geektime-go-orm/orm/db/dialect"
	"geektime-go-orm/orm/db/session"
	"geektime-go-orm/orm/orm_gen/data"
	"geektime-go-orm/orm/predicate"
	sql2 "geektime-go-orm/orm/sql"
)

type Executer interface {
	Execute(ctx context.Context) *orm.QueryResult
}

var _ Executer = &Inserter[data.User]{}

type Inserter[T any] struct {
	table          string
	columns        []predicate.Column
	values         []*T
	OnDuplicateKey *dialect.OnDuplicateKey
	sess           session.Session
	mdls           []aop.Middleware
	sql2.SQLBuilder
}

func (i *Inserter[T]) GetOnDuplicateKeyBuilder() *OnDuplicateKeyBuilder[T] {
	return NewOnDuplicateBuilder(i)
}

func (i *Inserter[T]) Execute(ctx context.Context) *orm.QueryResult {
	query, err := i.Build()
	if err != nil {
		return &orm.QueryResult{Err: err}
	}

	var root aop.Handler = func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult {
		return i.sess.ExecContext(ctx, query.SQL, query.Args...)
	}

	for _, mdl := range i.mdls {
		root = mdl(root)
	}

	return root(ctx, &orm.QueryContext{
		Model:   i.Model,
		Type:    "Insert",
		Builder: i,
	})
}

func (i *Inserter[T]) Build() (*orm.Query, error) {
	i.Sb.WriteString("Insert into ")
	if i.Model == nil {
		t := new(T)
		var err error
		i.Model, err = i.R.Get(t)
		if err != nil {
			return nil, err
		}
	}

	var table string
	if i.table == "" {
		table = i.Model.TableName
	} else {
		table = i.table
	}

	i.Sb.WriteString(table)

	if len(i.columns) > 0 {
		i.Sb.WriteString(" (")
		for idx, col := range i.columns {
			err := i.BuildColumn(col)
			if err != nil {
				return nil, err
			}
			if idx < len(i.columns)-1 {
				i.Sb.WriteString(",")
			}
		}
		i.Sb.WriteString(") ")
	}

	err := i.buildValues()
	if err != nil {
		return nil, err
	}

	err = i.buildOnDuplicateKey()
	if err != nil {
		return nil, err
	}

	i.Sb.WriteString(";")

	return &orm.Query{SQL: i.Sb.String(), Args: i.Args}, nil
}

func (i *Inserter[T]) buildOnDuplicateKey() error {
	if i.OnDuplicateKey != nil {
		i.Dialect.RegisterModel(i.Model)
		OnDuplicateKeyStatement, err := i.Dialect.BuildOnDuplicateKey(i.OnDuplicateKey)
		if err != nil {
			return err
		}
		i.Sb.WriteString(OnDuplicateKeyStatement.Sql)
		i.Args = append(i.Args, OnDuplicateKeyStatement.Args)
	}
	return nil
}

func (i *Inserter[T]) buildValues() error {
	if len(i.values) > 0 {
		i.Sb.WriteString(" Values ")
		for idx, row := range i.values {
			v := i.ValueCreator(row, i.Model)
			i.Sb.WriteString("(")

			if len(i.columns) > 0 {
				for jdx, col := range i.columns {
					val, err := v.Field(col.Name)
					if err != nil {
						return err
					}
					err = i.BuildValuer(predicate.Valuer{Value: val})
					if err != nil {
						return err
					}
					if jdx < len(i.columns)-1 {
						i.Sb.WriteString(",")
					}
				}
			} else {
				jdx := 0
				for colName := range i.Model.Fields {
					val, err := v.Field(colName)
					if err != nil {
						return err
					}
					err = i.BuildValuer(predicate.Valuer{Value: val})
					if err != nil {
						return err
					}
					if jdx < len(i.Model.Fields)-1 {
						i.Sb.WriteString(",")
					}
					jdx++
				}
			}

			i.Sb.WriteString(")")
			if idx < len(i.values)-1 {
				i.Sb.WriteString(",")
			}
		}
	}
	return nil
}

func (i *Inserter[T]) Table(table string) *Inserter[T] {
	i.table = table
	return i
}

func (i *Inserter[T]) Columns(cols ...predicate.Column) *Inserter[T] {
	i.columns = cols
	return i
}

func (i *Inserter[T]) Values(values ...*T) *Inserter[T] {
	i.values = values
	return i
}

func NewInserter[T any](sess session.Session, middlewares ...aop.Middleware) *Inserter[T] {
	i := &Inserter[T]{
		sess:       sess,
		SQLBuilder: sql2.SQLBuilder{Core: sess.GetCore()},
		mdls:       middlewares,
	}

	return i
}
