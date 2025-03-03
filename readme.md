# orm 框架

![image](https://github.com/skylar2826/go-orm/blob/main/overview.png)

##  sql builder

### query

- type Query struct {
	SQL  string
	Args []any
}
- type QueryBuilder interface {
	Build() (*Query, error)
}

### querier

- type Querier[T any] interface {
	// Get 列、聚合函数、子查询
	Get(ctx context.Context) (any, error)
	GetMutli(ctx context.Context) ([]any, error)
}

var _ Querier[data.User] = &Selector[data.User]{}

	- selector

### executer

- type Executer interface {
	Execute(ctx context.Context) *orm.QueryResult
}

var _ Executer = &Inserter[data.User]{}

	- insector

		- onduplicateKey

			- type OnDuplicateKeyBuilder[T any] struct {
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
			- 使用示例： NewInserter[data.User](db).GetOnDuplicateKeyBuilder().Update(predicate.Assign("Email", "更新重复1357"))

## predicate

### Expression

- // Expression 标记 表示 语句 或 语句的部分
type Expression interface {
	expr()
}

// 确保实现了Expression
var _ Expression = Column{}
var _ Expression = Aggregate{}
var _ Expression = RawExpr{}
var _ Expression = Predicate{}
var _ Expression = Valuer{}

### Selectable

- // Selectable 标记 canSelect 列、聚合函数、子查询、表达式
// Selectable 的部分都可以设置别名，即Aliasable
type Selectable interface {
	Selectable()
}

var _ Selectable = Column{}
var _ Selectable = Aggregate{}
var _ Selectable = RawExpr{}

	- selector中通过.Select(select selectable)使用

### Aggregate // 聚合函数 AVG(columnName)

### Assignable

- // Assignable 标记接口
// 实现该接口意味着用于赋值语句
// 用于UPDATE、UPSERT 语句
type Assignable interface {
	Assign()
}

	- insertor中常在BuildOnduplicateKey时使用

### Column

- 使用示例：
predicate.C("Id").Lt(predicate.Valuer{Value: val})

### Predicate

- type Predicate struct {
	Left  Expression
	Op    Op
	Right Expression
}

	- selector中.where(predicates...predicate.Predicate)

### RawExpr // 允许用户自定义输入

## db

### dialect

- type Dialect interface {
	Quoter() string
	BuildOnDuplicateKey(OnDuplicateKey *OnDuplicateKey) (*OnDuplicateKeyStatement, error)
}

var _ Dialect = &StandardSQL{}
var _ Dialect = &MysqlSQL{}

	- standardSql

		- quoter -> 双引号
duplicateKey语句： on conflict columnName do update set col1=exclued.col1

	- mysqlSql

		- quoter -> `
duplicateKey语句：on duplicate key update col1=values(col1)

- onduplicatekey

	- type OnDuplicateKey struct {
	Assigns         []predicate.Assignable // 支持两种形式 col=values(col)、col=具体值
	ConflictColumns []predicate.Column // 冲突列 standardSql语句构造使用
}

### session(上下文)

- // Session 是 db、tx的顶层抽象
type Session interface {
	GetCore() Core
	ExecContext(ctx context.Context, query string, args ...any) *sql2.QueryResult
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

var _ Session = &DB{}
var _ Session = &Tx{}

	- 使用示例：
func NewSelector[T any](sess session.Session) *Selector[T] {
	return &Selector[T]{
		sess:       sess,
		SQLBuilder: sql2.SQLBuilder{Core: sess.GetCore()},
	}
}


s := selector.NewSelector[data.User](tx) // 在selector中传入session(db | tx)在get、getMuti时调用相应的queryCtx方法执行语句拿到结果

- db

	- 支持doTx（闭包事务）、beginTx

- tx

### model

- 维护User和数据库中相应表、字段的关系
- register(注册中心)

	- 维护models，支持注册model、缓存model、get model

### valuer

- // Value 是对结构体实例的内部抽象
type Value interface {
	SetColumns(rows *sql.Rows) error
	Field(name string) (any, error)
}

type ValueCreator func(val any, model *register.Model) Value

	- 使用场景：
SetColumns(): selector中将db返回结果通过row.Scan塞进&user{}中

Field(): insertor中
遍历&user{Name: "zly", Birthdate: "1999-06-06"}获取每个字段及其值生成相应的sql语句

- ValueCreator 

	- 支持unsafe、reflect两种方式，unsafe性能更好，默认使用unsafe

## aop

### type Handler func(ctx context.Context, qc *orm.QueryContext) *orm.QueryResult

type Middleware func(next Handler) Handler


### 做数据观测、统计时用

## 代码自动生成技术

### 依靠ast生成常用模板语句
比如生成：
  func UserIdLt(val int) predicate.Predicate {
                return predicate.C("Id").Lt(predicate.Valuer{Value: val})

s := selector.NewSelector[data.User](tx).Where(data.UserIdLt(10))
            }

## 设计风格：
1. builder模式
2. 泛型约束

