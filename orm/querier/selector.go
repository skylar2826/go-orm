package querier

import (
	"context"
	"database/sql"
	"fmt"
	"geektime-go-orm/orm/db"
	"geektime-go-orm/orm/db/register"
	"geektime-go-orm/orm/db/valuer"
	"geektime-go-orm/orm/predicate"
	"reflect"
)

type Selector[T any] struct {
	table        string
	where        []predicate.Predicate
	db           *db.DB
	valueCreator valuer.Value
	selects      []predicate.Selectable
	model        *register.Model
	sqlBuilder   SQLBuilder
}

func (s *Selector[T]) Build() (*Query, error) {
	s.sqlBuilder.sb.WriteString("select ")
	if s.model == nil {
		t := new(T)
		var err error
		s.model, err = s.db.R.Get(t) // 此处只使用了类型信息，可以new(T);无需保证和Get()方法是同个t
		if err != nil {
			return nil, err
		}
		s.sqlBuilder.registerModel(s.model)
	}

	if len(s.selects) > 0 {
		for i, sl := range s.selects {
			err := s.sqlBuilder.buildExpression(sl.(predicate.Expression))
			if err != nil {
				return nil, err
			}
			s.sqlBuilder.buildAs(sl.(predicate.Aliasable))
			if i < len(s.selects)-1 {
				s.sqlBuilder.sb.WriteString(", ")
			}
		}
	} else {
		s.sqlBuilder.sb.WriteString("*")
	}
	s.sqlBuilder.sb.WriteString(" from ")

	var table string
	if s.table == "" {
		table = s.model.TableName
	} else {
		table = s.table
	}
	s.sqlBuilder.sb.WriteString(table)

	length := len(s.where)
	if length > 0 {
		exp := s.where[0]
		if length > 1 {
			for _, expr := range s.where[1:] {
				exp = exp.And(expr)
			}
		}
		s.sqlBuilder.sb.WriteString(" where ")
		err := s.sqlBuilder.buildExpression(exp)
		if err != nil {
			return nil, err
		}
	}

	s.sqlBuilder.sb.WriteString(";")

	return &Query{
		SQL:  s.sqlBuilder.sb.String(),
		Args: s.sqlBuilder.Args,
	}, nil
}

func (s *Selector[T]) handleScalar(rows *sql.Rows) (any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 只有一列时，有可能是聚合函数、子查询
	if len(columns) == 1 {
		if _, ok := s.model.ColumnMap[columns[0]]; !ok {
			var val any
			err = rows.Scan(&val)
			if err != nil {
				return nil, err
			}
			return val, nil
		}
	}
	return nil, nil
}

func (s *Selector[T]) handleColumns(t *T, rows *sql.Rows) (*T, error) {
	v := s.db.ValueCreator(t, s.model)
	err := v.SetColumns(rows)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Selector[T]) Get(ctx context.Context) (any, error) {
	t := new(T)
	typ := reflect.TypeOf(t).Elem()

	k := typ.Kind()
	if k != reflect.Struct {
		return nil, fmt.Errorf("T 需要是结构体， 实际是：%s\n", k)
	}

	query, err := s.Build()
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = s.db.DB.QueryContext(ctx, query.SQL, query.Args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	if !rows.Next() {
		return nil, fmt.Errorf("没有数据")
	}
	var val any
	val, err = s.handleScalar(rows)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return s.handleColumns(t, rows)
	}

	return val, fmt.Errorf("没有数据")
}

func (s *Selector[T]) GetMutli(ctx context.Context) ([]any, error) {
	var results []any

	t := new(T)
	typ := reflect.TypeOf(t).Elem()

	k := typ.Kind()
	if k != reflect.Struct {
		return nil, fmt.Errorf("T 需要是结构体， 实际是：%s\n", k)
	}

	query, err := s.Build()
	if err != nil {
		return nil, err
	}

	var rows *sql.Rows
	rows, err = s.db.DB.QueryContext(ctx, query.SQL, query.Args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var val any
		val, err = s.handleScalar(rows)
		if err != nil {
			return nil, err
		}
		if val == nil {
			col := new(T)
			col, err = s.handleColumns(col, rows)
			if err != nil {
				return nil, err
			}
			results = append(results, col)
		} else {
			results = append(results, val)
		}
	}

	return results, nil
}

func (s *Selector[T]) From(tableName string) *Selector[T] {
	s.table = tableName
	return s
}

func (s *Selector[T]) Where(predicates ...predicate.Predicate) *Selector[T] {
	s.where = predicates
	return s
}

func (s *Selector[T]) Select(columns ...predicate.Selectable) *Selector[T] {
	s.selects = columns
	return s
}

func NewSelector[T any](db *db.DB) *Selector[T] {
	return &Selector[T]{
		db: db,
	}
}
