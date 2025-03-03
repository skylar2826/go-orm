package main

import (
	"context"
	"database/sql"
	"fmt"
	"geektime-go-orm/orm/db/dialect"
	"geektime-go-orm/orm/db/session"
	"geektime-go-orm/orm/orm_gen/data"
	"geektime-go-orm/orm/sql/selector"
	"log"
)

func main() {
	dataSourceName := fmt.Sprint(Username, ":", Password, "@tcp(", Ip, ":", Port, ")/", DbName)
	db, err := session.Open("mysql", dataSourceName, session.WithDialect(dialect.NewMysqlSQL()))
	if err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()
	var tx *session.Tx
	ctx, tx, err = db.BeginTxV2(ctx, &sql.TxOptions{})

	s := selector.NewSelector[data.User](tx).Where(data.UserIdLt(10))
	var res []any
	res, err = s.GetMutli(ctx)
	fmt.Println(res)

}
