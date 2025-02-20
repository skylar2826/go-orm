package predicate

type Column struct {
	Name  string
	Alias string
}

func (c Column) Aliasable() string {
	return c.Alias
}

func (c Column) Selectable() {
}

func (c Column) expr() {
}

func (c Column) Eq(expr Expression) Predicate {
	return Predicate{
		Left:  c,
		Op:    opEq,
		Right: expr,
	}
}

func (c Column) Lt(expr Expression) Predicate {
	return Predicate{
		Left:  c,
		Op:    opLt,
		Right: expr,
	}
}

func (c Column) Gt(expr Expression) Predicate {
	return Predicate{
		Left:  c,
		Op:    opGt,
		Right: expr,
	}
}

func (c Column) As(alias string) Column {
	// 这种修改不会生效
	//c.Alias = alias
	return Column{
		Name:  c.Name,
		Alias: alias,
	}
}

// C C("id").Eq(1)
func C(name string) Column {
	return Column{Name: name}
}
