package errors

import "fmt"

var FieldNotFoundErr = func(str string) error {
	return fmt.Errorf("field not found: %s\n", str)
}
