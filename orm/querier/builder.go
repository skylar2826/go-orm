package querier

import (
	"fmt"
	"geektime-go-orm/orm/db/register"
	"geektime-go-orm/orm/errors"
	"geektime-go-orm/orm/predicate"
	"strings"
)

type SQLBuilder struct {
	sb    strings.Builder
	Args  []any
	model *register.Model
}

func (s *SQLBuilder) registerModel(model *register.Model) *SQLBuilder {
	s.model = model
	return s
}

func (s *SQLBuilder) buildAs(val predicate.Aliasable) {
	alias := val.Aliasable()
	if alias != "" {
		s.sb.WriteString(fmt.Sprintf(" AS `%s`", alias))
	}
}

func (s *SQLBuilder) buildAggregate(aggregate predicate.Aggregate) error {
	s.sb.WriteString(aggregate.Op.String())
	s.sb.WriteString("(")
	col := predicate.C(aggregate.Arg)
	err := s.buildExpression(col)
	if err != nil {
		return err
	}
	s.sb.WriteString(")")
	return nil
}

func (s *SQLBuilder) buildColumn(column predicate.Column) error {
	if s.model == nil {
		return fmt.Errorf("SQLBuilder: model不存在")
	}
	if name, ok := s.model.Fields[column.Name]; !ok {
		return errors.FieldNotFoundErr(column.Name)
	} else {
		s.sb.WriteString(fmt.Sprintf("`%s`", name.ColName))
	}
	return nil
}

func (s *SQLBuilder) buildPredicate(predicate2 predicate.Predicate) error {
	if predicate2.Left != nil {
		s.sb.WriteString("(")
		err := s.buildExpression(predicate2.Left)
		if err != nil {
			return err
		}
		s.sb.WriteString(")")
	}
	s.sb.WriteString(fmt.Sprintf(" %s ", predicate2.Op))
	if predicate2.Right != nil {
		s.sb.WriteString("(")
		err := s.buildExpression(predicate2.Right)
		if err != nil {
			return err
		}
		s.sb.WriteString(")")
	}
	return nil
}

func (s *SQLBuilder) buildExpression(exp predicate.Expression) error {
	switch expr := exp.(type) {
	case predicate.Column:
		err := s.buildColumn(expr)
		if err != nil {
			return err
		}
	case predicate.Aggregate:
		err := s.buildAggregate(expr)
		if err != nil {
			return err
		}
	case predicate.Valuer:
		s.sb.WriteString("?")
		s.Args = append(s.Args, expr.Value)
	case predicate.Predicate:
		err := s.buildPredicate(expr)
		if err != nil {
			return err
		}
	case predicate.RawExpr:
		s.sb.WriteString(expr.Sql)
		s.Args = append(s.Args, expr.Args...)
	default:
		return fmt.Errorf("无法识别的Expression： %s\n", expr)
	}
	return nil
}
