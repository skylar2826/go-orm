package dialect

import (
	"fmt"
	"geektime-go-orm/orm/db/register"
	"geektime-go-orm/orm/errors"
	"geektime-go-orm/orm/predicate"
	"strings"
)

type OnDuplicateKey struct {
	Assigns         []predicate.Assignable
	ConflictColumns []predicate.Column
}

type OnDuplicateKeyStatement struct {
	Sql  string
	Args []any
}

type Dialect interface {
	Quoter() string
	BuildOnDuplicateKey(OnDuplicateKey *OnDuplicateKey) (*OnDuplicateKeyStatement, error)
	RegisterModel(model *register.Model)
}

var _ Dialect = &StandardSQL{}
var _ Dialect = &MysqlSQL{}

type StandardSQL struct {
	sqlBuilder strings.Builder
	model      *register.Model
}

func (s *StandardSQL) Quoter() string {
	return "\""
}

func (s *StandardSQL) buildColumn(column predicate.Column) error {
	if s.model == nil {
		return fmt.Errorf("builder: model不存在")
	}
	if name, ok := s.model.Fields[column.Name]; !ok {
		return errors.FieldNotFoundErr(column.Name)
	} else {
		s.sqlBuilder.WriteString(fmt.Sprintf("%s%s%s", s.Quoter(), name.ColName, s.Quoter()))
	}
	return nil
}

func (s *StandardSQL) BuildOnDuplicateKey(OnDuplicateKey *OnDuplicateKey) (*OnDuplicateKeyStatement, error) {
	s.sqlBuilder.WriteString(" ON CONFLICT ")

	if len(OnDuplicateKey.ConflictColumns) > 0 {
		s.sqlBuilder.WriteString("(")
		for idx, col := range OnDuplicateKey.ConflictColumns {
			err := s.buildColumn(col)
			if err != nil {
				return nil, err
			}
			if idx < len(OnDuplicateKey.ConflictColumns)-1 {
				s.sqlBuilder.WriteString(",")
			}
		}
		s.sqlBuilder.WriteString(")")
	}

	s.sqlBuilder.WriteString(" DO UPDATE SET ")
	if len(OnDuplicateKey.Assigns) > 0 {
		for idx, assign := range OnDuplicateKey.Assigns {
			col, ok := assign.(predicate.Column)
			if !ok {
				return nil, fmt.Errorf("assignable 类型错误：%s\n", assign)
			}
			err := s.buildColumn(col)
			if err != nil {
				return nil, err
			}
			s.sqlBuilder.WriteString(" = Excluded.")
			s.sqlBuilder.WriteString(col.Name)

			if idx < len(OnDuplicateKey.Assigns)-1 {
				s.sqlBuilder.WriteString(",")
			}
		}
	}

	return &OnDuplicateKeyStatement{
		Sql: s.sqlBuilder.String(),
	}, nil
}

func (s *StandardSQL) RegisterModel(model *register.Model) {
	s.model = model
}

func NewStandardSQL() Dialect {
	return &StandardSQL{}
}

type MysqlSQL struct {
	StandardSQL
	Assigns []any
}

func (m *MysqlSQL) Quoter() string {
	return "`"
}

func (m *MysqlSQL) buildColumn(column predicate.Column) error {
	if m.model == nil {
		return fmt.Errorf("builder: model不存在")
	}
	if name, ok := m.model.Fields[column.Name]; !ok {
		return errors.FieldNotFoundErr(column.Name)
	} else {
		m.sqlBuilder.WriteString(fmt.Sprintf("%s%s%s", m.Quoter(), name.ColName, m.Quoter()))
	}
	return nil
}

func (m *MysqlSQL) BuildAssign(assign predicate.Assignable) error {
	switch expr := assign.(type) {
	case predicate.Assignment:
		c := predicate.C(expr.ColName)
		err := m.buildColumn(c)
		if err != nil {
			return err
		}
		m.sqlBuilder.WriteString(" = ?")
		m.Assigns = append(m.Assigns, expr.Val)
		return nil
	case predicate.Column:
		err := m.buildColumn(expr)
		if err != nil {
			return err
		}
		m.sqlBuilder.WriteString(" = Values(")
		err = m.buildColumn(expr)
		if err != nil {
			return err
		}
		m.sqlBuilder.WriteString(")")
	default:
		return fmt.Errorf("无法识别的Assignable： %builder\n", expr)
	}
	return nil
}

func (m *MysqlSQL) BuildOnDuplicateKey(OnDuplicateKey *OnDuplicateKey) (*OnDuplicateKeyStatement, error) {
	m.sqlBuilder.WriteString(" ON DUPLICATE KEY UPDATE ")
	for idx, assign := range OnDuplicateKey.Assigns {
		err := m.BuildAssign(assign)
		if err != nil {
			return nil, err
		}
		if idx < len(OnDuplicateKey.Assigns)-1 {
			m.sqlBuilder.WriteString(",")
		}
	}
	return &OnDuplicateKeyStatement{
		Sql:  m.sqlBuilder.String(),
		Args: m.Assigns,
	}, nil
}

func NewMysqlSQL() Dialect {
	return &MysqlSQL{}
}
