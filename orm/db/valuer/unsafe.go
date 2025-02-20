package valuer

import (
	"database/sql"
	"fmt"
	"geektime-go-orm/orm/db/register"
	"reflect"
	"unsafe"
)

type UnsafeValue struct {
	val   unsafe.Pointer
	model *register.Model
}

func (u *UnsafeValue) SetColumns(rows *sql.Rows) error {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if len(columns) > len(u.model.ColumnMap) {
		return fmt.Errorf("数据库表中的列数与元数据列数不相等")
	}

	colValues := make([]any, len(columns))
	for i, c := range columns {
		field := u.model.ColumnMap[c]
		addr := unsafe.Pointer(uintptr(u.val) + field.Offset)
		val := reflect.NewAt(field.Typ, addr)
		colValues[i] = val.Interface()
	}
	return rows.Scan(colValues...)
}

func NewUnsafeValue(val any, model *register.Model) Value {
	return &UnsafeValue{
		// reflect.ValueOf(val).Pointer() 返回的是uintptr(具体的地址值 比如0x8888)
		// 转成 unsafe.Pointer（指针）好处是：GC维护unsafe.Pointer，数据回收后GC会标记-复制数据到新地址（0x8888），更新unsafe.Pointer指向的地址(0x8888)
		// 不转成指针，使用uintptr,在数据回收处理更换地址(比如 0x0001)后，uintptr的地址值(0x8888)就不是数据(0x0001)现在的位置了
		val:   unsafe.Pointer(reflect.ValueOf(val).Pointer()),
		model: model,
	}
}

var _ ValueCreator = NewUnsafeValue
