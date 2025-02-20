package register

import (
	"fmt"
	"reflect"
	"sync"
)

type Register struct {
	Models map[string]*Model
	mutex  sync.Mutex
}

func (r *Register) Get(val any) (*Model, error) {
	if r.Models == nil {
		return nil, fmt.Errorf("models 不存在")
	}
	typ := reflect.TypeOf(val)
	name := typ.Name()
	r.mutex.Lock()
	m, ok := r.Models[name]
	r.mutex.Unlock()
	if !ok {
		r.mutex.Lock()
		m = &Model{}
		err := m.ParseModel(val)
		if err != nil {
			return nil, err
		}
		r.Models[name] = m
		r.mutex.Unlock()
	}
	return m, nil
}
