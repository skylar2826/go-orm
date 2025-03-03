package sql

import (
	"fmt"
	"geektime-go-orm/orm/db/session"
	"geektime-go-orm/orm/errors"
	"geektime-go-orm/orm/predicate"
	"strings"
)

type SQLBuilder struct {
	Sb   strings.Builder
	Args []any
	session.Core
}

func (s *SQLBuilder) BuildAs(val predicate.Aliasable) {
	alias := val.Aliasable()
	if alias != "" {
		s.Sb.WriteString(fmt.Sprintf(" AS %s%s%s", s.Dialect.Quoter(), alias, s.Dialect.Quoter()))
	}
}

func (s *SQLBuilder) buildAggregate(aggregate predicate.Aggregate) error {
	s.Sb.WriteString(aggregate.Op.String())
	s.Sb.WriteString("(")
	col := predicate.C(aggregate.Arg)
	err := s.BuildExpression(col)
	if err != nil {
		return err
	}
	s.Sb.WriteString(")")
	return nil
}

func (s *SQLBuilder) BuildColumn(column predicate.Column) error {
	if s.Model == nil {
		return fmt.Errorf("SQLBuilder: model不存在")
	}
	if name, ok := s.Model.Fields[column.Name]; !ok {
		return errors.FieldNotFoundErr(column.Name)
	} else {
		s.Sb.WriteString(fmt.Sprintf("%s%s%s", s.Dialect.Quoter(), name.ColName, s.Dialect.Quoter()))
	}
	return nil
}

func (s *SQLBuilder) buildPredicate(predicate2 predicate.Predicate) error {
	if predicate2.Left != nil {
		s.Sb.WriteString("(")
		err := s.BuildExpression(predicate2.Left)
		if err != nil {
			return err
		}
		s.Sb.WriteString(")")
	}
	s.Sb.WriteString(fmt.Sprintf(" %s ", predicate2.Op))
	if predicate2.Right != nil {
		s.Sb.WriteString("(")
		err := s.BuildExpression(predicate2.Right)
		if err != nil {
			return err
		}
		s.Sb.WriteString(")")
	}
	return nil
}

func (s *SQLBuilder) BuildValuer(valuer predicate.Valuer) error {
	s.Sb.WriteString("?")
	s.Args = append(s.Args, valuer.Value)
	return nil
}

func (s *SQLBuilder) BuildExpression(exp predicate.Expression) error {
	switch expr := exp.(type) {
	case predicate.Column:
		err := s.BuildColumn(expr)
		if err != nil {
			return err
		}
	case predicate.Aggregate:
		err := s.buildAggregate(expr)
		if err != nil {
			return err
		}
	case predicate.Valuer:
		err := s.BuildValuer(expr)
		if err != nil {
			return err
		}
	case predicate.Predicate:
		err := s.buildPredicate(expr)
		if err != nil {
			return err
		}
	case predicate.RawExpr:
		s.Sb.WriteString(expr.Sql)
		s.Args = append(s.Args, expr.Args...)
	default:
		return fmt.Errorf("无法识别的Expression： %s\n", expr)
	}
	return nil
}
