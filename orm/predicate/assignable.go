package predicate

// Assignable 标记接口
// 实现该接口意味着用于赋值语句
// 用于UPDATE、UPSERT 语句
type Assignable interface {
	Assign()
}

type Assignment struct {
	ColName string
	Val     any
}

func (a Assignment) expr() {
}

func (a Assignment) Assign() {
}

func Assign(colName string, val any) Assignment {
	return Assignment{
		ColName: colName,
		Val:     val,
	}
}

var _ Assignable = Assignment{}
var _ Assignable = Column{}
