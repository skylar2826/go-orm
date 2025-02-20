package main

import (
	"context"
	"fmt"
	"geektime-go-orm/orm/data"
	db2 "geektime-go-orm/orm/db"
	"geektime-go-orm/orm/querier"
	"log"
)

func main() {
	dataSourceName := fmt.Sprint(Username, ":", Password, "@tcp(", Ip, ":", Port, ")/", DbName)
	//db2.WithReflectValue()
	db, err := db2.Open("mysql", dataSourceName)
	if err != nil {
		log.Println(err)
		return
	}
	ctx := context.Background()

	//predicate := predicate2.C("Username").Eq(predicate2.Valuer{Value: "test2"})

	// select u_username, is_active, birthdate from users;
	//selects := []predicate2.Selectable{predicate2.C("Username"), predicate2.C("IsActive"), predicate2.C("Birthdate")}
	//s := selector.NewSelector[data.User](db).Select(selects...)

	//select * from users; // 检查json数据正常返回
	s := querier.NewSelector[data.User](db)
	var val any
	val, err = s.Get(ctx)
	user := val.(*data.User)
	fmt.Println(user)

	// select AVG(`u_id`) from users;
	//selects := []predicate2.Selectable{predicate2.AVG("Id")}
	//s := selector.NewSelector[data.User](db).Select(selects...)
	//var val any
	//val, err = s.Get(ctx)
	//str := string(val.([]byte))
	//fmt.Println(str)

	// select COUNT(DISTINCT `u_id`) from users where `is_active` = true;
	//s := selector.NewSelector[data.User](db).Select(predicate.Raw("COUNT(DISTINCT `u_id`)")).Where(predicate.Raw("`is_active` = ?", true).AsPredicate())
	//var val any
	//val, err = s.Get(ctx)
	//if val != nil {
	//	count := val.(int64)
	//	fmt.Println(count)
	//} else {
	//	fmt.Println(val)
	//}

	//s := selector.NewSelector[data.User](db)
	//var users []*data.User
	//users, err = s.GetMutli(ctx)
	//fmt.Println(users)

	// select `u_id` as `avg_id` from users where `u_id` > 1;
	// 注意 alias 不能放在where 里， where在数据聚合和分组前执行，不认识别名
	//col := predicate.C("Id").As("avg_id")
	//s := selector.NewSelector[data.User](db).Select(col).Where(col.Gt(predicate.Valuer{Value: 1}))
	////var val any
	////val, err = s.Get(ctx)
	////fmt.Println(val.(int64))
	//
	//var val2 []any
	//val2, err = s.GetMutli(ctx)
	//fmt.Println(val2)
}
