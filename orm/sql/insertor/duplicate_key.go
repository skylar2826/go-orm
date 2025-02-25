package insertor

import (
	"geektime-go-orm/orm/db/dialect"
	"geektime-go-orm/orm/predicate"
)

type OnDuplicateKeyBuilder[T any] struct {
	i               *Inserter[T]
	conflictColumns []predicate.Column
}

func (o *OnDuplicateKeyBuilder[T]) Update(assigns ...predicate.Assignable) *Inserter[T] {
	o.i.OnDuplicateKey = &dialect.OnDuplicateKey{
		Assigns:         assigns,
		ConflictColumns: o.conflictColumns,
	}
	return o.i
}

func (o *OnDuplicateKeyBuilder[T]) ConflictColumns(cols ...predicate.Column) *OnDuplicateKeyBuilder[T] {
	o.conflictColumns = cols
	return o
}

func NewOnDuplicateBuilder[T any](i *Inserter[T]) *OnDuplicateKeyBuilder[T] {
	return &OnDuplicateKeyBuilder[T]{
		i: i,
	}
}
