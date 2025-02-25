package dialect

import (
	"fmt"
	"geektime-go-orm/orm/predicate"
	"geektime-go-orm/orm/sqlCommonBuilder"
)

type OnDuplicateKey struct {
	Assigns         []predicate.Assignable
	ConflictColumns []predicate.Column
}

type Dialect interface {
	Quoter() string
	BuildOnDuplicateKey(builder *sqlCommonBuilder.SQLBuilder, OnDuplicateKey *OnDuplicateKey) error
}

var _ Dialect = &StandardSQL{}
var _ Dialect = &MysqlSQL{}

type StandardSQL struct {
}

func (s *StandardSQL) Quoter() string {
	return "\""
}

func (s *StandardSQL) BuildOnDuplicateKey(builder *sqlCommonBuilder.SQLBuilder, OnDuplicateKey *OnDuplicateKey) error {
	builder.Sb.WriteString(" ON CONFLICT ")

	if len(OnDuplicateKey.ConflictColumns) > 0 {
		builder.Sb.WriteString("(")
		for idx, col := range OnDuplicateKey.ConflictColumns {
			err := builder.BuildColumn(col)
			if err != nil {
				return err
			}
			if idx < len(OnDuplicateKey.ConflictColumns)-1 {
				builder.Sb.WriteString(",")
			}
		}
		builder.Sb.WriteString(")")
	}

	builder.Sb.WriteString(" DO UPDATE SET ")
	if len(OnDuplicateKey.Assigns) > 0 {
		for idx, assign := range OnDuplicateKey.Assigns {
			col, ok := assign.(predicate.Column)
			if !ok {
				return fmt.Errorf("assignable 类型错误：%s\n", assign)
			}
			err := builder.BuildColumn(col)
			if err != nil {
				return err
			}
			builder.Sb.WriteString(" = Excluded.")
			builder.Sb.WriteString(col.Name)

			if idx < len(OnDuplicateKey.Assigns)-1 {
				builder.Sb.WriteString(",")
			}
		}
	}

	return nil
}

func NewStandardSQL() *StandardSQL {
	return &StandardSQL{}
}

type MysqlSQL struct {
	StandardSQL
	Assigns []predicate.Assignable
}

func (m *MysqlSQL) Quoter() string {
	return "`"
}

func (m *MysqlSQL) BuildAssign(builder *sqlCommonBuilder.SQLBuilder, assign predicate.Assignable) error {
	switch expr := assign.(type) {
	case predicate.Assignment:
		c := predicate.C(expr.ColName)
		err := builder.BuildColumn(c)
		if err != nil {
			return err
		}
		builder.Sb.WriteString(" = ?")
		builder.Args = append(builder.Args, expr.Val)
		return nil
	case predicate.Column:
		err := builder.BuildColumn(expr)
		if err != nil {
			return err
		}
		builder.Sb.WriteString(" = Values(")
		err = builder.BuildColumn(expr)
		if err != nil {
			return err
		}
		builder.Sb.WriteString(")")
	default:
		return fmt.Errorf("无法识别的Assignable： %builder\n", expr)
	}
	return nil
}

func (m *MysqlSQL) BuildOnDuplicateKey(builder *sqlCommonBuilder.SQLBuilder, OnDuplicateKey *OnDuplicateKey) error {
	builder.Sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	for idx, assign := range OnDuplicateKey.Assigns {
		err := m.BuildAssign(builder, assign)
		if err != nil {
			return err
		}
		if idx < len(OnDuplicateKey.Assigns)-1 {
			builder.Sb.WriteString(",")
		}
	}
	return nil
}

func NewMysqlSQL() *MysqlSQL {
	return &MysqlSQL{}
}
