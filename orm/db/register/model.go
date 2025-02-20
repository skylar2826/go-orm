package register

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type field struct {
	ColName   string
	Typ       reflect.Type
	MetaName  string
	Offset    uintptr
	AliasName string
}

type TableName interface {
	TableName() string
}

type Model struct {
	TableName string
	Fields    map[string]*field // 元数据到字段的映射
	ColumnMap map[string]*field // 数据库列名到字段的映射
}

func (m *Model) ParseModel(val any) error {
	typ := reflect.TypeOf(val).Elem()
	if tableName, ok := val.(TableName); ok {
		m.TableName = tableName.TableName()
	} else {
		m.TableName = typ.Name()
	}

	k := typ.Kind()
	if k != reflect.Struct {
		return fmt.Errorf("model value 不是结构体指针, 实际是： %s\n", k)
	}

	numField := typ.NumField()
	m.Fields = make(map[string]*field, numField)
	m.ColumnMap = make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		f := typ.Field(i)
		tagName := f.Tag.Get("orm")
		if tagName != "" {
			fd := &field{
				ColName:  tagName,
				Typ:      f.Type,
				MetaName: f.Name,
				Offset:   f.Offset,
			}
			m.Fields[f.Name] = fd
			m.ColumnMap[tagName] = fd
		} else {
			colName := underlineCase(f.Name)
			fd := &field{
				ColName:  colName,
				Typ:      f.Type,
				MetaName: f.Name,
				Offset:   f.Offset,
			}
			m.Fields[f.Name] = fd
			m.ColumnMap[colName] = fd
		}

	}
	return nil
}

func underlineCase(val string) string {
	var b strings.Builder
	for i, v := range val {
		if i > 0 && unicode.IsUpper(v) {
			b.WriteString("_")
		}
		b.WriteRune(unicode.ToLower(v))
	}
	return b.String()
}
