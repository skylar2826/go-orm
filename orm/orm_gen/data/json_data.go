package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JsonData[T any] struct {
	Val   T
	Valid bool
}

// Scan 数据库类型转go类型
func (j JsonData[T]) Scan(value any) error {
	if value == nil {
		j.Valid = false
		return nil
	}

	var bytes []byte
	if val, ok := value.(string); ok {
		bytes = []byte(val)
	}

	switch val := value.(type) {
	case []byte:
		bytes = val
	case string:
		bytes = []byte(val)
	default:
		return fmt.Errorf("不支持的类型：%T\n", val)
	}

	return json.Unmarshal(bytes, &j.Val)
}

// Value go类型转数据库类型
func (j JsonData[T]) Value() (driver.Value, error) {
	if !j.Valid {
		return nil, nil
	}
	data, err := json.Marshal(j.Val)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}
