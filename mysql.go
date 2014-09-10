package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type DatabaseNameConfig map[string]string

func NewMySql(user string, password string, host string, dbNames DatabaseNameConfig) *MySql {
	mysql := new(MySql)
	con, err := sql.Open("mysql", user+":"+password+"@tcp("+host+":3306)/"+dbNames["SegundoDb"])

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mysql.db = con
	mysql.dbNames = dbNames

	return mysql
}

type MySql struct {
	db      *sql.DB
	dbNames DatabaseNameConfig
}

func (m *MySql) Close() {
	m.db.Close()
}
