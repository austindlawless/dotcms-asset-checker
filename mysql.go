package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

func NewMySql(user string, password string, host string, dbName string) *MySql {
	mysql := new(MySql)
	con, err := sql.Open("mysql", user+":"+password+"@tcp("+host+":3306)/"+dbName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mysql.db = con

	return mysql
}

type MySql struct {
	db *sql.DB
}

func (m *MySql) Close() {
	m.db.Close()
}
