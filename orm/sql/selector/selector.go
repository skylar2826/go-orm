package selector

import (
	"context"
	"database/sql"
	"fmt"
	"geektime-go-orm/orm"
	"geektime-go-orm/orm/db/session"
	"geektime-go-orm/orm/db/valuer"
	"geektime-go-orm/orm/orm_gen/data"
	"geektime-go-orm/orm/predicate"
	sql2 "geektime-go-orm/orm/sql"
	"reflect"
)

type Querier[T any] interface {
	// Get 列、聚合函数、子查询
	Get(ctx context.Context) (any, error)
	GetMutli(ctx context.Context) ([]any, error)
}

var _ Querier[data.User] = &Selector[data.User]{}

type Selector[T any] struct {
	table        string
	where        []predicate.Predicate
	sess         session.Session
	valueCreator valuer.Value
	selects      []predicate.Selectable

	sql2.SQLBuilder
}

func (s *Selector[T]) Build() (*orm.Query, error) {
	s.Sb.WriteString("select ")
	if s.Model == nil {
		t := new(T)
		var err error
		s.Model, err = s.R.Get(t) // 此处只使用了类型信息，可以new(T);无需保证和Get()方法是同个t
		if err != nil {
			return nil, err
		}
	}

	if len(s.selects) > 0 {
		for i, sl := range s.selects {
			err := s.BuildExpression(sl.(predicate.Expression))
			if err != nil {
				return nil, err
			}
			s.BuildAs(sl.(predicate.Aliasable))
			if i < len(s.selects)-1 {
				s.Sb.WriteString(", ")
			}
		}
	} else {
		s.Sb.WriteString("*")
	}
	s.Sb.WriteString(" from ")

	var table string
	if s.table == "" {
		table = s.Model.TableName
	} else {
		table = s.table
	}
	s.Sb.WriteString(table)

	length := len(s.where)
	if length > 0 {
		exp := s.where[0]
		if length > 1 {
			for _, expr := range s.where[1:] {
				exp = exp.And(expr)
			}
		}
		s.Sb.WriteString(" where ")
		err := s.BuildExpression(exp)
		if err != nil {
			return nil, err
		}
	}

	s.Sb.WriteString(";")

	return &orm.Query{
		SQL:  s.Sb.String(),
		Args: s.Args,
	}, nil
}

func (s *Selector[T]) handleScalar(rows *sql.Rows) (any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 只有一列时，有可能是聚合函数、子查询
	if len(columns) == 1 {
		if _, ok := s.Model.ColumnMap[columns[0]]; !ok {
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
	v := s.ValueCreator(t, s.Model)
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
	rows, err = s.sess.QueryContext(ctx, query.SQL, query.Args...)
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
	rows, err = s.sess.QueryContext(ctx, query.SQL, query.Args...)
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

func NewSelector[T any](sess session.Session) *Selector[T] {
	return &Selector[T]{
		sess:       sess,
		SQLBuilder: sql2.SQLBuilder{Core: sess.GetCore()},
	}
}
