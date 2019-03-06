package dbs

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var Conns *sql.DB

func init() {
	var err error
	Conns, err = sql.Open("mysql", "root:123456@tcp(192.168.33.11)/test")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = Conns.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
}
