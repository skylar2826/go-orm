package main

//dataSourceName := fmt.Sprint(Username, ":", Password, "@tcp(", Ip, ":", Port, ")/", DbName)
////db2.WithReflectValue()
//db, err := db2.Open("mysql", dataSourceName, db2.WithDialect(dialect.NewMysqlSQL()))
//if err != nil {
//log.Println(err)
//return
//}
//ctx := context.Background()
//user1 := &data.User{
//Username:  "zly",
//Id:        18,
//Email:     "135",
//Birthdate: "2024-06-29",
//IsActive:  true,
//BaseInfo: data.JsonData[data.BaseInfo]{
//Valid: true,
//Val: data.BaseInfo{
//Detail:      "111",
//Description: "5555",
//},
//},
//}
//
//user2 := &data.User{
//Username: "zly2",
//Id:       19,
//Email:    "136",
//}
//
//observer := aop.NewObserverMiddleBuilder(aop.WithTracer()).Build()
//
//// 测试插入部分列 Insert into users (`u_username`,`u_id`,`email`) Values (?,?,?),(?,?,?);
//i := sql2.NewInserter[data.User](db, observer).Columns(predicate.C("Username"), predicate.C("Id"), predicate.C("Email")).Values(user1, user2)
//
//// 测试插入全部列
////i := sql2.NewInserter[data.User](db).Values(user1)
//
////Insert into users (`u_username`,`u_id`,`email`)  Values (?,?,?) ON DUPLICATE KEY UPDATE `email` = ?,`u_id` = ?;
////i := querier.NewInserter[data.User](db).Columns(predicate.C("Username"), predicate.C("Id"), predicate.C("Email")).Values(user1).GetOnDuplicateKeyBuilder().Update(predicate.Assign("Email", "更新重复1357"), predicate.Assignment{ColName: "Id", Val: 12})
//
////Insert into users (`u_username`,`u_id`,`email`)  Values (?,?,?) ON DUPLICATE KEY UPDATE `birthdate` = Values(`birthdate`);
////i := sql2.NewInserter[data.User](db).Columns(predicate.C("Username"), predicate.C("Email"), predicate.C("Id"), predicate.C("Birthdate")).Values(user1).GetOnDuplicateKeyBuilder().Update(predicate.C("Birthdate"))
//
////i := sql2.NewInserter[data.User](db).Columns(predicate.C("Username"), predicate.C("Email"), predicate.C("Id"), predicate.C("Birthdate")).Values(user1).GetOnDuplicateKeyBuilder().ConflictColumns(predicate.C("Email")).Update(predicate.C("Birthdate"))
//
//res := i.Execute(ctx)
//if res.Err != nil {
//log.Println(err)
//return
//}
//
//var idx int64
//idx, err = res.Result.(sql.Result).LastInsertId()
//if err != nil {
//log.Println(err)
//}
//
//fmt.Println(idx)
