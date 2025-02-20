package predicate

const (
	opEq  = "="
	opLt  = "<"
	opGt  = ">"
	opAnd = "and"
	opOr  = "or"
	opNot = "not"
)

type Predicate struct {
	Left  Expression
	Op    Op
	Right Expression
}

func (p Predicate) expr() {
}

func (p Predicate) And(right Expression) Predicate {
	return Predicate{
		Left:  p,
		Op:    opAnd,
		Right: right,
	}
}

func (p Predicate) Or(right Expression) Predicate {
	return Predicate{
		Left:  p,
		Op:    opOr,
		Right: right,
	}
}

func Not(right Expression) Predicate {
	return Predicate{
		Op:    opNot,
		Right: right,
	}
}
