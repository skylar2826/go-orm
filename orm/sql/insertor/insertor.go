package insertor

import (
	"context"
	"database/sql"
	"geektime-go-orm/orm/data"
	"geektime-go-orm/orm/db"
	"geektime-go-orm/orm/db/dialect"
	"geektime-go-orm/orm/db/register"
	"geektime-go-orm/orm/predicate"
	sql2 "geektime-go-orm/orm/sql"
	sqlBuilder2 "geektime-go-orm/orm/sqlCommonBuilder"
)

type Executer interface {
	Execute(ctx context.Context) (sql.Result, error)
}

var _ Executer = &Inserter[data.User]{}

type Inserter[T any] struct {
	sqlBuilder     *sqlBuilder2.SQLBuilder
	model          *register.Model
	db             *db.DB
	table          string
	columns        []predicate.Column
	values         []*T
	OnDuplicateKey *dialect.OnDuplicateKey
}

func (i *Inserter[T]) GetOnDuplicateKeyBuilder() *OnDuplicateKeyBuilder[T] {
	return NewOnDuplicateBuilder(i)
}

func (i *Inserter[T]) Execute(ctx context.Context) (sql.Result, error) {
	query, err := i.Build()
	if err != nil {
		return nil, err
	}

	return i.db.DB.ExecContext(ctx, query.SQL, query.Args...)
}

func (i *Inserter[T]) Build() (*sql2.Query, error) {
	i.sqlBuilder.Sb.WriteString("Insert into ")
	if i.model == nil {
		t := new(T)
		var err error
		i.model, err = i.db.R.Get(t)
		if err != nil {
			return nil, err
		}
		i.sqlBuilder.RegisterModel(i.model)
	}

	var table string
	if i.table == "" {
		table = i.model.TableName
	} else {
		table = i.table
	}

	i.sqlBuilder.Sb.WriteString(table)

	if len(i.columns) > 0 {
		i.sqlBuilder.Sb.WriteString(" (")
		for idx, col := range i.columns {
			err := i.sqlBuilder.BuildColumn(col)
			if err != nil {
				return nil, err
			}
			if idx < len(i.columns)-1 {
				i.sqlBuilder.Sb.WriteString(",")
			}
		}
		i.sqlBuilder.Sb.WriteString(") ")
	}

	err := i.buildValues()
	if err != nil {
		return nil, err
	}

	if i.OnDuplicateKey != nil {
		err = i.db.Dialect.BuildOnDuplicateKey(i.sqlBuilder, i.OnDuplicateKey)
		if err != nil {
			return nil, err
		}
	}

	i.sqlBuilder.Sb.WriteString(";")

	return &sql2.Query{SQL: i.sqlBuilder.Sb.String(), Args: i.sqlBuilder.Args}, nil
}

func (i *Inserter[T]) buildValues() error {
	if len(i.values) > 0 {
		i.sqlBuilder.Sb.WriteString(" Values ")
		for idx, row := range i.values {
			v := i.db.ValueCreator(row, i.model)
			i.sqlBuilder.Sb.WriteString("(")

			if len(i.columns) > 0 {
				for jdx, col := range i.columns {
					val, err := v.Field(col.Name)
					if err != nil {
						return err
					}
					err = i.sqlBuilder.BuildValuer(predicate.Valuer{Value: val})
					if err != nil {
						return err
					}
					if jdx < len(i.columns)-1 {
						i.sqlBuilder.Sb.WriteString(",")
					}
				}
			} else {
				jdx := 0
				for colName := range i.model.Fields {
					val, err := v.Field(colName)
					if err != nil {
						return err
					}
					err = i.sqlBuilder.BuildValuer(predicate.Valuer{Value: val})
					if err != nil {
						return err
					}
					if jdx < len(i.model.Fields)-1 {
						i.sqlBuilder.Sb.WriteString(",")
					}
					jdx++
				}
			}

			i.sqlBuilder.Sb.WriteString(")")
			if idx < len(i.values)-1 {
				i.sqlBuilder.Sb.WriteString(",")
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

func NewInserter[T any](db *db.DB) *Inserter[T] {
	return &Inserter[T]{
		db:         db,
		sqlBuilder: &sqlBuilder2.SQLBuilder{Quote: db.Dialect.Quoter()},
	}
}
