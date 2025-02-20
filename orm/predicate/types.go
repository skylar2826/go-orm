package predicate

type Op string

func (o Op) String() string {
	return string(o)
}

// Expression 标记 表示 语句 或 语句的部分
type Expression interface {
	expr()
}

// 确保实现了Expression
var _ Expression = Column{}
var _ Expression = Aggregate{}
var _ Expression = RawExpr{}
var _ Expression = Predicate{}
var _ Expression = Valuer{}

// Selectable 标记 canSelect 列、聚合函数、子查询、表达式
// Selectable 的部分都可以设置别名，即Aliasable
type Selectable interface {
	Selectable()
}

var _ Selectable = Column{}
var _ Selectable = Aggregate{}
var _ Selectable = RawExpr{}

// Aliasable 标记 canAlias
type Aliasable interface {
	Aliasable() string
}

var _ Aliasable = Column{}
var _ Aliasable = Aggregate{}
var _ Aliasable = RawExpr{}
