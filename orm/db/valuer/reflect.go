package valuer

import (
	"database/sql"
	"fmt"
	"geektime-go-orm/orm/db/register"
	"reflect"
)

type reflectValue struct {
	val   reflect.Value
	model *register.Model
}

func (r *reflectValue) SetColumns(rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(columns) > len(r.model.ColumnMap) {
		return fmt.Errorf("数据库表列数与元数据列数不相等")
	}

	colValues := make([]any, len(columns))
	colElemValues := make([]reflect.Value, len(columns))

	for i, c := range columns {
		field := r.model.ColumnMap[c]
		val := reflect.New(field.Typ)
		colValues[i] = val.Interface()
		colElemValues[i] = val.Elem()
	}
	err = rows.Scan(colValues...)
	if err != nil {
		return err
	}
	for i, c := range columns {
		field := r.model.ColumnMap[c]
		fd := r.val.FieldByName(field.MetaName)
		if fd.CanSet() {
			fd.Set(colElemValues[i])
		}
	}
	return nil
}

func NewReflectValue(val any, model *register.Model) Value {
	return &reflectValue{
		val:   reflect.ValueOf(val).Elem(),
		model: model,
	}
}

var _ ValueCreator = NewReflectValue
