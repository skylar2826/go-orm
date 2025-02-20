package predicate

const (
	opAVG   = "AVG"
	opSUM   = "SUM"
	opMAX   = "MAX"
	opMIN   = "MIN"
	opCOUNT = "COUNT"
)

type Aggregate struct {
	Op    Op
	Arg   string
	Alias string
}

func (a Aggregate) Aliasable() string {
	return a.Alias
}

func (a Aggregate) expr() {
}

func (a Aggregate) Selectable() {
}

func (a Aggregate) As(alias string) Aggregate {
	return Aggregate{
		Op:    a.Op,
		Arg:   a.Arg,
		Alias: alias,
	}
}

func AVG(arg string) Aggregate {
	return Aggregate{
		Op:  opAVG,
		Arg: arg,
	}
}

// SUM 支持SUM(colName), 暂未支持SUM(a * b)
func SUM(arg string) Aggregate {
	return Aggregate{
		Op:  opSUM,
		Arg: arg,
	}
}

func MAX(arg string) Aggregate {
	return Aggregate{
		Op:  opMAX,
		Arg: arg,
	}
}

func MIN(arg string) Aggregate {
	return Aggregate{
		Op:  opMIN,
		Arg: arg,
	}
}

func COUNT(arg string) Aggregate {
	return Aggregate{
		Op:  opCOUNT,
		Arg: arg,
	}
}
