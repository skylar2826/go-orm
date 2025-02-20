package predicate

type RawExpr struct {
	Args  []any
	Sql   string
	Alias string
}

func (r RawExpr) Aliasable() string {
	return r.Alias
}

func (r RawExpr) Selectable() {
}

func (r RawExpr) expr() {
}

func (r RawExpr) AsPredicate() Predicate {
	return Predicate{
		Left: r,
	}
}

func (r RawExpr) As(alias string) RawExpr {
	return RawExpr{
		Args:  r.Args,
		Sql:   r.Sql,
		Alias: alias,
	}
}

func Raw(sql string, args ...any) RawExpr {
	return RawExpr{
		Sql:  sql,
		Args: args,
	}
}
